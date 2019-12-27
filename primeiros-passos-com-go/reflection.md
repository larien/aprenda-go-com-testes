# Reflection

[Do Twitter](https://twitter.com/peterbourgon/status/1011403901419937792?s=09)

> desafio golang: escreva uma função `percorre(x interface{}, fn func(string))` que recebe uma struct `x` e chama `fn` para todos os campos string encontrados dentro dela. nível de dificuldade: recursivamente.

Para fazer isso vamos precisar usar a `reflection`.

> A reflexão em computação é a habilidade de um programa examinar sua própria estrutura, particularmente através de tipos; é uma forma de metaprogramação. Também é uma ótima fonte de confusão.

De [The Go Blog: Reflection](https://blog.golang.org/laws-of-reflection)

## O que é `interface`?

Aproveitamos a segurança de tipos que o Go nos ofereceu em termos de funções que funcionam com tipos conhecidos, como `string`, `int` e nossos próprios tipos como `ContaBancaria`.

Isso significa que temos documentação de praxe e o compilador vai reclamar se você tentar passar o tipo errado para uma função.

Você pode se deparar com situações em que você quer escrever uma função, só que não sabe o tipo em tempo de compilação.

Go nos permite contornar isso com o tipo `interface{}` que você pode relacionar com _qualquer_ tipo.

Logo, `percorre(x interface{}, fn func(string))` aceitará qualquer valor para `x`.

### Então por que não usar `interface` para tudo e ter funções bem flexíveis?

* Como usuário de uma função que usa `interface`, você perde a segurança de tipos. E se você quisesse passar `Foo.bar` do tipo `string` para uma função, mas ao invés disso passa `Foo.baz` do tipo `int`? O compilador não vai ser capaz de te informar do seu erro. Você também não tem ideia _do que_ pode passar para uma função. Saber que uma função recebe um `ServicoDeUsuario`, por exemplo, é muito útil.

Resumindo, só use _reflection_ quando realmente precisar.

Se quiser funções polimórficas, considere desenvolvê-la em torno de uma interface (não `interface{}`, só para esclarecer) para que os usuários possam usar sua função com vários tipos se implementarem os métodos que você precisar para a sua função funcionar.

Nossa função vai precisar ser capaz de trabalhar com várias coisas diferentes. Como sempre, vamos usar uma abordagem iterativa, escrevendo testes para cada coisa nova que quisermos dar suporte e refatorando ao longo do caminho até finalizarmos.

# Escreva o teste primeiro

Vamos chamar nossa função com uma estrutura que tem um campo string dentro (`x`). Depois, podemos espiar a função (`fn`) passada para ela para ver se ela foi chamada.

```go
func TestPercorre(t *testing.T) {

    esperado := "Chris"
    var resultado []string

    x := struct {
        Nome string
    }{esperado}

    percorre(x, func(entrada string) {
        resultado = append(resultado, entrada)
    })

    if len(resultado) != 1 {
        t.Errorf("número incorreto de chamadas de função: resultado %d, esperado %d", len(resultado), 1)
    }
}
```

* Queremos armazenas um sice de strings (`resultado`) que armazena quais strings foram passadas dentro de `fn` pelo `percorre`. Algumas vezes, nos capítulos anteriores, criamos tipos dedicados para isso para espionar chamadas de função/método, mas nesse caso vamos apenas passá-lo em uma função anônima para `fn` que acaba em `resultado`.
* Usamos uma `struct` anônima com um campo `Nome` do tipo string para partir para caminho "feliz" e mais simples.
* Finalmente, chamamos `percorre` com `x` e o espião e por enquanto só verificamos o tamanho de `resultado`. Teremos mais precisão nas nossas verificações quando tivermos algo bem básico funcionando.
  
## Tente executar o teste

```text
./reflection_test.go:21:2: undefined: percorre
```

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Precisamos definir `percorre`.

```go
func percorre(x interface{}, fn func(entrada string)) {

}
```

Execute o teste novamente:

```text
=== RUN   TestPercorre
--- FAIL: TestPercorre (0.00s)
    reflection_test.go:19: número incorreto de chamadas de função: resultado 0, esperado 1
FAIL
```

### Escreva código o suficiente para fazer o teste passar

Agora podemos chamar o espião com qualquer string para fazer o teste passar.

```go
func percorre(x interface{}, fn func(entrada string)) {
    fn("Ainda não acredito que o Brasil perdeu de 7 a 1")
}
```

Agora o teste deve estar passando. A próxima coisa que vamos precisar fazer é criar uma verificação mais específica do que está sendo chamado dentro do nosso `fn`.

## Escreva o teste primeiro

Adicione o código a seguir para o teste existente para verificar se a string passada para `fn` está correta:

```go
if resultado[0] != esperado {
    t.Errorf("resultado '%s', esperado '%s'", resultado[0], esperado)
}
```

## Execute o teste

```text
=== RUN   TestPercorre
--- FAIL: TestPercorre (0.00s)
    reflection_test.go:23: resultado 'Ainda não acredito que o Brasil perdeu de 7 a 1', esperado 'Chris'
FAIL
```

### Escreva código o suficiente para fazer o teste passar

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := reflect.ValueOf(x) // ValorDe
	campo := valor.Field(0)     // Campo
	fn(campo.String())
}
```

Esse código está _pouco seguro e muito frágil_, mas lembre-se que nosso objetivo quando estamos no "vermelho" (os testes estão falhando) é escrever a menor quantidade de código possível. Depois escrevemos mais testes para resolver nossas lacunas.

Precisamos usar o reflection para verificar as propriedades de `x`.

No [pacote reflect](https://godoc.org/reflect) existe uma função chamada `ValueOf` que retorna um `Value` (valor) de determinada variável. Isso nos permite inspecionar um valor, inclusive seus campos usados nas próximas linhas.

Então podemos presumir coisas bem otimistas sobre o valor passado:

This code is _very unsafe and very naive_ but remembers our goal when we are in "red" \(the tests failing\) is to write the smallest amount of code possible. We then write more tests to address our concerns.

We need to use reflection to have a look at `x` and try and look at its properties.

The [reflect package]() has a function `ValueOf` which returns us a `Value` of a given variable. This has ways for us to inspect a value, including its fields which we use on the next line.

We then make some very optimistic assumptions about the value passed in

* We look at the first and only field, there may be no fields at all which would cause a panic
* We then call `String()` which returns the underlying value as a string but we know it would be wrong if the field was something other than a string.

## Refactor

Our code is passing for the simple case but we know our code has a lot of shortcomings.

We're going to be writing a number of tests where we pass in different values and checking the array of strings that `fn` was called with.

We should refactor our test into a table based test to make this easier to continue testing new scenarios.

```go
func TestPercorre(t *testing.T) {

    cases := []struct{
        Name string
        Input interface{}
        ExpectedCalls []string
    } {
        {
            "Struct with one string field",
            struct {
                Name string
            }{ "Chris"},
            []string{"Chris"},
        },
    }

    for _, test := range cases {
        t.Run(test.Name, func(t *testing.T) {
            var got []string
            walk(test.Input, func(input string) {
                got = append(got, input)
            })

            if !reflect.DeepEqual(got, test.ExpectedCalls) {
                t.Errorf("got %v, want %v", got, test.ExpectedCalls)
            }
        })
    }
}
```

Now we can easily add a scenario to see what happens if we have more than one string field.

## Write the test first

Add the following scenario to the `cases`.

```go
{
    "Struct with two string fields",
    struct {
        Name string
        City string
    }{"Chris", "London"},
    []string{"Chris", "London"},
}
```

## Try to run the test

```text
=== RUN   TestPercorre/Struct_with_two_string_fields
    --- FAIL: TestPercorre/Struct_with_two_string_fields (0.00s)
        reflection_test.go:40: got [Chris], want [Chris London]
```

## Write enough code to make it pass

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i:=0; i<val.NumField(); i++ {
        field := val.Field(i)
        fn(field.String())
    }
}
```

`value` has a method `NumField` which returns the number of fields in the value. This lets us iterate over the fields and call `fn` which passes our test.

## Refactor

It doesn't look like there's any obvious refactors here that would improve the code so let's press on.

The next shortcoming in `walk` is that it assumes every field is a `string`. Let's write a test for this scenario.

## Write the test first

Add the following case

```go
{
    "Struct with non string field",
    struct {
        Name string
        Age  int
    }{"Chris", 33},
    []string{"Chris"},
},
```

## Try to run the test

```text
=== RUN   TestPercorre/Struct_with_non_string_field
    --- FAIL: TestPercorre/Struct_with_non_string_field (0.00s)
        reflection_test.go:46: got [Chris <int Value>], want [Chris]
```

## Write enough code to make it pass

We need to check that the type of the field is a `string`.

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        if field.Kind() == reflect.String {
            fn(field.String())
        }
    }
}
```

We can do that by checking its [`Kind`](https://godoc.org/reflect#Kind).

## Refactor

Again it looks like the code is reasonable enough for now.

The next scenario is what if it isn't a "flat" `struct`? In other words, what happens if we have a `struct` with some nested fields?

## Write the test first

We have been using the anonymous struct syntax to declare types ad-hocly for our tests so we could continue to do that like so

```go
{
    "Nested fields",
    struct {
        Name string
        Profile struct {
            Age  int
            City string
        }
    }{"Chris", struct {
        Age  int
        City string
    }{33, "London"}},
    []string{"Chris", "London"},
},
```

But we can see that when you get inner anonymous structs the syntax gets a little messy. [There is a proposal to make it so the syntax would be nicer](https://github.com/golang/go/issues/12854).

Let's just refactor this by making a known type for this scenario and reference it in the test. There is a little indirection in that some of the code for our test is outside the test but readers should be able to infer the structure of the `struct` by looking at the initialisation.

Add the following type declarations somewhere in your test file

```go
type Person struct {
    Name    string
    Profile Profile
}

