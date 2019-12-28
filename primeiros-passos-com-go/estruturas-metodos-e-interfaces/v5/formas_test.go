package main

import (
	"testing"
)

func TestPerimetro(t *testing.T) {
	retangulo := Retangulo{10.0, 10.0}
	resultado := Perimetro(retangulo)
	esperado := 40.0

	if resultado != esperado {
		t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
	}
}

func TestArea(t *testing.T) {
	verificaArea := func(t *testing.T, forma Forma, esperado float64) {
		t.Helper()
		resultado := forma.Area()

		if resultado != esperado {
			t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
		}
	}

	t.Run("retângulos", func(t *testing.T) {
		retangulo := Retangulo{12.0, 6.0}
		verificaArea(t, retangulo, 72.0)
	})

	t.Run("círculos", func(t *testing.T) {
		circulo := Circulo{10}
		verificaArea(t, circulo, 314.1592653589793)
	})
}
