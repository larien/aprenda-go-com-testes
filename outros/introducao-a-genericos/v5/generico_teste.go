package main

import (
    "log"
)

func main() {
    VerificaIgual(1, 1)
    VerificaNaoIgual(1, 2)
    VerificaIgual("CJ", "CJ")
    VerificaNaoIgual(1, "1")

// Essa parte foi comentada para que rode os próximos testes, esses testes impedem dos outros de serem rodados    
//    VerificaIgual(50, 100) // isso deve falhar
//    VerificaNaoIgual(2, 2) // você não verá isso na impressão (print)

}

func VerificaIgual(recebido, esperado interface{}) {
    if recebido != esperado {
        log.Fatalf("resultado: recebido %+v, esperado %+v", recebido, esperado)
    } else {
        log.Printf("PASSOU: %+v é igual %+v\n", recebido, esperado)
    }
}

func VerificaNaoIgual(recebido, esperado interface{}) {
    if recebido == esperado {
        log.Fatalf("FALHOU: recebido %+v, esperado %+v", recebido, esperado)
    } else {
        log.Printf("PASSOU: %+v não é igual %+v\n", recebido, esperado)
    }

}