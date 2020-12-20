package poquer

import (
	"bufio"
	"io"
	"strings"
)

// CLI auxilia jogadores em um jogo de poquer
type CLI struct {
	armazenamentoJogador ArmazenamentoJogador
	in                   *bufio.Scanner
}

// NovoCLI cria uma CLI para jogar poquer
func NovoCLI(armazenamento ArmazenamentoJogador, in io.Reader) *CLI {
	return &CLI{
		armazenamentoJogador: armazenamento,
		in:                   bufio.NewScanner(in),
	}
}

// JogarPoquer inicia o jogo
func (cli *CLI) JogarPoquer() {
	userInput := cli.readLine()
	cli.armazenamentoJogador.GravarVitoria(extrairVencedor(userInput))
}

func extrairVencedor(userInput string) string {
	return strings.Replace(userInput, " venceu", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
