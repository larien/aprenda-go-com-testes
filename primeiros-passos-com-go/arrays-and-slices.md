# Arrays e slices

[**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/arrays)

Arrays te permitem armazenar diversos elementos do mesmo tipo em uma variável em uma ordem específica.

Quando você tem um array, é muito comum ter que percorrer sobre ele. Logo, vamos usar nosso [recém adquirido conhecimento de `for`](primeiros-passos-com-go/iteracao.md) para criar uma função `Soma`. `Soma` vai receber um array de números e retornar o total.

Também vamos praticar nossas habilidades em TDD.

## Escreva o teste primeiro

Em `soma_test.go`:

```go
package main

import "testing"

func TestSoma(t *testing.T) {

    numeros := [5]int{1, 2, 3, 4, 5}

    resultado := Soma(numeros)
    esperado := 15

    if esperado != recebido {
        t.Errorf("resultado %d, esperado %d, dado %v", resultado, esperado, numeros)
    }
}
```

Arrays têm uma _capacidade fixa_ que é definida quando você declara a variável. Podemos inicializar um array de duas formas:

-   [N]tipo{valor1, valor2, ..., valorN}, como `numeros := [5]int{1, 2, 3, 4, 5}`
-   [...]tipo{valor1, valor2, ..., valorN}, como `numbers := [...]int{1, 2, 3, 4, 5}`

Às vezes é útil também mostrarmos as entradas da função na mensagem de erro. Para isso estamos usando o formatador `%v`, que é o formato "padrão" e funciona bem com arrays.

