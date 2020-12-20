package poquer_test

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
	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/websockets/v2"
)

var (
	dummyGame = &JogoEspiao{}
	tenMS     = 10 * time.Millisecond
)

func deveFazerServidorJogador(t *testing.T, armazenamento poquer.ArmazenamentoJogador, partida poquer.Jogo) *poquer.ServidorJogador {
	servidor, err := poquer.NovoServidorJogador(armazenamento, partida)
	if err != nil {
		t.Fatal("problema ao criar o servidor do jogador", err)
	}
	return servidor
}

func TestObterJogadores(t *testing.T) {
	armazenamento := poquer.EsbocoDeArmazenamentoJogador{
		Pontuações: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	servidor := deveFazerServidorJogador(t, &armazenamento, dummyGame)

	t.Run("retorna pontuação de Pepper", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Pepper")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificaStatus(t, resposta, http.StatusOK)
		verificaCorpoDaResposta(t, resposta.Body.String(), "20")
	})

	t.Run("retorna pontuação do Floyd", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Floyd")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificaStatus(t, resposta, http.StatusOK)
		verificaCorpoDaResposta(t, resposta.Body.String(), "10")
	})

	t.Run("retorna 404 para jogadores em falta", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Apollo")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificaStatus(t, resposta, http.StatusNotFound)
	})
}

func TestArmazenarVitórias(t *testing.T) {
	armazenamento := poquer.EsbocoDeArmazenamentoJogador{
		Pontuações: map[string]int{},
	}
	servidor := deveFazerServidorJogador(t, &armazenamento, dummyGame)

	t.Run("grava vitória no POST", func(t *testing.T) {
		jogador := "Pepper"

		requisicao := novaRequisiçãoPostDeVitoria(jogador)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificaStatus(t, resposta, http.StatusAccepted)
		poquer.VerificaVitoriaDoVencedor(t, &armazenamento, jogador)
	})
}

func TestLiga(t *testing.T) {

	t.Run("retorna a tabela da Liga como JSON", func(t *testing.T) {
		ligaEsperada := []poquer.Jogador{
			{Nome: "Cleo", Vitorias: 32},
			{Nome: "Chris", Vitorias: 20},
			{Nome: "Tiest", Vitorias: 14},
		}

		armazenamento := poquer.EsbocoDeArmazenamentoJogador{Liga: ligaEsperada}
		servidor := deveFazerServidorJogador(t, &armazenamento, dummyGame)

		requisicao := novaRequisicaoDeLiga()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		obtido := obterLigaDaResposta(t, resposta.Body)

		verificaStatus(t, resposta, http.StatusOK)
		verificaLiga(t, obtido, ligaEsperada)
		verificaTipoDoConteudo(t, resposta, "application/json")

	})
}

func TestJogo(t *testing.T) {
	t.Run("GET /partida retorna 200", func(t *testing.T) {
		servidor := deveFazerServidorJogador(t, &poquer.EsbocoDeArmazenamentoJogador{}, dummyGame)

		requisicao := novaRequisicaoJogo()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificaStatus(t, resposta, http.StatusOK)
	})

	t.Run("start a partida with 3 jogadores, send some blind alerts down WS and declare Ruth the vencedor", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		vencedor := "Ruth"

		partida := &JogoEspiao{AlertaDeBlind: []byte(wantedBlindAlert)}
		servidor := httptest.NewServer(deveFazerServidorJogador(t, ArmazenamentoJogadorTosco, partida))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

		defer servidor.Close()
		defer ws.Close()

		escreverMensagemNoWebsocket(t, ws, "3")
		escreverMensagemNoWebsocket(t, ws, vencedor)

		verificaJogoComeçadoCom(t, partida, 3)
		verificaTerminosChamadosCom(t, partida, vencedor)
		within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
	})
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, esperado string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != esperado {
		t.Errorf(`obtido "%s", esperado "%s"`, string(msg), esperado)
	}
}

func tentarNovamenteAte(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed saida")
	case <-done:
	}
}

func escreverMensagemNoWebsocket(t *testing.T, conexão *websocket.Conn, mensagem string) {
	t.Helper()
	if err := conexão.WriteMessage(websocket.TextMessage, []byte(mensagem)); err != nil {
		t.Fatalf("não foi possível enviar mensagem na conexão websocket %v", err)
	}
}

func verificaTipoDoConteudo(t *testing.T, resposta *httptest.ResponseRecorder, esperado string) {
	t.Helper()
	if resposta.Header().Get("content-type") != esperado {
		t.Errorf("resposta não obteve content-type de %s, obtido %v", esperado, resposta.HeaderMap)
	}
}

func obterLigaDaResposta(t *testing.T, corpo io.Reader) []poquer.Jogador {
	t.Helper()
	liga, err := poquer.NovaLiga(corpo)

	if err != nil {
		t.Fatalf("Não foi possível fazer parse da resposta do servidor '%s' no slice de Jogador, '%v'", corpo, err)
	}

	return liga
}

func verificaLiga(t *testing.T, obtido, esperado []poquer.Jogador) {
	t.Helper()
	if !reflect.DeepEqual(obtido, esperado) {
		t.Errorf("obtido %v esperado %v", obtido, esperado)
	}
}

func verificaStatus(t *testing.T, obtido *httptest.ResponseRecorder, esperado int) {
	t.Helper()
	if obtido.Code != esperado {
		t.Errorf("não obteve o status correto, obtido %d, esperado %d", obtido.Code, esperado)
	}
}

func novaRequisicaoJogo() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/partida", nil)
	return req
}

func novaRequisicaoDeLiga() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/liga", nil)
	return req
}

func novaRequisicaoObterPontuacao(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func novaRequisiçãoPostDeVitoria(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func verificaCorpoDaResposta(t *testing.T, obtido, esperado string) {
	t.Helper()
	if obtido != esperado {
		t.Errorf("resposta corpo está incorreta, obtido '%s' esperado '%s'", obtido, esperado)
	}
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("não foi possível abrir uma conexão de websocket em %s %v", url, err)
	}

	return ws
}
