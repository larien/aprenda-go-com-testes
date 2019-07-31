# Injeção de dependência

[**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/injecao-de-dependencia)

Presume-se que você tenha lido a seção de `structs` antes, já que será necessário saber um pouco sobre interfaces para entender esse capítulo.

Há muitos mal entendidos relacionados à injeção de dependência na comunidade de programação. Se tudo der certo, esse guia vai te mostrar que:

-   Você não precisa de uma framework
-   Não torna seu design complexo demais
-   Facilita seus testes
-   Permite que você escreva funções ótimas para própósitos diversos.

Queremos criar uma função que cumprimenta alguém, assim como a que fizemos no capítulo [Olá, mundo](https://github.com/larien/learn-go-with-tests/tree/master/hello-world), mas dessa vez vamos testar o _print de verdade_.

Para recapitular, a função era parecida com isso:

```go
func Cumprimenta(nome string) {
    fmt.Printf("Olá, %s", nome)
}
```

Mas como podemos testar isso? Chamar `fmt.Printf` imprime na saída, o que torna a captura com a ferramenta de testes bem difícil para nós.

O que precisamos fazer é sermos capazes de **injetar** (que é só uma palavra chique para passar) a dependência de impressão.

**Nossa função não precisa se preocupar com** _**onde**_ **ou** _**como**_ **a impressão acontece, então vamos aceitar uma** _**interface**_ **ao invés de um tipo concreto.**

Se fizermos isso, podemos mudar a implementação para imprimir algo que controlamos para poder testá-lo. Na "vida real", você iria injetar em algo que escreve na saída.

Se dermos uma olhada no código fonte do `fmt.Printf`, podemos ver uma forma de começar:

```go
// Printf retorna o número de bytes escritos e algum erro de escrita encontrado.
func Printf(format string, a ...interface{}) (n int, err error) {
    return Fprintf(os.Stdout, format, a...)
}
```

Interessante! Por baixo dos panos, o `Printf` só chama o `Fprintf` passando o `os.Stdout`.

O que exatamente _é_ um `os.Stdout`? O que o `Fprintf` espera que passe para ele como primeiro argumento?

```go
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
    p := newPrinter()
    p.doPrintf(format, a)
    n, err = w.Write(p.buf)
    p.free()
    return
}
```

Um `io.Writer`:

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

Quanto mais você escreve código em Go, mais vai perceber que essa interface aparece bastante, pois é uma ótima interface de uso geral para "colocar esses dados em algum lugar".

Logo, sabemos que por baixo dos panos estamos usando o `Writer` para enviar nosso cumprimento para algum lugar. Vamos usar essa abstração existente para tornar nosso código testável e mais reutilizável.

## Escreva o teste primeiro

```go
func TestCumprimenta(t *testing.T) {
    buffer := bytes.Buffer{}
    Cumprimenta(&buffer, "Chris")

    resultado := buffer.String()
    esperado := "Olá, Chris"

    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }
}
```

O tipo `buffer` do pacote `bytes` implementa a interface `Writer`.

Logo, vamos utilizá-lo no nosso teste para enviá-lo como nosso `Writer` e depois podemos verificar o que foi escrito nele quando chamamos `Cumprimenta`.

## Execute o teste

O teste não vai compilar:

```bash
./id_test.go:10:7: too many arguments in call to Cumprimenta
    have (*bytes.Buffer, string)
    want (string)
```

```bash
./id_test.go:10:7: muitos argumentos na chamada de Cumprimenta
    obteve (*bytes.Buffer, string)
    esperado (string)
```

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

_Preste atenção no compilador_ e corrija o problema.

```go
func Cumprimenta(escritor io.Writer, nome string) {
    fmt.Printf("Olá, %s", nome)
}
```

`Olá, Chris id_test.go:16: resultado '', esperado 'Olá, Chris'`

O teste falha. Note que o nome está sendo impresso, mas está indo para a saída.

## Escreva código o suficiente para fazer o teste passar

Use o escritor para enviar o cumprimento para o buffer no nosso teste. Lembre-se que o `fmt.Fprintf` é parecido com o `fmt.Printf`, com a diferença de que leva um `Writer` em que a string é enviada, enquanto que o`fmt.Printf` redireciona para a saída por padrão.

```go
func Cumprimenta(escritor io.Writer, nome string) {
	fmt.Fprintf(escritor, "Olá, %s", nome)
}
```

Agora o teste vai passar.

## Refatoração

Antes, o compilador nos disse para passar um ponteiro para um `bytes.Buffer`. Isso está tecnicamente correto, mas não é muito útil.

Para demonstrar isso, tente utilizar a função `Cumprimenta` em uma aplicação Go onde queremos que imprima na saída.

```go
func main() {
    Cumprimenta(os.Stdout, "Elodie")
}
```

`./id.go:14:7: cannot use os.Stdout (type *os.File) as type *bytes.Buffer in argument to Cumprimenta`

`não é possível utilizar os.Stdout (tipo *os.File) como tipo *bytes.Buffer no argumento para Cumprimenta`

Como discutimos antes, o `fmt.Fprintf` te permite passar um `io.Writer`, que sabemos que o `os.Stdout` e `bytes.Buffer` implementam.

Se mudarmos nosso código para usar uma interface de propósito mais geral, podemos usá-la tanto nos testes quanto na nossa aplicação.

```go
package main

import (
    "fmt"
    "os"
    "io"
)

func Cumprimenta(escritor io.Writer, nome string) {
    fmt.Fprintf(escritor, "Olá, %s", nome)
}

func main() {
    Cumprimenta(os.Stdout, "Elodie")
}
```

## Mais sobre io.Writer

Quais outros lugares podemos escrever dados usando `io.Writer`? Para qual propósito geral nossa função `Cumprimenta` é feita?

### A internet

Execute o seguinte:

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
    Greet(w, "world")
}

