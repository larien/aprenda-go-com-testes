package corredor

import (
	"net/http"
)

// Corredor compara os tempos de resposta de a e b, retornando o mais rapido
func Corredor(a, b string) (vencedor string) {
	select {
	case <-ping(a):
		return a
	case <-ping(b):
		return b
	}
}

func ping(url string) chan bool {
	ch := make(chan bool)
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch
}
