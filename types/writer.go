package types

import (
	"context"
	"io"
	"sync"
)

type CustomWriter struct {
	io.Writer
	Context      context.Context
	BytesWritten int64
	Mu           sync.Mutex
}

func (w *CustomWriter) Write(p []byte) (int, error) {
	select {
	case <-w.Context.Done():
		return 0, w.Context.Err()
	default:
		n, err := w.Writer.Write(p)
		w.Mu.Lock()
		w.BytesWritten += int64(n)
		w.Mu.Unlock()
		return n, err
	}
}
