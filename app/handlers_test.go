package main

import (
	"bytes"
	"compress/gzip"
	"os"
	"testing"
)

func TestEcho(t *testing.T) {
	foo := "foo"
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write([]byte(foo))
	if err != nil {
		t.Errorf("Failed to write to gzip writer: %v", err)
	}
	err = zw.Close()
	if err != nil {
		t.Errorf("Failed to close gzip writer: %v", err)
	}
	fooCompressed := buf.String()

	testRequests := []struct {
		name     string
		request  Request
		expected Response
	}{
		{
			"path invalid empty",
			Request{
				Method: "GET",
				Target: "echo/",
				Header: Header{
					Host:      "localhost:4221",
					UserAgent: "curl/7.64.1",
					Accept:    "*/*",
				},
			},
			Response{
				StatusCode: 400,
				Message:    "path not cannot be empty",
			},
		},
		{
			"path invalid more than one",
			Request{
				Method: "GET",
				Target: "echo/foo/bar",
				Header: Header{
					Host:      "localhost:4221",
					UserAgent: "curl/7.64.1",
					Accept:    "*/*",
				},
			},
			Response{
				StatusCode: 400,
				Message:    "should have one path",
			},
		},
		{
			"path valid",
			Request{
				Method: "GET",
				Target: "echo/foo",
				Header: Header{
					Host:      "localhost:4221",
					UserAgent: "curl/7.64.1",
					Accept:    "*/*",
				},
			},
			Response{
				StatusCode:    200,
				Message:       "OK",
				ContentType:   "text/plain",
				ContentLength: 3,
				Body: Body{
					content: "foo",
				},
			},
		},
		{
			"compression valid",
			Request{
				Method: "GET",
				Target: "echo/foo",
				Header: Header{
					Host:           "localhost:4221",
					UserAgent:      "curl/7.64.1",
					Accept:         "*/*",
					AcceptEncoding: "gzip",
				},
			},
			Response{
				StatusCode:      200,
				Message:         "OK",
				ContentType:     "text/plain",
				ContentLength:   len(fooCompressed),
				ContentEncoding: "gzip",
				Body: Body{
					content: fooCompressed,
				},
			},
		},
		{
			"compression invalid encoding: return uncompressed",
			Request{
				Method: "GET",
				Target: "echo/foo",
				Header: Header{
					Host:           "localhost:4221",
					UserAgent:      "curl/7.64.1",
					Accept:         "*/*",
					AcceptEncoding: "invalid-encoding",
				},
			},
			Response{
				StatusCode:    200,
				Message:       "OK",
				ContentType:   "text/plain",
				ContentLength: 3,
				Body: Body{
					content: "foo",
				},
			},
		},
		{
			"compression multiple valid: select first valid",
			Request{
				Method: "GET",
				Target: "echo/foo",
				Header: Header{
					Host:           "localhost:4221",
					UserAgent:      "curl/7.64.1",
					Accept:         "*/*",
					AcceptEncoding: "invalid-encoding-1, gzip, invalid-encoding-2",
				},
			},
			Response{
				StatusCode:      200,
				Message:         "OK",
				ContentType:     "text/plain",
				ContentLength:   len(fooCompressed),
				ContentEncoding: "gzip",
				Body: Body{
					content: fooCompressed,
				},
			},
		},
		{
			"compression multiple invalid",
			Request{
				Method: "GET",
				Target: "echo/foo",
				Header: Header{
					Host:           "localhost:4221",
					UserAgent:      "curl/7.64.1",
					Accept:         "*/*",
					AcceptEncoding: "invalid-encoding-1, invalid-encoding-2",
				},
			},
			Response{
				StatusCode:    200,
				Message:       "OK",
				ContentType:   "text/plain",
				ContentLength: 3,
				Body: Body{
					content: "foo",
				},
			},
		},
	}

	for _, tr := range testRequests {
		t.Run(tr.name, func(t *testing.T) {
			response, err := echoHandler(&tr.request)
			if err != nil {
				t.Errorf("Error handling request: %v", err)
			}
			if response.StatusCode != tr.expected.StatusCode {
				t.Errorf("Expected status code %d, got %d", tr.expected.StatusCode, response.StatusCode)
			}
			if response.Message != tr.expected.Message {
				t.Errorf("Expected message %s, got %s", tr.expected.Message, response.Message)
			}
			if response.ContentType != tr.expected.ContentType {
				t.Errorf("Expected content type %s, got %s", tr.expected.ContentType, response.ContentType)
			}
			if response.ContentLength != tr.expected.ContentLength {
				t.Errorf("Expected content length %d, got %d", tr.expected.ContentLength, response.ContentLength)
			}
			if response.Body.content != tr.expected.Body.content {
				t.Errorf("Expected body content %s, got %s", tr.expected.Body.content, response.Body.content)
			}
		})
	}
}

func TestFiles(t *testing.T) {
	testRequests := []struct {
		name     string
		request  Request
		expected Response
	}{
		{
			"existing file",
			Request{
				Method: "GET",
				Target: "/files/foo",
				Header: Header{
					Host:      "localhost:4221",
					UserAgent: "curl/7.64.1",
					Accept:    "*/*",
				},
			},
			Response{
				StatusCode:    200,
				Message:       "OK",
				ContentType:   "application/octet-stream",
				ContentLength: 11,
				Body: Body{
					content: "foo_content",
				},
			},
		},
		{
			"non existing file",
			Request{
				Method: "GET",
				Target: "/files/bar",
				Header: Header{
					Host:      "localhost:4221",
					UserAgent: "curl/7.64.1",
					Accept:    "*/*",
				},
			},
			Response{
				StatusCode: 404,
				Message:    "Not Found",
			},
		},
		{
			"create new file",
			Request{
				Method: "POST",
				Target: "/files/baz",
				Header: Header{
					Host:          "localhost:4221",
					UserAgent:     "curl/7.64.1",
					Accept:        "*/*",
					ContentType:   "application/octet-stream",
					ContentLength: "11",
				},
				Body: Body{
					content: "baz_content",
				},
			},
			Response{
				StatusCode: 201,
				Message:    "Created",
			},
		},
	}

	for _, tr := range testRequests {
		t.Run(tr.name, func(t *testing.T) {
			response, err := filesHandler(&tr.request)
			if err != nil {
				t.Errorf("Error handling request: %v", err)
			}
			if response.StatusCode != tr.expected.StatusCode {
				t.Errorf("Expected status code %d, got %d", tr.expected.StatusCode, response.StatusCode)
			}
			if response.Message != tr.expected.Message {
				t.Errorf("Expected message %s, got %s", tr.expected.Message, response.Message)
			}
			if response.ContentType != tr.expected.ContentType {
				t.Errorf("Expected content type %s, got %s", tr.expected.ContentType, response.ContentType)
			}
			if response.ContentLength != tr.expected.ContentLength {
				t.Errorf("Expected content length %d, got %d", tr.expected.ContentLength, response.ContentLength)
			}
			if response.Body.content != tr.expected.Body.content {
				t.Errorf("Expected body content %s, got %s", tr.expected.Body.content, response.Body.content)
			}
		})
	}

	// Test that the file was created
	createdFileDir := PWD + "/baz"
	if _, err := os.Stat(createdFileDir); err != nil {
		t.Errorf("failed to create file: %v", err)
	} else {
		// Clean up the created file
		err := os.Remove(createdFileDir)
		if err != nil {
			t.Errorf("failed to remove file: %v", err)
		}
	}
}
