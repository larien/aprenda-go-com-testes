# Reflection

[**Você pode encontrar todo o código para esse capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/main/primeiros-passos-com-go/reflection)

[Do Twitter](https://twitter.com/peterbourgon/status/1011403901419937792?s=09)

> Desafio Golang: escreva uma função `percorre(x interface{}, fn func(string))` que recebe uma struct `x` e chama `fn` para todos os campos string encontrados dentro dela. nível de dificuldade: recursão.

Para fazer isso vamos precisar usar `reflection` (reflexão).

> A reflexão em computação é a habilidade de um programa examinar sua própria estrutura, particularmente através de tipos; é uma forma de metaprogramação. Também é uma ótima fonte de confusão.

De [The Go Blog: Reflection](https://blog.golang.org/laws-of-reflection)

## O que é `interface`?

Aproveitamos a segurança de tipos que o Go nos ofereceu em termos de funções que funcionam com tipos conhecidos, como `string`, `int` e nossos próprios tipos como `ContaBancaria`.

Isso significa que de praxe temos documentação e o compilador vai reclamar se você tentar passar o tipo errado para uma função.

Só que você pode se deparar com situações em que quer escrever uma função, mas não sabe o tipo da variável em tempo de compilação.

Go nos permite contornar isso com o tipo `interface{}`, que você pode relacionar com _qualquer_ tipo.

Logo, `percorre(x interface{}, fn func(string))` aceitará qualquer valor para `x`.

### Então por que não usar `interface` para tudo e ter funções bem flexíveis?

* Quando utiliza uma função que usa `interface`, você perde a segurança de tipos. E se você quisesse passar `Foo.bar` do tipo `string` para uma função, mas ao invés disso passa `Foo.baz` do tipo `int`? O compilador não vai ser capaz de informar seu erro. Você também não tem ideia _do que_ pode passar para uma função. Saber que uma função recebe um `ServicoDeUsuario`, por exemplo, é muito útil.

Resumindo, só use reflexão quando realmente precisar.

Se quiser funções polimórficas, considere desenvolvê-la em torno de uma interface (não `interface{}`, só para esclarecer) para que os usuários possam usar sua função com vários tipos se implementarem os métodos que você precisar para a sua função funcionar.

Nossa função vai precisar ser capaz de trabalhar com várias coisas diferentes. Como sempre, vamos usar uma abordagem iterativa, escrevendo testes para cada coisa nova que quisermos dar suporte e refatorando ao longo do caminho até finalizarmos.

## Escreva o teste primeiro

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

* Queremos armazenar um slice de strings (`resultado`) que armazena quais strings foram passadas dentro de `fn` pelo `percorre`. Algumas vezes, nos capítulos anteriores, criamos tipos dedicados para isso para espionar chamadas de função/método, mas nesse caso vamos apenas passá-lo em uma função anônima para `fn` que acaba em `resultado`.
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

* Podemos procurar pelo primeiro e único campo, mas pode não haver nenhum campo, o que causaria um pânico.
* Depois podemos chamar `String()` que tetorna o valor subjacente como string, mas sabemos que vai dar errado se o campo for de algum tipo que não uma string.

## Refatoração

Nosso código está passando pelo caso simples, mas sabemos que nosso código tem várias falhas.

Vamos escrever alguns testes onde passamos valores diferentes e verificaremos o array de strings com que `fn` foi chamado.

Precisamos refatorar nosso teste em um teste orientado por tabelas para tornar esse processo mais fácil para continuarmos testando novas situações.

```go
func TestPercorre(t *testing.T) {

	casos := []struct {
		Nome              string
		Entrada           interface{}
		ChamadasEsperadas []string
	}{
		{
			"Struct com um campo string",
			struct {
				Nome string
			}{"Chris"},
			[]string{"Chris"},
		},
	}

	for _, teste := range casos {
		t.Run(teste.Nome, func(t *testing.T) {
			var resultado []string
			percorre(teste.Entrada, func(entrada string) {
				resultado = append(resultado, entrada)
			})

			if !reflect.DeepEqual(resultado, teste.ChamadasEsperadas) {
				t.Errorf("resultado %v, esperado %v", resultado, teste.ChamadasEsperadas)
			}
		})
	}
}
```

Agora podemos adicionar uma situação facilmente para ver o que acontece se tivermos mais de um campo string.

## Escreva o teste primeiro

Adicione o cenário a seguir nos `casos`.

```go
{
	"Struct com dois campos tipo string",
	struct {
		Nome   string
		Cidade string
	}{"Chris", "Londres"},
	[]string{"Chris", "Londres"},
}
```

## Execute o teste

```text
=== RUN   TestPercorre/Struct_com_dois_campos_string
    --- FAIL: TestPercorre/Struct_com_dois_campos_string (0.00s)
        reflection_test.go:40: resultado [Chris], esperado [Chris Londres]
```

## Escreva código o suficiente para fazer o teste passar

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := reflect.ValueOf(x)

	for i := 0; i < valor.NumField(); i++ {
		campo := valor.Field(i)
		fn(campo.String())
	}
}
```

`valor` tem um método chamado `NumField` que retorna a quantidade de campos no valor. Isso nos permite iterar sobre os campos e chamar `fn`, o que faz nosso teste passar.

## Refatoração

Não parece haver nenhuma refatoração óbvia aqui que pode melhorar nosso código, então vamos continuar.

A próxima falha em `percorre` é que ela presume que todo campo é uma `string`. Vamos escrever um teste para esse caso.

## Escreva o teste primeiro

Inclua o seguinte cenário:

```go
{
	"Struct sem campo tipo string",
	struct {
		Nome  string
		Idade int
	}{"Chris", 33},
	[]string{"Chris"},
}
```

## Execute o teste

```text
=== RUN   TestPercorre/Struct_sem_campo_tipo_string
    --- FAIL: TestPercorre/Struct_with_noStruct_sem_campo_tipo_stringn_string_field (0.00s)
        reflection_test.go:46: resutado [Chris <int Value>], esperado [Chris]
```

## Escreva código o suficiente para fazer o teste passar

Precisamos verificar que o tipo do campo é uma `string`.

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := reflect.ValueOf(x)

	for i := 0; i < valor.NumField(); i++ {
		campo := valor.Field(i)

		if campo.Kind() == reflect.String { // Tipo
			fn(campo.String())
		}
	}
}
```

Podemos verificar seu tipo chamando a função [`Kind`](https://godoc.org/reflect#Kind).

## Refatoração

Parece que o código ainda está razoável por enquanto.

O próximo caso é: e se o valor não for uma `struct` "única"? Em outras palavras, o que acontece se tivermos uma `struct` com alguns campos aninhados?

## Escreva o teste primeiro

Estivemos usando a sintaxe de estrutura anônima para declarar tipos conforme precisávamos para nossos testes, então poderíamos continuar a fazer isso, como:

```go
{
    "Campos aninhados",
    struct {
        Nome string
        Perfil struct {
            Idade  int
            Cidade string
        }
    }{"Chris", struct {
        Idade  int
        Cidade string
    }{33, "Londres"}},
    []string{"Chris", "Londres"},
},
```

Mas podemos ver que quando você usa estruturas anônimas cada vez mais aninhadas, a sintaxe fica um pouco bagunçada. [Há uma proposta para fazer isso de forma que a sintaxe seja mais agradável](https://github.com/golang/go/issues/12854).

Vamos apenas refatorar isso criando um tipo conhecido para esse caso e referenciá-lo no nosso teste. Não é aconselhável colocar código do teste fora do teste, mas as pessoas devem ser capazes de encontrar essas estruturas procurando por sua definição.

Inclua as seguintes declarações de tipos no seu arquivo de teste:

```go
type Pessoa struct {
	Nome   string
	Perfil Perfil
}

type Perfil struct {
	Idade  int
	Cidade string
}
```

Agora podemos adicionar isso aos nossos casos ficarem bem mais legíveis que antes:

```go
{
	"Campos aninhados",
	Pessoa{
		"Chris",
		Perfil{33, "Londres"},
	},
	[]string{"Chris", "Londres"},
}
```

## Execute o teste

```text
=== RUN   TestPercorre/Campps_aninhados
    --- FAIL: TestPercorre/Campps_aninhados (0.00s)
        reflection_test.go:54: resultado [Chris], esperado [Chris Londres]
```

O problema é que estamos apenas iterando sobre os campos no primeiro nível da hierarquia de tipos.

## Escreva código o suficiente para fazer o teste passar

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := reflect.ValueOf(x)

	for i := 0; i < valor.NumField(); i++ {
		campo := valor.Field(i)

        if campo.Kind() == reflect.String {
            fn(campo.String())
        }

        if campo.Kind() == reflect.Struct {
            percorre(campo.Interface(), fn)
        }
    }
}
```

A solução é bem simples. Inspecionamos seu tipo novamente e se for uma estrutura apenas chamamos `percorre` novamente na nossa estrutura de dentro.

## Refatoração

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := reflect.ValueOf(x)

	for i := 0; i < valor.NumField(); i++ {
		campo := valor.Field(i)

		switch campo.Kind() {
		case reflect.String:
			fn(campo.String())
		case reflect.Struct:
			percorre(campo.Interface(), fn)
		}
	}
}
```

Quando você está fazendo uma comparação de mesmo valor mais de uma vez, _geralmente_ refatorar as condições dentro de um `switch` vai melhorar a legibilidade e tornar seu código mais fácil de estender.

E se o valor passado na estrutura for um ponteiro?

## Escreva o teste primeiro

Inclua esse caso:

```go
{
	"Ponteiros para coisas",
	&Pessoa{
		"Chris",
		Perfil{33, "Londres"},
	},
	[]string{"Chris", "Londres"},
}
```

## Execute o teste

```text
=== RUN   TestPercorre/Ponteiros_para_coisas
panic: reflect: call of reflect.Value.NumField on ptr Value [recovered]
    panic: reflect: call of reflect.Value.NumField on ptr Value
```

## Escreva código o suficiente para fazer o teste passar

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := reflect.ValueOf(x)

    if valor.Kind() == reflect.Ptr {
        valor = valor.Elem()
    }

    for i := 0; i < valor.NumField(); i++ {
		campo := valor.Field(i)

		switch campo.Kind() {
		case reflect.String:
			fn(campo.String())
		case reflect.Struct:
			percorre(campo.Interface(), fn)
		}
	}
}
```

Não é possível usar o `NumField` em um ponteiro `Value` e precisamos extrair o valor antes disso usando `Elem()`.

## Refatoração

Vamos encapsular a responsabilidade de extrair o `reflect.Value` de determinada `interface{}` para uma função.

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := obtemValor(x)

	for i := 0; i < valor.NumField(); i++ {
		campo := valor.Field(i)

		switch campo.Kind() {
		case reflect.String:
			fn(campo.String())
		case reflect.Struct:
			percorre(campo.Interface(), fn)
		}
	}
}

func obtemValor(x interface{}) reflect.Value {
	valor := reflect.ValueOf(x)

	if valor.Kind() == reflect.Ptr {
		valor = valor.Elem()
	}

	return valor
}
```

Isso acaba adicionando _mais_ código, mas me parece que o nível de abstração está correto.

* Obter o `reflect.Value` de `x` para que eu possa inspecioná-lo, não me importa de qual forma.
* Iterar pelos campos, fazendo o que for necessário dependendo de seu tipo.

Depois precisamos lidar com os slices.

## Escreva o teste primeiro

```go
{
	"Slices",
	[]Perfil{
		{33, "Londres"},
		{34, "Reykjavík"},
	},
	[]string{"Londres", "Reykjavík"},
}
```

## Execute o teste

```text
=== RUN   TestPercorre/Slices
panic: reflect: call of reflect.Value.NumField on slice Value [recovered]
    panic: reflect: call of reflect.Value.NumField on slice Value
```

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Esse caso se parece bastante com o do ponteiro acima, pois estamos chamar `NumField` em nosso `reflect.Value`, mas não há um por não ser uma struct.

## Escreva código o suficiente para fazer o teste passar

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := obtemValor(x)

    if valor.Kind() == reflect.Slice {
        for i:=0; i< valor.Len(); i++ {
            percorre(valor.Index(i).Interface(), fn)
        }
        return
    }

    for i := 0; i < valor.NumField(); i++ {
        campo := valor.Field(i)

        switch campo.Kind() {
        case reflect.String:
            fn(campo.String())
        case reflect.Struct:
            percorre(campo.Interface(), fn)
        }
    }
}
```

## Refatoração

Isso funciona, mas está bagunçado. Não se preocupe, pois temos cada pedaço de código coberto por testes e podemos brincar da forma que quisermos.

Se formos pensar um pouco abstradamente, queremos chamar `percorre` em:

* Cada campo de uma estrutura
* Cada _coisa_ de um slice

No momento nosso código faz isso, mas não reflete muito bem. Precisamos ter uma verificação no início da função para certificar se é um slice (com um `return` para parar a execução do restante do código) e se não for, só vamos presumir que é uma estrutura.

Vamos retrabalhar o código para verificar o tipo _primeiro_ para depois fazermos o que importa.

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := obtemValor(x)

    switch valor.Kind() {
    case reflect.Struct:
        for i:=0; i<valor.NumField(); i++ {
            percorre(valor.Field(i).Interface(), fn)
        }
    case reflect.Slice:
        for i:=0; i<valor.Len(); i++ {
            percorre(valor.Index(i).Interface(), fn)
        }
    case reflect.String:
        fn(valor.String())
    }
}
```

