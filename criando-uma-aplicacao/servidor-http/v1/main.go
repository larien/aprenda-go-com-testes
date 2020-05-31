package main

import (
	"log"
	"net/http"
)

func main() {
	tratador := http.HandlerFunc(JogadorServidor)
	if err := http.ListenAndServe(":5000", tratador); err != nil {
		log.Fatalf("não foi possível ouvir na porta 5000 %v", err)
	}
}
