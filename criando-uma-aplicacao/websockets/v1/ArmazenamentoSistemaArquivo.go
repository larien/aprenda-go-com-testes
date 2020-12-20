package poquer

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// FileSystemPlayerStore stores players in the filesystem
type FileSystemPlayerStore struct {
	baseDeDados *json.Encoder
	league      Liga
}

// NewFileSystemPlayerStore creates a FileSystemPlayerStore initialising the armazenamento if needed
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	err := initialisePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player armazenamento from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		baseDeDados: json.NewEncoder(&tape{file}),
		league:      league,
	}, nil
}

// FileSystemPlayerStoreFromFile creates a ArmazenamentoJogador from the contents of a JSON file found at path
func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	armazenamento, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player armazenamento, %v ", err)
	}

	return armazenamento, closeFunc, nil
}

func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

// ObterLiga retorna the Pontuações of all the players
func (f *FileSystemPlayerStore) ObterLiga() Liga {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Vitorias > f.league[j].Vitorias
	})
	return f.league
}

// ObtemPontuacaoDoJogador retrieves a player's pontuação
func (f *FileSystemPlayerStore) ObtemPontuacaoDoJogador(nome string) int {

	player := f.league.Find(nome)

	if player != nil {
		return player.Vitorias
	}

	return 0
}

// GravarVitoria will armazenamento a win for a player, incrementing venceu if already known
func (f *FileSystemPlayerStore) GravarVitoria(nome string) {
	player := f.league.Find(nome)

	if player != nil {
		player.Vitorias++
	} else {
		f.league = append(f.league, Jogador{nome, 1})
	}

	f.baseDeDados.Encode(f.league)
}
