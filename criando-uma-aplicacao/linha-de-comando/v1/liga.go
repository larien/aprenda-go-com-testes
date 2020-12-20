package poquer

import (
	"encoding/json"
	"fmt"
	"io"
)

// Liga armazena uma coleção de jogadores
type Liga []Jogador

// Encontrar tenta retornar um jogador de uma liga
func (l Liga) Encontrar(nome string) *Jogador {
	for i, p := range l {
		if p.Nome == nome {
			return &l[i]
		}
	}
	return nil
}

// NovaLiga cria uma liga de um JSON
func NovaLiga(leitor io.Reader) (Liga, error) {
	var liga []Jogador
	err := json.NewDecoder(leitor).Decode(&liga)

	if err != nil {
		err = fmt.Errorf("falha ao analizar a liga, %v", err)
	}

	return liga, err
}
