package corredor

import (
	"fmt"
	"net/http"
	"time"
)

var limiteDeDezSegundos = 10 * time.Second

// Corredor compara os tempos de resposta de a e b, retornando o mais rápido com tempo limite de 10s
func Corredor(a, b string) (vencedor string, error error) {
	return Configuravel(a, b, limiteDeDezSegundos)
}

// Configuravel compara os tempos de resposta de a e b, retornando o mais rápido
func Configuravel(a, b string, tempoLimite time.Duration) (vencedor string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(tempoLimite):
		return "", fmt.Errorf("tempo limite de espera excedido para %s e %s", a, b)
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
