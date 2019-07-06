# Ponteiros e erros

[**Você pode encontrar todos os códigos deste capítulo aqui**](https://github.com/quii/learn-go-with-tests/tree/master/pointers)


Nós aprendemos sobre estruturas na última seção, o que nos possibilitou capturar valores com conceito relacionado.

Em algum momento talvez você deseje utilizar estruturas para gerenciar valores, expondo métodos que permitam usuários muda-los de um jeito que você possa controlar.

**[Fintechs](https://www.infowester.com/fintech.php) amam Go** e uhh bitcoins? Então vamos mostrar um sistema bancário incrível que podemos construir.

Vamos construir uma estrutura de `Carteira` que possamos depositar `Bitcoin`.

## Escreva o primeiro teste

```go
func TestCarteira(t *testing.T) {

    carteira := Carteira{}

    carteira.Depositar(10)

    valor := carteira.Saldo()
    valorEsperado := 10

    if valor != valorEsperado {
        t.Errorf("valor %d valorEsperado %d", valor, valorEsperado)
    }
}
```

No [exemplo anterior](structs-methods-and-interfaces.md) nós acessamos campos diretamente pelo nome, entretanto na nossa _carteira super protegida_, nós não queremos expor o valor interno para o resto do mundo. Queremos controlar o acesso por meio de métodos.

## Tente rodar o teste

`./carteira_test.go:7:12: undefined: Carteira`

## Escreva o mínimo de código para o código executar e verifique a saída de erro do teste

O compilador não sabe o que uma `Carteira` é, então vamos declara-la.

```go
type Carteira struct { }
```

Agora que declaramos nossa carteira, tente rodar o teste novamente

```go
./carteira_test.go:9:8: carteira.Depositar undefined (type Carteira has no field or method Depositar)
./carteira_test.go:11:15: carteira.Saldo undefined (type Carteira has no field or method Saldo)
```
Nós precisamos definir estes métodos.

Lembre-se de apenas fazer o necessário para fazer os testes rodarem. Nós precisamos ter certeza que nossos testes falhem corretamente com uma mensagem de erro clara.

```go
func (c Carteira) Depositar(quantidade int) {

}

func (c Carteira) Saldo() int {
    return 0
}
```

Se essa sintaxe não é familiar, dê uma lida na seção de structs.

Os testes agora devem compilar e rodar

`carteira_test.go:15: valor 0 valorEsperado 10`

## Codifique o suficiente para fazer passar

Precisaremos de algum tipo de variável de _saldo_ em nossa estrutura para guardar o valor

```go
type Carteira struct {
    saldo int
}
```

Em Go, se uma variável, tipo, função e etc, começam com um símbolo minúsculo, então esta será privada para _outros pacotes que não seja o que a definiu_.

No nosso caso, noś queremos que apenas nossos métodos sejam capazes de manipular os valores.

Lembre-se, podemos acessar o valor interno do campo `saldo` usando a variável "receptora".

```go
func (c Carteira) Depositar(quantidade int) {
    c.saldo += quantidade
}

func (c Carteira) Saldo() int {
    return c.saldo
}
```

Com a nossa carreira em Fintechs segura, rode os testes para nos aquecermos para passarmos no teste.

`carteira_test.go:15: valor 0 valorEsperado 10`

### ????

Ok, isso é confuso. Parece que nosso código deveria funcionar, nós adicionamos nosso novo valor ao saldo e então o método saldo deveria retornar o valor atual.

Em Go, **quando uma função ou um método é invocado, os argumentos são** _**copiados**_.

Quando `func (c Carteira) Depositar(quantidade int)` é chamado, o `c` é uma cópia do valor de qualquer lugar que o método tenha sido chamado.

Não focando tanto em Ciências da Computação, quando criamos um valor - como uma carteira, este é alocado em algum lugar da memória. Você pode descobrir o _endereço_ desse bit de memória com `&meuValor`.

Experimente isso adicionando alguns prints no código

```go
func TestCarteira(t *testing.T) {

    carteira := Carteira{}

    carteira.Depositar(10)

    valor := carteira.Saldo()

    fmt.Printf("endereço do saldo no teste é %v \n", &carteira.saldo)

    valorEsperado := 10

    if valor != valorEsperado {
        t.Errorf("valor %d valorEsperado %d", valor, valorEsperado)
    }
}
```

```go
func (c Carteira) Depositar(quantidade int) {
    fmt.Printf("endereço do saldo no Depositar é %v \n", &c.saldo)
    c.saldo += quantidade
}
```

O `\n` é um caractere de escape, adiciona uma nova linha após imprimir o endereço de memória. Nós obtemos o ponteiro para algo com o símbolo de endereço: `&`.

Agora rode novamente o teste

```text
endereço do saldo no Depositar é 0xc420012268
endereço do saldo no teste é is 0xc420012260
```

Você pode ver que os endereços dos dois saldos são diferentes. Então, quando mudamos o valor de um dos saldos dentro do código, estamos trabalhando em uma cópia do que veio do teste. Portanto, o saldo no teste não é alterado.

Podemos consertar isso com _ponteiros_. [Ponteiros](https://gobyexample.com/pointers) nos permite _apontar_ para alguns valores e então mudá-los. Então, em vez de termos uma cópia da Carteira, nós pegamos um ponteiro para a carteira para que possamos alterá-la.

```go
func (c *Carteira) Depositar(quantidade int) {
    c.saldo += quantidade
}

func (c *Carteira) Saldo() int {
    return c.saldo
}
```

A diferença é que o tipo do argumento é `*Carteira` em vez de `Carteira` que você pode ler como "um ponteiro para uma carteira".

Rode novamente os testes e eles devem passar.

## Refatorar

Dissemos que estávamos fazendo uma carteira Bitcoin, mas até agora nós não os mecionamos. Estamos usando `int` porque é um bom tipo para contar coisas!

Parece um pouco exagerado criar uma `struct` para isso. `int` é o suficiente em termos de como funciona, mas não é descritivo o suficiente.

Go permite criarmos novos tipos a partir de tipos existentes.

A sintaxe é `type MeuNome TipoOriginal`

```go
type Bitcoin int

type Carteira struct {
    saldo Bitcoin
}

func (c *Carteira) Depositar(quantidade Bitcoin) {
    c.saldo += quantidade
}

func (c *Carteira) Saldo() Bitcoin {
    return c.saldo
}
```

```go
func TestCarteira(t *testing.T) {

    carteira := Carteira{}

    carteira.Depositar(Bitcoin(10))

    valor := carteira.Saldo()

    valorEsperado := Bitcoin(10)

    if valor != valorEsperado {
        t.Errorf("valor %d valorEsperado %d", valor, valorEsperado)
    }
}
```

Para criarmos `Bitcoin` basta usar a sintaxe `Bitcoin(999)`.

Ao fazermos isso, estamos criando um novo tipo e podemos declarar _métodos_ nele. Isto pode ser muito útil quando queremos adicionar funcionalidades de domínios específicos à tipos já existentes.

Vamos implementar [Stringer](https://golang.org/pkg/fmt/#Stringer) no Bitcoin

```go
type Stringer interface {
        String() string
}
```

Essa interface é definida no pacote `fmt` e permite definir como seu tipo é impresso quando utilizado com o operador de string `%s` em prints.

```go
func (b Bitcoin) String() string {
    return fmt.Sprintf("%d BTC", b)
}
```

Como podemos ver, a sintaxe para criar um método em um tipo definido por nós é a mesma que a utilizada em uma struct.

Agora precisamos atualizar nossas impressões de strings no teste para que usem `String()`.

```go
    if valor != valorEsperado {
        t.Errorf("valor %s valorEsperado %s", valor, valorEsperado)
    }
```

Para ver funcionando, quebre o teste de propósito para que possamos ver

`carteira_test.go:18: valor 10 BTC valorEsperado 20 BTC`

This makes it clearer what's going on in our test.

Isto deixa mais claro o que está acontecendo em nossos testes.

The next requirement is for a `Withdraw` function.

O próximo requisito é para a função `Retirar`.

## Write the test first

Pretty much the opposite of `Deposit()`

```go
func TestWallet(t *testing.T) {

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}

        wallet.Deposit(Bitcoin(10))

        got := wallet.Balance()

        want := Bitcoin(10)

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    })

    t.Run("Withdraw", func(t *testing.T) {
        wallet := Wallet{balance: Bitcoin(20)}

        wallet.Withdraw(Bitcoin(10))

        got := wallet.Balance()

        want := Bitcoin(10)

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    })
}
```

## Try to run the test

`./wallet_test.go:26:9: wallet.Withdraw undefined (type Wallet has no field or method Withdraw)`

## Write the minimal amount of code for the test to run and check the failing test output

```go
func (w *Wallet) Withdraw(amount Bitcoin) {

}
```

`wallet_test.go:33: got 20 BTC want 10 BTC`

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) {
    w.balance -= amount
}
```

## Refactor

There's some duplication in our tests, lets refactor that out.

```go
func TestWallet(t *testing.T) {

    assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
        t.Helper()
        got := wallet.Balance()

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    }

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw", func(t *testing.T) {
        wallet := Wallet{balance: Bitcoin(20)}
        wallet.Withdraw(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

}
```

What should happen if you try to `Withdraw` more than is left in the account? For now, our requirement is to assume there is not an overdraft facility.

How do we signal a problem when using `Withdraw` ?

In Go, if you want to indicate an error it is idiomatic for your function to return an `err` for the caller to check and act on.

Let's try this out in a test.

## Write the test first

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{startingBalance}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, startingBalance)

    if err == nil {
        t.Error("wanted an error but didn't get one")
    }
})
```

We want `Withdraw` to return an error _if_ you try to take out more than you have and the balance should stay the same.

We then check an error has returned by failing the test if it is `nil`.

`nil` is synonymous with `null` from other programming languages. Errors can be `nil` because the return type of `Withdraw` will be `error`, which is an interface. If you see a function that takes arguments or returns values that are interfaces, they can be nillable.

Like `null` if you try to access a value that is `nil` it will throw a **runtime panic**. This is bad! You should make sure that you check for nils.

## Try and run the test

`./wallet_test.go:31:25: wallet.Withdraw(Bitcoin(100)) used as value`

The wording is perhaps a little unclear, but our previous intent with `Withdraw` was just to call it, it will never return a value. To make this compile we will need to change it so it has a return type.

## Write the minimal amount of code for the test to run and check the failing test output

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
    w.balance -= amount
    return nil
}
```

Again, it is very important to just write enough code to satisfy the compiler. We correct our `Withdraw` method to return `error` and for now we have to return _something_ so let's just return `nil`.

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

    if amount > w.balance {
        return errors.New("oh no")
    }

    w.balance -= amount
    return nil
}
```

Remember to import `errors` into your code.

`errors.New` creates a new `error` with a message of your choosing.

## Refactor

Let's make a quick test helper for our error check just to help our test read clearer

```go
assertError := func(t *testing.T, err error) {
    t.Helper()
    if err == nil {
        t.Error("wanted an error but didnt get one")
    }
}
```

And in our test

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
    wallet := Wallet{Bitcoin(20)}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, Bitcoin(20))
    assertError(t, err)
})
```

Hopefully when returning an error of "oh no" you were thinking that we _might_ iterate on that because it doesn't seem that useful to return.

Assuming that the error ultimately gets returned to the user, let's update our test to assert on some kind of error message rather than just the existence of an error.

## Write the test first

Update our helper for a `string` to compare against.

```go
assertError := func(t *testing.T, got error, want string) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but wanted one")
    }

    if got.Error() != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
}
```

And then update the caller

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{startingBalance}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, startingBalance)
    assertError(t, err, "cannot withdraw, insufficient funds")
})
```

