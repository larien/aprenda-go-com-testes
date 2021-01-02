package main

import (
	"log"
	"net/http"
)

// ArmazenamentoJogadorEmMemoria armazena na memória os dados sobre os jogadores
type ArmazenamentoJogadorEmMemoria struct{}

// ObterPontuacaoJogador obtém as pontuações para um jogador
func (a *ArmazenamentoJogadorEmMemoria) ObterPontuacaoJogador(nome string) int {
	return 123
}

func main() {
	server := &ServidorJogador{&ArmazenamentoJogadorEmMemoria{}}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("não foi possível ouvir na porta 5000 %v", err)
	}
}