Parece muito melhor! Se for uma estrutura ou um slice, iteramos sobre seus valores chamando `percorre` para cada um. Por outro lado, se for um `reflect.String`, podemos apenas chamar `fn`.

Ainda assim me parece que poderia ficar melhor. Há repetição da operação de iterar sobre campos/valores e chamar `percorre` sendo que conceitualmente são a mesma coisa.

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := obtemValor(x)

	quantidadeDeValores := 0
	var obtemCampo func(int) reflect.Value

	switch valor.Kind() {
	case reflect.String:
		fn(valor.String())
	case reflect.Struct:
		quantidadeDeValores = valor.NumField()
		obtemCampo = valor.Field
	case reflect.Slice:
		quantidadeDeValores = valor.Len()
		obtemCampo = valor.Index
	}

	for i := 0; i < quantidadeDeValores; i++ {
		percorre(obtemCampo(i).Interface(), fn)
	}
}
```

Se o `valor` for um `reflect.String`, chamamos `fn` normalmente.

Se for outra coisa, nosso `switch` vai extrair duas coisas dependendo do tipo:

* Quantos campos existem
* Como extrair o `Value` (`Field` [campo] ou `Index` [índice])

Uma vez que determinamos esses pontos, podemos iterar pela `quantidadeDeValores` chamando `percorre` com o resultado da função `getField`.

A partir disso, lidar com arrays deve ser simples.

##  Escreva o teste primeiro

Inclua o caso:

```go
{
	"Arrays",
	[2]Perfil{
		{33, "Londres"},
		{34, "Reykjavík"},
	},
	[]string{"Londres", "Reykjavík"},
}
```

## Execute o teste

```text
=== RUN   TestPercorre/Arrays
    --- FAIL: TestPercorre/Arrays (0.00s)
        reflection_test.go:78: resultado [], esperado [Londres Reykjavík]
