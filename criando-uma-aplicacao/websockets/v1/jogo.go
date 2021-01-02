package poquer

// Jogo gerencia o estado de uma jogo
type Jogo interface {
	Começar(numeroDeJogadores int)
	Terminar(vencedor string)
}
