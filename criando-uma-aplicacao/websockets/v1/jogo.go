package poquer

// Game manages the state of a partida
type Game interface {
	Começar(numeroDeJogadores int)
	Terminar(vencedor string)
}
