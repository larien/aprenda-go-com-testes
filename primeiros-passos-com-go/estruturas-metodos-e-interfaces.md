# Estruturas, métodos e interfaces

[**Você pode encontrar todos os códigos desse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/estruturas-metodos-e-interfaces)

Suponha que precisamos de algum código de geometria para calcular o perímetro de um retângulo dado uma altura e largura. Podemos escrever uma função `Perimetro(largura float64, altura float64)`, onde `float64` representa números em ponto flutuante como `123.45`.

O ciclo de TDD deve ser mais familiar para você agora.

## Escreva o teste primeiro

```go
func TestPerimetro(t *testing.T) {
	resultado := Perimetro(10.0, 10.0)
	esperado := 40.0

	if resultado != esperado {
		t.Errorf("resultado %.2f esperado %.2f", resultado, esperado)
	}
}
```

Viu a nova string de formatação? O `f` é para nosso `float64` e o `.2` significa imprimir duas casas decimais.

## Execute o teste

`./formas_test.go:6:9: undefined: Perimetro`

`indefinido: Perimetro`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhando

```go
func Perimetro(largura float64, altura float64) float64 {
    return 0
}
```

Resulta em `formas_test.go:10: resultado 0, esperado 40`.

## Escreva código o suficiente para fazer o teste passar

```go
func Perimetro(largura float64, altura float64) float64 {
	return 2 * (largura + altura)
}
```

Por enquanto, tudo fácil. Agora vamos criar uma função chamada `Area(largura, altura float64)` que retorna a área de um retângulo.

Tente fazer isso sozinho, segundo o ciclo de TDD.

Você deve terminar com os testes como estes:

```go
func TestPerimetro(t *testing.T) {
	resultado := Perimetro(10.0, 10.0)
	esperado := 40.0

	if resultado != esperado {
		t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
	}
}

func TestArea(t *testing.T) {
	resultado := Area(12.0, 6.0)
	esperado := 72.0

	if resultado != esperado {
		t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
	}
}
```

E código como este:

```go
func Perimetro(largura float64, altura float64) float64 {
	return 2 * (largura + altura)
}

func Area(largura float64, altura float64) float64 {
	return largura * altura
}
```

## Refatoração

Nosso código faz o trabalho, mas não contém nada explícito sobre retângulos. Uma pessoa descuidada poderia tentar passar a largura e altura de um triângulo para esta função sem perceber que ela retornará uma resposta errada.

Podemos apenas dar para a função um nome mais específico como `AreaDoRetangulo`. Uma solução mais limpa é definir nosso próprio _tipo_ chamado `Retangulo` que encapsula este conceito para nós.

