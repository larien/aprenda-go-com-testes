# Mocking

[**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/mocks)

Te pediram para criar um programa que conta a partir de 3, imprimindo cada número em uma linha nova (com um segundo intervalo entre cada uma) e quando chega a zero, imprime "Vai!" e sai.

```text
3
2
1
Vai!
```

Vamos resolver isso escrevendo uma função chamada `Contagem` que vamos colocar dentro de um programa `main` e se parecer com algo assim:

```go
package main

func main() {
    Contagem()
}
```

Apesar de ser um programa simples, para testá-lo completamente vamos precisar, como de costume, de uma abordagem _iterativa_ e _orientada a testes_.

Mas o que quero dizer com iterativa? Precisamos ter certeza de que tomamos os menores passos que pudermos para ter um _software_ útil.

Não queremos passar muito tempo com código que vai funcionar hora ou outra após alguma implementação mirabolante, porque é assim que os desenvolvedores caem em armadilhas. **É importante ser capaz de dividir os requerimentos da menor forma que conseguir para você ter um** _**software funcionando**_**.**

Podemos separar essa tarefa da seguinte forma:

-   Imprimir 3
-   Imprimir de 3 para Vai!
-   Esperar um segundo entre cada linha

## Escreva o teste primeiro

Nosso software precisa imprimir para a saída. Vimos como podemos usar a injeção de dependência para facilitar nosso teste na [seção anterior](https://github.com/larien/learn-go-with-tests/tree/master/primeiros-passos-com-go/injecao-de-dependencia.md).

```go
func TestContagem(t *testing.T) {
    buffer := &bytes.Buffer{}

    Contagem(buffer)

    resultado := buffer.String()
    esperado := "3"

    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }
}
```

Se tiver dúvidas sobre o `buffer`, leia a [seção anterior](https://github.com/larien/learn-go-with-tests/tree/master/primeiros-passos-com-go/injecao-de-dependencia.md) novamente.

Sabemos que nossa função `Contagem` precisa escrever dados em algum lugar e o `io.Writer` é a forma de capturarmos essa saída como uma interface em Go.

-   Na `main`, vamos enviar o `os.Stdout` como parâmetro para nossos usuários verem a contagem regressiva impressa no terminal.
-   No teste, vamos enviar o `bytes.Buffer` como parâmetro para que nossos testes possam capturar que dado está sendo gerado.

## Execute o teste

`./contagem_test.go:11:2: undefined: Contagem`

`indefinido: Contagem`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Defina `Contagem`:

```go
func Contagem() {}
```

Tente novamente:

```go
./contagem_test.go:11:11: too many arguments in call to Countdown
    have (*bytes.Buffer)
    want ()
```

`argumentos demais na chamada para Contagem`

O compilador está te dizendo como a assinatura da função deve ser, então é só atualizá-la.

```go
func Contagem(saida *bytes.Buffer) {}
```

`contagem_test.go:17: resultado '', esperado '3'`

Perfeito!

## Escreva código o suficiente para fazer o teste passar

```go
func Contagem(saida *bytes.Buffer) {
    fmt.Fprint(saida, "3")
}
```

Estamos usando `fmt.Fprint`, o que significa que ele recebe um `io.Writer` (como `*bytes.Buffer`) e envia uma `string` para ele. O teste deve passar.

## Refatoração

Agora sabemos que, apesar do `*bytes.Buffer` funcionar, seria melhor ter uma interface de propósito geral ao invés disso.

```go
func Contagem(saida io.Writer) {
    fmt.Fprint(saida, "3")
}
```

Execute os testes novamente e eles devem passar.

Só para finalizar, vamos colocar nossa função dentro da `main` para que possamos executar o software para nos assegurarmos de que estamos progredindo.

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func Contagem(saida io.Writer) {
	fmt.Fprint(saida, "3")
}

func main() {
	Contagem(os.Stdout)
}
```

Execute o programa e surpreenda-se com seu trabalho.

Apesar de parecer simples, essa é a abordagem que recomendo para qualquer projeto. **Escolher uma pequena parte da funcionalidade e fazê-la funcionar do começo ao fim com apoio de testes.**

Depois, precisamos fazer o software imprimir 2, 1 e então "Vai!".

## Escreva o teste primeiro

Após investirmos tempo e esforço para fazer o principal funcionar, podemos iterar nossa solução com segurança e de forma simples. Não vamos mais precisar para parar e executar o programa novamente para ter confiança de que ele está funcionando, desde que a lógica esteja testada.

```go
func TestContagem(t *testing.T) {
    buffer := &bytes.Buffer{}

    Contagem(buffer)

    resultado := buffer.String()
    esperado := `3
2
1
Vai!`
    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }
}
```

A sintaxe de aspas simples é outra forma de criar uma `string`, mas te permite colocar coisas como linhas novas, o que é perfeito para nosso teste.

## Execute o teste

```bash
contagem_test.go:21: resultado '3', esperado '3
        2
        1
        Vai!'
```

## Escreva código o suficiente para fazer o teste passar

```go
func Contagem(saida io.Writer) {
    for i := 3; i > 0; i-- {
        fmt.Fprintln(saida, i)
    }
    fmt.Fprint(saida, "Go!")
}
```

Usamos um laço `for` fazendo contagem regressiva com `i--` e depois `fmt.Fprintln` para imprimir a `saida` com nosso número seguro por um caracter de nova linha. Finalmente, usamos o `fmt.Fprint` para enviar "Vai!" no final.

## Refatoração

Não há muito para refatorar além de transformar alguns valores mágicos em constantes com nomes descritivos.

```go
const ultimaPalavra = "Go!"
const inicioContagem = 3

func Contagem(saida io.Writer) {
    for i := inicioContagem; i > 0; i-- {
        fmt.Fprintln(saida, i)
    }
    fmt.Fprint(saida, ultimaPalavra)
}
```

Se executar o programa agora, você deve obter a saída de sejada, mas não tem uma contagem regressiva dramática com as pausas de 1 segundo.

Go te permite obter isso com `time.Sleep`. Tente adicionar essa função ao seu código.

```go
func Contagem(saida io.Writer) {
    for i := inicioContagem; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(saida, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(saida, ultimaPalavra)
}
```

Se você executar o programa, ele funciona conforme esperado.

## Mock

Os testes ainda vão passar e o software funciona como planejado, mas temos alguns problemas:

-   Nossos testes levam 4 segundos para rodar.
    -   Todo conteúdo gerado sobre desenvolvimento de software enfatiza a importância de loops de feedback rápidos.
    -   **Testes lentos arruinam a produtividade do desenvolvedor**.
    -   Imagine se os requerimentos ficam mais sofisticados, gerando a necessidade de mais testes. É viável adicionar 4s para cada teste novo de `Contagem`?
-   Não testamos uma propriedade importante da nossa função.

Temos uma dependência no `Sleep` que precisamos extrair para podermos controlá-la nos nossos testes.

Se conseguirmos _mockar_ o `time.Sleep`, podemos usar a _injeção de dependências_ para usá-lo ao invés de um `time.Sleep` "de verdade", e então podemos **verificar as chamadas** para certificar de que estão corretas.

## Escreva o teste primeiro

Vamos definir nossa dependência como uma interface. Isso nos permite usar um Sleeper _de verdade_ em `main` e um _sleeper spy_ nos nossos testes. Usar uma interface na nossa função `Contagem` é essencial para isso e dá certa flexibilidade à função que a chamar.

```go
type Sleeper interface {
    Sleep()
}
```

Tomei uma decisão de design que nossa função `Contagem` não seria responsável por quanto tempo o sleep leva. Isso simplifica um pouco nosso código, pelo menos por enquanto, e significa que um usuário da nossa função pode configurar a duração desse tempo como preferir.

Agora precisamos criar um _mock_ disso para usarmos nos nossos testes.

```go
type SleeperSpy struct {
    Chamadas int
}

func (s *SleeperSpy) Sleep() {
    s.Chamadas++
}
```

_Spies_ (espiões) são um tipo de _mock_ em que podemos gravar como uma dependência é usada. Eles podem gravar os argumentos definidos, quantas vezes são usados etc. No nosso caso, vamos manter o controle de quantas vezes `Sleep()` é chamada para verificá-la no nosso teste.

Atualize os testes para injetar uma dependência no nosso Espião e verifique se o sleep foi chamado 4 vezes.

```go
func TestContagem(t *testing.T) {
    buffer := &bytes.Buffer{}
    sleeperSpy := &SleeperSpy{}

    Contagem(buffer, sleeperSpy)

    resultado := buffer.String()
    esperado := `3
2
1
Vai!`

    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }

    if sleeperSpy.Chamadas != 4 {
        t.Errorf("não houve chamadas suficientes do sleeper, esperado 4, resultado %d", sleeperSpy.Chamadas)
    }
}
```

## Execute o teste

```bash
too many arguments in call to Contagem
    have (*bytes.Buffer, *SpySleeper)
    want (io.Writer)
```

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Precisamos atualizar a `Contagem` para aceitar nosso `Sleeper`:

```go
func Contagem(saida io.Writer, sleeper Sleeper) {
    for i := inicioContagem; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(saida, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(saida, ultimaPalavra)
}
```

Se tentar novamente, nossa `main` não vai mais compilar pelo menos motivo:

```text
./main.go:26:11: not enough arguments in call to Contagem
    have (*os.File)
    want (io.Writer, Sleeper)
```

Vamos criar um sleeper _de verdade_ que implementa a interface que precisamos:

```go
type SleeperPadrao struct {}

func (d *SleeperPadrao) Sleep() {
	time.Sleep(1 * time.Second)
}
```

Podemos usá-lo na nossa aplicação real, como:

```go
func main() {
    sleeper := &SleeperPadrao{}
    Contagem(os.Stdout, sleeper)
}
```

## Escreva código o suficiente para fazer o teste passar

Agora o teste está compilando, mas não passando. Isso acontece porque ainda estamos chamando o `time.Sleep` ao invés da injetada. Vamos arrumar isso.

The test is now compiling but not passing because we're still calling the `time.Sleep` rather than the injected in dependency. Let's fix that.

```go
func Contagem(saida io.Writer, sleeper Sleeper) {
    for i := inicioContagem; i > 0; i-- {
        sleeper.Sleep()
        fmt.Fprintln(saida, i)
    }

    sleeper.Sleep()
    fmt.Fprint(saida, ultimaPalavra)
}
```

O teste deve passar sem levar 4 segundos.

### Ainda temos alguns problemas

Ainda há outra propriedade importante que não estamos testando.

A `Contagem` deve ter uma pausa para cada impressão, como por exemplo:

-   `Pausa`
-   `Imprime N`
-   `Pausa`
-   `Imprime N-1`
-   `Pausa`
-   `Imprime Vai!`
-   etc

Nossa alteração mais recente só verifica se o software teve 4 pausas, mas essas pausas poderiam ocorrer fora de ordem.

Quando escrevemos testes, se não estiver confiante de que seus testes estão te dando confiança o suficiente, quebre-o (mas certifique-se de que você salvou suas alterações antes)! Mude o código para o seguinte:

```go
func Contagem(saida io.Writer, sleeper Sleeper) {
    for i := inicioContagem; i > 0; i-- {
        sleeper.Pausa()
        fmt.Fprintln(saida, i)
    }

    for i := inicioContagem; i > 0; i-- {
        fmt.Fprintln(saida, i)
    }

    sleeper.Pausa()
    fmt.Fprint(saida, ultimaPalavra)
}
```

Se executar seus testes, eles ainda vão passar, apesar da implementação estar errada.

Vamos usar o spy novamente com um novo teste para verificar se a ordem das operações está correta.

Temos duas dependências diferentes e queremos gravar todas as operações delas em uma lista. Logo, vamos criar _um spy para ambas_.

```go
type SpyContagemOperacoes struct {
    Chamadas []string
}

func (s *SpyContagemOperacoes) Pausa() {
    s.Chamadas = append(s.Chamadas, pausa)
}

func (s *SpyContagemOperacoes) Write(p []byte) (n int, err error) {
    s.Chamadas = append(s.Chamadas, escrita)
    return
}

const escrita = "escrita"
const pausa = "pausa"
```

Nosso `SpyContagemOperacoes` implementa tanto o `io.Writer` quanto o `Sleeper`, gravando cada chamada em um slice. Nesse teste, temos preocupação apenas na ordem das operações, então apenas gravá-las em uma lista de operações nomeadas é suficiente.

Agora podemos adicionar um subteste no nosso conjunto de testes.

```go
t.Run("pausa antes de cada impressão", func(t *testing.T) {
        spyImpressoraSleep := &SpyContagemOperacoes{}
        Contagem(spyImpressoraSleep, spyImpressoraSleep)

        esperado := []string{
            pausa,
            escrita,
            pausa,
            escrita,
            pausa,
            escrita,
            pausa,
            escrita,
        }

        if !reflect.DeepEqual(esperado, spyImpressoraSleep.Chamadas) {
            t.Errorf("esperado %v chamadas, resultado %v", esperado, spyImpressoraSleep.Chamadas)
        }
    })
```

Esse teste deve falhar. Volte o código que quebramos para a versão correta e agora o novo teste deve passar.

Agora temos dois spies no `Sleeper`. O próximo passo é refatorar nosso teste para que um teste o que está sendo impresso e o outro se certifique de que estamos pausando entre as impressões. Por fim, podemos apagar nosso primeiro spy, já que não é mais utilizado.

```go
func TestContagem(t *testing.T) {

    t.Run("imprime 3 até Vai!", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        Contagem(buffer, &SpyContagemOperacoes{})

        resultado := buffer.String()
        esperado := `3
2
1
Vai!`

        if resultado != esperado {
            t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
        }
    })

    t.Run("pausa antes de cada impressão", func(t *testing.T) {
        spyImpressoraSleep := &SpyContagemOperacoes{}
        Contagem(spyImpressoraSleep, spyImpressoraSleep)

        esperado := []string{
            pausa,
            escrita,
            pausa,
            escrita,
            pausa,
            escrita,
            pausa,
            escrita,
        }

        if !reflect.DeepEqual(esperado, spyImpressoraSleep.Chamadas) {
            t.Errorf("esperado %v chamadas, resultado %v", esperado, spyImpressoraSleep.Chamadas)
        }
    })
}
```

Agora temos nossa função e suas duas propriedades testadas adequadamente.

## Extendendo o Sleeper para se tornar configurável

Uma funcionalidadee legal seria o `Sleeper` seja configurável.

### Escreva o teste primeiro

Agora vamos criar um novo tipo para `SleeperConfiguravel` que aceita o que precisamos para configuração e teste.

```go
type ConfigurableSleeper struct {
    duration time.Duration
    sleep    func(time.Duration)
}
```

We are using `duration` to configure the time slept and `sleep` as a way to pass in a sleep function. The signature of `sleep` is the same as for `time.Sleep` allowing us to use `time.Sleep` in our real implementation and a spy in our tests.

```go
type SpyTime struct {
    durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
    s.durationSlept = duration
}
```

With our spy in place, we can create a new test for the configurable sleeper.

```go
func TestConfigurableSleeper(t *testing.T) {
    sleepTime := 5 * time.Second

    spyTime := &SpyTime{}
    sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
    sleeper.Sleep()

    if spyTime.durationSlept != sleepTime {
        t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
    }
}
```

There should be nothing new in this test and it is setup very similar to the previous mock tests.

### Try and run the test

```text
sleeper.Sleep undefined (type ConfigurableSleeper has no field or method Sleep, but does have sleep)
```

You should see a very clear error message indicating that we do not have a `Sleep` method created on our `ConfigurableSleeper`.

### Write the minimal amount of code for the test to run and check failing test output

```go
func (c *ConfigurableSleeper) Sleep() {
}
```

With our new `Sleep` function implemented we have a failing test.

```text
countdown_test.go:56: should have slept for 5s but slept for 0s
```

### Write enough code to make it pass

All we need to do now is implement the `Sleep` function for `ConfigurableSleeper`.

```go
func (c *ConfigurableSleeper) Sleep() {
    c.sleep(c.duration)
}
```

With this change all of the test should be passing again.

### Cleanup and refactor

The last thing we need to do is to actually use our `ConfigurableSleeper` in the main function.

```go
func main() {
    sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
    Countdown(os.Stdout, sleeper)
}
```

If we run the tests and the program manually, we can see that all the behavior remains the same.

Since we are using the `ConfigurableSleeper`, it is safe to delete the `DefaultSleeper` implementation. Wrapping up our program.

## But isn't mocking evil?

You may have heard mocking is evil. Just like anything in software development it can be used for evil, just like [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself).

People normally get in to a bad state when they don't _listen to their tests_ and are _not respecting the refactoring stage_.

If your mocking code is becoming complicated or you are having to mock out lots of things to test something, you should _listen_ to that bad feeling and think about your code. Usually it is a sign of

-   The thing you are testing is having to do too many things
    -   Break the module apart so it does less
-   Its dependencies are too fine-grained
    -   Think about how you can consolidate some of these dependencies into one meaningful module
-   Your test is too concerned with implementation details
    -   Favour testing expected behaviour rather than the implementation

Normally a lot of mocking points to _bad abstraction_ in your code.

**What people see here is a weakness in TDD but it is actually a strength**, more often than not poor test code is a result of bad design or put more nicely, well-designed code is easy to test.

### But mocks and tests are still making my life hard!

Ever run into this situation?

-   You want to do some refactoring
-   To do this you end up changing lots of tests
-   You question TDD and make a post on Medium titled "Mocking considered harmful"

This is usually a sign of you testing too much _implementation detail_. Try to make it so your tests are testing _useful behaviour_ unless the implementation is really important to how the system runs.

It is sometimes hard to know _what level_ to test exactly but here are some thought processes and rules I try to follow:

-   **The definition of refactoring is that the code changes but the behaviour stays the same**. If you have decided to do some refactoring in theory you should be able to do make the commit without any test changes. So when writing a test ask yourself
    -   Am i testing the behaviour I want or the implementation details?
    -   If i were to refactor this code, would I have to make lots of changes to the tests?
-   Although Go lets you test private functions, I would avoid it as private functions are to do with implementation.
-   I feel like if a test is working with **more than 3 mocks then it is a red flag** - time for a rethink on the design
-   Use spies with caution. Spies let you see the insides of the algorithm you are writing which can be very useful but that means a tighter coupling between your test code and the implementation. **Be sure you actually care about these details if you're going to spy on them**

As always, rules in software development aren't really rules and there can be exceptions. [Uncle Bob's article of "When to mock"](https://8thlight.com/blog/uncle-bob/2014/05/10/WhenToMock.html) has some excellent pointers.

## Wrapping up

### More on TDD approach

-   When faced with less trivial examples, break the problem down into "thin vertical slices". Try to get to a point where you have _working software backed by tests_ as soon as you can, to avoid getting in rabbit holes and taking a "big bang" approach.
-   Once you have some working software it should be easier to _iterate with small steps_ until you arrive at the software you need.

> "When to use iterative development? You should use iterative development only on projects that you want to succeed."

Martin Fowler.

### Mocking

-   **Without mocking important areas of your code will be untested**. In our case we would not be able to test that our code paused between each print but there are countless other examples. Calling a service that _can_ fail? Wanting to test your system in a particular state? It is very hard to test these scenarios without mocking.
-   Without mocks you may have to set up databases and other third parties things just to test simple business rules. You're likely to have slow tests, resulting in **slow feedback loops**.
-   By having to spin up a database or a webservice to test something you're likely to have **fragile tests** due to the unreliability of such services.

Once a developer learns about mocking it becomes very easy to over-test every single facet of a system in terms of the _way it works_ rather than _what it does_. Always be mindful about **the value of your tests** and what impact they would have in future refactoring.

In this post about mocking we have only covered **Spies** which are a kind of mock. There are different kind of mocks. [Uncle Bob explains the types in a very easy to read article](https://8thlight.com/blog/uncle-bob/2014/05/14/TheLittleMocker.html). In later chapters we will need to write code that depends on others for data, which is where we will show **Stubs** in action.
