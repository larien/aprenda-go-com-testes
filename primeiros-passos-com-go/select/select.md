# Select

[**Você pode encontrar todos os códigos desse capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/master/primeiros-passos-com-go/select)

Te pediram para fazer uma função chamada `Corredor` que recebe duas URLs que "competirão" entre si através de uma chamada HTTP GET onde a primeira URL a responder será retornada. Se nenhuma delas responder dentro de 10 segundos a função deve retornar um `erro`. 

Para isso, vamos utilizar:

* `net/http` para chamadas HTTP.
* `net/http/httptest` para nos ajudar a testar.
* goroutines.
* `select` para sincronizar processos.

## Escreva o teste primeiro

Vamos começar com algo simples.

```go
func TestCorredor(t *testing.T) {
    URLLenta := "http://www.facebook.com"
    URLRapida := "http://www.quii.co.uk"

    esperado := URLRapida
    resultado := Corredor(URLLenta, urlRapida)

    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }
}
```

Sabemos que não está perfeito e que existem problemas, mas é um bom início. É importante não perder tanto tempo deixando as coisas perfeitas de primeira.

## Execute o teste

`./corredor_test.go:14:9: undefined: Corredor`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

```go
func Corredor(a, b string) (vencedor string) {
    return
}
```

`corredor_test.go:25: resultado '', esperado 'http://www.quii.co.uk'`

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
Para cada URL:

