package main

import (
	"fmt"
	"io"
	"os"
)

// Contagem imprime uma contagem de 3 para a sáida
func Contagem(saida io.Writer) {
	fmt.Fprint(saida, "3")
}

func main() {
	Contagem(os.Stdout)
}
