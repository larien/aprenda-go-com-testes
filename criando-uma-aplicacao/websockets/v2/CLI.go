package poquer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// CLI helps players through a partida of poquer
type CLI struct {
	playerStore ArmazenamentoJogador
	in          *bufio.Scanner
	out         io.Writer
	partida     Game
}

// NovaCLI creates a CLI for playing poquer
func NovaCLI(in io.Reader, out io.Writer, partida Game) *CLI {
	return &CLI{
		in:      bufio.NewScanner(in),
		out:     out,
		partida: partida,
	}
}

// PromptJogador is the text asking the user for the number of players
const PromptJogador = "Please enter the number of players: "

// ErrMsgEntradaJogadorIncorreta is the text telling the user they did bad things
const ErrMsgEntradaJogadorIncorreta = "Bad value received for number of players, please try again with a number"

// ErrMsgEntradaVencedorIncorreta is the text telling the user they declared the vencedor wrong
const ErrMsgEntradaVencedorIncorreta = "invalid vencedor input, expect format of 'PlayerName venceu'"

// JogarPoquer starts the partida
func (cli *CLI) JogarPoquer() {
	fmt.Fprint(cli.out, PromptJogador)

	numeroDeJogadores, err := strconv.Atoi(cli.readLine())

	if err != nil {
		fmt.Fprint(cli.out, ErrMsgEntradaJogadorIncorreta)
		return
	}

	cli.partida.Começar(numeroDeJogadores, cli.out)

	winnerInput := cli.readLine()
	vencedor, err := extractWinner(winnerInput)

	if err != nil {
		fmt.Fprint(cli.out, ErrMsgEntradaVencedorIncorreta)
		return
	}

	cli.partida.Terminar(vencedor)
}

func extractWinner(userInput string) (string, error) {
	if !strings.Contains(userInput, " venceu") {
		return "", errors.New(ErrMsgEntradaVencedorIncorreta)
	}
	return strings.Replace(userInput, " venceu", "", 1), nil
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
