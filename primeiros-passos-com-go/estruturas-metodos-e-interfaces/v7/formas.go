package main

import "math"

// Forma é implementado por qualquer coisa que possa dizer qual é sua área
type Forma interface {
	Area() float64
}

// Retangulo tem as dimensões de um retângulo
type Retangulo struct {
	Largura float64
	Altura  float64
}

// Area retorna a área de um retângulo
func (r Retangulo) Area() float64 {
	return r.Largura * r.Altura
}

// Perimetro retorna o perímetro de um retângulo
func Perimetro(retangulo Retangulo) float64 {
	return 2 * (retangulo.Largura + retangulo.Altura)
}

// Circulo representa um círculo.
type Circulo struct {
	Raio float64
}

// Area retorna a área de um círculo
func (c Circulo) Area() float64 {
	return math.Pi * c.Raio * c.Raio
}

// Triangulo representa as dimensões de um triângulo
type Triangulo struct {
	Base   float64
	Altura float64
}

// Area retorna a área de um triângulo
func (t Triangulo) Area() float64 {
	return (t.Base * t.Altura) * 0.5
}
