package poquer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

// ArmazenamentoJogador stores pontuação information about jogadores
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

// ServidorJogador is a HTTP interface for jogador information
type ServidorJogador struct {
	armazenamento ArmazenamentoJogador
	http.Handler
	template *template.Template
}

const tipoConteudoJSON = "application/json"
const htmlTemplatePath = "partida.html"

// NovoServidorJogador cria um ServidorJogador com rotas configuradas
func NovoServidorJogador(armazenamento ArmazenamentoJogador) (*ServidorJogador, error) {
	p := new(ServidorJogador)

	tmpl, err := template.ParseFiles("partida.html")

	if err != nil {
		return nil, fmt.Errorf("problema ao abrir %s %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.armazenamento = armazenamento

	router := http.NewServeMux()
	router.Handle("/liga", http.HandlerFunc(p.manipulaLiga))
	router.Handle("/jogadores/", http.HandlerFunc(p.manipulaJogadores))
	router.Handle("/partida", http.HandlerFunc(p.partida))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
}

var atualizadorDeWebsocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
	conexão, _ := atualizadorDeWebsocket.Upgrade(w, r, nil)
	_, winnerMsg, _ := conexão.ReadMessage()
	p.armazenamento.GravarVitoria(string(winnerMsg))
}

func (p *ServidorJogador) partida(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}

func (p *ServidorJogador) manipulaLiga(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", tipoConteudoJSON)
	json.NewEncoder(w).Encode(p.armazenamento.ObterLiga())
}

func (p *ServidorJogador) manipulaJogadores(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]

	switch r.Method {
	case http.MethodPost:
		p.processarVitoria(w, jogador)
	case http.MethodGet:
		p.mostrarPontuacao(w, jogador)
	}
}

func (p *ServidorJogador) mostrarPontuacao(w http.ResponseWriter, jogador string) {
	pontuação := p.armazenamento.ObtemPontuacaoDoJogador(jogador)

	if pontuação == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuação)
}

func (p *ServidorJogador) processarVitoria(w http.ResponseWriter, jogador string) {
	p.armazenamento.GravarVitoria(jogador)
	w.WriteHeader(http.StatusAccepted)
}
