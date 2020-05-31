package main

import (
	"fmt"
	"net/http"
)

// JogadorArmazenamento armazena as pontuacoes dos jogadores
type JogadorArmazenamento interface {
	ObterPontuacaoJogador(nome string) int
}

// JogadorServidor é uma interface HTTP para os dados dos jogadores
type JogadorServidor struct {
	armazenamento JogadorArmazenamento
}

func (js *JogadorServidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		js.registrarVitoria(w)
	case http.MethodGet:
		js.mostrarPontuacao(w, r)
	}

}

func (p *JogadorServidor) mostrarPontuacao(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]

	pontuacao := p.armazenamento.ObterPontuacaoJogador(jogador)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}

func (p *JogadorServidor) registrarVitoria(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}
