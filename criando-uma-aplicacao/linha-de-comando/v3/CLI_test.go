package poquer_test

import (
	"io"
	"strings"
	"testing"

	poquer "github.com/larien/aprenda-go-com-testes/criando-uma-aplicacao/linha-de-comando/v3"
)

func TestCLI(t *testing.T) {

	t.Run("recorda vencedor chris digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Chris venceu\n")
		armazenamentoJogador := &poquer.EsbocoArmazenamentoJogador{}

		cli := poquer.NovoCLI(armazenamentoJogador, in)
		cli.JogarPoquer()

		poquer.VerificaVitoriaJogador(t, armazenamentoJogador, "Chris")
	})

	t.Run("recorda vencedor cleo digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Cleo venceu\n")
		armazenamentoJogador := &poquer.EsbocoArmazenamentoJogador{}

		cli := poquer.NovoCLI(armazenamentoJogador, in)
		cli.JogarPoquer()

		poquer.VerificaVitoriaJogador(t, armazenamentoJogador, "Cleo")
	})

	t.Run("não ler além da primeira nova linha", func(t *testing.T) {
		in := LeitorDeFalhaAoTerminar{
			t,
			strings.NewReader("Chris venceu\n E ai"),
		}

		armazenamentoJogador := &poquer.EsbocoArmazenamentoJogador{}

		cli := poquer.NovoCLI(armazenamentoJogador, in)
		cli.JogarPoquer()
	})

}

type LeitorDeFalhaAoTerminar struct {
	t      *testing.T
	leitor io.Reader
}

func (m LeitorDeFalhaAoTerminar) Read(p []byte) (n int, err error) {

	n, err = m.leitor.Read(p)

	if n == 0 || err == io.EOF {
		m.t.Fatal("Leu o até o fim quando não precisava")
	}

	return n, err
}
