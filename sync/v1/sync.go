package v1

// Contador will increment a number
type Contador struct {
	valor int
}

// Inc the count
func (c *Contador) Inc() {
	c.valor++
}

// Valor returns the current count
func (c *Contador) Valor() int {
	return c.valor
}
