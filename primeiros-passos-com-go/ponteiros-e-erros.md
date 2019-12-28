# Ponteiros e erros

[**Você pode encontrar todos os códigos deste capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/ponteiros)

Aprendemos sobre estruturas na última seção, o que nos possibilita capturar valores com conceito relacionado.

Em algum momento talvez você deseje utilizar estruturas para gerenciar valores, expondo métodos que permita aos usuários mudá-los de um jeito que você possa controlar.

**[Fintechs](https://www.infowester.com/fintech.php) amam Go** e uhh bitcoins? Então vamos mostrar um sistema bancário incrível que podemos construir.

Vamos construir uma estrutura de `Carteira` que possamos depositar `Bitcoin`.

## Escreva o teste primeiro

```go
func TestCarteira(t *testing.T) {
    carteira := Carteira{}

    carteira.Depositar(10)

    resultado := carteira.Saldo()
    esperado := 10

    if resultado != esperado {
		t.Errorf("resultado %s, esperado %s", resultado, esperado)
	}
}
```

No [exemplo anterior](../estruturas-metodos-e-interfaces/estruturas-metodos-e-interfaces.md) acessamos campos diretamente pelo nome. Entretanto, na nossa _carteira super protegida_, não queremos expor o valor interno para o resto do mundo. Queremos controlar o acesso por meio de métodos.

## Execute o teste

`./carteira_test.go:7:12: undefined: Carteira`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

O compilador não sabe o que uma `Carteira` é, então vamos declará-la.

```go
type Carteira struct { }
```

Agora que declaramos nossa carteira, tente rodar o teste novamente:

```go
./carteira_test.go:9:8: carteira.Depositar undefined (type Carteira has no field or method Depositar)
./carteira_test.go:11:15: carteira.Saldo undefined (type Carteira has no field or method Saldo)
```

Precisamos definir estes métodos.

Lembre-se de apenas fazer o necessário para fazer os testes rodarem. Precisamos ter certeza que nossos testes falhem corretamente com uma mensagem de erro clara.

```go
func (c Carteira) Depositar(quantidade int) {

}

func (c Carteira) Saldo() int {
    return 0
}
```

Se essa sintaxe não for familiar, dê uma lida na seção de estruturas.

Os testes agora devem compilar e rodar:

`carteira_test.go:15: resultado 0, esperado 10`

## Escreva código o suficiente para fazer o teste passar

Precisaremos de algum tipo de variável de _saldo_ em nossa estrutura para guardar o valor:

```go
type Carteira struct {
    saldo int
}
```

Em Go, se uma variável, tipo, função e etc, começam com uma letra minúsculo, então esta será privada para _outros pacotes que não seja o que a definiu_.

No nosso caso, queremos que apenas nossos métodos sejam capazes de manipular os valores.

Lembre-se que podemos acessar o valor interno do campo `saldo` usando a variável "receptora".

```go
func (c Carteira) Depositar(quantidade int) {
    c.saldo += quantidade
}

func (c Carteira) Saldo() int {
    return c.saldo
}
```

Com a nossa carreira em Fintechs segura, rode os testes para nos aquecermos para passarmos no teste.

`carteira_test.go:15: resultado 0, esperado 10`

### ????

Ok, isso é confuso. Parece que nosso código deveria funcionar, pois adicionamos nosso novo valor ao saldo e o método Saldo deveria retornar o valor atual.

Em Go, **quando uma função ou um método é invocado, os argumentos são** _**copiados**_.

Quando `func (c Carteira) Depositar(quantidade int)` é chamado, o `c` é uma cópia do valor de qualquer lugar que o método tenha sido chamado.

Sem focar em Ciência da Computação, quando criamos um valor (como uma carteira), esse valor é alocado em algum lugar da memória. Você pode descobrir o _endereço_ desse bit de memória usando `&meuValor`.

Experimente isso adicionando alguns prints no código:

```go
func TestCarteira(t *testing.T) {
    carteira := Carteira{}

    carteira.Depositar(10)

    resultado := carteira.Saldo()

    fmt.Printf("O endereço do saldo no teste é %v \n", &carteira.saldo)

    esperado := 10

    if resultado != esperado {
        t.Errorf("resultado %d, esperado %d", resultado, esperado)
    }
}
```

```go
func (c Carteira) Depositar(quantidade int) {
    fmt.Printf("O endereço do saldo no Depositar é %v \n", &c.saldo)
    c.saldo += quantidade
}
```

O `\n` é um caractere de escape queeadiciona uma nova linha após imprimir o endereço de memória. Conseguimos acessar o ponteiro para algo com o símbolo de endereço `&`.

Agora rode o teste novamente:

```text
O endereço do saldo no Depositar é 0xc420012268
O endereço do saldo no teste é is 0xc420012260
```

Podemos ver que os endereços dos dois saldos são diferentes. Então, quando mudamos o valor de um dos saldos dentro do código, estamos trabalhando em uma cópia do que veio do teste. Portanto, o saldo no teste não é alterado.

Podemos consertar isso com _ponteiros_. [Ponteiros](https://gobyexample.com/pointers) nos permitem _apontar_ para alguns valores e então mudá-los. Então, em vez de termos uma cópia da Carteira, usamos um ponteiro para a carteira para que possamos alterá-la.

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

## Refatoração

Dissemos que estávamos fazendo uma carteira Bitcoin, mas até agora não os mencionamos. Estamos usando `int` porque é um bom tipo para contar coisas!

Parece um pouco exagerado criar uma `struct` para isso. `int` é o suficiente nesse contexto, mas não é descritivo o suficiente.

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

    resultado := carteira.Saldo()

    esperado := Bitcoin(10)

    if resultado != esperado {
			t.Errorf("resultado %d, esperado %d", resultado, esperado)
		}
}
```

Para criarmos `Bitcoin`, basta usar a sintaxe `Bitcoin(999)`.

Ao fazermos isso, estamos criando um novo tipo e podemos declarar _métodos_ nele. Isto pode ser muito útil quando queremos adicionar funcionalidades de domínios específicos a tipos já existentes.

Vamos implementar um [Stringer](https://golang.org/pkg/fmt/#Stringer) para o Bitcoin:

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
    if resultado != esperado {
		t.Errorf("resultado %s, esperado %s", resultado, esperado)
	}
```

Para ver funcionando, quebre o teste de propósito para que possamos ver:

`carteira_test.go:18: resultado 10 BTC, esperado 20 BTC`

Isto deixa mais claro o que está acontecendo em nossos testes.

O próximo requisito é criar uma função de `Retirar`.

## Escreva o teste primeiro

É basicamente o aposto da função `Depositar()`:

```go
func TestCarteira(t *testing.T) {
    t.Run("Depositar", func(t *testing.T) {
        carteira := Carteira{}

        carteira.Depositar(Bitcoin(10))

        resultado := carteira.Saldo()

        esperado := Bitcoin(10)

        if resultado != esperado {
			t.Errorf("resultado %s, esperado %s", resultado, esperado)
		}
    })

    t.Run("Retirar", func(t *testing.T) {
        carteira := Carteira{saldo: Bitcoin(20)}

        carteira.Retirar(Bitcoin(10))

        resultado := carteira.Saldo()

        esperado := Bitcoin(10)

        if resultado != esperado {
			t.Errorf("resultado %s, esperado %s", resultado, esperado)
		}
    })
}
```

## Execute o teste

`./carteira_test.go:26:9: carteira.Retirar undefined (type Carteira has no field or method Retirar)`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

```go
func (c *Carteira) Retirar(quantidade Bitcoin) {
}
```

`carteira_test.go:33: resultado 20 BTC, esperado 10 BTC`

## Escreva código o suficiente para fazer o teste passar

```go
func (c *Carteira) Retirar(quantidade Bitcoin) {
    c.saldo -= quantidade
}
```

## Refatoração

Há algumas duplicações em nossos testes, vamos refatorar isso.

```go
func TestCarteira(t *testing.T) {
    confirmaSaldo := func(t *testing.T, carteira Carteira, valorEsperado Bitcoin) {
        t.Helper()
		resultado := carteira.Saldo()

		if resultado != esperado {
			t.Errorf("resultado %s, esperado %s", resultado, esperado)
		}
    }

    t.Run("Depositar", func(t *testing.T) {
		carteira := Carteira{}
		carteira.Depositar(Bitcoin(10))
		confirmaSaldo(t, carteira, Bitcoin(10))
	})

	t.Run("Retirar", func(t *testing.T) {
		carteira := Carteira{saldo: Bitcoin(20)}
		carteira.Retirar(10)
		confirmaSaldo(t, carteira, Bitcoin(10))
	})
}
```

O que aconteceria se você tentasse `Retirar` mais do que há de saldo na conta? Por enquanto, nossos requisitos são assumir que não há nenhum tipo de cheque-especial.

Como sinalizamos um problema quando estivermos usando `Retirar` ?

Em Go, se você quiser indicar um erro, sua função deve retornar um `err` para que quem a chamou possar verificá-lo e tratá-lo.

Vamos tentar fazer isso em um teste.

## Escreva o teste primeiro

```go
t.Run("Retirar com saldo insuficiente", func(t *testing.T) {
    saldoInicial := Bitcoin(20)
	carteira := Carteira{saldoInicial}
	erro := carteira.Retirar(Bitcoin(100))

	confirmaSaldo(t, carteira, saldoInicial)

    if erro == nil {
        t.Error("Esperava um erro mas nenhum ocorreu")
    }
})
```

Queremos que `Retirar` retorne um erro se tentarmos retirar mais do que temos e o saldo deverá continuar o mesmo.

Verificamos se um erro foi retornado falhando o teste se o valor for `nil`.

`nil` é a mesma coisa que `null` de outras linguagens de programação.

Erros podem ser `nil`, porque o tipo do retorno de `Retirar` vai ser `error`, que é uma interface. Se você vir uma função que tem argumentos ou retornos que são interfaces, eles podem ser nulos.

Do mesmo jeito que `null`, se tentarmos acessar um valor que é `nil`, isso irá disparar um **pânico em tempo de execução**. Isso é ruim! Devemos ter certeza que tratamos os valores nulos.

## Execute o teste

`./carteira_test.go:31:25: carteira.Retirar(Bitcoin(100)) used as value`

Talvez não esteja tão claro, mas nossa intenção era apenas invocar a função `Retirar` e ela nunca irá retornar um valor pois o saldo será diretamente subtraído com o ponteiro e a função deve apenas retornar o erro (se houver). Para fazer compilar, precisaremos mudar a função para que retorne um tipo.

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

```go
func (c *Carteira) Retirar(quantidade Bitcoin) error {
    c.saldo -= quantidade
    return nil
}
```

Novamente, é muito importante escrever apenas o suficiente para compilar. Corrigimos o método `Retirar` para retornar `error` e por enquanto temos que retornar _alguma coisa_, então vamos apenas retornar `nil` .

## Escreva código o suficiente para fazer o teste passar

```go
func (c *Carteira) Retirar(quantidade Bitcoin) error {
    if quantidade > c.saldo {
        return errors.New("eita")
    }

    c.saldo -= quantidade
    return nil
}
```

Lembre-se de importar `errors`.

`errors.New` cria um novo `error` com a mensagem escolhida.

## Refatoração

Vamos fazer um método auxiliar de teste para nossa verificação de erro para deixar nosso teste mais legível.

```go
confirmaErro := func(t *testing.T, erro error) {
	t.Helper()
	if erro == nil {
		t.Error("esperava um erro, mas nenhum ocorreu.")
	}
}
```

E em nosso teste:

```go
t.Run("Retirar com saldo insuficiente", func(t *testing.T) {
    saldoInicial := Bitcoin(20)
	carteira := Carteira{saldoInicial}
	erro := carteira.Retirar(Bitcoin(100))

	confirmaSaldo(t, carteira, saldoInicial)
	confirmaErro(t, erro)
})
```

Espero que, ao retornamos um erro do tipo "eita", você pense que _devêssemos_ deixar mais claro o que ocorreu, já que esta não parece uma informação útil para nós.

Assumindo que o erro enfim foi retornado para o usuário, vamos atualizar nosso teste para verificar o tipo espcífico de mensagem de erro ao invés de apenas verificar se um erro existe.

## Escreva o teste primeiro

Atualize nosso helper para comparar com uma `string`:

```go
confirmarErro := func(t *testing.T, valor error, valorEsperado string) {
    t.Helper()
	if resultado == nil {
		t.Fatal("esperava um erro, mas nenhum ocorreu")
	}

	if resultado != esperado {
		t.Errorf("resultado %s, esperado %s", resultado, esperado)
	}
}
```

E então atualize o invocador:

```go
t.Run("Retirar saldo insuficiente", func(t *testing.T) {
    saldoInicial := Bitcoin(20)
	carteira := Carteira{saldoInicial}
	erro := carteira.Retirar(Bitcoin(100))

	confirmaSaldo(t, carteira, saldoInicial)
    confirmaErro(t, erro, "não é possível retirar: saldo insuficiente")
})
```

Usamos o `t.Fatal` que interromperá o teste se for chamado. Isso é feito porque não queremos fazer mais asserções no erro retornado, se não houver um. Sem isso, o teste continuaria e causaria erros por causa do ponteiro `nil`.

## Execute o teste

`carteira_test.go:61: erro resultado 'eita', erro esperado 'não é possível retirar: saldo insuficiente'`

## Escreva código o suficiente para fazer o teste passar

```go
func (c *Carteira) Retirar(quantidade Bitcoin) error {

    if quantidade > c.saldo {
        return errors.New("não é possível retirar: saldo insuficiente")
    }

    c.saldo -= quantidade
    return nil
}
```

## Refatoração

Temos duplicação da mensagem de erro tanto no código de teste quanto no código de `Retirar`.

Seria chato se o teste falhasse por alguém ter mudado a mensagem do erro e é muito detalhe para o nosso teste. Nós não _necessariamente_ nos importamos qual mensagem é exatamente, apenas que algum tipo de erro significativo sobre a função é retornado dada uma certa condição.

Em Go, erros são valores, então podemos refatorar isso para ser uma variável e termos apenas uma fonte da verdade.

```go
var ErroSaldoInsuficiente = errors.New("não é possível retirar: saldo insuficiente")

func (c *Carteira) Retirar(amount Bitcoin) error {

    if quantidade > c.saldo {
        return ErroSaldoInsuficiente
    }

    c.saldo -= quantidade
    return nil
}
```

A palavra-chave `var` no escopo do arquivo nos permite definir valores globais para o pacote.

Esta é uma mudança positiva, pois agora nossa função `Retirar` parece mais limpa.

Agora, podemos refatorar nosso código para usar este valor ao invés de uma string específica.

```go
func TestCarteira(t *testing.T) {
	t.Run("Depositar", func(t *testing.T) {
		carteira := Carteira{}
		carteira.Depositar(Bitcoin(10))

		confirmaSaldo(t, carteira, Bitcoin(10))
	})

	t.Run("Retirar com saldo suficiente", func(t *testing.T) {
		carteira := Carteira{Bitcoin(20)}
		erro := carteira.Retirar(Bitcoin(10))

		confirmaSaldo(t, carteira, Bitcoin(10))
		confirmaErroInexistente(t, erro)
	})

	t.Run("Retirar com saldo insuficiente", func(t *testing.T) {
		saldoInicial := Bitcoin(20)
		carteira := Carteira{saldoInicial}
		erro := carteira.Retirar(Bitcoin(100))

		confirmaSaldo(t, carteira, saldoInicial)
		confirmaErro(t, erro, ErroSaldoInsuficiente)
	})
}

func confirmaSaldo(t *testing.T, carteira Carteira, esperado Bitcoin) {
	t.Helper()
	resultado := carteira.Saldo()

	if resultado != esperado {
		t.Errorf("resultado %s, esperado %s", resultado, esperado)
	}
}

func confirmaErro(t *testing.T, resultado error, esperado error) {
	t.Helper()
	if resultado == nil {
		t.Fatal("esperava um erro, mas nenhum ocorreu")
	}

	if resultado != esperado {
		t.Errorf("erro resultado %s, erro esperado %s", resultado, esperado)
	}
}
```

Agora está mais fácil dar continuidade ao nosso teste.

Nós apenas movemos os métodos auxiliares para fora da função principal de teste. Logo, quando alguém abrir o arquivo, começará lendo nossas asserções primeiro ao invés desses métodos auxiliares.

Outra propriedade útil de testes é que eles nos ajudam a entender o uso _real_ do nosso código e assim podemos fazer códigos mais compreensivos. Podemos ver aqui que um desenvolvedor pode simplesmente chamar nosso código e fazer uma comparação de igualdade a `ErroSaldoInsuficiente`, e então agir de acordo.

### Erros não verificados

Embora o compilador do Go ajude bastante, há coisas que você pode acabar errando e o tratamento de erro pode se tornar complicado.

Há um cenário que nós não testamos. Para descobri-lo, execute o comando a seguir no terminal para instalar o `errcheck`, um dos muitos linters disponíveis em Go.

`go get -u github.com/kisielk/errcheck`

Então, dentro do diretório do seu código, execute `errcheck .`.

Você deve receber algo assim:

`carteira_test.go:17:18: carteira.Retirar(Bitcoin(10))`

O que isso está nos dizendo é que não verificamos o erro sendo retornado naquela linha de código. Aquela linha de código, no meu computador, corresponde para o nosso cenário normal de retirada, porque não verificamos que se `Retirar` é bem sucedido quando um erro _não_ é retornado.

Aqui está o código de teste final que resolve isto.

```go
func TestCarteira(t *testing.T) {
	t.Run("Depositar", func(t *testing.T) {
		carteira := Carteira{}
		carteira.Depositar(Bitcoin(10))

		confirmaSaldo(t, carteira, Bitcoin(10))
	})

	t.Run("Retirar com saldo suficiente", func(t *testing.T) {
		carteira := Carteira{Bitcoin(20)}
		erro := carteira.Retirar(Bitcoin(10))

		confirmaSaldo(t, carteira, Bitcoin(10))
		confirmaErroInexistente(t, erro)
	})

	t.Run("Retirar com saldo insuficiente", func(t *testing.T) {
		saldoInicial := Bitcoin(20)
		carteira := Carteira{saldoInicial}
		erro := carteira.Retirar(Bitcoin(100))

		confirmaSaldo(t, carteira, saldoInicial)
		confirmaErro(t, erro, ErroSaldoInsuficiente)
	})
}

func confirmaSaldo(t *testing.T, carteira Carteira, esperado Bitcoin) {
	t.Helper()
	resultado := carteira.Saldo()

	if resultado != esperado {
		t.Errorf("resultado %s, esperado %s", resultado, esperado)
	}
}

func confirmaErroInexistente(t *testing.T, resultado error) {
	t.Helper()
	if resultado != nil {
		t.Fatal("erro inesperado recebido")
	}
}

func confirmaErro(t *testing.T, resultado error, esperado error) {
	t.Helper()
	if resultado == nil {
		t.Fatal("esperava um erro, mas nenhum ocorreu")
	}

	if resultado != esperado {
		t.Errorf("erro resultado %s, erro esperado %s", resultado, esperado)
	}
}
```

## Resumo

### Ponteiros

* Go copia os valores quando são passados para funções/métodos. Então, se estiver escrevendo uma função que precise mudar o estado, você precisará de um ponteiro para o valor que você quer mudar.
* O fato de que Go pega um cópia dos valores é muito útil na maior parte do tempo, mas às vezes você não vai querer que o seu sistema faça cópia de alguma coisa. Nesse caso, você precisa passar uma referência. Podemos, por exemplo, ter dados muito grandes,  ou coisas que você talvez pretenda ter apenas uma instância \(como conexões a banco de dados\).

### nil

* Ponteiros podem ser `nil`.
* Quando uma função retorna um ponteiro para algo, você precisa ter certeza de verificar se ele é `nil` ou isso vai gerar uma exceção em tempo de execução, já que o compilador não te consegue te ajudar nesses casos.
* Útil para quando você quer descrever um valor que pode estar faltando.
  
### Erros

* Erros são a forma de sinalizar falhas na execução de uma função/método.
* Analisando nossos testes, concluímos que buscar por uma string em um erro poderia resultar em um teste não muito confiável. Então, refatoramos para usar um valor significativo, que resultou em um código mais fácil de ser testado e concluímos que também seria mais fácil para usuários de nossa API.
* Este não é o fim do assunto de tratamento de erros. Você pode fazer coisas mais sofisticadas, mas esta é apenas uma introdução. Capítulos posteriores vão abordar mais estratégias.
* [Não somente verifique os erros, trate-os graciosamente](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

### Crie novos tipos a partir de existentes

* Útil para adicionar domínios mais específicos a valores
* Permite implementar interfaces

Ponteiros e erros são uma grande parte de escrita em Go que você precisa estar confortável. Por sorte, _na maioria das vezes_ o compilador irá ajudar se você fizer algo errado. É só tomar um tempinho lendo a mensagem de erro.

