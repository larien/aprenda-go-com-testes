# Maps

[**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/maps)

Em [arrays e slices](primeiros-passos-com-go/arrays-e-slices.md), vimos como armazenas valores em ordem. Agora, vamos descobrir uma forma de armazenar itens por uma `chave` (chave) e procurar por ela rapidamente.

Maps te permitem armazenar itens de forma parecida com a de um dicionário. Você pode pensar na `chave` como a palavra e o `valor` como a definição. E que forma melhor de aprender sobre Maps do que criar seu próprio dicionário?

Primeiro, vamos presumir que já temos algumas palavras com suas definições no dicionário. Se procurarmos por uma palavra, ele deve retornar sua definição.

## Escreva o teste primeiro

Em `dicionario_test.go`

```go
package main

import "testing"

func TestBusca(t *testing.T) {
    dicionario := map[string]string{"teste": "isso é apenas um teste"}

    resultado := Busca(dictionary, "teste")
    esperado := "isso é apenas um teste"

    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s', dado '%s'", resultado, esperado, "test")
    }
}
```

Declarar um Map é bem parecido com um array. A diferença é que começa com a palavra-chave `map` e requer dois tipos. O primeiro é o tipo da chave, que é escrito dentro de `[]`. O segundo é o tipo do valor, que vai logo após o `[]`.

