package iteration

const repeatCount = 5

// Repeat retorna o caracter repetido 5 vezes
func Repeat(character string) string {
	var repeated string
	for i := 0; i < repeatCount; i++ {
		repeated += character
	}
	return repeated
}
