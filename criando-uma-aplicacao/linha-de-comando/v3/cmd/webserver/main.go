package main

import (
	"log"
	"net/http"

	poker "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/command-line/v3"
)

const dbFileName = "game.db.json"

func main() {
	armazenamento, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(armazenamento)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("nao foi possivel escutar na porta 5000 %v", err)
	}
}
