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

Try and run the program and be amazed at your handywork.

Yes this seems trivial but this approach is what I would recommend for any project. **Take a thin slice of functionality and make it work end-to-end, backed by tests.**

Next we can make it print 2,1 and then "Go!".

## Write the test first

By investing in getting the overall plumbing working right, we can iterate on our solution safely and easily. We will no longer need to stop and re-run the program to be confident of it working as all the logic is tested.

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}

    Countdown(buffer)

    got := buffer.String()
    want := `3
2
1
Go!`

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

The backtick syntax is another way of creating a `string` but lets you put things like newlines which is perfect for our test.

## Try and run the test

```text
countdown_test.go:21: got '3' want '3
        2
        1
        Go!'
```

## Write enough code to make it pass

```go
func Countdown(out io.Writer) {
    for i := 3; i > 0; i-- {
        fmt.Fprintln(out, i)
    }
    fmt.Fprint(out, "Go!")
}
```

Use a `for` loop counting backwards with `i--` and use `fmt.Fprintln` to print to `out` with our number followed by a newline character. Finally use `fmt.Fprint` to send "Go!" aftward.

## Refactor

There's not much to refactor other than refactoring some magic values into named constants.

```go
const finalWord = "Go!"
const countdownStart = 3

func Countdown(out io.Writer) {
    for i := countdownStart; i > 0; i-- {
        fmt.Fprintln(out, i)
    }
    fmt.Fprint(out, finalWord)
}
```

If you run the program now, you should get the desired output but we don't have it as a dramatic countdown with the 1 second pauses.

Go let's you achieve this with `time.Sleep`. Try adding it in to our code.

```go
func Countdown(out io.Writer) {
    for i := countdownStart; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(out, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(out, finalWord)
}
```

If you run the program it works as we want it to.

## Mocking

The tests still pass and the software works as intended but we have some problems:

-   Our tests take 4 seconds to run.
    -   Every forward thinking post about software development emphasises the importance of quick feedback loops.
    -   **Slow tests ruin developer productivity**.
    -   Imagine if the requirements get more sophisticated warranting more tests. Are we happy with 4s added to the test run for every new test of `Countdown`?
-   We have not tested an important property of our function.

We have a dependency on `Sleep`ing which we need to extract so we can then control it in our tests.

If we can _mock_ `time.Sleep` we can use _dependency injection_ to use it instead of a "real" `time.Sleep` and then we can **spy on the calls** to make assertions on them.

## Write the test first

Let's define our dependency as an interface. This lets us then use a _real_ Sleeper in `main` and a _spy sleeper_ in our tests. By using an interface our `Countdown` function is oblivious to this and adds some flexibility for the caller.

```go
type Sleeper interface {
    Sleep()
}
```

I made a design decision that our `Countdown` function would not be responsible for how long the sleep is. This simplifies our code a little for now at least and means a user of our function can configure that sleepiness however they like.

Now we need to make a _mock_ of it for our tests to use.

```go
type SpySleeper struct {
    Calls int
}

func (s *SpySleeper) Sleep() {
    s.Calls++
}
```

_Spies_ are a kind of _mock_ which can record how a dependency is used. They can record the arguments sent in, how many times, etc. In our case, we're keeping track of how many times `Sleep()` is called so we can check it in our test.

Update the tests to inject a dependency on our Spy and assert that the sleep has been called 4 times.

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}
    spySleeper := &SpySleeper{}

    Countdown(buffer, spySleeper)

    got := buffer.String()
    want := `3
2
1
Go!`

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }

    if spySleeper.Calls != 4 {
        t.Errorf("not enough calls to sleeper, want 4 got %d", spySleeper.Calls)
    }
}
```

## Try and run the test

```text
too many arguments in call to Countdown
    have (*bytes.Buffer, *SpySleeper)
    want (io.Writer)
```

## Write the minimal amount of code for the test to run and check the failing test output

We need to update `Countdown` to accept our `Sleeper`

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(out, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(out, finalWord)
}
```

