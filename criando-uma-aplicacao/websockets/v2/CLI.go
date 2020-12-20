package poquer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// CLI ajuda jogadores em uma partida de pôquer
type CLI struct {
	armazenamentoJogador ArmazenamentoJogador
	entrada              *bufio.Scanner
	saida                io.Writer
	partida              Jogo
}

// NovaCLI cria uma CLI para jogar pôquer
func NovaCLI(entrada io.Reader, saida io.Writer, partida Jogo) *CLI {
	return &CLI{
		entrada: bufio.NewScanner(entrada),
		saida:   saida,
		partida: partida,
	}
}

// PromptJogador é o texto pedindo o número de jogadores para o usuário
const PromptJogador = "Favor entrar o número de jogadores: "

// ErrMsgEntradaJogadorIncorreta is the text telling the user they did bad things
const ErrMsgEntradaJogadorIncorreta = "Bad value received for number of jogadores, please try again with a number"

// ErrMsgEntradaVencedorIncorreta is the text telling the user they declared the vencedor wrong
const ErrMsgEntradaVencedorIncorreta = "invalid vencedor input, expect format of 'PlayerName venceu'"

// JogarPoquer starts the partida
func (cli *CLI) JogarPoquer() {
	fmt.Fprint(cli.saida, PromptJogador)

	numeroDeJogadores, err := strconv.Atoi(cli.readLine())

	if err != nil {
		fmt.Fprint(cli.saida, ErrMsgEntradaJogadorIncorreta)
		return
	}

	cli.partida.Começar(numeroDeJogadores, cli.saida)

	winnerInput := cli.readLine()
	vencedor, err := extractWinner(winnerInput)

	if err != nil {
		fmt.Fprint(cli.saida, ErrMsgEntradaVencedorIncorreta)
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
	cli.entrada.Scan()
	return cli.entrada.Text()
}
