package main

import "testing"

func TestOla(t *testing.T) {
	verificaMensagemCorreta := func(t *testing.T, resultado, esperado string) {
		t.Helper()
		if resultado != esperado {
			t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
		}
	}
	t.Run("diz olá para uma pessoa", func(t *testing.T) {
		resultado := Ola("Chris", "")
		esperado := "Olá, Chris"
		verificaMensagemCorreta(t, resultado, esperado)
	})

	t.Run("diz olá, mundo quando uma string vazia é passada como parâmetro", func(t *testing.T) {
		resultado := Ola("", "")
		esperado := "Olá, Mundo"
		verificaMensagemCorreta(t, resultado, esperado)
	})

	t.Run("diz olá em espanhol", func(t *testing.T) {
		resultado := Ola("Elodie", espanhol)
		esperado := "Hola, Elodie"
		verificaMensagemCorreta(t, resultado, esperado)
	})

	t.Run("diz olá em francês", func(t *testing.T) {
		resultado := Ola("Lauren", frances)
		esperado := "Bonjour, Lauren"
		verificaMensagemCorreta(t, resultado, esperado)
	})
}
