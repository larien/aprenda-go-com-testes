package poquer_test

import (
	"fmt"
	"testing"
	"time"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/websockets/v1"
)

func TestJogo_Começar(t *testing.T) {
	t.Run("agenda alertas em jogos que começam com 5 jogadores", func(t *testing.T) {
		alertadorDeBlind := &poquer.AlertadorDeBlindEspiao{}
		jogo := poquer.NovoTexasHoldem(alertadorDeBlind, ArmazenamentoJogadorTosco)

		jogo.Começar(5)

		cases := []poquer.AlertaAgendado{
			{Em: 0 * time.Second, Quantia: 100},
			{Em: 10 * time.Minute, Quantia: 200},
			{Em: 20 * time.Minute, Quantia: 300},
			{Em: 30 * time.Minute, Quantia: 400},
			{Em: 40 * time.Minute, Quantia: 500},
			{Em: 50 * time.Minute, Quantia: 600},
			{Em: 60 * time.Minute, Quantia: 800},
			{Em: 70 * time.Minute, Quantia: 1000},
			{Em: 80 * time.Minute, Quantia: 2000},
			{Em: 90 * time.Minute, Quantia: 4000},
			{Em: 100 * time.Minute, Quantia: 8000},
		}

		verificaCasosAgendados(cases, t, alertadorDeBlind)
	})

	t.Run("agenda alertas em jogos que começam com 7 jogadores", func(t *testing.T) {
		alertadorDeBlind := &poquer.AlertadorDeBlindEspiao{}
		jogo := poquer.NovoTexasHoldem(alertadorDeBlind, ArmazenamentoJogadorTosco)

		jogo.Começar(7)

		cases := []poquer.AlertaAgendado{
			{Em: 0 * time.Second, Quantia: 100},
			{Em: 12 * time.Minute, Quantia: 200},
			{Em: 24 * time.Minute, Quantia: 300},
			{Em: 36 * time.Minute, Quantia: 400},
		}

		verificaCasosAgendados(cases, t, alertadorDeBlind)
	})

}

func TestJogo_Terminar(t *testing.T) {
	armazenamento := &poquer.EsbocoDeArmazenamentoJogador{}
	jogo := poquer.NovoTexasHoldem(AlertadorDeBlindTosco, armazenamento)
	vencedor := "Ruth"

	jogo.Terminar(vencedor)
	poquer.VerificaVitoriaDoVencedor(t, armazenamento, vencedor)
}

func verificaCasosAgendados(cases []poquer.AlertaAgendado, t *testing.T, alertadorDeBlind *poquer.AlertadorDeBlindEspiao) {
	for i, esperado := range cases {
		t.Run(fmt.Sprint(esperado), func(t *testing.T) {

			if len(alertadorDeBlind.Alertas) <= i {
				t.Fatalf("alerta %d não foi agendado %v", i, alertadorDeBlind.Alertas)
			}

			obtido := alertadorDeBlind.Alertas[i]
			verificaAlertaAgendado(t, obtido, esperado)
		})
	}
}