type Profile struct {
    Age  int
    City string
}
```

Now we can add this to our cases which reads a lot clearer than before

```go
{
    "Nested fields",
    Person{
        "Chris",
        Profile{33, "London"},
    },
    []string{"Chris", "London"},
},
```

## Try to run the test

```text
=== RUN   TestPercorre/Nested_fields
    --- FAIL: TestPercorre/Nested_fields (0.00s)
        reflection_test.go:54: got [Chris], want [Chris London]
```

The problem is we're only iterating on the fields on the first level of the type's hierarchy.

## Write enough code to make it pass

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        if field.Kind() == reflect.String {
            fn(field.String())
        }

        if field.Kind() == reflect.Struct {
            walk(field.Interface(), fn)
        }
    }
}
```

The solution is quite simple, we again inspect its `Kind` and if it happens to be a `struct` we just call `walk` again on that inner `struct`.

## Refactor

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

When you're doing a comparison on the same value more than once _generally_ refactoring into a `switch` will improve readability and make your code easier to extend.

What if the value of the struct passed in is a pointer?

## Write the test first

Add this case

```go
{
    "Pointers to things",
    &Person{
        "Chris",
        Profile{33, "London"},
    },
    []string{"Chris", "London"},
},
```

## Try to run the test

```text
=== RUN   TestPercorre/Pointers_to_things
panic: reflect: call of reflect.Value.NumField on ptr Value [recovered]
    panic: reflect: call of reflect.Value.NumField on ptr Value
```

