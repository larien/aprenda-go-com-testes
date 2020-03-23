# Sync

[**Você pode encontrar todo o código para esse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/primeiros-passos-com-go/sync)

Queremos fazer um contador que é seguro para ser usado concorrentemente.

Vamos começar com um contador não seguro e verificar se seu comportamento funciona em um ambiente com apenas uma _thread_.

Em seguida, vamos testar sua falta de segurança com várias *goroutines* tentando usar o contador dentro dos testes e consertar essa falha.

## Escreva o teste primeiro

Queremos que nossa API nos dê um método para incrementar o contador e depois recupere esse valor.

```go
func TestContador(t *testing.T) {
    t.Run("incrementar o contador 3 vezes resulta no valor 3", func(t *testing.T) {
        contador := Contador{}
        contador.Incrementa()
        contador.Incrementa()
        contador.Incrementa()

        if contador.Valor() != 3 {
		t.Errorf("resultado %d, esperado %d", contador.Valor(), 3)
	}
    })
}
```

## Tente rodar o teste

```text
./sync_test.go:9:14: undefined: Contador
```

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Vamos definir `Contador`.

```go
type Contador struct {

}
```

Tente rodar o teste de novo e ele falhará com o seguinte erro:

```text
./sync_test.go:14:10: contador.Incrementa undefined (type Contador has no field or method Incrementa)
./sync_test.go:18:13: contador.Valor undefined (type Contador has no field or method Valor)
```

Então, para finalmente fazer o teste rodar, podemos definir esses métodos:

```go
func (c *Contador) Incrementa() {

}

func (c *Contador) Valor() int {
    return 0
}
```

Agora tudo deve rodar e falhar:

```text
=== RUN   TestContador
=== RUN   TestContador/incrementar_o_contador_3_vezes_resulta_no_valor_3
--- FAIL: TestContador (0.00s)
    --- FAIL: TestContador/incrementar_o_contador_3_vezes_resulta_no_valor_3 (0.00s)
        sync_test.go:27: resultado 0, esperado 3
```

## Escreva código o suficiente para fazer o teste passar

Isso deve ser simples para _experts_ em Go como nós. Precisamos criar uma instância do tipo Contador e incrementá-lo com cada chamada de `Incrementa`.


```go
type Contador struct {
    valor int
}

func (c *Contador) Incrementa() {
    c.valor++
}

func (c *Contador) Valor() int {
    return c.valor
}
```

## Refatoração

Não há muito o que refatorar, mas já que iremos escrever mais testes em torno do `Contador`, vamos escrever uma pequena função de asserção `verificaContador` para que o teste fique um pouco mais legível.


```go
t.Run("incrementar o contador 3 vezes resulta no valor 3", func(t *testing.T) {
	contador := Contador{}
	contador.Incrementa()
	contador.Incrementa()
	contador.Incrementa()

	verificaContador(t, contador, 3)
})

func verificaContador(t *testing.T, resultado Contador, esperado int) {
	t.Helper()
	if resultado.Valor() != esperado {
		t.Errorf("resultado %d, esperado %d", resultado.Valor(), esperado)
	}
}
```

## Próximos passos

Isso foi muito fácil, mas agora temos um requerimento que é: o programa precisa ser seguro o suficiente para ser usado em um ambiente com acesso concorrente. Vamos precisar criar um teste para exercitar isso.

## Escreva o teste primeiro

```go
t.Run("roda concorrentemente em segurança", func(t *testing.T) {
	contagemEsperada := 1000
	contador := Contador{}

	var wg sync.WaitGroup
	wg.Add(contagemEsperada)

	for i := 0; i < contagemEsperada; i++ {
		go func(w *sync.WaitGroup) {
			contador.Incrementa()
			w.Done()
		}(&wg)
	}
	wg.Wait()

	verificaContador(t, contador, contagemEsperada)
})
```

Isso vai iterar até a nossa `contagemEsperada` e disparar uma *goroutine* para chamar `contador.Incrementa()` a cada iteração.

