package poquer_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/websockets/v1"
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
	t.Run("começa partida com 3 jogadores e termina partida com 'Chris' como vencedor", func(t *testing.T) {
		partida := &JogoEspiao{}
		saida := &bytes.Buffer{}

		in := usuarioEnvia("3", "Chris venceu")
		cli := poquer.NovaCLI(in, saida, partida)

		cli.JogarPoquer()

		verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador)
		verificaJogoComeçadoCom(t, partida, 3)
		verificaTerminosChamadosCom(t, partida, "Chris")
	})

	t.Run("começa partida com 8 jogadores e grava 'Cleo' como vencedor", func(t *testing.T) {
		partida := &JogoEspiao{}

		in := usuarioEnvia("8", "Cleo venceu")
		cli := poquer.NovaCLI(in, SaidaTosca, partida)

		cli.JogarPoquer()

		verificaJogoComeçadoCom(t, partida, 8)
		verificaTerminosChamadosCom(t, partida, "Cleo")
	})

	t.Run("imprime um erro quando um valor não numérico é inserido e não começa a partida", func(t *testing.T) {
		partida := &JogoEspiao{}

		saida := &bytes.Buffer{}
		in := usuarioEnvia("tortas")

		cli := poquer.NovaCLI(in, saida, partida)
		cli.JogarPoquer()

		verificaPartidaNaoIniciada(t, partida)
		verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador, poquer.ErrMsgEntradaJogadorIncorreta)
	})

	t.Run("imprime um erro quando o vencedor é declarado incorretamente", func(t *testing.T) {
		partida := &JogoEspiao{}
		saida := &bytes.Buffer{}

		in := usuarioEnvia("8", "Lloyd é incrível")
		cli := poquer.NovaCLI(in, saida, partida)

		cli.JogarPoquer()

		verificaPartidaNaoFinalizada(t, partida)
		verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador, poquer.ErrMsgEntradaVencedorIncorreta)
	})
}

func verificaJogoComeçadoCom(t *testing.T, partida *JogoEspiao, numeroDeJogadoresDesejados int) {
	t.Helper()
	if partida.ComecouASerChamadoCom != numeroDeJogadoresDesejados {
		t.Errorf("esperava Começar chamado com %d mas obteve %d", numeroDeJogadoresDesejados, partida.ComecouASerChamadoCom)
	}
}

func verificaPartidaNaoFinalizada(t *testing.T, partida *JogoEspiao) {
	t.Helper()
	if partida.TerminouDeSerChamado {
		t.Errorf("partida não deveria ter finalizado")
	}
}

func verificaPartidaNaoIniciada(t *testing.T, partida *JogoEspiao) {
	t.Helper()
	if partida.ComecouASerChamado {
		t.Errorf("partida não deveria ter começado")
	}
}

func verificaTerminosChamadosCom(t *testing.T, partida *JogoEspiao, vencedor string) {
	t.Helper()
	if partida.TerminouDeSerChamadoCom != vencedor {
		t.Errorf("esperava chamada de término com '%s' mas obteve '%s' ", vencedor, partida.TerminouDeSerChamadoCom)
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
