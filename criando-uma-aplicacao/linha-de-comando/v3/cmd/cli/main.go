package main

import (
	"fmt"
	"log"
	"os"

	poquer "github.com/larien/aprenda-go-com-testes/criando-uma-aplicacao/linha-de-comando/v3"
)

const nomeArquivoBD = "jogo.db.json"

func main() {
	armazenamento, close, err := poquer.ArmazenamentoSistemaDeArquivoJogadorAPartirDeArquivo(nomeArquivoBD)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Vamos jogar poquer")
	fmt.Println("Digite {Nome} venceu para registrar uma vitoria")
	poquer.NovoCLI(armazenamento, os.Stdin).JogarPoquer()
}
