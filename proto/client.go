package proto

import (
	"io"
)

func (req *Request) Send(conn io.Writer) error {
	var w Writer
	(&RequestOrResponse{T: req}).Marshal(&w)
	_, err := conn.Write(w.B)
	return err
}

func (resp *Response) Receive(conn io.Reader) error {
	r := Reader{B: make([]byte, 4)}
	if _, err := conn.Read(r.B); err != nil {
		return err
	}
	size := int(r.ReadInt32())
	r.Grow(size)
	if _, err := io.ReadAtLeast(conn, r.B[4:], size); err != nil {
		return err
	}
	r.Reset()
	(&RequestOrResponse{T: resp}).Unmarshal(&r)
	return r.Err
}
