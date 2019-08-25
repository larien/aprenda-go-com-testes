package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

// League armazenda uma coleção de jogadores
type League []Player

// Find tenta retornar um jogador de uma liga
func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

// NewLeague cria uma liga(league) de um JSON
func NewLeague(rdr io.Reader) (League, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)

	if err != nil {
		err = fmt.Errorf("falha ao analizar a liga, %v", err)
	}

	return league, err
}
