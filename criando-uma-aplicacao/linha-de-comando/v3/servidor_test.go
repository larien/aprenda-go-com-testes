package poquer

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	armazenamento := EsbocoArmazenamentoJogador{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := NovoServidorJogador(&armazenamento)

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
	armazenamento := EsbocoArmazenamentoJogador{
		map[string]int{},
		nil,
		nil,
	}
	server := NovoServidorJogador(&armazenamento)

	t.Run("armazena vitórias com POST", func(t *testing.T) {
		jogador := "Pepper"

		request := newPostWinRequest(jogador)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
		VerificaVitoriaJogador(t, &armazenamento, jogador)
	})
}

func TestLeague(t *testing.T) {

	t.Run("retorna a tabela da liga como JSON", func(t *testing.T) {
		wantedLeague := []Jogador{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := EsbocoArmazenamentoJogador{nil, nil, wantedLeague}
		server := NovoServidorJogador(&armazenamento)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		obtido := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		verificaLiga(t, obtido, wantedLeague)
		assertContentType(t, response, tipoConteudoJSON)

	})
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, esperado string) {
	t.Helper()
	if response.Header().Get("content-type") != esperado {
		t.Errorf("resposta não tem o content-type igual a %s, recebi %v", esperado, response.HeaderMap)
	}
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []Jogador {
	t.Helper()
	liga, err := NovaLiga(body)

	if err != nil {
		t.Fatalf("Incapaz de converter a resposta do servidor '%s' em forma de Jogador, '%v'", body, err)
	}

	return liga
}

func verificaLiga(t *testing.T, obtido, esperado []Jogador) {
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
	req, _ := http.NewRequest(http.MethodGet, "/liga", nil)
	return req
}

func newGetScoreRequest(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func newPostWinRequest(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func assertResponseBody(t *testing.T, obtido, esperado string) {
	t.Helper()
	if obtido != esperado {
		t.Errorf("corpo da resposta incorreto, recebi '%s' esperava '%s'", obtido, esperado)
	}
}