We've introduced `t.Fatal` which will stop the test if it is called. This is because we don't want to make any more assertions on the error returned if there isn't one around. Without this the test would carry on to the next step and panic because of a nil pointer.

## Try to run the test

`wallet_test.go:61: got err 'oh no' want 'cannot withdraw, insufficient funds'`

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

    if amount > w.balance {
        return errors.New("cannot withdraw, insufficient funds")
    }

    w.balance -= amount
    return nil
}
```

## Refactor

We have duplication of the error message in both the test code and the `Withdraw` code.

It would be really annoying for the test to fail if someone wanted to re-word the error and it's just too much detail for our test. We don't _really_ care what the exact wording is, just that some kind of meaningful error around withdrawing is returned given a certain condition.

In Go, errors are values, so we can refactor it out into a variable and have a single source of truth for it.

```go
var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {

    if amount > w.balance {
        return ErrInsufficientFunds
    }

    w.balance -= amount
    return nil
}
```

The `var` keyword allows us to define values global to the package.

This is a positive change in itself because now our `Withdraw` function looks very clear.

Next we can refactor our test code to use this value instead of specific strings.

```go
func TestWallet(t *testing.T) {

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw with funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        wallet.Withdraw(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw insufficient funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        err := wallet.Withdraw(Bitcoin(100))

        assertBalance(t, wallet, Bitcoin(20))
        assertError(t, err, ErrInsufficientFunds)
    })
}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
    t.Helper()
    got := wallet.Balance()

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}

