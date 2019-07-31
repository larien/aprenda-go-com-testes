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

## Escreva o teste primeiro

```go
func TestAdiciona(t *testing.T) {
	t.Run("palavra nova", func(t *testing.T) {
		dicionario := Dicionario{}
		palavra := "teste"
		definicao := "isso é apenas um teste"

		err := dicionario.Adiciona(palavra, definicao)

		comparaErro(t, err, nil)
		comparaDefinicao(t, dicionario, palavra, definicao)
	})

	t.Run("palavra existente", func(t *testing.T) {
		palavra := "teste"
		definicao := "isso é apenas um teste"
		dicionario := Dicionario{palavra: definicao}
		err := dicionario.Add(palavra, "teste novo")

		comparaErro(t, err, ErrPalavraExistente)
		comparaDefinicao(t, dicionario, palavra, definicao)
	})
}
```

Para esse teste, fizemos `Adiciona` devolver um erro, que estamos validando com uma nova variável de erro, `ErrPalavraExistente`. Também modificamos o teste anterior para verificar um erro `nil`.

## Execute o teste

Agora o compilador vai falhar porque não estamos devolvendo um valor para `Adiciona`.

```bash
./dicionario_test.go:30:13: dicionario.Adiciona(palavra, definicao) used as value
./dicionario_test.go:41:13: dicionario.Adiciona(palavra, "teste novo") used as value
```

`usado como valor`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhado

Em `dicionario.go`

```go
var (
    ErrNaoEncontrado = errors.New("não foi possível encontrar a palavra que você procura")
    ErrPalavraExistente = errors.New("não é possível adicionar a palavra pois ela já existe")
)

func (d Dicionario) Adiciona(palavra, definicao string) error {
    d[palavra] = definicao
    return nil
}
```

Agora temos mais dois erros. Ainda estamos modificando o valor e retornando um erro `nil`.

Now we get two more errors. We are still modifying the value, and returning a `nil` error.

```bash
dicionario_test.go:43: resultado erro '%!s(<nil>)', esperado 'não é possível adicionar a palavra pois ela já existe'
dicionario_test.go:44: resultado 'teste novo', esperado 'isso é apenas um teste'
```

## Escreva código o suficiente para fazer o teste passar

```go
func (d Dicionario) Adiciona(palavra, definicao string) error {
    _, err := d.Busca(palavra)
    switch err {
    case ErrNaoEncontrado:
        d[palavra] = definicao
    case nil:
        return ErrPalavraExistente
    default:
        return err

    }

    return nil
}
```

Aqui estamos usando a declaração `switch` para coincidir com o erro. Usar o `switch` dessa forma dá uma segurança a mais, no caso de `Busca` retornar um erro diferente de `ErrNaoEncontrado`.

## Refatoração

Não temos muito o que refatorar, mas já que nossos erros estão aumentando, podemos fazer algumas modificações.

```go
const (
    ErrNaoEncontrado = ErrDicionario("não foi possível encontrar a palavra que você procura")
    ErrPalavraExistente = ErrDicionario("não é possível adicionar a palavra pois ela já existe")
)

type ErrDicionario string

func (e ErrDicionario) Error() string {
    return string(e)
}
```

