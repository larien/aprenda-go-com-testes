package main

import "errors"

// Dicionario armazena definições de palavras
type Dicionario map[string]string

// ErrNaoEncontrado significa que a definição não pôde ser encontrada para determinada palavra
var ErrNaoEncontrado = errors.New("não foi possível encontrar a palavra que você procura")

// Busca encontra uma palavra no dicionário
func (d Dicionario) Busca(palavra string) (string, error) {
	definicao, existe := d[palavra]
	if !existe {
		return "", ErrNaoEncontrado
	}

	return definicao, nil
}

// Adiciona insere uma palavra e definição no dicionário
func (d Dicionario) Adiciona(palavra, definicao string) {
	d[palavra] = definicao
}
