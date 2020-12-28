package poquer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// ArmazenamentoJogador armazena informação de pontuação sobre jogadores
type ArmazenamentoJogador interface {
	ObtemPontuacaoDoJogador(nome string) int
	GravarVitoria(nome string)
	ObterLiga() Liga
}

// Jogador armazena um nome com um número de vitórias
type Jogador struct {
	Nome     string
	Vitorias int
}

// ServidorJogador é uma interface HTTP para informações de jogador
type ServidorJogador struct {
	armazenamento ArmazenamentoJogador
	http.Handler
	template *template.Template
	partida  Jogo
}

const tipoConteudoJSON = "application/json"
const caminhoTemplateHTML = "jogo.html"

// NovoServidorJogador cria um ServidorJogador com rotas configuradas
func NovoServidorJogador(armazenamento ArmazenamentoJogador, partida Jogo) (*ServidorJogador, error) {
	p := new(ServidorJogador)

	tmpl, err := template.ParseFiles("jogo.html")

	if err != nil {
		return nil, fmt.Errorf("problema ao abrir %s %v", caminhoTemplateHTML, err)
	}

	p.partida = partida
	p.template = tmpl
	p.armazenamento = armazenamento

	roteador := http.NewServeMux()
	roteador.Handle("/liga", http.HandlerFunc(p.manipulaLiga))
	roteador.Handle("/jogadores/", http.HandlerFunc(p.manipulaJogadores))
	roteador.Handle("/partida", http.HandlerFunc(p.jogarJogo))
	roteador.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = roteador

	return p, nil
}

var atualizadorDeWebsocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := novoWebsocketServidorJogador(w, r)

	mensagemNumeroDeJogadores := ws.EsperarPelaMensagem()
	numeroDeJogadores, _ := strconv.Atoi(mensagemNumeroDeJogadores)
	p.partida.Começar(numeroDeJogadores, ws)

	vencedor := ws.EsperarPelaMensagem()
	p.partida.Terminar(vencedor)
}

func (p *ServidorJogador) jogarJogo(w http.ResponseWriter, r *http.Request) {
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
