package main

import (
	"reflect"
	"testing"
)

func TestSoma(t *testing.T) {
	t.Run("coleção de qualquer tamanho", func(t *testing.T) {

		numeros := []int{1, 2, 3}

		resultado := Soma(numeros)
		esperado := 6

		if resultado != esperado {
			t.Errorf("resultado %d, esperado %d, dado, %v", resultado, esperado, numeros)
		}
	})
}

func TestSomaTodoOResto(t *testing.T) {

	t.Run("faz as somas de alguns slices", func(t *testing.T) {
		resultado := SomaTodoOResto([]int{1, 2}, []int{0, 9})
		esperado := []int{2, 9}

		if !reflect.DeepEqual(resultado, esperado) {
			t.Errorf("resultado %v, esperado %v", resultado, esperado)
		}
	})

	t.Run("soma slices vazios de forma segura", func(t *testing.T) {
		resultado := SomaTodoOResto([]int{}, []int{3, 4, 5})
		esperado := []int{0, 9}

		if !reflect.DeepEqual(resultado, esperado) {
			t.Errorf("resultado %v, esperado %v", resultado, esperado)
		}
	})

}
