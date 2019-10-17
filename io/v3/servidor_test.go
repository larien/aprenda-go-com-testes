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
	pontuacoes   map[string]int
	chamadasDeVitoria []string
	liga   []Jogador
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

func TestaPEGAJogadores(t *testing.T) {
	armazenamento := EsbocoDoArmazenamentoDoJogador{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	servidor := NovoServidorDoJogador(&armazenamento)

	t.Run("returna a pontuacao de Pepper'", func(t *testing.T) {
		requisicao := novaRequisicaoPegaPontuacao("Pepper")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		defineStatus(t, resposta.Code, http.StatusOK)
		defineCorpodeResposta(t, resposta.Body.String(), "20")
	})

	t.Run("returna a pontuacao de Floyd", func(t *testing.T) {
		requisicao := novaRequisicaoPegaPontuacao("Floyd")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		defineStatus(t, resposta.Code, http.StatusOK)
		defineCorpoDeResposta(t, resposta.Body.String(), "10")
	})

	t.Run("returna 404 para jogadores nao existentes", func(t *testing.T) {
		requisicao := novaRequisicaoPegaPontuacao("Apollo")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		defineStatus(t, resposta.Code, http.StatusNotFound)
	})
}

func TesteArmazenamentoDeVitorias(t *testing.T) {
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
			t.Fatalf("recebido %d chamadas para RecordeDeVitoria esperado %d", len(armazenamento.chamadasDeVitoria), 1)
		}

		if armazenamento.chamadasDeVitoria[0] != jogador {
			t.Errorf("nao armazenou vencedor correto recebido '%s' esperado '%s'", armazenamento.chamadasDeVitoria[0], jogador)
		}
	})
}

func defineTipoDeConteudo(t *testing.T, resposta *httptest.ResponseRecorder, esperado string) {
	t.Helper()
	if resposta.Header().Get("content-type") != esperado {
		t.Errorf("resposta nao tinha content-type de %s, recebido %v", esperado, resposta.HeaderMap)
	}
}

func pegaLigaDaResposta(t *testing.T, corpo io.Reader) []Jogador {
	t.Helper()
	liga, err := NovaLiga(corpo)

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

func requisicaoNovaLiga() *http.requisicao {
	req, _ := http.NewRequest(http.MethodGet, "/liga", nil)
	return req
}

func novaRequisicaoPegaPontuacao(nome string) *http.requisicao {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func novaRequisicaoPostaVitoria(nome string) *http.requisicao {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func defineCorpodeResposta(t *testing.T, recebido, esperado string) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("corpo da resposta esta errado, recebido '%s' esperado '%s'", recebido, esperado)
	}
}
