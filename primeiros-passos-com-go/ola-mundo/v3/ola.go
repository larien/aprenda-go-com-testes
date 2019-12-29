package main

import "fmt"

// Ola retorna uma saudação personalizada
func Ola(nome string) string {
	return "Olá, " + nome
}

func main() {
	fmt.Println(Ola("mundo"))
}
