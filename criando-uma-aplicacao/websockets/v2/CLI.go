package poquer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// CLI ajuda jogadores em uma jogo de pôquer
type CLI struct {
	armazenamentoJogador ArmazenamentoJogador
	entrada              *bufio.Scanner
	saida                io.Writer
	jogo                 Jogo
}

// NovaCLI cria uma CLI para jogar pôquer
func NovaCLI(entrada io.Reader, saida io.Writer, jogo Jogo) *CLI {
	return &CLI{
		entrada: bufio.NewScanner(entrada),
		saida:   saida,
		jogo:    jogo,
	}
}

// PromptJogador é o texto pedindo o número de jogadores para o usuário
const PromptJogador = "Favor entrar o número de jogadores: "

// ErrMsgEntradaJogadorIncorreta representa o texto dizendo ao usuário que ele inseriu um valor incorreto
const ErrMsgEntradaJogadorIncorreta = "Valor inválido recebido para número de jogadores, favor tentar novamente com um número"

// ErrMsgEntradaVencedorIncorreta representa o texto dizendo ao usuário que a declaração de vencedor foi errada
const ErrMsgEntradaVencedorIncorreta = "entrada de vencedor incorreta, espera-se formato de 'NomeDoJogador venceu'"

// JogarPoquer começa a jogo
func (cli *CLI) JogarPoquer() {
	fmt.Fprint(cli.saida, PromptJogador)

	numeroDeJogadores, err := strconv.Atoi(cli.lerLinha())

	if err != nil {
		fmt.Fprint(cli.saida, ErrMsgEntradaJogadorIncorreta)
		return
	}

	cli.jogo.Começar(numeroDeJogadores, cli.saida)

	entradaVencedor := cli.lerLinha()
	vencedor, err := extrairJogador(entradaVencedor)

	if err != nil {
		fmt.Fprint(cli.saida, ErrMsgEntradaVencedorIncorreta)
		return
	}

	cli.jogo.Terminar(vencedor)
}

func extrairJogador(userInput string) (string, error) {
	if !strings.Contains(userInput, " venceu") {
		return "", errors.New(ErrMsgEntradaVencedorIncorreta)
	}
	return strings.Replace(userInput, " venceu", "", 1), nil
}

func (cli *CLI) lerLinha() string {
	cli.entrada.Scan()
	return cli.entrada.Text()
}
