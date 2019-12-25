package corredor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCorredor(t *testing.T) {

	t.Run("compara a velocidade de servidores, retornando o endereço do mais rapido", func(t *testing.T) {
		servidorLento := criarServidorDemorado(20 * time.Millisecond)
		servidorRapido := criarServidorDemorado(0 * time.Millisecond)

		defer servidorLento.Close()
		defer servidorRapido.Close()

		urlLenta := servidorLento.URL
		urlRapida := servidorRapido.URL

		esperado := urlRapida
		obteve, err := Corredor(urlLenta, urlRapida)

		if err != nil {
			t.Fatalf("não esperava um erro, mas obteve um %v", err)
		}

		if obteve != esperado {
			t.Errorf("obteve '%s', esperado '%s'", obteve, esperado)
		}
	})

	t.Run("retorna um erro se o servidor não responder dentro de 10s", func(t *testing.T) {
		servidor := criarServidorDemorado(25 * time.Millisecond)

		defer servidor.Close()

		_, err := CorredorConfiguravel(servidor.URL, servidor.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("esperava um erro, mas não obtive um.")
		}
	})
}

func criarServidorDemorado(demora time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(demora)
		w.WriteHeader(http.StatusOK)
	}))
}
