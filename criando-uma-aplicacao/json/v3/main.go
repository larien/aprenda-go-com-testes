package main

import (
	"log"
	"net/http"
)

func main() {
	servidor := NovoServidorJogador(NovoArmazenamentoDeJogadorNaMemoria())

	if err := http.ListenAndServe(":5000", servidor); err != nil {
		log.Fatalf("não foi possível ouvir na porta 5000 %v", err)
	}
}
