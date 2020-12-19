package main

import (
	"log"
	"net/http"
)

func main() {
	servidor := NovoServidorDoJogador(NovoArmazenamentoDeJogadorNaMemoria())

	if err := http.ListenAndServe(":5000", servidor); err != nil {
		log.Fatalf("Não foi possivel ouvir na porta 5000 %v", err)
	}
}
