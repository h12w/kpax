package proto

import (
	"io"

	"h12.me/wipro"
)

func (req *Request) Send(conn io.Writer) error {
	var w wipro.Writer
	(&RequestOrResponse{M: req}).Marshal(&w)
	if _, err := conn.Write(w.B); err != nil {
		return ErrConn
	}
	return nil
}

func (resp *Response) Receive(conn io.Reader) error {
	r := wipro.Reader{B: make([]byte, 4)}
	if _, err := conn.Read(r.B); err != nil {
		return ErrConn
	}
	size := int(r.ReadInt32())
	r.Grow(size)
	if _, err := io.ReadAtLeast(conn, r.B[4:], size); err != nil {
		return ErrConn
	}
	r.Reset()
	(&RequestOrResponse{M: resp}).Unmarshal(&r)
	return r.Err
}
