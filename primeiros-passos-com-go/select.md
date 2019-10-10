# Select

[**Você pode encontrar todos os códigos desse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/select)

Te pediram para fazer uma função chamada `Corredor` que recebe duas URLs que "competirão" entre si através de uma chamada HTTP GET onde a primeira URL a responder será retornada. Se nenhuma delas responder dentro de 10 segundos, então a função deve retornar um `erro`. 

Para isso, vamos utilizar:

* `net/http` para chamadas HTTP.
* `net/http/httptest` para nos ajudar a testar.
* goroutines.
* `select` para sincronizar processos.

## Escreva o teste primeiro

Vamos começar com algo simples.

```go
func TestCorredor(t *testing.T) {
    urlLenta := "http://www.facebook.com"
    urlRapida := "http://www.quii.co.uk"

    esperado := urlRapida
    obteve := Corredor(urlLenta, urlRapida)

    if obteve != esperado {
        t.Errorf("obteve '%s', esperado '%s'", obteve, esperado)
    }
}
```

Sabemos que não está perfeito e que existem problemas, mas é um bom início. É importante não perder tanto tempo deixando as coisas perfeitas de primeira.

## Tente rodar o teste

`./corredor_test.go:14:9: undefined: Corredor`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

```go
func Corredor(a, b string) (vencedor string) {
    return
}
```

`corredor_test.go:25: obteve '', esperado 'http://www.quii.co.uk'`

## Escreva código suficiente para que o teste passe

```go
func Corredor(a, b string) (vencedor string) {
    inicioA := time.Now()
    http.Get(a)
    duracaoA := time.Since(inicioA)

    inicioB := time.Now()
    http.Get(b)
    duracaoB := time.Since(inicioB)

    if duracaoA < duracaoB {
        return a
    }

    return b
}
```
para cada URL:

