# Linha de comando e estrutura de pacotes

[**Você pode encontrar os exemplos deste capítulo aqui**](/command-line)

Nosso gerente de produto quer [_pivotar_](https://pt.wikipedia.org/wiki/Startup#Dicion%C3%A1rio_com_os_termos_mais_usados_pelas_startups)
 e introduzir uma segunda aplicação - uma aplicação de 
linha de comando. 

Inicialmente, ela vai apenas ser capaz de gravar o que um jogador vence quando o usuário digita `Ruth venceu`. 
A intenção é eventualmente criar uma ferramenta para ajudar usuários a jogar poker.

O gerente de produto quer que o banco de dados seja compartilhado entre as duas aplicações para que a `liga` atualize 
de acordo com as vitórias gravadas nessa nova aplicação.

## Lembrando do código

Nós temos uma aplicação com um arquivo `main.go` que inicia um servidor HTTP. O servidor HTTP não é nosso interesse neste
exercício mas a abstração usada é. Ele depende de `PlayerStore`.

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
    GetLeague() League
}
```

No capítulo anterior, criamos um `FileSystemPlayerStore` que implementa essa mesma interface. Temos que poder reutilizar
parte dela para a nossa nova aplicação. 

## Primeiro vamos [refatorar](https://pt.wikipedia.org/wiki/Refatora%C3%A7%C3%A3o) um pouco

Nosso projeto precisa criar dois executáveis, nosso existente servidor web e o app de linha de comand. 

Antes de nos entretermos no nosso novo código, precisamos estruturar nosso projeto melhor para suportar isso.

Até agora todos os códigos foram colocador em uma única pasta, em uma estrutura parecida com essa

`$GOPATH/src/github.com/seu-nome/meu-app`

Para fazer qualquer aplicação em Go, é necessário uma função `main` dentro de um `package main`. Até agora todo nosso
código viveu dentro de `package main` e a função `func main` pode referenciar tudo. 

Isso foi legal e é uma boa prática não sair gerando estrutura com pacotes logo de início. Se você olhar dentro da biblioteca
padrão você vai ver bem pouco a utilização de pastas e estruturas.

Felizmente é bem fácil adicionar uma estrutura _quando precisar dela_.

Dentro do projeto existente crie uma pasta `cmd` com uma chamada `webserver` dentro dela \(ex: `mkdir -p cmd/webserver`\).

Mova o arquivo `main.go` para dentro dessa pasta.

Se você tiver o comand `tree` instalado você pode executar sua estrutura de pastas tem que parecer

```text
.
├── FileSystemStore.go
├── FileSystemStore_test.go
├── cmd
│   └── webserver
│       └── main.go
├── league.go
├── server.go
├── server_integration_test.go
├── server_test.go
├── tape.go
└── tape_test.go
```

Agora temos uma separação efetiva entre nossa aplicação e o código da biblioteca mas agora temos que mudar alguns nomes
de pacotes(package). Lembre-se que ao construir uma aplicação Go seu nome _deve_  ser `main`.

Mude todos os outros códigos para ter um pacote chamado `poker`.

Finalmente, temos que importar esse pacote no `main.go` para utilizá-lo na criação de nosso servidor web. Então podemos
usar nossa biblioteca chamando `poker.NomeDaFunção`.

Os caminhos de diretórios vão ser diferentes no seu computador, mas deveria parecer com isso: 

```go
package main

import (
    "log"
    "net/http"
    "os"
    "github.com/larien/learn-go-with-tests/command-line/v1"
)

const dbFileName = "game.db.json"

