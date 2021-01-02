package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestObterJogadores(t *testing.T) {
	requisicao, _ := http.NewRequest(http.MethodGet, "/", nil)
	resposta := httptest.NewRecorder()

	ServidorJogador(resposta, requisicao)

	t.Run("obter pontuação de Maria", func(t *testing.T) {
		retornado := resposta.Body.String()
		esperado := "20"

		if retornado != esperado {
			t.Errorf("retornado '%s', esperado '%s'", retornado, esperado)
		}
	})

}
