package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type EsbocoDoArmazenamentoDoJogador struct {
	pontuacoes        map[string]int
	chamadasDeVitoria []string
	liga              []Jogador
}

func (s *EsbocoDoArmazenamentoDoJogador) PegaPontuacaoDoJogador(nome string) int {
	pontuacao := s.pontuacoes[nome]
	return pontuacao
}

func (s *EsbocoDoArmazenamentoDoJogador) SalvaVitoria(nome string) {
	s.chamadasDeVitoria = append(s.chamadasDeVitoria, nome)
}

func (s *EsbocoDoArmazenamentoDoJogador) PegaLiga() Liga {
	return s.liga
}

func TestPegaJogadores(t *testing.T) {
	armazenamento := EsbocoDoArmazenamentoDoJogador{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	servidor := NovoServidorDoJogador(&armazenamento)

	t.Run("retorna a pontuacao de Pepper'", func(t *testing.T) {
		requisicao := novaRequisicaoPegaPontuacao("Pepper")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		defineStatus(t, resposta.Code, http.StatusOK)
		definecorpodeResposta(t, resposta.Body.String(), "20")
	})

	t.Run("retorna a pontuacao de Floyd", func(t *testing.T) {
		requisicao := novaRequisicaoPegaPontuacao("Floyd")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		defineStatus(t, resposta.Code, http.StatusOK)
		definecorpodeResposta(t, resposta.Body.String(), "10")
	})

	t.Run("retorna 404 para jogadores nao existentes", func(t *testing.T) {
		requisicao := novaRequisicaoPegaPontuacao("Apollo")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		defineStatus(t, resposta.Code, http.StatusNotFound)
	})
}

func TestArmazenamentoDeVitorias(t *testing.T) {
	armazenamento := EsbocoDoArmazenamentoDoJogador{
		map[string]int{},
		nil,
		nil,
	}
	servidor := NovoServidorDoJogador(&armazenamento)

	t.Run("salva vitorias no POST", func(t *testing.T) {
		jogador := "Pepper"

		requisicao := novaRequisicaoPostaVitoria(jogador)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		defineStatus(t, resposta.Code, http.StatusAccepted)

		if len(armazenamento.chamadasDeVitoria) != 1 {
			t.Fatalf("recebido %d chamadas para salvaVitoria esperado %d", len(armazenamento.chamadasDeVitoria), 1)
		}

		if armazenamento.chamadasDeVitoria[0] != jogador {
			t.Errorf("nao armazenou vencedor correto recebido '%s' esperado '%s'", armazenamento.chamadasDeVitoria[0], jogador)
		}
	})
}

func TestLiga(t *testing.T) {

	t.Run("retorna tabela de liga como JSON", func(t *testing.T) {
		ligaDesejada := []Jogador{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := EsbocoDoArmazenamentoDoJogador{nil, nil, ligaDesejada}
		server := NovoServidorDoJogador(&armazenamento)

		requisicao := requisicaoNovaLiga()
		resposta := httptest.NewRecorder()

		server.ServeHTTP(resposta, requisicao)

		recebido := pegaLigaDaResposta(t, resposta.Body)

		defineStatus(t, resposta.Code, http.StatusOK)
		defineLiga(t, recebido, ligaDesejada)
		defineTipoDeConteudo(t, resposta, jsonContentType)

	})
}

func defineTipoDeConteudo(t *testing.T, resposta *httptest.ResponseRecorder, esperado string) {
	t.Helper()
	if resposta.Header().Get("content-type") != esperado {
		t.Errorf("resposta nao tinha content-type de %s, recebido %v", esperado, resposta.HeaderMap)
	}
}

func pegaLigaDaResposta(t *testing.T, body io.Reader) []Jogador {
	t.Helper()
	liga, err := NovaLiga(body)

	if err != nil {
		t.Fatalf("Nao foi possivel passar resposta do servidor '%s' em pedaco de jogador, '%v'", body, err)
	}

	return liga
}

func defineLiga(t *testing.T, recebido, esperado []Jogador) {
	t.Helper()
	if !reflect.DeepEqual(recebido, esperado) {
		t.Errorf("recebido %v esperado %v", recebido, esperado)
	}
}

func defineStatus(t *testing.T, recebido, esperado int) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("nao recebeu status correto, recebido %d, esperado %d", recebido, esperado)
	}
}

func requisicaoNovaLiga() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/liga", nil)
	return req
}

func novaRequisicaoPegaPontuacao(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func novaRequisicaoPostaVitoria(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func definecorpodeResposta(t *testing.T, recebido, esperado string) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("corpo da resposta esta errado, recebido '%s' esperado '%s'", recebido, esperado)
	}
}
