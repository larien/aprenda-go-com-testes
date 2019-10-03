package v1

import (
	"sync"
	"testing"
)

func TestContador(t *testing.T) {

	t.Run("incrementar o contador 3 vezes o deixa com valor 3", func(t *testing.T) {
		contador := NovoContador()
		contador.Inc()
		contador.Inc()
		contador.Inc()

		assertContador(t, contador, 3)
	})

	t.Run("roda concorrentemente em seguranca", func(t *testing.T) {
		contadorDesejado := 1000
		contador := NovoContador()

		var wg sync.WaitGroup
		wg.Add(contadorDesejado)

		for i := 0; i < contadorDesejado; i++ {
			go func(w *sync.WaitGroup) {
				contador.Inc()
				w.Done()
			}(&wg)
		}
		wg.Wait()

		assertContador(t, contador, contadorDesejado)
	})

}

func assertContador(t *testing.T, recebido *Contador, desejado int) {
	t.Helper()
	if recebido.Value() != desejado {
		t.Errorf("recebido %d, desejado %d", recebido.Value(), desejado)
	}
}
