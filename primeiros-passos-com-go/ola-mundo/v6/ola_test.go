package main

import "testing"

func TestOla(t *testing.T) {
	verificaMensagemCorreta := func(t testing.TB, resultado, esperado string) {
		t.Helper()
		if resultado != esperado {
			t.Errorf("resultado %q, esperado %q", resultado, esperado)
		}
	}

	t.Run("para uma pessoa", func(t *testing.T) {
		resultado := Ola("Chris", "")
		esperado := "Olá, Chris"
		verificaMensagemCorreta(t, resultado, esperado)
	})

	t.Run("string vazia", func(t *testing.T) {
		resultado := Ola("", "")
		esperado := "Olá, Mundo"
		verificaMensagemCorreta(t, resultado, esperado)
	})

	t.Run("em espanhol", func(t *testing.T) {
		resultado := Ola("Elodie", espanhol)
		esperado := "Hola, Elodie"
		verificaMensagemCorreta(t, resultado, esperado)
	})

	t.Run("em francês", func(t *testing.T) {
		resultado := Ola("Lauren", frances)
		esperado := "Bonjour, Lauren"
		verificaMensagemCorreta(t, resultado, esperado)
	})
}
