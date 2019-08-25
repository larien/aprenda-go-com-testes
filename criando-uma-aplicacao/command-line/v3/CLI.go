package poker

import (
	"bufio"
	"io"
	"strings"
)

// CLI auxilia jogadores em um jogo de poker
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
}

// NewCLI cria uma CLI para jogar poker
func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
	}
}

// PlayPoker inicia o jogo
func (cli *CLI) PlayPoker() {
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " venceu", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
