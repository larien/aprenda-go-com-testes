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

func deveFazerServidorJogador(t *testing.T, armazenamento ArmazenamentoJogador) *ServidorJogador {
	servidor, err := NovoServidorJogador(armazenamento)
	if err != nil {
		t.Fatal("problema ao criar o servidor do jogador", err)
	}
	return servidor
}

func TestObterJogadores(t *testing.T) {
	armazenamento := EsbocoDeArmazenamentoJogador{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	servidor := deveFazerServidorJogador(t, &armazenamento)

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
	armazenamento := EsbocoDeArmazenamentoJogador{
		map[string]int{},
		nil,
		nil,
	}
	servidor := deveFazerServidorJogador(t, &armazenamento)

	t.Run("grava vitória no POST", func(t *testing.T) {
		jogador := "Pepper"

		requisicao := novaRequisiçãoPostDeVitoria(jogador)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificaStatus(t, resposta, http.StatusAccepted)
		VerificaVitoriaDoVencedor(t, &armazenamento, jogador)
	})
}

func TestLiga(t *testing.T) {

	t.Run("retorna a tabela da Liga como JSON", func(t *testing.T) {
		ligaEsperada := []Jogador{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := EsbocoDeArmazenamentoJogador{nil, nil, ligaEsperada}
		servidor := deveFazerServidorJogador(t, &armazenamento)

		requisicao := novaRequisicaoDeLiga()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		obtido := obterLigaDaResposta(t, resposta.Body)

		verificaStatus(t, resposta, http.StatusOK)
		verificaLiga(t, obtido, ligaEsperada)
		verificaTipoDoConteudo(t, resposta, tipoConteudoJSON)

	})
}

func TestJogo(t *testing.T) {
	t.Run("GET /jogo retorna 200", func(t *testing.T) {
		servidor := deveFazerServidorJogador(t, &EsbocoDeArmazenamentoJogador{})

		requisicao := novaRequisicaoJogo()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificaStatus(t, resposta, http.StatusOK)
	})

	t.Run("quando recebemos uma mensagem de um websocket que é vencedor da jogo", func(t *testing.T) {
		armazenamento := &EsbocoDeArmazenamentoJogador{}
		vencedor := "Ruth"
		servidor := httptest.NewServer(deveFazerServidorJogador(t, armazenamento))
		defer servidor.Close()

		wsURL := "ws" + strings.TrimPrefix(servidor.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("não foi possível abrir uma conexão de websocket em %s %v", wsURL, err)
		}
		defer ws.Close()

		escreverMensagemNoWebsocket(t, ws, vencedor)

		time.Sleep(10 * time.Millisecond)
		VerificaVitoriaDoVencedor(t, armazenamento, vencedor)
	})
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

func obterLigaDaResposta(t *testing.T, corpo io.Reader) []Jogador {
	t.Helper()
	liga, err := NovaLiga(corpo)

	if err != nil {
		t.Fatalf("Não foi possível fazer parse da resposta do servidor '%s' no slice de Jogador, '%v'", corpo, err)
	}

	return liga
}

func verificaLiga(t *testing.T, obtido, esperado []Jogador) {
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
	req, _ := http.NewRequest(http.MethodGet, "/jogo", nil)
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
