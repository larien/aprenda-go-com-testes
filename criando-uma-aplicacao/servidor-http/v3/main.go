package main

import (
	"log"
	"net/http"
)

// JogadorArmazenamentoEmMemoria armazena na memória os dados sobre os jogadores
type JogadorArmazenamentoEmMemoria struct{}

// ObterPontuacaoJogador obtém as pontuações para um jogador
func (i *JogadorArmazenamentoEmMemoria) ObterPontuacaoJogador(nome string) int {
	return 123
}

func main() {
	server := &JogadorServidor{&JogadorArmazenamentoEmMemoria{}}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("não foi possível ouvir na porta 5000 %v", err)
	}
}
