package main

import (
	"fmt"
	"io"
	"os"
)

const ultimaPalavra = "Vai!"
const inicioContagem = 3

// Contagem imprime uma contagem de 3 para a sáida
func Contagem(saida io.Writer) {
	for i := inicioContagem; i > 0; i-- {
		fmt.Fprintln(saida, i)
	}
	fmt.Fprint(saida, ultimaPalavra)
}

func main() {
	Contagem(os.Stdout)
}
