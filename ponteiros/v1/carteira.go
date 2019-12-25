package main

import "fmt"

// Bitcoin representa o número de Bitcoins
type Bitcoin int

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

// Carteira armazena o número de bitcoins que uma pessoa tem
type Carteira struct {
	saldo Bitcoin
}

// Depositar vai adicionar Bitcoins à carteira
func (c *Carteira) Depositar(quantidade Bitcoin) {
	c.saldo += quantidade
}

// Saldo retorna o número de Bitcoins que uma carteira tem
func (c *Carteira) Saldo() Bitcoin {
	return c.saldo
}
