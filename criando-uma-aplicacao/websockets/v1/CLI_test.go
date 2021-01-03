package poquer_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poquer "github.com/larien/aprenda-go-com-testes/criando-uma-aplicacao/websockets/v1"
)

var AlertadorDeBlindTosco = &poquer.AlertadorDeBlindEspiao{}
var ArmazenamentoJogadorTosco = &poquer.EsbocoDeArmazenamentoJogador{}
var EntradaTosca = &bytes.Buffer{}
var SaidaTosca = &bytes.Buffer{}

type JogoEspiao struct {
	ComecouASerChamado    bool
	ComecouASerChamadoCom int

	TerminouDeSerChamado    bool
	TerminouDeSerChamadoCom string
}

func (j *JogoEspiao) Começar(numeroDeJogadores int) {
	j.ComecouASerChamado = true
	j.ComecouASerChamadoCom = numeroDeJogadores
}

func (j *JogoEspiao) Terminar(vencedor string) {
	j.TerminouDeSerChamado = true
	j.TerminouDeSerChamadoCom = vencedor
}

func usuarioEnvia(mensagens ...string) io.Reader {
	return strings.NewReader(strings.Join(mensagens, "\n"))
}

func TestCLI(t *testing.T) {
	t.Run("começa jogo com 3 jogadores e termina jogo com 'Chris' como vencedor", func(t *testing.T) {
		jogo := &JogoEspiao{}
		saida := &bytes.Buffer{}

		entrada := usuarioEnvia("3", "Chris venceu")
		cli := poquer.NovaCLI(entrada, saida, jogo)

		cli.JogarPoquer()

		verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador)
		verificaJogoComeçadoCom(t, jogo, 3)
		verificaTerminosChamadosCom(t, jogo, "Chris")
	})

	t.Run("começa jogo com 8 jogadores e grava 'Cleo' como vencedor", func(t *testing.T) {
		jogo := &JogoEspiao{}

		entrada := usuarioEnvia("8", "Cleo venceu")
		cli := poquer.NovaCLI(entrada, SaidaTosca, jogo)

		cli.JogarPoquer()

		verificaJogoComeçadoCom(t, jogo, 8)
		verificaTerminosChamadosCom(t, jogo, "Cleo")
	})

	t.Run("imprime um erro quando um valor não numérico é inserido e não começa a jogo", func(t *testing.T) {
		jogo := &JogoEspiao{}

		saida := &bytes.Buffer{}
		entrada := usuarioEnvia("tortas")

		cli := poquer.NovaCLI(entrada, saida, jogo)
		cli.JogarPoquer()

		verificaPartidaNaoIniciada(t, jogo)
		verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador, poquer.ErrMsgEntradaJogadorIncorreta)
	})

	t.Run("imprime um erro quando o vencedor é declarado incorretamente", func(t *testing.T) {
		jogo := &JogoEspiao{}
		saida := &bytes.Buffer{}

		entrada := usuarioEnvia("8", "Lloyd é incrível")
		cli := poquer.NovaCLI(entrada, saida, jogo)

		cli.JogarPoquer()

		verificaPartidaNaoFinalizada(t, jogo)
		verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador, poquer.ErrMsgEntradaVencedorIncorreta)
	})
}

func verificaJogoComeçadoCom(t *testing.T, jogo *JogoEspiao, numeroDeJogadoresDesejados int) {
	t.Helper()
	if jogo.ComecouASerChamadoCom != numeroDeJogadoresDesejados {
		t.Errorf("esperava Começar chamado com %d mas obteve %d", numeroDeJogadoresDesejados, jogo.ComecouASerChamadoCom)
	}
}

func verificaPartidaNaoFinalizada(t *testing.T, jogo *JogoEspiao) {
	t.Helper()
	if jogo.TerminouDeSerChamado {
		t.Errorf("jogo não deveria ter finalizado")
	}
}

func verificaPartidaNaoIniciada(t *testing.T, jogo *JogoEspiao) {
	t.Helper()
	if jogo.ComecouASerChamado {
		t.Errorf("jogo não deveria ter começado")
	}
}

func verificaTerminosChamadosCom(t *testing.T, jogo *JogoEspiao, vencedor string) {
	t.Helper()
	if jogo.TerminouDeSerChamadoCom != vencedor {
		t.Errorf("esperava chamada de término com '%s' mas obteve '%s' ", vencedor, jogo.TerminouDeSerChamadoCom)
	}
}

func verificaMensagensEnviadasParaUsuario(t *testing.T, saida *bytes.Buffer, mensagens ...string) {
	t.Helper()
	esperado := strings.Join(mensagens, "")
	obtido := saida.String()
	if obtido != esperado {
		t.Errorf("obtido '%s' enviado para saida mas esperava %+v", obtido, mensagens)
	}
}

func verificaAlertaAgendado(t *testing.T, obtido, esperado poquer.AlertaAgendado) {
	t.Helper()
	if obtido != esperado {
		t.Errorf("obtido %+v, esperado %+v", obtido, esperado)
	}
}
