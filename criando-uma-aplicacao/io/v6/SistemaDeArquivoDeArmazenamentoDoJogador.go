package main

import (
	"encoding/json"
	"io"
)

// SistemaDeArquivoDeArmazenamentoDoJogador armazena jogadores no sistema de arquivos
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
	bancoDeDados io.ReaderWriteSeeker
	liga Liga
}

// NovoSistemaDeArquivoDeArmazenamentoDoJogador cria um SistemaDeArquivoDeArmazenamentoDoJogador
func NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados io.ReadWriteSeeker) *SistemaDeArquivoDeArmazenamentoDoJogador {
	bancoDeDados.Seek(0, 0)
	liga, _ := NovaLiga(bancoDeDados)

	return &SistemaDeArquivoDeArmazenamentoDoJogador{
		bancoDeDados: bancoDeDados,
		liga:   liga,
	}
}

// PegaLiga retorna a pontuacao de todos os jogadores
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() Liga {
	return f.liga
}

// PegaPontuacaoDoJogador retorna a pontuacao do jogador
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaPontuacaoDoJogador(nome string) int {

	jogador := f.PegaLiga().Find(nome)

	if jogador != nil {
		return jogador.Vitorias
	}

	return 0
}

// SalvaVitoria vai armazenar uma vitoria para o jogador, aumentando se ja for conhecido
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {
	jogador := liga.Find(nome)

	if jogador != nil {
		jogador.Vitorias++
	} else {
		f.liga = append(f.liga, Jogador{nome, 1})
	}

	f.bancoDeDados.Seek(0, 0)
	json.NewEncoder(f.bancoDeDados).Encode(f.liga)
}
