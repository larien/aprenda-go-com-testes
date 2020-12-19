package poker

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) ObterPontuacaoDeJogador(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) ObterLiga() League {
	return s.league
}

func TestGETPlayers(t *testing.T) {
	armazenamento := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&armazenamento)

	t.Run("retorna os pontos da Pepper", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("retorna os pontos do Floyd", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("retorna 404 para jogadores que não existem", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	armazenamento := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&armazenamento)

	t.Run("armazenda vitórias com POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(armazenamento.winCalls) != 1 {
			t.Fatalf("recebi %d chamadas de RecordWin esperava %d", len(armazenamento.winCalls), 1)
		}

		if armazenamento.winCalls[0] != player {
			t.Errorf("não armazenou o vencedor correto recebi '%s' esperava '%s'", armazenamento.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {

	t.Run("retorna a tabela da liga como JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&armazenamento)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		obtido := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		verificaLiga(t, obtido, wantedLeague)
		assertContentType(t, response, jsonContentType)

	})
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, esperado string) {
	t.Helper()
	if response.Header().Get("content-type") != esperado {
		t.Errorf("resposta não tem o content-type igual a %s, recebi %v", esperado, response.HeaderMap)
	}
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []Player {
	t.Helper()
	league, err := NewLeague(body)

	if err != nil {
		t.Fatalf("Incapaz de converter a resposta do servidor '%s' em forma de Player, '%v'", body, err)
	}

	return league
}

func verificaLiga(t *testing.T, obtido, esperado []Player) {
	t.Helper()
	if !reflect.DeepEqual(obtido, esperado) {
		t.Errorf("recebi %v esperava %v", obtido, esperado)
	}
}

func assertStatus(t *testing.T, obtido, esperado int) {
	t.Helper()
	if obtido != esperado {
		t.Errorf("não pegou o estado correto, recebi %d, esperava %d", obtido, esperado)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t *testing.T, obtido, esperado string) {
	t.Helper()
	if obtido != esperado {
		t.Errorf("corpo da resposta incorreto, recebi '%s' esperava '%s'", obtido, esperado)
	}
}
