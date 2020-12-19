package poker

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// FileSystemPlayerStore armazena os jogadores no sistema de arquivos
type FileSystemPlayerStore struct {
	baseDeDados *json.Encoder
	league      League
}

// NovoArmazenamentoSistemaDeArquivodeJogador cria uma FileSystemPlayerStore inicializando o armazenamento se necessário
func NovoArmazenamentoSistemaDeArquivodeJogador(file *os.File) (*FileSystemPlayerStore, error) {

	err := initialisePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("falha ao inicializar o arquivo bando de dados do jogador, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("falha lendo o armazenamento do jogador a partir do arquivo %s, %v", file.Nome(), err)
	}

	return &FileSystemPlayerStore{
		baseDeDados: json.NewEncoder(&tape{file}),
		league:      league,
	}, nil
}

func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("falha ao pegar informacoes sobre o arquivo do arquivo %s, %v", file.Nome(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

// ObterLiga retorna a pontuação de todos os jogadores
func (f *FileSystemPlayerStore) ObterLiga() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Vitorias > f.league[j].Vitorias
	})
	return f.league
}

// ObterPontuacaoDeJogador consulta os pontos do jogador
func (f *FileSystemPlayerStore) ObterPontuacaoDeJogador(name string) int {

	player := f.league.Find(name)

	if player != nil {
		return player.Vitorias
	}

	return 0
}

// RecordWin vai armazenar uma vitória para o jogador, incrementa o número de vitórias se já existir
func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Vitorias++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.baseDeDados.Encode(f.league)
}
