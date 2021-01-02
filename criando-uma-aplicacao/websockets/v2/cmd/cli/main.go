package main

import (
	"fmt"
	"log"
	"os"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/websockets/v2"
)

const nomeArquivoBaseDeDados = "jogo.db.json"

func main() {
	armazenamento, close, err := poquer.SistemaArquivoArmazenamentoJogadorDoArquivo(nomeArquivoBaseDeDados)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	jogo := poquer.NovoTexasHoldem(poquer.AlertadorDeBlindFunc(poquer.Alertador), armazenamento)
	cli := poquer.NovaCLI(os.Stdin, os.Stdout, jogo)

	fmt.Println("Vamos jogar pôquer")
	fmt.Println("Digite o nome para gravar uma vitória")
	cli.JogarPoquer()
}
