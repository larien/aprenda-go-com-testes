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

## Desenvolvimento iterativo orientado a testes

Faremos uma abordagem iterativa em que estamos sempre dando passos simples e seguros em direção ao nosso objetivo.

Isso exige que interrompamos nosso trabalho, mas devemos ter cuidado para não cair na armadilha de adotar uma abordagem ["de baixo para cima"](https://pt.wikipedia.org/wiki/Abordagem_top-down_e_bottom-up).

Não devemos confiar em nossa imaginação hiperativa quando começamos a trabalhar. Poderíamos ser tentados a fazer algum tipo de abstração que só é validada quando juntamos tudo, como algum tipo de `BlogPublicacaoArquivoAnalisador`.

Isso _não_ é iterativo e está perdendo os loops de feedback menor que o TDD deveria nos trazer.

Kent Beck diz:

> O otimismo é um risco ocupacional da programação. O feedback é o tratamento.

Em vez disso, nossa abordagem deve se esforçar para estar o mais perto possível de entregar valor _real_ ao consumidor o mais rápido possível (geralmente chamado de "caminho feliz"). Depois de entregar uma pequena quantidade de valor ao consumidor de ponta a ponta, a iteração posterior do restante dos requisitos é geralmente direta.
