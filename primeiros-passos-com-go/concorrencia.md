# Concorrência

[**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/concurrency)

A questão é a seguinte: um colega escreveu uma função, `VerificaWebsites`, que verifica o status de uma lista de URLs.

```go
package concorrencia

type VerificadorWebsite func(string) bool

func VerificaWebsites(vw VerificadorWebsite, urls []string) map[string]bool {
    resultados := make(map[string]bool)

    for _, url := range urls {
        resultados[url] = vw(url)
    }

    return resultados
}
```

Ela retorna um map de cada URL verificado com um valor booleano - `true` para uma boa resposta, `false` para uma resposta ruim.

Você também tem que passar um `VerificadorWebsite` como parâmetro, que leva um URL e retorna um boleano. Isso é usado pela função que verifica todos os websites.

Usando a [injeção de dependência](dependency-injection.md), conseguimos testar a função sem fazer chamadas HTTP de verdade, tornando o teste seguro e rápido.

Aqui está o teste que escreveram:

```go
package concurrency

import (
    "reflect"
    "testing"
)

ffunc mockVerificadorWebsite(url string) bool {
    if url == "waat://furhurterwe.geds" {
        return false
    }
    return true
}

func TestVerificaWebsites(t *testing.T) {
    websites := []string{
        "http://google.com",
        "http://blog.gypsydave5.com",
        "waat://furhurterwe.geds",
    }

    esperado := map[string]bool{
        "http://google.com":          true,
        "http://blog.gypsydave5.com": true,
        "waat://furhurterwe.geds":    false,
    }

    resultado := VerificaWebsites(mockVerificadorWebsite, websites)

    if !reflect.DeepEqual(esperado, resultado) {
        t.Fatalf("esperado %v, resultado %v", esperado, resultado)
    }
}
```

A função que está em produção está sendo usada para verificar centenas de websites. Só que seu colega começou a reclamar que está lento demais, e pediram sua ajuda para melhorar a velocidade dele.

## Escreva o teste primeiro

Vamos usar um teste de benchmark para testar a velocidade de `VerificaWebsites` para que possamos ver o efeito das nossas alterações.

```go
package concorrencia

import (
    "testing"
    "time"
)

func slowStubVerificadorWebsite(_ string) bool {
    time.Sleep(20 * time.Millisecond)
    return true
}

func BenchmarkVerificaWebsites(b *testing.B) {
    urls := make([]string, 100)
    for i := 0; i < len(urls); i++ {
        urls[i] = "uma url"
    }

    for i := 0; i < b.N; i++ {
        VerificaWebsites(slowStubVerificadorWebsite, urls)
    }
}
```

O benchmark testa `VerificaWebsites` usando um slice de 100 URLs e usa uma nova implementação falsa de `VerificadorWebsite`. `slowStubVerificadorWebsite` é intencionalmente lento. Ele usa um `time.Sleep` para esperar exatamente 20 milissegundos e então retorna verdadeiro.

Quando executamos o benchmark com `go test -bench=.` (ou, se estiver no Powershell do WIndows, `go test -bench="."`):

```bash
pkg: github.com/larien/learn-go-with-tests/concorrencia/v1
BenchmarkVerificaWebsites-4               1        2249228637 ns/op
PASS
ok      github.com/larien/learn-go-with-tests/concorrencia/v1        2.268s
```

`VerificaWebsites` teve uma marca de 2249228637 nanosegundos - pouco mais de dois segundos.

Vamos torná-lo mais rápido.

### Escreva código o suficiente para fazer o teste passar

Agora finalmente podemos falar sobre concorrência que, apenas para fins dessa situação, significa "fazer mais do que uma coisa ao mesmo tempo". Isso é algo que fazemos naturalmente todo dia.

Por exemplo, hoje de manhã fiz uma xícara de chá. Coloquei a chaleira no foco e, enquanto esperava a água ferver, tirei o leite da geladeira, tirei o chá do armário, encontrei minha xícara favorita, coloquei o saquinho do chá e, quando a chaleira ferveu a água, coloquei a água na xícara.

