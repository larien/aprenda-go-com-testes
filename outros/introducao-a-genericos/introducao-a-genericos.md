# Introdução a genéricos

[**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/main/outros/introducao-a-genericos)

(No momento da escrita) Go não possui suporte para genéricos definidos pelo usuário, mas [a proposta de inclusão](https://blog.golang.org/generics-proposal) [foi aceita](https://github.com/golang/go/issues/43651#issuecomment-776944155) e será incluído na versão 1.18.

No entanto, existem maneiras de experimentar _hoje_ a implementação futura usando o [playground go2go](https://go2goplay.golang.org/). Portanto, para trabalhar neste capítulo, você terá que deixar seu precioso editor de escolha e, em vez disso, fazer o trabalho dentro do playground.

Este capítulo lhe dará uma introdução aos genéricos, dissipará os receios que você possa ter sobre eles e dará uma ideia de como simplificar parte do seu código no futuro. Após ler isso, você saberá como escrever:

- Uma função que recebe argumentos genéricos
- Uma estrutura de dados genérica

## Configurando o playground

No playground _go2go_ não podemos executar o comando `go test`. Como vamos escrever testes para explorar o código genérico?

O playground _não_ nos permite executar código e, como somos programadores, isso significa que podemos contornar a falta de um executor de teste **criando o nosso próprio**.

## Nossos próprios auxiliares de teste (`VerificaIgual` (AssertEqual),`VerificaNaoIgual`(AssertNotEqual))

Para explorar os genéricos, escreveremos alguns auxiliares de teste que matarão o programa e imprimirão algo útil se um teste falhar.

### Verificar inteiros

Vamos começar com algo básico e iterar em direção ao nosso objetivo

```go
package main

import (
    "log"
)

func main() {
    VerificaIgual(1, 1)
    VerificaNaoIgual(1, 2)

    VerificaIgual(50, 100) // isso deve falhar

    VerificaNaoIgual(2, 2) // você não verá isso na impressão (print)
}

func VerificaIgual(recebido, esperado int) {
    if recebido != esperado {
        log.Fatalf("resultado: recebido %d, esperado %d", recebido, esperado)
    } else {
        log.Printf("PASSOU: %d é igual %d\n", recebido, esperado)
    }
}

func VerificaNaoIgual(recebido, esperado int) {
    if recebido == esperado {
        log.Fatalf("FALHOU: recebido %d, esperado %d", recebido, esperado)
    } else {
        log.Printf("PASSOU: %d não é igual %d\n", recebido, esperado)
    }
}
```

[This program prints](https://go2goplay.golang.org/p/YoJZhWMYO):

```
2009/11/10 23:00:00 PASSOU: 1 é igual 1
2009/11/10 23:00:00 PASSOU: 1 não é igual 2
2009/11/10 23:00:00 resultado: recebido 50, esperado 100
```

### Verificar textos

Ser capaz de verificar a igualdade de inteiros é ótimo, mas e se quisermos verificar algum `texto`?

```go
func main() {
    VerificaIgual("CJ", "CJ")
}
```

Você obterá um erro

```
type checking failed for main
prog.go2:15:16: cannot use "CJ" (untyped string constant) as int value in argument to VerificaIgual
prog.go2:15:22: cannot use "CJ" (untyped string constant) as int value in argument to VerificaIgual
```
_ Traduzindo a mensagem do log: "não pode usar "CJ" (constante de texto não digitada) como valor de int (inteiro) no argumento para VerificaIgual"_

Se você parar para ler o erro com calma, verá que o compilador está reclamando que estamos tentando passar um `texto` para uma função que espera um valor `inteiro`.

#### Recapitulação sobre segurança de tipos

Se você leu os capítulos anteriores deste livro, ou tem experiência com linguagens estaticamente tipadas, isso não deve surpreendê-lo. O compilador de Go espera que você escreva suas funções, estruturas e etc, descrevendo com quais tipos deseja trabalhar.

Você não pode passar um `texto` para uma função que espera um valor `inteiro`.

Embora possa parecer uma cerimônia, pode ser extremamente útil. Ao descrever essas restrições, você:

- Simplifica a implementação da função. Ao descrever ao compilador com quais tipos você trabalha, você **restringe o número de possíveis implementações válidas**. Você não pode "adicionar" uma `Pessoa` e uma `Conta bancária`. Você não pode colocar um `inteiro` em maiúscula. Em software, as restrições costumam ser extremamente úteis.
- Não poderá passar acidentalmente dados para uma função que você não pretendia.

Go atualmente oferece uma maneira de ser mais abstrato com seus tipos com interfaces, para que você possa projetar funções que não aceitam tipos concretos, mas sim tipos que oferecem o comportamento de que você precisa. Isso dá a você alguma flexibilidade enquanto mantém a segurança do tipo.

### Uma função que recebe um texto ou um inteiro? (ou de fato, outras coisas)

A outra opção que Go oferece _atualmente_ é declarar o tipo de seu argumento como `interface{}` que significa "qualquer coisa".

Tente alterar as assinaturas dos metódo para usar este tipo.

```go
func VerificaIgual(recebido, esperado interface{}) {

func VerificaNaoIgual(recebido, esperado interface{}) {

```

Os testes agora devem ser compilados e aprovados. A saída será um pouco confusa porque estamos usando a string de formato inteiro `%d` para imprimir nossas mensagens, então mude-as para o formato geral `%+v` para uma melhor saída de qualquer tipo de valor.

### Trocas feitas sem genéricos

As funções do nosso `VerificaAlgo` são bastante ingênuas, mas conceitualmente não são muito diferentes de como outras [bibliotecas populares oferecem essa funcionalidade](https://github.com/matryer/is/blob/master/is.go#L150)

```go
func (is *I) VerificaIgual(a, b interface{}) {
```

Então qual é o problema?

Ao usar a `interface{}`, o compilador não pode nos ajudar a escrever nosso código, porque não estamos dizendo a ele nada de útil sobre os tipos de coisas passadas para a função. Volte para o _go2go playground_ e tente comparar dois tipos diferentes,

```go
VerificaNaoIgual(1, "1")
```

Nesse caso, escapamos impunes; o teste compila e falha como esperávamos, embora a mensagem de erro `recebido 1, esperado 1` não esteja clara; mas queremos ser capazes de comparar strings com inteiros? Que tal comparar uma `Pessoa` com um `Aeroporto`?

Escrever funções que usam `interface{}` pode ser extremamente desafiador e sujeito a bugs porque _perdemos_ nossas restrições, e não temos nenhuma informação em tempo de compilação sobre os tipos de dados com os quais estamos lidando.

Frequentemente, os desenvolvedores precisam refletir para implementar essas funções *ahem (som de limpando a garganta)* genéricas, o que geralmente é doloroso e pode prejudicar o desempenho do seu programa.

## Nossos próprios auxiliares de teste com genéricos

Idealmente, não queremos ter que fazer funções `VerificaAlgo` específicas para cada tipo com que lidamos. Gostaríamos de poder ter _uma_ função `VerificaIgual` que funcione com _qualquer_ tipo, mas não permite que você compare [maçãs e laranjas](https://en.wikipedia.org/wiki/Apples_and_oranges).

Os genéricos nos oferecem uma nova maneira de fazer abstrações (como interfaces), permitindo-nos **descrever nossas restrições** de maneiras que não podemos fazer atualmente.

```go
package main

import (
    "log"
)

func main() {
    VerificaIgual(1, 1)
    VerificaIgual("1", "1")
    VerificaNaoIgual(1, 2)
 // VerificaIgual(1, "1") - descomente-me para ver o erro de compilação
}

func VerificaIgual[T comparable](recebido, esperado T) {
    if recebido != esperado {
        log.Fatalf("FALHOU: recebido %+v, esperado %+v", recebido, esperado)
    } else {
        log.Printf("PASSOU: %+v é igual %+v\n", recebido, esperado)
    }
}

func VerificaNaoIgual[T comparable](recebido, esperado T) {
    if recebido == esperado {
        log.Fatalf("FALHOU: recebido %+v, esperado %+v", recebido, esperado)
    } else {
        log.Printf("PASSOU: %+v não é igual  %+v\n", recebido, esperado)
    }
}
```

[link do go2go playground](https://go2goplay.golang.org/p/k8UujhlduJU)

Para escrever funções genéricas em Go, você precisa fornecer "parâmetros com tipos", que é apenas uma maneira elegante de dizer "descreva seu tipo genérico e dê a ele um rótulo".

Em nosso caso, o tipo de nosso parâmetro de tipo é [`comparable`](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md#comparable-types-in-constraints) e demos a ele o rótulo de `T`. Este rótulo nos permite descrever os tipos de argumentos para nossa função (`recebido, esperado T`).

Estamos usando `comparable` porque queremos descrever para o compilador que desejamos usar os operadores `==`e `!=` Em coisas do tipo `T` em nossa função, queremos comparar! Se você tentar mudar o tipo para `any`,

```go
func VerificaNaoIgual[T any](recebido, esperado T) {
```

Você obterá o seguinte erro:

```
prog.go2:24:8: cannot compare recebido == esperado (operator == not defined for T)
```

O que faz muito sentido, porque você não pode usar esses operadores em todos (ou `any`) tipo.

### [`Any`](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md#the-constraint) é o mesmo que `interface{}` ?

Considere duas funções

```go
func GenericoFoo[T any](x, y T)
```

```go
func InterfaceFoo(x, y interface{})
```

Qual é o objetivo dos genéricos aqui? `any` não descreve ... nada?

Em termos de restrições, `any` significa "qualquer coisa" assim como `interface{}`. A diferença com a versão genérica é _você ainda está descrevendo um tipo específico_ e o que isso significa é que ainda restringimos esta função para funcionar apenas com _um_ tipo.

Isso significa que você pode chamar `InterfaceFoo` com qualquer combinação de tipos (por exemplo, `InterfaceFoo(maçã, laranja)`). No entanto, `GenericoFoo` ainda oferece algumas restrições porque dissemos que ele só funciona com _um_ tipo, `T`.

Válido:

- `GenericoFoo(maçã1, maçã2)`
- `GenericoFoo(laranja1, laranja2)`
- `GenericoFoo(1, 2)`
- `GenericoFoo("um", "dois")`

Não é válido (falha na compilação):

- `GenericoFoo(maçã1, laranja1)`
- `GenericoFoo("1", 1)`

`any` é especialmente útil ao criar tipos de dados onde você deseja que funcionem com vários tipos, mas você não _utiliza_ o tipo em sua própria estrutura de dados (normalmente, você está apenas armazenando-o). Coisas como `Set` e `LinkedList` são boas candidatas para usar `any`.

## Próximo Tópico: Tipos de dados genéricos

Vamos criar um tipo de dados [stack (pilha)](https://en.wikipedia.org/wiki/Stack_(abstract_data_type)). As pilhas (stacks) devem ser bastante simples de entender do ponto de vista dos requisitos. Eles são uma coleção de itens onde você pode `Empilhar (Push)` itens para o "topo" e para obter os itens de volta você `Desempilhar (Pop)` itens do topo (LIFO (last in, first out) - último a entrar, primeiro a sair).

Para ser breve, omiti o processo TDD que me chegou ao [seguinte código](https://go2goplay.golang.org/p/HghXymv1OKm) para uma pilha de `inteiro`s, e uma pilha de `texto`s.

```go
package main

import (
	"log"
)

type PilhaDeInteiros struct {
	valores []int
}

func (p *PilhaDeInteiros) Empilhar(valor int) {
	p.valores = append(p.valores, valor)
}

func (p *PilhaDeInteiros) EstaVazio() bool {
	return len(p.valores) == 0
}

func (p *PilhaDeInteiros) Desempilhar() (int, bool) {
	if p.EstaVazio() {
		return 0, false
	}

	indice := len(p.valores) - 1
	el := p.valores[indice]
	p.valores = p.valores[:indice]
	return el, true
}

type PilhaDeTextos struct {
	valores []string
}

func (p *PilhaDeTextos) Empilhar(valor string) {
	p.valores = append(p.valores, valor)
}

func (p *PilhaDeTextos) EstaVazio() bool {
	return len(p.valores) == 0
}

func (p *PilhaDeTextos) Desempilhar() (string, bool) {
	if p.EstaVazio() {
		return "", false
	}

	indice := len(p.valores) - 1
	el := p.valores[indice]
	p.valores = p.valores[:indice]
	return el, true
}

func main() {
	// PILHA DE INTEIROS

	minhaPilhaDeInteiros := new(PilhaDeInteiros)

	// verifica se a pilha está vazia
	VerificaVerdadeiro(minhaPilhaDeInteiros.EstaVazio())

	// adiciona alguma coisa e em seguida, verifica se a pilha não está vazia
	minhaPilhaDeInteiros.Empilhar(123)
	VerificaFalso(minhaPilhaDeInteiros.EstaVazio())

	// adiciona outra coisa e em seguida, desempilhe a pilha
	minhaPilhaDeInteiros.Empilhar(456)
	valor, _ := minhaPilhaDeInteiros.Desempilhar()
	VerificaIgual(valor, 456)
	valor, _ = minhaPilhaDeInteiros.Desempilhar()
	VerificaIgual(valor, 123)
	VerificaVerdadeiro(minhaPilhaDeInteiros.EstaVazio())

	// PILHA DE TEXTOS

	minhaPilhaDeTextos := new(PilhaDeTextos)

	// verifica se a pilha está vazia
	VerificaVerdadeiro(minhaPilhaDeTextos.EstaVazio())

	// adiciona alguma coisa e em seguida, verifica se a pilha não está vazia
	minhaPilhaDeTextos.Empilhar("um dois tres")
	VerificaFalso(minhaPilhaDeTextos.EstaVazio())

	// adiciona outra coisa e em seguida, desempilhe a pilha
	minhaPilhaDeTextos.Empilhar("quatro cinco seis")
	valorTexto, _ := minhaPilhaDeTextos.Desempilhar()
	VerificaIgual(valorTexto, "quatro cinco seis")
	valorTexto, _ = minhaPilhaDeTextos.Desempilhar()
	VerificaIgual(valorTexto, "um dois tres")
	VerificaVerdadeiro(minhaPilhaDeTextos.EstaVazio())
}

func VerificaVerdadeiro(algo bool) {
    if algo {
        log.Printf("PASSOU: Esperava-se que fosse verdade e foi\n")
    } else {
        log.Fatalf("FALHOU: Esperava-se que fosse verdadeiro, mas foi falso")
    }
}

func VerificaFalso(algo bool) {
    if !algo {
        log.Printf("PASSOU: Esperava-se que fosse falso e foi\n")
    } else {
        log.Fatalf("FALHOU: Esperava-se que fosse falso mas foi verdadeiro")
    }
}

func VerificaIgual[T comparable](recebido, esperado T) {
    if recebido != esperado {
        log.Fatalf("FALHOU: recebido %+v, esperado %+v", recebido, esperado)
    } else {
        log.Printf("PASSOU: %+v é igual  %+v\n", recebido, esperado)
    }
}

func VerificaNaoIgual[T comparable](recebido, esperado T) {
    if recebido == esperado {
        log.Fatalf("FALHOU: recebido %+v, esperado %+v", recebido, esperado)
    } else {
        log.Printf("PASSOU: %+v não é igual  %+v\n", recebido, esperado)
    }
}
```

### Problemas

- O código para `PilhaDeTextos` e `PilhaDeInteiros` são quase idênticos. Embora a duplicação nem sempre seja o fim do mundo, isso não parece bom e aumenta o custo de manutenção.
- Como estamos duplicando a lógica em dois tipos, tivemos que duplicar os testes também.

Realmente queremos capturar a _idéia_ de uma pilha em um tipo, e ter um conjunto de testes para eles. Devemos usar nosso chapéu de refatoração agora, o que significa que não devemos mudar os testes porque queremos manter o mesmo comportamento.

Pré-genéricos, isso é o que _podemos_ fazer

```go
type PilhaDeInteiros = Pilha
type PilhaDeTextos = Pilha

type Pilha struct {
	valores []interface{}
}

func (p *Pilha) Empilhar(valor interface{}) {
	p.valores = append(p.valores, valor)
}

func (p *Pilha) EstaVazio() bool {
	return len(p.valores) == 0
}

func (p *Pilha) Desempilhar() (interface{}, bool) {
	if p.EstaVazio() {
		var zero interface{}
		return zero, false
	}

	indice := len(p.valores) - 1
	el := p.valores[indice]
	p.valores = p.valores[:indice]
	return el, true
}
```

- Estamos alterando nossas implementações anteriores de `PilhaDeInteiros` e `PilhaDeTextos` para um novo tipo unificado `Pilha`
- Removemos a segurança de tipo da `Pilha`, tornando-o de forma que os `valores` sejam uma [slice](https://github.com/larien/aprenda-go-com-testes/blob/main/primeiros-passos-com-go/arrays-e-slices/arrays-e-slices.md) da `interface{}`

... E nossos testes ainda passam. Quem precisa de genéricos?

### O problema de descartar o tipo de segurança

O primeiro problema é o mesmo que vimos com nosso `VerificaIgual` - perdemos a segurança de tipo. Agora posso 'Empurrar' as maçãs para uma pilha de laranjas.

Mesmo se tivermos a disciplina para não fazer isso, o código ainda é desagradável de trabalhar porque quando os métodos **retornam `interface{}` eles são horríveis de se trabalhar**.

Adicione o seguinte teste,

```go
minhaPilhaDeInteiros.Empilhar(1)
minhaPilhaDeInteiros.Empilhar(2)
primeiroNum, _ := minhaPilhaDeInteiros.Desempilhar()
segundoNum, _ := minhaPilhaDeInteiros.Desempilhar()
VerificaIgual(primeiroNum+segundoNum, 3)
```

Você obtém um erro do compilador, mostrando a fraqueza de perder a segurança de tipo:

```go
prog.go2:77:16: invalid operation: operator + not defined for primeiroNum (variable of type interface{})
```

Quando `Desempilhar` retorna uma `interface{}`, significa que o compilador não tem informações sobre o que são os dados e portanto, limita severamente o que podemos fazer. Ele não pode saber que deve ser um inteiro, então não nos permite usar o operador `+`.

Para contornar isso, o chamador deve fazer uma [asserção de tipo](https://golang.org/ref/spec#Type_assertions) para cada valor.

```go
minhaPilhaDeInteiros.Empilhar(1)
minhaPilhaDeInteiros.Empilhar(2)
primeiroNum, _ := minhaPilhaDeInteiros.Desempilhar()
segundoNum, _ := minhaPilhaDeInteiros.Desempilhar()

// obtenha inteiros da nossa interface{}
realmentePrimeiroNum, ok := primeiroNum.(int)
VerificaVerdadeiro(ok) // precisamos verificar definitivamente se obtivemos um inteiro da interface{}

realmenteSegundoNum, ok := segundoNum.(int)
VerificaVerdadeiro(ok) // e novamente!

VerificaIgual(realmentePrimeiroNum+realmenteSegundoNum, 3)
```

Just like you can define generic arguments to functions, you can define generic data structures.

Here's our new `Stack` implementation, featuring a generic data type and the tests, showing them working how we'd like them to work, with full type-safety. ([Full code listing here](https://go2goplay.golang.org/p/xAWcaMelgQV))

O desagrado que irradia deste teste seria repetido para cada usuário potencial de nossa implementação de `Pilha`, eca.

### Estruturas de dados genéricas para o resgate

Assim como você pode definir argumentos genéricos para funções, você pode definir estruturas de dados genéricas.

Aqui está nossa nova implementação `Pilha`, apresentando um tipo de dado genérico e os testes, mostrando-os funcionando como gostaríamos que funcionassem, com total segurança de tipo. ([Lista completa do código aqui](https://go2goplay.golang.org/p/g_miVq844Aq))

```go
package main

import (
    "log"
)

type Pilha[T any] struct {
    valores []T
}

func (p *Pilha[T]) Empilhar(valor T) {
    p.valores = append(p.valores, valor)
}

func (p *Pilha[T]) EstaVazio() bool {
    return len(p.valores)==0
}

func (p *Pilha[T]) Desempilhar() (T, bool) {
    if p.EstaVazio() {
        var zero T
        return zero, false
    }

    indice := len(p.valores) -1
    el := p.valores[indice]
    p.valores = p.valores[:indice]
    return el, true
}

func main() {
    minhaPilhaDeInteiros := new(Pilha[int])

    // verifica se a pilha está vazia
    VerificaVerdadeiro(minhaPilhaDeInteiros.EstaVazio())

    // adiciona alguma coisa e em seguida, verifica se a pilha não está vazia
    minhaPilhaDeInteiros.Empilhar(123)
    VerificaFalso(minhaPilhaDeInteiros.EstaVazio())

    // adiciona outra coisa e em seguida, desempilhe a pilha
    minhaPilhaDeInteiros.Empilhar(456)
    valor, _ := minhaPilhaDeInteiros.Desempilhar()
    VerificaIgual(valor, 456)
    valor, _ = minhaPilhaDeInteiros.Desempilhar()
    VerificaIgual(valor, 123)
    VerificaVerdadeiro(minhaPilhaDeInteiros.EstaVazio())

    // pode obter os números que colocamos como números, e não como interface {}
    minhaPilhaDeInteiros.Empilhar(1)
    minhaPilhaDeInteiros.Empilhar(2)
    primeiroNum, _ := minhaPilhaDeInteiros.Desempilhar()
    segundoNum, _ := minhaPilhaDeInteiros.Desempilhar()
    VerificaIgual(primeiroNum+segundoNum, 3)
}
```

Você notará que a sintaxe para definir estruturas de dados genéricas é consistente com a definição de argumentos genéricos para funções.

```go
type Pilha[T any] struct {
    valores []T
}
```

É _quase_ o mesmo que antes, mas o que estamos dizendo é que o **tipo de pilha restringe os tipos de valores com os quais você pode trabalhar**.

Depois de criar uma `Pilha[Laranja]` ou uma `Pilha[Maca]`, os métodos definidos em nossa pilha só permitirão que você passe e retornará apenas o tipo particular da pilha com a qual está trabalhando:

```go
func (p *Pilha[T]) Desempilhar() (T, bool) {
```

Você pode imaginar os tipos de implementação estão sendo gerados de alguma forma para você, dependendo do tipo de pilha que você criar:

```go
func (p *Pilha[Laranja]) Desempilhar() (Laranja, bool) {
```

```go
func (p *Pilha[Maca]) Desempilhar() (Maca, bool) {
```

Agora que fizemos essa refatoração, podemos remover com segurança o teste de pilha de strings porque não precisamos provar a mesma lógica repetidamente.

Usando um tipo de dados genérico, temos:

- Redução da duplicação de lógicas importantes.
- Faz o `Desempilhar` retornar `T` de forma que se criarmos uma `Pilha[int]` nós, na prática, obteremos de volta o `inteiro` do `Desempilhar`; agora podemos usar `+` sem a necessidade de ginástica de asserção de tipo.
- Evita o uso indevido em tempo de compilação. Você não pode `Empilhar` laranjas para uma pilha de maçã.

## Concluindo

Este capítulo deve ter lhe dado um gostinho da sintaxe dos genéricos e algumas idéias de por que os genéricos podem ser úteis. Escrevemos nossas próprias funções `Verifica`, que podemos reutilizar com segurança para experimentar outras idéias em torno dos genéricos, e implementamos uma estrutura de dados simples para armazenar qualquer tipo de dados que desejarmos, de maneira segura.

### Genéricos são mais simples do que usar `interface{}` na maioria dos casos

Se você não tem experiência com linguagens tipadas estaticamente, o ponto dos genéricos pode não ser imediatamente óbvio, mas espero que os exemplos neste capítulo tenham ilustrado onde a linguagem Go não é tão expressiva quanto gostaríamos. Em particular, usar `interface{}` torna seu código:

- Menos seguro (misturar maçãs e laranjas), requer mais tratamento de erros
- Menos expressivo, a `interface{}` não diz nada sobre os dados
- É mais provável confiar na [reflexão](https://github.com/larien/aprenda-go-com-testes/blob/main/primeiros-passos-com-go/reflection/reflection.md), asserções de tipo e etc, o que torna seu código mais difícil de trabalhar e mais sujeito a erros, pois empurra as verificações do tempo de compilação para o tempo de execução

Usar linguagens tipadas estaticamente é um ato de descrever restrições. Se você fizer isso bem, criará um código que não só é seguro e simples de usar, mas também mais simples de escrever porque o espaço de solução possível é menor.

Os genéricos nos fornecem uma nova maneira de expressar restrições em nosso código, o que, conforme demonstrado, nos permitirá consolidar e simplificar o código que não é possível fazer hoje.

### Os genéricos transformarão o Go em Java?

- Não.

Há muito [FUD (medo, incerteza e dúvida)](https://en.wikipedia.org/wiki/Fear,_uncertainty,_and_doubt) na comunidade Go sobre genéricos que levam a abstrações aterrorizantes e bases de código desconcertantes. Isso geralmente é advertido com "eles devem ser usados com cautela".

Embora isso seja verdade, não é um conselho especialmente útil porque isso é verdade para qualquer recurso de linguagem.

Poucas pessoas reclamam de nossa capacidade de definir interfaces que, como os genéricos, é uma forma de descrever as restrições em nosso código. Quando você descreve uma interface, está fazendo uma escolha de design que _pode ser pobre_, os genéricos não são os únicos em sua capacidade de tornar o código confuso e irritante.

### Você já está usando genéricos

Quando você considera que se você usou matrizes (arrays), slices ou mapas (maps); você _já foi um consumidor de código genérico_.

```go
var minhasMacas []Macas
// Você não pode fazer isso!
append(minhasMacas, Laranja{})
```

### Abstração não é um palavrão

É fácil mergulhar em [AbstractSingletonProxyFactoryBean](https://docs.spring.io/spring-framework/docs/current/javadoc-api/org/springframework/aop/framework/AbstractSingletonProxyFactoryBean.html), mas não vamos usar uma base de código sem nenhuma abstração também não é ruim. É sua função _conhecer_ conceitos relacionados quando apropriado, para que seu sistema seja mais fácil de entender e mudar; em vez de ser uma coleção de funções e tipos díspares com falta de clareza.

### [Faça funcionar, fazer certo, torná-lo rápido](https://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast#:~:text=%22Make%20it%20work%2C%20make%20it,to%20DesignForPerformance%20ahead%20of%20time.)

As pessoas têm problemas com os genéricos quando estão abstraindo muito rapidamente, sem informações suficientes para tomar boas decisões de design.

O ciclo TDD de red, green, refactor significa que você tem mais orientação sobre qual código você _realmente precisa_ para entregar seu comportamento, **em vez de imaginar abstrações antecipadamente**; mas você ainda precisa ter cuidado.

Não há regras rígidas e rápidas aqui, mas resista a tornar as coisas genéricas até que você possa ver que tem uma generalização útil. Quando criamos as várias implementações `Pilha`, começamos de maneira importante com o comportamento _concreto_ como `PilhaDeTextos` e `PilhaDeInteiros` apoiados por testes. A partir de nosso código _real_, podemos começar a ver padrões reais e, apoiados por nossos testes, podemos explorar a refatoração em direção a uma solução de propósito mais geral.

As pessoas geralmente aconselham você a generalizar apenas quando vir o mesmo código três vezes, o que parece uma boa regra inicial.

Um caminho comum que tomei em outras linguagens de programação foi:

- Um ciclo TDD para conduzir algum comportamento
- Outro ciclo TDD para exercitar alguns outros cenários relacionados

> Hmm, essas coisas parecem semelhantes - mas um pouco de duplicação é melhor do que o acoplamento a uma abstração ruim

- Pense mais sobre isso
- Outro ciclo TDD

> OK, gostaria de tentar ver se posso generalizar isso. Graças a Deus eu sou uma pessoa inteligente e uma pessoa bonita porque uso TDD, então posso refatorar sempre que quiser, e o processo me ajudou a entender qual comportamento eu realmente preciso antes de projetar muito.

- Essa abstração é agradável! Os testes ainda estão passando e o código é mais simples
- Agora posso excluir uma série de testes, capturei a _essência_ do comportamento e removi detalhes desnecessários
