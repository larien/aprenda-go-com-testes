package main

// Dicionario armazena definições de palavras
type Dicionario map[string]string

// Busca encontra uma palavra no dicionário
func (d Dicionario) Busca(palavra string) string {
	return d[palavra]
}
