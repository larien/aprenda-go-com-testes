package main

import (
	"log"
	"net/http"
	"os"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/linha-de-comando/v2"
)

const nomeArquivoBD = "jogo.db.json"

func main() {
	db, err := os.OpenFile(nomeArquivoBD, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problema ao abrir %s %v", nomeArquivoBD, err)
	}

	armazenamento, err := poquer.NovoArmazenamentoSistemaDeArquivodeJogador(db)

	if err != nil {
		log.Fatalf("problema criando armazenamento de sistema de arquivo de jogador, %v ", err)
	}

	server := poquer.NovoServidorJogador(armazenamento)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
