package poker

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

	server := NewPlayerServer(armazenamento)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		obtido := getLeagueFromResponse(t, response.Body)
		esperado := []Player{
			{"Pepper", 3},
		}
		verificaLiga(t, obtido, esperado)
	})
}
