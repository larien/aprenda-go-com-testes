# Olá, mundo

**[Você pode encontrar os códigos abordados nesse capítulo aqui](https://github.com/larien/learn-go-with-tests)**

É comum o primeiro programa em uma nova linguagem ser um _Olá, mundo_.

No [capítulo anterior](install-go.md#go-environment) discutimos sobre como Go pode ser dogmático como onde você coloca seus arquivos.

Crie um diretório no seguinte caminho `$GOPATH/src/github.com/{seu-lindo-nome-de-usuario}/ola`.

Se você estiver num ambiente baseado em unix e seu nome de usuário do Sistema Operacional for "bob" e você está motivado em seguir as convenções do Go sobre `$GOPATH` (que é a maneira mais fácil de configurar) você pode rodar `mkdir -p $GOPATH/src/github.com/bob/ola`.

Crie um arquivo chamado `ola.go` no diretório mencionado e escreva o seguinte código. Para rodá-lo, basta digitar no console `go run ola.go`.

```go
package main

import "fmt"

func main() {
    fmt.Println("Olá, mundo")
}
```

## Como isso funciona?

Quando você escreve um programa em Go, você irá ter um pacote `main` definido com uma função(`func`) `main` dentro dele. Os pacotes são maneiras de agrupar códigos escritos em Go.

A palavra reservada `func` é utilizada para que você defina uma função com um nome e um corpo.

Usando `import "fmt"` nós estamos importando um pacote que contém a função `Println` que será utilizada para imprimir (escrever) um valor na tela.

## Como testar isso?

Como você testaria isso? É bom separar seu "domínio"(suas regras de negócio) do resto do mundo \(efeitos colaterais\). A função `fmt.Println` é um efeito colateral \(que está imprimindo um valor no _**stdout** [saída padrão do terminal]_ \) e a string, nós estamos enviando dentro do seu próprio domínio.

Então, vamos separar essas referências para ficar mais fácil para testarmos.

```go
package main

import "fmt"

func Ola() string {
    return "Olá, mundo"
}

func main() {
    fmt.Println(Ola())
}
```

Nós criamos uma nova função usando `func`, mas dessa vez nós adicionamos outra palavra reservada `string` na sua definição. Isso significa que essa função irá ter como retorno uma `string` (_cadeia de caracteres_).

Agora, criaremos outro arquivo chamado `ola_test.go` onde nós iremos escrever um teste para nossa função `Ola`.

```go
package main

import "testing"

func TestOla(t *testing.T) {
    obtido := Ola()
    esperado := "Olá, mundo"

    if got != want {
        t.Errorf("obtido '%s' esperado '%s'", obtido, esperado)
    }
}
```

Antes de explicar, vamos rodar o código. Execute `go test` no seu terminal. Isso deve passar! Para checar, tente quebrar de alguma forma o teste mudando a string `esperado`.

Perceba que você nào precisa usar várias frameworks (ou bibliotecas) de testes e ficar se complicando tentando instalá-las. Tudo o que você precisa está pronto na linguagem e a sintaxe é a mesma para o resto dos códigos que você irá escrever.

### Escrevendo testes

Escrever um teste é como escrever uma função, com algumas regras

* Ele precisa estar num arquivo com um nome parecido com `xxx_test.go`
* A função de teste precisa começar com a palavra `Test`
* A função de teste recebe apenas um único argumento `t *testing.T`

Por agora é o bastante para saber que o nosso `t` do tipo `*testing.T` é o nosso "hook"(gancho) dentro do framework de testes e assim você poderá utilizar o `t.Fail()` quando você precisar relatar um erro.

### Abordando alguns novos tópicos:

#### `if`
Instruções `If` em Go são muito parecidas com a de outras linguagens.

#### Declarando Variáveis

Nós estamos declarando algumas variáveis com a sintaxe  `nomeDaVariavel := valor`, que nos permite reutilizar alguns valores nos nossos testes de maneira legível.

#### `t.Errorf`

Nós estamos chamando o _método_ `Errorf` em nosso `t` que irá imprimir uma mensagem e falhar o teste. O sufixo `f` significa que podemos formatar e montar uma string com valores inseridos dentro de valores de preenchimentos `%s`. Quando fazemos um teste falhar, devemos ser bastante claros como isso tudo aconteceu.

Nós iremos explorar mais na frente a diferença entre métodos e funções.

### Go doc

Outra funcionalidade importante de Go é sua documentação. Você pode rodar a documentação localmente rodando `godoc -http :8000`. Se você for para [localhost:8000/pkg](http://localhost:8000/pkg) irá ver todos os pacotes instalados no seu sistema.

A vasta biblioteca padrão da linguagem tem uma documentação excelente e com exemplos. Navegando para [http://localhost:8000/pkg/testing/](http://localhost:8000/pkg/testing/) vale a pena dar uma olhada para verificar o que está disponível para você.

### Olá, VOCÊ

Agora que temos um teste, nós podemos iterar sobre nosso software de maneira segura.

No último exemplo, nós escrevemos o teste somente _depois_ do código ser escrito, apenas para que você pudesse ter um exemplo de como escrever um teste e declarar uma função. A partir de agora, estamos _escrevendo os testes primeiro_.

Nosso próximo requisito é nos deixar especificar quem recebe a saudação.

Vamos começar especificando esses requisitos em um teste. Estamos fazendo um TDD(desenvolvimento orientado a testes) bastante simples e que nos permite ter certeza que nosso teste está _testando_ o que nós precisamos. Quando você escreve testes retroativamente existe o risco que seu teste pode continuar passando mesmo que o código não esteja funcionando como esperado.

```go
package main

import "testing"

func TestOla(t *testing.T) {
    obtido := Ola("Chris")
    esperado := "Olá, Chris"

    if obtido != esperado {
        t.Errorf("obtido '%s' esperado '%s'", obtido, esperado)
    }
}
```

Agora, rodando `go test`, deve ter aparecido um erro de compilação

```text
./ola_test.go:6:18: too many arguments in call to Ola
    have (string)
    want ()
```

Quando você está usando uma linguagem estaticamente tipada como Go, é importante _escutar o compilador_. O compilador entende como seu código deve se encaixar, não delegando essa função a você.

Neste caso, o compilador está te falando o que você precisa fazer para continuar. Nós temos que mudar a nossa função `Ola` para receber apenas um argumento.

Edite a função `Ola` para que seja aceito um argumento do tipo string

```go
func Ola(nome string) string {
    return "Olá, mundo"
}
```

Se você tentar rodar seus testes novamente, seu arquivo `main.go` irá falhar durante a compilação por que você não está passando um argumento. Passe "mundo" como argumento para fazer o teste passar.

```go
func main() {
    fmt.Println(Ola("mundo"))
}
```

Agora, quando você for rodar seus testes você verá algo parecido com isso

```text
ola_test.go:10: got 'Olá, mundo' want 'Olá, Chris''
```

Agora, finalmente temos um programa que compila mas não está satisfazendo os requisitos de acordo com o teste.

Vamos então fazer o teste passar usando o argumento `nome` e concatenar com `Olá,`

```go
func Ola(name string) string {
    return "Olá, " + name
}
```

Quando você rodar os testes eles irão passar. É comum como parte do ciclo do TDD _refatorar_ o nosso código agora.

### Uma nota sobre versionamento de código

Nesse ponto, se você estiver usando um versionamento de código \(que você deveria estar fazendo!\) eu faria um `commit` do código no estado atual. Agora, temos um software funcional suportado por um teste.

Apesar de que eu _não faria_ um push para a master, por que eu planejo refatorar em breve. É legal fazer um commit nesse ponto porque você pode se perder com a refatoração, fazendo um commit você pode sempre voltar para a última versão funcional do seu software.

Não tem muita coisa para refatorar aqui, mas nós podemos introduzir outro recurso da linguagem: _constantes_.

### Constantes

Constantes podem ser definidas como o exemplo abaixo:

```go
const prefixoOlaPortugues = "Olá, "
```

Agora, podemos refatorar nosso código

```go
const prefixoOlaPortugues = "Olá, "

func Ola(nome string) string {
    return prefixoOlaPortugues + nome
}
```

Depois da refatoração, rode novamente os seus testes para ter certeza que você não quebrou nada.

Constantes devem melhorar a performance da nossa aplicação assim como evitam com que você crie uma string `"Ola, "` para cada vez que `Ola` é chamado.

Sendo mais claro, o aumento de performance é incrivelmente insignificante para esse exemplo! Mas vale a pena pensar em criar constantes para capturar o significado dos valores e, às vezes, para ajudar no desempenho.

## Olá, mundo... novamente

O próximo requisito é: quando nossa função for chamada com uma string vazia, ela precisa imprimir o valor padrão "Olá, mundo", ao invés de "Olá, ".

Começaremos escrevendo um novo teste que irá falhar

```go
func TestOla(t *testing.T) {

    t.Run("diga olá para as pessoas", func(t *testing.T) {
        obtido := Ola("Chris")
        esperado := "Olá, Chris"

        if obtido != esperado {
            t.Errorf("obtido '%s' esperado '%s'", got, want)
        }
    })

    t.Run("diga 'Olá, mundo' quando uma string vazia for passada", func(t *testing.T) {
        obtido := Ola("")
        esperado := "Olá, mundo"

        if obtido != esperado {
            t.Errorf("obtido '%s' esperado '%s'", obtido, esperado)
        }
    })

}
```

Aqui nós estamos introduzindo outra ferramenta em nosso arsenal de testes, _subtestes_. Às vezes, é útil agrupar testes em torno de uma "coisa" e, em seguida, ter _subtestes_ descrevendo diferentes cenários.

O benefício dessa abordagem é que você poderá construir um código que pode ser compartilhado por outros testes.

Há um código repetido quando verificamos se a mensagem é o que esperamos.

A refatoração não é _apenas_ o código de produção!

É importante que seus testes _sejam especificações claras_ do que o código precisa fazer.

Podemos e devemos refatorar nossos testes.

```go
func TestOla(t *testing.T) {

    verificaMensagemCorreta := func(t *testing.T, got, want string) {
        t.Helper()
        if obtido != esperado {
            t.Errorf("obtido '%s' esperado '%s'", got, want)
        }
    }

    t.Run("dizendo ola para as pessoas", func(t *testing.T) {
        obtido := Ola("Chris")
        esperado := "Olá, Chris"
        verificaMensagemCorreta(t, obtido, esperado)
    })

    t.Run("'Mundo' como padrão para 'string' vazia", func(t *testing.T) {
        obtido := Ola("")
        want := "Olá, Mundo"
        verificaMensagemCorreta(t, obtido, esperado)
    })

}
```

### O que fizemos aqui?

Refatoramos nossa asserção em uma função. Isso reduz a duplicação e melhora a legibilidade de nossos testes. No Go, você pode declarar funções dentro de outras funções e atribuí-las a variáveis. Você pode chamá-las, assim como as funções normais. Precisamos passar em `t * testing.T` para que possamos dizer ao código de teste que falhará quando necessário.

`t.Helper ()` é necessário para dizermos ao conjunto de testes que este é método auxiliar. Ao fazer isso, quando o teste falhar, o número da linha relatada estará em nossa chamada de função, e não dentro do nosso auxiliar de teste. Isso ajudará outros desenvolvedores a rastrear os problemas com maior facilidade. Se você ainda não entender, comente, faça um teste falhar e observe a saída do teste.

Agora que temos um teste bem escrito falhando, vamos corrigir o código, usando um `if`.

```go
const prefixoOlaPortugues = "Olá, "

func Ola(nome string) string {
    if nome == "" {
        nome = "Mundo"
    }
    return prefixoOlaPortugues + nome
}
```

Se executarmos nossos testes, veremos que ele satisfaz o novo requisito e não quebramos acidentalmente a outra funcionalidade.

### De volta controle de versão

Agora, estamos felizes com o código. Eu adicionaria mais um commit ao anterior, então apenas verifique o quão adorável ficou o nosso código com os testes.

### Disciplina

Vamos repassar o ciclo novamente

* Escreva um teste
* Compile o código
* Rode o teste, e veja o teste falhar, depois verifique a mensagem de erro
* Escreva um código mínimo necessário para o teste passar
* Refatore

Este ciclo pode parecer tedioso, mas se manter nesse ciclo de feedback é importante.

Ele não apenas garante que você tenha _testes relevantes_, como também ajuda a _projetar um bom software_ refatorando-o com a segurança dos testes.

Ver a falha no teste é uma verificação importante porque também permite que você veja como é a mensagem de erro. Para quem programa, pode ser muito difícil trabalhar com uma base de código que, quando há falha nos testes, não fornece uma idéia clara de qual é o problema.

Assegurando que seus testes sejam rápidos e configurando suas ferramentas para que a execução de testes seja simples, você pode entrar em um estado de fluxo ao escrever seu código.

Ao não escrever testes, você está comprometendo-se a verificar manualmente seu código executando o software que interrompe seu estado de fluxo e não economiza tempo, especialmente a longo prazo.

## Continue! Mais requisitos

Meu Deus, temos mais requisitos. Agora precisamos suportar um segundo parâmetro, especificando o idioma da saudação. Se for passado um idioma que não reconhecemos, use como padrão o português.

Devemos ter certeza de que podemos usar o TDD para aprimorar essa funcionalidade facilmente!

Escreva um teste para um usuário, passando espanhol. Adicione-o ao conjunto de testes existente.

```go
    t.Run("em Espanhol", func(t *testing.T) {
        obtido := Ola("Elodie", "Espanhol")
        esperado := "Hola, Elodie"
        verificaMensagemCorreta(t, obtido, esperado)
    })
```

Lembre-se de não trapacear! _Primeiro os Testes_. Quando você tenta executar o teste, o compilador deve reclamar porque está sendo chamando `Ola` com dois argumentos ao invés de um.

```text
./ola_test.go:27:19: too many arguments in call to Ola
    have (string, string)
    want (string)
```

Acerte os problemas de compilação, adicionando um novo argumento do tipo `string` ao método `Ola`

```go
func Ola(nome string, idioma string) string {
    if nome == "" {
        nome = "Mundo"
    }
    return prefixoOlaPortugues + nome
}
```

Quando você tentar executar o teste novamente, ele se queixará de não ter sido passado argumentos suficientes para `Ola` nos seus outros testes em `ola.go`

```text
./ola.go:15:19: not enough arguments in call to Ola
    have (string)
    want (string, string)
```

Corrija-os passando `strings` vazia. Agora todos os seus testes devem compilar _e_ passar, além do nosso novo cenário

```text
ola_test.go:29: got 'Olá, Elodie' want 'Hola, Elodie'
```

Podemos usar `if` aqui para verificar se o idioma é igual a "espanhol" e, em caso afirmativo, alterar a mensagem

```go
func Ola(nome string, idioma string) string {
    if nome == "" {
        nome = "Mundo"
    }

    if idioma == "Espanhol" {
        return "Hola, " + nome
    }

    return prefixoOlaPortugues + nome
}
```

Os testes devem passar agora.

Agora é hora de _refatorar_. Você verá alguns problemas no código, seqüências de caracteres "mágicas", algumas das quais são repetidas. Tente refatorar você mesmo, a cada alteração, execute novamente os testes para garantir que sua refatoração não esteja quebrando nada.

```go
const espanhol = "Espanhol"
const prefixoOlaPortugues = "Olá, "
const prefixoOlaEspanhol = "Hola, "

func Ola(nome string, idioma string) string {
    if nome == "" {
        nome = "Mundo"
    }

    if idioma == espanhol {
        return prefixoOlaEspanhol + nome
    }

    return prefixoOlaPortugues + nome
}
```

### Francês

* Escreva um teste que verifique que quando passamos o idioma `"Francês"` obtemos `"Bonjour, "`
* Veja o teste falhar, verifique que a mensagem de erro é fácil de ler
* Faça a mínima alteração de código o suficiente para que o teste passe

Você pode ter escrito algo parecido com isto.

```go
func Ola(nome string, idioma string) string {
    if nome == "" {
        nome = "Mundo"
    }

    if idioma == espanhol {
        return prefixoOlaEspanhol + nome
    }

    if idioma == frances {
        return prefixoOlaFrances + nome
    }

    return prefixoOlaPortugues + nome
}
```

## `switch`

Quando você tem muitas instruções `if` verificando um valor específico, é comum usar uma instrução` switch`. Podemos usar o `switch` para refatorar o código facilitando a leitura e a sua extensão, caso desejarmos adicionar suporte a mais idiomas posteriormente.

```go
func Ola(nome string, idioma string) string {
    if nome == "" {
        nome = "Mundo"
    }

    prefixo := prefixoOlaPortugues

    switch idioma {
    case frances:
        prefixo = prefixoOlaFrances
    case espanhol:
        prefixo = prefixoOlaEspanhol
    }

    return prefixo + nome
}
```

Faça um teste para incluir agora uma saudação no idioma de sua escolha e você deve ver como é simples estender nossa _fantástica_ função.

### uma ... última ... refatoração?

Você poderia argumentar que talvez nossa função esteja ficando um pouco grande. A refatoração mais simples para isso seria extrair algumas funcionalidades para outra função.

```go
func Ola(nome string, idioma string) string {
    if nome == "" {
        name = "Mundo"
    }

    return prefixoDaSaudacao(idioma) + nome
}

func prefixoDaSaudacao(idioma string) (prefixo string) {
    switch idioma {
    case frances:
        prefixo = prefixoOlaFrances
    case espanhol:
        prefixo = prefixoOlaEspanhol
    default:
        prefixo = prefixoOlaPortugues
    }
    return
}
```

Alguns novos conceitos:

* Em nossa assinatura de função, criamos um valor de retorno nomeado`(prefixo string)`.
* Isso criará uma variável chamada `prefixo` na nossa função.
  * Lhe será atribuído o valor "zero". Isso dependendo do tipo, por exemplo, `int`s serão 0 para strings serão ` "" `.
    * Você pode retornar o que quer que esteja definido, apenas chamando `return` ao invés de `return prefixo`.
  * Isso será exibido no Go Doc para sua função, para que possa tornar mais clara a intenção do seu código.
* `default` será escolhido caso nenhuma das outras instruções `case` do `switch` corresponder.
* O nome da função começa com uma letra minúscula. As funções públicas em _Go_ começam com uma letra maiúscula e as privadas, com minúsculas. Não queremos que as partes internas do nosso algoritmo sejam expostas ao mundo, portanto tornamos essa função privada.

## Resumindo

Quem imaginaria que você poderia tirar tanto proveito de um `Olá, mundo`?

Até agora você deve ter alguma compreensão de:

### Algumas das sintaxes da linguagem _Go_ para:

* Escrever testes
* Declarar funções, com argumentos e tipos de retorno
* `if`, `const` e `switch`
* Declarar variáveis e constantes

### O processo TDD e _por que_ as etapas são importantes

* _Escreva um teste com falha e veja-o falhar_, para que saibamos que escrevemos um teste _relevante_ para nossos requisitos e vimos que ele produz uma _descrição da falha fácil de entender_
* Escrevendo a menor quantidade de código para fazer passar, para que saibamos que temos software funcionando
* _Em seguida_, refatorar, apoiando-se na segurança de nossos testes para garantir que tenhamos um código bem feito e fácil de trabalhar

No nosso caso, passamos de `Ola()` para `Ola("nome")`, para `Ola ("nome"," Francês ")` em etapas pequenas e fáceis de entender.

Naturalmente, isso é trivial comparado ao software do "mundo real", mas os princípios ainda permanecem. O TDD é uma habilidade que precisa de prática para se desenvolver, mas, ao ser capaz de dividir os problemas em componentes menores que você pode testar, você terá muito mais facilidade em escrever software.
