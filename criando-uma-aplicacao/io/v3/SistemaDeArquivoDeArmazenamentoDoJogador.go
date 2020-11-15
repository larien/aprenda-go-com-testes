package main

import (
	"encoding/json"
	"io"
)

// SistemaDeArquivoDeArmazenamentoDoJogador armazena jogadores no sistema de arquivos
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
    bancoDeDados io.Reader
}

// PegaLiga retorna a pontuacao de todos os jogadores
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() Liga {
	f.bancoDeDados.Seek(0, 0)
    liga, _ := NovaLiga(f.bancoDeDados)
    return liga
}

// PegaPontuacaoDoJogador retorna a pontuacao de um jogador
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaPontuacaoDoJogador(nome string) int {

	jogador := f.PegaLiga().Find(nome)

	if jogador != nil {
		return jogador.Vitorias
	}

	return 0
}
// SalvaVitoria vai armazenar uma vitoria para o jogador, aumentando se ja for conhecido
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {
	liga := f.PegaLiga()
	jogador := liga.Find(nome)

	if jogador != nil {
		jogador.Vitorias++
	}

	f.bancoDeDados.Seek(0, 0)
	json.NewEncoder(f.bancoDeDados).Encode(liga)
}
