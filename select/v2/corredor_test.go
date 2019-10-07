package corredor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCorredor(t *testing.T) {

	servidorLento := criarServidorDemorado(20 * time.Millisecond)
	servidorRapido := criarServidorDemorado(0 * time.Millisecond)

	defer servidorLento.Close()
	defer servidorRapido.Close()

	urlLenta := servidorLento.URL
	urlRapida := servidorRapido.URL

	esperado := urlRapida
	obteve := Corredor(urlLenta, urlRapida)

	if obteve != esperado {
		t.Errorf("obteve '%s', esperado '%s'", obteve, esperado)
	}
}

func criarServidorDemorado(demora time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(demora)
		w.WriteHeader(http.StatusOK)
	}))
}
