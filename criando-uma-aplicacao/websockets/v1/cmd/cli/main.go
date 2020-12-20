package main

import (
	"fmt"
	"log"
	"os"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/websockets/v1"
)

const dbFileName = "partida.db.json"

func main() {
	armazenamento, close, err := poquer.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	partida := poquer.NovoTexasHoldem(poquer.AlertadorDeBlindFunc(poquer.SaidaAlertador), armazenamento)
	cli := poquer.NovaCLI(os.Stdin, os.Stdout, partida)

	fmt.Println("Let's play poquer")
	fmt.Println("Type {Nome} venceu to record a win")
	cli.JogarPoquer()
}
