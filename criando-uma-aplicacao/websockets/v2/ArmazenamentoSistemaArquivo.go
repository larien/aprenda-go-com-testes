package poquer

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// SistemaArquivoArmazenamentoJogador armazena jogadores no sistema de arquivos
type SistemaArquivoArmazenamentoJogador struct {
	baseDeDados *json.Encoder
	liga        Liga
}

// NovoSistemaArquivoArmazenamentoJogador cria um SistemaArquivoArmazenamentoJogador iniciaizando o armazenamento se necessário
func NovoSistemaArquivoArmazenamentoJogador(arquivo *os.File) (*SistemaArquivoArmazenamentoJogador, error) {

	err := inicializaArquivoDBJogador(arquivo)

	if err != nil {
		return nil, fmt.Errorf("problema ao inicializar o arquivo de base de dados do jogador, %v", err)
	}

	liga, err := NovaLiga(arquivo)

	if err != nil {
		return nil, fmt.Errorf("problema ao carregar o armazenamento do jogador do arquivo %s, %v", arquivo.Name(), err)
	}

	return &SistemaArquivoArmazenamentoJogador{
		baseDeDados: json.NewEncoder(&Tape{arquivo}),
		liga:        liga,
	}, nil
}

// SistemaArquivoArmazenamentoJogadorDoArquivo cria um ArmazenamentoJogador dos conteúdos de um arquivo JSON encontrado em um caminho
func SistemaArquivoArmazenamentoJogadorDoArquivo(path string) (*SistemaArquivoArmazenamentoJogador, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problema ao abrir %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	armazenamento, err := NovoSistemaArquivoArmazenamentoJogador(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problema ao criar sistema de arquivo de armazenamento do jogador, %v ", err)
	}

	return armazenamento, closeFunc, nil
}

func inicializaArquivoDBJogador(arquivo *os.File) error {
	arquivo.Seek(0, 0)

	info, err := arquivo.Stat()

	if err != nil {
		return fmt.Errorf("problema ao obter informações do arquivo do arquivo %s, %v", arquivo.Name(), err)
	}

	if info.Size() == 0 {
		arquivo.Write([]byte("[]"))
		arquivo.Seek(0, 0)
	}

	return nil
}

// ObterLiga retorna as Pontuações de todos os jogadores
func (s *SistemaArquivoArmazenamentoJogador) ObterLiga() Liga {
	sort.Slice(s.liga, func(i, j int) bool {
		return s.liga[i].Vitorias > s.liga[j].Vitorias
	})
	return s.liga
}

// ObtemPontuacaoDoJogador retorna a pontuação de um jogador
func (s *SistemaArquivoArmazenamentoJogador) ObtemPontuacaoDoJogador(nome string) int {

	jogador := s.liga.Encontrar(nome)

	if jogador != nil {
		return jogador.Vitorias
	}

	return 0
}

// GravarVitoria armazena uma vitória para um jogador, incrementando as vitórias já conhecidas
func (s *SistemaArquivoArmazenamentoJogador) GravarVitoria(nome string) {
	jogador := s.liga.Encontrar(nome)

	if jogador != nil {
		jogador.Vitorias++
	} else {
		s.liga = append(s.liga, Jogador{nome, 1})
	}

	s.baseDeDados.Encode(s.liga)
}
