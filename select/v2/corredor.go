package corredor

import (
	"net/http"
)

// Corredor compara os tempos de resposta de a e b, retornando o mais rápido
func Corredor(a, b string) (vencedor string) {
	select {
	case <-ping(a):
		return a
	case <-ping(b):
		return b
	}
}

func ping(URL string) chan bool {
	ch := make(chan bool)
	go func() {
		http.Get(URL)
		close(ch)
	}()
	return ch
}
