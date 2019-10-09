package v1

// Contador incrementa um número
type Contador struct {
	valor int
}

// Incrementa o Contador
func (c *Contador) Inc() {
	c.valor++
}

// Valor retorna o contador atual
func (c *Contador) Valor() int {
	return c.valor
}
