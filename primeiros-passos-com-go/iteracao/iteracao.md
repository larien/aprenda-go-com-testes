# Iteração

**[Você pode encontrar todo o código desse capítulo aqui](https://github.com/larien/learn-go-with-tests/tree/master/primeiros-passos-com-go/iteracao)**

Para fazer coisas repetidamente em Go, você precisará do `for`. Go não possui nenhuma palavra chave do tipo `while`, `do` ou `until`. Você pode usar apenas `for`, o que é uma coisa boa!

Vamos escrever um teste para uma função que repete um caractere 5 vezes.

Não há nenhuma novidade até aqui, então tente escrever você mesmo para praticar.

## Escreva o teste primeiro

```go
package iteracao

import "testing"

func TestRepetir(t *testing.T) {
    repeticoes := Repetir("a")
    esperado := "aaaaa"

    if repeticoes != esperado {
        t.Errorf("esperado '%s' mas obteve '%s'", esperado, repeticoes)
    }
}
```

## Execute o teste

`./repetir_test.go:6:14: undefined: Repetir`

## Escreva a quantidade mínima de código para o teste rodar e verifique o erro na saída

_Mantenha a disciplina!_ Você não precisa saber nada de diferente agora para fazer o teste falhar apropriadamente.

Tudo o que foi feito até agora é o suficiente para compilar, para que você possa verificar se escreveu o teste corretamente.

```go
package iteracao

func Repetir(caractere string) string {
    return ""
}
```

Não é legal saber que você já conhece o bastante em Go para escrever testes para problemas simples? Isso significa que agora você pode mexer no código de produção o quanto quiser sabendo que ele se comportará da maneira que você desejar.

`repetir_test.go:10: esperado 'aaaaa' mas obteve ''`

## Escreva código o suficiente para fazer o teste passar

A sintaxe do `for` é muito fácil de lembrar e segue a maioria das linguagens baseadas em `C`:

```go
func Repetir(caractere string) string {
    var repeticoes string
    for i := 0; i < 5; i++ {
        repeticoes = repeticoes + caractere
    }
    return repeticoes
}
```

Ao contrário de outras linguagens como `C`, `Java` ou `Javascript`, não há parênteses ao redor dos três componentes do `for`. No entanto, as chaves `{ }` são obrigatórias.

Execute o teste e ele deverá passar.

Variações adicionais do loop `for` podem ser vistas [aqui](https://gobyexample.com/for).

## Refatoração

Agora é hora de refatorarmos e apresentarmos outro operador de atribuição: o `+=`.

```go
const quantidadeRepeticoes = 5

func Repetir(caractere string) string {
    var repeticoes string
    for i := 0; i < quantidadeRepeticoes; i++ {
        repeticoes += caractere
    }
    return repeticoes
}
```

O operador adicionar & atribuir `+=` adiciona o valor que está à direita no valor que esta à esquerda e atribui o resultado ao valor da esquerda. Também funciona com outros tipos, como por exemplo, inteiros (`integer`).

### Benchmarking

Escrever [benchmarks](https://golang.org/pkg/testing/#hdr-Benchmarks) em Go é outro recurso disponível nativamente na linguagem e é tão facil quanto escrever testes.

```go
func BenchmarkRepetir(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Repetir("a")
    }
}
```

Você notará que o código é muito parecido com um teste.

O `testing.B` dará a você acesso a `b.N`.

Quando o benchmark é rodado, ele executa `b.N` vezes e mede quanto tempo leva.

A quantidade de vezes que o código é executado não deve importar para você. O framework irá determinar qual valor é "bom" para que você consiga ter resultados decentes.

Para executar o benchmark, digite `go test -bench=.` no terminal (ou se estiver executando do PowerShell do Windows, `go test-bench="."`)

```bash
goos: darwin
goarch: amd64
pkg: github.com/larien/learn-go-with-tests/primeiros-passos-com-go/iteracao/v4
10000000           136 ns/op
PASS
```

`136 ns/op` significa que nossa função demora cerca de 136 nanossegundos para ser executada (no meu computador). E isso é ótimo! Para chegar a esse resultado ela foi executada 10000000 (10 milhões de vezes) vezes.

_NOTA_ por padrão, o benchmark é executado sequencialmente.

## Exercícios para praticar

-   Altere o teste para que a função possa especificar quantas vezes o caractere deve ser repetido e então corrija o código para passar no teste.
-   Escreva `ExampleRepetir` para documentar sua função.
-   Veja também o pacote [strings](https://golang.org/pkg/strings). Encontre funções que você considera serem úteis e experimente-as escrevendo testes como fizemos aqui. Investir tempo aprendendo a biblioteca padrão irá te recompensar com o tempo.

## Resumindo

-   Mais praticás de TDD
-   Aprendemos o `for`
-   Aprendemos como escrever benchmarks
