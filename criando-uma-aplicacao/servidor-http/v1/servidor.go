package main

import (
	"fmt"
	"net/http"
)

// ServidorJogador retorna valor fixo "20" para qualquer chamada
func ServidorJogador(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
