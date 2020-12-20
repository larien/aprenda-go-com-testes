package poquer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ArmazenamentoJogador stores pontuacao information about jogadores
type ArmazenamentoJogador interface {
	ObterPontuacaoDeJogador(nome string) int
	GravarVitoria(nome string)
	ObterLiga() Liga
}

// Jogador stores a nome with a number of wins
type Jogador struct {
	Nome              string
	ChamadasDeVitoria int
}

// ServidorJogador is a HTTP interface for jogador information
type ServidorJogador struct {
	armazenamento ArmazenamentoJogador
	http.Handler
}

const tipoConteudoJSON = "application/json"

// NovoServidorJogador creates a ServidorJogador with routing configured
func NovoServidorJogador(armazenamento ArmazenamentoJogador) *ServidorJogador {
	p := new(ServidorJogador)

	p.armazenamento = armazenamento

	router := http.NewServeMux()
	router.Handle("/liga", http.HandlerFunc(p.manipuladorLiga))
	router.Handle("/jogadores/", http.HandlerFunc(p.manipuladorJogadores))

	p.Handler = router

	return p
}

func (p *ServidorJogador) manipuladorLiga(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", tipoConteudoJSON)
	json.NewEncoder(w).Encode(p.armazenamento.ObterLiga())
}

func (p *ServidorJogador) manipuladorJogadores(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]

	switch r.Method {
	case http.MethodPost:
		p.processarVitoria(w, jogador)
	case http.MethodGet:
		p.mostrarResultado(w, jogador)
	}
}

func (p *ServidorJogador) mostrarResultado(w http.ResponseWriter, jogador string) {
	pontuacao := p.armazenamento.ObterPontuacaoDeJogador(jogador)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}

func (p *ServidorJogador) processarVitoria(w http.ResponseWriter, jogador string) {
	p.armazenamento.GravarVitoria(jogador)
	w.WriteHeader(http.StatusAccepted)
}