[Leia mais sobre formatação de strings aqui](https://golang.org/pkg/fmt/)

## Execute o teste

Ao executar `go test`, o compilador vai falhar com `./soma_test.go:10:15: undefined: Soma`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhado

Em `soma.go`:

```go
package main

func Soma(numeros [5]int) int {
    return 0
}
```

Agora seu teste deve falhar com uma _mensagem clara de erro_:

`soma_test.go:13: resultado 0, esperado 15, dado [1 2 3 4 5]`

## Escreva código o suficiente para fazer o teste passar

```go
func Soma(numeros [5]int) int {
    soma := 0
    for i := 0; i < 5; i++ {
        soma += numeros[i]
    }
    return soma
}
```

Para receber o valor de um array em uma posição específica, basta usar a sintaxe `array[índice]`. Nesse caso, estamos usando o `for` para percorrer cada posição do array (que tem 5 posições) e somar cada valor na variável `soma`.

## Refatoração

Vamos apresentar o [`range`](https://gobyexample.com/range) para nos ajudar a limpar o código:

```go
func Soma(numeros [5]int) int {
    soma := 0
    for _, numero := range numeros {
        soma += numero
    }
    return soma
}
```

O `range` permite que você percorra um array. Sempre que é chamado, retorna dois valores: o índice e o valor. Decidimos ignorar o valor índice usando `_` [_blank identifier_](https://golang.org/doc/effective_go.html#blank).

### Arrays e seus tipos

Uma propriedade interessante dos arrays é que seu tamanho é relacionado ao seu tipo. Se tentar passar um `[4]int` dentro da função que espera `[5]int`, ela não vai compilar. Elas são de tipos diferentes e é a mesma coisa que tentar passar uma `string` para uma função que espera um `int`.

Você pode estar pensando que é bastante complicado que arrays tenham tamanho fixo, não é? Só que na maioria das vezes, você provavelmente não vai usá-los!

O Go tem _slices_, em que você não define o tamanho da coleção e, graças a isso, pode ter qualquer tamanho.

O próprio requerimento será somar coleções de tamanhos variados.

## Escreva o teste primeiro

Agora vamos usar o [tipo slice](https://golang.org/doc/effective_go.html#slices) que nos permite ter coleções de qualquer tamanho. A sintaxe é bem parecida com a dos arrays e você só precisa omitir o tamanho quando declará-lo.

`meuSlice := []int{1,2,3}` ao invés de `meuArray := [3]int{1,2,3}`

```go
func TestSoma(t *testing.T) {

    t.Run("coleção de 5 números", func(t *testing.T) {
        numeros := [5]int{1, 2, 3, 4, 5}

        resultado := Soma(numeros)
        esperado := 15

        if resultado != want {
            t.Errorf("resultado %d, want %d, dado %v", resultado, esperado, numeros)
        }
    })

    t.Run("coleção de qualquer tamanho", func(t *testing.T) {
        numeros := []int{1, 2, 3}

        resultado := Soma(numeros)
        esperado := 6

        if resultado != want {
            t.Errorf("resultado %d, esperado %d, dado %v", resultado, esperado, numeros)
        }
    })

}
```

## Execute o teste

Isso não vai compilar.

`./soma_test.go:22:13: cannot use numbers (type []int) as type [5]int in argument to Soma`

`não é possível usar números (tipo []int) como tipo [5]int no argumento para Soma`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhado

Para resolver o problema, podemos:

-   Alterar a API existente mudando o argumento de `Soma` para um slice ao invés de um array.Quando fazemos isso, vamos saber que podemos ter arruinado do dia de alguém, porque nosso _outro_ teste não vai compilar!

-   Criar uma nova função

No nosso caso, mais ninguém está usando nossa função. Logo, ao invés de ter duas funções para manter, vamos usar apenas uma.

```go
func Soma(numeros []int) int {
    soma := 0
    for _, numero := range numeros {
        soma += numero
    }
    return soma
}
```

Se tentar rodar os testes eles ainda não vão compilar. Você vai ter que alterar o primeiro teste e passar um slice ao invés de um array.

## Escreva código o suficiente para fazer o teste passar

Nesse caso, para arrumar os problemas de compilação, tudo o que precisamos fazer aqui é fazer os testes passarem!

## Refatoração

Nós já refatoramos a função `Soma` e tudo o que fizemos foi mudar os arrays para slices. Logo, não há muito o que fazer aqui. Lembre-se que não devemos abandonar nosso código de teste na etapa de refatoração e precisamos fazer alguma coisa aqui.

```go
func TestSoma(t *testing.T) {

    t.Run("coleção de 5 números", func(t *testing.T) {
        numeros := []int{1, 2, 3, 4, 5}

        resultado := Soma(numeros)
        esperado := 15

        if resultado != esperado {
            t.Errorf("resultado %d, esperado %d, dado, %v", resultado, esperado, numeros)
        }
    })

    t.Run("coleção de qualquer tamanho", func(t *testing.T) {
        numeros := []int{1, 2, 3}

        resultado := Soma(numeros)
        esperado := 6

        if resultado != esperado {
            t.Errorf("resultado %d, esperado %d, dado %v", resultado, esperado, numeros)
        }
    })

}
```

É importante questionar o valor dos seus testes. Ter o máximo de testes possível não deve ser o objetivo e sim ter o máximo de _confiança_ possível na sua base de código. Ter testes demais pode se tornar um problema real e só adiciona mais peso na manutenção. **Todo teste tem um custo**.

No nosso caso, dá para perceber que ter dois testes para essa função é redundância. Se funciona para um soice de determindo tamanho, é muito provável que funciona para um slice de qualquer tamanho (dentro desse escopo).

A ferramenta de testes nativa do Go tem a funcionalidade de [cobertura de código](https://blog.golang.org/cover) que te ajuda a identificar áreas do seu código que você não cobriu. Já adianto que ter 100% de cobertura não deve ser seu objetivo; é apenas uma ferramenta para te dar uma ideia da sua cobertura. De qualquer forma, se você aplicar o TDD, é bem provável que chegue bem perto dos 100% de cobertura.

Tente executar `go test -cover` no terminal.

Você deve ver:

```bash
PASS
coverage: 100.0% of statements
```

Agora apague um dos testes e verifique a cobertura novamente.

Agora que estamos felizes com nossa função bem testada, você deve salvar seu trabalho incrível com um commit antes de partir para o próximo desafio.

Precisamos de uma nova função chamada `SomaTudo`, que vai receber uma quantidade variável de slices e devolver um novo slice contendo as somas de cada slice recebido.

Por exemplo:

`SomaTudo([]int{1,2}, []int{0,9})` deve retornar `[]int{3, 9}`

ou

`SomaTudo([]int{1,1,1})` deve retornar `[]int{3}`

## Escreva o teste primeiro

```go
func TestSomaTudo(t *testing.T) {

    resultado := SomaTudo([]int{1,2}, []int{0,9})
    esperado := []int{3, 9}

    if resultado != esperado {
        t.Errorf("resultado %v esperado %v", resultado, esperado)
    }
}
```

## Execute o teste

`./soma_test.go:23:9: undefined: SomaTudo`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhado

Precisamos definir o SomaTudo de acordo com o que nosso teste precisa.

O Go te permite escrever [_funções variádicas_](https://gobyexample.com/variadic-functions) em que a quantidade de argumentos podem variar.

```go
func SomaTudo(numerosParaSomar ...[]int) (somas []int) {
    return
}
```

Pode tentar compilar, mas nossos testes não vão funcionar!

`./soma_test.go:26:9: invalid operation: got != want (slice can only be compared to nil)`

`operação inválida: recebido != esperado (slice só pode ser comparado a nil`

O Go não te deixa usar operadores de igualdade com slices. _É possível_ escrever uma função que percorre cada slice `recebido` e `esperado` e verificar seus valores, mas por praticidade podemos usar o [`reflect.DeepEqual`](https://golang.org/pkg/reflect/#DeepEqual) que é útil para verificar se _duas variáveis_ são iguais.

```go
func TestSomaTudo(t *testing.T) {

    recebido := SomaTudo([]int{1,2}, []int{0,9})
    esperado := []int{3, 9}

    if !reflect.DeepEqual(recebido, esperado) {
        t.Errorf("recebido %v esperado %v", recebido, esperado)
    }
}
```

(coloque `import reflect` no topo do seu arquivo para ter acesso ao `DeepEqual`)

É importante saber que o `reflect.DeepEqual` não tem "segurança de tipos", ou seja, o código vai compilar mesmo se você tiver feito algo estranho. Para ver isso em ação, altere o teste temporariamente para:

```go
func TestSomaTudo(t *testing.T) {

    recebido := SomaTudo([]int{1,2}, []int{0,9})
    esperado := "joao"

    if !reflect.DeepEqual(recebido, esperado) {
        t.Errorf("recebido %v, esperado %v", recebido, esperado)
    }
}
```

O que fizemos aqui foi comparar um `slice` com uma `string`. Isso não faz sentido, mas o teste compila! Logo, apesar de ser uma forma simples de comparar slices (e outras coisas), você deve tomar cuidado quando for usar o `reflect.DeepEqual`.

Volte o teste da forma como estava e execute-o. Você deve ter a saída do teste com uma mensagem tipo:

`soma_test.go:30: recebido [], esperado [3 9]`

## Escreva código o suficiente para fazer o teste passar

O que precisamos fazer é percorrer as variáveis recebidas como argumento, calcular a soma com nossa função `Soma` de antes e adicioná-la ao slice que vamos retornar:

```go
func SomaTudo(numerosParaSomar ...[]int) []int {
    quantidadeDeNumeros := len(numerosParaSomar)
    somas := make([]int, quantidadeDeNumeros)

    for i, numeros := range numerosParaSomar {
        somas[i] = Soma(numeros)
    }

    return somas
}
```

Muitas coisas novas para aprender!

Há uma nova forma de criar um slice. O `make` te permite criar um slice com uma capacidade inicial de `len` de `numerosParaSomar` que precisamos percorrer.

Você pode indexar slices como arrays com `meuSlice[N]` para obter seu valor ou designá-lo a um novo valor com `=`.

Agora o teste deve passar.

## Refatoração

Como mencionado, slices têm uma capacidade. Se você tiver um slice com uma capacidade de 2 e tentar fazer uma atribuição como `meuSlice[10] = 1`, vai receber um erro em _tempo de execução_.

No entanto, você pode usar a função `append`, que recebe um slice e um novo valor e retorna um novo slice com todos os itens dentro dele.

```go
func SomaTudo(numerosParaSomar ...[]int) []int {
    var somas []int
    for _, numeros := range numerosParaSomar {
        somas = append(somas, Soma(numeros))
    }

    return somas
}
```

Nessa implementação, nos preocupamos menos sobre capacidade. Começamos com um slice vazio `somas` e o anexamos ao resultado de `Soma` enquanto percorremos as variáveis recebidas como argumento.

Nosso próprio requisito é alterar o `SomaTudo` para `SomaTodosOsFinais`, onde agora calcula os totais de todos os "finais" de cada slice. O final de uma coleção é todos os itens com exceção do primeiro (a "cabeça").

## Escreva o teste primeiro

```go
func TestSomaTodosOsFinais(t *testing.T) {
    resultado := SomaTodosOsFinais([]int{1,2}, []int{0,9})
    esperado := []int{2, 9}

    if !reflect.DeepEqual(resultado, esperado) {
        t.Errorf("resultado %v, esperado %v", resultado, esperado)
    }
}
```

## Execute o teste

`./soma_test.go:26:9: undefined: SomaTodosOsFinais`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhado

Renomeie a função para `` e reexecute o teste.

Rename the function to `SomaTodosOsFinais` and re-run the test

`sum_test.go:30: resultado [3 9], esperado [2 9]`

## Escreva código o suficiente para fazer o teste passar

```go
func SomaTodosOsFinais(numerosParaSomar ...[]int) []int {
    var somas []int
    for _, numeros := range numerosParaSomar {
        final := numeros[1:]
        somas = append(somas, Soma(final))
    }

    return somas
}
```

Slices podem ser "fatiados"! A sintaxe usada é `slice[inicio:final]`. Se você omitir o valor de um dos lados dos `:` ele captura tudo do lado omitido. No nosso caso, quando usamos `numeros[1:]`, estamos dizendo "pegue da posição 1 até o final". É uma boa ideia investir um tempo escrevend outros testes com slices e brincar com o operador slice para criar mais familiaridade com ele.

## Refatoração

Não tem muito o que refatorar dessa vez.

O que acha que aconteceria se você passar um slice vazio para a nossa função? Qual é o "final" de um slice vazio? O que acontece quando você fala para o Go capturar todos os elementos de `meuSliceVazio[1:]`?

## Escreva o teste primeiro

```go
func TestSomaTodosOsFinais(t *testing.T) {

    t.Run("faz as somas de alguns slices", func(t *testing.T) {
        resultado := SomaTodosOsFinais([]int{1,2}, []int{0,9})
        esperado := []int{2, 9}

        if !reflect.DeepEqual(resultado, esperado) {
            t.Errorf("resultado %v, esperado %v", resultado, esperado)
        }
    })

    t.Run("soma slices vazios de forma segura", func(t *testing.T) {
        resultado := SomaTodosOsFinais([]int{}, []int{3, 4, 5})
        esperado := []int{0, 9}

        if !reflect.DeepEqual(resultado, esperado) {
            t.Errorf("resultado %v, esperado %v", resultado, esperado)
        }
    })

}
```

## Execute o teste

```bash
panic: runtime error: slice bounds out of range [recovered]
    panic: runtime error: slice bounds out of range
```

`pânico: erro em tempo de execução: fora da capacidade do slice`

Oh, não! É importante perceber que o test _foi compilado_, esse é um erro em tempo de execução. Erros em tempo de compilação são nossos amigos, porque nos ajudam a escrever softwares que funcionam. Erros em tempo de execução são nosso inimigos, porque afetam nossos usuários.

## Escreva código o suficiente para fazer o teste passar

```go
func SomaTodosOsFinais(numerosParaSomar ...[]int) []int {
    var somas []int
    for _, numeros := range numerosParaSomar {
        if len(numeros) == 0 {
            somas = append(somas, 0)
        } else {
            final := numeros[1:]
            somas = append(somas, Sum(final))
        }
    }

    return somas
}
```

## Refatoração

Our tests have some repeated code around assertion again, let's extract that into a function

```go
func TestSomaAllTails(t *testing.T) {

    checkSums := func(t *testing.T, got, want []int) {
        t.Helper()
        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v want %v", got, want)
        }
    }

    t.Run("make the sums of tails of", func(t *testing.T) {
        got := SumAllTails([]int{1, 2}, []int{0, 9})
        want := []int{2, 9}
        checkSums(t, got, want)
    })

    t.Run("safely sum empty slices", func(t *testing.T) {
        got := SumAllTails([]int{}, []int{3, 4, 5})
        want := []int{0, 9}
        checkSums(t, got, want)
    })

}
```

A handy side-effect of this is this adds a little type-safety to our code. If a silly developer adds a new test with `checkSums(t, got, "dave")` the compiler will stop them in their tracks.

```bash
$ go test
./sum_test.go:52:21: cannot use "dave" (type string) as type []int in argument to checkSums
```

## Wrapping up

We have covered

-   Arrays
-   Slices

    -   The various ways to make them
    -   How they have a _fixed_ capacity but you can create new slices from old ones

        using `append`

    -   How to slice, slices!

-   `len` to get the length of an array or slice
-   Test coverage tool
-   `reflect.DeepEqual` and why it's useful but can reduce the type-safety of your code

We've used slices and arrays with integers but they work with any other type too, including arrays/slices themselves. So you can declare a variable of `[][]string` if you need to.

[Check out the Go blog post on slices](https://blog.golang.org/go-slices-usage-and-internals) for an in-depth look into slices. Try writing more tests to demonstrate what you learn from reading it.

Another handy way to experiment with Go other than writing tests is the Go playground. You can try most things out and you can easily share your code if you need to ask questions. [I have made a go playground with a slice in it for you to experiment with.](https://play.golang.org/p/ICCWcRGIO68)
