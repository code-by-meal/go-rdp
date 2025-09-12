package serverdata

import "io"

type Response struct{}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) Read(stream io.Writer) error {
	return nil
}
