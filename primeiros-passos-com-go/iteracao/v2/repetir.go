package iteracao

const quantidadeRepeticoes = 5

// Repetir retorna o caractere repetido 5 vezes
func Repetir(caractere string) string {
	var repeticoes string
	for i := 0; i < quantidadeRepeticoes; i++ {
		repeticoes += caractere
	}
	return repeticoes
}
