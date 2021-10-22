package main

import (
    "log"
)

func main() {
    VerificaIgual(1, 1)
    VerificaIgual("1", "1")
    VerificaNaoIgual(1, 2)
    //VerificaIgual(1, "1") - descomente-me para ver o erro de compilação

}

func VerificaIgual[T comparable](recebido, esperado T) {
    if recebido != esperado {
        log.Fatalf("FALHOU: recebido %+v, esperado %+v", recebido, esperado)
    } else {
        log.Printf("PASSOU: %+v é igual %+v\n", recebido, esperado)
    }
}

func VerificaNaoIgual[T comparable](recebido, esperado T) {
    if recebido == esperado {
        log.Fatalf("FALHOU: recebido %+v, esperado %+v", recebido, esperado)
    } else {
        log.Printf("PASSOU: %+v não é igual  %+v\n", recebido, esperado)
    }
}