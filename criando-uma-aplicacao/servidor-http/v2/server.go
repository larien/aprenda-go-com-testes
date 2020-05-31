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
	jogador := r.URL.Path[len("/jogadores/"):]

	pontuacao := js.armazenamento.ObterPontuacaoJogador(jogador)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}
