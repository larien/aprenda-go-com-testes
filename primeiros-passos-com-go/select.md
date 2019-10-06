# Select

[**Você pode encontrar todos os códigos desse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/select)

Te pediram para fazer uma função chamada `WebsiteRacer` que recebe duas URLs que "competirão" entre elas através de uma chamada HTTP GET em cada, devolvendo a URL que retornar primeiro. Se nenhuma delas retornar dentro de 10 segundos, então a função deve retornar `error` 

Para isso, vamos utilizar:

* `net/http` para chamadas HTTP.
* `net/http/httptest` para nos ajudar a testar.
* goroutines.
* `select` para sincronizar processos.

## Escreva o teste primeiro

Vamos pegar com algo simples pra começar.

```go
func TestRacer(t *testing.T) {
    slowURL := "http://www.facebook.com"
    fastURL := "http://www.quii.co.uk"

    want := fastURL
    got := Racer(slowURL, fastURL)

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
}
```

Sabemos que não está perfeito e que existem problemas, mas é um bom início. É importante não perder tanto tempo deixando as coisas perfeitas de primeira.

## Tente rodar o teste

`./racer_test.go:14:9: undefined: Racer`

## Escreva a menor quantidade de código para rodar o teste e verifique a saída do teste que falhou

```go
func Racer(a, b string) (winner string) {
    return
}
```

`racer_test.go:25: got '', want 'http://www.quii.co.uk'`

## Escreva código suficiente para que o teste passe

```go
func Racer(a, b string) (winner string) {
    startA := time.Now()
    http.Get(a)
    aDuration := time.Since(startA)

    startB := time.Now()
    http.Get(b)
    bDuration := time.Since(startB)

    if aDuration < bDuration {
        return a
    }

    return b
}
```
para cada URL:

