# Sync

[**Você pode encontrar todo o código para esse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/sync)

Queremos fazer um contador que é seguro para se usar concorrentemente.

Vamos começar com um contador inseguro e verificar se seu comportamento funciona em um ambiente com apenas uma thread.

Depois vamos exercitar sua insegurança com múltiplas *goroutines* tentando usar
ele via teste e consertar essa falha.

## Escreva o teste primeiro

Queremos que nossa API nos dê um método para incrementar o contador e depois recuperar esse valor.

```go
func TestContador(t *testing.T) {
    t.Run("incrementar o contador 3 vezes o deixa com valor 3", func(t *testing.T) {
        contador := Contador{}
        contador.Inc()
        contador.Inc()
        contador.Inc()

        if contador.Valor() != 3 {
            t.Errorf("recebido %d, desejado %d", contador.Valor(), 3)
        }
    })
}
```

## Tente rodar o teste

```text
./sync_test.go:9:14: undefined: Contador
```

## Escreva a quantidade mínima de código para o teste rodar e verifique a saída do teste que falhou

Vamos definir `Contador`.

```go
type Contador struct {

}
```

Tente de novo e ele falhará com o seguinte

```text
./sync_test.go:14:10: contador.Inc undefined (type Contador has no field or method Inc)
./sync_test.go:18:13: contador.Valor undefined (type Contador has no field or method Valor)
```

Então, pra finalmente fazer o teste rodar, podemos definir esses métodos

```go
func (c *Contador) Inc() {

}

func (c *Contador) Valor() int {
    return 0
}
```

Agora tudo deve rodar e falhar

```text
=== RUN   TestContador
=== RUN   TestContador/incrementar_o_contador_3_vezes_o_deixa_com_valor_3
--- FAIL: TestContador (0.00s)
    --- FAIL: TestContador/incrementar_o_contador_3_vezes_o_deixa_com_valor_3 (0.00s)
        sync_test.go:27: recebido 0, desejado 3
```

## Escreva código o suficiente para fazer o teste passar

Isso deve ser trivial para experts em Go como nós. Precisamos manter algum
estado do contador no nosso datatype e daí incrementá-lo em cada chamada do
`Inc`.


```go
type Contador struct {
    valor int
}

func (c *Contador) Inc() {
    c.valor++
}

func (c *Contador) Valor() int {
    return c.valor
}
```

## Refatoração

Não há muito o que refatorar, mas, dado que iremos escrever mais testes em
torno do `Contador`, vamos escrever uma pequena função de asserção `assertCount`
para que o teste fique um pouco mais legível.


```go
t.Run("incrementar o contador 3 vezes o deixa com valor 3", func(t *testing.T) {
    contador := Contador{}
    contador.Inc()
    contador.Inc()
    contador.Inc()

    assertContador(t, contador, 3)
})

func assertContador(t *testing.T, recebido Contador, desejado int)  {
    t.Helper()
    if recebido.Valor() != desejado {
        t.Errorf("recebido %d, quero receber %d", recebido.Valor(), desejado)
    }
}
```

## Próximos passos

Isso foi muito fácil, mas agora nós temos uma requisição que é: ele precisa
ser seguro o suficiente para usar em um ambiente com acesso concorrente. Vamos precisar
criar um teste pra exercitar isso.

## Escreva o teste primeiro

```go
t.Run("roda concorrentemente em segurança", func(t *testing.T) {
    contadorDesejado := 1000
    contador := Contador{}

    var wg sync.WaitGroup
    wg.Add(contadorDesejado)

    for i:=0; i<contadorDesejado; i++ {
        go func(w *sync.WaitGroup) {
            contador.Inc()
            w.Done()
        }(&wg)
    }
    wg.Wait()

    assertContador(t, contador, contadorDesejado)
})
```

Isso vai iterar pelo nosso `contadorDesejado` e disparar uma *goroutine* pra chamar `contador.Inc()`.

