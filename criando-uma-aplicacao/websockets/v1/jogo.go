package poquer

// Jogo gerencia o estado de uma partida
type Jogo interface {
	Começar(numeroDeJogadores int)
	Terminar(vencedor string)
}