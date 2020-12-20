package main

import (
	"log"
	"net/http"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/linha-de-comando/v3"
)

const nomeArquivoBD = "jogo.db.json"

func main() {
	armazenamento, close, err := poquer.ArmazenamentoSistemaDeArquivoJogadorAPartirDeArquivo(nomeArquivoBD)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poquer.NovoServidorJogador(armazenamento)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("nao foi possivel escutar na porta 5000 %v", err)
	}
}
