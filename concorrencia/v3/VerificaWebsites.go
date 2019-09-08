package concorrencia

// VerificadorWebsite verifica uma URL, retornando uma booleana
type VerificadorWebsite func(string) bool
type resultado struct {
	string
	bool
}

// VerificaWebsites recebe um VerificadorWebsite e um slice de URLs e retorna um map
// de URLs com o resultado da verificação de cada URL com a função VerificadorWebsite
func VerificaWebsites(vw VerificadorWebsite, urls []string) map[string]bool {
	resultados := make(map[string]bool)
	canalResultado := make(chan resultado)

	for _, url := range urls {
		go func(u string) {
			canalResultado <- resultado{u, vw(u)}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		resultado := <-canalResultado
		resultados[resultado.string] = resultado.bool
	}

	return resultados
}