## Write enough code to make it pass

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

You can't use `NumField` on a pointer `Value`, we need to extract the underlying value before we can do that by using `Elem()`.

## Refactor

Let's encapsulate the responsibility of extracting the `reflect.Value` from a given `interface{}` into a function.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}

func getValue(x interface{}) reflect.Value {
    val := reflect.ValueOf(x)

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    return val
}
```

This actually adds _more_ code but I feel the abstraction level is right.

* Get the `reflect.Value` of `x` so I can inspect it, I don't care how.
* Iterate over the fields, doing whatever needs to be done depending on its type.

Next, we need to cover slices.

## Write the test first

```go
{
    "Slices",
    []Profile {
        {33, "London"},
        {34, "Reykjavík"},
    },
    []string{"London", "Reykjavík"},
},
```

## Try to run the test

```text
=== RUN   TestPercorre/Slices
panic: reflect: call of reflect.Value.NumField on slice Value [recovered]
    panic: reflect: call of reflect.Value.NumField on slice Value
```

## Write the minimal amount of code for the test to run and check the failing test output

This is similar to the pointer scenario before, we are trying to call `NumField` on our `reflect.Value` but it doesn't have one as it's not a struct.

## Write enough code to make it pass

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    if val.Kind() == reflect.Slice {
        for i:=0; i< val.Len(); i++ {
            walk(val.Index(i).Interface(), fn)
        }
        return
    }

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

## Refactor

This works but it's yucky. No worries, we have working code backed by tests so we are free to tinker all we like.

If you think a little abstractly, we want to call `walk` on either

* Each field in a struct
* Each _thing_ in a slice

Our code at the moment does this but doesn't reflect it very well. We just have a check at the start to see if it's a slice \(with a `return` to stop the rest of the code executing\) and if it's not we just assume it's a struct.

Let's rework the code so instead we check the type _first_ and then do our work.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    switch val.Kind() {
    case reflect.Struct:
        for i:=0; i<val.NumField(); i++ {
            walk(val.Field(i).Interface(), fn)
        }
    case reflect.Slice:
        for i:=0; i<val.Len(); i++ {
            walk(val.Index(i).Interface(), fn)
        }
    case reflect.String:
        fn(val.String())
    }
}
```

