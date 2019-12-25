package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestContagem(t *testing.T) {

	t.Run("imprime 3 até Vai!", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		Contagem(buffer, &SpyContagemOperacoes{})

		resultado := buffer.String()
		esperado := `3
2
1
Vai!`

		if resultado != esperado {
			t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
		}
	})

	t.Run("pausa antes de cada impressão", func(t *testing.T) {
		spyImpressoraSleep := &SpyContagemOperacoes{}
		Contagem(spyImpressoraSleep, spyImpressoraSleep)

		esperado := []string{
			pausa,
			escrita,
			pausa,
			escrita,
			pausa,
			escrita,
			pausa,
			escrita,
		}

		if !reflect.DeepEqual(esperado, spyImpressoraSleep.Chamadas) {
			t.Errorf("esperado %v chamadas, resultado %v", esperado, spyImpressoraSleep.Chamadas)
		}
	})
}

type SpyContagemOperacoes struct {
	Chamadas []string
}

func (s *SpyContagemOperacoes) Pausa() {
	s.Chamadas = append(s.Chamadas, pausa)
}

func (s *SpyContagemOperacoes) Write(p []byte) (n int, err error) {
	s.Chamadas = append(s.Chamadas, escrita)
	return
}

const escrita = "escrita"
const pausa = "pausa"
