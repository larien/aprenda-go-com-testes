# Concorrência

[**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/master/primeiros-passos-com-go/concorrencia)

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

Usando a [injeção de dependência](../injecao-de-dependencia/injecao-de-dependencia.md), conseguimos testar a função sem fazer chamadas HTTP de verdade, tornando o teste seguro e rápido.

Aqui está o teste que escreveram:

```go
package concurrency

import (
    "reflect"
    "testing"
)

func mockVerificadorWebsite(url string) bool {
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

Quando executamos o benchmark com `go test -bench=.` (ou, se estiver no Powershell do Windows, `go test -bench="."`):

```bash
pkg: github.com/larien/aprenda-go-com-testes/concorrencia/v1
BenchmarkVerificaWebsites-4               1        2249228637 ns/op
PASS
ok      github.com/larien/aprenda-go-com-testes/concorrencia/v1        2.268s
```

`VerificaWebsites` teve uma marca de 2249228637 nanosegundos - pouco mais de dois segundos.

Vamos torná-lo mais rápido.

### Escreva código o suficiente para fazer o teste passar

Agora finalmente podemos falar sobre concorrência que, apenas para fins dessa situação, significa "fazer mais do que uma coisa ao mesmo tempo". Isso é algo que fazemos naturalmente todo dia.

Por exemplo, hoje de manhã fiz uma xícara de chá. Coloquei a chaleira no fogo e, enquanto esperava a água ferver, tirei o leite da geladeira, tirei o chá do armário, encontrei minha xícara favorita, coloquei o saquinho do chá e, quando a chaleira ferveu a água, coloquei a água na xícara.

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

Já que a única forma de começar uma goroutine é colocar `go` na frente da chamada de função, costumamos usar _funções anônimas_ quando queremos iniciar uma goroutine. Uma função anônima literal é bem parecida com uma declaração de função normal, mas (obviamente) sem um nome. Você pode ver uma acima no corpo do laço `for`.

Funções anônimas têm várias funcionalidades que as tornam úteis, duas das quais estamos usando acima. Primeiramente, elas podem ser executadas assim que fazemos sua declaração - que é o `()` no final da função anônima. Em segundo lugar, elas mantém acesso ao escopo léxico em que são definidas - todas as variáveis que estão disponíveis no ponto em que a função anônima é declarada também estão disponíveis no corpo da função.

O corpo da função anônima acima é quase o mesmo da função no laço utilizada anteriormente. A única diferença é que cada iteração do loop vai iniciar uma nova goroutine, concorrente com o processo atual (a função `VerificadorWebsite`), e cada uma vai adicionar seu resultado ao map de resultados.

```bash
--- FAIL: TestVerificaWebsites (0.00s)
        VerificaWebsites_test.go:31: esperado map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], resultado map[]
FAIL
exit status 1
FAIL    github.com/larien/aprenda-go-com-testes/concorrencia/v2        0.010s
```

### Uma breve visita ao universo paralelo...

Você pode não ter obtido esse resultado. Você pode obter uma mensagem de pânico, que vamos falar sobre em breve. Não se preocupe se isso aparecer para você, basta você executar o teste até você _de fato_ receber o resultado acima. Ou faça de conta que você recebeu. Escolha sua. Boas vindas à concorrência: quando não for trabalhada da forma correta, é difícil prever o que vai acontecer. Não se preocupe, é por isso que estamos escrevendo testes: para nos ajudar a saber quando estamos trabalhando com concorrência de forma previsível.

### ... e estamos de volta.

Acabou que os testes originais do `VerificadorWebsite` agora estão devolvendo um map vazio. O que deu de errado?

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
FAIL    github.com/larien/aprenda-go-com-testes/concorrencia/v1        0.010s
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
ok      github.com/larien/aprenda-go-com-testes/concorrencia/v1        2.012s
```

No entanto, se não tiver sorte (isso é mais provável se estiver rodando o código com o benchmark, já que haverá mais tentativas):

