package main

import (
	"io"
)

// SistemaDeArquivoDeArmazenamentoDoJogador armazena jogadores no sistema de arquivos
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
	bancoDeDados io.ReadSeeker
}

// PegaLiga retorna a pontuacao de todos os jogadores
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() []Jogador {
	liga, _ := NovaLiga(f.bancoDeDados)
	return liga
}
