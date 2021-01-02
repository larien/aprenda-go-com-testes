package poquer

import "time"

// TexasHoldem gerencia um jogo de pôquer
type TexasHoldem struct {
	alertador     AlertadorDeBlind
	armazenamento ArmazenamentoJogador
}

// NovoTexasHoldem retorna um novo jogo
func NovoTexasHoldem(alertador AlertadorDeBlind, armazenamento ArmazenamentoJogador) *TexasHoldem {
	return &TexasHoldem{
		alertador:     alertador,
		armazenamento: armazenamento,
	}
}

// Começar armazena alertas de blind dependendo do número de jogadores
func (p *TexasHoldem) Começar(numeroDeJogadores int) {
	incrementoDeBlind := time.Duration(5+numeroDeJogadores) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	horarioDoBlind := 0 * time.Second
	for _, blind := range blinds {
		p.alertador.AgendarAlertaPara(horarioDoBlind, blind)
		horarioDoBlind = horarioDoBlind + incrementoDeBlind
	}
}

// Terminar finaliza o jogo, gravando o vencedor
func (p *TexasHoldem) Terminar(vencedor string) {
	p.armazenamento.GravarVitoria(vencedor)
}
