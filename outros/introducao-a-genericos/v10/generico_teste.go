package main

import (
	"log"
)

type PilhaDeInteiros = Pilha
type PilhaDeTextos = Pilha

type Pilha struct {
	valores []interface{}
}

func (p *Pilha) Empilhar(valor interface{}) {
	p.valores = append(p.valores, valor)
}

func (p *Pilha) EstaVazio() bool {
	return len(p.valores) == 0
}

func (p *Pilha) Desempilhar() (interface{}, bool) {
	if p.EstaVazio() {
		var zero interface{}
		return zero, false
	}

	indice := len(p.valores) - 1
	el := p.valores[indice]
	p.valores = p.valores[:indice]
	return el, true
}

func main() {
	// PILHA DE INTEIROS

	minhaPilhaDeInteiros := new(PilhaDeInteiros)

	// verifica se a pilha está vazia
	VerificaVerdadeiro(minhaPilhaDeInteiros.EstaVazio())

	// adiciona alguma coisa e em seguida, verifica se a pilha não está vazia
	minhaPilhaDeInteiros.Empilhar(123)
	VerificaFalso(minhaPilhaDeInteiros.EstaVazio())

	// adiciona outra coisa e em seguida, desempilhe a pilha
	minhaPilhaDeInteiros.Empilhar(456)
	valor, _ := minhaPilhaDeInteiros.Desempilhar()
	VerificaIgual(valor, 456)
	valor, _ = minhaPilhaDeInteiros.Desempilhar()
	VerificaIgual(valor, 123)
	VerificaVerdadeiro(minhaPilhaDeInteiros.EstaVazio())

	// PILHA DE TEXTOS

	minhaPilhaDeTextos := new(PilhaDeTextos)

	// verifica se a pilha está vazia
	VerificaVerdadeiro(minhaPilhaDeTextos.EstaVazio())

	// adiciona alguma coisa e em seguida, verifica se a pilha não está vazia
	minhaPilhaDeTextos.Empilhar("um dois tres")
	VerificaFalso(minhaPilhaDeTextos.EstaVazio())

	// adiciona outra coisa e em seguida, desempilhe a pilha
	minhaPilhaDeTextos.Empilhar("quatro cinco seis")
	valorTexto, _ := minhaPilhaDeTextos.Desempilhar()
	VerificaIgual(valorTexto, "quatro cinco seis")
	valorTexto, _ = minhaPilhaDeTextos.Desempilhar()
	VerificaIgual(valorTexto, "um dois tres")
	VerificaVerdadeiro(minhaPilhaDeTextos.EstaVazio())
	
	minhaPilhaDeInteiros.Empilhar(1)
	minhaPilhaDeInteiros.Empilhar(2)
	primeiroNum, _ := minhaPilhaDeInteiros.Desempilhar()
	segundoNum, _ := minhaPilhaDeInteiros.Desempilhar()
	VerificaIgual(primeiroNum+segundoNum, 3)
}


func VerificaVerdadeiro(algo bool) {
	if algo {
		log.Printf("PASSOU: Esperava-se que fosse verdade e foi\n")
	} else {
		log.Fatalf("FALHOU: Esperava-se que fosse verdadeiro, mas foi falso")
	}
}

func VerificaFalso(algo bool) {
	if !algo {
		log.Printf("PASSOU: Esperava-se que fosse falso e foi\n")
	} else {
		log.Fatalf("FALHOU: Esperava-se que fosse falso mas foi verdadeiro")
	}
}

func VerificaIgual[T comparable](recebido, esperado T) {
	if recebido != esperado {
		log.Fatalf("FALHOU: recebido %+v, esperado %+v", recebido, esperado)
	} else {
		log.Printf("PASSOU: %+v é igual  %+v\n", recebido, esperado)
	}
}

func VerificaNaoIgual[T comparable](recebido, esperado T) {
	if recebido == esperado {
		log.Fatalf("FALHOU: recebido %+v, esperado %+v", recebido, esperado)
	} else {
		log.Printf("PASSOU: %+v não é igual  %+v\n", recebido, esperado)
	}
}