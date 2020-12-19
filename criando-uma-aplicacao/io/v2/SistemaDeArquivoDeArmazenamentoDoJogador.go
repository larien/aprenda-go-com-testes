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
	f.bancoDeDados.Seek(0, 0)
	liga, _ := NovaLiga(f.bancoDeDados)
	return liga
}

// PegaPontuacaoDoJogador coleta a pontuacao de um jogador
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaPontuacaoDoJogador(nome string) int {

	var vitorias int

	for _, jogador := range f.PegaLiga() {
		if jogador.Nome == nome {
			vitorias = jogador.Vitorias
			break
		}
	}

	return vitorias
}