```bash
fatal error: concurrent map writes

goroutine 37 [running]:
runtime.throw(0x6d74f3, 0x15)
    /usr/local/go/src/runtime/panic.go:608 +0x72 fp=0xc000034718 sp=0xc0000346e8 pc=0x42d4e2
runtime.mapassign_faststr(0x67dbe0, 0xc000082660, 0x6d33cb, 0x7, 0x0)
    /usr/local/go/src/runtime/map_faststr.go:275 +0x3bf fp=0xc000034780 sp=0xc000034718 pc=0x4139ff
github.com/larien/aprenda-go-com-testes/concorrencia/v2.VerificaWebsites.func1(0x6e6580, 0xc000082660, 0x6d33cb, 0x7)
    /home/larien/go/src/github.com/larien/aprenda-go-com-testes/concorrencia/v2/VerificaWebsites.go:17 +0x7f fp=0xc0000347c0 sp=0xc000034780 pc=0x64035f
runtime.goexit()
    /usr/local/go/src/runtime/asm_amd64.s:1333 +0x1 fp=0xc0000347c8 sp=0xc0000347c0 pc=0x45c661
created by github.com/larien/aprenda-go-com-testes/concorrencia/v2.VerificaWebsites
	/home/larien/go/src/github.com/larien/aprenda-go-com-testes/concorrencia/v2/VerificaWebsites.go:16 +0xa9

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
  github.com/larien/aprenda-go-com-testes/concorrencia/v2.TestVerificaWebsites()
      /home/larien/go/src/github.com/larien/aprenda-go-com-testes/concorrencia/v2/VerificaWebsites_test.go:30 +0x1ad
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:827 +0x162

Previous write at 0x00c000120089 by goroutine 8:
  github.com/larien/aprenda-go-com-testes/concorrencia/v2.VerificaWebsites.func1()
      /home/larien/go/src/github.com/larien/aprenda-go-com-testes/concorrencia/v2/VerificaWebsites.go:17 +0x97

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
  github.com/larien/aprenda-go-com-testes/concorrencia/v2.VerificaWebsites()
      /home/larien/go/src/github.com/larien/aprenda-go-com-testes/concorrencia/v2/VerificaWebsites.go:16 +0xb2
  github.com/larien/aprenda-go-com-testes/concorrencia/v2.TestVerificaWebsites()
      /home/larien/go/src/github.com/larien/aprenda-go-com-testes/concorrencia/v2/VerificaWebsites_test.go:28 +0x17f
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:827 +0x162
==================
```

Os detalhes ainda assim são bem difíceis de serem lidos - mas o `WARNING: DATA RACE` (CUIDADO: CONDIÇÃO DE CORRIDA) é bem claro. Lendo o corpo do erro podemos ver duas goroutines diferentes performando escritas em um map:

`Write at 0x00c000120089 by goroutine 6:`

está escrevendo no mesmo bloco de memória que:

`Previous write at 0x00c000120089 by goroutine 8:`

Além disso, conseguimos ver a linha de código onde a escrita está acontecendo:

`/home/larien/go/src/github.com/larien/aprenda-go-com-testes/concorrencia/v2/VerificaWebsites.go:17 +0x97`

e a linha de código onde as goroutines 6 e 7 foram iniciadas:

`/home/larien/go/src/github.com/larien/aprenda-go-com-testes/concorrencia/v2/VerificaWebsites.go:16 +0xb2`

Tudo o que você precisa saber está impresso no seu terminal - tudo o que você tem que fazer é ser paciente o bastante para lê-lo.

### Canais

Podemos resolver essa condição de corrida coordenando nossas goroutines usando _canais_. Canais são uma estrutura de dados em Go que pode receber e enviar valores. Essas operações, junto de seus detalhes, permitem a comunicação entre processos diferentes.

Nesse caso, queremos pensar sobre a comunicação entre o processo pai e cada uma das goroutines criadas por ele de forma que façam o trabalho de executar a função `VerificadorWebsite` com a URL.

```go
package concurrency

type VerificadorWebsite func(string) bool
type resultado struct {
    string
    bool
}

func VerificaWebsites(vw VerificadorWebsite, urls []string) map[string]bool {
    resultados := make(map[string]bool)
    canalResultado := make(chan resultado)

    for _, url := range urls {
        go func(u string) {
            canalResultado <- resultado{u, vw(u)}
        }(url)
    }

    for i := 0; i < len(urls); i++ {
        resultado := <-canalResultado
        resultados[resultado.string] = resultado.bool
    }

    return resultados
}
```

