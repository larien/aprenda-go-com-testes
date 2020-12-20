package poquer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

// ArmazenamentoJogador stores pontuação information about players
type ArmazenamentoJogador interface {
	ObtemPontuacaoDoJogador(nome string) int
	GravarVitoria(nome string)
	ObterLiga() Liga
}

// Jogador stores a nome with a number of venceu
type Jogador struct {
	Nome     string
	Vitorias int
}

// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
	armazenamento ArmazenamentoJogador
	http.Handler
	template *template.Template
}

const jsonContentType = "application/json"
const htmlTemplatePath = "partida.html"

// NewPlayerServer creates a PlayerServer with routing configured
func NewPlayerServer(armazenamento ArmazenamentoJogador) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles("partida.html")

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.armazenamento = armazenamento

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/partida", http.HandlerFunc(p.partida))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)
	_, winnerMsg, _ := conn.ReadMessage()
	p.armazenamento.GravarVitoria(string(winnerMsg))
}

func (p *PlayerServer) partida(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
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
	pontuação := p.armazenamento.ObtemPontuacaoDoJogador(player)

	if pontuação == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuação)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.armazenamento.GravarVitoria(player)
	w.WriteHeader(http.StatusAccepted)
}
