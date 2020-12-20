package poquer

import (
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	arquivo, clean := criarArquivoTemporario(t, "12345")
	defer clean()

	tape := &tape{arquivo}

	tape.Write([]byte("abc"))

	arquivo.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(arquivo)

	obtido := string(newFileContents)
	esperado := "abc"

	if obtido != esperado {
		t.Errorf("recebi '%s' esperava '%s'", obtido, esperado)
	}
}
