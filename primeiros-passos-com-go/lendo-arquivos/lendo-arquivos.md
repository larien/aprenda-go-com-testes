# Lendo arquivos

- [**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/main/primeiros-passos-com-go/lendo-arquivos)
- [Aqui está um vídeo meu trabalhando no problema e respondendo à perguntas na Twitch (conteúdo em inglês)](https://www.youtube.com/watch?v=nXts4dEJnkU)

Neste capítulo, aprenderemos como ler alguns arquivos, obter alguns dados deles e fazer algo útil.

Finja que está trabalhando com seu amigo para criar algum software de blog. A ideia é que um autor escreva suas publicações em markdown, com alguns metadados no topo do arquivo. Na inicialização, o servidor web lerá uma pasta para criar algumas `publicações` e, em seguida, uma função `NovoTratamento` separada usará essas `publicações` como uma fonte de dados para o servidor web do blog.

Fomos solicitados a criar o pacote que converte uma determinada pasta de arquivos de publicação de blog em uma coleção de publicações.

### Exemplo de dado

ola-mundo.md

```markdown
Titulo: Olá, mundo TDD!
Descricao: Primeira publicação em nosso maravilhoso blog
Tags: tdd, go

---

Olá mundo!

O corpo das publicações começa após o `---`
```

### Dado esperado

```go
type Publicação struct {
	Titulo, Descricao, Corpo string
	Tags []string
}
```
