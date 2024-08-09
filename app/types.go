package main

import "fmt"

type Header struct {
	Host           string
	UserAgent      string
	Accept         string
	ContentLength  string
	ContentType    string
	AcceptEncoding string
}

type Body struct {
	content string
}

type Request struct {
	Method string
	Target string
	Header Header
	Body   Body
}

type Response struct {
	StatusCode      int
	Message         string
	ContentType     string
	ContentLength   int
	ContentEncoding string

	Body Body
}

func (r *Response) String() string {
	if r.ContentLength == 0 {
		return fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", r.StatusCode, r.Message)
	}

	return fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", r.StatusCode, r.Message, r.ContentType, r.ContentLength, r.Body.content)
}
