package poquer

import (
	"testing"
)

func TestCLI(t *testing.T) {
	armazenamentoJogador := &EsbocoArmazenamentoJogador{}

	cli := &CLI{armazenamentoJogador}
	cli.JogarPoquer()

	if len(armazenamentoJogador.ChamadasDeVitoria) < 1 {
		t.Fatal("esperando uma chamada de vitoria mas nao recebi nenhuma")
	}
}
