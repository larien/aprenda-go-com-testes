package poquer

import (
	"io"
	"time"
)

// TexasHoldem manages a partida of poquer
type TexasHoldem struct {
	alerter       AlertadorDeBlind
	armazenamento ArmazenamentoJogador
}

// NovoTexasHoldem retorna a new partida
func NovoTexasHoldem(alerter AlertadorDeBlind, armazenamento ArmazenamentoJogador) *TexasHoldem {
	return &TexasHoldem{
		alerter:       alerter,
		armazenamento: armazenamento,
	}
}

// Começar will schedule blind alerts dependant on the number of jogadores
func (p *TexasHoldem) Começar(numeroDeJogadores int, destinoDosAlertas io.Writer) {
	incrementoDeBlind := time.Duration(5+numeroDeJogadores) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.alerter.AgendarAlertaPara(blindTime, blind, destinoDosAlertas)
		blindTime = blindTime + incrementoDeBlind
	}
}

// Terminar ends the partida, recording the vencedor
func (p *TexasHoldem) Terminar(vencedor string) {
	p.armazenamento.GravarVitoria(vencedor)
}
