package main

// CriarJogadorArmazenamentoNaMemoria cria um JogadorArmazenamento vazio
func CriarJogadorArmazenamentoNaMemoria() *JogadorArmazenamentoNaMemoria {
	return &JogadorArmazenamentoNaMemoria{map[string]int{}}
}

// JogadorArmazenamentoEmMemoria armazena na memória os dados sobre os jogadores
type JogadorArmazenamentoNaMemoria struct {
	armazenamento map[string]int
}

// RegistrarVitoria irá registrar uma vitoria
func (ja *JogadorArmazenamentoNaMemoria) RegistrarVitoria(nome string) {
	ja.armazenamento[nome]++
}

// ObterPontuacaoJogador obtém as pontuações para um jogador
func (ja *JogadorArmazenamentoNaMemoria) ObterPontuacaoJogador(nome string) int {
	return ja.armazenamento[nome]
}
