package main

import (
	"github.com/larien/learn-go-with-tests/criando-uma-aplicacao/time/v1"
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("falha ao abrir %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("falha ao criar sistema de arquivos para armazenar jogadores, %v ", err)
	}

	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("nao foi possivel escutar na porta 5000 %v", err)
	}
}
