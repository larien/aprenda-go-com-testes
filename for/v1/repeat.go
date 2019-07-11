package iteration

// Repeat retorna o caracter repetido 5 vezes
func Repeat(character string) string {
	var repeated string
	for i := 0; i < 5; i++ {
		repeated = repeated + character
	}
	return repeated
}