O que eu _não fiz_ foi colocar a chaleira no fogo e então ficar sem fazer nada só esperando a chaleira ferver a água, para depois fazer todo o restante quando a água tivesse fervido.

Se conseguir entender por que é mais rápido fazer chá da primeira forma, então você é capaz de entender como vamos tornar o `VerificaWebsites` mais rápido. Ao invés de esperar por um website responder antes de enviar uma requisição para o próximo website, vamos dizer para nosso computador fazer a próxima requisição enquanto espera pela primeira.

Normalmente, em Go, quando chamamos uma função `fazAlgumaCoisa()`, esperamos que ela retorne alguma coisa (mesmo se não tiver valor para retornar, ainda esperamos que ela termine). Chamamos essa operação de _bloqueante_ - espera algo acabar para terminar seu trabalho. Uma operação que não bloqueia no Go vai rodar em um _processo_ separado, chamado de _goroutine_. Pense no processo como uma leitura de uma página de código Go de cima para baixo, 'entrando' em cada função quando é chamado para ler o que essa página faz. Quando um processo separado começa, é como se outro leitor começasse a ler o interior da função, deixando o leitor original continuar lendo a página.

Para dizer ao Go começar uma nova goroutine, transformamos a chamada de função em uma declaração `go` colocando a palavra-chave `go` na frente da função: `go fazAlgumaCoisa()`.

```go
package concurrency

type VerificadorWebsite func(string) bool

func VerificaWebsites(vw VerificadorWebsite, urls []string) map[string]bool {
    resultados := make(map[string]bool)

    for _, url := range urls {
        go func() {
            resultados[url] = vw(url)
        }()
    }

    return resultados
}
```

Já que a única forma de começar uma goroutine é colocar `go` na frente da chamada de função, costumamos usar _funções anônimas_ quando queremos iniciar uma goroutine. Uma função anônima literal é bem parecida com uma declaração de função normal, mas (obviamente) sem um nome. Você ṕde ver uma acima no corpo do laço `for`.

Funções anônimas têm várias funcionalidades que as torna útil, duas das quais estamos usando acima. Primeiramente, elas podem ser executadas assim que fazemos sua declaração - que é o `()` no final da função anônima. Em segundo lugar, elas mantém acesso ao escopo léxico em que são definidas - todas as variáveis que estão disponíveis no ponto em que a função anônima é declarada também estão variáveis no corpo da função.

O corpo da função anônima acima é quase o mesmo da função no laço utilizada anteriormente. A única diferença é que cada iteração do loop vai iniciar uma nova goroutine, concorrente com o processo atual (a função `VerificadorWebsite`), e cada uma vai adicionar seu resultado ao map de resultados.

```bash
--- FAIL: TestVerificaWebsites (0.00s)
        VerificaWebsites_test.go:31: esperado map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], resultado map[]
FAIL
exit status 1
FAIL    github.com/larien/learn-go-with-tests/concorrencia/v2        0.010s
```

### Uma breve visita ao universo paralelo...

Você pode não ter obtido esse resultado. Você pode obter uma mesnagem de pânico, que vamos falar sobre em breve. Não se preocupe se isso aparecer para você, basta você executar o teste até você _de fato_ receber o resultado acima. Ou faça de conta que você recebeu. Escolha sua. Boas vindas à concorrência: quando não for trabalhada da forma correta, é difícil prever o que vai acontecer. Não se preocupe, é por isso que estamos escrevendo testes: para nos ajudar a saber quando estamos trabalhando com concorrência de forma previsível.

### ... e estamos de volta.

Acabou que os testes originais do `VerificadorWebsite` agora está devolvendo um map vazio. O que deu de errado?

Nenhuma das goroutines que nosso loop `for` iniciou teve tempo de adicionar seu resultado ao map `resultados`; a função `VerificadorWebsite` é rápida demais para eles, e por isso retorna o map vazio.

