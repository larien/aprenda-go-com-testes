package corredor

import (
	"net/http"
	"time"
)

// Corredor compara os tempos de resposta de a e b, retornando o mais rapido
func Corredor(a, b string) (vencedor string) {
	duracaoA := medirTempoResposta(a)
	duracaoB := medirTempoResposta(b)

	if duracaoA < duracaoB {
		return a
	}

	return b
}

func medirTempoResposta(url string) time.Duration {
	inicio := time.Now()
	http.Get(url)
	return time.Since(inicio)
}
