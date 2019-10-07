package v1

import (
	"testing"
)

func TestContador(t *testing.T) {

	t.Run("incrementar o contador 3 vezes o deixa com valor 3", func(t *testing.T) {
		contador := Contador{}
		contador.Inc()
		contador.Inc()
		contador.Inc()

		assertContador(t, contador, 3)
	})
}

func assertContador(t *testing.T, recebido Contador, desejado int) {
	t.Helper()
	if recebido.Valor() != desejado {
		t.Errorf("recebido %d, desejado %d", recebido.Valor(), desejado)
	}
}
