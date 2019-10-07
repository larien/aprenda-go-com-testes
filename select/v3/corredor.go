package corredor

import (
	"net/http"
	"time"
)

var limiteDezSegundos = 10 * time.Second

// Corredor compara os tempos de resposta de a e b, retornando o mais rapido
func Corredor(a, b string) (vencedor string, error error) {
	return CorredorConfiguravel(a, b, limiteDezSegundos)
}

// CorredorConfiguravel compara os tempos de resposta de a e b, retornando o mais rapido
func CorredorConfiguravel(a, b string, tempoLimite time.Duration) (vencedor string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
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
