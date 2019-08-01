# Structs, métodos e interfaces

[**Você pode encontrar todos os códigos desse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/structs)

Supondo que precisamos de algum código geométrico para calcular o perímetro de um retângulo dado uma altura e largura. Podemos escrever uma função `Perimeter(width float64, height float64)`, onde `float64` é para números em ponto flutuante como `123.45`.

O ciclo de TDD deve ser mais familiar agora para você.

## Escreva o teste primeiro

```go
func TestPerimeter(t *testing.T) {
    got := Perimeter(10.0, 10.0)
    want := 40.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}
```

Viu a nova string de formatação? O `f` é para nosso `float64` e o `.2` significa imprima 2 casas decimais.

## Execute o teste

`./shapes_test.go:6:9: undefined: Perimeter`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhando

```go
func Perimeter(width float64, height float64) float64 {
    return 0
}
```

Resulta em `shapes_test.go:10: got 0 want 40`.

## Escreva código suficiente para fazer o teste passar

```go
func Perimeter(width float64, height float64) float64 {
    return 2 * (width + height)
}
```

Por enquanto, tudo fácil. Agora vamos criar uma função chamada `Area(width, height float64)` que retorna a área de um retângulo.

Tente fazer isso sozinho, segundo o ciclo de TDD.

Você deveria terminar com os testes como estes:

```go
func TestPerimeter(t *testing.T) {
    got := Perimeter(10.0, 10.0)
    want := 40.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}

func TestArea(t *testing.T) {
    got := Area(12.0, 6.0)
    want := 72.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}
```

E código como este:

```go
func Perimeter(width float64, height float64) float64 {
    return 2 * (width + height)
}

func Area(width float64, height float64) float64 {
    return width * height
}
```

## Refatoração

Nosso código faz o trabalho, mas não contém nada explícito sobre retângulos. Uma pessoa descuidada poderia tentar passar a largura e altura de um triângulo para esta função sem perceber que ela retornará uma resposta errada.

Podemos apenas dar para a função um nome mais específico como `RectangleArea`. Uma solução mais limpa é definir nosso próprio _tipo_ chamado `Rectangle` que encapsula este conceito para nós.

Podemos criar um tipo simples usando uma **struct** (estrutura). [Uma struct](https://golang.org/ref/spec#Struct_types) é apenas uma coleção nomeada de campos onde você pode armazenar dados.

Declare uma `struct` assim:

```go
type Rectangle struct {
    Width float64
    Height float64
}
```

Agora vamos refatorar os testes para usar `Rectangle` em vez de simples `float64`.

```go
func TestPerimeter(t *testing.T) {
    rectangle := Rectangle{10.0, 10.0}
    got := Perimeter(rectangle)
    want := 40.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}

func TestArea(t *testing.T) {
    rectangle := Rectangle{12.0, 6.0}
    got := Area(rectangle)
    want := 72.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}
```

Lembre de rodar seus testes antes de tentar corrigir. Você deve ter erro útil como:

```text
./shapes_test.go:7:18: not enough arguments in call to Perimeter
    have (Rectangle)
    want (float64, float64)
```

Você pode acessar os campos de uma `struct` com a sintaxe `myStruct.field`.

Mude as duas funções para corrigir o teste.

```go
func Perimeter(rectangle Rectangle) float64 {
    return 2 * (rectangle.Width + rectangle.Height)
}

func Area(rectangle Rectangle) float64 {
    return rectangle.Width * rectangle.Height
}
```

Eu espero que você concorde que passando um `Rectangle` para a função comunica com mais clareza nossa intenção, mas existem mais benefícios em usar `structs` que já vamos entender.

Nosso próximo requisito é escrever uma função `Area` para círculos.

## Escreva o teste primeiro

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        got := Area(rectangle)
        want := 72.0

        if got != want {
            t.Errorf("got %.2f want %.2f", got, want)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        got := Area(circle)
        want := 314.1592653589793

        if got != want {
            t.Errorf("got %.2f want %.2f", got, want)
        }
    })

}
```

## Execute o teste

`./shapes_test.go:28:13: undefined: Circle`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhando

Precisamos definir nosso tipo `Circle`.

```go
type Circle struct {
    Radius float64
}
```

Agora rode os testes novamente.

`./shapes_test.go:29:14: cannot use circle (type Circle) as type Rectangle in argument to Area`

Algumas linguagens de programação permitem você fazer algo como:

```go
func Area(circle Circle) float64 { ... }
func Area(rectangle Rectangle) float64 { ... }
```

Mas em Go você não pode

`./shapes.go:20:32: Area redeclared in this block`

Temos duas escolhas:

* Podemos ter funções com o mesmo nome declaradas em _pacotes_ diferentes. Então poderíamos criar nossa `Area(Circle)` em um novo _pacote_, mas isso parece um exagero aqui.
* Em vez disso, podemos definir [_métodos_](https://golang.org/ref/spec#Method_declarations) em nosso mais novo tipo definido.

### O que são métodos?

Até agora só escrevemos _funções_, mas temos usado alguns métodos. Quando chamamos `t.Errorf`, nós chamamos o método `Errorf` na instância de nosso `t` \(`testing.T`\).

Um método é uma função com um receptor. Uma declaração de método víncula um identificador, o nome do método, a um método e associa o método com o tipo base do receptor.*

Métodos são muito similares a funções, mas, são chamados invocando eles em uma instância de um tipo específico.
Enquanto você chamar funções onde quiser, como por exemplo `Area(rectangle)`, você só pode chamar métodos em "coisas".

Um exemplo ajudará. Então vamos mudar nossos testes primeiro para chamar métodos em vez de funções, e, em seguida, corrigir o código.

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        got := rectangle.Area()
        want := 72.0

        if got != want {
            t.Errorf("got %.2f want %.2f", got, want)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        got := circle.Area()
        want := 314.1592653589793

        if got != want {
            t.Errorf("got %f want %f", got, want)
        }
    })

}
```