If you try again, your `main` will no longer compile for the same reason

```text
./main.go:26:11: not enough arguments in call to Countdown
    have (*os.File)
    want (io.Writer, Sleeper)
```

Let's create a _real_ sleeper which implements the interface we need

```go
type DefaultSleeper struct {}

func (d *DefaultSleeper) Sleep() {
    time.Sleep(1 * time.Second)
}
```

We can then use it in our real application like so

```go
func main() {
    sleeper := &DefaultSleeper{}
    Countdown(os.Stdout, sleeper)
}
```

## Write enough code to make it pass

The test is now compiling but not passing because we're still calling the `time.Sleep` rather than the injected in dependency. Let's fix that.

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        sleeper.Sleep()
        fmt.Fprintln(out, i)
    }

    sleeper.Sleep()
    fmt.Fprint(out, finalWord)
}
```

The test should pass and no longer taking 4 seconds.

### Still some problems

There's still another important property we haven't tested.

`Countdown` should sleep before each print, e.g:

-   `Sleep`
-   `Print N`
-   `Sleep`
-   `Print N-1`
-   `Sleep`
-   `Print Go!`
-   etc

Our latest change only asserts that it has slept 4 times, but those sleeps could occur out of sequence.

When writing tests if you're not confident that your tests are giving you sufficient confidence, just break it! \(make sure you have committed your changes to source control first though\). Change the code to the following

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        sleeper.Sleep()
    }

    for i := countdownStart; i > 0; i-- {
        fmt.Fprintln(out, i)
    }

    sleeper.Sleep()
    fmt.Fprint(out, finalWord)
}
```

If you run your tests they should still be passing even though the implementation is wrong.

Let's use spying again with a new test to check the order of operations is correct.

We have two different dependencies and we want to record all of their operations into one list. So we'll create _one spy for them both_.

```go
type CountdownOperationsSpy struct {
    Calls []string
}

func (s *CountdownOperationsSpy) Sleep() {
    s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
    s.Calls = append(s.Calls, write)
    return
}

const write = "write"
const sleep = "sleep"
```

Our `CountdownOperationsSpy` implements both `io.Writer` and `Sleeper`, recording every call into one slice. In this test we're only concerned about the order of operations, so just recording them as list of named operations is sufficient.

We can now add a sub-test into our test suite.

```go
t.Run("sleep before every print", func(t *testing.T) {
    spySleepPrinter := &CountdownOperationsSpy{}
    Countdown(spySleepPrinter, spySleepPrinter)

    want := []string{
        sleep,
        write,
        sleep,
        write,
        sleep,
        write,
        sleep,
        write,
    }

    if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
        t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
    }
})
```

This test should now fail. Revert it back and the new test should pass.

We now have two tests spying on the `Sleeper` so we can now refactor our test so one is testing what is being printed and the other one is ensuring we're sleeping in between the prints. Finally we can delete our first spy as it's not used anymore.

```go
func TestCountdown(t *testing.T) {

    t.Run("prints 3 to Go!", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        Countdown(buffer, &CountdownOperationsSpy{})

        got := buffer.String()
        want := `3
2
1
Go!`

        if got != want {
            t.Errorf("got '%s' want '%s'", got, want)
        }
    })

    t.Run("sleep before every print", func(t *testing.T) {
        spySleepPrinter := &CountdownOperationsSpy{}
        Countdown(spySleepPrinter, spySleepPrinter)

        want := []string{
            sleep,
            write,
            sleep,
            write,
            sleep,
            write,
            sleep,
            write,
        }

        if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
            t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
        }
    })
}
```

We now have our function and its 2 important properties properly tested.

## Extending Sleeper to be configurable

A nice feature would be for the `Sleeper` to be configurable.

### Write the test first

Let's first create a new type for `ConfigurableSleeper` that accepts what we need for configuration and testing.

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
