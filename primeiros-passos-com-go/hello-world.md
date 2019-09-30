# Olá, mundo

**[Você pode encontrar os códigos abordados nesse capítulo aqui](https://github.com/larienmf/learn-go-with-tests/tree/master/hello-world)**

É comum o primeiro programa em uma nova linguagem ser um Olá, mundo.

No [capítulo anterior](install-go.md#go-environment) discutimos sobre como Go pode ser dogmático como onde você coloca seus arquivos.

Crie um diretório no seguinte caminho `$GOPATH/src/github.com/{seu-lindo-nome-de-usuario}/hello`.

Se você estiver num ambiente baseado em unix e seu nome de usuário do SO for "bob" e você está motivado em seguir as convenções do Go sobre `$GOPATH` (que é a maneira mais fácil de configurar) você pode rodar `mkdir -p $GOPATH/src/github.com/bob/hello`.

Crie um arquivo no diretório mencionado chamado `hello.go` e escreva o seguinte código. Para rodar-lo, basta executar `go run hello.go`.

```go
package main

import "fmt"

func main() {
    fmt.Println("Olá, mundo")
}
```

## Como isso funciona?

Quando tu escreves um programa em Go, você irá ter um pacote `main` definido com uma função(`func`) `main` dentro disso. Os pacotes são maneiras de agrupar códigos de Go juntos.

A palavra reservada `func` é como você define uma função com um nome e um corpo.

Usando `import "fmt"` nós estamos a importar um pacote que contém a função `Println` que irá ser utilizada para imprimir um valor na tela.

## Como testar isso?

Como você testaria isso? É bom separar seu "domínio"(seu código) do resto do mundo \(side-effects\). A função `fmt.Println` é um side effect \(que está imprimindo um valor no stdout\) e a string, nós estamos enviando dentro do seu próprio domínio.

Então, vamos separar essas referencias para ficar mais fácil para testarmos

```go
package main

import "fmt"

func Hello() string {
    return "Olá, mundo"
}

func main() {
    fmt.Println(Hello())
}
```

Nós criamos uma nova função usando `func` mas dessa vez nós adicionamos outra palavra reservada `string` na definição. Isso significa que essa função irá ter como retorno uma `string`.

Agora, criaremos outro arquivo chamado `hello_test.go` onde nós iremos escrever um teste para nossa função `hello`.

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello()
    want := "Olá, mundo"

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

Antes de explicar, vamos rodar o código. Rode `go test` no seu terminal. Isso deve passar! Para checar, tente quebrar de alguma forma o teste mudando a string `want`.

Perceba que você nào precisa usar várias frameworks de testes e ficar se complicando tentando instalar-las. Tudo o que você precisa é feito na mesma linguagem e a sintaxe é a mesma para o resto dos códigos que você irá escrever.

### Escrevendo testes

Escrever um teste é como escrever uma função, com algumas regras

* Ele precisa estar num arquivo com um nome parecido com `xxx_test.go`
* A função de teste precisa começar com a palavra `Test`
* A função de teste recebe apenas um único argumento `t *testing.T`

Por agora é o bastante para saber que o nosso `t` do tipo `*testing.T` é o nosso "hook"(gancho) dentro da framework de testes e assim você poderá fazer coisas como `t.Fail()` quando você precisar testar um erro.

We've covered some new topics:

#### `if`
Instruções If em Go são muito parecidas com a de outras linguagens.

#### Declarando Variáveis

Nós estamos declarando algumas variáveis com a sintaxe  `varName := value`, que nos permite reutilizar alguns valores nos nossos testes de maneira legível.

#### `t.Errorf`

Nós estamos chamando o _method_(método) `Errorf` no nosso `t` que irá imprimir uma mensagem e falhar o teste. O prefixo `f` significa que podemos formatar e montar uma string com valores inseridos dentro de valores de preenchimentos `%s`. Quando tu fizeres um teste falhar, deves ser bastante claro como isso tudo aconteceu.

Nós iremos explorar mais na frente a diferença entre métodos e funções.

### Go doc

Outra funcionalidade importante de Go é sua documentação. Você pode rodar a documentação localmente rodando `godoc -http :8000`. Se você for para [localhost:8000/pkg](http://localhost:8000/pkg) irá ver todos os pacotes instalados no seu sistema.

A vasta biblioteca padrão da linguagem tem uma documentação excelente e com exemplos. Navegando para [http://localhost:8000/pkg/testing/](http://localhost:8000/pkg/testing/) vale a pena dar uma olhada para verificar o que está disponível para você.

### Olá, VOCÊ

Agora que temos um teste, nós podemos iterar sobre nosso software de maneira segura.

No último exemplo, nós escrevemos o teste somente _depois_ do código ser escrito, apenas para que você pudesse ter um exemplo de como escrever um teste e declarar uma função. A partir de agora, estamos _escrevendo os testes primeiro_.

Nosso próximo requisito é nos deixar especificar quem recebe a saudação.

Vamos começar especificando esses requisitos em um teste. Estamos fazendo um TDD(desenvolvimento orientado a testes) bastante simples e que nos permite ter certeza que nosso teste está _testando_ o que nós precisamos. Quando você escreve testes retroativamente existe o risco que seu test pode continuar passando mesmo que o código não esteja funcionando como esperado.

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello("Chris")
    want := "Olá, Chris"

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

Agora, rodando `go test`, deve ter aparecido um erro de compilação

```text
./hello_test.go:6:18: too many arguments in call to Hello
    have (string)
    want ()
```

Quando você está usando uma linguagem estaticamente tipada como Go, é importante _escutar o compilador_. O compilador entende como seu código deve se encaixar, não delegando essa função a você.

Neste caso, o compilador está te falando o que você precisa fazer para continuar. Nós temos que mudar a nossa função `Hello` para receber apenas um argumento.

Edite a função `Hello` para que seja aceito um argumento do tipo string

```go
func Hello(name string) string {
    return "Olá, mundo"
}
```

Se você tentar rodar seus testes novamente, seu arquivo `main.go` irá falhar durante a compilação por que você não está passando um argumento. Passe "mundo" como argumento para fazer o teste passar.

```go
func main() {
    fmt.Println(Hello("mundo"))
}
```

Agora, quando você for rodar seus testes você verá algo parecido com isso

```text
hello_test.go:10: got 'Olá, mundo' want 'Olá, Chris''
```

Agora, finalmente temos um programa que compila mas não está satisfazendo os requisitos de acordo com o teste.

Vamos então fazer o teste passar usando o argumento `name` e concatenar com `Hello,`

```go
func Hello(name string) string {
    return "Hello, " + name
}
```

Quando você rodar os testes eles irão passar. É comum como parte do ciclo do TDD _refatorar_ o nosso código agora.

### Uma nota sobre versionamento de código

Nesse ponto, se você estiver usando um versionamento de código \(que você deveria estar fazendo!\) Eu faria um `commit` do código no estado atual. Agora, temos um software funcional suportado por um teste.

Apesar de que eu _não faria_ um push para a master, por que eu planejo refatorar em breve. É legal fazer um commit nesse ponto porque você pode se perder com o refactoring, fazendo um commit você pode sempre voltar para a última versão funcional do seu software.

Não tem muita coisa para refatorar aqui, mas nós podemos introduzir outro recurso da linguagem: _constantes_.

### Constantes

Constantes podem ser definidas como o exemplo abaixo:

```go
const englishHelloPrefix = "Olá, "
```

Agora, podemos refatorar nosso código

```go
const portugueseHelloPrefix = "Olá, "

func Hello(name string) string {
    return portugueseHelloPrefix + name
}
```

Depois da refatoração, rode novamente os seus testes para ter certeza que você não quebrou nada.

Constantes melhoraram a performance da nossa aplicação assim como evitam com que você crie uma string `"Hello, "` para cada vez que `Hello` é chamado.

Sendo mais claro, o aumento de performance é incrivelmente insignificante para esse exemplo! Mas vale a pena pensar em criar constantes para capturar o significado dos valores e, às vezes, para ajudar no desempenho.

## Olá, mundo... novamente

O próximo requisito é: quando nossa função for chamada com uma string vazia, ela precisa imprimir o valor padrão "Olá, mundo", ao invés de "Olá, ".

Começaremos escrevendo um novo teste que irá falhar

```go
func TestHello(t *testing.T) {

    t.Run("diga olá para as pessoas", func(t *testing.T) {
        got := Hello("Chris")
        want := "Olá, Chris"

        if got != want {
            t.Errorf("got '%s' want '%s'", got, want)
        }
    })

    t.Run("diga 'Olá, mundo' quando uma string vazia for passada", func(t *testing.T) {
        got := Hello("")
        want := "Olá, mundo"

        if got != want {
            t.Errorf("got '%s' want '%s'", got, want)
        }
    })

}
```

Aqui nós estamos introduzindo outra ferramenta em nosso arsenal de testes, _subtestes_. Às vezes, é útil agrupar testes em torno de uma "coisa" e, em seguida, ter _subtestes_ descrevendo diferentes cenários.

O benefício dessa abordagem é que você poderá construir um código que pode ser compartilhado por outros testes.

Há um código repetido quando verificamos se a mensagem é o que esperamos.

A refatoração não é _apenas_ o código de produção!

É importante que seus testes _sejam especificações claras_ do que o código precisa fazer.

Podemos e devemos refatorar nossos testes.

```go
func TestHello(t *testing.T) {

    assertCorrectMessage := func(t *testing.T, got, want string) {
        t.Helper()
        if got != want {
            t.Errorf("got '%s' want '%s'", got, want)
        }
    }

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        want := "Hello, Chris"
        assertCorrectMessage(t, got, want)
    })

    t.Run("empty string defaults to 'World'", func(t *testing.T) {
        got := Hello("")
        want := "Olá, mundo"
        assertCorrectMessage(t, got, want)
    })

}
```

O que fizemos aqui?

Refatoramos nossa asserção em uma função. Isso reduz a duplicação e melhora a legibilidade de nossos testes. No Go, você pode declarar funções dentro de outras funções e atribuí-las a variáveis. Você pode chamá-las, assim como as funções normais. Precisamos passar em `t * testing.T` para que possamos dizer ao código de teste que falhará quando necessário.

`t.Helper ()` é necessário para dizermos ao conjunto de testes que este é método auxiliar. Ao fazer isso, quando o teste falhar, o número da linha relatada estará em nossa chamada de função, e não dentro do nosso auxiliar de teste. Isso ajudará outros desenvolvedores a rastrear os problemas com maior facilidade. Se você ainda não entender, comente, faça um teste falhar e observe a saída do teste.

Now that we have a well-written failing test, let's fix the code, using an `if`.

Agora que temos um teste bem escrito falhando, vamos corrigir o código, usando um `if`.

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
    if name == "" {
        name = "World"
    }
    return englishHelloPrefix + name
}
```

Se executarmos nossos testes, veremos que ele satisfaz o novo requisito e não quebramos acidentalmente a outra funcionalidade.

### De volta controle de versão

Agora, estamos felizes com o código. Eu adicionaria mais um commit ao anterior, então apenas verifique o quão adorável ficou o nosso código com os testes.

### Disciplina

Vamos repassar o ciclo novamente

* Write a test
* Make the compiler pass
* Run the test, see that it fails and check the error message is meaningful
* Write enough code to make the test pass
* Refactor

On the face of it this may seem tedious but sticking to the feedback loop is important.

Not only does it ensure that you have _relevant tests_, it helps ensure _you design good software_ by refactoring with the safety of tests.

Seeing the test fail is an important check because it also lets you see what the error message looks like. As a developer it can be very hard to work with a codebase when failing tests do not give a clear idea as to what the problem is.

By ensuring your tests are _fast_ and setting up your tools so that running tests is simple you can get in to a state of flow when writing your code.

By not writing tests you are committing to manually checking your code by running your software which breaks your state of flow and you won't be saving yourself any time, especially in the long run.

## Keep going! More requirements

Goodness me, we have more requirements. We now need to support a second parameter, specifying the language of the greeting. If a language is passed in that we do not recognise, just default to English.

We should be confident that we can use TDD to flesh out this functionality easily!

Write a test for a user passing in Spanish. Add it to the existing suite.

```go
    t.Run("in Spanish", func(t *testing.T) {
        got := Hello("Elodie", "Spanish")
        want := "Hola, Elodie"
        assertCorrectMessage(t, got, want)
    })
```

Remember not to cheat! _Test first_. When you try and run the test, the compiler _should_ complain because you are calling `Hello` with two arguments rather than one.

```text
./hello_test.go:27:19: too many arguments in call to Hello
    have (string, string)
    want (string)
```

Fix the compilation problems by adding another string argument to `Hello`

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }
    return englishHelloPrefix + name
}
```

When you try and run the test again it will complain about not passing through enough arguments to `Hello` in your other tests and in `hello.go`

```text
./hello.go:15:19: not enough arguments in call to Hello
    have (string)
    want (string, string)
```

Fix them by passing through empty strings. Now all your tests should compile _and_ pass, apart from our new scenario

```text
hello_test.go:29: got 'Hello, Elodie' want 'Hola, Elodie'
```

We can use `if` here to check the language is equal to "Spanish" and if so change the message

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == "Spanish" {
        return "Hola, " + name
    }

    return englishHelloPrefix + name
}
```

The tests should now pass.

Now it is time to _refactor_. You should see some problems in the code, "magic" strings, some of which are repeated. Try and refactor it yourself, with every change make sure you re-run the tests to make sure your refactoring isn't breaking anything.

```go
const spanish = "Spanish"
const englishHelloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "

func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == spanish {
        return spanishHelloPrefix + name
    }

    return englishHelloPrefix + name
}
```

### French

* Write a test asserting that if you pass in `"French"` you get `"Bonjour, "`
* See it fail, check the error message is easy to read
* Do the smallest reasonable change in the code

You may have written something that looks roughly like this

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == spanish {
        return spanishHelloPrefix + name
    }

    if language == french {
        return frenchHelloPrefix + name
    }

    return englishHelloPrefix + name
}
```

## `switch`

When you have lots of `if` statements checking a particular value it is common to use a `switch` statement instead. We can use `switch` to refactor the code to make it easier to read and more extensible if we wish to add more language support later

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    prefix := englishHelloPrefix

    switch language {
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    }

    return prefix + name
}
```

Write a test to now include a greeting in the language of your choice and you should see how simple it is to extend our _amazing_ function.

### one...last...refactor?

You could argue that maybe our function is getting a little big. The simplest refactor for this would be to extract out some functionality into another function.

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
    switch language {
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    default:
        prefix = englishHelloPrefix
    }
    return
}
```

