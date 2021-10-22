package blogpublicacoes_test

import (
	"reflect"
	"testing"
	"testing/fstest"

	blogpublicacoes "github.com/larien/aprenda-go-com-testes/primeiros-passos-com-go/lendo-arquivos"
)

func TestNovasPublicacoesBlog(t *testing.T) {
	const (
		primeiroCorpo = `Titulo: Publicação 1
Descricao: Descrição 1
Tags: tdd, go
---
Olá
Mundo`
		segundoCorpo = `Titulo: Publicação 2
Descricao: Descrição 2
Tags: rust, borrow-checker
---
B
L
M`
	)

	sa := fstest.MapFS{
		"ola-mundo.md":  {Data: []byte(primeiroCorpo)},
		"ola-mundo2.md": {Data: []byte(segundoCorpo)},
	}

	publicacoes, err := blogpublicacoes.NovasPublicacoesDoSA(sa)

	verificaErro(t, err)

	verificaPublicacoes(t, publicacoes, sa)

	verificaPublicacao(t, publicacoes[0], blogpublicacoes.Publicacao{
		Titulo:    "Publicação 1",
		Descricao: "Descrição 1",
		Tags:      []string{"tdd", "go"},
		Corpo: `Olá
Mundo`,
	})
}

func verificaErro(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func verificaPublicacoes(t *testing.T, publicacoes []blogpublicacoes.Publicacao, sa fstest.MapFS) {
	t.Helper()
	if len(publicacoes) != len(sa) {
		t.Errorf("obteve %d publicacoes, esperado %d publicacoes", len(publicacoes), len(sa))
	}
}

func verificaPublicacao(t *testing.T, resultado blogpublicacoes.Publicacao, esperado blogpublicacoes.Publicacao) {
	t.Helper()
	if !reflect.DeepEqual(resultado, esperado) {
		t.Errorf("resultado %+v, esperado %+v", resultado, esperado)
	}
}