1. Usamos `time.Now()` para marcar o tempo antes de tentarmos pegar a `URL`.
2. Então usamos [`http.Get`](https://golang.org/pkg/net/http/#Client.Get) para tentar capturar os conteúdos da `URL`. Essa função retorna [`http.Response`](https://golang.org/pkg/net/http/#Response) e um `erro` mas não estamos interessados nesses valores.
3. `time.Since` pega o tempo inicial e retorna a diferença como `time.Duration`

Uma vez que fazemos isso, podemos simplesmente comparar as durações e ver qual é mais rápida.

### Problemas

Isso pode ou não fazer com que o teste passe para você. O problema é que estamos acessando sites reais para testar nossa lógica.

Testar códigos que usam HTTP é tão comum que Go tem ferramentas na biblioteca padrão para te ajudar a testá-los.

Nos capítulos de mock e injeção de dependências, cobrimos como idealmente não queremos depender de serviços externos para testar nosso código pois:

* Podem ser lentos
* Podem ser inconsistentes (Flaky)
* Não podemos testar casos extremos

Na biblioteca padrão, existe um pacote chamado [`net/http/httptest`](https://golang.org/pkg/net/http/httptest/) onde é possível simular facilmente um servidor HTTP.

Vamos alterar nosso teste para usar essas simulações, assim teremos servidores confiáveis para testar sob nosso controle.

```go
func TestCorredor(t *testing.T) {

    servidorLento := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(20 * time.Millisecond)
        w.WriteHeader(http.StatusOK)
    }))

    servidorRapido := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    urlLenta := servidorLento.URL
    urlRapida := servidorRapido.URL

    esperado := urlRapida
    obteve := Corredor(urlLenta, urlRapida)

    if obteve != esperado {
        t.Errorf("obteve '%s', esperado '%s'", obteve, esperado)
    }

    servidorLento.Close()
    servidorRapido.Close()
}
```

A sintaxe pode parecer um tanto complicada mas não tenha pressa.

`httptest.NewServer` recebe um `http.HandlerFunc` que estaremos enviando através de uma função _anonima_.

`http.HandlerFunc` é um tipo que se parece com isso: `type HandlerFunc func(ResponseWriter, *Requisicao)`.

Tudo que está realmente dizendo é que ele precisa de uma função que recebe um `ResponseWriter` e uma `Requisição`, o que não é muito surpreendente para um servidor HTTP.

Acontece que não existe nenhuma mágica aqui, **também é assim que você escreveria um servidor HTTP** __**real**__ **em Go**. A única diferença é que estamos utilizando ele dentro de um `httptest.NewServer` o que nos facilita de usá-lo em testes, por ele encontrar uma porta aberta para escutar e você poder fechá-lo quando os teste estiverem concluídos.

Dentro de nossos dois servidores, fazemos com que um deles tenha um `time.Sleep` quando receber a requisição para fazê-lo mais lento que o outro. Ambos servidores então devolvem uma resposta `OK` com `w.WriteHeader(http.StatusOK)` a quem realizou a chamada.

Se você rodar o teste novamente agora ele definitivamente irá passar e deve ser mais rápido. Brinque com os __sleeps__ para quebrar o teste propositalmente.

## Refatorar

Temos algumas duplicações tanto em nosso código de produção quanto em nosso código de teste.

```go
func Corredor(a, b string) (vencedor string) {
    duracaoA := medirTempoResposta(a)
    duracaoB := medirTempoResposta(b)

    if duracaoA < duracaoB {
        return a
    }

    return b
}

func medirTempoResposta(url string) time.Duration {
    inicio := time.Now()
    http.Get(url)
    return time.Since(inicio)
}
```

Essa "enxugada" torna nosso código `Corredor` bem mais legível.

```go
func TestCorredor(t *testing.T) {

    servidorLento := criarServidorDemorado(20 * time.Millisecond)
    servidorRapido := criarServidorDemorado(0 * time.Millisecond)

    defer servidorLento.Close()
    defer servidorRapido.Close()

    urlLenta := servidorLento.URL
    urlRapida := servidorRapido.URL

    esperado := urlRapida
    obteve := Corredor(urlLenta, urlRapida)

    if obteve != esperado {
        t.Errorf("obteve '%s', esperado '%s'", obteve, esperado)
    }
}

func criarServidorDemorado(demora time.Duration) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(demora)
        w.WriteHeader(http.StatusOK)
    }))
}
```

Nós refatoramos criando nossos servidores falsos numa função chamada `criarServidorDemorado`para remover alguns códigos desnecessários do nosso teste e reduzir repetições.

### `defer`

Ao chamar uma função com o prefixo `defer`, ela será chamada _ao final da função que a contém_.

As vezes você vai precisar liberar recursos, como fechar um arquivo ou, como no nosso caso, fechar um servidor para que esse não continue escutando a uma porta.

Você deseja que isso seja executado ao final de uma função, mas mantenha a instrução próxima de onde foi criado o servidor para o benefício de futuros leitores do código.

Nossa refatoração é uma melhoria e uma solução razoável dados os recursos de Go que vimos até aqui, mas podemos deixar essa solução ainda mais simples.

### Sincronizando processos
* Por quê estamos testando a velocidade dos sites sequencialmente quando Go é ótimo com concorrência? Devemos conseguir verificar ambos ao mesmo tempo.
* Não nos preocupamos com o _tempo exato de resposta_ das requisições, apenas queremos saber qual retorna primeiro.

Para fazer isso, vamos introduzir uma nova construção chamada `select` que nos ajudará a sincronizar os processos de forma mais fácil e clara.

```go
func Corredor(a, b string) (vencedor string) {
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

Definimos a função `ping` que cria um `chan bool` e a retorna.

No nosso caso, não nos _importamos_ com o tipo enviado no canal, _só queremos enviar um sinal_ para dizer que terminamos então booleanos já servem.

Dentro da mesma função, iniciamos a goroutine que enviará um sinal a esse canal uma vez que esteja completa a `http.Get(url)`.

#### `select`

Se você se lembrar do capítulo de concorrência, é possível esperar os valores serem enviados a um canal com `variavel := <-ch`. Isso é uma chamada _bloqueante_, pois está aguardando por um valor.

O que o `select` te premite fazer é aguardar _múltiplos_ canais. O primeiro a enviar um valor "vence" e o código abaixo do `case` é executado.

Nós usamos `ping` em nosso `select` para configurar um canal para cada uma de nossas `URL`s. Qualquer um que escrever a esse canal primeiro vai ter seu código executado no `select`, que resultará nessa `URL` sendo retornada \(e consequentemente a vencedora\).

Após essas mudanças a intenção por trás de nosso código é bem clara e sua implementação efetivamente mais simples.

### Limites de tempo

Nosso último requisito era retornar um erro se o `Corredor` demorar mais que 10 segundos.

## Escreva o teste primeiro

```go
t.Run("retorna um erro se o teste não responder dentro de 10s", func(t *testing.T) {
    servidorA := criarServidorDemorado(11 * time.Second)
    servidorB := criarServidorDemorado(12 * time.Second)

    defer servidorA.Close()
    defer servidorB.Close()

    _, err := Corredor(servidorA.URL, servidorB.URL)

    if err == nil {
        t.Error("esperava um erro, mas não obtive um.")
    }
})
```

Fizemos nossos servidores de teste demorarem mais que 10s para retornar para exercitar esse cenário e estamos esperando que `Corredor` retorne dois valores agora, a URL vencedora \(que ignoramos nesse teste com `_`\) e um `erro`.

## Tente rodar o teste

`./corredor_test.go:37:10: assignment mismatch: 2 variables but 1 values`

## Escreva a menor quantidade de código para rodar o teste e verifique a saída do teste que falhou

```go
func Corredor(a, b string) (vencedor string, erro error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    }
}
```

Altere a assinatura de `Corredor` para retornar o vencedor e um `erro`. Retorne `nil` para nossos casos de sucesso.

O compilador vai reclamar sobre seu _ primeiro teste_ apenas olhando para um valor, então altere essa linha para `obteve, _ := Corredor(urlLenta, urlRapida)`, sabendo disso devemos verificar se _não_ obteremos um erro em nosso cenário de sucesso.

Se executar isso agora, após 11 segundos irá falhar.

```text
--- FAIL: TestCorredor (12.00s)
    --- FAIL: TestCorredor/retorna_um_erro_se_o_teste_não_responder_dentro_de_10s (12.00s)
        corredor_test.go:40: esperava um erro, mas não obtive um.
```

## Escreva código o suficiente para fazer o teste passar

```go
func Corredor(a, b string) (vencedor string, erro error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(10 * time.Second):
        return "", fmt.Errorf("tempo limite de espera excedido para %s e %s", a, b)
    }
}
```
`time.After` é uma função muito conveniente quando usamos `select`. Embora não ocorra em nosso caso, você pode escrever um código que bloqueia para sempre se os canais que estiver ouvindo nunca retornarem um valor.
 `time.After` retorna um `chan` \(como `ping`\) e te enviará um sinal após a quantidade de tempo definida.

 Para nós isso é perfeito; se `a` ou `b` conseguir retornar vencerá, mas se chegar a 10 segundos então nosso `time.After` nos enviará um sinal e retornaremos um `erro`.

### Testes lentos

O problema que temos é que esse teste demora 10 segundos para rodar. Para uma lógica tão simples, isso não parece ótimo.

O que podemos fazer é deixar esse esgotamento de tempo configurável. Então, em nosso teste, podemos ter um tempo bem curto e, quando utilizado no mundo real, esse tempo ser definido para 10 segundos.

```go
func Corredor(a, b string, tempoLimite time.Duration) (vencedor string, erro error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(tempoLimite):
        return "", fmt.Errorf("tempo limite de espera excedido para %s e %s", a, b)
    }
}
```

Nosso teste não irá compilar pois não fornecemos um tempo de expiração.

Antes de nos apressar para adicionar esse valor padrão a ambos os testes, vamos _ouvi-los_.

* Nos importamos com o tempo excedido em nosso teste "feliz"?
* Os requisitos foram explícitos sobre o tempo limite?

Dado esse conhecimento, vamos fazer uma pequena refatoração para ser simpático aos nossos testes e aos usuários de nosso código.

```go
var limiteDezSegundos = 10 * time.Second

