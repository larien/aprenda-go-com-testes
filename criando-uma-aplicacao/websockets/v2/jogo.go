package poquer

import "io"

// Jogo gerencia o estado de uma jogo
type Jogo interface {
	Começar(numeroDeJogadores int, destinoDosAlertas io.Writer)
	Terminar(vencedor string)
}
