package proto

import "errors"

var (
	ErrUnexpectedEOF = errors.New("proto: unexpected EOF")
)

type T interface {
	Marshal(*Writer)
	Unmarshal(*Reader)
}

type Writer struct {
	B []byte
}

type Reader struct {
	B      []byte
	Offset int
	Err    error
}

func (w *Writer) WriteInt8(i int8) {
	w.B = append(w.B, byte(i))
}

func (r *Reader) ReadInt8() int8 {
	if r.Err != nil {
		return 0
	}
	if r.Offset >= len(r.B) {
		r.Err = ErrUnexpectedEOF
		return 0
	}
	i := r.Offset
	r.Offset++
	return int8(r.B[i])
}

func (w *Writer) WriteInt16(i int16) {
	w.B = append(w.B, byte(i>>8), byte(i))
}

func (r *Reader) ReadInt16() int16 {
	if r.Err != nil {
		return 0
	}
	if r.Offset >= len(r.B) {
		r.Err = ErrUnexpectedEOF
		return 0
	}
	i := r.Offset
	r.Offset += 2
	return int16(r.B[i])<<8 | int16(r.B[i+1])
}

func (w *Writer) WriteInt32(i int32) {
	w.B = append(w.B, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}

func (r *Reader) ReadInt32() int32 {
	if r.Err != nil {
		return 0
	}
	if r.Offset >= len(r.B) {
		r.Err = ErrUnexpectedEOF
		return 0
	}
	i := r.Offset
	r.Offset += 4
	return int32(r.B[i])<<24 | int32(r.B[i+1])<<16 | int32(r.B[i+2])<<8 | int32(r.B[i+3])
}

func (w *Writer) WriteInt64(i int64) {
	w.B = append(w.B, byte(i>>56), byte(i>>48), byte(i>>40), byte(i>>32), byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}

func (r *Reader) ReadInt64() int64 {
	if r.Err != nil {
		return 0
	}
	if r.Offset >= len(r.B) {
		r.Err = ErrUnexpectedEOF
		return 0
	}
	i := r.Offset
	r.Offset += 8
	return int64(r.B[i])<<56 | int64(r.B[i+1])<<48 | int64(r.B[i+2])<<40 | int64(r.B[i+3])<<32 |
		int64(r.B[i+4])<<24 | int64(r.B[i+5])<<16 | int64(r.B[i+6])<<8 | int64(r.B[i+7])
}

func (w *Writer) WriteString(s string) {
	w.WriteInt16(int16(len(s)))
	w.B = append(w.B, s...)
}

func (r *Reader) ReadString() string {
	if r.Err != nil {
		return ""
	}
	if r.Offset >= len(r.B) {
		r.Err = ErrUnexpectedEOF
		return ""
	}
	l := int(r.ReadInt16())
	i := r.Offset
	r.Offset += l
	return string(r.B[i : i+l])
}

func (w *Writer) WriteBytes(bs []byte) {
	w.WriteInt32(int32(len(bs)))
	w.B = append(w.B, bs...)
}

func (r *Reader) ReadBytes() []byte {
	if r.Err != nil {
		return nil
	}
	if r.Offset >= len(r.B) {
		r.Err = ErrUnexpectedEOF
		return nil
	}
	l := int(r.ReadInt32())
	i := r.Offset
	r.Offset += l
	return r.B[i : i+l]
}

func (w *Writer) SetInt32(offset int, i int32) {
	w.B[offset] = byte(i >> 24)
	w.B[offset+1] = byte(i >> 16)
	w.B[offset+2] = byte(i >> 8)
	w.B[offset+3] = byte(i)
}

func (w *Writer) SetUint32(offset int, i uint32) {
	w.B[offset] = byte(i >> 24)
	w.B[offset+1] = byte(i >> 16)
	w.B[offset+2] = byte(i >> 8)
	w.B[offset+3] = byte(i)
}