Se rodarmos os testes agora, recebemos:

```text
./shapes_test.go:19:19: rectangle.Area undefined (type Rectangle has no field or method Area)
./shapes_test.go:29:16: circle.Area undefined (type Circle has no field or method Area)
```

> type Circle has no field or method Area

Gostaria de reforçar quão grandioso é o compilador aqui. É muito importante ter tempo para ler lentamente as mensagens de erro que você recebe, isso te ajudará a longo prazo.

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhando

Vamos adicionar alguns métodos para nossos tipos:

```go
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64  {
    return 0
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64  {
    return 0
}
```

A sintaxe para declaração de métodos é quase a mesma que usamos para funções e isto é porque eles são tão parecidos. A única diferença é a sintaxe para o método receptor `func (receiverName RecieverType) MethodName(args)`.

Quando seu método é chamado em uma variável deste tipo, você tem sua referência para o dado através da variável `receiverName`. Em muitas outras linguagens de programação isto é feito implicitamente e você acessa o receptor através de `this`.

É uma convenção em Go que a variável receptora seja a primeira letra do tipo.

```go
r Rectangle
```

Se você executar novamente os testes, agora eles devem compilar e dar alguma saída do teste falhando.

## Escreva código suficiente para fazer o teste passar

Agora vamos fazer nossos testes de retângulo passarem corrigindo nosso novo método.

```go
func (r Rectangle) Area() float64  {
    return r.Width * r.Height
}
```

Se você executar novamente os testes, aqueles de retângulo devem passar, mas, os de círculo ainda falham.

Para fazer a função `Area` de círculo passar, nós vamos emprestar a constante `Pi` do pacote `math` \(lembre-se de importá-lo\).

```go
func (c Circle) Area() float64  {
    return math.Pi * c.Radius * c.Radius
}
```

## Refatoração

Existe alguma duplicação em nossos testes.

Tudo o que queremos fazer é pegar uma coleção de _formas_, chamar o método `Area()` e então verificar o resultado.

Nós queremos ser capazes de escrever um tipo de função `checkArea` que permita passar tanto `Rectangle` quanto `Circle`, mas falhe ao compilar se tentarmos passar algo que não seja uma _forma_.

Com Go, podemos codificar esta intenção com **interfaces**.

