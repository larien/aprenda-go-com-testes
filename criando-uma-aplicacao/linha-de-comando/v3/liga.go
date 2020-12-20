package poquer

import (
	"encoding/json"
	"fmt"
	"io"
)

// Liga stores a collection of jogadores
type Liga []Jogador

// Encontrar tries to return a jogador from a liga
func (l Liga) Encontrar(nome string) *Jogador {
	for i, p := range l {
		if p.Nome == nome {
			return &l[i]
		}
	}
	return nil
}

// NovaLiga creates a liga from JSON
func NovaLiga(leitor io.Reader) (Liga, error) {
	var liga []Jogador
	err := json.NewDecoder(leitor).Decode(&liga)

	if err != nil {
		err = fmt.Errorf("problem parsing liga, %v", err)
	}

	return liga, err
}
