package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// Liga armazena uma colecao de jogadores
type Liga []Jogador

// Busca tenta retornar um jogador de uma liga
func (l Liga) Busca(nome string) *Jogador {
	for i, p := range l {
		if p.Nome == nome {
			return &l[i]
		}
	}
	return nil
}

// NovaLiga cria uma liga a partir do JSON
func NovaLiga(rdr io.Reader) (Liga, error) {
    var liga []Jogador
	err := json.NewDecoder(rdr).Decode(&liga)

	if err != nil {
		err = fmt.Errorf("problema parseando liga, %v", err)
	}

	return liga, err
}
