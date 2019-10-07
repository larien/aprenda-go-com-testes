package corredor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func CorredorTeste(t *testing.T) {

	servidorLento := criarServidorDemorado(20 * time.Millisecond)
	servidorRapido := criarServidorDemorado(0 * time.Millisecond)

	defer servidorLento.Close()
	defer servidorRapido.Close()

	urlLenta := servidorLento.URL
	urlRapida := servidorRapido.URL

	quer := urlRapida
	obteve := Corredor(urlLenta, urlRapida)

	if obteve != quer {
		t.Errorf("obteve '%s', quer '%s'", obteve, quer)
	}
}

func criarServidorDemorado(demora time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(demora)
		w.WriteHeader(http.StatusOK)
	}))
}
