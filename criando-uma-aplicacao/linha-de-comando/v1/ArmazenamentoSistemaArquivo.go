package poquer

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// SistemaDeArquivoArmazenamentoJogador armazena os jogadores no sistema de arquivos
type SistemaDeArquivoArmazenamentoJogador struct {
	baseDeDados *json.Encoder
	liga        Liga
}

// NovoArmazenamentoSistemaDeArquivodeJogador cria uma SistemaDeArquivoArmazenamentoJogador inicializando o armazenamento se necessário
func NovoArmazenamentoSistemaDeArquivodeJogador(arquivo *os.File) (*SistemaDeArquivoArmazenamentoJogador, error) {

	err := inicializarArquivoDBJogador(arquivo)

	if err != nil {
		return nil, fmt.Errorf("falha ao inicializar o arquivo bando de dados do jogador, %v", err)
	}

	liga, err := NovaLiga(arquivo)

	if err != nil {
		return nil, fmt.Errorf("falha lendo o armazenamento do jogador a partir do arquivo %s, %v", arquivo.Name(), err)
	}

	return &SistemaDeArquivoArmazenamentoJogador{
		baseDeDados: json.NewEncoder(&tape{arquivo}),
		liga:        liga,
	}, nil
}

func inicializarArquivoDBJogador(arquivo *os.File) error {
	arquivo.Seek(0, 0)

	info, err := arquivo.Stat()

	if err != nil {
		return fmt.Errorf("falha ao pegar informacoes sobre o arquivo do arquivo %s, %v", arquivo.Name(), err)
	}

	if info.Size() == 0 {
		arquivo.Write([]byte("[]"))
		arquivo.Seek(0, 0)
	}

	return nil
}

// ObterLiga retorna a pontuação de todos os jogadores
func (s *SistemaDeArquivoArmazenamentoJogador) ObterLiga() Liga {
	sort.Slice(s.liga, func(i, j int) bool {
		return s.liga[i].ChamadasDeVitoria > s.liga[j].ChamadasDeVitoria
	})
	return s.liga
}

// ObterPontuacaoDeJogador consulta os pontos do jogador
func (s *SistemaDeArquivoArmazenamentoJogador) ObterPontuacaoDeJogador(nome string) int {

	jogador := s.liga.Encontrar(nome)

	if jogador != nil {
		return jogador.ChamadasDeVitoria
	}

	return 0
}

// GravarVitoria vai armazenar uma vitória para o jogador, incrementa o número de vitórias se já existir
func (s *SistemaDeArquivoArmazenamentoJogador) GravarVitoria(nome string) {
	jogador := s.liga.Encontrar(nome)

	if jogador != nil {
		jogador.ChamadasDeVitoria++
	} else {
		s.liga = append(s.liga, Jogador{nome, 1})
	}

	s.baseDeDados.Encode(s.liga)
}
