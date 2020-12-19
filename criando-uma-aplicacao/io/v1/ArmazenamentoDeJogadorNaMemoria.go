package main

// NovoArmazenamentoDeJogadorNaMemoria  inicializa um armazenamento vazio de jogador
func NovoArmazenamentoDeJogadorNaMemoria() *ArmazenamentoDeJogadorNaMemoria {
	return &ArmazenamentoDeJogadorNaMemoria{map[string]int{}}
}

// ArmazenamentoDeJogadorNaMemoria coleta dados sobre os jogadores na memoria
type ArmazenamentoDeJogadorNaMemoria struct {
	armazenamento map[string]int
}

// PegaLiga retorna uma colecao de jogadores
func (i *ArmazenamentoDeJogadorNaMemoria) PegaLiga() []Jogador {
	var liga []Jogador
	for nome, vitorias := range i.armazenamento {
		liga = append(liga, Jogador{nome, vitorias})
	}
	return liga
}

// SalvaVitoria armazena uma vitoria do jogador
func (i *ArmazenamentoDeJogadorNaMemoria) SalvaVitoria(nome string) {
	i.armazenamento[nome]++
}

// PegaPontuacaoDoJogador retorna pontuacoes de um dado jogador
func (i *ArmazenamentoDeJogadorNaMemoria) PegaPontuacaoDoJogador(nome string) int {
	return i.armazenamento[nome]
}