func main() {
    db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        log.Fatalf("falha ao abrir %s %v", dbFileName, err)
    }

    store, err := poker.NewFileSystemPlayerStore(db)

    if err != nil {
        log.Fatalf("falha ao criar sistema de arquivos para armazenar jogadores, %v ", err)
    }

    server := poker.NewPlayerServer(store)

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("nao foi possivel escutar na porta 5000 %v", err)
    }
}
```

O caminho da pasta pode parecer chocante, mas essa é a forma para importar _qualquer_ biblioteca pública no seu código.

Separando nosso código em um pacote isolado e enviando para um repositório público como o GitHub qualquer desenvolvedor
Go pode escrever código que importe esse pacote com as funcionalidades que disponibilizarmos. A primeira vez que você
tentar e executar ele vai reclamar que o pacote não existe mas tudo que precisa ser feito é executar `go get`.

[Além disso, usuários podem ver a documentação em godoc.org](https://godoc.org/github.com/larien/learn-go-with-tests/command-line/v1).

### Verificações finais

* Dentro do diretório raiz rode `go test` e valide que ainda está passando
* Vá dentro de `cmd/webserver` e rode `go run main.go`
  * Abra `http://localhost:5000/league` e veja que ainda está funcionando

### Estrutura inicial

Antes de escrever os testes, vamos adicionar uma nova aplicação que nosso projeto vai construir. Crie outro diretório
 dentro de `cmd` chamado `cli` \(command line interface\) e adicione um arquivo `main.go` com

```go
package main

import "fmt"

func main() {
    fmt.Println("Vamos jogar poker")
}
```

O primeiro requisito que vamos discutir is como gravar uma vitória quando o usuário digitar `{NomeDoJogador} venceu`.

## Escreva o teste antes

Sabemos que temos que escrever algo chamado `CLI` que vai nos permitir `Play(Jogar)` poker. Isso vai precisar ler o que
o usuário digita e então gravar a vitória no armazenamento `PlayerStore`.  

Antes de irmos muito longe, vamos apenas escrever um teste para verificar a integração com a `PlayerStore` funciona como
gostaríamos. 

Dentro de `CLI_test.go` \(no diretório raiz do projeto, não dentro de `cmd`\)

```go
func TestCLI(t *testing.T) {
    playerStore := &StubPlayerStore{}
    cli := &CLI{playerStore}
    cli.PlayPoker()

    if len(playerStore.winCalls) !=1 {
        t.Fatal("esperando uma chamada de vitoria mas nao recebi nenhuma")
    }
}
```

* Podemos usar nossa `StubPlayerStore` de outros testes
* Passamos nossa dependência dentro do nosso ainda não existente tipo `CLI`
* Iniciamos o jogo chamando um método que chamaremos de `PlayPoker`
* Validamos se a vitória foi registrada

## Tente rodar o teste

```text
# github.com/larien/learn-go-with-tests/command-line/v2
./cli_test.go:25:10: undefined: CLI
```

## Escreva o mínimo código para o teste rodar e verificarmos o próximo error

Neste ponto, você deveria estar confortável para criar nossa nova `CLI` struct(estrutura de dados) com os respectivos
campos necessários para nossa dependência e adicionar um método.

Você deveria acabar com um código como esse

```go
type CLI struct {
    playerStore PlayerStore
}

func (cli *CLI) PlayPoker() {}
```

Lembre-se que estamos apenas tentando fazer o teste rodar para validarmos que ele falha como esperamos

```text
--- FAIL: TestCLI (0.00s)
    cli_test.go:30: esperando uma chamada de vitoria mas nao recebi nenhuma
FAIL
```

## Escreva código suficiente para fazer ele passar

```go
func (cli *CLI) PlayPoker() {
    cli.playerStore.RecordWin("Cleo")
}
```

Isso deve fazer ele passar.

Agora, precisamos simular leando isso from `Stdin` \(o que o usuário digita\) para que fique registrado vitórias para 
jogadores específicos.

Vamos incrementar nosso teste para exercitar essa condição.

## Escreva o teste antes

```go
func TestCLI(t *testing.T) {
    in := strings.NewReader("Chris venceu\n")
    playerStore := &StubPlayerStore{}

    cli := &CLI{playerStore, in}
    cli.PlayPoker()

    if len(playerStore.winCalls) < 1 {
        t.Fatal("esperando uma chamada de vitoria mas nao recebi nenhuma")
    }

    got := playerStore.winCalls[0]
    want := "Chris"

    if got != want {
        t.Errorf("nao armazenou o vencedor correto, recebi '%s', esperava '%s'", got, want)
    }
}
```

`os.Stdin` é o que vamos usar no `main` para capturar o que for digitado pelo usuário. Ele é um `*File` por trás dos panos
 o que siginifica que implementa `io.Reader` o qual sabemos ser um jeito útil de capturar texto.