Junto do map `resultados`, agora temos um `canalResultados`, que criamos da mesma forma usando `make`. O `chan resultado` é o tipo do canal - um canal de `resultado`. O tipo novo, `resultado`, foi criado para associar o retorno de `VerificadorWebsite` com a URL sendo verificada - é uma estrutura que contém uma `string` e um `bool`. Já que não precisamos que nenhum valor tenha um nome, cada um deles é anônimo dentro da struct; isso pode ser útil quando for difícil saber que nome dar a um valor.

Agora que iteramos pelas URLs, ao invés de escrever no `map` diretamente, enviamos uma struct `resultado` para cada chamada de `vw` para o `canalResultado` com uma _sintaxe de envio_. Essa sintaxe usa o operador `<-`, usando um canal à esquerda e um valor à direita:

```go
// Sintaxe de envio
canalResultado <- resultado{u, vw(u)}
```

O próximo laço `for` itera uma vez sobre cada uma das URLs. Dentro, estamos usando uma _expressão de recebimento_, que atribui um valor recebido por um canal a uma variável. Essa expressão também usa o operador `<-`, mas com os dois operandos ao posições invertidas: o canal agora fica à direita e a variável que está recebendo o valor dele fica à esquerda:

```go
// Expressão recebida
resultado := <-canalResultado
```

E depois usamos o `resultado` recebido para atualizar o map.

Ao enviar os resultados para um canal, podemos controlar o timing de cada escrita dentro do map `resultados`, garantindo que só aconteça uma por vez. Apesar de cada uma das chamadas de `vw` e cada envio ao canal resultado estar acontecendo em paralelo dentro de seu próprio processo, cada resultado está sendo resolvido de cada vez enquanto tiramos o valor do canal resultado com a expressão recebida.

Paralelizamos um pedaço do código que queríamos tornar mais rápida, enquanto mantivemos a parte que não pode acontecer em paralelo ainda acontecendo linearmente. E comunicamos diversos processos envolvidos utilizando canais.

Agora podemos executar o benchmark:

```bash
pkg: github.com/larien/aprenda-go-com-testes/concorrencia/v3
BenchmarkVerificaWebsites-8             100          23406615 ns/op
PASS
ok      github.com/larien/aprenda-go-com-testes/concorrencia/v3        2.377s
```

23406615 nanossegundos - 0.023 segundos, cerca de 100 vezes mais rápida que a função original. Um sucesso enorme.

## Resumo

Esse exercício foi um pouco mais leve na parte do TDD que o restante. Levamos um bom tempo refatorando a função `VerificaWebsites`; as entradas e saídas não mudaram, ela apenas ficou mais rápida. Mas, com os testes que já tinhamos escrito, assim como com o benchmark que escrevemos, fomos capazes de refatorar o `VerificaWebsites` de forma que mantivéssemos a confiança de que o software ainda estava funcionando, enquanto demonstramos que ela realmente havia ficado mais rápida.

Tornando as coisas mais rápidas, aprendemos sobre:

-   _goroutines_, a unidade básica de concorrência em Go, que nos permite verificar mais do que um site ao mesmo tempo.

-   _funções anônimas_, que usamos para iniciar cada um dos processos concorrentes que verificam os sites.

-   _canais_, para nos ajudar a organizar e controlar a comunicação entre diferentes processos, nos permitindo evitar um bug de _condição de corrida_.

-   _o detector de corrida_, que nos ajudou a desvendar problemas com código concorrente.

### Torne-o rápido

Uma formulação da forma ágil de desenvolver software, erroneamente atribuida a Kent Beck, é:

> [Faça funcionar, faça da forma certa, torne-o rápido](http://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast) (em inglês)

Onde 'funcionar' é fazer os testes passarem, 'forma certa' é refatorar o código e 'tornar rápido' é otimizar o código para, por exemplo, tornar sua execução rápida. Só podemos 'torná-lo rápido' quando fizermos funcionar da forma certa. Tivemos sorte que o código que estudamos já estava funcionando e não precisava ser refatorado. Nunca devemos tentar 'torná-lo rápido' antes das outras duas etapas terem sido feitas, porque:

> [Otimização prematura é a raiz de todo o mal](http://wiki.c2.com/?PrematureOptimization) -- Donald Knuth
