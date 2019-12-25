package main

import "fmt"

const prefixoOlaPortugues = "Olá, "

// Ola retorna uma saudação personalizada, com "Olá, Mundo" como padrão de um nome vazio for passado
func Ola(nome string) string {
	if nome == "" {
		nome = "Mundo"
	}
	return prefixoOlaPortugues + nome
}

func main() {
	fmt.Println(Ola("mundo"))
}
