package inteiros

import (
	"fmt"
	"testing"
)

func TestAdicionador(t *testing.T) {
	soma := Adiciona(2, 2)
	esperado := 4

	if soma != esperado {
		t.Errorf("esperado '%d', resultado '%d'", esperado, soma)
	}
}

func ExemploAdiciona() {
	soma := Adiciona(1, 5)
	fmt.Println(soma)
	// Saida: 6
}