```

## Escreva código o suficiente para fazer o teste passar

Podemos resolver o caso dos arrays da mesma forma que os slices, basta adicioná-los com uma vírgula:

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := obtemValor(x)

	quantidadeDeValores := 0
	var obtemCampo func(int) reflect.Value

	switch valor.Kind() {
	case reflect.String:
		fn(valor.String())
	case reflect.Struct:
		quantidadeDeValores = valor.NumField()
		obtemCampo = valor.Field
	case reflect.Slice, reflect.Array:
		quantidadeDeValores = valor.Len()
		obtemCampo = valor.Index
	}

	for i := 0; i < quantidadeDeValores; i++ {
		percorre(obtemCampo(i).Interface(), fn)
	}
}
```

O último tipo que queremos lidar é o `map`.

## Escreva o teste primeiro

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

## Execute o teste

```text
=== RUN   TestPercorre/Maps
    --- FAIL: TestPercorre/Maps (0.00s)
        reflection_test.go:86: resultado [], esperado [Bar Boz]
```

## Escreva código o suficiente para fazer o teste passar

Novamente, se pensar um pouco de forma abstrata, percebe-se que o `map` é bem parecido com a `struct`, mas as chaves são desconhecidas em tempo de compilação.

Again if you think a little abstractly you can see that `map` is very similar to `struct`, it's just the keys are unknown at compile time.

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := obtemValor(x)

    quantidadedeValores := 0
    var obtemCampo func(int) reflect.Value

    switch valor.Kind() {
    case reflect.String:
        fn(valor.String())
    case reflect.Struct:
        quantidadeDeValores = valor.NumField()
        obtemCampo = valor.Field
    case reflect.Slice, reflect.Array:
        quantidadeDeValores = valor.Len()
        obtemCampo = valor.Index
    case reflect.Map:
        for _, chave := range valor.MapKeys() {
            percorre(valor.MapIndex(chave).Interface(), fn)
        }
    }

    for i := 0; i< quantidadeDeValores; i++ {
        percorre(obtemCampo(i).Interface(), fn)
    }
}
```

No entanto, por design, não é possível obter os valores de um map por índice. Só é possível fazer isso pela _chave_, que, caramba, acaba com a nossa abstração.

## Refatoração

Como se sente agora? Parecia que essa era uma boa abstração naquele momento, mas agora o código parece um pouco bagunçado.

_Está tudo bem!_ Refatoração é uma jornada e às vezes vamos cometer erros. Um ponto importante do TDD é que ele nos dá a liberdade de testar esse tipo de coisa.

Graças aos testes implmentados a cada etapa, essa situação não é irreversível de forma alguma. Vamos apenas voltar a como estava antes da refatoração.

```go
func percorre(x interface{}, fn func(entrada string)) {
	valor := obtemValor(x)

	percorreValor := func(valor reflect.Value) {
		percorre(valor.Interface(), fn)
	}

	switch valor.Kind() {
	case reflect.String:
		fn(valor.String())
	case reflect.Struct:
		for i := 0; i < valor.NumField(); i++ {
			percorreValor(valor.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < valor.Len(); i++ {
			percorreValor(valor.Index(i))
		}
	case reflect.Map:
		for _, chave := range valor.MapKeys() {
			percorreValor(valor.MapIndex(chave))
		}
	}
}
```

Apresentamos o `percorreValor`, que encapsula chamadas para `percorre` dentro do nosso `switch` para que só tenham que extrair os `reflect.Value` de `valor`.

### Um último problema

Lembre que maps em Go não têm ordem garantida. Logo, às vezes os testes irão falhar porque verificamos as chamadas de `fn` em uma ordem específica.

Para arrumar isso, precisaremos mover nossa verificação com os maps para um novo teste onde não nos importamos com a ordem.

```go
t.Run("com maps", func(t *testing.T) {
	mapA := map[string]string{
		"Foo": "Bar",
		"Baz": "Boz",
    }
    
	var resultado []string
	percorre(mapA, func(entrada string) {
		resultado = append(resultado, entrada)
    })
    
	verificaSeContem(t, resultado, "Bar")
	verificaSeContem(t, resultado, "Boz")
})
```

Essa é a definição de `verificaSeContem`:

```go
func verificaSeContem(t *testing.T, palheiro []string, agulha string) {
	contem := false
	for _, x := range palheiro {
		if x == agulha {
			contem = true
		}
	}
	if !contem {
		t.Errorf("esperava-se que %+v contivesse '%s', mas não continha", palheiro, agulha)
	}
}
```

## Resumo

* Apresentamos alguns dos conceitos do pacote `reflect`.
* Usamos recursão para percorrer estruturas de dados arbitrárias.
* Houve uma reflexão quanto a uma refatoração ruim, mas não há por que se preocupar muito com isso. Isso não deve ser um problema muito grande se trabalharmos com testes de forma iterativa.
* Esse capítulo só cobre um aspecto pequeno de reflexão. [O blog do Go tem um artigo excelente cobrindo mais detalhes](https://blog.golang.org/laws-of-reflection).
* Agora que você tem conhecimento sobre reflexão, faça o possível para evitá-lo.

