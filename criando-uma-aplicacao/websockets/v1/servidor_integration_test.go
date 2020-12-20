package poquer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[]`)
	defer limparBaseDeDados()
	armazenamento, err := NewFileSystemPlayerStore(baseDeDados)

	verificaSemErro(t, err)

	server := mustMakePlayerServer(t, armazenamento)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get pontuação", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get Liga", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response, http.StatusOK)

		obtido := getLeagueFromResponse(t, response.Body)
		esperado := []Jogador{
			{"Pepper", 3},
		}
		verificaLiga(t, obtido, esperado)
	})
}