Estamos usando [`sync.WaitGroup`](https://golang.org/pkg/sync/#WaitGroup), que é uma maneira simples de sincronizar processos concorrentes.

> Um WaitGroup aguarda por uma coleção de *goroutines* terminar seu processamento. A *goroutine* principal faz a chamada para o `Add` definir o número de *goroutines* que serão esperadas. Então, cada uma das *goroutines* é executada e chama `Done` quando termina sua execução. Ao mesmo tempo, `Wait` pode ser usado para bloquear a execução até que todas as *goroutines* tenham terminado.

Ao esperar por `wg.Wait()` terminar sua execução antes de fazer nossas asserções, podemos ter certeza que todas as nossas *goroutines* tentaram chamar o `Incrementa` no `Contador`.

## Tente rodar o teste

```text
=== RUN   TestContador/roda_concorrentemente_em_seguranca
--- FAIL: TestContador (0.00s)
    --- FAIL: TestContador/roda_concorrentemente_em_seguranca (0.00s)
        sync_test.go:26: resultado 939, esperado 1000
FAIL
```

O teste _provavelmente_ vai falhar com um número diferente, mas de qualquer forma demonstra que não roda corretamente quando várias *goroutines* tentam mudar o valor do contador ao mesmo tempo.

## Escreva código o suficiente para fazer o teste passar

Uma solução simples é adicionar uma trava ao nosso `Contador`, um [`Mutex`](https://golang.org/pkg/sync/#Mutex).

> Um Mutex é uma trava de exclusão mútua. O valor zero de um Mutex é um Mutex destravado.

```go
type Contador struct {
    mu sync.Mutex
    valor int
}

func (c *Contador) Incrementa() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.valor++
}
```

Isso significa que qualquer *goroutine* chamando `Incrementa` vai receber a trava em `Contador` se for a primeira chamando essa função. Todas as outras *goroutines* vão ter que esperar por essa primeira execução até que ele esteja `Unlock`, ou destravado, antes de ganhar o acesso à instância de `Contador` alterada pela primeira chamada de função.

Agora, se você rodar o teste novamente, ele deve funcionar porque cada uma das *goroutines* tem que esperar até que seja sua vez antes de fazer alguma mudança.

## Já vi outros exemplos em que o `sync.Mutex` está embutido dentro da struct.

Você pode ver exemplos como esse:

```go
type Contador struct {
    sync.Mutex
    valor int
}
```

Há quem diga que isso torna o código um pouco mais elegante.

```go
func (c *Contador) Incrementa() {
    c.Lock()
    defer c.Unlock()
    c.valor++
}
```

Isso _parece_ legal, mas, apesar de programação ser uma área altamente subjetiva, isso é **feio e errado**.

Às vezes as pessoas esquecem que tipos embutidos significam que os métodos daquele tipo se tornam _parte da interface pública_; e você geralmente não quer isso. Não se esqueçam que devemos ter muito cuidado com as nossas APIs públicas. O momento que tornamos algo público é o momento que outros códigos podem acoplar-se a ele e queremos evitar acoplamentos desnecessários.

Expôr `Lock` e `Unlock` é, no seu melhor caso, muito confuso e, no seu pior caso, potencialmente perigoso para o seu software se quem chamar o seu tipo começar a chamar esses métodos diretamente.

![Demonstração de como um usuário dessa API pode chamar erroneamente o estado da trava](https://i.imgur.com/SWYNpwm.png)

_Isso parece uma péssima ideia._

## Copiando mutexes

Nossos testes passam, mas nosso código ainda é um pouco perigoso.

Se você rodar `go vet` no seu código, deve receber um erro similar ao seguinte:

```text
sync/v2/sync_test.go:16: call of verificaContador copies lock valor: v1.Contador contains sync.Mutex
sync/v2/sync_test.go:39: verificaContador passes lock by valor: v1.Contador contains sync.Mutex
```

Uma rápida olhada na documentação do [`sync.Mutex`](https://golang.org/pkg/sync/#Mutex) nos diz o porquê:

> Um Mutex não deve ser copiado depois do primeiro uso.

Quando passamos nosso `Contador` \(por valor\) para `verificaContador`, ele vai tentar criar uma cópia do mutex.

Para resolver isso, devemos passar um ponteiro para o nosso `Contador`. Vamos, então, mudar a assinatura de `verificaContador`.

```go
func verificaContador(t *testing.T, resultado *Contador, esperado int)
```

Nossos testes não vão mais compilar porque estamos tentando passar um `Contador` ao invés de um `*Contador`. Para resolver isso, é melhor criar um construtor que mostra aos usuários da nossa API que seria melhor ele mesmo não inicializar seu tipo.


```go
func NovoContador() *Contador {
    return &Contador{}
}
```

Use essa função em seus testes quando for inicializar o `Contador`.

## Resumo

Falamos sobre algumas coisas do [pacote sync](https://golang.org/pkg/sync/):

* `Mutex` nos permite adicionar travas aos nossos dados
* `WaitGroup` é uma maneira de esperar as *goroutines* terminarem suas tarefas

### Quando usar travas em vez de *channels* e *goroutines*?

[Anteriormente falamos sobre *goroutines* no primeiro capítulo sobre concorrência](../concorrencia/concorrencia.md)
que nos permite escrever código concorrente e seguro, então por que usar travas?
[A wiki do Go tem uma página dedicada para esse tópico: Mutex ou Channel?](https://github.com/golang/go/wiki/MutexOrChannel)

> Um erro comum de um iniciante em Go é usar demais os *channels* e *goroutines* apenas porque é possível e/ou porque é divertido. Não tenha medo de usar um `sync.Mutex` se for uma solução melhor para o seu problema. Go é pragmático em deixar você escolher as ferramentas que melhor resolvem o seu problema e não te força em um único estilo de código.

Resumindo:

* **Use channels quando for passar a propriedade de um dado**
* **Use mutexes para gerenciar estados**

### go vet

Não se esqueça de usar `go vet` nos seus scripts de _build_ porque ele pode te alertar a respeito de bugs mais sutis no seu código antes que eles atinjam seus pobres usuários.

### Não use códigos embutidos apenas porque é conveniente

* Pense a respeito do efeito que embutir códigos tem na sua API pública.
* Você _realmente_ quer expôr esses métodos e ter pessoas acoplando o código próprio delas a ele?
* Mutexes podem se tornar um desastre de maneiras muito imprevisíveis e estranhas. Imagine um código inesperado destravando um mutex quando não deveria? Isso causaria erros muito estranhos que seriam muito difíceis de encontrar.

