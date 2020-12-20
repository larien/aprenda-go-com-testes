package main

import (
	"log"
	"net/http"
	"os"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/websockets/v2"
)

const nomeArquivoBaseDeDados = "partida.db.json"

func main() {
	db, err := os.OpenFile(nomeArquivoBaseDeDados, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problema ao abrir %s %v", nomeArquivoBaseDeDados, err)
	}

	armazenamento, err := poquer.NovoSistemaArquivoArmazenamentoJogador(db)

	if err != nil {
		log.Fatalf("problema ao criar sistema de arquivo de armazenamento do jogador, %v ", err)
	}

	partida := poquer.NovoTexasHoldem(poquer.AlertadorDeBlindFunc(poquer.Alertador), armazenamento)

	servidor, err := poquer.NovoServidorJogador(armazenamento, partida)

	if err != nil {
		log.Fatalf("problema ao criar o servidor do jogador %v", err)
	}

	if err := http.ListenAndServe(":5000", servidor); err != nil {
		log.Fatalf("não foi possível ouvir na porta 5000 %v", err)
	}
}