Podemos criar um tipo simples usando uma **struct** (estrutura). [Uma struct](https://golang.org/ref/spec#Struct_types) é apenas uma coleção nomeada de campos onde você pode armazenar dados.

Declare uma `struct` assim:

```go
type Retangulo struct {
	Largura float64
	Altura  float64
}
```

Agora vamos refatorar os testes para usar `Retangulo` em vez de um simples `float64`.

```go
func TestPerimetro(t *testing.T) {
	retangulo := Retangulo{10.0, 10.0}
	resultado := Perimetro(retangulo)
	esperado := 40.0

	if resultado != esperado {
		t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
	}
}

func TestArea(t *testing.T) {
	retangulo := Retangulo{12.0, 6.0}
	resultado := Area(retangulo)
	esperado := 72.0

	if resultado != esperado {
		t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
	}
}
```

Lembre de rodar seus testes antes de tentar corrigir. Você deve ter erro útil como:

```text
./formas_test.go:7:18: not enough arguments in call to Perimetro
    have (Retangulo)
    esperado (float64, float64)
```

Você pode acessar os campos de uma `struct` com a sintaxe `minhaStruct.campo`.

Mude as duas funções para corrigir o teste.

```go
func Perimetro(retangulo Retangulo) float64 {
	return 2 * (retangulo.Largura + retangulo.Altura)
}

func Area(retangulo Retangulo) float64 {
	return retangulo.Largura * retangulo.Altura
}
```

Espero que você concorde que passar um `Retangulo` para a função mostra nossa intenção com mais clareza, mas existem mais benefícios em usar `structs` que já vamos entender.

Nosso próximo requisito é escrever uma função `Area` para círculos.

## Escreva o teste primeiro

```go
func TestArea(t *testing.T) {
	t.Run("retângulos", func(t *testing.T) {
		retangulo := Retangulo{12.0, 6.0}
		resultado := Area(retangulo)
		esperado := 72.0

		if resultado != esperado {
			t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
		}
	})

	t.Run("círculos", func(t *testing.T) {
		circulo := Circulo{10}
		resultado := Area(circulo)
		esperado := 314.1592653589793

		if resultado != esperado {
			t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
		}
	})
}
```

## Execute o teste

`./formas_test.go:28:13: undefined: Circulo`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhando

Precisamos definir nosso tipo `Circulo`.

```go
type Circulo struct {
	Raio float64
}
```

Agora rode os testes novamente.

`./formas_test.go:29:14: cannot use circulo (type Circulo) as type Retangulo in argument to Area`

Algumas linguagens de programação permitem você fazer algo como:

```go
func Area(circulo Circulo) float64 { ... }
func Area(retangulo Retangulo) float64 { ... }
```

Mas em Go você não pode:

`./formas.go:20:32: Area redeclared in this block`

Temos duas escolhas:

* Podemos ter funções com o mesmo nome declaradas em _pacotes_ diferentes. Então, poderíamos criar nossa `Area(Circulo)` em um novo _pacote_, só que isso parece um exagero aqui.
* Em vez disso, podemos definir [_métodos_](https://golang.org/ref/spec#Method_declarations) em nosso mais novo tipo definido.

### O que são métodos?

Até agora só escrevemos _funções_, mas temos usado alguns métodos. Quando chamamos `t.Errorf`, nós chamamos o método `Errorf` na instância de nosso `t` \(`testing.T`\).

Um método é uma função com um receptor. Uma declaração de método vincula um identificador e o nome do método a um método e associa o método com o tipo base do receptor.

Métodos são muito parecidos com funções, mas são chamados invocando-os em uma instância de um tipo específico.

Enquanto você chama funções onde quiser, como por exemplo em `Area(retangulo)`, você só pode chamar métodos em "coisas" específicas.

Um exemplo ajudará. Então, vamos mudar nossos testes primeiro para chamar métodos em vez de funções, e, em seguida, corrigir o código.

```go
func TestArea(t *testing.T) {
	t.Run("retângulos", func(t *testing.T) {
		retangulo := Retangulo{12.0, 6.0}
		resultado := retangulo.Area()
		esperado := 72.0

		if resultado != esperado {
			t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
		}
	})

	t.Run("círculos", func(t *testing.T) {
		circulo := Circulo{10}
		resultado := circulo.Area()
		esperado := 314.1592653589793

		if resultado != esperado {
			t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
		}
	})
}
```

Se rodarmos os testes agora, recebemos:

```text
./formas_test.go:19:19: retangulo.Area undefined (type Retangulo has no field or method Area)
./formas_test.go:29:16: circulo.Area undefined (type Circulo has no field or method Area)
```

> type Circulo has no field or method Area

Gostaria de reforçar o quão grandioso o compilador é. É muito importante ter tempo para ler lentamente as mensagens de erro que você recebe, pois isso te ajudará a longo prazo.

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhando

Vamos adicionar alguns métodos para nossos tipos:

```go
type Retangulo struct {
    Largura  float64
    Altura float64
}

func (r Retangulo) Area() float64  {
    return 0
}

type Circulo struct {
    Raio float64
}

func (c Circulo) Area() float64  {
    return 0
}
```

A sintaxe para declaração de métodos é quase a mesma que usamos para funções e isso acontece porque eles são muito parecidos. A única diferença é a sintaxe para o método receptor: `func (nomeDoReceptor TipoDoReceptor) NomeDoMetodo(argumentos)`.

Quando seu método é chamado em uma variável desse tipo, você tem sua referência para o dado através da variável `nomeDoReceptor`. Em muitas outras linguagens de programação isto é feito implicitamente e você acessa o receptor através de `this`.

É uma convenção em Go que a variável receptora seja a primeira letra do tipo em minúsculo.

```go
r Retangulo
```

Se você executar novamente os testes, eles devem compilar e dar alguma saída do teste falhando.

## Escreva código suficiente para fazer o teste passar

Agora vamos fazer nossos testes de retângulo passarem corrigindo nosso novo método.

```go
func (r Retangulo) Area() float64  {
    return r.Largura * r.Altura
}
```

Se você executar novamente os testes, aqueles de retângulo devem passar, mas os de círculo ainda falham.

Para fazer a função `Area` de círculo passar, vamos emprestar a constante `Pi` do pacote `math` \(lembre-se de importá-lo\).

```go
func (c Circulo) Area() float64  {
    return math.Pi * c.Raio * c.Raio
}
```

## Refatoração

Existe duplicação em nossos testes.

Tudo o que queremos fazer é pegar uma coleção de _formas_, chamar o método `Area()` e então verificar o resultado.

Queremos ser capazes de escrever um tipo de função `verificaArea` que permita passar tanto `Retangulo` quanto `Circulo`, mas falhe ao compilar se tentarmos passar algo que não seja uma _forma_.

Com Go, podemos trabalhar dessa forma com **interfaces**.

[Interfaces](https://golang.org/ref/spec#Interface_types) são um conceito muito poderoso em linguagens de programação estaticamente tipadas, como Go, porque permitem que você crie funções que podem ser usadas com diferentes tipos e permite a criação de código altamente desacoplado, mantendo ainda a segurança de tipos.

Vamos apresentar isso refatorando nossos testes.

```go
func TestArea(t *testing.T) {
	verificaArea := func(t *testing.T, forma Forma, esperado float64) {
		t.Helper()
		resultado := forma.Area()

		if resultado != esperado {
			t.Errorf("resultado %.2f, esperado %.2f", resultado, esperado)
		}
	}

	t.Run("retângulos", func(t *testing.T) {
		retangulo := Retangulo{12.0, 6.0}
		verificaArea(t, retangulo, 72.0)
	})

	t.Run("círculos", func(t *testing.T) {
		circulo := Circulo{10}
		verificaArea(t, circulo, 314.1592653589793)
	})
}
```

Estamos criando uma função auxiliar como fizemos em outros exercícios, mas desta vez estamos pedindo que uma `Forma` seja passada. Se tentarmos chamá-la com algo que não seja uma _forma_, não vai compilar.

Como algo se torna uma _forma_? Precisamos apenas falar para o Go o que é uma `Forma` usando uma declaração de interface.

```go
type Forma interface {
	Area() float64
}
```

Estamos criando um novo `tipo`, assim como fizemos com `Retangulo` e `Circulo`, mas desta vez é uma `interface` em vez de uma `struct`.

Uma vez adicionado isso ao código, os testes passarão.

### Peraí, como assim?

A interface em Go bem diferente das interfaces na maioria das outras linguagens de programação. Normalmente você tem que escrever um código para dizer que `meu tipo Foo implementa a interface Bar`.

Só que no nosso caso:

* `Retangulo` tem um método chamado `Area` que retorna um `float64`, então satisfaz a interface `Forma`.
* `Circulo` tem um método chamado `Area` que retorna um `float64`, então satisfaz a interface `Forma`.
* `string` não tem esse método, então não satisfaz a interface.
* etc.

Em Go a **resolução de interface é implícita**. Se o tipo que você passar combinar com o que a interface está esperando, o código será compilado.

### Desacoplando

Veja como nossa função auxiliar não precisa se preocupar se a _forma_ é um `Retangulo` ou um `Circulo` ou um `Triangulo`. Ao declarar uma interface, a função auxiliar está _desacoplada_ de tipos concretos e tem apenas o método que precisa para fazer o trabalho.

Este tipo de abordagem - de usar interfaces para declarar **somente o que você precisa** - é muito importante no desenvolvimento de software e será coberto mais detalhadamente nas próximas seções.

## Refatoração adicional

Agora que você conhece as `structs`, podemos apresentar os "table driven tests" (testes orientados por tabela).

[Table driven tests](https://github.com/golang/go/wiki/TableDrivenTests) são úteis quando você quer construir uma lista de casos de testes que podem ser testados da mesma forma.

```go
func TestArea(t *testing.T) {
	testesArea := []struct {
		forma    Forma
		esperado float64
	}{
		{Retangulo{12, 6}, 72.0},
		{Circulo{10}, 314.1592653589793},
	}

	for _, tt := range testesArea {
		resultado := tt.forma.Area()
		if resultado != tt.esperado {
			t.Errorf("resultado %.2f, esperado %.2f", resultado, tt.esperado)
		}
	}
}
```

A única sintaxe nova aqui é a criação de uma "struct anônima", `testesArea`. Estamos declarando um slice de structs usando `[]struct` com dois campos, o `forma` e o `esperado`. Então preenchemos o slice com os casos.

Depois iteramos sobre eles assim como fazemos com qualquer outro slice, usando os campos da struct para executar nossos testes.

Dá para perceber como será muito fácil para uma pessoa inserir uma nova forma, implementar `Area` e então adicioná-la nos casos de teste. Além disso, se for encontrada uma falha em `Area`, é muito fácil adicionar um novo caso de teste para verificar antes de corrigi-la.

_Testes baseados em tabela_ podem ser um item valioso em sua caixa de ferramentas, mas tenha certeza de que você precisa da sintaxe extra nos testes. Se você deseja testar várias implementações de uma interface ou se o dado passado para uma função tem muitos requisitos diferentes que precisam de testes, eles podem servir bem.

Vamos demonstrar tudo isso adicionando e testando outra forma; um triângulo.

## Escreva o teste primeiro

Adicionar um teste para nossa nova forma é muito fácil. Simplesmente adicione `{Triangulo{12, 6}, 36.0},` à nossa lista.

```go
func TestArea(t *testing.T) {
	testesArea := []struct {
		forma    Forma
		esperado float64
	}{
		{Retangulo{12, 6}, 72.0},
		{Circulo{10}, 314.1592653589793},
		{Triangulo{12, 6}, 36.0},
	}

	for _, tt := range testesArea {
		resultado := tt.forma.Area()
		if resultado != tt.esperado {
			t.Errorf("resultado %.2f, esperado %.2f", resultado, tt.esperado)
		}
	}
}
```

## Execute o teste

Lembre-se, continue tentando executar o teste e deixe o compilador guiá-lo em direção a solução.

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste falhando

`./formas_test.go:25:4: undefined: Triangulo`

Ainda não definimos `Triangulo`:

```go
type Triangulo struct {
    Base   float64
    Altura float64
}
```

Tente novamente:

```text
./formas_test.go:25:8: cannot use Triangulo literal (type Triangulo) as type Forma in field value:
    Triangulo does not implement Forma (missing Area method)
```
`Triangulo não implementa Forma (método Area faltando)`

Isso nos diz que não podemos usar um `Triangulo` como uma `Forma` porque ele não tem um método `Area()`, então adicione uma implementação vazia para fazermos o teste funcionar:

```go
func (t Triangulo) Area() float64 {
    return 0
}
```

Finalmente o código compilou e temos o nosso erro:

`formas_test.go:31: resultado 0.00, esperado 36.00`

## Escreva código suficiente para fazer o teste passar

```go
func (t Triangulo) Area() float64 {
    return (t.Base * t.Altura) * 0.5
}
```

E nossos testes passaram!

## Refatoração

Novamente, a implementação está boa, mas nossos testes podem ser melhorados.

Quando você lê isso:

```go
{Retangulo{12, 6}, 72.0},
{Circulo{10}, 314.1592653589793},
{Triangulo{12, 6}, 36.0},
```

Não está tão claro o que todos os números representam e você deve ter o objetivo de escrever testes que sejam fáceis de entender.

Até agora você viu uma sintaxe para criar instâncias de structs como `MinhaStruct{valor1, valor2}`, mas você pode opcionalmente nomear esses campos.

Vamos ver como isso funciona:

```go
        {forma: Retangulo{largura: 12, altura: 6}, esperado: 72.0},
        {forma: Circulo{Raio: 10}, esperado: 314.1592653589793},
        {forma: Triangulo{Base: 12, altura: 6}, esperado: 36.0},
```

Em [Test-Driven Development by Example](https://g.co/kgs/yCzDLF) Kent Beck refatora alguns testes para um ponto e afirma:

> O teste é lido de forma mais clara, como se fosse uma afirmação da verdade, **não uma sequência de operações**

\(ênfase minha\)

Agora nossos testes \(pelo menos a lista de casos\) fazem afirmações da verdade sobre formas e suas áreas.

## Garanta que a saída do seu teste seja útil

Lembra anteriormente quando implementamos `Triangulo` e tivemos um teste falhando? Ele imprimiu `formas_test.go:31: resultado 0.00 esperado, 36.00`.

Nós sabíamos que estava relacionado ao `Triangulo` porque estávamos trabalhando nisso, mas e se uma falha escorregasse para o sistema em um dos 20 casos na tabela? Como alguém saberia qual caso falhou? Não parece ser uma boa experiência. Ela teria que olhar caso a caso para encontrar qual deles está falhando de fato.

Podemos mudar nossa mensagem de erro para `%#v resultado %.2f, esperado %.2f`. A string de formatação `%#v` irá imprimir nossa struct com os valores em seu campo para que as pessoas possam ver imediatamente as propriedades que estão sendo testadas.

Para melhorar a legibilidade de nossos futuros casos de teste, podemos renomear o campo `esperado` para algo mais descritivo como `temArea`.

Uma dica final com testes guiados por tabela é usar `t.Run` e renomear os casos de teste.

Envolvendo cada caso em um `t.Run` você terá uma saída de testes mais limpa em caso de falhas, além de imprimir o nome do caso.

```text
--- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/Retangulo (0.00s)
        formas_test.go:33: main.Retangulo{Largura:12, Altura:6} resultado 72.00, esperado 72.10
```

E você pode rodar testes específicos dentro de sua tabela com `go test -run TestArea/Retangulo`.

Aqui está o código final do nosso teste que captura isso:

```go
func TestArea(t *testing.T) {
	testesArea := []struct {
		nome    string
		forma   Forma
		temArea float64
	}{
		{nome: "Retângulo", forma: Retangulo{Largura: 12, Altura: 6}, temArea: 72.0},
		{nome: "Círculo", forma: Circulo{Raio: 10}, temArea: 314.1592653589793},
		{nome: "Triângulo", forma: Triangulo{Base: 12, Altura: 6}, temArea: 36.0},
	}

	for _, tt := range testesArea {
		t.Run(tt.nome, func(t *testing.T) {
			resultado := tt.forma.Area()
			if resultado != tt.temArea {
				t.Errorf("%#v resultado %.2f, esperado %.2f", tt.forma, resultado, tt.temArea)
			}
		})
	}
}
```

## Resumo

Esta foi mais uma prática de TDD, iterando em nossas soluções para problemas matemáticos básicos e aprendendo novos recursos da linguagem motivados por nossos testes.

* Declarar structs para criar seus próprios tipos de dados permite agrupar dados relacionados e torna a intenção do seu código mais clara.
* Declarar interfaces permite que você possa definir funções que podem ser usadas por diferentes tipos \([polimorfismo paramétrico](https://pt.wikipedia.org/wiki/Polimorfismo_paramétrico)\).
* Adicionar métodos permite que você possa adicionar funcionalidades aos seus tipos de dados e implementar interfaces.
* Testes baseados em tabela permite que você torne suas asserções mais claras e seus testes mais fáceis de estender e manter.

Este foi um capítulo importante porque agora começamos a definir nossos próprios tipos. Em linguagens estaticamente tipadas como Go, conseguir projetar seus próprios tipos é essencial para construir software que seja fácil de entender, compilar e testar.

Interfaces são uma ótima ferramenta para ocultar a complexidade de outras partes do sistema. Em nosso caso, o _código_ de teste auxiliar não precisou conhecer a forma exata que estava afirmando, apenas como "pedir" pela sua área.

Conforme você se familiariza com Go, começa a ver a força real das interfaces e da biblioteca padrão.

Você aprenderá sobre as interfaces definidas na biblioteca padrão que são usadas _em todo lugar_ e, implementando-as em relação aos seus próprios tipos, você pode reutilizar rapidamente muitas das ótimas funcionalidades.
