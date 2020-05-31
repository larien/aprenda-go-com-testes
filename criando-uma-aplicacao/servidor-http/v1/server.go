package main

import (
	"fmt"
	"net/http"
)

// JogadorServidor retorna valor fixo "20" para qualquer chamada
func JogadorServidor(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
