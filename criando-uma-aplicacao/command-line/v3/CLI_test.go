package poker_test

import (
	"github.com/larien/learn-go-with-tests/criando-uma-aplicacao/command-line/v3"
	"io"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {

	t.Run("recorda vencedor chris digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Chris venceu\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("recorda vencedor cleo digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Cleo venceu\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

	t.Run("não ler além da primeira nova linha", func(t *testing.T) {
		in := failOnEndReader{
			t,
			strings.NewReader("Chris wins\n E ai"),
		}

		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()
	})

}

type failOnEndReader struct {
	t   *testing.T
	rdr io.Reader
}

func (m failOnEndReader) Read(p []byte) (n int, err error) {

	n, err = m.rdr.Read(p)

	if n == 0 || err == io.EOF {
		m.t.Fatal("Ler o até o fim quando não precisava")
	}

	return n, err
}
