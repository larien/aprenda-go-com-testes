package main

import (
	"fmt"
	"net/http"
)

// JogadorServidor currently returns Hello, world given _any_ request
func JogadorServidor(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
