package src

import (
	"io"
	"strings"
)

type Reader interface {
	Read(p []byte) (int, error)
	ReadAll(bufSize int) (string, error)
	BytesRead() int64
}

type CountingToLowerReaderImpl struct {
	Reader         io.Reader
	TotalBytesRead int64
}

func (cr *CountingToLowerReaderImpl) Read(p []byte) (int, error) {

	n, err := cr.Reader.Read(p)
	cr.TotalBytesRead += int64(n)
	// copy(p, bytes.ToLower(p))
	toLowerInplace(p)
	return n, err
}

func (cr *CountingToLowerReaderImpl) ReadAll(bufSize int) (string, error) {
	strBuilder := strings.Builder{}
	buf := make([]byte, bufSize)

	n, err := cr.Read(buf)
	for ; err == nil; n, err = cr.Read(buf) {
		strBuilder.Write(buf[:n])
	}

	if err == io.EOF {
		return strBuilder.String(), nil
	}

	return strBuilder.String(), err
}

func (cr *CountingToLowerReaderImpl) BytesRead() int64 {
	return cr.TotalBytesRead
}

func NewCountingReader(r io.Reader) *CountingToLowerReaderImpl {
	return &CountingToLowerReaderImpl{
		Reader: r,
	}
}

// only ASCII
func toLowerInplace(p []byte) {
	shift := byte('a' - 'A')
	for i, v := range p {
		if 'A' <= v && v <= 'Z' {
			p[i] = v + shift
		}
	}
}
