package main

import (
	"log"
	"net/http"
	"os"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/linha-de-comando/v1"
)

const nomeArquivoBD = "jogo.db.json"

func main() {
	db, err := os.OpenFile(nomeArquivoBD, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("falha ao abrir %s %v", nomeArquivoBD, err)
	}

	armazenamento, err := poquer.NovoArmazenamentoSistemaDeArquivodeJogador(db)

	if err != nil {
		log.Fatalf("falha ao criar sistema de arquivos para armazenar jogadores, %v ", err)
	}

	server := poquer.NovoServidorJogador(armazenamento)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("nao foi possivel escutar na porta 5000 %v", err)
	}
}
