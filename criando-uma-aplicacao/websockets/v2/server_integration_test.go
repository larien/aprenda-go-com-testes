package poquer_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/linha-de-comando/v1"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	baseDeDados, limparBaseDeDados := criarArquivoTemporario(t, `[]`)
	defer limparBaseDeDados()
	armazenamento, err := poquer.NewFileSystemPlayerStore(baseDeDados)

	verificaSemErro(t, err)

	server := mustMakePlayerServer(t, armazenamento, dummyGame)
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
		esperado := []poquer.Jogador{
			{Nome: "Pepper", Vitorias: 3},
		}
		verificaLiga(t, obtido, esperado)
	})
}
