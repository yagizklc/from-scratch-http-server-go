package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	PORT = "4221"
	HOST = "localhost"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	if err := runServer(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runServer() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", HOST, PORT))
	if err != nil {
		return fmt.Errorf("failed to bind to port 4221: %v", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %v", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	request, err := parseRequest(reader)
	if err != nil {
		response := handleErr(err)
		conn.Write([]byte(response.String()))
		return
	}

	endpoint := strings.Split(request.Target, "/")[1]
	response := &Response{}

	switch endpoint {
	case "":
		response = &Response{
			StatusCode:    200,
			Message:       "OK",
			ContentType:   "text/plain",
			ContentLength: 2,
		}

	case "echo":
		response, err = echoHandler(request)
		if err != nil {
			response = handleErr(err)
		}

	case "files":
		response, err = filesHandler(request)
		if err != nil {
			response = handleErr(err)
		}

	default:
		response = &Response{
			StatusCode: 404,
			Message:    "Not Found",
		}
	}

	conn.Write([]byte(response.String()))
}

func handleErr(err error) *Response {
	fmt.Println("error: ", err)
	return &Response{
		StatusCode: 500,
		Message:    "Internal Server Error",
	}
}
