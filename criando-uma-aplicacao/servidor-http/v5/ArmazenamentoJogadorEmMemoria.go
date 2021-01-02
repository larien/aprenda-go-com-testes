package main

// NovoArmazenamentoJogadorEmMemoria cria um ArmazenamentoJogador vazio
func NovoArmazenamentoJogadorEmMemoria() *ArmazenamentoJogadorEmMemoria {
	return &ArmazenamentoJogadorEmMemoria{map[string]int{}}
}

// ArmazenamentoJogadorEmMemoria armazena na memória os dados sobre os jogadores
type ArmazenamentoJogadorEmMemoria struct {
	armazenamento map[string]int
}

// RegistrarVitoria irá registrar uma vitoria
func (a *ArmazenamentoJogadorEmMemoria) RegistrarVitoria(nome string) {
	a.armazenamento[nome]++
}

// ObterPontuacaoJogador obtém as pontuações para um jogador
func (a *ArmazenamentoJogadorEmMemoria) ObterPontuacaoJogador(nome string) int {
	return a.armazenamento[nome]
}
