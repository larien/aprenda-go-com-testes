package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestaSalvarERetornarVitorias(t *testing.T) {
	bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, "")
	defer limpaBancoDeDados()
	armazenamento := &SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}
	servidor := NovoServidorDoJogador(armazenamento)
	jogador := "Pepper"

	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoPostaVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoPostaVitoria(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoPostaVitoria(jogador))

	t.Run("get score", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoPegaPontuacao(player))
		defineStatus(t, resposta.Code, http.StatusOK)

		defineCorpoDeReposta(t, resposta.Body.String(), "3")
	})

	t.Run("pega liga", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, requisicaoNovaLiga())
		defineStatus(t, resposta.Code, http.StatusOK)

		recebido := pegaLigaDaResposta(t, resposta.Body)
		esperado := []Jogador{
			{"Pepper", 3},
		}
		defineLiga(t, recebido, esperado)
	})
}
