package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseRequest(reader *bufio.Reader) (*Request, error) {
	var buf bytes.Buffer
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading connection: ", err.Error())
			os.Exit(1)
		}
		if line == "\r\n" {
			break // End of headers
		}
		buf.WriteString(line)
	}

	request, err := _parseHeader(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse header: %v", err)
	}

	if request.Header.ContentLength != "" {
		contentLength, err := strconv.Atoi(request.Header.ContentLength)
		if err != nil {
			return nil, fmt.Errorf("failed to parse content length: %v", err)
		}

		body := make([]byte, contentLength)
		_, err = reader.Read(body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %v", err)
		}

		request.Body = Body{
			content: string(body),
		}
	}

	return request, nil
}

func _parseHeader(buf *bytes.Buffer) (*Request, error) {
	lines := strings.Split(buf.String(), "\r\n")
	requestLineParts := strings.Split(lines[0], " ")
	method := requestLineParts[0]
	target := requestLineParts[1]
	httpVersion := requestLineParts[2]
	if httpVersion != "HTTP/1.1" {
		return nil, fmt.Errorf("unsupported HTTP version: %s", httpVersion)
	}

	// Parse headers
	header := Header{}
	for _, line := range lines[1:] {
		if line == "" {
			break // End of headers
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue // Invalid header
		}
		key := parts[0]
		value := parts[1]

		switch key {
		case "Host":
			header.Host = value
		case "User-Agent":
			header.UserAgent = value
		case "Accept":
			header.Accept = value
		case "Content-Type":
			header.ContentType = value
		case "Content-Length":
			header.ContentLength = value
		case "Accept-Encoding":
			header.AcceptEncoding = value
		}
	}

	return &Request{
		Method: method,
		Target: target,
		Header: header,
	}, nil
}
