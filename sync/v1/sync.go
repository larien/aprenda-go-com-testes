package v1

// Contador incrementa um número
type Contador struct {
	valor int
}

// Incrementa o contador
func (c *Contador) Incrementa() {
	c.valor++
}

// Valor retorna a contagem atual
func (c *Contador) Valor() int {
	return c.valor
}
