package v1

import (
	"testing"
)

func TestContador(t *testing.T) {
	t.Run("incrementar o contador 3 vezes resulta no valor 3", func(t *testing.T) {
		contador := Contador{}
		contador.Incrementa()
		contador.Incrementa()
		contador.Incrementa()

		verificaContador(t, contador, 3)
	})
}

func verificaContador(t *testing.T, resultado Contador, esperado int) {
	t.Helper()
	if resultado.Valor() != esperado {
		t.Errorf("resultado %d, esperado %d", resultado.Valor(), esperado)
	}
}
