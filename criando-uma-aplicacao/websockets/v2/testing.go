package poquer

import (
	"fmt"
	"io"
	"testing"
	"time"
)

// EsbocoDeArmazenamentoJogador implementa ArmazenamentoJogador para propósitos de teste
type EsbocoDeArmazenamentoJogador struct {
	Pontuações        map[string]int
	ChamadasDeVitoria []string
	Liga              []Jogador
}

// ObtemPontuacaoDoJogador retorna uma pontuação de Pontuações
func (s *EsbocoDeArmazenamentoJogador) ObtemPontuacaoDoJogador(nome string) int {
	pontuação := s.Pontuações[nome]
	return pontuação
}

// GravarVitoria grava uma vitória para ChamadasDeVitoria
func (s *EsbocoDeArmazenamentoJogador) GravarVitoria(nome string) {
	s.ChamadasDeVitoria = append(s.ChamadasDeVitoria, nome)
}

// ObterLiga retorna Liga
func (s *EsbocoDeArmazenamentoJogador) ObterLiga() Liga {
	return s.Liga
}

// VerificaVitoriaDoVencedor te permite espionar as chamadas ao armazenamento de GravarVitoria
func VerificaVitoriaDoVencedor(t *testing.T, armazenamento *EsbocoDeArmazenamentoJogador, vencedor string) {
	t.Helper()

	if len(armazenamento.ChamadasDeVitoria) != 1 {
		t.Fatalf("obtido %d chamadas paraGravarVitoria esperado %d", len(armazenamento.ChamadasDeVitoria), 1)
	}

	if armazenamento.ChamadasDeVitoria[0] != vencedor {
		t.Errorf("não armazenou o vencedor correto obtido '%s' esperado '%s'", armazenamento.ChamadasDeVitoria[0], vencedor)
	}
}

// AlertaAgendado contém informações sobre quando um alerta é agendado
type AlertaAgendado struct {
	Em      time.Duration
	Quantia int
}

func (s AlertaAgendado) String() string {
	return fmt.Sprintf("%d chips em %v", s.Quantia, s.Em)
}

// AlertadorDeBlindEspiao te permite espionar em chamadas AgendarAlertaPara
type AlertadorDeBlindEspiao struct {
	Alertas []AlertaAgendado
}

// AgendarAlertaPara grava alertas que foram agendados
func (s *AlertadorDeBlindEspiao) AgendarAlertaPara(em time.Duration, quantia int, para io.Writer) {
	s.Alertas = append(s.Alertas, AlertaAgendado{em, quantia})
}
