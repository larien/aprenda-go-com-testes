package poquer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[]`)
	defer limparBaseDeDados()
	armazenamento, err := NovoArmazenamentoSistemaDeArquivodeJogador(baseDeDados)

	verificaSemErro(t, err)

	server := NovoServidorJogador(armazenamento)
	jogador := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(jogador))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(jogador))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(jogador))

	t.Run("retorna os pontos", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(jogador))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("retorna a liga", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		obtido := getLeagueFromResponse(t, response.Body)
		esperado := []Jogador{
			{"Pepper", 3},
		}
		verificaLiga(t, obtido, esperado)
	})
}
