# Mocks

[**Você pode encontrar todos os códigos para esse capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/master/primeiros-passos-com-go/mocks)

Te pediram para criar um programa que conta a partir de 3, imprimindo cada número em uma linha nova (com um segundo de intervalo entre cada uma) e quando chega a zero, imprime "Vai!" e sai.

```text
3
2
1
Vai!
```

Vamos resolver isso escrevendo uma função chamada `Contagem` que vamos colocar dentro de um programa `main` e se parecer com algo assim:

```go
package main

func main() {
    Contagem()
}
```

Apesar de ser um programa simples, para testá-lo completamente vamos precisar, como de costume, de uma abordagem _iterativa_ e _orientada a testes_.

Mas o que quero dizer com iterativa? Precisamos ter certeza de que tomamos os menores passos que pudermos para ter um _software_ útil.

Não queremos passar muito tempo com código que vai funcionar hora ou outra após alguma implementação mirabolante, porque é assim que os desenvolvedores caem em armadilhas. **É importante ser capaz de dividir os requerimentos da menor forma que conseguir para você ter um** _**software funcionando**_**.**

Podemos separar essa tarefa da seguinte forma:

-   Imprimir 3
-   Imprimir de 3 para Vai!
-   Esperar um segundo entre cada linha

## Escreva o teste primeiro

Nosso software precisa imprimir para a saída. Vimos como podemos usar a injeção de dependência para facilitar nosso teste na [seção anterior](../injecao-de-dependencia/injecao-de-dependencia.md).

```go
func TestContagem(t *testing.T) {
    buffer := &bytes.Buffer{}

    Contagem(buffer)

    resultado := buffer.String()
    esperado := "3"

    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }
}
```

Se tiver dúvidas sobre o `buffer`, leia a [seção anterior](../injecao-de-dependencia/injecao-de-dependencia.md) novamente.

Sabemos que nossa função `Contagem` precisa escrever dados em algum lugar e o `io.Writer` é a forma de capturarmos essa saída como uma interface em Go.

-   Na `main`, vamos enviar o `os.Stdout` como parâmetro para nossos usuários verem a contagem regressiva impressa no terminal.
-   No teste, vamos enviar o `bytes.Buffer` como parâmetro para que nossos testes possam capturar que dado está sendo gerado.

## Execute o teste

`./contagem_test.go:11:2: undefined: Contagem`

`indefinido: Contagem`

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Defina `Contagem`:

```go
func Contagem() {}
```

Tente novamente:

```go
./contagem_test.go:11:11: too many arguments in call to Countdown
    have (*bytes.Buffer)
    want ()
```

`argumentos demais na chamada para Contagem`

O compilador está te dizendo como a assinatura da função deve ser, então é só atualizá-la.

```go
func Contagem(saida *bytes.Buffer) {}
```

`contagem_test.go:17: resultado '', esperado '3'`

Perfeito!

## Escreva código o suficiente para fazer o teste passar

```go
func Contagem(saida *bytes.Buffer) {
    fmt.Fprint(saida, "3")
}
```

Estamos usando `fmt.Fprint`, o que significa que ele recebe um `io.Writer` (como `*bytes.Buffer`) e envia uma `string` para ele. O teste deve passar.

## Refatoração

Agora sabemos que, apesar do `*bytes.Buffer` funcionar, seria melhor ter uma interface de propósito geral ao invés disso.

```go
func Contagem(saida io.Writer) {
    fmt.Fprint(saida, "3")
}
```

Execute os testes novamente e eles devem passar.

Só para finalizar, vamos colocar nossa função dentro da `main` para que possamos executar o software para nos assegurarmos de que estamos progredindo.

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func Contagem(saida io.Writer) {
	fmt.Fprint(saida, "3")
}

func main() {
	Contagem(os.Stdout)
}
```

Execute o programa e surpreenda-se com seu trabalho.

Apesar de parecer simples, essa é a abordagem que recomendo para qualquer projeto. **Escolher uma pequena parte da funcionalidade e fazê-la funcionar do começo ao fim com apoio de testes.**

Depois, precisamos fazer o software imprimir 2, 1 e então "Vai!".

## Escreva o teste primeiro

Após investirmos tempo e esforço para fazer o principal funcionar, podemos iterar nossa solução com segurança e de forma simples. Não vamos mais precisar parar e executar o programa novamente para ter confiança de que ele está funcionando, desde que a lógica esteja testada.

```go
func TestContagem(t *testing.T) {
    buffer := &bytes.Buffer{}

    Contagem(buffer)

    resultado := buffer.String()
    esperado := `3
2
1
Vai!`
    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }
}
```

A sintaxe de aspas simples é outra forma de criar uma `string`, mas te permite colocar coisas como linhas novas, o que é perfeito para nosso teste.

## Execute o teste

```bash
contagem_test.go:21: resultado '3', esperado '3
        2
        1
        Vai!'
