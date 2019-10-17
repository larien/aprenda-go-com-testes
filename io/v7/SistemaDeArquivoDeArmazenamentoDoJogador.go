package main

import (
	"encoding/json"
	"os"
)

// SistemaDeArquivoDeArmazenamentoDoJogador armazena jogadores no sistema de arquivos
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
	bancoDeDados *json.Encoder
	liga Liga
}

// NovoSistemaDeArquivoDeArmazenamentoDoJogador cria um SistemaDeArquivoDeArmazenamentoDoJogador
func NovoSistemaDeArquivoDeArmazenamentoDoJogador(arquivo *os.File) *SistemaDeArquivoDeArmazenamentoDoJogador {
	arquivo.Seek(0, 0)
	liga, _ := NovaLiga(arquivo)

	return &SistemaDeArquivoDeArmazenamentoDoJogador{
		bancoDeDados: json.NewEncoder(&fita{arquivo}),
		liga:   liga,
	}
}

// PegaLiga retorna a pontuacao de todos os jogadores
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() Liga {
	return f.liga
}

// PegaPontuacaoDoJogador retorna a pontuacao do jogador
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaPontuacaoDoJogador(nome string) int {

	jogador := f.liga.Find(nome)

	if jogador != nil {
		return jogador.Vitorias
	}

	return 0
}

// SalvaVitoria vai armazenar uma vitoria para o jogador, aumentando se ja for conhecido
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {
	jogador := f.liga.Find(nome)

	if jogador != nil {
		jogador.Vitorias++
	} else {
		f.liga = append(f.liga, Jogador{nome, 1})
	}

	f.bancoDeDados.Encode(f.liga)
}
