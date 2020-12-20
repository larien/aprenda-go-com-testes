package poquer_test

import (
	"io/ioutil"
	"testing"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/websockets/v2"
)

func TestTape_Write(t *testing.T) {
	arquivo, limpar := criarArquivoTemporario(t, "12345")
	defer limpar()

	tape := &poquer.Tape{File: arquivo}

	tape.Write([]byte("abc"))

	arquivo.Seek(0, 0)
	novosConteudosDeArquivo, _ := ioutil.ReadAll(arquivo)

	obtido := string(novosConteudosDeArquivo)
	esperado := "abc"

	if obtido != esperado {
		t.Errorf("obtido '%s' esperado '%s'", obtido, esperado)
	}
}
