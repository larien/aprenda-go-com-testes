package concorrencia

// VerificadorWebsite verifica uma URL, retornando uma booleana
type VerificadorWebsite func(string) bool

// VerificaWebsites recebe um VerificadorWebsite e um slice de URLs e retorna um map
// de URLs com o resultado da verificação de cada URL com a função VerificadorWebsite
func VerificaWebsites(vw VerificadorWebsite, urls []string) map[string]bool {
	resultados := make(map[string]bool)

	for _, url := range urls {
		resultados[url] = vw(url)
	}

	return resultados
}
