package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PlayerStore stores score information about players
type PlayerStore interface {
	ObterPontuacaoDeJogador(name string) int
	RecordWin(name string)
	ObterLiga() League
}

// Player stores a name with a number of wins
type Player struct {
	Nome     string
	Vitorias int
}

// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
	armazenamento PlayerStore
	http.Handler
}

const jsonContentType = "application/json"

// NewPlayerServer creates a PlayerServer with routing configured
func NewPlayerServer(armazenamento PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.armazenamento = armazenamento

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.armazenamento.ObterLiga())
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.armazenamento.ObterPontuacaoDeJogador(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.armazenamento.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
