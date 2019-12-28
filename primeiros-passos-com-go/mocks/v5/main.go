package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Sleeper te permite definir pausas
type Sleeper interface {
	Pausa()
}

// SleeperConfiguravel é uma implementação de Sleepr com uma pausa definida
type SleeperConfiguravel struct {
	duracao time.Duration
	pausa   func(time.Duration)
}

// Pausa vai pausar a execução pela Duração definida
func (s *SleeperConfiguravel) Pausa() {
	s.pausa(s.duracao)
}

const ultimaPalavra = "Vai!"
const inicioContagem = 3

// Contagem imprime uma contagem de 3 para a saída com um atraso determinado por um Sleeper
func Contagem(saida io.Writer, sleeper Sleeper) {
	for i := inicioContagem; i > 0; i-- {
		sleeper.Pausa()
		fmt.Fprintln(saida, i)
	}

	sleeper.Pausa()
	fmt.Fprint(saida, ultimaPalavra)
}

func main() {
	sleeper := &SleeperConfiguravel{1 * time.Second, time.Sleep}
	Contagem(os.Stdout, sleeper)
}