func main() {
    http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler))
}
```

Run the program and go to [http://localhost:5000](http://localhost:5000). You'll see your greeting function being used.

HTTP servers will be covered in a later chapter so don't worry too much about the details.

When you write an HTTP handler, you are given an `http.ResponseWriter` and the `http.Request` that was used to make the request. When you implement your server you _write_ your response using the writer.

You can probably guess that `http.ResponseWriter` also implements `io.Writer` so this is why we could re-use our `Greet` function inside our handler.

## Resumo

Our first round of code was not easy to test because it wrote data to somewhere we couldn't control.

_Motivated by our tests_ we refactored the code so we could control _where_ the data was written by **injecting a dependency** which allowed us to:

-   **Test our code** If you can't test a function _easily_, it's usually because of dependencies hard-wired into a function _or_ global state. If you have a global database connection pool for instance that is used by some kind of service layer, it is likely going to be difficult to test and they will be slow to run. DI will motivate you to inject in a database dependency \(via an interface\) which you can then mock out with something you can control in your tests.
-   **Separate our concerns**, decoupling _where the data goes_ from _how to generate it_. If you ever feel like a method/function has too many responsibilities \(generating data _and_ writing to a db? handling HTTP requests _and_ doing domain level logic?\) DI is probably going to be the tool you need.
-   **Allow our code to be re-used in different contexts** The first "new" context our code can be used in is inside tests. But further on if someone wants to try something new with your function they can inject their own dependencies.

### What about mocking? I hear you need that for DI and also it's evil

Mocking will be covered in detail later \(and it's not evil\). You use mocking to replace real things you inject with a pretend version that you can control and inspect in your tests. In our case though, the standard library had something ready for us to use.

### The Go standard library is really good, take time to study it

By having some familiarity with the `io.Writer` interface we are able to use `bytes.Buffer` in our test as our `Writer` and then we can use other `Writer`s from the standard library to use our function in a command line app or in web server.

The more familiar you are with the standard library the more you'll see these general purpose interfaces which you can then re-use in your own code to make your software reusable in a number of contexts.

This example is heavily influenced by a chapter in [The Go Programming language](https://www.amazon.co.uk/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440), so if you enjoyed this, go buy it!
