package poquer

import (
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := criarArquivoTemporario(t, "12345")
	defer clean()

	tape := &tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)

	obtido := string(newFileContents)
	esperado := "abc"

	if obtido != esperado {
		t.Errorf("obtido '%s' esperado '%s'", obtido, esperado)
	}
}
