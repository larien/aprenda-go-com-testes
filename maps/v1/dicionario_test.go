package main

import "testing"

func TestBusca(t *testing.T) {
    dicionario := map[string]string{"teste": "isso é apenas um teste"}

    resultado := Busca(dicionario, "teste")
    esperado := "isso é apenas um teste"

    compararStrings(t, resultado, esperado)
}

func compararStrings(t *testing.T, resultado, esperado string) {
	t.Helper()

	if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s', dado '%s'", resultado, esperado, "test")
    }
}
