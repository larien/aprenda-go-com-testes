# Inteiros

**[Você pode encontrar todos os códigos desse capítulo aqui](https://github.com/larien/learn-go-with-tests/tree/master/inteiros)**

Inteiros funcionam como é de se esperar. Vamos escrever uma função de soma para testar algumas coisas. Crie um arquivo de teste chamado `adicionador_test.go` e escreva o seguinte código.

**nota**: Os arquivos-fonte de Go devem ter apenas um `package`(pacote) por diretório, verifique se os arquivos estão organizados separadamente. [Aqui tem uma boa explicação sobre isso (em inglês).](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project)

## Escreva o teste primeiro

```go
package inteiros

import "testing"

func TestAdicionador(t *testing.T) {
    soma := Adiciona(2, 2)
    esperado := 4

    if soma != esperado {
        t.Errorf("esperado '%d', resultado '%d'", esperado, soma)
    }
}
```

Você deve ter notado que estamos usando `%d` como string de formatação, em vez de `%s`. Isso porque queremos que ele imprima um valor inteiro e não uma string.
Observe também que não estamos mais usando o pacote `main`, em vez disso, definimos um pacote chamado `inteiros`, pois o nome sugere que ele agrupará funções para trabalhar com números inteiros, como Adiciona.

## Tente executar o teste

Execute o test com `go test`

Inspecione o erro de compilação

`./adicionador_test.go:6:9: undefined: Adiciona`

## Escreva a quantidade mínima de código para o teste rodar e verifique o erro na saída do teste

Escreva apenas o suficiente de código para satisfazer o compilador _e nada mais_ - lembre-se de que queremos verificar se nossos testes falham pelo motivo certo.

```go
package inteirosinteiros

func Adiciona(x, y int) int {
    return 0
}
```

Quando você tem mais de um argumento do mesmo tipo (no nosso caso dois inteiros) ao invés de ter `(x int, y int)` você pode encurtá-lo para `(x, y int)`.

Agora execute os testes. Devemos ficar felizes que o teste esteja relatando corretamente o que está errado.

`adicionador_test.go:10: esperado '4', resultado '0'`

Você deve ter percebido que aprendemos sobre o _valor de retorno nomeado_ na [última](hello-world.md#one...last...refactor?) seção, mas não estamos usando aqui. Ele geralmente deve ser usado quando o significado do resultado não está claro no contexto. No nosso caso, é muito claro que a função `Adiciona` irá adicionar os parâmetros. Você pode consultar [esta](https://github.com/golang/go/wiki/CodeReviewComments#named-result-parameters) wiki para mais detalhes.

## Escreva código o suficiente para fazer o teste passar

No sentido estrito de TDD, devemos escrever a _quantidade mínima de código para fazer o teste passar_. Uma pessoa pretenciosa pode fazer isso:

```go
func Adiciona(x, y int) int {
    return 4
}
```

Ah hah! Frustração mais uma vez! TDD é uma farsa, né?

Poderíamos escrever outro teste, com números diferentes para forçar o teste a falhar, mas isso parece um jogo de gato e rato.

Quando estivermos mais familiarizados com a sintaxe do Go, apresentarei uma técnica chamada Testes Baseados em Propriedade, que interromperá a irritação das pessoas e ajudará a encontrar bugs.

Por enquanto, vamos corrigi-lo corretamente:

```go
func Adiciona(x, y int) int {
    return x + y
}
```

Se você executar os testes novamente, eles devem passar.

## Refatoração

Não há muitas melhorias que possamos fazer aqui.

Anteriormente, vimos como nomear o argumento de retorno que aparece na documentação e também na maioria dos editores de código.

Isso é ótimo porque ajuda na usabilidade do código que você está escrevendo. É preferível que um usuário possa entender o uso de seu código apenas observando a assinatura de tipo e a documentação.

Você pode adicionar documentação em funções escrevendo comentários, e elas aparecerão no Go Doc como quando você olha a documentação da biblioteca padrão.

```go
// Adiciona recebe dois inteiros e retorna a soma deles
func Adiciona(x, y int) int {
    return x + y
}
```

### Exemplos

Se realmente quer ir além, você pode fazer [exemplos](https://blog.golang.org/examples). Você encontrará muitos exemplos na documentação da biblioteca padrão.

Muitas vezes, exemplos encontrados fora da base de código, como um arquivo readme, ficam desatualizados e incorretos em comparação com o código real, porque eles não são verificados.

Os exemplos de Go são executados da mesma forma que os testes, para que você possa ter certeza de que eles refletem o que o código realmente faz.

Exemplos são compilados \(e opcionalmente executados\) como parte do conjunto de testes de um pacote.

Como nos testes comuns, os exemplos são funções que residem nos arquivos \_test.go de um pacote. Adicione a seguinte função ExemploAdiciona no arquivo `adicionador_test.go`.

```go
func ExemploAdiciona() {
    soma := Adiciona(1, 5)
    fmt.Println(soma)
    // Saida: 6
}
```
> obs: As palavras `Exemplo` e `Saida` devem ser alteradas para `Example` e `Output` para a execução correta do código.

(Se o seu editor não importar os pacotes automaticamente, a etapa de compilação irá falhar porque você não colocou o `import "fmt"` no `adicionador_test.go`. É altamente recomendável que você pesquise como ter esses tipos de erros corrigidos automaticamente em qualquer editor que esteja usando.)

Se o seu código mudar fazendo com que o exemplo não seja mais válido, você vai ter um erro de compilação.

Executando os testes do pacote, podemos ver que a função de exemplo é executada sem a necessidade de ajustes:

```bash
$ go test -v
=== RUN   TestAdicionador
--- PASS: TestAdicionador (0.00s)
=== RUN   ExemploAdiciona
--- PASS: ExemploAdiciona (0.00s)
```

Note que a função de exemplo não será executada se você remover o comentário "// Output: 6". Embora a função seja compilada, ela não será executada.

Ao adicionar este trecho de código, o exemplo aparecerá na documentação dentro do `godoc`, tornando seu código ainda mais acessível.

Para ver como isso funciona, execute `godoc -http=:6060` e navegue para `http://localhost:6060/pkg/`

Aqui você vai ver uma lista de todos os pacotes em seu `$GOPATH`. Então, supondo que tenha escrito esse código em algum lugar como `$GOPATH/src/github.com/{seu_id}`, você poderá encontrar uma documentação com seus exemplos.

Se você publicar seu código com exemplos em uma URL pública, poderá compartilhar a documentação em [godoc.org](https://godoc.org). Por exemplo, aqui está a API finalizada deste capítulo [https://godoc.org/github.com/larien/learn-go-with-tests/inteiros/v2](https://godoc.org/github.com/larien/learn-go-with-tests/inteiros/v2).

## Resumindo

Falamos sobre:

-   Mais práticas do fluxo de trabalho de TDD
-   Inteiros, adição
-   Escrever melhores documentações para que os usuários do nosso código possam entender seu uso rapidamente
-   Exemplos de como usar nosso código, que são verificados como parte de nossos testes
