package poquer

import "testing"

// EsbocoArmazenamentoJogador implementa ArmazenamentoJogador para propósitos de testes
type EsbocoArmazenamentoJogador struct {
	Pontuacoes        map[string]int
	ChamadasDeVitoria []string
	Liga              []Jogador
}

// ObterPontuacaoDeJogador retorna uma pontuacao de Pontuacoes
func (s *EsbocoArmazenamentoJogador) ObterPontuacaoDeJogador(nome string) int {
	pontuacao := s.Pontuacoes[nome]
	return pontuacao
}

// GravarVitoria grava uma vitória para ChamadasDeVitoria
func (s *EsbocoArmazenamentoJogador) GravarVitoria(nome string) {
	s.ChamadasDeVitoria = append(s.ChamadasDeVitoria, nome)
}

// ObterLiga retorna a Liga
func (s *EsbocoArmazenamentoJogador) ObterLiga() Liga {
	return s.Liga
}

// VerificaVitoriaJogador permite que você espione as chamadas para GravarVitoria do armazenamento
func VerificaVitoriaJogador(t *testing.T, armazenamento *EsbocoArmazenamentoJogador, vencedor string) {
	t.Helper()

	if len(armazenamento.ChamadasDeVitoria) != 1 {
		t.Fatalf("recebi %d chamadas de GravarVitoria, esperava %d", len(armazenamento.ChamadasDeVitoria), 1)
	}

	if armazenamento.ChamadasDeVitoria[0] != vencedor {
		t.Errorf("não armazenou o vencedor correto recebi '%s', esperava '%s'", armazenamento.ChamadasDeVitoria[0], vencedor)
	}
}
