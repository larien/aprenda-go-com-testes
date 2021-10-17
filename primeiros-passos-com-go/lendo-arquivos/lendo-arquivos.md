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
type Publicacao struct {
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

## Escreva o teste primeiro

Devemos manter o escopo pequeno e útil o quanto for possível. Se provarmos que podemos ler todos os arquivos em um diretório, será um bom começo. Isso nos dará confiança no software que estamos escrevendo. Podemos verificar se a contagem de `[]Publicacao` retornada é igual ao número de arquivos em nosso sistema de arquivos falsos.

Crie um novo projeto para começarmos.

- `mkdir blogpublicacoes`
- `cd blogpublicacoes`
- `go mod init github.com/{seu-nome}/blogpublicacoes`
- `touch blogpublicacoes_test.go`

```go
package blogpublicacoes_test

import (
	"testing"
	"testing/fstest"
)

func TestNovasPublicacoesBlog(t *testing.T) {
	sa := fstest.MapFS{
		"ola-mundo.md":  {Data: []byte("oi")},
		"ola-mundo2.md": {Data: []byte("hola")},
	}

	publicacoes := blogpublicacoes.NovasPublicacoesDoSA(sa)

	if len(publicacoes) != len(sa) {
		t.Errorf("obteve %d publicacoes, esperado %d publicacoes", len(publicacoes), len(sa))
	}
}
```

Observe que o pacote do nosso teste é `blogpublicacoes_test`. Lembre-se, quando o TDD é bem praticado, adotamos uma abordagem _orientada para o consumidor_: não queremos testar detalhes internos porque os _consumidores_ não se importam com eles. Ao anexar `_test` ao nome de nosso pacote pretendido, nós apenas acessamos membros exportados de nosso pacote - assim como um usuário real de nosso pacote.

Importamos [`testing/fstest`](https://golang.org/pkg/testing/fstest/) que nos dá acesso ao tipo [`fstest.MapFS`](https://golang.org/pkg/testing/fstest/#MapFS). Nosso falso sistema de arquivos irá passar `fstest.MapFS` para o nosso pacote.

> Um MapFS é um sistema de arquivos simples na memória para uso em testes, representado como um map de nomes de caminhos (argumentos para Abrir) para informações sobre os arquivos ou diretórios que eles representam.

Isso parece mais simples do que manter uma pasta de arquivos de teste. E será executado mais rápido.

Por fim, codificamos o uso de nossa API do ponto de vista do consumidor e verificamos se ela cria o número correto de publicações.

## Tente executar o teste

```
./blogpublicacoes_test.go:14:16: undefined: blogpublicacoes
```

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

O pacote não existe. Crie um novo arquivo `blogpublicacoes.go` e coloque `pacote blogpublicacoes` dentro dele. Você precisará importar esse pacote para seus testes. Para mim, as importações ficaram assim:

```go
import (
	"testing"
	"testing/fstest"

	blogpublicacoes "github.com/larien/aprenda-go-com-testes/primeiros-passos-com-go/lendo-arquivos"
)
```

Agora os testes não compilarão porque nosso novo pacote não possui uma função `NovasPublicacoesDoSA`, que retorna algum tipo de coleção.

```
./blogpublicacoes_test.go:16:16: undefined: blogpublicacoes.NovasPublicacoesDoSA
```

Isso nos força a fazer o esqueleto de nossa função para fazer o teste funcionar. Lembre-se de não pensar demais no código neste ponto; estamos apenas tentando fazer um teste em execução e garantir que ele falhe conforme o esperado. Se pularmos esta etapa, podemos pular suposições e escrever um teste que não é útil.

```go
package blogpublicacoes

import "testing/fstest"

type Publicacao struct {

}

func NovasPublicacoesDoSA(sistemaArquivos fstest.MapFS) []Publicacao {
	return nil
}
```

O teste agora deve falhar corretamente

```
--- FAIL: TestNovasPublicacoesBlog (0.00s)
    blogpublicacoes_test.go:19: obteve 0 publicacoes, esperado 2 publicacoes
```

## Escreva código o suficiente para fazer o teste passar

Nós _podemos_ slime\* para fazê-lo passar:

```go
func NovasPublicacoesDoSA(sistemaArquivos fstest.MapFS) []Publicacao {
	return []Publicacao{{}, {}}
}
```

Mas, como Denise Yu escreveu:

> Sliming\* é útil para dar um “esqueleto” ao seu objeto. Projetar uma interface e executar a lógica são duas preocupações, e os sliming\* testes estrategicamente permitem que você se concentre em um de cada vez.

Já temos nossa estrutura. Então, o que fazemos em vez disso?

Como cortamos o escopo, tudo o que precisamos fazer é ler o diretório e criar uma publicação para cada arquivo que encontrarmos. Não precisamos nos preocupar em abrir arquivos e analisá-los ainda.

```go
func NovasPublicacoesDoSA(sistemaArquivos fstest.MapFS) []Publicacao {
	dir, _ := fs.ReadDir(sistemaArquivos, ".")
	var publicacoes []Publicacao
	for range dir {
		publicacoes = append(publicacoes, Publicacao{})
	}
	return publicacoes
}
```

[`fs.ReadDir`](https://golang.org/pkg/io/fs/#ReadDir) lê um diretório dentro de um determinado `fs.FS` retornando [`[]DirEntry`](https://golang.org/pkg/io/fs/#DirEntry).

Nossa visão idealizada do mundo já foi frustrada porque erros podem acontecer, mas lembre-se agora que nosso foco é _passar no teste_, não alterar o design, então ignoraremos o erro por enquanto.

O resto do código é direto: itere sobre as entradas, crie uma `Publicacao` para cada uma e retorne `publicacoes`.