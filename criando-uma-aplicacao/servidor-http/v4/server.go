package main

import (
	"fmt"
	"net/http"
)

// JogadorArmazenamento armazena as pontuacoes dos jogadores
type JogadorArmazenamento interface {
	ObterPontuacaoJogador(nome string) int
	RegistrarVitoria(nome string)
}

// JogadorServidor é uma interface HTTP para os dados dos jogadores
type JogadorServidor struct {
	armazenamento JogadorArmazenamento
}

func (js *JogadorServidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]

	switch r.Method {
	case http.MethodPost:
		js.registrarVitoria(w, jogador)
	case http.MethodGet:
		js.mostrarPontuacao(w, jogador)
	}
}

func (js *JogadorServidor) mostrarPontuacao(w http.ResponseWriter, jogador string) {
	pontuacao := js.armazenamento.ObterPontuacaoJogador(jogador)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}

func (p *JogadorServidor) registrarVitoria(w http.ResponseWriter, jogador string) {
	p.armazenamento.RegistrarVitoria(jogador)
	w.WriteHeader(http.StatusAccepted)
}