func assertError(t *testing.T, got error, want error) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but wanted one")
    }

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
}
```

And now the test is easier to follow too.

I have moved the helpers out of the main test function just so when someone opens up a file they can start reading our assertions first, rather than some helpers.

Another useful property of tests is that they help us understand the _real_ usage of our code so we can make sympathetic code. We can see here that a developer can simply call our code and do an equals check to `ErrInsufficientFunds` and act accordingly.

### Unchecked errors

Whilst the Go compiler helps you a lot, sometimes there are things you can still miss and error handling can sometimes be tricky.

There is one scenario we have not tested. To find it, run the following in a terminal to install `errcheck`, one of many linters available for Go.

`go get -u github.com/kisielk/errcheck`

Then, inside the directory with your code run `errcheck .`

You should get something like

`wallet_test.go:17:18: wallet.Withdraw(Bitcoin(10))`

What this is telling us is that we have not checked the error being returned on that line of code. That line of code on my computer corresponds to our normal withdraw scenario because we have not checked that if the `Withdraw` is successful that an error is _not_ returned.

Here is the final test code that accounts for this.

```go
func TestWallet(t *testing.T) {

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))

        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw with funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        err := wallet.Withdraw(Bitcoin(10))

        assertBalance(t, wallet, Bitcoin(10))
        assertNoError(t, err)
    })

    t.Run("Withdraw insufficient funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        err := wallet.Withdraw(Bitcoin(100))

        assertBalance(t, wallet, Bitcoin(20))
        assertError(t, err, ErrInsufficientFunds)
    })
}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
    t.Helper()
    got := wallet.Balance()

    if got != want {
        t.Errorf("got %s want %s", got, want)
    }
}

