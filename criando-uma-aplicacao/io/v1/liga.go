package main

import (
	"encoding/json"
	"io"
)

// NovaLiga cria uma liga a partir do JSON
func NovaLiga(rdr io.Reader) ([]Jogador, error) {
	var liga []Jogador
	err := json.NewDecoder(rdr).Decode(&liga)
	return liga, err
}