func Corredor(a, b string) (vencedor string, erro error) {
    return CorredorConfiguravel(a, b, limiteDezSegundos)
}

func CorredorConfiguravel(a, b string, tempoLimite time.Duration) (vencedor string, erro error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(tempoLimite):
        return "", fmt.Errorf("tempo limite de espera excedido para %s e %s", a, b)
    }
}
```

Nossos usuários e nosso primeiro teste podem utilizar `Corredor` \(que usa `CorredorConfiguravel` por baixo dos panos\) e nosso caminho triste pode usar `CorredorConfiguravel`.

```go
func TestCorredor(t *testing.T) {

    t.Run("compara a velocidade de servidores, retornando o endereço do mais rapido", func(t *testing.T) {
        servidorLento := criarServidorDemorado(20 * time.Millisecond)
        servidorRapido := criarServidorDemorado(0 * time.Millisecond)

        defer servidorLento.Close()
        defer servidorRapido.Close()

        urlLenta := servidorLento.URL
        urlRapida := servidorRapido.URL

        esperado := urlRapida
        obteve, err := Corredor(urlLenta, urlRapida)

        if err != nil {
            t.Fatalf("não esperava um erro, mas obteve um %v", err)
        }

        if obteve != esperado {
            t.Errorf("obteve '%s', esperado '%s'", obteve, esperado)
        }
    })

    t.Run("retorna um erro se o servidor não responder dentro de 10s", func(t *testing.T) {
        servidor := criarServidorDemorado(25 * time.Millisecond)

        defer servidor.Close()

        _, err := CorredorConfiguravel(servidor.URL, servidor.URL, 20*time.Millisecond)

        if err == nil {
            t.Error("esperava um erro, mas não obtive um.")
        }
    })
}
```
Adicionei uma verificação final ao primeiro teste para saber se não pegamos um `erro`.

## Resumindo

### `select`

* Ajuda você a esperar em vários canais.
* As vezes você gostaria de incluir `time.After` em um de seus `cases` para prevenir que seu sistema bloqueie para sempre.

### `httptest`

* Uma forma conveniente de criar servidores de teste para que se tenha testes confiáveis e controláveis.
* Usando as mesmas interfaces que servidores `net/http` reais, o que é consistente e menos para você aprender.