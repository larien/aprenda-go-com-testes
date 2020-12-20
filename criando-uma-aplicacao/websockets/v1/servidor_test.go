package poquer

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func mustMakePlayerServer(t *testing.T, armazenamento ArmazenamentoJogador) *PlayerServer {
	server, err := NewPlayerServer(armazenamento)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func TestGETPlayers(t *testing.T) {
	armazenamento := EsbocoDeArmazenamentoJogador{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, &armazenamento)

	t.Run("retorna Pepper's pontuação", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("retorna Floyd's pontuação", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("retorna 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	armazenamento := EsbocoDeArmazenamentoJogador{
		map[string]int{},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, &armazenamento)

	t.Run("it records venceu on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusAccepted)
		VerificaVitoriaDoVencedor(t, &armazenamento, player)
	})
}

func TestLeague(t *testing.T) {

	t.Run("it retorna the Liga table as JSON", func(t *testing.T) {
		wantedLeague := []Jogador{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := EsbocoDeArmazenamentoJogador{nil, nil, wantedLeague}
		server := mustMakePlayerServer(t, &armazenamento)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		obtido := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response, http.StatusOK)
		verificaLiga(t, obtido, wantedLeague)
		assertContentType(t, response, jsonContentType)

	})
}

func TestGame(t *testing.T) {
	t.Run("GET /partida retorna 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &EsbocoDeArmazenamentoJogador{})

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
	})

	t.Run("when we get a message over a websocket it is a vencedor of a partida", func(t *testing.T) {
		armazenamento := &EsbocoDeArmazenamentoJogador{}
		vencedor := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, armazenamento))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
		}
		defer ws.Close()

		writeWSMessage(t, ws, vencedor)

		time.Sleep(10 * time.Millisecond)
		VerificaVitoriaDoVencedor(t, armazenamento, vencedor)
	})
}

func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, esperado string) {
	t.Helper()
	if response.Header().Get("content-type") != esperado {
		t.Errorf("response did not have content-type of %s, obtido %v", esperado, response.HeaderMap)
	}
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []Jogador {
	t.Helper()
	league, err := NewLeague(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server '%s' into slice of Jogador, '%v'", body, err)
	}

	return league
}

func verificaLiga(t *testing.T, obtido, esperado []Jogador) {
	t.Helper()
	if !reflect.DeepEqual(obtido, esperado) {
		t.Errorf("obtido %v esperado %v", obtido, esperado)
	}
}

func assertStatus(t *testing.T, obtido *httptest.ResponseRecorder, esperado int) {
	t.Helper()
	if obtido.Code != esperado {
		t.Errorf("did not get correct status, obtido %d, esperado %d", obtido.Code, esperado)
	}
}

func newGameRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/partida", nil)
	return req
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func newGetScoreRequest(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", nome), nil)
	return req
}

func newPostWinRequest(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", nome), nil)
	return req
}

func assertResponseBody(t *testing.T, obtido, esperado string) {
	t.Helper()
	if obtido != esperado {
		t.Errorf("response body is wrong, obtido '%s' esperado '%s'", obtido, esperado)
	}
}
