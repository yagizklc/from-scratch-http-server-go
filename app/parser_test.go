package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestParseRequest(t *testing.T) {
	testRequests := []struct {
		name     string
		url      string
		expected Request
	}{
		{
			"simple get",
			"GET / HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			Request{
				Method: "GET",
				Target: "/",
				Header: Header{
					Host:      "localhost:4221",
					UserAgent: "curl/7.64.1",
					Accept:    "*/*",
				},
			},
		},
		{
			"target: extract path",
			"GET /echo/abc HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
			Request{
				Method: "GET",
				Target: "/echo/abc",
				Header: Header{
					Host:      "localhost:4221",
					UserAgent: "curl/7.64.1",
					Accept:    "*/*",
				},
			},
		},
		{
			"header: extract user-agent",
			"GET /user-agent HTTP/1.1 \r\nHost: localhost:4221\r\nUser-Agent: foobar/1.2.3\r\nAccept: */*\r\n\r\n",
			Request{
				Method: "GET",
				Target: "/user-agent",
				Header: Header{
					Host:      "localhost:4221",
					UserAgent: "foobar/1.2.3",
					Accept:    "*/*",
				},
			},
		},
		{
			"header	: extract Accept-Encoding",
			"GET /echo/foo HTTP/1.1 \r\nHost: localhost:4221\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nAccept-Encoding: gzip\r\n\r\n",
			Request{
				Method: "GET",
				Target: "/echo/foo",
				Header: Header{
					Host:           "localhost:4221",
					UserAgent:      "curl/7.81.0",
					Accept:         "*/*",
					AcceptEncoding: "gzip",
				},
			},
		},
		{
			"header	: extract Accept-Encoding multiple",
			"GET /echo/foo HTTP/1.1 \r\nHost: localhost:4221\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\nAccept-Encoding: gzip, somethingelse\r\n\r\n",
			Request{
				Method: "GET",
				Target: "/echo/foo",
				Header: Header{
					Host:           "localhost:4221",
					UserAgent:      "curl/7.81.0",
					Accept:         "*/*",
					AcceptEncoding: "gzip, somethingelse",
				},
			},
		},
		{
			"body: extract body",
			"POST /files/foo HTTP/1.1 \r\nHost: localhost:4221\r\nUser-Agent: foobar/1.2.3\r\nAccept: */*\r\nContent-Type: application/octet-stream\r\nContent-Length: 11\r\n\r\nfoo_content",
			Request{
				Method: "POST",
				Target: "/files/foo",
				Header: Header{
					Host:          "localhost:4221",
					UserAgent:     "foobar/1.2.3",
					Accept:        "*/*",
					ContentType:   "application/octet-stream",
					ContentLength: "11",
				},
				Body: Body{
					content: "foo_content",
				},
			},
		},
	}

	for _, tr := range testRequests {
		var buf bytes.Buffer
		buf.WriteString(tr.url)
		reader := bufio.NewReader(&buf)
		request, err := parseRequest(reader)
		if err != nil {
			t.Errorf("Failed to parse request: %v", err)
		}

		if *request != tr.expected {
			t.Errorf("Expected %s, got %s", tr.expected, *request)
		}
	}
}