[Interfaces](https://golang.org/ref/spec#Interface_types) são um conceito muito poderoso em linguagens de programação estaticamente tipadas, como Go, porque permitem que você crie funções que podem ser usadas com diferentes tipos e cria código altamente desacoplado, mantendo ainda a segurança de tipo.

Vamos apresentar isso refatorando nossos testes.

```go
func TestArea(t *testing.T) {

    checkArea := func(t *testing.T, shape Shape, want float64) {
        t.Helper()
        got := shape.Area()
        if got != want {
            t.Errorf("got %.2f want %.2f", got, want)
        }
    }

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        checkArea(t, rectangle, 72.0)
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        checkArea(t, circle, 314.1592653589793)
    })

}
```

Estamos criando uma função auxiliar como fizemos em outros exercícios, mas desta vez, estamos pedindo que um `Shape` seja passado. Se tentarmos chamá-la com algo que não seja uma _forma_, não vai compilar.

Como algo se torna uma _forma_? Nós apenas falamos para o Go o que é um `Shape` usando uma declaração de interface.

```go
type Shape interface {
    Area() float64
}
```

Estamos criando um novo `tipo` assim como fizemos com `Rectangle` e `Circle`, mas desta vez é uma `interface` em vez de uma `struct`.

Uma vez adicionado isso ao código, os testes passarão.

### Espera, o que?

Isso é bem diferente das interfaces na maioria das outras linguagens de programação. Normalmente você tem que escrever um código para dizer `Meu tipo Foo implementa a interface Bar`.

Mas em nosso caso:

* `Rectangle` tem um método chamado `Area` que retorna um `float64`, então isso satisfaz a interface `Shape`.
* `Circle` tem um método chamado `Area` que retorna um `float64`, então isso satisfaz a interface `Shape`.
* `string` não tem tal método, então isso não satisfaz a interface.
* etc.

Em Go **resolução de interface é implícita**. Se o tipo que você passar combinar com o que a interface está esperando, o código será compilado.

### Desacoplando

Veja como nossa função auxiliar não precisa se preocupar se a _forma_ é um `Rectangle` ou um `Circle` ou um `Triangle`. Ao declarar uma interface, a função auxiliar está _desacoplada_ de tipos concretos e tem apenas o método que precisa para fazer o trabalho.

Este tipo de abordagem - de usar interfaces para declarar **somente o que você precisa** - é muito importante em desenho de software e será coberto mais detalhadamente nas próximas seções.

## Refatoração adicional

Agora que você tem algum entendimento sobre `structs`, podemos apresentar "table driven tests" (testes guiados por tabela)*

[Table driven tests](https://github.com/golang/go/wiki/TableDrivenTests) são úteis quando você quer construir uma lista de casos de testes que podem ser testados da mesma forma.

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        shape Shape
        want  float64
    }{
        {Rectangle{12, 6}, 72.0},
        {Circle{10}, 314.1592653589793},
    }

    for _, tt := range areaTests {
        got := tt.shape.Area()
        if got != tt.want {
            t.Errorf("got %.2f want %.2f", got, tt.want)
        }
    }

}
```

A única sintaxe nova aqui é a criação de uma "struct anônima", `areaTests`. Estamos declarando um slice de structs usando `[]struct` com dois campos, o `shape` e o `want`. Então preenchemos o slice com os casos.

Então iteramos sobre eles assim como fazemos com qualquer outro slice, usando os campos da struct para executar nossos testes.

Você pode ver como será muito fácil para uma pessoa introduzir uma nova forma, implementar `Area` e então adicioná-la nos casos de teste. Além disso, se for encontrada uma falha com `Area`, é muito fácil adicionar um novo caso de teste para exercitar antes de corrigí-la.

_Testes baseados em tabela_ podem ser um item valioso em sua caixa de ferramentas, mas, tenha certeza de que você precisa do ruído extra nos testes. Se você deseja testar várias implementações de uma interface ou se o dado passado para uma função tem muitos requisitos diferentes que precisam de testes, então eles se encaixam bem.

Vamos demonstrar tudo isso adicionando e testando outra forma; um triângulo.

## Escreva o teste primeiro

Adicionar um teste para nossa nova forma é muito fácil. Simplesmente adicione `{Triangle{12, 6}, 36.0},` à nossa lista.

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        shape Shape
        want  float64
    }{
        {Rectangle{12, 6}, 72.0},
        {Circle{10}, 314.1592653589793},
        {Triangle{12, 6}, 36.0},
    }

    for _, tt := range areaTests {
        got := tt.shape.Area()
        if got != tt.want {
            t.Errorf("got %.2f want %.2f", got, tt.want)
        }
    }

}
```

## Execute o teste

Lembre-se, continue tentando executar o teste e deixe o compilador guiá-lo em direção a solução.

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhando

`./shapes_test.go:25:4: undefined: Triangle`

Ainda não definimos `Triangle`:

```go
type Triangle struct {
    Base   float64
    Height float64
}
```

Tente novamente:

```text
./shapes_test.go:25:8: cannot use Triangle literal (type Triangle) as type Shape in field value:
    Triangle does not implement Shape (missing Area method)
```
`Triangle não implementa Shape (método Area ausente)`

Isso nos diz que não podemos usar um `Triangle` como uma forma porque ele não tem um método `Area()`, então adicione uma implementação vazia para termos o teste funcionando:

```go
func (t Triangle) Area() float64 {
    return 0
}
```

Finalmente o código compilou e temos o nosso erro:

`shapes_test.go:31: got 0.00 want 36.00`

## Escreva código suficiente para fazer o teste passar

```go
func (t Triangle) Area() float64 {
    return (t.Base * t.Height) * 0.5
}
```

