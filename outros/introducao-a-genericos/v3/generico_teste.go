package main

import (
    "log"
)

func main() {
    VerificaIgual(1, 1)
    VerificaNaoIgual(1, 2)

    VerificaIgual(50, 100) // isso deve falhar

    VerificaNaoIgual(2, 2) // você não verá isso na impressão (print)

    VerificaIgual("CJ", "CJ")
}

func VerificaIgual(recebido, esperado interface{}) {
    if recebido != esperado {
        log.Fatalf("resultado: recebido %d, esperado %d", recebido, esperado)
    } else {
        log.Printf("PASSOU: %d é igual %d\n", recebido, esperado)
    }
}

func VerificaNaoIgual(recebido, esperado interface{}) {
    if recebido == esperado {
        log.Fatalf("FALHOU: recebido %d, esperado %d", recebido, esperado)
    } else {
        log.Printf("PASSOU: %d não é igual %d\n", recebido, esperado)
    }

}