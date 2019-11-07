package concorrencia

import "net/http"

// VerificaWebsite retorna verdadeiro se a URL retornar um status code 200, senão falso
func VerificaWebsite(url string) bool {
	resposta, err := http.Head(url)
	if err != nil {
		return false
	}

	if resposta.StatusCode != http.StatusOK {
		return false
	}

	return true
}
