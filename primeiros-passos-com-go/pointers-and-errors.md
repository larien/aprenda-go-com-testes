# Ponteiros e erros

[**Você pode encontrar todos os códigos deste capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/pointers)


Nós aprendemos sobre estruturas na última seção, o que nos possibilitou capturar valores com conceito relacionado.

Em algum momento talvez você deseje utilizar estruturas para gerenciar valores, expondo métodos que permitam usuários mudá-los de um jeito que você possa controlar.

**[Fintechs](https://www.infowester.com/fintech.php) amam Go** e uhh bitcoins? Então vamos mostrar um sistema bancário incrível que podemos construir.

Vamos construir uma estrutura de `Carteira` que possamos depositar `Bitcoin`.

## Escreva o teste primeiro

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

No [exemplo anterior](structs-methods-and-interfaces.md) acessamos campos diretamente pelo nome. Entretanto, na nossa _carteira super protegida_, não queremos expor o valor interno para o resto do mundo. Queremos controlar o acesso por meio de métodos.

## Tente rodar o teste

`./carteira_test.go:7:12: undefined: Carteira`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

O compilador não sabe o que uma `Carteira` é, então vamos declará-la.

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

## Escreva código o suficiente para fazer o teste passar

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

Podemos consertar isso com _ponteiros_. [Ponteiros](https://gobyexample.com/pointers) nos permitem _apontar_ para alguns valores e então mudá-los. Então, em vez de termos uma cópia da Carteira, nós pegamos um ponteiro para a carteira para que possamos alterá-la.

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

Dissemos que estávamos fazendo uma carteira Bitcoin, mas até agora nós não os mencionamos. Estamos usando `int` porque é um bom tipo para contar coisas!

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

Isto deixa mais claro o que está acontecendo em nossos testes.

O próximo requisito é para a função `Retirar`.

## Escreva o teste primeiro

Basicamente o aposto da função `Depositar()`

```go
func TestCarteira(t *testing.T) {

    t.Run("Depositar", func(t *testing.T) {
        carteira := Carteira{}

        carteira.Depositar(Bitcoin(10))

        valor := carteira.Balance()

        valorEsperado := Bitcoin(10)

        if valor != valorEsperado {
            t.Errorf("valor %s valorEsperado %s", valor, valorEsperado)
        }
    })

    t.Run("Retirar", func(t *testing.T) {
        carteira := Carteira{saldo: Bitcoin(20)}

        carteira.Retirar(Bitcoin(10))

        valor := carteira.Balance()

        valorEsperado := Bitcoin(10)

        if valor != valorEsperado {
            t.Errorf("valor %s valorEsperado %s", valor, valorEsperado)
        }
    })
}
```

## Tente rodar o teste

`./wallet_test.go:26:9: carteira.Retirar undefined (type Carteira has no field or method Retirar)`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

```go
func (c *Carteira) Retirar(quantidade Bitcoin) {

}
```

`wallet_test.go:33: valor 20 BTC valorEsperado 10 BTC`

## Escreva código o suficiente para fazer o teste passar

```go
func (c *Carteira) Retirar(quantidade Bitcoin) {
    c.saldo -= quantidade
}
```

## Refatoração

Há algumas duplicações em nossos testes, vamos refatorar isto.

```go
func TestCarteira(t *testing.T) {

    assertBalance := func(t *testing.T, carteira Carteira, valorEsperado Bitcoin) {
        t.Helper()
        valor := carteira.Saldo()

        if valor != ValorEsperado {
            t.Errorf("valor %s valorEsperado %s", valor, valorEsperado)
        }
    }

    t.Run("Depositar", func(t *testing.T) {
        carteira := Carteira{}
        carteira.Depositar(Bitcoin(10))
        assertBalance(t, carteira, Bitcoin(10))
    })

    t.Run("Retirar", func(t *testing.T) {
        carteira := Carteira{saldo: Bitcoin(20)}
        carteira.Retirar(Bitcoin(10))
        assertBalance(t, carteira, Bitcoin(10))
    })

}
```

O que aconteceria se você tentasse `Retirar` mais do que há de saldo na conta? Por enquanto, nossos requisitos é assumir que não há nenhum tipo de cheque-especial.

Como sinalizamos um problema quando estivermos usando `Retirar` ?

Em Go, se você quiser indicar um erro, sua função deve retornar um `err` para que quem a chamou possar checar e tratar.

Vamos tentar isto em um teste.

## Escreva o primeiro teste

```go
t.Run("Retirar saldo insuficiente", func(t *testing.T) {
    saldoInicial := Bitcoin(20)
    carteira := Carteira{saldoInicial}
    erro := carteira.Retirar(Bitcoin(100))

    assertBalance(t, carteira, saldoInicial)

    if erro == nil {
        t.Error("Esperava um erro mas nenhum ocorreu.")
    }
})
```

Nós queremos que `Retirar` retorne um erro se tentarmos retirar mais do que temos, e o saldo deverá continuar o mesmo.

Nós checamos se um erro foi retornado falhando o teste se o valor for `nil`.

`nil` é sinônimo de `null` de outras linguagens de programação.
Erros podem ser `nil`, porque o tipo do retorno de `Retirar` vai ser `error`, que é uma interface. Se você ver uma função que tem argumentos ou retornos que são interfaces, eles podem ser nulos.

Do mesmo jeito que `null`, se tentarmos acessar um valor que é `nil`, isto irá disparar um **runtime panic**. Isto é ruim! Devemos ter certeza que tratamos os valores nulos.

## Execute o teste

`./wallet_test.go:31:25: carteira.Retirar(Bitcoin(100)) used as value`

The wording is perhaps a little unclear, but our previous intent with `Retirar` was just to call it, it will never return a value. To make this compile we will need to change it so it has a return type.

Talvez não esteja tão claro, mas nossa intenção era apenas invocar a função `Retirar`, ela nunca irá retornar um valor. Para fazer compilar, precisaremos mudar a função para que retorne um tipo.

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

```go
func (c *Carteira) Retirar(quantidade Bitcoin) error {
    c.saldo -= quantidade
    return nil
}
```

Novamente, é muito importante escrever apenas o suficiente para compilar. Nós corrigimos o método `Retirar` para retornar `error` e por agora temos que retornar _alguma coisa_, então vamos apenas retornar `nil` .

## Escreva código o suficiente para fazer o teste passar

```go
func (c *Carteira) Retirar(quantidade Bitcoin) error {

    if quantidade > c.saldo {
        return errors.New("Ah não!")
    }

    c.saldo -= quantidade
    return nil
}
```

Lembre-se de importar `errors`.

`errors.New` cria um novo `error` com a mensagem escolhida.

## Refatorando

Vamos fazer um rápido helper de teste para nossa checagem de erro, para deixar nosso teste mais legível.

```go
assertError := func(t *testing.T, err error) {
    t.Helper()
    if err == nil {
        t.Error("Esperava um erro mas nenhum ocorreu.")
    }
}
```

E em nosso teste

```go
t.Run("Retirar saldo insuficiente", func(t *testing.T) {
    carteira := Carteira{Bitcoin(20)}
    erro := carteira.Retirar(Bitcoin(100))

    assertBalance(t, carteira, Bitcoin(20))
    assertError(t, erro)
})
```

Acredito, que quando retornamos um erro "oh no", você deve estar pensando que _devessemos_ ponderar melhor, aliás isto não parece tão útil para ser retornado.

Assumindo que o erro enfim foi retornado para o usuário, vamos atualizar nosso teste para verificar em algum tipo de mensagem de erro em vez de apenas checar a existência de um erro.

## Escreva o primeiro teste

Atualize nosso helper para comparar com uma `string`.

```go
assertError := func(t *testing.T, valor error, valorEsperado string) {
    t.Helper()
    if valor == nil {
        t.Fatal("Esperava um erro mas nenhum ocorreu.")
    }

    if valor.Error() != valorEsperado {
        t.Errorf("valor '%s', valorEsperado '%s'", valor, valorEsperado)
    }
}
```

E então atualize o *invocador

```go
t.Run("Retirar saldo insuficiente", func(t *testing.T) {
    saldoInicial := Bitcoin(20)
    carteira := Carteira{saldoInicial}
    erro := carteira.Retirar(Bitcoin(100))

    assertBalance(t, carteira, saldoInicial)
    assertError(t, erro, "Não pode retirar. Saldo insuficiente")
})
```

Nós apresentamos o `t.Fatal` que interromperá o teste se for chamado.
Isto se deve ao fato de que não queremos fazer mais asserções no erro retornado, se não há um. Sem isto, o teste continuaria e causaria erros por causa do ponteiro `nil`.

## Execute o teste

`wallet_test.go:61: valor err 'Ah não' valorEsperado 'Não pode retirar. Saldo insuficiente'`

## Escreva código o suficiente para fazer o teste passar

```go
func (c *Carteira) Retirar(quantidade Bitcoin) error {

    if quantidade > c.saldo {
        return errors.New("Não pode retirar. Saldo insuficiente")
    }

    c.saldo -= quantidade
    return nil
}
```

## Refatorando

Nós temos duplicação da mensagem de erro tanto no código de teste, quanto no código de `Retirar`.

Seria chato se o teste falhasse porque alguém ter mudado a mensagem do erro e é muito detalhe para o nosso test. Nós não _necessariamente_ nos importamos qual mensagem é exatamente, apenas que algum tipo de erro significativo sobre a função é retornado dado uma certa condição.

Em Go, erros são valores, então podemos refatorar isso para ser uma variável e termos apenas uma fonte da verdade.

```go
var ErroSaldoInsuficiente = errors.New("Não pode retirar. Saldo insuficiente")

func (c *Carteira) Retirar(amount Bitcoin) error {

    if quantidade > c.saldo {
        return ErroSaldoInsuficiente
    }

    c.saldo -= quantidade
    return nil
}
```

A palavra-chave `var` nos permite definir valores globais para o pacote.

Está uma é uma mudança positiva porque agora nossa função `Retirar` parece mais limpa.

Agora, nós podemos refatorar nosso código para usar este valor em vez de uma string específica.

```go
func TestCarteira(t *testing.T) {

    t.Run("Depositar", func(t *testing.T) {
        carteira := Carteira{}
        carteira.Depositar(Bitcoin(10))
        assertBalance(t, carteira, Bitcoin(10))
    })

    t.Run("Retirar com saldo suficiente", func(t *testing.T) {
        carteira := Carteira{Bitcoin(20)}
        carteira.Retirar(Bitcoin(10))
        assertBalance(t, carteira, Bitcoin(10))
    })

    t.Run("Retirar saldo insuficiente", func(t *testing.T) {
        carteira := Carteira{Bitcoin(20)}
        erro := carteira.Retirar(Bitcoin(100))

        assertBalance(t, carteira, Bitcoin(20))
        assertError(t, erro, ErroSaldoInsuficiente)
    })
}

func assertBalance(t *testing.T, carteira Carteira, valorEsperado Bitcoin) {
    t.Helper()
    valorEsperado := carteira.Saldo()

    if valor != valorEsperado {
        t.Errorf("valor '%s' valorEsperado '%s'", valor, valorEsperado)
    }
}

func assertError(t *testing.T, valor error, valoresperado error) {
    t.Helper()
    if valor == nil {
        t.Fatal("Esperava um erro mas nenhum ocorreu.")
    }

    if valor != valorEsperado {
        t.Errorf("valor '%s', valorEsperado '%s'", valor, valorEsperado)
    }
}
```

Agora nosso teste está mais fácil para dar continuidade.

Nós apenas movemos os helpers para fora da função principal de teste, então, quando alguém abrir o arquivo, começara lendo nossas asserções primeiro em vez de alguns helpers.

Outra propriedade útil de testes, é que eles nos ajudam a entender o uso _real_ do nosso código, e assim podemos fazer códigos mais compreensivos. Podemos ver aqui que um desenvolvedor pode simplesmente chamar nosso código e fazer uma comparação de igualdade a `ErroSaldoInsuficiente`, e então agir de acordo.

### Erros não checados

Embora o compilador do Go ajude bastante, as vezes há coisas que você pode errar e o tratamento de erro pode ser complicado.

Há um cenário que nós não testamos. Para descobri-lo, execute o comando a seguir no terminal para instalar o `errcheck`, um dos muitos linters disponíveis em Go.

`go get -u github.com/kisielk/errcheck`

Então, dentro do diretório do seu código execute `errcheck .`

Você deve receber algo assim

`wallet_test.go:17:18: carteira.Retirar(Bitcoin(10))`

O que isso está nos dizendo é que nós não checamos o erro sendo retornado naquela linha de código. Aquela linha de código, no meu computador, corresponde para o nosso cenário normal de retirada, porque nós não checamos que se `Retirar` é bem sucedido, um erro _não_ é retornado.

Aqui está o código de teste final que resolve isto.

```go
func TestCarteira(t *testing.T) {

    t.Run("Depositar", func(t *testing.T) {
        carteira := Carteira{}
        carteira.Depositar(Bitcoin(10))

        assertBalance(t, carteira, Bitcoin(10))
    })

    t.Run("Retirar com saldo suficiente", func(t *testing.T) {
        carteira := Carteira{Bitcoin(20)}
        erro := carteira.Retirar(Bitcoin(10))

        assertBalance(t, carteira, Bitcoin(10))
        assertNoError(t, erro)
    })

    t.Run("Retirar com saldo insuficiente", func(t *testing.T) {
        carteira := Carteira{Bitcoin(20)}
        erro := carteira.Retirar(Bitcoin(100))

        assertBalance(t, carteira, Bitcoin(20))
        assertError(t, erro, ErroSaldoInsuficiente)
    })
}

func assertBalance(t *testing.T, carteira Carteira, valorEsperado Bitcoin) {
    t.Helper()
    valorEsperado := carteira.Saldo()

    if valor != valorEsperado {
        t.Errorf("valor %s valorEsperado %s", valor, valorEsperado)
    }
}

func assertNoError(t *testing.T, valor error) {
    t.Helper()
    if valor != nil {
        t.Fatal("Esperava um erro mas nenhum ocorreu.")
    }
}

func assertError(t *testing.T, valor error, valorEsperado error) {
    t.Helper()
    if valor == nil {
        t.Fatal("Esperava um erro mas nenhum ocorreu.")
    }

    if valor != valorEsperado {
        t.Errorf("valor %s, valrEsperado %s", valor, valorEsperado)
    }
}
```

## Resumindo

### Ponteiros

* Go copia os valores quando são passados para funções/métodos, então, se você está escrevendo uma função que precise mudar o estado, você precisará de um ponteiro para o valor que você quer mudar.
* O fato que Go pega um cópia dos valores é muito útil na maior parte dos tempos, mas as vezes você não vai querer que o seu sistema faça cópia de alguma coisa, nesse caso você precisa passar uma referência. Podemos ser dados muito grandes por exemplo, ou talvez coisas que você pretende ter apenas uma instância \(como conexões a banco de dados\).

### nil

* Ponteiros podem ser nil
* Quando uma função retorna um ponteiro para algo, você precisa ter certeza de checar se é nil ou você precisa disparar uma exceção em tempo de execução, o compilador não te ajudará aqui.
* Útil para quando você quer descrever um valor que pode estar faltando.
### Erros

* Erros são a forma de sinalizar falhas quando executar um função/método.
* Analisando nossos testes, concluímos que buscando por uma string em um erro poderia resultar em um teste não muito confiável. Então, nós refatoramos para usar um valor significativo, e isto resultou em um código mais fácil de ser testado, e concluímos que seria mais fácil para usuários de nossa API também.
* Este não é o fim do assunto de tratamento de erros, você pode fazer coisas mais sofisticadas, mas esta é apenas uma introdução. Seções posteriores vão abordar mais estratégias.
* [Não cheque erros apenas, trate os graciosamente](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

### Crie novos tipos a partir de existentes

* Útil para adicionar domínios mais específicos a valores
* Permite implementar interfaces

Ponteiros e erros são uma grande parte de escrita em Go que você precisa estar confortável. Por sorte, _na maioria das vezes_, o compilador irá ajudar se você fizer algo errado, apenas tire um tempo e leia a mensagem de erro.