1. Usamos `time.Now()` para marcar o tempo antes de tentarmos pegar a `URL`.
2. Então usamos [`http.Get`](https://golang.org/pkg/net/http/#Client.Get) para tentar capturar os conteúdos da `URL`. Essa função retorna [`http.Response`](https://golang.org/pkg/net/http/#Response) e um `error`mas não estamos interessados nesses valores.
3. `time.Since` pega o tempo inicial e retorna a diferença como `time.Duration`

Uma vez que fazemos isso, podemos simplesmente comparar as durações e ver qual é mais rápida.

### Problemas

Isso pode ou não fazer com que o teste passe para você. O problema é que estamos acessando sites reais para testar nossa lógica.

Testar códigos que usam HTTP é tão comum que Go tem ferramentas na biblioteca padrão para te ajudar a testá-los.

Nos capítulos de mock e injeção de dependências, cobrimos como, idealmente, não queremos depender de serviços externos para testar nosso código pois eles podem ser:

* Lentos
* Flaky
* Não podemos testar casos extremos

Na biblioteca padrão, existe um pacote chamado [`net/http/httptest`](https://golang.org/pkg/net/http/httptest/) onde é possível simular facilmente um servidor HTTP.

Vamos alterar nosso teste para usar essas simulações, assim teremos servidores confiáveis para testar sob nosso controle.

```go
func TestRacer(t *testing.T) {

    slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(20 * time.Millisecond)
        w.WriteHeader(http.StatusOK)
    }))

    fastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    slowURL := slowServer.URL
    fastURL := fastServer.URL

    want := fastURL
    got := Racer(slowURL, fastURL)

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }

    slowServer.Close()
    fastServer.Close()
}
```

A sintaxe pode parecer um tanto complicada mas não tenha pressa.

`httptest.NewServer` recebe um `http.HandlerFunc` que estaremos enviando através de uma função _anonima_.

`http.HandlerFunc` é um tipo que se parece com isso: `type HandlerFunc func(ResponseWriter, *Request)`.

Tudo que está realmente dizendo é que ele precisa de uma função que recebe um `ResponseWriter` e uma `Request`, o que não é muito surpreendente para um servidor HTTP.

Acontece que não existe nenhuma mágica aqui, **também é assim que você escreveria um servidor HTTP** __**real**__ **em Go**. A única diferença é que estamos utilizando ele dentro de um `httptest.NewServer` o que nos facilita de usá-lo em testes, por ele encontrar uma porta aberta para escutar e você poder fechá-lo quando os teste estiverem concluídos.

Dentro de nossos dois servidores, fazemos com que o mais lento tenha um `time.Sleep` quando recebe a requisição para fazê-lo mais lento que o outro. Ambos servidores então devolvem uma resposta `OK` com `w.WriteHeader(http.StatusOK)` a quem realizou a chamada.

Se você rodar o teste novamente agora ele definitivamente irá passar e deve ser mais rápido. Brinque com os __sleeps__ para quebrar o teste propositalmente.

## Refatorar

Temos algumas duplicações tanto em nosso código de produção quanto em nosso código de teste.

```go
func Racer(a, b string) (winner string) {
    aDuration := measureResponseTime(a)
    bDuration := measureResponseTime(b)

    if aDuration < bDuration {
        return a
    }

    return b
}

func measureResponseTime(url string) time.Duration {
    start := time.Now()
    http.Get(url)
    return time.Since(start)
}
```

Essa "secagem" torna nosso código `Racer` bem mais legível.

```go
func TestRacer(t *testing.T) {

    slowServer := makeDelayedServer(20 * time.Millisecond)
    fastServer := makeDelayedServer(0 * time.Millisecond)

    defer slowServer.Close()
    defer fastServer.Close()

    slowURL := slowServer.URL
    fastURL := fastServer.URL

    want := fastURL
    got := Racer(slowURL, fastURL)

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(delay)
        w.WriteHeader(http.StatusOK)
    }))
}
```

Nós refatoramos criando nossos servidores falsos numa função chamada `makeDelayedServer`para remover alguns códigos desnecessários do nosso teste e reduzir repetições.

### `defer`

Ao chama uma função com o prefixo `defer`, ela será chamada _ao final da função que a contém_

As vezes você vai precisar liberar recursos, como fechar um arquivo ou, como no nosso caso, fechar um servidor para que esse não continue escutando uma porta.

Você deseja que isso seja executado ao final de uma função, mas mantenha a instrução próxima de onde foi criado o servidor para o benefício de futuros leitores do código.

Nossa refatoração é uma melhoria e uma solução razoável dados os recursos de Go que vimos até aqui, mas podemos deixar essa solução mais simples.

### Sincronizando processos
* Por quê estamos testando a velocidade dos sites sequencialmente quando GO é ótimo com concorrência? Devemos conseguir verificar ambos ao mesmo tempo.
* Não nos preocupamos com o _tempo exato de resposta_ das requisições, apenas queremos saber qual retorna primeiro.

Para fazer isso, vamos introduzir uma nova construção chamada `select` que nos ajudará a sincronizar processes de forma mais fácil e clara.

```go
func Racer(a, b string) (winner string) {
    select {
    case <-ping(a):
        return a
    case <-ping(b):
        return b
    }
}

func ping(url string) chan bool {
    ch := make(chan bool)
    go func() {
        http.Get(url)
        ch <- true
    }()
    return ch
}
```

#### `ping`

Definimos a função `ping` que cria a `chan bool` e a retorna.

No nosso caso, não nos _importamos_ com o tipo enviado no canal, _só queremos enviar um sinal_ para dizer que terminamos então booleans já servem.

Dentro da mesma função, iniciamos a goroutine que enviará um sinal a esse canal uma vez que esteja completa a `http.Get(url)`.

#### `select`

Se você se lembrar do capítulo de concorrência, é possível esperar os valores serem enviados a um canal com `myVar := <-ch`. Isso é um chamada _bloqueante_, pois está aguardando por um valor.

O que o `select`te premite fazer é agardar _multiplos_ canais. O primeiro a enviar um valor "vence" e o código abaixo do `case` é executado.

Nós usamos `ping` em nosso `select` para configurar um canal para cada uma de nossas `URL`s. Qualquer um que escrever a esse canal primeiro vai ter seu código executado no `select`, que resultará nessa `URL`sendo retornada \(e sendo a vencedora\).

Após essas mudanças, a intenção por trás de nosso código é bem clara e sua implementação efetivamente mais simples.

### Timeouts

Nosso último requisito era retornar um erro se o `Racer` demorar mais que 10 segundos.

## Escreva o teste primeiro

```go
t.Run("retorna um erro se o teste não responder dentro de 10s", func(t *testing.T) {
    serverA := makeDelayedServer(11 * time.Second)
    serverB := makeDelayedServer(12 * time.Second)

    defer serverA.Close()
    defer serverB.Close()

    _, err := Racer(serverA.URL, serverB.URL)

    if err == nil {
        t.Error("esperava um erro, mas não consegui um.")
    }
})
```

Fizemos nossos servidores de teste demorarem mais que 10s para retornar para exercitar esse cenário e estamos esperando que `Racer` retorne dois valores agora, a URL vencedora \(que ignoramos nesse teste com `_`\) e um `erro`.

## Tente rodar o teste

`./racer_test.go:37:10: assignment mismatch: 2 variables but 1 values`

## Escreva a menor quantidade de código para rodar o teste e verifique a saída do teste que falhou

```go
func Racer(a, b string) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    }
}
```

Altere a assinatura de `Racer` para retornar o vencedor e um `erro`. Retorne `nil` para nossos casos felizes.

O compilador vai reclamar sobre seu _ primeiro teste_ apenas olhando para um valor, então altere essa linha para `got, _ := Racer(slowURL, fastURL)`, sabendo disso devemos verificar se _não_ obteremos um erro em nosso cenário feliz.

Se executar isso agora, após 11 segundos irá falhar.

```text
--- FAIL: TestRacer (12.00s)
    --- FAIL: TestRacer/returns_an_error_if_a_server_doesn't_respond_within_10s (12.00s)
        racer_test.go:40: expected an error but didn't get one
```

## Escreva código suficiente para que o teste passe

```go
func Racer(a, b string) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(10 * time.Second):
        return "", fmt.Errorf("tempo de espera excedido para %s e %s", a, b)
    }
}
```
`time.After` é uma função muito conveniente quando usamos `select`. Embora não ocorra em nosso caso, você pode escrever um código que bloqueia pra sempre se os canais que você estiver ouvindo nunca retornarem um valor.
 `time.After` retorna um `chan` \(como `ping`\) e te enviará um sinal após a quantidade de tempo definida.

 Para nós isso é perfeito; se `a` ou `b` conseguir retornar vencerá, mas se chegar a 10 segundos então nosso `time.After` nos enviará um sinal e retornaremos um `erro`.

### Testes lentos

O problema que temos é que esses teste demora 10 segundos para rodar. Para uma lógica tão simples, isso nẽo parece ótimo.

O que podemos fazer é deixar esse esgotamento de tempo configurável. Então em nosso teste, podemos ter um tempo bem curto e quando utilizado no mundo real esse tempo ser definido em 10 segundos.

```go
func Racer(a, b string, timeout time.Duration) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("tempo de espera excedido para %s e %s", a, b)
    }
}
```

Nosso teste não irá compilar pois não fonecemos um tempo.

Antes de nos apressar para adicionar esse valor padrão a ambos os testes, vamos _ouvi-los_.

* Nos importamos com o tempo excedido em nosso teste "feliz"?
* Os requisitos foram explícitos sobre o tempo limite?

Dado esse conhecimento, vamos fazer uma pequena refatoração para ser simpático aos nossos testes e aos usuários de nosso código.

```go
var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
    return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("tempo de espera excedido para %s e %s", a, b)
    }
}
```

Nossos usuários e nosso primeiro teste podem utilizar `Racer` \(que usa `ConfigurableRacer` por baixo dos panos\) e nosso caminho triste pode usar `ConfigurableRacer`.

```go
func TestRacer(t *testing.T) {

    t.Run("compara a velocidade de servidores, retornando o endereço do mais rapido", func(t *testing.T) {
        slowServer := makeDelayedServer(20 * time.Millisecond)
        fastServer := makeDelayedServer(0 * time.Millisecond)

        defer slowServer.Close()
        defer fastServer.Close()

        slowURL := slowServer.URL
        fastURL := fastServer.URL

        want := fastURL
        got, err := Racer(slowURL, fastURL)

        if err != nil {
            t.Fatalf("did not expect an error but got one %v", err)
        }

        if got != want {
            t.Errorf("got '%s', want '%s'", got, want)
        }
    })

    t.Run("retorna um erro se o teste não responder dentro de 10s", func(t *testing.T) {
        server := makeDelayedServer(25 * time.Millisecond)

        defer server.Close()

        _, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

        if err == nil {
            t.Error("expected an error but didn't get one")
        }
    })
}
```
Adicionei uma verificação final no primeiro teste para saber se não pegamos um `erro`.

## Resumindo

### `select`

* Ajuda você a esperar em vários canais.
* As vezes você gostaria de incluir `time.After` em um de seus `cases` para prevenir que seu sistema bloqueie pra sempre.

### `httptest`

* Uma forma conveniente de criar servidores de teste para que se tenha testes confiáveis e controláveis.
* Usando as mesmas interfaces que servidores `net/http` reais, o que é consistente e menos para você aprender.