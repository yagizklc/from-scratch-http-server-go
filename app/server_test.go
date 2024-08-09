package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestConcurrentConnections(t *testing.T) {
	// Start the server in a goroutine
	go func() {
		err := runServer()
		if err != nil {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Number of concurrent connections
	noc := 50

	// Use a WaitGroup to coordinate goroutines
	var wg sync.WaitGroup
	wg.Add(noc)

	// Channel to collect errors
	errorChan := make(chan error, noc)

	// Create a context with a 10-second timeout for the entire test
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Establish all connections
	for i := 0; i < noc; i++ {
		go func(index int) {
			defer wg.Done()

			var d net.Dialer
			conn, err := d.DialContext(ctx, "tcp", "localhost:4221")
			if err != nil {
				errorChan <- fmt.Errorf("connection %d failed: %v", index, err)
				return
			}
			defer conn.Close()

			// Set a 2-second timeout for individual operations
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				errorChan <- fmt.Errorf("failed to set deadline for connection %d: %v", index, err)
				return
			}

			_, err = conn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n"))
			if err != nil {
				errorChan <- fmt.Errorf("failed to write to connection %d: %v", index, err)
				return
			}

			buffer := make([]byte, 1024)
			_, err = conn.Read(buffer)
			if err != nil {
				errorChan <- fmt.Errorf("failed to read from connection %d: %v", index, err)
				return
			}

			if !strings.Contains(string(buffer), "HTTP/1.1 200 OK") {
				errorChan <- fmt.Errorf("unexpected response from connection %d: %s", index, string(buffer))
				return
			}
		}(i)
	}

	// Wait for all goroutines to finish or context to timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		t.Fatalf("Test timed out after 10 seconds")
	case <-done:
		// All connections completed successfully
	}

	close(errorChan)

	// Check for any errors
	for err := range errorChan {
		t.Error(err)
	}
}