E nossos testes passaram!

## Refatoração

Novamente, a implementação está boa, mas, nossos testes podem ser melhorados.

Quando você lê isso:

```go
{Rectangle{12, 6}, 72.0},
{Circle{10}, 314.1592653589793},
{Triangle{12, 6}, 36.0},
```

Não está imediatamente claro o que todos os números representam e você deve mirar para que seus testes sejam fáceis de entender.

Até agora você viu uma sintaxe para criar instâncias de structs `MyStruct{val1, val2}`, mas você pode opcionalmente nomear os campos.

Vamos ver como isso parece:

```go
        {shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
        {shape: Circle{Radius: 10}, want: 314.1592653589793},
        {shape: Triangle{Base: 12, Height: 6}, want: 36.0},
```

Em [Test-Driven Development by Example](https://g.co/kgs/yCzDLF) Kent Beck refatora alguns testes para um ponto e afirma:

> O teste nos fala mais claramente, como se fosse uma afirmação da verdade, **não uma sequência de operações**

\(ênfase minha\)

Agora nossos testes \(pelo menos a lista de casos\) fazem afirmações da verdade sobre formas e suas áreas.

## Garanta que a saída do seu teste seja útil

Lembra anteriormente quando implementamos `Triangle` e tivemos um teste falhando? Ele imprimiu `shapes_test.go:31: got 0.00 want 36.00`.

Nós sabíamos que estava relacionado ao `Triangle` porque estávamos trabalhando nisso, mas e se uma falha escorregasse para o sistema em um dos 20 casos na tabela? Como uma pessoa saberia qual caso falhou? Esta não é uma boa experiência. Eles teriam que olhar manualmente através dos casos para encontrar qual deles está falhando de fato.

Podemos mudar nossa mensagem de erro para `%#v got %.2f want %.2f`. A string de formatação `%#v` irá imprimir nossa struct com os valores em seu campo, então as pessoas podem ver imediatamente as propriedades que estão sendo testadas.

Para melhorar a legibilidade de nossos futuros casos de teste, podemos renomear o campo `want` para algo mais descritivo como `hasArea`.

Uma dica final com testes guiados por tabela é usar `t.Run` e renomear os casos de teste.

Envolvendo cada caso em um `t.Run` você terá uma saída de testes mais limpa em caso de falhas, além de imprimir o nome do caso.

```text
--- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/Rectangle (0.00s)
        shapes_test.go:33: main.Rectangle{Width:12, Height:6} got 72.00 want 72.10
```

E você pode rodar testes específicos dentro de sua tabela com `go test -run TestArea/Rectangle`.

Aqui está o código final do nosso teste que captura isso:

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        name    string
        shape   Shape
        hasArea float64
    }{
        {name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
        {name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
        {name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
    }

    for _, tt := range areaTests {
        // using tt.name from the case to use it as the `t.Run` test name
        t.Run(tt.name, func(t *testing.T) {
            got := tt.shape.Area()
            if got != tt.hasArea {
                t.Errorf("%#v got %.2f want %.2f", tt.shape, got, tt.hasArea)
            }
        })

    }

}
```

## Resumindo

Esta foi mais uma prática de TDD, iterando em nossas soluções para problemas matemáticos básicos e aprendendo novos recursos da linguagem motivados por nossos testes.

* Declarando structs para criar seus próprios tipos de dados, o que permite agrupar dados relacionados e torna a intenção do seu código mais clara.
* Declarando interfaces para que você possa definir funções que podem ser usadas por diferentes tipos \([polimorfismo paramétrico](https://pt.wikipedia.org/wiki/Polimorfismo_paramétrico)\).
* Adicionando métodos para que você possa adicionar funcionalidades aos seus tipos de dados e implementar interfaces.
* Testes baseados em tabela para tornar suas asserções mais claras e suas suítes mais fáceis de estender e manter.

Este foi um capítulo importante porque agora começamos a definir nossos próprios tipos. Em linguagens estaticamente tipadas como Go, conseguir projetar seus próprios tipos é essencial para construir software que seja fácil de entender, compilar e testar.

Interfaces são uma ótima ferramenta para ocultar a complexidade de outras partes do sistema. Em nosso caso, o _código_ de teste auxiliar não precisou conhecer a forma exata que estava afirmando, apenas como "pedir" pela sua área.

Conforme você se familiariza com Go, começa a ver a força real das interfaces e da biblioteca padrão.
Você aprenderá sobre as interfaces definidas na biblioteca padrão que são usadas _em todo lugar_ e, implementando-as em relação aos seus próprios tipos, você pode reutilizar rapidamente muitas das ótimas funcionalidades.
