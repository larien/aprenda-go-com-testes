package v1

import "sync"

// Contador vai incrementar um número
type Contador struct {
	mu    sync.Mutex
	valor int
}

// NovoContador retorna um novo Contador
func NovoContador() *Contador {
	return &Contador{}
}

// Incrementa o contador
func (c *Contador) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.valor++
}

// Value returns the current count
// Valor retorna o atual contador
func (c *Contador) Valor() int {
	return c.valor
}