Nós criamos um `io.Reader` no nosso teste usando `strings.NewReader`, preenchendo ele com o que esperamos que o usuário digite.

## Tente rodar o teste

`./CLI_test.go:12:32: too many values in struct initializer`

## Escreva o mínimo código para o teste rodar e verificarmos o próximo error

Precisamos adicionar nossa nova dependência dentro de `CLI`.

```go
type CLI struct {
    playerStore PlayerStore
    in io.Reader
}
```

## Escreva código suficiente para fazer ele passar

```text
--- FAIL: TestCLI (0.00s)
    CLI_test.go:23: nao armazenou o vencedor correto, recebi 'Cleo', esperava 'Chris'
FAIL
```

Lembre-se de primeiro fazer o que for mais fácil

```go
func (cli *CLI) PlayPoker() {
    cli.playerStore.RecordWin("Chris")
}
```

O teste vai passar. Depois nós vamos adicionar outro teste que vai nos forçar a escrever mais código, mas antes, vamos 
[refatorar](https://pt.wikipedia.org/wiki/Refatora%C3%A7%C3%A3o).

## [Refatoração](https://pt.wikipedia.org/wiki/Refatora%C3%A7%C3%A3o)

No `server_test` anteriormente fizemos validações para saber se uma vitória é armazenada assim como temos aqui. Vamos mover
essa validação para dentro de um helper e manter o código [DRY](https://pt.wikipedia.org/wiki/Don%27t_repeat_yourself).

```go
func assertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
    t.Helper()

    if len(store.winCalls) != 1 {
        t.Fatalf("recebi %d chamadas de RecordWin esperava %d", len(store.winCalls), 1)
    }

    if store.winCalls[0] != winner {
        t.Errorf("nao armazenou o vencedor correto, recebi '%s' esperava '%s'", store.winCalls[0], winner)
    }
}
```

Agora troque a validação em ambos os arquivos `server_test.go` e `CLI_test.go`.

O teste deve agora parecer com

```go
func TestCLI(t *testing.T) {
    in := strings.NewReader("Chris venceu\n")
    playerStore := &StubPlayerStore{}

    cli := &CLI{playerStore, in}
    cli.PlayPoker()

    assertPlayerWin(t, playerStore, "Chris")
}
```

Agora vamos escrever _outro_ teste com uma variação do que o usuário digitou nos forçando a ler de verdade.

## Escreva o teste antes

```go
func TestCLI(t *testing.T) {

    t.Run("recorda vencedor chris digitado pelo usuario", func(t *testing.T) {
        in := strings.NewReader("Chris venceu\n")
        playerStore := &StubPlayerStore{}

        cli := &CLI{playerStore, in}
        cli.PlayPoker()

        assertPlayerWin(t, playerStore, "Chris")
    })

    t.Run("recorda vencedor cleo digitado pelo usuario", func(t *testing.T) {
        in := strings.NewReader("Cleo venceu\n")
        playerStore := &StubPlayerStore{}

        cli := &CLI{playerStore, in}
        cli.PlayPoker()

        assertPlayerWin(t, playerStore, "Cleo")
    })

}
```

## Tente rodar o teste

```text
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/recorda_vencedor_chris_digitado_pelo_usuario
    --- PASS: TestCLI/recorda_vencedor_chris_digitado_pelo_usuario (0.00s)
=== RUN   TestCLI/recorda_vencedor_cleo_digitado_pelo_usuario
    --- FAIL: TestCLI/recorda_vencedor_cleo_digitado_pelo_usuario (0.00s)
        CLI_test.go:27: nao armazenou o vencedor correto, recebi 'Chris' esperava 'Cleo'
FAIL
```

## Escreva código suficiente para fazer ele passar

Vamos usar o [`bufio.Scanner`](https://golang.org/pkg/bufio/) para ler o que foi digitado no `io.Reader`.

> O pacote bufio implementa [buffered](https://pt.wikipedia.org/wiki/Buffer_(ci%C3%AAncia_da_computa%C3%A7%C3%A3o)) [I/O](https://pt.wikipedia.org/wiki/Entrada/sa%C3%ADda).
 Ele encapsula um objeto io.Reader ou io.Writer, criando um outro objeto \(Reader ou Writer\) que também implementa a interface mas prover buffering e ajuda com entradas/saídas de textos.

Atualize o código para

```go
type CLI struct {
    playerStore PlayerStore
    in          io.Reader
}

func (cli *CLI) PlayPoker() {
    reader := bufio.NewScanner(cli.in)
    reader.Scan()
    cli.playerStore.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
    return strings.Replace(userInput, " venceu", "", 1)
}
```

O teste agora vai passar.

* `Scanner.Scan()` vai ler até o carácter de nova linha.
* Só então usamos `Scanner.Text()` para returnar a `string` lida pelo scanner.

Agora que temos alguns testes passando, devemos amarrar isso ao nosso `main`. Lembre-se que devemos sempre almejar ter
o código funcionando totalmente integrado o mais rápido que pudermos.

No `main.go` adicione o seguinte e execute. \(você pode ter que ajustar o caminho da segunda dependência para refletir 
o que tem no seu computador\)

```go
package main

import (
    "fmt"
    "github.com/larien/learn-go-with-tests/command-line/v3"
    "log"
    "os"
)

const dbFileName = "game.db.json"

func main() {
    fmt.Println("Vamos jogar poker")
    fmt.Println("Digite {Nome} venceu para registrar uma vitoria")

    db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        log.Fatalf("falha ao abrir %s %v", dbFileName, err)
    }

    store, err := poker.NewFileSystemPlayerStore(db)

    if err != nil {
        log.Fatalf("falha ao criar sistema de arquivos para armazenar jogadores, %v ", err)
    }

    game := poker.CLI{store, os.Stdin}
    game.PlayPoker()
}
```

Você deve receber um erro

```text
command-line/v3/cmd/cli/main.go:32:25: implicit assignment of unexported field 'playerStore' in poker.CLI literal
command-line/v3/cmd/cli/main.go:32:34: implicit assignment of unexported field 'in' in poker.CLI literal
```

O que está acontecendo é que por causa da tentativa de associar os campos `playerStore` e `in` na `CLI`. Eles são campos 
não exportados\(privados\). Nós _podemos_ fazer isso nos nossos testes porque o teste está no mesmo pacote da `CLI` \(`poker`\). 
Mas nosso `main` é um pacote `main` portanto não tem acesso.

Isso enfatiza a importância de _integrar seu código_. Nós definimos corretamente as dependências da `CLI` como privada 
\(porque não queremos expô-las para os usuários da `CLI`\) mas não criamos uma forma para os usuário construí-las.  

Existe alguma forma de identificarmos esse problema antes?

### `package mypackage_test`

Nos exemplos usados até agora, quando nós fazemos um arquivo para testes nós declaramos ele como pertencendo ao mesmo pacote
que estamos testando.

Tudo bem e fazer isso significa no pior dos casos que queremos testar algo que é pertecente somente aquele pacote
conseguimos acesso aos tipos não exportados.

Mas considerando que, _em geral_, advogamos para _não_ se fazer testes de coisas internas, como Go pode garantir isso?
 E se pudéssemos testar nosso código aonde somente temos acesso aos tipos exportados \(como em nossp `main`\)?

Quando você escreve um project com múltiplos pacotes eu recomendo fortmente que o nome to seu pacote tenha o sufixo `_test`.
 Fazendo isso você somente ter acesso aos tipos públicos no seu pacote. Isso ajuda nesse caso especificamente mas também
 ajuda a disciplinar o teste somente de APIs públicas. Se ainda assim você precisar testar coisa interna você pode criar
 um teste separado com o nome de pacote igual ao do que você quer testar.

A máxima do TDD é que se você não pode testar o seu código então provávelmente vai ser difícil para os usuários do seu
 código de integrar com ele. Fazendo uso de `package foo_test` vai forçar você à testar seu código como se você estivesse
 importando ele como vão fazer aqueles que importarem o seu pacote.  

Antes de consertar o `main` vamos mudar o nome de pacote do nosso teste dentro de `CLI_test.go` para `poker_test`.

Se sua IDE estiver bem configurada você vai de repente ver um monte de vermelho! Se você rodar o compilador vocês vai ver
os seguintes errors:

```text
./CLI_test.go:12:19: undefined: StubPlayerStore
./CLI_test.go:17:3: undefined: assertPlayerWin
./CLI_test.go:22:19: undefined: StubPlayerStore
./CLI_test.go:27:3: undefined: assertPlayerWin
```

Nós agora tropeçamenos nos problemas de desenho do pacote. Para testar nosso código nós criamos algumas funções auxíliares
 e tipos emulados sem exportá-los e portanto não estão mais disponíveis para uso no nosso `CLI_test` porque eles foram
 definidos somente nos arquivos com `_test.go` no pacote `poker`.
 
#### Queremos ter as funções auxíliares e tipos emulados disponível publicamente?

Está é uma discussão subjetiva. One argumento é que não queremos poluir a API do nosso pacote só para ter código que
 facilitam os tests.

Na apresentação ["Testes avançados em Go"](https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=53) do
 Mitchell Hashimoto, é descrito como eles advogam na HashiCorp isso para que usuários do pacote possam escrever testes
 sem ter que reinventar a roda escrevendo tipos emulados. No nosso caso, isso significa que qualquer um usando nosso
 pacote `poker` não tem que criar seus próprios `PlayStore` emulados se eles quiserem usar nosso código.  

Informalmente eu tenho usado esta técnica em outros pacotes compartilhados e tem se provado extremamente útil em termos
 de economizar tempo dos usuários quando eles integram com nossos pacotes.

Então vamos criar um arquivo chamado `testing.go` e adicionar nossos cógidos auxiliares nele.

```go
package poker

import "testing"

type StubPlayerStore struct {
    scores   map[string]int
    winCalls []string
    league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
    score := s.scores[name]
    return score
}

func (s *StubPlayerStore) RecordWin(name string) {
    s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
    return s.league
}

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
    t.Helper()

    if len(store.winCalls) != 1 {
        t.Fatalf("recebi %d chamadas de RecordWin esperava %d", len(store.winCalls), 1)
    }

    if store.winCalls[0] != winner {
        t.Errorf("nao armazenou o vencedor correto, recebi '%s' esperava '%s'", store.winCalls[0], winner)
    }
}

// tarega para você - adicionar os códigos restantes
```

Você precisar tornar essas funções públicas \(lembre-se que exportar em Go é feito apenas colocando a primeira letra em
maíusculo\) se você quiser que elas sejam expostas para quem importar esse pacote.

No nosso teste `CLI` você precisa chamar o código como se fosse usando de um pacote diferente.

```go
func TestCLI(t *testing.T) {

    t.Run("recorda vencedor chris digitado pelo usuario", func(t *testing.T) {
        in := strings.NewReader("Chris venceu\n")
        playerStore := &poker.StubPlayerStore{}

        cli := &poker.CLI{playerStore, in}
        cli.PlayPoker()

        poker.AssertPlayerWin(t, playerStore, "Chris")
    })

    t.Run("recorda vencedor cleo digitado pelo usuario", func(t *testing.T) {
        in := strings.NewReader("Cleo venceu\n")
        playerStore := &poker.StubPlayerStore{}

        cli := &poker.CLI{playerStore, in}
        cli.PlayPoker()

        poker.AssertPlayerWin(t, playerStore, "Cleo")
    })

}
```

Você vai ver que agora temos o mesmo problema que tivemos na `main`

```text
./CLI_test.go:15:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:15:39: implicit assignment of unexported field 'in' in poker.CLI literal
./CLI_test.go:25:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:25:39: implicit assignment of unexported field 'in' in poker.CLI literal
```

O jeito mais fácil de resolver isso é fazer um construtor como temos para outros tipos. Nós também vamos mudar o `CLI`
 para que ele armazene a `bufio.Scanner` ao invés do leitor pois ele vai ser automaticamente encapsulado no momento da
 construção.  

```go
type CLI struct {
    playerStore PlayerStore
    in          *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
    return &CLI{
        playerStore: store,
        in:          bufio.NewScanner(in),
    }
}
```

Fazendo isso, podemos simplificar e refatorar no código do leitor

```go
func (cli *CLI) PlayPoker() {
    userInput := cli.readLine()
    cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
    return strings.Replace(userInput, " venceu", "", 1)
}

func (cli *CLI) readLine() string {
    cli.in.Scan()
    return cli.in.Text()
}
```

Mude o teste para usar o esse construtor e valtamos a ter nossos testes passando.

Por último, podemos voltar para o nosso `main.go` e usar o construtor que acabamos de criar

```go
game := poker.NewCLI(store, os.Stdin)
```

Tente executar ele, digite "Bob venceu".

### Refatoração

Nós temos alguma repetição nas nossas respectivas aplicações aonde estamos abrindo um arquivo e criando um `FileSystemStore`
 a partir do seu conteúdo. Isso parece uma pequena fraqueza no desenho do nosso pacote então deveríamos fazer uma função
 nele para encapsular a abertura de arquivos dado um caminho e retornar a `PlayerStore`.

```go
func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
    db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        return nil, nil, fmt.Errorf("falha ao abrir %s %v", path, err)
    }

    closeFunc := func() {
        db.Close()
    }

    store, err := NewFileSystemPlayerStore(db)

    if err != nil {
        return nil, nil, fmt.Errorf("falha ao criar sistema de arquivos para armazenar jogadores, %v ", err)
    }

    return store, closeFunc, nil
}
```

Agora refatorando ambas aplicações para usar a função de criar o armazenamento.

#### Código da aplicação CLI

```go
package main

import (
    "github.com/larien/learn-go-with-tests/command-line/v3"
    "log"
    "os"
    "fmt"
)

const dbFileName = "game.db.json"

func main() {
    store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

    if err != nil {
        log.Fatal(err)
    }
    defer close()

    fmt.Println("Vamos jogar poker")
    fmt.Println("Digite {Nome} venceu para registrar uma vitoria")
    poker.NewCLI(store, os.Stdin).PlayPoker()
}
```

#### Código da aplicação do servidor Web

```go
package main

import (
    "github.com/larien/learn-go-with-tests/command-line/v3"
    "log"
    "net/http"
)

const dbFileName = "game.db.json"

func main() {
    store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

    if err != nil {
        log.Fatal(err)
    }
    defer close()

    server := poker.NewPlayerServer(store)

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("nao foi possivel escutar na porta 5000 %v", err)
    }
}
```

Note a simetria: mesmo sendo diferente interfaces de usuário o setup é quase idêntico. Isso dá impressão de uma boa
 validação do nosso desenho. E note também que `FileSystemPlayerStoreFromFile` retorna uma função `close` (fechar), que
 podemos encerrar o arquivo fundamental assim que terminarmos de usar o armazenamento.

## Resumindo

### Estrutura do pacote

Esse capítulo pretendia criar duas aplicações, reusar o código de domínio que escrevemos até agora. Para fazer isso,
 nós precisamos atualizar a estrutura do nosso pacote para que ela tivesse pastas separadas para nossos respectivos
 `main`s.

Fazendo isso nós enfrentamos problemas de integração devido a valores não exportados então demostrando o valor de trabalhar
 em pequenas "etapas" e integrar com frequência.

Aprendemos como `mypackage_test` ajudou a criar um ambiente de testes que prover a mesma experiência de outros pacotes
 integrando com nosso código, assim ajudando você a pegar problemas de integração e ver o quão fácil \(ou não\) é de usar
 seu código. 

### Lendo a entrada do usuário

Vimos como lendo do `os.Stdin` é muito fácil de usar pois ele implementa o `io.Reader`. Nós usamos `bufio.Scanner` para
 facilitar a leitura linha à linha do que o usuário digita.

### Abstração simples leva à simples reutilização de código

Quase não nos esforçamos para integrar a `PlayerStore` na nossa aplicação \(assim que fizemos alguns ajustes no pacode\)
 e subsequente testar foi muito fácil tambem porque nós decidimos também expor a versão emulada. 
It was almost no effort to integrate `PlayerStore` into our new application \(once we had made the package adjustments\) 
and subsequently testing was very easy too because we decided to expose our stub version too.
