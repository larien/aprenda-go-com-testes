package main

import "testing"

func TestOla(t *testing.T) {
	resultado := Ola()
	esperado := "Olá, mundo"

	if resultado != esperado {
		t.Errorf("resultado %q, esperado %q", resultado, esperado)
	}
}
