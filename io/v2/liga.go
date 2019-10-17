package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// NovaLiga cria uma liga a partir do JSON
func NovaLiga(rdr io.Reader) ([]Jogador, error) {
    var liga []Jogador
	err := json.NewDecoder(rdr).Decode(&liga)

	if err != nil {
		err = fmt.Errorf("problema parseando liga, %v", err)
	}

	return liga, err
}
