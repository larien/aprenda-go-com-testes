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

## Pensando no tipo de teste que queremos ver

Vamos nos lembrar de nossa mentalidade e objetivos ao começar:

- **Escreva o teste que queremos ver**. Pense em como gostaríamos de usar o código que vamos escrever do ponto de vista de quem irá usar.
- Concentre-se em _o que_ e _por que_, mas não se distraia com _como_.

Nosso pacote precisa oferecer uma função que pode ser apontada para uma pasta e nos retornar algumas publicações.

```go
var publicacoes []blogpublicacoes.Publicacao
publicacoes = blogpublicacoes.NovasPublicacoesDoSA("alguma-pasta")
```

Para escrever um teste em torno disso, precisaríamos de algum tipo de pasta de teste com alguns exemplos de publicações nela. _Não há nada de terrivelmente errado com isso_, mas você está fazendo algumas trocas:

- para cada teste, você pode precisar criar novos arquivos para testar um comportamento específico
- algum comportamento será difícil de testar, como falha ao carregar arquivos
- os testes serão executados um pouco mais devagar porque eles precisarão acessar o sistema de arquivos

Também estamos desnecessariamente nos acoplando a uma implementação específica do sistema de arquivos.

### Abstrações do sistema de arquivos introduzidas no Go 1.16

Go 1.16 introduziu uma abstração para sistemas de arquivos; o pacote [io/fs](https://golang.org/pkg/io/fs/).

> O pacote fs define interfaces básicas para um sistema de arquivos. Um sistema de arquivos pode ser fornecido pelo sistema operacional host, mas também por outros pacotes.

Isso nos permite tirar nosso acoplamento a um sistema de arquivos específico, o que nos permite injetar diferentes implementações de acordo com nossas necessidades.

> [No lado do produtor da interface, o novo tipo embed.FS implementa fs.FS, assim como zip.Reader. A nova função os.DirFS fornece uma implementação de fs.FS apoiada por uma árvore de arquivos do sistema operacional.](Https://golang.org/doc/go1.16#fs)

Se usarmos essa interface, os usuários de nosso pacote terão várias opções incorporadas à biblioteca padrão para usar. Aprender a aproveitar as interfaces definidas na biblioteca padrão do Go (por exemplo, `io.fs`, [`io.Reader`](https://golang.org/pkg/io/#Reader), [`io.Writer`](https://golang.org/pkg/io/#Writer)), é vital para escrever pacotes fracamente acoplados. Esses pacotes podem então ser reutilizados em contextos diferentes daqueles que você imaginou, com o mínimo de alarde de seus consumidores.

Em nosso caso, talvez nosso consumidor queira que as publicações sejam incorporadas ao binário Go em vez de arquivos em um sistema de arquivos "real"? De qualquer forma, _nosso código não precisa se preocupar_.

Para nossos testes, o pacote [testing/fstest](https://golang.org/pkg/testing/fstest/) nos oferece uma implementação de [io/FS](https://golang.org/pkg/io/fs/#FS) para usar, semelhante às ferramentas com as quais estamos familiarizados em [net/http/httptest](https://golang.org/pkg/net/http/httptest/).

Com essas informações, o seguinte parece uma abordagem melhor.

```go
var publicacoes blogpublicacoes.Publicacao
publicacoes = blogpublicacoes.NovasPublicacoesDoSA(algumSA)
```
