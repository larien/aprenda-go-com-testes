package main

import "errors"

// Dicionario armazena definições de palavras
type Dicionario map[string]string

// ErrNaoEncontrado é a definição para não ter encontrado determinada palavra
var ErrNaoEncontrado = errors.New("não foi possível encontrar a palavra procurada")

// Busca encontra uma palavra no dicionário
func (d Dicionario) Busca(palavra string) (string, error) {
	definicao, ok := d[palavra]
	if !ok {
		return "", ErrNaoEncontrado
	}

	return definicao, nil
}
