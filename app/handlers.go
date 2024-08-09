package main

import (
	"bytes"
	"compress/gzip"
	"os"
	"path"
	"slices"
	"strings"
)

// echo target string as body
func echoHandler(request *Request) (*Response, error) {
	target := request.Target
	targetParts := strings.Split(target, "/")
	if len(targetParts) != 2 {
		return &Response{
			StatusCode: 400,
			Message:    "should have one path",
		}, nil
	}

	var pathValue string

	slices.Reverse(targetParts)
	pathValue = targetParts[0]
	if pathValue == "" {
		return &Response{
			StatusCode: 400,
			Message:    "path not cannot be empty",
		}, nil
	}

	encoding := getEncoding(request.Header.AcceptEncoding)
	if encoding == "gzip" {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err := zw.Write([]byte(pathValue))
		if err != nil {
			return nil, err
		}
		err = zw.Close()
		if err != nil {
			return nil, err
		}
		pathValue = buf.String()
	}

	return &Response{
		StatusCode:    200,
		Message:       "OK",
		ContentType:   "text/plain",
		ContentLength: len(pathValue),
		Body: Body{
			content: pathValue,
		},
	}, nil
}

func getEncoding(acceptEncoding string) string {
	acceptEncodingCleared := strings.ReplaceAll(acceptEncoding, " ", "")
	isMultipleAcceptEncoding := strings.Contains(acceptEncodingCleared, ",")
	if !isMultipleAcceptEncoding {
		return acceptEncodingCleared
	}

	encodings := strings.Split(acceptEncodingCleared, ",")
	for _, enc := range encodings {
		if enc == "gzip" {
			return "gzip"
		}
	}
	return ""
}

func filesHandler(request *Request) (*Response, error) {
	target := request.Target
	body := request.Body.content
	method := request.Method

	if !strings.HasPrefix(target, "/files") {
		return &Response{
			StatusCode: 400,
			Message:    "path should start with /files",
		}, nil
	}
	target = strings.TrimPrefix(target, "/files")
	target = path.Join(PWD, target)

	// Check if the file exists
	if _, err := os.Stat(target); os.IsNotExist(err) {
		if method != "POST" {
			return &Response{
				StatusCode: 404,
				Message:    "Not Found",
			}, nil
		}

		// create a new file
		_, err := os.Create(target)
		if err != nil {
			return nil, err
		}
		// write to the file
		err = os.WriteFile(target, []byte(body), 0644)
		if err != nil {
			return nil, err
		}

		return &Response{
			StatusCode: 201,
			Message:    "Created",
		}, nil
	}

	// Read the file
	file, err := os.ReadFile(target)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode:    200,
		Message:       "OK",
		ContentType:   "application/octet-stream",
		ContentLength: len(file),
		Body: Body{
			content: string(file),
		},
	}, nil
}