```

## Escreva código o suficiente para fazer o teste passar

```go
func Contagem(saida io.Writer) {
    for i := 3; i > 0; i-- {
        fmt.Fprintln(saida, i)
    }
    fmt.Fprint(saida, "Go!")
}
```

Usamos um laço `for` fazendo contagem regressiva com `i--` e depois `fmt.Fprintln` para imprimir a `saida` com nosso número seguro por um caracter de nova linha. Finalmente, usamos o `fmt.Fprint` para enviar "Vai!" no final.

## Refatoração

Não há muito para refatorar além de transformar alguns valores mágicos em constantes com nomes descritivos.

```go
const ultimaPalavra = "Go!"
const inicioContagem = 3

func Contagem(saida io.Writer) {
    for i := inicioContagem; i > 0; i-- {
        fmt.Fprintln(saida, i)
    }
    fmt.Fprint(saida, ultimaPalavra)
}
```

Se executar o programa agora, você deve obter a saída desejada, mas não tem uma contagem regressiva dramática com as pausas de 1 segundo.

Go te permite obter isso com `time.Sleep`. Tente adicionar essa função ao seu código.

```go
func Contagem(saida io.Writer) {
    for i := inicioContagem; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(saida, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(saida, ultimaPalavra)
}
```

Se você executar o programa, ele funciona conforme esperado.

## Mock

Os testes ainda vão passar e o software funciona como planejado, mas temos alguns problemas:

-   Nossos testes levam 4 segundos para rodar.
    -   Todo conteúdo gerado sobre desenvolvimento de software enfatiza a importância de loops de feedback rápidos.
    -   **Testes lentos arruinam a produtividade do desenvolvedor**.
    -   Imagine se os requerimentos ficam mais sofisticados, gerando a necessidade de mais testes. É viável adicionar 4s para cada teste novo de `Contagem`?
-   Não testamos uma propriedade importante da nossa função.

Temos uma dependência no `Sleep` que precisamos extrair para podermos controlá-la nos nossos testes.

Se conseguirmos _mockar_ o `time.Sleep`, podemos usar a _injeção de dependências_ para usá-lo ao invés de um `time.Sleep` "de verdade", e então podemos **verificar as chamadas** para certificar de que estão corretas.

## Escreva o teste primeiro

Vamos definir nossa dependência como uma interface. Isso nos permite usar um Sleeper _de verdade_ em `main` e um _sleeper spy_ nos nossos testes. Usar uma interface na nossa função `Contagem` é essencial para isso e dá certa flexibilidade à função que a chamar.

```go
type Sleeper interface {
    Sleep()
}
```

Tomei uma decisão de design que nossa função `Contagem` não seria responsável por quanto tempo o sleep leva. Isso simplifica um pouco nosso código, pelo menos por enquanto, e significa que um usuário da nossa função pode configurar a duração desse tempo como preferir.

Agora precisamos criar um _mock_ disso para usarmos nos nossos testes.

```go
type SleeperSpy struct {
    Chamadas int
}

func (s *SleeperSpy) Sleep() {
    s.Chamadas++
}
```

_Spies_ (espiões) são um tipo de _mock_ em que podemos gravar como uma dependência é usada. Eles podem gravar os argumentos definidos, quantas vezes são usados etc. No nosso caso, vamos manter o controle de quantas vezes `Sleep()` é chamada para verificá-la no nosso teste.

Atualize os testes para injetar uma dependência no nosso Espião e verifique se o sleep foi chamado 4 vezes.

```go
func TestContagem(t *testing.T) {
    buffer := &bytes.Buffer{}
    sleeperSpy := &SleeperSpy{}

    Contagem(buffer, sleeperSpy)

    resultado := buffer.String()
    esperado := `3
2
1
Vai!`

    if resultado != esperado {
        t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
    }

    if sleeperSpy.Chamadas != 4 {
        t.Errorf("não houve chamadas suficientes do sleeper, esperado 4, resultado %d", sleeperSpy.Chamadas)
    }
}
```

## Execute o teste

```bash
too many arguments in call to Contagem
    have (*bytes.Buffer, *SpySleeper)
    want (io.Writer)
```

## Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

Precisamos atualizar a `Contagem` para aceitar nosso `Sleeper`:

```go
func Contagem(saida io.Writer, sleeper Sleeper) {
    for i := inicioContagem; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(saida, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(saida, ultimaPalavra)
}
```

Se tentar novamente, nossa `main` não vai mais compilar pelo mesmo motivo:

```text
./main.go:26:11: not enough arguments in call to Contagem
    have (*os.File)
    want (io.Writer, Sleeper)
```

Vamos criar um sleeper _de verdade_ que implementa a interface que precisamos:

```go
type SleeperPadrao struct {}

func (d *SleeperPadrao) Sleep() {
	time.Sleep(1 * time.Second)
}
```

Podemos usá-lo na nossa aplicação real, como:

```go
func main() {
    sleeper := &SleeperPadrao{}
    Contagem(os.Stdout, sleeper)
}
```

## Escreva código o suficiente para fazer o teste passar

Agora o teste está compilando, mas não passando. Isso acontece porque ainda estamos chamando o `time.Sleep` ao invés da injetada. Vamos arrumar isso.

```go
func Contagem(saida io.Writer, sleeper Sleeper) {
    for i := inicioContagem; i > 0; i-- {
        sleeper.Sleep()
        fmt.Fprintln(saida, i)
    }

    sleeper.Sleep()
    fmt.Fprint(saida, ultimaPalavra)
}
```

O teste deve passar sem levar 4 segundos.

### Ainda temos alguns problemas

Ainda há outra propriedade importante que não estamos testando.

A `Contagem` deve ter uma pausa para cada impressão, como por exemplo:

-   `Pausa`
-   `Imprime N`
-   `Pausa`
-   `Imprime N-1`
-   `Pausa`
-   `Imprime Vai!`
-   etc

Nossa alteração mais recente só verifica se o software teve 4 pausas, mas essas pausas poderiam ocorrer fora de ordem.

Quando escrevemos testes, se não estiver confiante de que seus testes estão te dando confiança o suficiente, quebre-o (mas certifique-se de que você salvou suas alterações antes)! Mude o código para o seguinte:

```go
func Contagem(saida io.Writer, sleeper Sleeper) {
    for i := inicioContagem; i > 0; i-- {
        sleeper.Pausa()
        fmt.Fprintln(saida, i)
    }

    for i := inicioContagem; i > 0; i-- {
        fmt.Fprintln(saida, i)
    }

    sleeper.Pausa()
    fmt.Fprint(saida, ultimaPalavra)
}
```

Se executar seus testes, eles ainda vão passar, apesar da implementação estar errada.

Vamos usar o spy novamente com um novo teste para verificar se a ordem das operações está correta.

Temos duas dependências diferentes e queremos gravar todas as operações delas em uma lista. Logo, vamos criar _um spy para ambas_.

```go
type SpyContagemOperacoes struct {
    Chamadas []string
}

func (s *SpyContagemOperacoes) Pausa() {
    s.Chamadas = append(s.Chamadas, pausa)
}

func (s *SpyContagemOperacoes) Write(p []byte) (n int, err error) {
    s.Chamadas = append(s.Chamadas, escrita)
    return
}

const escrita = "escrita"
const pausa = "pausa"
```

Nosso `SpyContagemOperacoes` implementa tanto o `io.Writer` quanto o `Sleeper`, gravando cada chamada em um slice. Nesse teste, temos preocupação apenas na ordem das operações, então apenas gravá-las em uma lista de operações nomeadas é suficiente.

Agora podemos adicionar um subteste no nosso conjunto de testes.

```go
t.Run("pausa antes de cada impressão", func(t *testing.T) {
        spyImpressoraSleep := &SpyContagemOperacoes{}
        Contagem(spyImpressoraSleep, spyImpressoraSleep)

        esperado := []string{
            pausa,
            escrita,
            pausa,
            escrita,
            pausa,
            escrita,
            pausa,
            escrita,
        }

        if !reflect.DeepEqual(esperado, spyImpressoraSleep.Chamadas) {
            t.Errorf("esperado %v chamadas, resultado %v", esperado, spyImpressoraSleep.Chamadas)
        }
    })
```

Esse teste deve falhar. Volte o código que quebramos para a versão correta e agora o novo teste deve passar.

Agora temos dois spies no `Sleeper`. O próximo passo é refatorar nosso teste para que um teste o que está sendo impresso e o outro se certifique de que estamos pausando entre as impressões. Por fim, podemos apagar nosso primeiro spy, já que não é mais utilizado.

```go
func TestContagem(t *testing.T) {

    t.Run("imprime 3 até Vai!", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        Contagem(buffer, &SpyContagemOperacoes{})

        resultado := buffer.String()
        esperado := `3
2
1
Vai!`

        if resultado != esperado {
            t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
        }
    })

    t.Run("pausa antes de cada impressão", func(t *testing.T) {
        spyImpressoraSleep := &SpyContagemOperacoes{}
        Contagem(spyImpressoraSleep, spyImpressoraSleep)

        esperado := []string{
            pausa,
            escrita,
            pausa,
            escrita,
            pausa,
            escrita,
            pausa,
            escrita,
        }

        if !reflect.DeepEqual(esperado, spyImpressoraSleep.Chamadas) {
            t.Errorf("esperado %v chamadas, resultado %v", esperado, spyImpressoraSleep.Chamadas)
        }
    })
}
```

Agora temos nossa função e suas duas propriedades testadas adequadamente.

## Extendendo o Sleeper para se tornar configurável

Uma funcionalidade legal seria o `Sleeper` ser configurável.

### Escreva o teste primeiro

Agora vamos criar um novo tipo para `SleeperConfiguravel` que aceita o que precisamos para configuração e teste.

```go
type SleeperConfiguravel struct {
	duracao time.Duration
	pausa   func(time.Duration)
}
```

Estamos usando a `duracao` para configurar o tempo de pausa e `pausa` como forma de passar uma função de pausa. A assinatura de `sleep` é a mesma de `time.Sleep`, nos permitindo usar `time.Sleep` na nossa implementação real e um spy nos nossos testes.

```go
type TempoSpy struct {
	duracaoPausa time.Duration
}

func (t *TempoSpy) Pausa(duracao time.Duration) {
	t.duracaoPausa = duracao
}
```

Definindo nosso spy, podemos criar um novo teste para o sleeper configurável.

```go
func TestSleeperConfiguravel(t *testing.T) {
    tempoPausa := 5 * time.Second

    tempoSpy := &TempoSpy{}
    sleeper := SleeperConfiguravel{tempoPausa, tempoSpy.Pausa}
    sleeper.Pausa()

    if tempoSpy.duracaoPausa != tempoPausa {
        t.Errorf("deveria ter pausado por %v, mas pausou por %v", tempoPausa, tempoSpy.duracaoPausa)
    }
}
```

Não há nada de novo nesse teste e seu funcionamento é bem semelhante aos testes com mock anteriores.

### Execute o teste

```bash
sleeper.Pausa undefined (type SleeperConfiguravel has no field or method Pausa, but does have pausa)
```

`sleeper.Pausa não definido (tipo SleeperConfiguravel não tem campo ou método Pausa, mas tem o método sleep`

Você deve ver uma mensagem de erro bem clara indicando que não temos um método `Pausa` criado no nosso `SleeperConfiguravel`.

### Escreva o mínimo de código possível para fazer o teste rodar e verifique a saída do teste que tiver falhado

```go
func (c *SleeperConfiguravel) Pausa() {
}
```

Com nossa nova função `Pausa` implementada, ainda há um teste falhando.

```bash
contagem_test.go:56: deveria ter pausado por 5s, mas pausou por 0s
```

### Escreva código o suficiente para fazer o teste passar

Tudo o que precisamos fazer agora é implementar a função `Pausa` para o `SleeperConfiguravel`.

```go
func (s *SleeperConfiguravel) Pausa() {
    s.pausa(s.duracao)
}
```

Com essa mudança, todos os testes devem voltar a passar.

### Limpeza e refatoração

A última coisa que precisamos fazer é de fato usar nosso `SleeperConfiguravel` na função main.

```go
func main() {
    sleeper := &SleeperConfiguravel{1 * time.Second, time.Sleep}
    Contagem(os.Stdout, sleeper)
}
```

Se executarmos os testes e o programa manualmente, podemos ver que todo o comportamento permanece o mesmo.

Já que estamos usando o `SleeperConfiguravel`, é seguro deletar o `SleeperPadrao`.

## Mas o mock não é do demonho?

Você já deve ter ouvido que o mock é do mal. Quase qualquer coisa no desenvolvimento de software pode ser usada para o mal, assim como o [DRY](https://pt.wikipedia.org/wiki/Don%27t_repeat_yourself).

As pessoas acabam chegando numa fase ruim em que não _dão atenção aos próprios testes_ e _não respeitam a etapa de refatoração_.

Se seu código de mock estiver ficando complicado ou você tem que mockar muita coisa para testar algo, você deve _prestar mais atenção_ a essa sensação ruim e pensar sobre o seu código. Geralmente isso é sinal de que:

-   A coisa que você está testando está tendo que fazer coisas demais
    -   Modularize a função para que faça menos coisas
-   Suas dependências estão muito desacopladas
    -   Pense e uma forma de consolidar algumas das dependências em um módulo útil
-   Você está se preocupando demais com detalhes de implementação
    -   Dê prioridade em testar o comportamento esperado ao invés da implementação

Normalmente, muitos pontos de mock são sinais de _abstração ruim_ no seu código.

**As pessoas costumam pensar que essa é uma fraqueza no TDD, mas na verdade é um ponto forte**. Testes mal desenvolvidos são resultado de código ruim. Código bem desenvolvido é fácil de ser testado.

### Só que mocks e testes ainda estão dificultando minha vida!

Já se deparou com a situação a seguir?

-   Você quer refatorar algo
-   Para isso, você precisa mudar vários testes
-   Você duvida do TDD e cria um post no Medium chamado "Mock é prejudicial"

Isso costuma ser um sinal de que você está testando muito _detalhe de implementação_. Tente fazer de forma que esteja testando _comportamentos úteis_, a não ser que a implementação seja tão importante que a falta dela possa fazer o sistema quebrar.

Às vezes é difícil saber _qual nível_ testar exatamente, então aqui vai algumas ideias e regras que tento seguir:

-**A definição de refatoração é que o código muda, mas o comportamento permanece o mesmo**. Se você decidiu refatorar alguma coisa, na teoria você deve ser capaz de salvar seu código sem que o teste mude. Então, quando estiver escrevendo um teste, pergunte para si: - Estou testando o comportamento que quero ou detalhes de implementação? - Se fosse refatorar esse código, eu teria que fazer muitas mudanças no meu teste?

-   Apesar do Go te deixar testar funções privadas, eu evitaria fazer isso, já que funções privadas costumam ser detalhes de implementação.
-   Se o teste estiver com **3 mocks, esse é um sinal de alerta** - hora de repensar no design.
-   Use spies com cuidado. Spies te deixam ver a parte interna do algoritmo que você está escrevendo, o que pode ser bem útil, mas significa que há um acoplamento maior entre o código do teste e a implementação. **Certifique-se de que você realmente precisa desses detalhes se você vai colocar um spy neles**.

Como sempre, regras no desenvolvimento de software não são realmente regras e podem haver exceções. [O artigo do Uncle Bob sobre "Quando mockar"](https://8thlight.com/blog/uncle-bob/2014/05/10/WhenToMock.html) (em inglês) tem alguns pontos excelentes.

## Resumo

### Mais sobre abordagem TDD

-   Quando se deparar com exemplos menos comuns, divida o problema em "linhas verticais finas". Tente chegar em um ponto onde você tem _software em funcionamento com o apoio de testes_ o mais rápido possível, para evitar cair em armadilhas e se perder.
-   Quando tiver uma parte do software em funcionamento, deve ser mais fácil _iterar com etapas pequenas_ até chegar no software que você precisa.

> "Quando usar o desenvolvimento iterativo? Apenas em projetos que você quer obter sucesso."

Martin Fowler.

### Mock

-   **Sem o mock, partes importantes do seu código não serão testadas**. No nosso caso, não seríamos capazes de testar se nosso código pausava em cada impressão, mas existem inúmeros exemplos. Chamar um serviço que _pode_ falhar? Querer testar seu sistema em um estado em particular? É bem difícil testar esses casos sem mock.
-   Sem mocks você pode ter que definir bancos de dados e outras dependências externas só para testar regras de negócio simples. Seus testes provavelmente ficarão mais lentos, resultando em **loops de feedback lentos**.
-   Ter que se conectar a um banco de dados ou webservice para testar algo vai tornar seus testes **frágeis** por causa da falta de segurança nesses serviços.

Uma vez que a pessoa aprende a mockar, é bem fácil testar pontos demais de um sistema em termos da _forma que ele funciona_ ao invés _do que ele faz_. Sempre tenha em mente o **valor dos seus testes** e qual impacto eles teriam em uma refatoração futura.

Nesse artigo sobre mock, falamos sobre **spies**, que são um tipo de mock. Aqui estão diferentes tipos de mocks. [O Uncle Bob explica os tipos em um artigo bem fácil de ler](https://8thlight.com/blog/uncle-bob/2014/05/14/TheLittleMocker.html) (em inglês). Nos próximos capítulos, vamos precisar escrever código que depende de outros para obter dados, que é aonde vou mostrar os **Stubs** em ação.