func assertNoError(t *testing.T, got error) {
    t.Helper()
    if got != nil {
        t.Fatal("got an error but didnt want one")
    }
}

func assertError(t *testing.T, got error, want error) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but wanted one")
    }

    if got != want {
        t.Errorf("got %s, want %s", got, want)
    }
}
```

## Wrapping up

### Pointers

* Go copies values when you pass them to functions/methods so if you're writing a function that needs to mutate state you'll need it to take a pointer to the thing you want to change.
* The fact that Go takes a copy of values is useful a lot of the time but sometimes you wont want your system to make a copy of something, in which case you need to pass a reference. Examples could be very large data or perhaps things you intend only to have one instance of \(like database connection pools\).

### nil

* Pointers can be nil
* When a function returns a pointer to something, you need to make sure you check if it's nil or you might raise a runtime exception, the compiler wont help you here.
* Useful for when you want to describe a value that could be missing

### Errors

* Errors are the way to signify failure when calling a function/method.
* By listening to our tests we concluded that checking for a string in an error would result in a flaky test. So we refactored to use a meaningful value instead and this resulted in easier to test code and concluded this would be easier for users of our API too.
* This is not the end of the story with error handling, you can do more sophisticated things but this is just an intro. Later sections will cover more strategies.
* [Don’t just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

### Create new types from existing ones

* Useful for adding more domain specific meaning to values
* Can let you implement interfaces

Pointers and errors are a big part of writing Go that you need to get comfortable with. Thankfully the compiler will _usually_ help you out if you do something wrong, just take your time and read the error.