Tornamos os erros constantes; para isso, tivemos que criar nosso próprio tipo `ErrDicionario` que implementa a interface `error`. Você pode ler mais sobre nesse [artigo excelente escrito por Dave Cheney](https://dave.cheney.net/2016/04/07/constant-errors) (em inglês). Resumindo, isso torna os erros mais reutilizáveis e imutáveis.

Agora, vamos criar uma função que `Atualiza` a definição de uma palavra.

## Escreva o teste primeiro

```go
func TestUpdate(t *testing.T) {
	palavra := "teste"
	definicao := "isso é apenas um teste"
	dicionario := Dicionario{palavra: definicao}
	novaDefinicao := "nova definição"

	dicionario.Atualiza(palavra, novaDefinicao)

	comparaDefinicao(t, dicionario, palavra, novaDefinicao)
}
```

`Atualiza` é bem parecido com `Adiciona` e será nossa próxima implementação.

## Execute o teste

```bash
./dicionario_test.go:53:2: dicionario.Atualiza undefined (type Dicionario has no field or method Atualiza)
```

`dicionario.Atualiza não definido (tipo Dicionario não tem nenhum campo ou método chamado Atualiza`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Já sabemos como lidar com um erro como esse. Precisamos definir nossa função.

```go
func (d Dicionario) Atualiza(palavra, definicao string) {}
```

Feito isso, somos capazes de ver o que precisamos para mudar a definição da palavra.

With that in place, we are able to see that we need to change the definition of the word.

```bash
dicionario_test.go:55: resultado 'isso é apenas um teste', esperado 'nova definição'
```

## Escreva código o suficiente para fazer o teste passar

Já vimos como fazer essa implementação quando corrigimos o problema com `Adiciona`. Logo, vamos implementar algo bem parecido a `Adiciona`.

```go
func (d Dicionario) Atualiza(palavra, definicao string) {
	d[palavra] = definicao
}
```

Não há refatoração necessária, já que foi uma mudança simples. No entanto, agora temos o mesmo problema com `Adiciona`. Se passarmos uma palavra nova, `Atualiza` vai adicioná-la no dicionário.

## Escreva o teste primeiro

```go
    t.Run("palavra existente", func(t *testing.T) {
        palavra := "teste"
        definicao := "isso é apenas um teste"
        novaDefinicao := "nova definição"
        dicionario := Dicionario{palavra: definicao}
        err := dicionario.Atualiza(palavra, novaDefinicao)

        comparaErro(t, err, nil)
        comparaDefinicao(t, dicionario, palavra, novaDefinicao)
    })

    t.Run("palavra nova", func(t *testing.T) {
        palavra := "teste"
        definicao := "isso é apenas um teste"
        dicionario := Dicionario{}

        err := dicionario.Atualiza(palavra, definicao)

        comparaErro(t, err, ErrPalavraInexistente)
    })
```

Criamos um outro tipo de erro para quando a palavra não existe. Também modificamos o `Atualiza` para retornar um valor `error`.

## Execute o teste

```bash
./dicionario_test.go:53:16: dicionario.Atualiza(palavra, "teste novo") used as value
./dicionario_test.go:64:16: dicionario.Atualiza(palavra, definicao) used as value
./dicionario_test.go:66:23: undefined: ErrPalavraInexistente
```

Agora recebemos três erros, mas sabemos como lidar com eles.

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

```go
const (
    ErrNaoEncontrado = ErrDicionario("não foi possível encontrar a palavra que você procura")
    ErrPalavraExistente = ErrDicionario("não é possível adicionar a palavra pois ela já existe")
    ErrPalavraInexistente = ErrDicionario("não foi possível atualizar a palavra pois ela não existe")
)

func (d Dicionario) Atualiza(palavra, definicao string) error {
    d[palavra] = definicao
    return nil
}
```

Adicionamos nosso próprio tipo erro e retornamos um erro `nil`.

Com essas mudanças, agora temos um erro muito mais claro:

```bash
dicionario_test.go:66: resultado erro '%!s(<nil>)', esperado 'não foi possível atualizar a palavra pois ela não existe'
```

## Escreva código o suficiente para fazer o teste passar

```go
func (d Dicionario) Atualiza(palavra, definicao string) error {
	_, err := d.Busca(palavra)
	switch err {
	case ErrNaoEncontrado:
		return ErrPalavraInexistente
	case nil:
		d[palavra] = definicao
	default:
		return err

	}

	return nil
}
```

Essa função é quase idêntica à `Adiciona`, com exceção de que trocamos quando atualizamos o `dicionario` e quando retornamos um erro.

### Nota sobre a declaração de um novo erro para Atualiza

Poderíamos reutilizar `ErrNaoEncontrado` e não criar um novo erro. No entanto, geralmente é melhor ter um erro preciso para quando uma atualização falhar.

Ter erros específicos te dá mais informação sobre o que deu errado. Segue um exemplo em uma aplicação web:

> Você pode redirecionar o usuário quando o `ErrNaoEncontrado` é encontrado, mas mostrar uma mensagem de erro só quando `ErrPalavraInexistente` é encontrado.

Agora, vamos criar uma função que `Deleta` uma palavra no dicionário.

## Escreva o teste primeiro

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

## Escreva código o suficiente para fazer o teste passar

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
