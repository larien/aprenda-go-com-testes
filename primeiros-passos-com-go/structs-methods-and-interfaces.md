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

## Refactor

There is some duplication in our tests.

All we want to do is take a collection of _shapes_, call the `Area()` method on them and then check the result.

We want to be able to write some kind of `checkArea` function that we can pass both `Rectangle`s and `Circle`s to, but fail to compile if we try to pass in something that isn't a shape.

With Go, we can codify this intent with **interfaces**.

[Interfaces](https://golang.org/ref/spec#Interface_types) are a very powerful concept in statically typed languages like Go because they allow you to make functions that can be used with different types and create highly-decoupled code whilst still maintaining type-safety.

Let's introduce this by refactoring our tests.

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

We are creating a helper function like we have in other exercises but this time we are asking for a `Shape` to be passed in. If we try to call this with something that isn't a shape, then it will not compile.

How does something become a shape? We just tell Go what a `Shape` is using an interface declaration

```go
type Shape interface {
    Area() float64
}
```

We're creating a new `type` just like we did with `Rectangle` and `Circle` but this time it is an `interface` rather than a `struct`.

Once you add this to the code, the tests will pass.

### Wait, what?

This is quite different to interfaces in most other programming languages. Normally you have to write code to say `My type Foo implements interface Bar`.

But in our case

* `Rectangle` has a method called `Area` that returns a `float64` so it satisfies the `Shape` interface
* `Circle` has a method called `Area` that returns a `float64` so it satisfies the `Shape` interface
* `string` does not have such a method, so it doesn't satisfy the interface
* etc.

In Go **interface resolution is implicit**. If the type you pass in matches what the interface is asking for, it will compile.

### Decoupling

Notice how our helper does not need to concern itself with whether the shape is a `Rectangle` or a `Circle` or a `Triangle`. By declaring an interface the helper is _decoupled_ from the concrete types and just has the method it needs to do its job.

This kind of approach of using interfaces to declare **only what you need** is very important in software design and will be covered in more detail in later sections.

## Further refactoring

Now that you have some understanding of structs we can now introduce "table driven tests".

[Table driven tests](https://github.com/golang/go/wiki/TableDrivenTests) are useful when you want to build a list of test cases that can be tested in the same manner.

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

The only new syntax here is creating an "anonymous struct", areaTests. We are declaring a slice of structs by using `[]struct` with two fields, the `shape` and the `want`. Then we fill the slice with cases.

We then iterate over them just like we do any other slice, using the struct fields to run our tests.

You can see how it would be very easy for a developer to introduce a new shape, implement `Area` and then add it to the test cases. In addition, if a bug is found with `Area` it is very easy to add a new test case to exercise it before fixing it.

Table based tests can be a great item in your toolbox but be sure that you have a need for the extra noise in the tests. If you wish to test various implementations of an interface, or if the data being passed in to a function has lots of different requirements that need testing then they are a great fit.

Let's demonstrate all this by adding another shape and testing it; a triangle.

## Write the test first

Adding a new test for our new shape is very easy. Just add `{Triangle{12, 6}, 36.0},` to our list.

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

## Try to run the test

Remember, keep trying to run the test and let the compiler guide you toward a solution.

## Write the minimal amount of code for the test to run and check the failing test output

`./shapes_test.go:25:4: undefined: Triangle`

We have not defined Triangle yet

```go
type Triangle struct {
    Base   float64
    Height float64
}
```

Try again

```text
./shapes_test.go:25:8: cannot use Triangle literal (type Triangle) as type Shape in field value:
    Triangle does not implement Shape (missing Area method)
```

It's telling us we cannot use a Triangle as a shape because it does not have an `Area()` method, so add an empty implementation to get the test working

```go
func (t Triangle) Area() float64 {
    return 0
}
```

Finally the code compiles and we get our error

`shapes_test.go:31: got 0.00 want 36.00`

## Write enough code to make it pass

```go
func (t Triangle) Area() float64 {
    return (t.Base * t.Height) * 0.5
}
```

And our tests pass!

## Refactor

Again, the implementation is fine but our tests could do with some improvement.

When you scan this

```go
{Rectangle{12, 6}, 72.0},
{Circle{10}, 314.1592653589793},
{Triangle{12, 6}, 36.0},
```

It's not immediately clear what all the numbers represent and you should be aiming for your tests to easily understood.

So far you've only been shown one syntax for creating instances of structs `MyStruct{val1, val2}` but you can optionally name the fields.

Let's see what it looks like

```go
        {shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
        {shape: Circle{Radius: 10}, want: 314.1592653589793},
        {shape: Triangle{Base: 12, Height: 6}, want: 36.0},
```

In [Test-Driven Development by Example](https://g.co/kgs/yCzDLF) Kent Beck refactors some tests to a point and asserts:

> The test speaks to us more clearly, as if it were an assertion of truth, **not a sequence of operations**

\(emphasis mine\)

Now our tests \(at least the list of cases\) make assertions of truth about shapes and their areas.

## Make sure your test output is helpful

Remember earlier when we were implementing `Triangle` and we had the failing test? It printed `shapes_test.go:31: got 0.00 want 36.00`.

We knew this was in relation to `Triangle` because we were just working with it, but what if a bug slipped in to the system in one of 20 cases in the table? How would a developer know which case failed? This is not a great experience for the developer, they will have to manually look through the cases to find out which case actually failed.

We can change our error message into `%#v got %.2f want %.2f`. The `%#v` format string will print out our struct with the values in its field, so the developer can see at a glance the properties that are being tested.

To increase the readability of our test cases further we can rename the `want` field into something more descriptive like `hasArea`.

One final tip with table driven tests is to use `t.Run` and to name the test cases.

By wrapping each case in a `t.Run` you will have clearer test output on failures as it will print the name of the case

```text
--- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/Rectangle (0.00s)
        shapes_test.go:33: main.Rectangle{Width:12, Height:6} got 72.00 want 72.10
```

And you can run specific tests within your table with `go test -run TestArea/Rectangle`.

Here is our final test code which captures this

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

## Wrapping up

This was more TDD practice, iterating over our solutions to basic mathematic problems and learning new language features motivated by our tests.

* Declaring structs to create your own data types which lets you bundle related data together and make the intent of your code clearer
* Declaring interfaces so you can define functions that can be used by different types \([parametric polymorphism](https://en.wikipedia.org/wiki/Parametric_polymorphism)\)
* Adding methods so you can add functionality to your data types and so you can implement interfaces
* Table based tests to make your assertions clearer and your suites easier to extend & maintain

This was an important chapter because we are now starting to define our own types. In statically typed languages like Go, being able to design your own types is essential for building software that is easy to understand, to piece together and to test.

Interfaces are a great tool for hiding complexity away from other parts of the system. In our case our test helper _code_ did not need to know the exact shape it was asserting on, only how to "ask" for it's area.

As you become more familiar with Go you start to see the real strength of interfaces and the standard library. You'll learn about interfaces defined in the standard library that are used _everywhere_ and by implementing them against your own types you can very quickly re-use a lot of great functionality.

