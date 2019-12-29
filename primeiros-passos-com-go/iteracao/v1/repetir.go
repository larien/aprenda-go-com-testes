package iteracao

// Repetir retorna o caractere repetido 5 vezes
func Repetir(caractere string) string {
	var repeticoes string
	for i := 0; i < 5; i++ {
		repeticoes = repeticoes + caractere
	}
	return repeticoes
}
