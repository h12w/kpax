package proto

import (
	"io"

	"h12.me/wipro"
)

func (req *Request) Send(conn io.Writer) error {
	return wipro.Send(&RequestOrResponse{M: req}, conn)
}

func (resp *Response) Receive(conn io.Reader) error {
	return wipro.Receive(conn, &RequestOrResponse{M: resp})
}