O tipo da chave é especial. Só pode ser um tipo comparável, porque sem a habilidade de dizer se duas chaves são iguais, não temos como certificar de que estamos obtendo o valor correto. Tipos comparáveis são explicados com detalhes na [especificação da linguagem](https://golang.org/ref/spec#Comparison_operators).

O tipo do valor, por outro lado, pode ser o tipo que quiser. Pode até ser outro map.

O restante do teste já deve ser familiar para você.

## Execute o teste

Ao executar `go test`, o compilador vai falhar com `./dicionario_test.go:8:9: undefined: Busca`.

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhado

Em `dicionario.go`:

```go
package main

func Busca(dicionario map[string]string, palavra string) string {
    return ""
}
```

Agora seu teste vai falhar com uma _mensagem de erro clara_:

`dicionario_test.go:12: resultado '', esperado 'isso é apenas um teste', dado 'teste'`.

## Escreva código o suficiente para fazer o teste passar

```go
func Busca(dicionario map[string]string, palavra string) string {
    return dicionario[palavra]
}
```

Obter um valor de um Map é igual a obter um valor de um Array: `map[chave]`.

## Refatoração

```go
func TestBusca(t *testing.T) {
    dicionario := map[string]string{"teste": "isso é apenas um teste"}

    resultado := Busca(dicionario, "teste")
    esperado := "isso é apenas um teste"

    compararStrings(t, resultado, esperado)
}

func compararStrings(t *testing.T, resultado, esperado string) {
	t.Helper()

	if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s', dado '%s'", resultado, esperado, "test")
    }
}
```

Decidi criar um helper `compararStrings` para tornar a implementação mais genérica.

### Usando um tipo personalizado

Podemos melhorar o uso do nosso dicionário criando um novo tipo baseado no map e fazendo a `Busca` virar um método.

Em `dicionario_test.go`:

```go
func TestBusca(t *testing.T) {
    dicionario := Dictionary{"teste": "isso é apenas um teste"}

    resultado := dictionary.Busca("teste")
    esperado := "isso é apenas um teste"

    compararStrings(t, resultado, esperado)
}
```

Começamos a usar o tipo `Dicionario`, que ainda não definimos. Depois disso, chamamos `Busca` da instância de `Dicionario`.

Não precisamos mudar o `comparaStrings`.

Em `dicionario.go`:

```go
type Dicionario map[string]string

func (d Dicionario) Busca(palavra string) string {
	return d[palavra]
}
```

Aqui criamos um tipo `Dicionario` que trabalha em cima da abstração de `map`. Com o tipo personalizado definido, podemos criar o método `Busca`.

## Escreva o teste primeiro

A busca básica foi bem fácil de implementar, mas o que acontece se passarmos uma palavra que não está no nosso dicionário?

Como o código está agora, não recebemos nada de volta. Isso é bom porque o programa continua a ser executado, mas há uma abordagem melhor. A função pode reportar que a palavra não está no dicionário. Dessa forma, o usuário não fica se perguntando se a palavra não existe ou se apenas não existe definição para ela (isso pode não parecer tão útil para um dicionário. No entanto, é um caso que pode ser essencial em outros casos de usos).

```go
func TestBusca(t *testing.T) {
	dicionario := Dicionario{"teste": "isso é apenas um teste"}

	t.Run("palavra conhecida", func(t *testing.T) {
		resultado, _ := dicionario.Busca("teste")
		esperado := "isso é apenas um teste"

		comparaStrings(t, resultado, esperado)
	})

	t.Run("palavra desconhecida", func(t *testing.T) {
		_, resultado := dicionario.Busca("desconhecida")

		if err == nil {
            t.Fatal("é esperado que um erro seja obtido.")
        }
	})
}
```

A forma de lidar com esse caso no Go é retornar um segundo argumento que é do tipo `Error`.

Erros podem ser convertidos para uma string com o método `.Error()`, o que podemos fazer quando passarmos para a asserção. Também estamos protegendo o `comparaStrings` com `if` para certificar que não chamemos `.Error()` quando o erro for `nil`.

## Execute o teste

Isso não vai compilar.

This does not compile

`./dictionary_test.go:18:10: assignment mismatch: 2 variables but 1 values`

`incompatibilidade de atribuição: 2 variáveis, mas 1 valor`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhado

```go
func (d Dicionario) Busca(palavra string) (string, error) {
    return d[palavra], nil
}
```

Agora seu teste deve falhar com uma mensagem de erro muito mais clara.

`dictionary_test.go:22: expected to get an error.`

`erro esperado.`

## Escreva código o suficiente para fazer o teste passar

```go
func (d Dicionario) Busca(palavra string) (string, error) {
    definicao, existe := d[palavra]
    if !existe {
        return "", errors.New("não foi possível encontrar a palavra que você procura")
    }

    return definicao, nil
}
```

Para fazê-lo passar, estamos usando uma propriedade interessante ao percorrer o map. Ele pode retornar dois valores. O segundo valor é uma boleana que indica se a chave foi encontrada com sucesso.

Essa propriedade nos permite diferenciar entre uma palavra que não existe e uma palavra que só não tem uma definição.

## Refatoração

```go
var ErrNaoEncontrado = errors.New("não foi possível encontrar a palavra que você procura")

func (d Dicionario) Busca(palavra string) (string, error) {
    definicao, existe := d[palavra]
    if !existe {
        return "", ErrNaoEncontrado
    }

    return definicao, nil
}
```

Podemos nos livrar do "erro mágico" na nossa função de `Busca` extraindo-o para dentro de uma variável. Isso também nos permite ter um teste melhor.

```go
t.Run("palavra desconhecida", func(t *testing.T) {
    _, resultado := dicionario.Busca("desconhecida")

    comparaErro(t, resultado, ErrNotFound)
})

func comparaErro(t *testing.T, resultado, esperado error) {
    t.Helper()

    if resultado != esperado {
        t.Errorf("resultado erro '%s', esperado '%s'", resultado, esperado)
    }
}
```

Conseguimos simplificar nosso teste criando um novo helper e começando a usar nossa variável `ErrNaoEncontrado` para que nosso teste não falhe se mudarmos o texto do erro no futuro.

## Escreva o teste primeiro

Temos uma ótima maneira de buscar no dicionário. No entanto, não temos como adicionar novas palavras nele.

```go
func TestAdiciona(t *testing.T) {
    dicionario := Dicionario{}
    dicionario.Adiciona("teste", "isso é apenas um teste")

    esperado := "isso é apenas um teste"
    resultado, err := dicionario.Busca("teste")
    if err != nil {
        t.Fatal("não foi possível encontrar a palavra adicionada:", err)
    }

    if esperado != resultado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }
}
```

Nesse teste, estamos utilizando nossa função `Busca` para tornar a validação do dicionário um pouco mais fácil.

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhado

Em `dicionario.go`

```go
func (d Dicionario) Adiciona(palavra, definicao string) {
}
```

Agora seu teste deve falhar.

```bash
dicionario_test.go:31: deveria ter encontrado palavra adicionada: não foi possível encontrar a palavra que você procura
```

## Escreva código o suficiente para fazer o teste passar

```go
func (d Dicionario) Adiciona(palavra, definicao string) {
	d[palavra] = definicao
}
```

Adicionar coisas a um map também é bem semelhante a um array. Você só precisar especificar uma chave e definir qual é seu valor.

### Tipos Referência

Uma propriedade interessante dos maps é que você pode modificá-los sem passá-los como ponteiro. Isso é porque o `map` é um tipo referência. Isso significa que ele contém uma referência à estrutura de dado subjacente, assim como um ponteiro. A estrutura de data subjacente é uma `tabela de dispersão` ou `mapa de hash`, e você pode ler mais sobre [aqui](https://pt.wikipedia.org/wiki/Tabela_de_dispers%C3%A3o).

É muito bom referenciar um map, porque não importa o tamanho do map, só vai haver uma cópia.

Um conceito que os tipos referência apresentam é que maps podem ser um valor `nil`. Um map `nil` se comporta como um map vazio durante a leitura,mas tentar inserir coisas em um map `nil` gera um panic em tempo de execução. Você pode saber mais sobre maps [aqui](https://blog.golang.org/go-maps-in-action) (em inglês).

Além disso, você nunca deve inicializar um map vazio:

```go
var m map[string]string
```

Ao invés disso, você pode inicializar um map vazio como fizemos lá em cima, ou usando a palavra-chave `make` para criar um map para você:` keyword to create a map for you:

```go
dicionario = map[string]string{}

// OU

dicionario = make(map[string]string)
```

Ambas as abordagens criam um `hash map` vazio e apontam um `dicionario` para ele. Assim, nos certificamos que você nunca vai obter um panic em tempo de execução.

## Refatoração

Não há muito para refatorar na nossa implementação, mas podemos simplificar o teste um pouco.

```go
func TestAdiciona(t *testing.T) {
	dicionario := Dicionario{}
	palavra := "teste"
	definicao := "isso é apenas um teste"

	dicionario.Adiciona(palavra, definicao)

	comparaDefinicao(t, dicionario, palavra, definicao)
}

func comparaDefinicao(t *testing.T, dicionario Dicionario, palavra, definicao string) {
	t.Helper()

	resultado, err := dicionario.Busca(palavra)
	if err != nil {
		t.Fatal("deveria ter encontrado palavra adicionada:", err)
	}

	if definicao != resultado {
		t.Errorf("resultado '%s',  esperado '%s'", resultado, definicao)
	}
}
```

Criamos variáveis para palavra e definição e movemos a comparação da definição para sua própria função auxiliar.

Nosso `Adiciona` está bom. No entanto, não consideramos o que acontece quando o valor que estamos tentando adicionar já existe!

O map não vai mostrar um erro se o valor já existe. Ao invés disso, elas vão sobrescrever o valor com o novo recebido. Isso pode ser conveniente na prática, mas torna o nome da nossa função muito menos preciso. `Adiciona` não deve modificar valores existentes. Só deve adicionar palavras novas ao nosso dicionário.

## Write the test first

```go
func TestAdd(t *testing.T) {
    t.Run("new word", func(t *testing.T) {
        dictionary := Dictionary{}
        word := "test"
        definition := "this is just a test"

        err := dictionary.Add(word, definition)

        assertError(t, err, nil)
        assertDefinition(t, dictionary, word, definition)
    })

    t.Run("existing word", func(t *testing.T) {
        word := "test"
        definition := "this is just a test"
        dictionary := Dictionary{word: definition}
        err := dictionary.Add(word, "new test")

        assertError(t, err, ErrWordExists)
        assertDefinition(t, dictionary, word, definition)
    })
}
```

For this test, we modified `Add` to return an error, which we are validating against a new error variable, `ErrWordExists`. We also modified the previous test to check for a `nil` error.

## Try to run test

The compiler will fail because we are not returning a value for `Add`.

```text
./dictionary_test.go:30:13: dictionary.Add(word, definition) used as value
./dictionary_test.go:41:13: dictionary.Add(word, "new test") used as value
```

## Write the minimal amount of code for the test to run and check the output

In `dictionary.go`

```go
var (
    ErrNotFound   = errors.New("could not find the word you were looking for")
    ErrWordExists = errors.New("cannot add word because it already exists")
)

func (d Dictionary) Add(word, definition string) error {
    d[word] = definition
    return nil
}
```

Now we get two more errors. We are still modifying the value, and returning a `nil` error.

```text
dictionary_test.go:43: got error '%!s(<nil>)' want 'cannot add word because it already exists'
dictionary_test.go:44: got 'new test' want 'this is just a test'
```

## Write enough code to make it pass

```go
func (d Dictionary) Add(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        d[word] = definition
    case nil:
        return ErrWordExists
    default:
        return err
    }

    return nil
}
```

Here we are using a `switch` statement to match on the error. Having a `switch` like this provides an extra safety net, in case `Search` returns an error other than `ErrNotFound`.

## Refactor

We don't have too much to refactor, but as our error usage grows we can make a few modifications.

```go
const (
    ErrNotFound   = DictionaryErr("could not find the word you were looking for")
    ErrWordExists = DictionaryErr("cannot add word because it already exists")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
    return string(e)
}
```

We made the errors constant; this required us to create our own `DictionaryErr` type which implements the `error` interface. You can read more about the details in [this excellent article by Dave Cheney](https://dave.cheney.net/2016/04/07/constant-errors). Simply put, it makes the errors more reusable and immutable.

Next, let's create a function to `Update` the definition of a word.

## Write the test first

```go
func TestUpdate(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{word: definition}
    newDefinition := "new definition"

    dictionary.Update(word, newDefinition)

    assertDefinition(t, dictionary, word, newDefinition)
}
```

`Update` is very closely related to `Add` and will be our next implementation.

## Try and run the test

```text
./dictionary_test.go:53:2: dictionary.Update undefined (type Dictionary has no field or method Update)
```

## Write minimal amount of code for the test to run and check the failing test output

We already know how to deal with an error like this. We need to define our function.

```go
func (d Dictionary) Update(word, definition string) {}
```

With that in place, we are able to see that we need to change the definition of the word.

```text
dictionary_test.go:55: got 'this is just a test' want 'new definition'
```

## Write enough code to make it pass

We already saw how to do this when we fixed the issue with `Add`. So let's implement something really similar to `Add`.

```go
func (d Dictionary) Update(word, definition string) {
    d[word] = definition
}
```

There is no refactoring we need to do on this since it was a simple change. However, we now have the same issue as with `Add`. If we pass in a new word, `Update` will add it to the dictionary.

## Write the test first

```go
t.Run("existing word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    newDefinition := "new definition"
    dictionary := Dictionary{word: definition}

    err := dictionary.Update(word, newDefinition)

    assertError(t, err, nil)
    assertDefinition(t, dictionary, word, newDefinition)
})

t.Run("new word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{}

    err := dictionary.Update(word, definition)

    assertError(t, err, ErrWordDoesNotExist)
})
```

We added yet another error type for when the word does not exist. We also modified `Update` to return an `error` value.

## Try and run the test

```text
./dictionary_test.go:53:16: dictionary.Update(word, "new test") used as value
./dictionary_test.go:64:16: dictionary.Update(word, definition) used as value
./dictionary_test.go:66:23: undefined: ErrWordDoesNotExists
```

We get 3 errors this time, but we know how to deal with these.

## Write the minimal amount of code for the test to run and check the failing test output

```go
const (
    ErrNotFound         = DictionaryErr("could not find the word you were looking for")
    ErrWordExists       = DictionaryErr("cannot add word because it already exists")
    ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

func (d Dictionary) Update(word, definition string) error {
    d[word] = definition
    return nil
}
```

We added our own error type and are returning a `nil` error.

With these changes, we now get a very clear error:

```text
dictionary_test.go:66: got error '%!s(<nil>)' want 'cannot update word because it does not exist'
```

## Write enough code to make it pass

```go
func (d Dictionary) Update(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        return ErrWordDoesNotExist
    case nil:
        d[word] = definition
    default:
        return err
    }

    return nil
}
```

This function looks almost identical to `Add` except we switched when we update the `dictionary` and when we return an error.

### Note on declaring a new error for Update

We could reuse `ErrNotFound` and not add a new error. However, it is often better to have a precise error for when an update fails.

Having specific errors gives you more information about what went wrong. Here is an example in a web app:

> You can redirect the user when `ErrNotFound` is encountered, but display an error message when `ErrWordDoesNotExist` is encountered.

Next, let's create a function to `Delete` a word in the dictionary.

## Write the test first

```go
func TestDelete(t *testing.T) {
    word := "test"
    dictionary := Dictionary{word: "test definition"}

    dictionary.Delete(word)

    _, err := dictionary.Search(word)
    if err != ErrNotFound {
        t.Errorf("Expected '%s' to be deleted", word)
    }
}
```

Our test creates a `Dictionary` with a word and then checks if the word has been removed.

## Try to run the test

By running `go test` we get:

```text
./dictionary_test.go:74:6: dictionary.Delete undefined (type Dictionary has no field or method Delete)
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
func (d Dictionary) Delete(word string) {

}
```

After we add this, the test tells us we are not deleting the word.

```text
dictionary_test.go:78: Expected 'test' to be deleted
```

## Write enough code to make it pass

```go
func (d Dictionary) Delete(word string) {
    delete(d, word)
}
```

Go has a built-in function `delete` that works on maps. It takes two arguments. The first is the map and the second is the key to be removed.

The `delete` function returns nothing, and we based our `Delete` method on the same notion. Since deleting a value that's not there has no effect, unlike our `Update` and `Add` methods, we don't need to complicate the API with errors.

## Wrapping up

In this section, we covered a lot. We made a full CRUD \(Create, Read, Update and Delete\) API for our dictionary. Throughout the process we learned how to:

-   Create maps
-   Search for items in maps
-   Add new items to maps
-   Update items in maps
-   Delete items from a map
-   Learned more about errors
    -   How to create errors that are constants
    -   Writing error wrappers