Para consertar isso, podemos apenas esperar enquanto todas as goroutines fazem seu trabalho, para depois retornar. Dois segundos devem servir, certo?

```go
package concurrency

import "time"

type VerificadorWebsite func(string) bool

func VerificaWebsites(vw VerificadorWebsite, urls []string) map[string]bool {
    resultados := make(map[string]bool)

    for _, url := range urls {
        go func() {
            resultados[url] = vw(url)
        }()
    }

    time.Sleep(2 * time.Second)

    return resultados
}

```

Agora, quando os testes forem executados, você vai ver (ou não - leia a mensagem no início do tópico):

```bash
--- FAIL: TestVerificaWebsites (0.00s)
        VerificaWebsites_test.go:31: esperado map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], resultado map[waat://furhurterwe.geds:false]
FAIL
exit status 1
FAIL    github.com/larien/learn-go-with-tests/concorrencia/v1        0.010s
```

Isso não é muito bom - por que só um resultado? Podemos arrumar isso aumentando o tempo de espera - pode tentar se preferir. Não vai funcionar. O problema aqui é que a variável `url` é reutilizada para cada iteração do laço `for` - ele recebe um valor novo de `urls` a cada vez. Mas cada uma das goroutines tem uma referência para a variável `url` - eles não têm sua própria cópia independente. Logo, _todas_ estão escrevendo o valor que `url` tem no final da iteração - o último URL. E é por isso que o resultado que obtemos é a última URL.

Para consertar isso:

```go
package concorrencia

import (
    "time"
)

type VerificadorWebsite func(string) bool

func VerificaWebsites(vw VerificadorWebsite, urls []string) map[string]bool {
    resultados := make(map[string]bool)

    for _, url := range urls {
        go func(u string) {
            resultados[u] = vw(u)
        }(url)
    }

    time.Sleep(2 * time.Second)

    return resultados
}
```

Ao passar cada função anônima como parâmetro para a URL - como `u` - e chamar a função anônima com `url` como argumento, nos certificamos de que o valor de `u` está fixado como o valor de `url` para cada iteração do laço de `url` e não pode ser modificado.

Agora, se você tiver sorte, vai obter:

```bash
PASS
ok      github.com/larien/learn-go-with-tests/concorrencia/v1        2.012s
```

No entanto, se não tiver sorte (isso é mais provável se estiver rodando o código com o benchmark, já que haverá mais tentativas):

```bash
fatal error: concurrent map writes

goroutine 37 [running]:
runtime.throw(0x6d74f3, 0x15)
    /usr/local/go/src/runtime/panic.go:608 +0x72 fp=0xc000034718 sp=0xc0000346e8 pc=0x42d4e2
runtime.mapassign_faststr(0x67dbe0, 0xc000082660, 0x6d33cb, 0x7, 0x0)
    /usr/local/go/src/runtime/map_faststr.go:275 +0x3bf fp=0xc000034780 sp=0xc000034718 pc=0x4139ff
github.com/larien/learn-go-with-tests/concorrencia/v2.VerificaWebsites.func1(0x6e6580, 0xc000082660, 0x6d33cb, 0x7)
    /home/larien/go/src/github.com/larien/learn-go-with-tests/concorrencia/v2/VerificaWebsites.go:17 +0x7f fp=0xc0000347c0 sp=0xc000034780 pc=0x64035f
runtime.goexit()
    /usr/local/go/src/runtime/asm_amd64.s:1333 +0x1 fp=0xc0000347c8 sp=0xc0000347c0 pc=0x45c661
created by github.com/larien/learn-go-with-tests/concorrencia/v2.VerificaWebsites
	/home/larien/go/src/github.com/larien/learn-go-with-tests/concorrencia/v2/VerificaWebsites.go:16 +0xa9

        ... e mais um monte de linhas assustadoras ...
```

