package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/command-line/v3"
)

const dbFileName = "game.db.json"

func main() {
	armazenamento, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Vamos jogar poker")
	fmt.Println("Digite {Nome} venceu para registrar uma vitoria")
	poker.NewCLI(armazenamento, os.Stdin).PlayPoker()
}
