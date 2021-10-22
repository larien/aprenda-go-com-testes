# Lendo arquivos

- [**Você pode encontrar todo o código deste capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/main/primeiros-passos-com-go/lendo-arquivos)
- [Aqui está um vídeo meu trabalhando no problema e respondendo à perguntas na Twitch (conteúdo em inglês)](https://www.youtube.com/watch?v=nXts4dEJnkU)

Neste capítulo, aprenderemos como ler alguns arquivos, obter alguns dados deles e fazer algo útil.

Finja que está trabalhando com seu amigo para criar algum software de blog. A ideia é que um autor escreva suas publicações em markdown, com alguns metadados no topo do arquivo. Na inicialização, o servidor web lerá uma pasta para criar algumas `publicações` e, em seguida, uma função `NovoTratamento` separada usará essas `publicações` como uma fonte de dados para o servidor web do blog.

Fomos solicitados a criar o pacote que converte uma determinada pasta de arquivos de publicação de blog em uma coleção de publicações.

### Exemplo de dado

ola-mundo.md

```
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

## Refatoração

Mesmo que nossos testes estejam passando, não podemos usar nosso novo pacote fora deste contexto, porque ele está acoplado a uma implementação concreta `fstest.MapFS`. Mas não precisa ser assim. Mude o argumento para nossa função `NovasPublicacoesDoSA` para aceitar a interface da biblioteca padrão.

```go
func NovasPublicacoesDoSA(sistemaArquivos fs.FS) []Publicacao {
	dir, _ := fs.ReadDir(sistemaArquivos, ".")
	var publicacoes []Publicacao
	for range dir {
		publicacoes = append(publicacoes, Publicacao{})
	}
	return publicacoes
}
```

Execute novamente os testes: tudo deve estar funcionando.

### Manipulação de erros

Deixamos de lado o tratamento de erros anteriormente, quando nos concentramos em fazer o caminho feliz funcionar. Antes de continuar a iterar na funcionalidade, devemos reconhecer que podem ocorrer erros ao trabalhar com arquivos. Além de ler o diretório, podemos ter problemas ao abrir arquivos individuais. Vamos mudar nossa API (primeiro por meio de nossos testes, naturalmente) para que ela possa retornar um `erro`.

```go
func TestNovasPublicacoesBlog(t *testing.T) {
	sa := fstest.MapFS{
		"ola-mundo.md":  {Data: []byte("oi")},
		"ola-mundo2.md": {Data: []byte("hola")},
	}

	publicacoes, err := blogpublicacoes.NovasPublicacoesDoSA(sa)

	if err != nil {
		t.Fatal(err)
	}

	if len(publicacoes) != len(sa) {
		t.Errorf("obteve %d publicacoes, esperado %d publicacoes", len(publicacoes), len(sa))
	}
}
```

Execute o teste: ele deve reclamar do número errado de valores de retorno. Consertar o código é simples.

```go
func NovasPublicacoesDoSA(sistemaArquivos fs.FS) ([]Publicacao, error) {
	dir, err := fs.ReadDir(sistemaArquivos, ".")
	if err != nil {
		return nil, err
	}
	var publicacoes []Publicacao
	for range dir {
		publicacoes = append(publicacoes, Publicacao{})
	}
	return publicacoes, nil
}
```

Isso fará com que o teste passe. O praticante de TDD em você pode ficar irritado por não vermos um teste com falha antes de escrever o código para propagar o erro de `fs.ReadDir`. Para fazer isso "apropriadamente", precisaríamos de um novo teste onde injetamos um test-double (dublê de teste) falho de `fs.FS` para fazer `fs.ReadDir` retornar um `erro`.

```go
type StubFalhoSA struct {
}

func (s StubFalhoSA) Abrir(nome string) (fs.File, error) {
	return nil, errors.New("oh não, eu sempre falho")
}

// later
_, err := blogpublicacoes.NovasPublicacoesDoSA(StubFalhoSA{})
```

**nota**: O prefixo `Stub` significa que o objeto tem um comportamento fixo e previsível.

Isso deve dar a você confiança em nossa abordagem. A interface que estamos usando tem um método, o que torna trivial a criação de dublês de teste para testar diferentes cenários.

Em alguns casos, testar o tratamento de erros é a coisa sensata a se fazer, mas, em nosso caso, não estamos fazendo nada _interessante_ com o erro, estamos apenas propagando-o, então não vale a pena escrever um novo teste.

Logicamente, nossa próxima iteração será em torno de expandir nosso tipo `Publicacao` para que ele tenha alguns dados úteis.

## Escreva o teste primeiro

Começaremos com a primeira linha do exemplo de publicação do blog proposto, o campo de título.

Precisamos alterar o conteúdo dos arquivos de teste para que correspondam ao que foi especificado e, então, podemos afirmar de que foi analisado corretamente.

```go
func TestNovasPublicacoesBlog(t *testing.T) {
	sa := fstest.MapFS{
		"ola-mundo.md":  {Data: []byte("Titulo: Publicação 1")},
		"ola-mundo2.md": {Data: []byte("Titulo: Publicação 2")},
	}

	publicacoes, err := blogpublicacoes.NovasPublicacoesDoSA(sa)

	if err != nil {
		t.Fatal(err)
	}

	if len(publicacoes) != len(sa) {
		t.Errorf("obteve %d publicacoes, esperado %d publicacoes", len(publicacoes), len(sa))
	}

	got := publicacoes[0]
	want := blogpublicacoes.Publicacao{Titulo: "Publicação 1"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("resultado %+v, esperado %+v", got, want)
	}
}
```

## Execute o teste

```
./blogpublicacoes_test.go:28:37: unknown field 'Titulo' in struct literal of type blogpublicacoes.Publicacao
```

Isso significa que o campo `Titulo` não existe no tipo `Publicacao`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Adicione o novo campo ao nosso tipo `Publicacao` para que o teste seja executado

```go
type Publicacao struct {
	Titulo string
}
```

Execute novamente o teste e você deverá vê-lo falhando

```
--- FAIL: TestNovasPublicacoesBlog (0.00s)
    blogpublicacoes_test.go:31: resultado {Titulo:}, esperado {Titulo:Publicação 1}
```

## Escreva código o suficiente para fazer o teste passar

Precisamos abrir cada arquivo e extrair o título

```go
func NovasPublicacoesDoSA(sistemaArquivos fs.FS) ([]Publicacao, error) {
	dir, err := fs.ReadDir(sistemaArquivos, ".")
	if err != nil {
		return nil, err
	}

	var publicacoes []Publicacao

	for _, a := range dir {
		publicacao, err := obterPublicacao(sistemaArquivos, a)
		if err != nil {
			return nil, err //todo: se um arquivo falhar, devemos parar ou apenas ignorar?
		}

		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func obterPublicacao(sistemaArquivos fs.FS, a fs.DirEntry) (Publicacao, error) {
	publicacaoArquivo, err := sistemaArquivos.Open(a.Name())
	if err != nil {
		return Publicacao{}, err
	}

	defer publicacaoArquivo.Close()

	publicacaoDados, err := io.ReadAll(publicacaoArquivo)
	if err != nil {
		return Publicacao{}, err
	}

	publicacao := Publicacao{Titulo: string(publicacaoDados)[8:]}

	return publicacao, nil
}
```

Lembre-se de que nosso foco neste ponto não é escrever um código elegante, mas apenas chegar a um ponto em que tenhamos um programa funcionando.

Mesmo que pareça um passo pequeno adiante, ainda exigiu que escrevêssemos uma boa quantidade de código e fizéssemos algumas suposições com relação ao tratamento de erros. Este seria um ponto em que você deveria conversar com seus colegas e decidir a melhor abordagem.

A abordagem iterativa nos deu um feedback rápido de que nossa compreensão dos requisitos está incompleta.

`fs.FS` nos dá uma maneira de abrir um arquivo dentro dele, pelo nome, com seu método `Open`. A partir daí, lemos os dados do arquivo e, por enquanto, não precisamos de nenhuma análise sofisticada, apenas excluindo o texto `Título:` ao cortar a string.

## Refatoração

Separar o 'código que abre o arquivo' do 'código que ler o arquivo' tornará o código mais simples de entender e trabalhar.

```go
func obterPublicacao(sistemaArquivos fs.FS, a fs.DirEntry) (Publicacao, error) {
	publicacaoArquivo, err := sistemaArquivos.Open(a.Name())
	if err != nil {
		return Publicacao{}, err
	}

	defer publicacaoArquivo.Close()

	return novaPublicacao(publicacaoArquivo)
}

func novaPublicacao(publicacaoArquivo fs.File) (Publicacao, error) {
	publicacaoDados, err := io.ReadAll(publicacaoArquivo)
	if err != nil {
		return Publicacao{}, err
	}

	publicacao := Publicacao{Titulo: string(publicacaoDados)[8:]}

	return publicacao, nil
}
```

Ao refatorar novas funções ou métodos, tome cuidado e pense sobre os argumentos. Você que está projetando aqui e está livre para pensar profundamente sobre o que é apropriado, pois os testes estão passando. Pense em acoplamento e coesão. Neste caso, você deve se perguntar:

> A função `novaPublicacao` precisa ser acoplada a um `fs.File`? Estamos usando todos os métodos e dados desse tipo? O que nós _realmente_ precisamos?

Em nosso caso, nós apenas o usamos como um argumento para `io.ReadAll` que precisa de um `io.Reader`. Portanto, devemos desacoplar nossa função e solicitar um `io.Reader`.

```go
func novaPublicacao(publicacaoArquivo io.Reader) (Publicacao, error) {
	publicacaoDados, err := io.ReadAll(publicacaoArquivo)
	if err != nil {
		return Publicacao{}, err
	}

	publicacao := Publicacao{Titulo: string(publicacaoDados)[8:]}

	return publicacao, nil
}
```

Você pode fazer um argumento semelhante para nossa função `obterPublicacao`, que recebe um argumento `fs.DirEntry`, mas chamamos somente o método `Name()` para obter o nome do arquivo. Não precisamos de tudo isso; vamos nos separar desse tipo e passar o nome do arquivo como uma string. Aqui está o código totalmente refatorado:

```go
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

func novaPublicacao(publicacaoArquivo io.Reader) (Publicacao, error) {
	publicacaoDados, err := io.ReadAll(publicacaoArquivo)
	if err != nil {
		return Publicacao{}, err
	}

	publicacao := Publicacao{Titulo: string(publicacaoDados)[8:]}

	return publicacao, nil
}
```

A partir de agora, a maioria dos nossos esforços pode ser perfeitamente contida em `novaPublicacao`. As preocupações de abrir e iterar sobre os arquivos estão resolvidas e agora podemos nos concentrar em extrair os dados para nosso tipo `Publicacao`. Embora não seja tecnicamente necessário, os arquivos são uma boa maneira de agrupar logicamente coisas relacionadas, então movi o tipo `Publicacao` e` novaPublicacao` para um novo arquivo `publicacao.go`.

### Ajudante de teste

Devemos cuidar de nossos testes também. Faremos muitas afirmações de `Publicações`, então devemos escrever algum código para ajudar com isso

```go
func verificaPublicacao(t *testing.T, resultado blogpublicacoes.Publicacao, esperado blogpublicacoes.Publicacao) {
	t.Helper()
	if !reflect.DeepEqual(resultado, esperado) {
		t.Errorf("resultado %+v, esperado %+v", resultado, esperado)
	}
}
```

```go
verificaPublicacao(t, publicacoes[0], blogpublicacoes.Publicacao{Titulo: "Publicação 1"})
```

## Escreva o teste primeiro

Vamos estender nosso teste ainda mais para extrair a próxima linha do arquivo, a descrição. Fazer passar agora deve parecer confortável e familiar.

```go
func TestNovasPublicacoesBlog(t *testing.T) {
	const (
		primeiroCorpo = `Titulo: Publicação 1
		Descricao: Descrição 1`
		segundoCorpo = `Titulo: Publicação 2
		Descricao: Descrição 2`
	)

	sa := fstest.MapFS{
		"ola-mundo.md":  {Data: []byte(primeiroCorpo)},
		"ola-mundo2.md": {Data: []byte(segundoCorpo)},
	}

  // código escondido

	verificaPublicacao(t, publicacoes[0], blogpublicacoes.Publicacao{
		Titulo:    "Publicação 1",
		Descricao: "Descrição 1",
	})
}
```

## Execute o teste

```
./blogpublicacoes_test.go:36:3: unknown field 'Descricao' in struct literal of type blogpublicacoes.Publicacao
```

Isso significa que o campo `Descricao` não existe no tipo `Publicacao`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Adicione o novo campo em `Publicacao`.

```go
type Publicacao struct {
	Titulo    string
	Descricao string
}
```

Os testes agora devem ser compilados e falhar.

```
--- FAIL: TestNovasPublicacoesBlog (0.00s)
    blogpublicacoes_test.go:34: resultado {Titulo:Publicação 1
                        Descricao: Descrição 1 Descricao:}, esperado {Titulo:Publicação 1 Descricao:Descrição 1}
```

## Escreva código o suficiente para fazer o teste passar

A biblioteca padrão possui uma biblioteca útil para ajudá-lo a digitalizar os dados, linha por linha; [`bufio.Scanner`](https://golang.org/pkg/bufio/#Scanner)

> Scanner fornece uma interface conveniente para leitura de dados, como um arquivo de linhas de texto delimitadas por uma nova linha.

```go
func novaPublicacao(publicacaoArquivo io.Reader) (Publicacao, error) {
	escaner := bufio.NewScanner(publicacaoArquivo)

	escaner.Scan()
	tituloLinha := escaner.Text()

	escaner.Scan()
	descricaoLinha := escaner.Text()

	return Publicacao{Titulo: tituloLinha[8:], Descricao: descricaoLinha[13:]}, nil
}
```

Convenientemente, `NewScanner` também necessita de um `io.Reader` para ler, então não precisamos alterar os argumentos de nossa função. (Obrigado desacoplamento!)

Chame `Scan` para ler uma linha e extraia os dados usando `Text`.

Esta função nunca poderia retornar um `erro`. Seria tentador, neste ponto, removê-lo do tipo de retorno, mas sabemos que teremos que lidar com estruturas de arquivo inválidas mais tarde, portanto, podemos também deixá-lo.

## Refatoração

Estamos repetindo o processo de escanear uma linha e depois ler o texto. Sabemos que faremos essa operação pelo menos mais uma vez, é uma simples refatoração para DRY ("Don't Repeat Yourself" ou "não se repita", em português), então vamos começar com isso.

```go
func novaPublicacao(publicacaoArquivo io.Reader) (Publicacao, error) {
	escaner := bufio.NewScanner(publicacaoArquivo)

	lerLinha := func() string {
		escaner.Scan()
		return escaner.Text()
	}

	titulo := lerLinha()[8:]
	descricao := lerLinha()[13:]

	return Publicacao{Titulo: titulo, Descricao: descricao}, nil
}
```

Isso quase não salvou nenhuma linha de código, mas raramente é o ponto da refatoração. O que estou tentando fazer aqui é apenas separar o _o que_ do _como_ das linhas de leitura para tornar o código um pouco mais declarativo para o leitor.

Embora os números mágicos de 7 e 13 façam o trabalho, eles não são muito descritivos.

```go
const (
	tituloSeparador    = "Titulo: "
	descricaoSeparador = "Descricao: "
)

func novaPublicacao(publicacaoArquivo io.Reader) (Publicacao, error) {
	escaner := bufio.NewScanner(publicacaoArquivo)

	lerLinha := func() string {
		escaner.Scan()
		return escaner.Text()
	}

	titulo := lerLinha()[len(tituloSeparador):]
	descricao := lerLinha()[len(descricaoSeparador):]

	return Publicacao{Titulo: titulo, Descricao: descricao}, nil
}
```

Agora que estou olhando para o código com minha mente criativa de refatoração, gostaria de tentar fazer com que nossa função lerLinha se encarregasse de remover a tag. Também existe uma maneira mais legível de remover um prefixo de uma string com a função `strings.TrimPrefix`.

```go
func novaPublicacao(publicacaoCorpo io.Reader) (Publicacao, error) {
	escaner := bufio.NewScanner(publicacaoCorpo)

	lerLinha := func(tagNome string) string {
		escaner.Scan()
		return strings.TrimPrefix(escaner.Text(), tagNome)
	}

	return Publicacao{
		Titulo:    lerLinha(tituloSeparador),
		Descricao: lerLinha(descricaoSeparador),
	}, nil
}
```

Você pode gostar ou não dessa ideia, mas eu gosto. A questão é que, no estado de refatoração, estamos livres para brincar com os detalhes internos e você pode continuar executando seus testes para verificar se as coisas ainda se comportam corretamente. Sempre podemos voltar aos estados anteriores se não estivermos satisfeitos. A abordagem TDD nos dá essa licença para experimentar ideias com frequência, então temos mais possibilidades de escrever bons códigos.

O próximo requisito é extrair as tags da publicação. Se você estiver acompanhando, recomendo tentar implementá-lo sozinho antes de continuar a leitura. Agora você deve ter um ritmo bom e se sentir confiante para extrair a próxima linha e analisar os dados.

Para resumir, não vou passar pelas etapas de TDD, mas aqui está o teste com tags adicionadas.

```go
func TestNovasPublicacoesBlog(t *testing.T) {
	const (
		primeiroCorpo = `Titulo: Publicação 1
Descricao: Descrição 1
Tags: tdd, go`
		segundoCorpo = `Titulo: Publicação 2
Descricao: Descrição 2
Tags: rust, borrow-checker`
	)

    // código escondido

	verificaPublicacao(t, publicacoes[0], blogpublicacoes.Publicacao{
		Titulo:    "Publicação 1",
		Descricao: "Descrição 1",
		Tags:      []string{"tdd", "go"},
	})
}
```

Você só está se enganando se apenas copiar e colar o que escrevo. Para ter certeza de que estamos todos na mesma página, aqui está o meu código, que inclui a extração das tags.

```go
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
	}, nil
}
```

Sem surpresas aqui. Pudemos reutilizar `lerLinha` para obter a próxima linha para as tags e então dividi-las usando `strings.Split`.

A última iteração em nosso caminho feliz é extrair o corpo.

Aqui está um lembrete do formato de arquivo proposto.

```
Titulo: Olá, mundo TDD!
Descricao: Primeira publicação em nosso maravilhoso blog
Tags: tdd, go
---
Olá mundo!

O corpo das publicações começa após o `---`
```

Já lemos as primeiras 3 linhas. Precisamos então ler mais uma linha, descartá-la e o restante do arquivo conterá o corpo da publicação.

## Escreva o teste primeiro

Altere os dados de teste para que tenham o separador e um corpo com algumas novas linhas para verificar se pegamos todo o conteúdo.

```go
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
```

Adicione ao nosso `verificaPublicacao`

```go
	verificaPublicacao(t, publicacoes[0], blogpublicacoes.Publicacao{
		Titulo:    "Publicação 1",
		Descricao: "Descrição 1",
		Tags:      []string{"tdd", "go"},
		Corpo: `Olá
Mundo`,
	})
```

## Execute o teste

```
./blogpublicacoes_test.go:47:3: unknown field 'Corpo' in struct literal of type blogpublicacoes.Publicacao
```

Isso significa que o campo `Corpo` não existe no tipo `Publicacao`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Adicione `Corpo` ao `Publicacao` e o teste deve falhar.

```
--- FAIL: TestNovasPublicacoesBlog (0.00s)
    blogpublicacoes_test.go:43: resultado {Titulo:Publicação 1 Descricao:Descrição 1 Tags:[tdd go] Corpo:}, esperado {Titulo:Publicação 1 Descricao:Descrição 1 Tags:[tdd go] Corpo:Olá
        Mundo}
```

## Escreva código o suficiente para fazer o teste passar

1. Leia a próxima linha para ignorar o separador `---`.
2. Continue digitalizando até que não haja mais nada para digitalizar.

```go
func novaPublicacao(publicacaoCorpo io.Reader) (Publicacao, error) {
	escaner := bufio.NewScanner(publicacaoCorpo)

	lerLinha := func(tagNome string) string {
		escaner.Scan()
		return strings.TrimPrefix(escaner.Text(), tagNome)
	}

	titulo := lerLinha(tituloSeparador)
	descricao := lerLinha(descricaoSeparador)
	tags := strings.Split(lerLinha(tagsSeparador), ", ")

	escaner.Scan() // ignorar uma linha

	buf := bytes.Buffer{}
	for escaner.Scan() {
		fmt.Fprintln(&buf, escaner.Text())
	}
	corpo := strings.TrimSuffix(buf.String(), "\n")

	return Publicacao{
		Titulo:    titulo,
		Descricao: descricao,
		Tags:      tags,
		Corpo:     corpo,
	}, nil
}
```

- `escaner.Scan()` retorna um `bool` que indica se há mais dados para escanear, então podemos usar isso com um loop `for` para continuar lendo os dados até o final.
- Após cada `Scan()`, escrevemos os dados no buffer usando `fmt.Fprintln`. Usamos a versão que adiciona uma nova linha porque o escaner remove as novas linhas de cada linha, mas precisamos mantê-las.
- Precisamos apagar a última nova linha, para que não tenhamos um espaço em branco no final.

## Refatoração

Encapsular a ideia de obter o resto dos dados em uma função ajudará aos futuros leitores a entender rapidamente _o que_ está acontecendo em `novaPublicacao`, sem ter que se preocupar com especificações de implementação.

```go
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
```

## Iterando ainda mais

Criamos os caminhos de execução mais importantes, tomando o caminho mais curto para chegar ao nosso caminho feliz, mas claramente há um longo caminho a percorrer antes de estar pronto para produção.

Criamos as principais funcionalidades, tomando o caminho mais curto para chegar ao nosso caminho feliz, mas claramente há um longo caminho a percorrer antes de estar pronto para produção.

Não tratamos de:

- quando o formato do arquivo não está correto
- o arquivo não é um `.md`
- e se a ordem dos campos de metadados for diferente? Isso deveria ser permitido? Devemos ser capazes de lidar com isso?

Porém, o mais importante é que temos um software funcionando e definimos nossa interface. Os pontos acima são apenas mais iterações, mais testes para escrever e direcionar nosso comportamento. Para suportar qualquer um dos itens acima, não devemos alterar nosso _design_, apenas os detalhes de implementação.

Manter o foco no objetivo significa que tomamos as decisões importantes e as validamos em relação ao comportamento desejado, em vez de nos prendermos a questões que não afetarão o design geral.

## Resumo

`fs.FS`, e as outras mudanças no Go 1.16 nos fornecem algumas maneiras elegantes de ler dados de sistemas de arquivos e testá-los de forma simples.

Se você deseja experimentar o código "de verdade":

- Crie uma pasta `cmd` dentro do projeto, adicione um arquivo `main.go`
- Adicione o seguinte código

```go
import (
    blogposts "github.com/quii/fstest-spike"
    "log"
    "os"
)

func main() {
	publicacoes, err := blogposts.NewPostsFromFS(os.DirFS("publicacoes"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(publicacoes)
}
```

- Adicione alguns arquivos markdown em uma pasta `publicacoes` e execute o programa!

Observe a simetria entre o código de produção

```go
publicacoes, err := blogposts.NewPostsFromFS(os.DirFS("publicacoes"))
```

**nota**: O código está vindo de um repositório do autor do livro, por isso a função está em inglês

E os testes

```go
publicacoes, err := blogpublicacoes.NovasPublicacoesDoSA(sa)
```

É quando o TDD de cima para baixo, orientado para o consumidor, parece estar correto.

Um usuário de nosso pacote pode olhar nossos testes e rapidamente se familiarizar com o que ele deve fazer e como usá-lo. Como mantenedores, podemos estar _confiantes de que nossos testes são úteis porque eles são do ponto de vista do consumidor_. Não estamos testando detalhes de implementação ou outros detalhes incidentais, então podemos estar razoavelmente confiantes de que nossos testes nos ajudarão, ao invés de nos atrapalhar durante a refatoração.

Contando com boas práticas de engenharia de software, como [**injeção de dependência**](https://larien.gitbook.io/aprenda-go-com-testes/v/main/primeiros-passos-com-go/injecao-de-dependencia), nosso código é simples de testar e reutilizar.

Ao criar pacotes, mesmo que sejam apenas internos ao projeto, prefira uma abordagem de cima para baixo voltada para o consumidor. Isso o impedirá de imaginar designs excessivos e de fazer abstrações de que talvez nem precise e ajudará a garantir que os testes que você escreve sejam úteis.

A abordagem iterativa manteve cada etapa pequena, e o feedback contínuo nos ajudou a descobrir requisitos pouco claros, possivelmente mais cedo do que com outras abordagens mais ad-hoc.

### Escrita?

É importante notar que esses novos recursos só têm operações para _leitura_ de arquivos. Se o seu trabalho precisa ser escrito, você precisará procurar em outro lugar. Lembre-se de continuar pensando sobre o que a biblioteca padrão oferece atualmente, se você estiver escrevendo dados, você provavelmente deve dar uma olhada em interfaces existentes, como `io.Writer` para manter seu código pouco acoplado e reutilizável.

### Leitura adicional

- Esta foi uma introdução leve para `io/fs`. [Ben Congdon fez um excelente artigo (conteúdo em inglês)](https://benjamincongdon.me/blog/2021/01/21/A-Tour-of-Go-116s-iofs-package/) que foi de muita ajuda para escrever este capítulo.
- [Discussão sobre as interfaces do sistema de arquivos (conteúdo em inglês)](https://github.com/golang/go/issues/41190)
