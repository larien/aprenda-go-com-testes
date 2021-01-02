package poquer

import (
	"encoding/json"
	"fmt"
	"io"
)

// Liga armazena uma coleção de jogadores
type Liga []Jogador

// Encontrar tenta retornar um jogador de uma Liga
func (l Liga) Encontrar(nome string) *Jogador {
	for i, p := range l {
		if p.Nome == nome {
			return &l[i]
		}
	}
	return nil
}

// NovaLiga cria uma liga do JSON
func NovaLiga(rdr io.Reader) (Liga, error) {
	var liga []Jogador
	err := json.NewDecoder(rdr).Decode(&liga)

	if err != nil {
		err = fmt.Errorf("problema ao fazer parse da Liga, %v", err)
	}

	return liga, err
}
