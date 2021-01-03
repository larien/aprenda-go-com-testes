# Injeção de dependência

[**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/master/primeiros-passos-com-go/injecao-de-dependencia)

Presume-se que você tenha lido a seção de `structs` antes, já que será necessário saber um pouco sobre interfaces para entender este capítulo.

Há muitos mal entendidos relacionados à injeção de dependência na comunidade de programação. Se tudo der certo, esse guia vai te mostrar que:

-   Você não precisa de uma framework
-   Não torna seu design complexo demais
-   Facilita seus testes
-   Permite que você escreva funções ótimas para propósitos diversos.

Queremos criar uma função que cumprimenta alguém, assim como a que fizemos no capítulo [Olá, mundo](../ola-mundo/ola-mundo.md), mas dessa vez vamos testar o _print de verdade_.

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

func Cumprimenta(escritor io.Writer, nome string) {
    fmt.Fprintf(escritor, "Olá, %s", nome)
}

func HandlerMeuCumprimento(w http.ResponseWriter, r *http.Request) {
    Cumprimenta(w, "mundo")
}

func main() {
    err := http.ListenAndServe(":5000", http.HandlerFunc(HandlerMeuCumprimento))

    if err != nil {
        fmt.Println(err)
    }
}
```

Execute o programa e vá para [http://localhost:5000](http://localhost:5000). Você verá sua função de cumprimento ser utilizada.

Falaremos sobre servidores HTTP em um próximo capítulo, então não se preocupe muito com os detalhes.

Quando se cria um handler HTTP, você recebe um `http.ResponseWriter` e o `http.Request` que é usado para fazer a requisição. Quando implementa seu servidor, você _escreve_ sua resposta usando o escritor.

Você deve ter adivinhado que o `http.ResponseWriter` também implementa o `io.Writer` e é por isso que podemos reutilizar nossa função `Cumprimenta` dentro do nosso handler.

## Resumo

Nossa primeira rodada de código não foi fácil de testar porque escrevemos dados em algum lugar que não podíamos controlar.

_Graças aos nossos testes_, refatoramos o código para que pudéssemos controlar para _onde_ os dados eram escritos **injetando uma dependência** que nos permitiu:

-   **Testar nosso código**: se você não consegue testar uma função _de forma simples_, geralmente é porque dependências estão acopladas em uma função _ou_ estado global. Se você tem um pool de conexão global da base de dados, por exemplo, é provável que seja difícil testar e vai ser lento para ser execudado. A injeção de dependência te motiva a injetar em uma dependência da base de dados (através de uma interface), para que você possa criar um mock com algo que você possa controlar nos seus testes.
-   _Separar nossas preocupações_, desacoplando _onde os dados vão_ de _como gerá-los_. Se você já achou que um método/função tem responsabilidades demais (gerando dados _e_ escrevendo na base de dados? Lidando com requisições HTTP _e_ aplicando lógica a nível de domínio?), a injeção de dependência provavelmente será a ferramenta que você precisa.
-   **Permitir que nosso código seja reutilizado em contextos diferentes**: o primeiro contexto "novo" do nosso código pode ser usado dentro dos testes. No entanto, se alguém quiser testar algo novo com nossa função, a pessoa pode injetar suas próprias dependências.

### Mas e o mock? Ouvi falar que precisa disso para trabalhar com injeção de dependência e que também é do demonho

Vamos falar mais sobre mocks depois (e não é do demonho). Você mocka para substituir coisas reais que você injeta com uma versão falsa que você pode controlar e examinar nos seus testes. No entanto, no nosso caso a biblioteca padrão já tinha algo pronto para usarmos.

### A biblioteca padrão do Go é muito boa, leve um tempo para estudá-la

Ao termos familiaridade com a interface `io.Writer`, somos capazes de usar `bytes.Buffer` no nosso teste como nosso `Writer` para que depois possamos usar outros `Writer` da biblioteca padrão para usar na nossa função em uma aplicação de linha de comando ou em um servidor web.

Quanto mais familiar você for com a biblioteca padrão, mais vai ver essas interfaces de propósito geral que você pode reutilizar no seu próprio código para tornar o software reutilizável em vários contextos diferentes.

Esse exemplo teve grande influência de um capítulo de [A Linguagem de Programação Go](https://www.amazon.com.br/dp/8575225464/ref=cm_sw_r_tw_dp_U_x_0HIqDbYP7VSN5). Logo, se gostou, vá adquiri-lo!
