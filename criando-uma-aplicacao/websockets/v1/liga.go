package poquer

import (
	"encoding/json"
	"fmt"
	"io"
)

// Liga stores a collection of players
type Liga []Jogador

// Find tries to return a player from a Liga
func (l Liga) Find(nome string) *Jogador {
	for i, p := range l {
		if p.Nome == nome {
			return &l[i]
		}
	}
	return nil
}

// NewLeague creates a Liga from JSON
func NewLeague(rdr io.Reader) (Liga, error) {
	var league []Jogador
	err := json.NewDecoder(rdr).Decode(&league)

	if err != nil {
		err = fmt.Errorf("problem parsing Liga, %v", err)
	}

	return league, err
}
