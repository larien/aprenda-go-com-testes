package main

import (
	"fmt"
	"net/http"
)

// JogadorArmazenamento armazena as pontuacoes dos jogadores
type JogadorArmazenamento interface {
	ObterPontuacaoJogador(name string) int
	RegistrarVitoria(name string)
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

func (js *JogadorServidor) mostrarPontuacao(w http.ResponseWriter, player string) {
	pontuacao := js.armazenamento.ObterPontuacaoJogador(player)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}

func (p *JogadorServidor) registrarVitoria(w http.ResponseWriter, player string) {
	p.armazenamento.RegistrarVitoria(player)
	w.WriteHeader(http.StatusAccepted)
}
