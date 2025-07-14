package common

import (
	"io"
)

// MockReader 自定义的模拟读取器
type MockReader struct {
	Data []byte
	Pos  int
}

// Read 实现 io.Reader 接口
func (mr *MockReader) Read(p []byte) (n int, err error) {
	if mr.Pos >= len(mr.Data) {
		return 0, io.EOF
	}
	n = copy(p, mr.Data[mr.Pos:])
	mr.Pos += n
	return n, nil
}