Nós estamos usando [`sync.WaitGroup`](https://golang.org/pkg/sync/#WaitGroup)
que é uma maneira conveniente de sincronizar processos concorrentes.

> Um WaitGroup aguarda por uma coleção *goroutines* terminar seu processamento.
A *goroutine* principal faz a chamada para o Add definir o número de *goroutines*
que serão esperadas. Então, cada uma das *goroutines* roda novamente e chama
Done quando terminar sua execução. Ao mesmo tempo, Wait pode ser usada para
bloquear até que todas as *goroutines* tenham terminado.

Ao esperar por `wg.Wait()` terminar sua execução antes de fazer nossas asserções, nós
podemos ter certeza que todas as nossas *goroutines* tentaram `Inc` o `Contador`.

## Tente rodar o teste

```text
=== RUN   TestContador/roda_concorrentemente_em_seguranca
--- FAIL: TestContador (0.00s)
    --- FAIL: TestContador/roda_concorrentemente_em_seguranca (0.00s)
        sync_test.go:26: recebido 939, desejado 1000
FAIL
```

O teste _provavelmente_ vai falhar com um número diferente, mas de toda forma
ele demonstra que não roda corretamente quando várias *goroutines* tentam
mudar o valor do contador ao mesmo tempo.

## Escreva código o suficiente para fazer o teste passar

Uma solução simples é adicionar uma trava ao nosso `Contador`, um
[`Mutex`](https://golang.org/pkg/sync/#Mutex)

> Um Mutex é uma trava de exclusão mútua. O valor zero de um Mutex é um Mutex destravado.

```go
type Contador struct {
    mu sync.Mutex
    valor int
}

func (c *Contador) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.valor++
}
```

Isso significa que qualquer *goroutine* chamando `Inc` vai receber a trava em `Contador`
se for a primeira. Todas as outras *goroutines* vão ter que esperar por ele até
que ele esteja `Unlock`, ou destravado, antes de ganhar o acesso.

Agora se você rodar o teste novamente, ele deve funcionar porque cada uma das *goroutines*
têm que esperar até que seja sua vez antes de fazer alguma mudança.

## Eu vi outros exemplos nos quais `sync.Mutex` está embutido dentro da struct.

Você pode ver exemplos como esse

```go
type Contador struct {
    sync.Mutex
    valor int
}
```

É discutido que isso pode tornar o código um pouco mais elegante.

```go
func (c *Contador) Inc() {
    c.Lock()
    defer c.Unlock()
    c.valor++
}
```

Isso _parece_ legal, mas enquanto programação é uma disciplina altamente
subjetiva, isso é **feio e errado**.

Às vezes as pessoas esquecem que tipos embutidos significam que os métodos
daquele tipo se tornam _parte da interface pública_; e você geralmente não
quer isso. Não se esqueçam que devemos ser muito cuidadosos com as nossas APIs
públicas. O momento que tornamos algo público é o momento que outros códigos
podem acoplar-se a ele e nós queremos evitar acoplamentos desnecessários.

Expor `Lock` e `Unlock` é, no seu melhor caso, muito confuso e, no seu pior
caso, potencialmente perigoso para o seu software se quem chamar o seu tipo
começar a chamar esses métodos diretamente.

![Demonstração de como um usuário dessa API pode chamar erroneamente o estado da trava](https://i.imgur.com/SWYNpwm.png)

_Isso parece uma péssima ideia_

## Copiando mutexes

Nossos testes passam, mas nosso código ainda é um pouco perigoso.

Se você rodar `go vet` no seu código, deve receber um erro similar ao seguinte:

```text
sync/v2/sync_test.go:16: call of assertContador copies lock valor: v1.Contador contains sync.Mutex
sync/v2/sync_test.go:39: assertContador passes lock by valor: v1.Contador contains sync.Mutex
```

Uma rápida olhada na documentação do [`sync.Mutex`](https://golang.org/pkg/sync/#Mutex)
nos diz o porquê

> Um Mutex não deve ser copiado depois do primeiro uso.

Quando passamos nosso `Contador` \(por valor\) para `assertContador` ele vai tentar criar uma cópia do mutex.

Para resolver isso, devemos passar um ponteiro para o nosso `Contador`. Vamos então mudar a assinatura de
`assertContador`.

```go
func assertContador(t *testing.T, recebido *Contador, desejado int)
```

Nossos testes não vão mais compilar porque estamos tentando passar um `Contador` em vez de um `*Contador`.
Para resolver isso, é preferível criar um construtor que mostra aos usuários da nossa API que seria
melhor não inicializar o tipo ele mesmo.


```go
func NovoContador() *Contador {
    return &Contador{}
}
```

Use essa função em seus teste quando for inicializar o `Contador`.

## Resumindo

Cobrimos algumas coisas no [pacote sync](https://golang.org/pkg/sync/):

* `Mutex` que nos permite adicionar travas aos nossos dados
* `Waitgroup` que é uma maneira de esperar as *goroutines* terminarem suas tarefas

### Quando usar travas em vez de *channels* e *goroutines*?

[Anteriormente cobrimos *goroutines* no primeiro capítulo de concorrência](concurrency.md)
que nos permite escrever código concorrente e seguro, então por que usar travas?
[A wiki do go tem uma página dedicada para esse tópico: Mutex ou Channel?](https://github.com/golang/go/wiki/MutexOrChannel)

> Um erro comum de um iniciante em Go é usar demais os *channels* e *goroutines* apenas
porque é possível e/ou porque é divertido. Não tenha medo de usar um `sync.Mutex` se
ele se encaixa melhor no seu problema. Go é pragmático em deixar você escolher as
ferramentas que melhor resolvem o seu problema e não te forçar em um único estilo
de código.

Paraphrasing:

* **Use channels quando for passar a propridade de um dado**
* **Use mutexes pra gerenciar estados**

### go vet

Não se esqueça de usar `go vet` nos seus scripts de build porque ele pode te alertar a respeito de bugs mais
sutis no seu código antes que eles atinjam seus pobres usuários.

### Não use códigos embutidos apenas porque é conveniente

* Pense a respeito do efeito que embutir códigos tem na sua API pública.
* Você _realmente_ quer expor esses métodos e ter pessoas acoplando o código
próprio delas a ele?
* No que diz respeito a mutexes, pode ser potencialmente um desastre de maneiras
muito imprevisíveis e estranhas. Imagine algum código obscuro destravando um
mutex quando não deveria; isso causaria erros muito estranhos e que seriam bastante
difíceis de encontrar.

