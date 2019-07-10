# Tipos de erro

[**Você pode encontrar todo o código aqui**](https://github.com/quii/learn-go-with-tests/tree/master/q-and-a/error-types)

**Criar os seus próprios tipos de erros pode ser uma forma elegante de organizar seu código, deixando ele fácil de usar e de testar.**

Pedro perguntou no Slack do Gophers:

> Se estou criando um erro como `fmt.Errorf("%s must be foo, got %s", bar, baz)`, existe alguma forma de validar que os valores são iguais sem fazer uma comparação do valor de string?

Vamos criar uma função para explorar essa ideia.

```go
// DumbGetter retorna o corpo da resposta em string da url se conseguir um 200 OK
func DumbGetter(url string) (string, error) {
    res, err := http.Get(url)

    if err != nil {
        return "", fmt.Errorf("erro ao buscar de %s, %v", url, err)
    }

    if res.StatusCode != http.StatusOK {
        return "", fmt.Errorf("não retornou 200 de %s, teve %d", url, res.StatusCode)
    }

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body) // ignorando erros aqui para resumir

    return string(body), nil
}
```

Não é incomum escrever uma função que possa fahar por mais de um motivo, e nós queremos ter certeza de que controlamos cada cenário corretamente.

Como Pedro disse, nós _poderíamos_ escrever um teste para o erro de status assim:

```go
t.Run("quando você não consegue 200 você tem um status de erro"), func(t *testing.T) {

    svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.WriteHeader(http.StatusTeapot)
    }))
    defer svr.Close()

    _, err := DumbGetter(svr.URL)

    if err == nil {
        t.Fatal("expected an error")
    }

    want := fmt.Sprintf("não retornou 200 de %s, teve %d", svr.URL, http.StatusTeapot)
    got := err.Error()

    if got != want {
        t.Errorf(`retornou "%v", queria "%v"`, got, want)
    }
})
```

Este teste cria um servidor que sempre retorno um `StatusTeapot` e então usamos sua URL como argumento para `DumbGetter`, para que então ele gerencie respostas que não são 200 OK corretamente.

## Problemas com essa forma de testar

Este livro busca enfatizar que você _escute seus testes_ e este teste aqui não _se sente_ bem:

* Estamos construindo a mesma string assim como o código de produção faz para testar
* É chata de ler e escrever
* É com a exatidão da mensagem de erro que nós _realmente estamos preocupados_?

O que isso nos diz? A ergonomia do nosso teste pode ser refletida em outra parte do código que tente usar nosso código.

Como um usuário do nosso código reagiria ao tipo específico de erro que retornamos? O melhor que pode fazer é olhar para a string de erro, a qual é extremamente favorável à erros e é horrível de escrever.

## O que devemos fazer

Com TDD temos o benefício de entrar no seguinte pensamento:

> Como _eu_ gostaria de usar esse código?

O que podemos fazer por `DumbGetter` é oferecer uma forma aos usuários para usar o sistema de tipos para entender que tipo de erro aconteceu.

E se `DumbGetter` pudesse retornar algo do tipo:

```go
type BadStatusError struct {
    URL    string
    Status int
}
```

Ao invés de uma string mágica, temos agora _dados_ de verdade para trabalhar.

Vamos mudar nosso teste para refletir essa necessidade:

```go
t.Run("quando você não consegue 200 você tem um status de erro"), func(t *testing.T) {

    svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.WriteHeader(http.StatusTeapot)
    }))
    defer svr.Close()

    _, err := DumbGetter(svr.URL)

    if err == nil {
        t.Fatal("esperado um erro")
    }

    got, isStatusErr := err.(BadStatusError)

    if !isStatusErr {
        t.Fatalf("não foi um BadStatusError, foi %T", err)
    }

    want := BadStatusError{URL:svr.URL, Status:http.StatusTeapot}

    if got != want {
        t.Errorf(`retornou "%v", queria "%v"`, got, want)
    }
})
```

Temos que fazer com que `BadStatusError` implemente a interface de error.

```go
func (b BadStatusError) Error() string {
    return fmt.Sprintf("não retornou 200 de %s, teve %d", b.URL, b.StatusTeapot)
}
```

### O que esse teste faz?

Ao invés de checar a string do erro, estamos fazendo uma [asserção de tipo](https://tour.golang.org/methods/15) no erro para ver se ele é um `BadStatusError`. Isso reflete o nosso desejo pelo _tipo_ exato de erro de forma mais clara. Assumindo que a asserção passe, nós podemos então conferir se as propriedades do erro estão corretas.

Quando executamos o teste, o retorno nos diz que não retornamos o tipo correto de erro:

```text
--- FAIL: TestDumbGetter (0.00s)
t.Run("quando você não consegue 200 você tem um status de erro"), func(t *testing.T) {
        error-types_test.go:56: não foi um BadStatusError, foi *errors.errorString
```

Vamos corrigir `DumbGetter` atualizando nosso código de controle de erro para usar o nosso tipo:

```go
if res.StatusCode != http.StatusOK {
    return "", BadStatusError{URL: url, Status: res.StatusCode}
}
```

Essa mudança tem alguns _reais efeitos positivos_:

* Nossa função `DumbGetter` ficou mais simples, sem se preocupar com a complexidade da string de erro, criando só um `BadStatusError`.
* Nossos testes agora refletem \(e documentam\) o que um usuário do nosso código _poderia_ fazer caso decida ter um tratamento de erro mais sofisticado do que só escrever logs. Basta fazer uma asserção do tipo e então você tem acesso fácil às propriedades do erro.
* Ainda trata-se de "só" um `error` então, caso quiserem, os  usuários podem passar o erro pela pilha de chamadas or gerar um log como qualquer outro `error`.

## Resumindo

Se você está fazendo testes para múltiplas condições de erro, não caia na armadilha de comparar as mensagens de erro.

Isso leva à fragilidade e dificuldade para ler e escrever testes e reflete as dificuldades que usuários do seu código terão caso também precisem começar a fazer coisas de forma diferente, dependendo do tipo de erros que encontrarem.

Certifique-se sempre de que seus testes refletem como _você_ gostaria de usar o seu código, então dessa forma considere criar tipos de erro para encapsular seus tipos de erro. Isso torna o tratamento de tipos de erros mais fácil para os usuários do seu código e também faz o seu tratamento de erro mais simples e fácil de ler.

