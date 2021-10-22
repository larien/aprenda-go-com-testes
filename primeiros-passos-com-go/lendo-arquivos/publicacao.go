package blogpublicacoes

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Publicacao struct {
	Titulo    string
	Descricao string
	Tags      []string
	Corpo     string
}

const (
	tituloSeparador    = "Titulo: "
	descricaoSeparador = "Descricao: "
	tagsSeparador      = "Tags: "
)

func novaPublicacao(publicacaoCorpo io.Reader) (Publicacao, error) {
	escaner := bufio.NewScanner(publicacaoCorpo)

	lerLinha := func(tagNome string) string {
		escaner.Scan()
		return strings.TrimPrefix(escaner.Text(), tagNome)
	}

	return Publicacao{
		Titulo:    lerLinha(tituloSeparador),
		Descricao: lerLinha(descricaoSeparador),
		Tags:      strings.Split(lerLinha(tagsSeparador), ", "),
		Corpo:     lerCorpo(escaner),
	}, nil
}

func lerCorpo(escaner *bufio.Scanner) string {
	escaner.Scan() // ignorar uma linha
	buf := bytes.Buffer{}
	for escaner.Scan() {
		fmt.Fprintln(&buf, escaner.Text())
	}
	return strings.TrimSuffix(buf.String(), "\n")
}
