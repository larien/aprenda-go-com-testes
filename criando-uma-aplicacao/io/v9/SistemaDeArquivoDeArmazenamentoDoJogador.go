package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// SistemaDeArquivoDeArmazenamentoDoJogador armazena jogadores no sistema de arquivos
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
	bancoDeDados *json.Encoder
	liga         Liga
}

// NovoSistemaDeArquivoDeArmazenamentoDoJogador cria um SistemaDeArquivoDeArmazenamentoDoJogador
func NovoSistemaDeArquivoDeArmazenamentoDoJogador(arquivo *os.File) (*SistemaDeArquivoDeArmazenamentoDoJogador, error) {

	err := iniciaArquivoBDDoJogador(arquivo)

	if err != nil {
		return nil, fmt.Errorf("problem iniciando arquivo bd do jogador, %v", err)
	}

	liga, err := NovaLiga(arquivo)

	if err != nil {
		return nil, fmt.Errorf("problema carregando armazenamento de jogador do arquivo %s, %v", arquivo.Name(), err)
	}

	return &SistemaDeArquivoDeArmazenamentoDoJogador{
		bancoDeDados: json.NewEncoder(&fita{arquivo}),
		liga:         liga,
	}, nil
}

func iniciaArquivoBDDoJogador(arquivo *os.File) error {
	arquivo.Seek(0, 0)

	info, err := arquivo.Stat()

	if err != nil {
		return fmt.Errorf("problema procurando informacao de arquivo do arquivo %s, %v", arquivo.Name(), err)
	}

	if info.Size() == 0 {
		arquivo.Write([]byte("[]"))
		arquivo.Seek(0, 0)
	}

	return nil
}

// PegaLiga retorna a pontuacao de todos os jogadores
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() Liga {
	sort.Slice(f.liga, func(i, j int) bool {
		return f.liga[i].Vitorias > f.liga[j].Vitorias
	})
	return f.liga
}

// PegaPontuacaoDoJogador retorna a pontuacao do jogador
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaPontuacaoDoJogador(nome string) int {

	jogador := f.liga.Busca(nome)

	if jogador != nil {
		return jogador.Vitorias
	}

	return 0
}

// SalvaVitoria vai armazenar uma vitoria para o jogador, aumentando se ja for conhecido
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {
	jogador := f.liga.Busca(nome)

	if jogador != nil {
		jogador.Vitorias++
	} else {
		f.liga = append(f.liga, Jogador{nome, 1})
	}

	f.bancoDeDados.Encode(f.liga)
}
