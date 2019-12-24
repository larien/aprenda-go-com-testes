package main

import "fmt"

const prefixoOlaPortugues = "Olá, "

// Ola retorna uma saudação personalizada
func Ola(nome string) string {
	return prefixoOlaPortugues + nome
}

func main() {
	fmt.Println(Ola("mundo"))
}
