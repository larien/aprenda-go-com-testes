package poker

// CLI auxilia jogadores em um jogo de poker
type CLI struct {
	playerStore PlayerStore
}

// PlayPoker inicia o jogo
func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Cleo")
}