1. Usamos `time.Now()` para marcar o tempo antes de tentarmos pegar a `URL`.
2. Então usamos [`http.Get`](https://golang.org/pkg/net/http/#Client.Get) para tentar capturar os conteúdos da `URL`. Essa função retorna [`http.Response`](https://golang.org/pkg/net/http/#Response) e um `erro`, mas não temos interesse nesses valores.
3. `time.Since` pega o tempo inicial e retorna a diferença na forma de `time.Duration`.

Feito isso, podemos simplesmente comparar as durações e ver qual é mais rápida.

### Problemas

Isso pode ou não fazer com que o teste passe para você. O problema é que estamos acessando sites reais para testar nossa lógica.

Testar códigos que usam HTTP é tão comum que Go tem ferramentas na biblioteca padrão para te ajudar a testá-los.

Nos capítulos de [mock](../mocks/mocks.md) e [injeção de dependências](../injecao-de-dependencia/injecao-de-dependencia.md), falamos sobre como idealmente não queremos depender de serviços externos para testar nosso código, pois:

* Podem ser lentos
* Podem ser inconsistentes
* Não conseguimos testar casos extremos

Na biblioteca padrão, existe um pacote chamado [`net/http/httptest`](https://golang.org/pkg/net/http/httptest/) onde é possível simular um servidor HTTP facilmente.

Vamos alterar nosso teste para usar essas simulações para termos servidores confiáveis para testar sob nosso controle.

```go
func TestCorredor(t *testing.T) {

    servidorLento := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(20 * time.Millisecond)
        w.WriteHeader(http.StatusOK)
    }))

    servidorRapido := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    URLLenta := servidorLento.URL
    URLRapida := servidorRapido.URL

    esperado := URLRapida
    resultado := Corredor(URLLenta, URLRapida)

    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }

    servidorLento.Close()
    servidorRapido.Close()
}
```

A sintaxe pode parecer um pouco complicada, mas não tenha pressa.

`httptest.NewServer` recebe um `http.HandlerFunc` que vamos enviar para uma função _anônima_.

`http.HandlerFunc` é um tipo que se parece com isso: `type HandlerFunc func(ResponseWriter, *Requisicao)`.

Tudo o que assinatura diz é que ela precisa de uma função que recebe um `ResponseWriter` e uma `Requisição`, o que não é novidade para um servidor HTTP.

Acontece que não existe nenhuma mágica aqui, **também é assim que você escreveria um servidor HTTP** __**real**__ **em Go**. A única diferença é que estamos utilizando ele dentro de um `httptest.NewServer` ,o que facilita seu uso em testes por ele encontrar uma porta aberta para escutar e você poder fechá-la quando estiverem concluídos dentro dos próprios testes.

Dentro de nossos dois servidores, fazemos com que um deles tenha um `time.Sleep` quando receber a requisição para torná-lo propositalmente mais lento que o outro. Ambos os servidores, então, devolvem uma resposta `OK` com `w.WriteHeader(http.StatusOK)` a quem realizou a chamada.

Se você rodar o teste novamente, ele definitivamente irá passar e deve ser mais rápido. Brinque com os __sleeps__ para quebrar o teste propositalmente.

## Refatoração

Temos algumas duplicações tanto em nosso código de produção quanto em nosso código de teste.

```go
func Corredor(a, b string) (vencedor string) {
	duracaoA := medirTempoDeResposta(a)
	duracaoB := medirTempoDeResposta(b)

	if duracaoA < duracaoB {
		return a
	}

	return b
}

func medirTempoDeResposta(URL string) time.Duration {
	inicio := time.Now()
	http.Get(URL)
	return time.Since(inicio)
}
```

Essa "enxugada" torna nosso código `Corredor` bem mais legível.

```go
func TestCorredor(t *testing.T) {

	servidorLento := criarServidorComAtraso(20 * time.Millisecond)
	servidorRapido := criarServidorComAtraso(0 * time.Millisecond)

	defer servidorLento.Close()
	defer servidorRapido.Close()

	URLLenta := servidorLento.URL
	URLRapida := servidorRapido.URL

	esperado := URLRapida
	resultado := Corredor(URLLenta, URLRapida)

	if resultado != esperado {
		t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
	}
}

func criarServidorComAtraso(atraso time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(atraso)
		w.WriteHeader(http.StatusOK)
	}))
}
```

Fizemos a refatoração criando nossos servidores falsos numa função chamada `criarServidorComAtraso` para remover alguns códigos desnecessários do nosso teste e reduzir repetições.

### `defer`

Ao chamar uma função com o prefixo `defer`, ela será chamada _após o término da função que a contém_.

Às vezes você vai precisar liberar recursos, como fechar um arquivo ou, como no nosso caso, fechar um servidor para que esse não continue escutando a uma porta.

Utilizamos o `defer` quando queremos que a função seja executada no final de uma função, mas mantendo essa instrução próxima de onde o servidor foi criado para facilitar a vida das pessoas que forem ler o código futuramente.

Nossa refatoração é uma melhoria e uma solução razoável dados os recursos de Go que vimos até aqui, mas podemos deixar essa solução ainda mais simples.

### Sincronizando processos
* Por que estamos testando a velocidade dos sites sequencialmente quando Go é ótimo com concorrência? Devemos conseguir verificar ambos ao mesmo tempo.
* Não nos preocupamos com o _tempo exato de resposta_ das requisições, apenas queremos saber qual retorna primeiro.

Para fazer isso, vamos apresentar uma nova construção chamada `select` que nos ajudará a sincronizar os processos de forma mais fácil e clara.

```go
func Corredor(a, b string) (vencedor string) {
    select {
    case <-ping(a):
        return a
    case <-ping(b):
        return b
    }
}

func ping(URL string) chan bool {
    ch := make(chan bool)
    go func() {
        http.Get(URL)
        ch <- true
    }()
    return ch
}
```

#### `ping`

Definimos a função `ping` que cria um `chan bool` e a retorna.

No nosso caso, não nos _importamos_ com o tipo enviado no canal, _só queremos enviar um sinal_ para dizer que terminamos, então booleanos já servem.

Dentro da mesma função, iniciamos a goroutine que enviará um sinal a esse canal uma vez que a função `http.Get(URL)` tenha sido finalizada.

#### `select`

Se você se lembrar do capítulo de [concorrência](../concorrencia/concorrencia.md), é possível esperar os valores serem enviados a um canal com `variavel := <-ch`. Isso é uma chamada _bloqueante_, pois está aguardando por um valor.

O que o `select` te permite fazer é aguardar _múltiplos_ canais. O primeiro a enviar um valor "vence" e o código abaixo do `case` é executado.

Nós usamos `ping` em nosso `select` para configurar um canal para cada uma de nossas `URL`s. Qualquer um que enviar para esse canal primeiro vai ter seu código executado no `select`, que resultará nessa `URL` sendo retornada \(que consequentemente será a vencedora\).

Após essas mudanças, a intenção por trás de nosso código fica bem clara e sua implementação efetivamente mais simples.

### Limites de tempo

Nosso último requisito era retornar um erro se o `Corredor` demorar mais que 10 segundos.

## Escreva o teste primeiro

```go
t.Run("retorna um erro se o servidor não responder dentro de 10s", func(t *testing.T) {
    servidorA := criarServidorComAtraso(11 * time.Second)
    servidorB := criarServidorComAtraso(12 * time.Second)

    defer servidorA.Close()
    defer servidorB.Close()

    _, err := Corredor(servidorA.URL, servidorB.URL)

    if err == nil {
        t.Error("esperava um erro, mas não obtive um")
    }
})
```

Fizemos nossos servidores de teste demorarem mais que 10s para retornar para exercitar esse cenário e agora estamos esperando que `Corredor` retorne dois valores: a URL vencedora \(que ignoramos nesse teste com `_`\) e um `erro`.

## Execute o teste

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

Alteramos a assinatura de `Corredor` para retornar o vencedor e um `erro`. Retornamos `nil` para nossos casos de sucesso.

O compilador vai reclamar sobre seu _primeiro teste_ esperar apenas um valor, então altere essa linha para `obteve, _ := Corredor(urlLenta, urlRapida)`. Sabendo disso devemos verificar se _não_ obteremos um erro em nosso caso de sucesso.

Se executar isso agora, o teste irá falhar após 11 segundos.

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
`time.After` é uma função muito útil quando usamos `select`. Embora não ocorra em nosso caso, você pode escrever um código que bloqueia para sempre se os canais que o `select` estiver ouvindo nunca retornarem um valor.
 `time.After` retorna um `chan` \(como `ping`\) e te enviará um sinal após a quantidade de tempo definida.

 Para nós isso é perfeito; se `a` ou `b` conseguir retornar teremos um vencedor, mas se chegar a 10 segundos nosso `time.After` nos enviará um sinal e retornaremos um `erro`.

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

* Nos importamos com o tempo excedido em nosso caso de teste de sucesso?
* Os requisitos foram explícitos sobre o tempo limite?

Dado esse conhecimento, vamos fazer uma pequena refatoração para ser simpático aos nossos testes e aos usuários de nosso código.

```go
var limiteDeDezSegundos = 10 * time.Second

func Corredor(a, b string) (vencedor string, erro error) {
    return Configuravel(a, b, limiteDeDezSegundos)
}

func Configuravel(a, b string, tempoLimite time.Duration) (vencedor string, erro error) {
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

Nossos usuários e nosso primeiro teste podem utilizar `Corredor` \(que usa `Configuravel` por baixo dos panos\) e nosso caminho triste pode usar `Configuravel`.

```go
func TestCorredor(t *testing.T) {
	t.Run("compara a velocidade de servidores, retornando o endereço do mais rápido", func(t *testing.T) {
		servidorLento := criarServidorComAtraso(20 * time.Millisecond)
		servidorRapido := criarServidorComAtraso(0 * time.Millisecond)

		defer servidorLento.Close()
		defer servidorRapido.Close()

		URLLenta := servidorLento.URL
		URLRapida := servidorRapido.URL

		esperado := URLRapida
		resultado, err := Corredor(URLLenta, URLRapida)

		if err != nil {
			t.Fatalf("não esperava um erro, mas obteve um %v", err)
		}

		if resultado != esperado {
			t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
		}
	})

	t.Run("retorna um erro se o servidor não responder dentro de 10s", func(t *testing.T) {
		servidor := criarServidorComAtraso(25 * time.Millisecond)

		defer servidor.Close()

		_, err := Configuravel(servidor.URL, servidor.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("esperava um erro, mas não obtive um")
		}
	})
}
```

Adicionei uma verificação final ao primeiro teste para saber se não pegamos um `erro`.

## Resumo

### `select`

* Ajuda você a escutar vários canais.
* Às vezes você pode precisar incluir `time.After` em um de seus `cases` para prevenir que seu sistema fique bloqueado para sempre.

### `httptest`

* Uma forma conveniente de criar servidores de teste para que se tenha testes confiáveis e controláveis.
* Usa as mesmas interfaces que servidores `net/http` reais, o que torna seu sistema consistente e gera menos coisas para você aprender.
