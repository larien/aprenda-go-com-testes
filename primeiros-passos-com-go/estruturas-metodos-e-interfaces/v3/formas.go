package main

// Retangulo tem as dimensões de um retângulo
type Retangulo struct {
	Largura float64
	Altura  float64
}

// Perimetro retorna o perímetro de um retângulo
func Perimetro(retangulo Retangulo) float64 {
	return 2 * (retangulo.Largura + retangulo.Altura)
}

// Area retorna a área de um retângulo
func Area(retangulo Retangulo) float64 {
	return retangulo.Largura * retangulo.Altura
}
