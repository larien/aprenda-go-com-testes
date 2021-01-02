package poquer_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	poquer "github.com/larien/learn-go-with-tests/criando-uma-aplicacao/websockets/v2"
)

var AlertadorDeBlindTosco = &poquer.AlertadorDeBlindEspiao{}
var ArmazenamentoJogadorTosco = &poquer.EsbocoDeArmazenamentoJogador{}
var EntradaTosca = &bytes.Buffer{}
var SaidaTosca = &bytes.Buffer{}

type JogoEspiao struct {
	ComecouASerChamado    bool
	ComecouASerChamadoCom int
	AlertaDeBlind         []byte

	TerminouDeSerChamado    bool
	TerminouDeSerChamadoCom string
}

func (j *JogoEspiao) Começar(numeroDeJogadores int, saida io.Writer) {
	j.ComecouASerChamado = true
	j.ComecouASerChamadoCom = numeroDeJogadores
	saida.Write(j.AlertaDeBlind)
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

		poquer.NovaCLI(entrada, saida, jogo).JogarPoquer()

		verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador)
		verificaJogoComeçadoCom(t, jogo, 3)
		verificaTerminosChamadosCom(t, jogo, "Chris")
	})

	t.Run("começa jogo com 8 jogadores e grava 'Cleo' como vencedor", func(t *testing.T) {
		jogo := &JogoEspiao{}

		entrada := usuarioEnvia("8", "Cleo venceu")

		poquer.NovaCLI(entrada, SaidaTosca, jogo).JogarPoquer()

		verificaJogoComeçadoCom(t, jogo, 8)
		verificaTerminosChamadosCom(t, jogo, "Cleo")
	})

	t.Run("imprime um erro quando um valor não numérico é inserido e não começa a jogo", func(t *testing.T) {
		jogo := &JogoEspiao{}

		saida := &bytes.Buffer{}
		entrada := usuarioEnvia("tortas")

		poquer.NovaCLI(entrada, saida, jogo).JogarPoquer()

		verificaPartidaNaoIniciada(t, jogo)
		verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador, poquer.ErrMsgEntradaJogadorIncorreta)
	})

	t.Run("imprime um erro quando o vencedor é declarado incorretamente", func(t *testing.T) {
		jogo := &JogoEspiao{}

		saida := &bytes.Buffer{}
		entrada := usuarioEnvia("8", "Lloyd é incrível")

		poquer.NovaCLI(entrada, saida, jogo).JogarPoquer()

		verificaPartidaNaoFinalizada(t, jogo)
		verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador, poquer.ErrMsgEntradaVencedorIncorreta)
	})
}

func verificaJogoComeçadoCom(t *testing.T, jogo *JogoEspiao, numeroDeJogadoresDesejados int) {
	t.Helper()

	passou := tentarNovamenteAte(500*time.Millisecond, func() bool {
		return jogo.ComecouASerChamadoCom == numeroDeJogadoresDesejados
	})

	if !passou {
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

	passou := tentarNovamenteAte(500*time.Millisecond, func() bool {
		return jogo.TerminouDeSerChamadoCom == vencedor
	})

	if !passou {
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
