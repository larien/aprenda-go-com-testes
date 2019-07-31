package main

import (
	"fmt"
	"io"
	"net/http"
)

// Cumprimenta envia um cumprimento personalizado ao escritor
func Cumprimenta(escritor io.Writer, nome string) {
	fmt.Fprintf(escritor, "Olá, %s", nome)
}

// HandlerMeuCumprimento diz Olá, mundo via HTTP
func HandlerMeuCumprimento(w http.ResponseWriter, r *http.Request) {
	Cumprimenta(w, "mundo")
}

func main() {
	err := http.ListenAndServe(":5000", http.HandlerFunc(HandlerMeuCumprimento))

	if err != nil {
		fmt.Println(err)
	}
}