Looking much better! If it's a struct or a slice we iterate over its values calling `walk` on each one. Otherwise, if it's a `reflect.String` we can call `fn`.

Still, to me it feels like it could be better. There's repetition of the operation of iterating over fields/values and then calling `walk` but conceptually they're the same.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice:
        numberOfValues = val.Len()
        getField = val.Index
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

If the `value` is a `reflect.String` then we just call `fn` like normal.

Otherwise, our `switch` will extract out two things depending on the type

* How many fields there are
* How to extract the `Value` \(`Field` or `Index`\)

Once we've determined those things we can iterate through `numberOfValues` calling `walk` with the result of the `getField` function.

Now we've done this, handling arrays should be trivial.

## Write the test first

Add to the cases

```go
{
    "Arrays",
    [2]Profile {
        {33, "London"},
        {34, "Reykjavík"},
    },
    []string{"London", "Reykjavík"},
},
```

## Try to run the test

```text
=== RUN   TestPercorre/Arrays
    --- FAIL: TestPercorre/Arrays (0.00s)
        reflection_test.go:78: got [], want [London Reykjavík]
```

## Write enough code to make it pass

Arrays can be handled the same way as slices, so just add it to the case with a comma

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

The final type we want to handle is `map`.

## Write the test first

```go
{
    "Maps",
    map[string]string{
        "Foo": "Bar",
        "Baz": "Boz",
    },
    []string{"Bar", "Boz"},
},
```

## Try to run the test

```text
=== RUN   TestPercorre/Maps
    --- FAIL: TestPercorre/Maps (0.00s)
        reflection_test.go:86: got [], want [Bar Boz]
```

## Write enough code to make it pass

Again if you think a little abstractly you can see that `map` is very similar to `struct`, it's just the keys are unknown at compile time.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walk(val.MapIndex(key).Interface(), fn)
        }
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

However, by design you cannot get values out of a map by index. It's only done by _key_, so that breaks our abstraction, darn.

## Refactor

How do you feel right now? It felt like maybe a nice abstraction at the time but now the code feels a little wonky.

_This is OK!_ Refactoring is a journey and sometimes we will make mistakes. A major point of TDD is it gives us the freedom to try these things out.

By taking small steps backed by tests this is in no way an irreversible situation. Let's just put it back to how it was before the refactor.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    walkValue := func(value reflect.Value) {
        walk(value.Interface(), fn)
    }

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        for i := 0; i< val.NumField(); i++ {
            walkValue(val.Field(i))
        }
    case reflect.Slice, reflect.Array:
        for i:= 0; i<val.Len(); i++ {
            walkValue(val.Index(i))
        }
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walkValue(val.MapIndex(key))
        }
    }
}
```

We've introduced `walkValue` which DRYs up the calls to `walk` inside our `switch` so that they only have to extract out the `reflect.Value`s from `val`.

### One final problem

Remember that maps in Go do not guarantee order. So your tests will sometimes fail because we assert that the calls to `fn` are done in a particular order.

To fix this, we'll need to move our assertion with the maps to a new test where we do not care about the order.

```go
t.Run("with maps", func(t *testing.T) {
    aMap := map[string]string{
        "Foo": "Bar",
        "Baz": "Boz",
    }

    var got []string
    walk(aMap, func(input string) {
        got = append(got, input)
    })

    assertContains(t, got, "Bar")
    assertContains(t, got, "Boz")
})
```

Here is how `assertContains` is defined

```go
func assertContains(t *testing.T, haystack []string, needle string)  {
    contains := false
    for _, x := range haystack {
        if x == needle {
            contains = true
        }
    }
    if !contains {
        t.Errorf("expected %+v to contain '%s' but it didnt", haystack, needle)
    }
}
```

## Wrapping up

* Introduced some of the concepts from the `reflect` package.
* Used recursion to traverse arbitrary data structures.
* Did an in retrospect bad refactor but didn't get too upset about it. By working iteratively with tests it's not such a big deal.
* This only covered a small aspect of reflection. [The Go blog has an excellent post covering more details](https://blog.golang.org/laws-of-reflection).
* Now that you know about reflection, do your best to avoid using it.

