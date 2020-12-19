package main

import (
	"os"
)

type fita struct {
	arquivo *os.File
}

func (t *fita) Write(p []byte) (n int, err error) {
	t.arquivo.Truncate(0)
	t.arquivo.Seek(0, 0)
	return t.arquivo.Write(p)
}
