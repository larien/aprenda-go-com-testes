package poquer

import (
	"os"
)

type tape struct {
	arquivo *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.arquivo.Truncate(0)
	t.arquivo.Seek(0, 0)
	return t.arquivo.Write(p)
}
