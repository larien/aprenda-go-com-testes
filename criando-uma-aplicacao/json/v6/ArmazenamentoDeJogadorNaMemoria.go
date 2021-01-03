package main

// NovoArmazenamentoDeJogadorNaMemoria inicializa um armazenamento de jogador vazio
func NovoArmazenamentoDeJogadorNaMemoria() *ArmazenamentoDeJogadorNaMemoria {
	return &ArmazenamentoDeJogadorNaMemoria{map[string]int{}}
}

// ArmazenamentoDeJogadorNaMemoria coleta dados sobre jogadores em memória
type ArmazenamentoDeJogadorNaMemoria struct {
	armazenamento map[string]int
}

// ObterLiga retorna uma coleção de Jogadores
func (a *ArmazenamentoDeJogadorNaMemoria) ObterLiga() []Jogador {
	var liga []Jogador
	for nome, vitórias := range a.armazenamento {
		liga = append(liga, Jogador{nome, vitórias})
	}
	return liga
}

// GravarVitoria grava a vitória de um jogador
func (a *ArmazenamentoDeJogadorNaMemoria) GravarVitoria(nome string) {
	a.armazenamento[nome]++
}

// ObtemPontuacaoDoJogador retorna pontuações para determinado jogador
func (a *ArmazenamentoDeJogadorNaMemoria) ObtemPontuacaoDoJogador(nome string) int {
	return a.armazenamento[nome]
}
