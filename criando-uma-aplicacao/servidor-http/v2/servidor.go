package main

import (
	"fmt"
	"net/http"
)

// ArmazenamentoJogador armazena as pontuacoes dos jogadores
type ArmazenamentoJogador interface {
	ObterPontuacaoJogador(nome string) int
}

// ServidorJogador é uma interface HTTP para os dados dos jogadores
type ServidorJogador struct {
	armazenamento ArmazenamentoJogador
}

func (s *ServidorJogador) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]

	pontuacao := s.armazenamento.ObterPontuacaoJogador(jogador)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}
