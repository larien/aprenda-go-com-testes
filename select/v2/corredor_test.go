package corredor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCorredor(t *testing.T) {

	servidorLento := criarServidorComAtraso(20 * time.Millisecond)
	servidorRapido := criarServidorComAtraso(0 * time.Millisecond)

	defer servidorLento.Close()
	defer servidorRapido.Close()

	URLLenta := servidorLento.URL
	URLRapida := servidorRapido.URL

	esperado := URLRapida
	resultado := Corredor(URLLenta, URLRapida)

	if resultado != esperado {
		t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
	}
}

func criarServidorComAtraso(atraso time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(atraso)
		w.WriteHeader(http.StatusOK)
	}))
}
