package common

import (
	"io"
)

// WrappedWriter 实现行包装功能
type WrappedWriter struct {
	W    io.Writer
	Wrap int
	Cnt  int
}

func (w *WrappedWriter) Write(p []byte) (n int, err error) {
	var total int
	for i := 0; i < len(p); i++ {
		if w.Cnt > 0 && w.Cnt%w.Wrap == 0 {
			if _, err := w.W.Write([]byte{'\n'}); err != nil {
				return total, err
			}
			w.Cnt = 0
		}
		if _, err := w.W.Write(p[i : i+1]); err != nil {
			return total, err
		}
		w.Cnt++
		total++
	}
	return total, nil
}
