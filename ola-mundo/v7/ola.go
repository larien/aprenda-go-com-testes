package main

import "fmt"

const espanhol = "espanhol"
const frances = "francês"
const prefixoOlaPortugues = "Olá, "
const prefixoOlaEspanhol = "Hola, "
const prefixoOlaFrances = "Bonjour, "

// Ola retorna uma saudação personalizada em determinado idiota
func Ola(nome string, idioma string) string {
	if nome == "" {
		nome = "Mundo"
	}

	prefixo := prefixoOlaPortugues

	switch idioma {
	case frances:
		prefixo = prefixoOlaFrances
	case espanhol:
		prefixo = prefixoOlaEspanhol
	}

	return prefixo + nome
}

func main() {
	fmt.Println(Ola("mundo", ""))
}
