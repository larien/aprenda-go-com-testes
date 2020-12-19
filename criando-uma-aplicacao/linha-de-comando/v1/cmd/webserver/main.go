package main

import (
	"log"
	"net/http"
	"os"

	poker "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/linha-de-comando/v1"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("falha ao abrir %s %v", dbFileName, err)
	}

	armazenamento, err := poker.NovoArmazenamentoSistemaDeArquivodeJogador(db)

	if err != nil {
		log.Fatalf("falha ao criar sistema de arquivos para armazenar jogadores, %v ", err)
	}

	server := poker.NewPlayerServer(armazenamento)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("nao foi possivel escutar na porta 5000 %v", err)
	}
}