Isso pode ser enorme e assustador, mas tudo o que precisamos fazer é respirar com calma e ler o stacktrace: `fatal error: concurrent map writes` (erro fatal: escrita concorrente no map). Às vezes, quando executamos nossos testes, duas das goroutines escrevem no map `resultados` ao mesmo tempo. Maps em Go não gostam quando mais de uma coisa tenta escrever algo neles ao mesmo tempo, então o `erro fatal` é gerado.

Essa é uma _condição de corrida_, um bug que aparece quando a saída do nosso software depende do timing e da sequência de eventos que não temos controle sobre. Por não termos controle exato sobre quando cada goroutine escreve no map `resultados`, ficamos vulneráveis à situação de duas goroutines escreverem nele ao mesmo tempo.

O Go nos ajuda a encontrar condições de corrida com seu [_detector de corrida_](https://blog.golang.org/race-detector) nativo. Para habilitar essa funcionalidade, execute os testes com a flag `race`: `go test -race`.

Você deve ver uma saída parecida com essa:

```bash
==================
WARNING: DATA RACE
Write at 0x00c000120089 by goroutine 6:
  reflect.typedmemmove()
      /usr/local/go/src/runtime/mbarrier.go:177 +0x0
  reflect.Value.MapIndex()
      /usr/local/go/src/reflect/value.go:1124 +0x2ae
  reflect.deepValueEqual()
      /usr/local/go/src/reflect/deepequal.go:118 +0x13be
  reflect.DeepEqual()
      /usr/local/go/src/reflect/deepequal.go:196 +0x2f0
  github.com/larien/learn-go-with-tests/concorrencia/v2.TestVerificaWebsites()
      /home/larien/go/src/github.com/larien/learn-go-with-tests/concorrencia/v2/VerificaWebsites_test.go:30 +0x1ad
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:827 +0x162

Previous write at 0x00c000120089 by goroutine 8:
  github.com/larien/learn-go-with-tests/concorrencia/v2.VerificaWebsites.func1()
      /home/larien/go/src/github.com/larien/learn-go-with-tests/concorrencia/v2/VerificaWebsites.go:17 +0x97

Goroutine 6 (running) created at:
  testing.(*T).Run()
      /usr/local/go/src/testing/testing.go:878 +0x659
  testing.runTests.func1()
      /usr/local/go/src/testing/testing.go:1119 +0xa8
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:827 +0x162
  testing.runTests()
      /usr/local/go/src/testing/testing.go:1117 +0x4ee
  testing.(*M).Run()
      /usr/local/go/src/testing/testing.go:1034 +0x2ee
  main.main()
      _testmain.go:44 +0x221

Goroutine 8 (finished) created at:
  github.com/larien/learn-go-with-tests/concorrencia/v2.VerificaWebsites()
      /home/larien/go/src/github.com/larien/learn-go-with-tests/concorrencia/v2/VerificaWebsites.go:16 +0xb2
  github.com/larien/learn-go-with-tests/concorrencia/v2.TestVerificaWebsites()
      /home/larien/go/src/github.com/larien/learn-go-with-tests/concorrencia/v2/VerificaWebsites_test.go:28 +0x17f
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:827 +0x162
==================
```

Os detalhes ainda assim são bem difíceis de serem lidos - mas o `WARNING: DATA RACE` (CUIDADO: CONDIÇÃO DE CORRIDA) é bem claro. Lendo o corpo do erro podemos ver duas goroutines diferentes performando escritas em um map:

`Write at 0x00c000120089 by goroutine 6:`

está escrevendo no mesmo bloco de memória que:

`Previous write at 0x00c000120089 by goroutine 8:`

Além disso, conseguimos ver a linha de código onde a escrita está acontecendo:

`/home/larien/go/src/github.com/larien/learn-go-with-tests/concorrencia/v2/VerificaWebsites.go:17 +0x97`

e a linha de código onde as goroutines 6 e 7 foram iniciadas:

`/home/larien/go/src/github.com/larien/learn-go-with-tests/concorrencia/v2/VerificaWebsites.go:16 +0xb2`

Tudo o que você precisa saber está impresso no seu terminal - tudo o que você tem que fazer é ser paciente o bastante para lê-lo.

### Canais

Podemos resolver essa condição de corrida coordenando nossas goroutines usando _canais_. Canais são uma estrutura de dados em Go que pode receber e enviar valores. Essas operações, junto de seus detalhes, permitem a comunicação entre processos diferentes.

Nesse caso, queremos pensar sobre a comunicação entre o processo pai e cada uma das goroutines criadas por ele de orma que façam o trabalho de executar a função `VerificadorWebsite` com a URL.

```go
package concurrency

type WebsiteChecker func(string) bool
type result struct {
    string
    bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)
    resultChannel := make(chan result)

    for _, url := range urls {
        go func(u string) {
            resultChannel <- result{u, wc(u)}
        }(url)
    }

    for i := 0; i < len(urls); i++ {
        result := <-resultChannel
        results[result.string] = result.bool
    }

    return results
}
```

Alongside the `results` map we now have a `resultChannel`, which we `make` in the same way. `chan result` is the type of the channel - a channel of `result`. The new type, `result` has been made to associate the return value of the `WebsiteChecker` with the url being checked - it's a struct of `string` and `bool`. As we don't need either value to be named, each of them is anonymous within the struct; this can be useful in when it's hard to know what to name a value.

Now when we iterate over the urls, instead of writing to the `map` directly we're sending a `result` struct for each call to `wc` to the `resultChannel` with a _send statement_. This uses the `<-` operator, taking a channel on the left and a value on the right:

```go
// Send statement
resultChannel <- result{u, wc(u)}
```

The next `for` loop iterates once for each of the urls. Inside we're using a _receive expression_, which assigns a value received from a channel to a variable. This also uses the `<-` operator, but with the two operands now reversed: the channel is now on the right and the variable that we're assigning to is on the left:

```go
// Receive expression
result := <-resultChannel
```

We then use the `result` received to update the map.

By sending the results into a channel, we can control the timing of each write into the results map, ensuring that it happens one at a time. Although each of the calls of `wc`, and each send to the result channel, is happening in parallel inside its own process, each of the results is being dealt with one at a time as we take values out of the result channel with the receive expression.

We have paralellized the part of the code that we wanted to make faster, while making sure that the part that cannot happen in parallel still happens linearly. And we have communicated across the multiple processes involved by using channels.

When we run the benchmark:

```bash
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v2
BenchmarkCheckWebsites-8             100          23406615 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        2.377s
```

23406615 nanoseconds - 0.023 seconds, about one hundred times as fast as original function. A great success.

## Wrapping up

This exercise has been a little lighter on the TDD than usual. In a way we've been taking part in one long refactoring of the `CheckWebsites` function; the inputs and outputs never changed, it just got faster. But the tests we had in place, as well as the benchmark we wrote, allowed us to refactor `CheckWebsites` in a way that maintained confidence that the software was still working, while demonstrating that it had actually become faster.

In making it faster we learned about

-   _goroutines_, the basic unit of concurrency in Go, which let us check more

    than one website at the same time.

-   _anonymous functions_, which we used to start each of the concurrent processes

    that check websites.

-   _channels_, to help organize and control the communication between the

    different processes, allowing us to avoid a _race condition_ bug.

-   _the race detector_ which helped us debug problems with concurrent code

### Make it fast

One formulation of an agile way of building software, often misattributed to Kent Beck, is:

> [Make it work, make it right, make it fast](http://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast)

Where 'work' is making the tests pass, 'right' is refactoring the code, and 'fast' is optimizing the code to make it, for example, run quickly. We can only 'make it fast' once we've made it work and made it right. We were lucky that the code we were given was already demonstrated to be working, and didn't need to be refactored. We should never try to 'make it fast' before the other two steps have been performed because

> [Premature optimization is the root of all evil](http://wiki.c2.com/?PrematureOptimization) -- Donald Knuth
