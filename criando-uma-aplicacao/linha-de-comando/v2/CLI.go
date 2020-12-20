package poquer

// CLI auxilia jogadores em um jogo de poquer
type CLI struct {
	armazenamentoJogador ArmazenamentoJogador
}

// JogarPoquer inicia o jogo
func (cli *CLI) JogarPoquer() {
	cli.armazenamentoJogador.GravarVitoria("Cleo")
}
