package poquer_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	poquer "github.com/larien/aprenda-go-com-testes/criando-uma-aplicacao/websockets/v2"
)

func TestGravaVitoriasEAsRetorna(t *testing.T) {
	baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[]`)
	defer limparBaseDeDados()
	armazenamento, err := poquer.NovoSistemaArquivoArmazenamentoJogador(baseDeDados)

	verificaSemErro(t, err)

	servidor := deveFazerServidorJogador(t, armazenamento, jogoTosco)
	jogador := "Pepper"

	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))

	t.Run("obterpontuação", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
		verificaStatus(t, resposta, http.StatusOK)

		verificaCorpoDaResposta(t, resposta.Body.String(), "3")
	})

	t.Run("obterLiga", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoDeLiga())
		verificaStatus(t, resposta, http.StatusOK)

		obtido := obterLigaDaResposta(t, resposta.Body)
		esperado := []poquer.Jogador{
			{Nome: "Pepper", Vitorias: 3},
		}
		verificaLiga(t, obtido, esperado)
	})
}