A few new concepts:

* In our function signature we have made a _named return value_ `(prefix string)`.
* This will create a variable called `prefix` in your function.
  * It will be assigned the "zero" value. This depends on the type, for example `int`s are 0 and for strings it is `""`.
    * You can return whatever it's set to by just calling `return` rather than `return prefix`.
  * This will display in the Go Doc for your function so it can make the intent of your code clearer.
* `default` in the switch case will be branched to if none of the other `case` statements match.
* The function name starts with a lowercase letter. In Go public functions start with a capital letter and private ones start with a lowercase. We don't want the internals of our algorithm to be exposed to the world, so we made this function private.

## Wrapping up

Who knew you could get so much out of `Olá, mundo`?

By now you should have some understanding of:

### Some of Go's syntax around

* Writing tests
* Declaring functions, with arguments and return types
* `if`, `const` and `switch`
* Declaring variables and constants

### The TDD process and _why_ the steps are important

* _Write a failing test and see it fail_ so we know we have written a _relevant_ test for our requirements and seen that it produces an _easy to understand description of the failure_
* Writing the smallest amount of code to make it pass so we know we have working software
* _Then_ refactor, backed with the safety of our tests to ensure we have well-crafted code that is easy to work with

In our case we've gone from `Hello()` to `Hello("name")`, to `Hello("name", "French")` in small, easy to understand steps.

This is of course trivial compared to "real world" software but the principles still stand. TDD is a skill that needs practice to develop but by being able to break problems down into smaller components that you can test you will have a much easier time writing software.