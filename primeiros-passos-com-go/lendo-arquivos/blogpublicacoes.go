package blogpublicacoes

import (
	"io/fs"
)

func NovasPublicacoesDoSA(sistemaArquivos fs.FS) ([]Publicacao, error) {
	dir, err := fs.ReadDir(sistemaArquivos, ".")
	if err != nil {
		return nil, err
	}

	var publicacoes []Publicacao

	for _, a := range dir {
		publicacao, err := obterPublicacao(sistemaArquivos, a.Name())
		if err != nil {
			return nil, err //todo: se um arquivo falhar, devemos parar ou apenas ignorar?
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func obterPublicacao(sistemaArquivos fs.FS, arquivoNome string) (Publicacao, error) {
	publicacaoArquivo, err := sistemaArquivos.Open(arquivoNome)
	if err != nil {
		return Publicacao{}, err
	}

	defer publicacaoArquivo.Close()

	return novaPublicacao(publicacaoArquivo)
}
