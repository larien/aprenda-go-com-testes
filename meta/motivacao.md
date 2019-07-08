# Por que criar testes unitários e como fazê-los dar certo

[Vejam um vídeo meu falando sobre esse assunto](https://www.youtube.com/watch?v=Kwtit8ZEK7U)

Se não gostar muito de vídeos, aqui vai o artigo relacionado a isso.

## Software

A promessa do software é que ele pode mudar. É por isso que é chamado de \_soft_ware: é mais maleável se comparado ao hardware. Uma boa equipe de engenharia deve ser um componente incrível para uma empresa, criando sistemas que podem evoluir com um negócio para manter seu valor de entrega.

Então por que somos tão ruins nisso? Quantos projetos que você ouve falar sobre que ultrapassam o nível da falha? Ou viram "legado" e precisam ser totalmente recriados (e a reescrita também acaba falhando)!

Mas como é que um software "falha"? Não dá para ele apenas ser modificado até estar correto? É isso que prometemos!

Muita gente costuma escolher o Go para criar sistemas porque a linguagem teve várias decisões que evitam que o software vire legado.

-   Comparado à minha antiga vida de Scala onde [descrevi como é fácil acabar se dando mal com a linguagem](http://www.quii.co.uk/Scala_-_Just_enough_rope_to_hang_yourself), o Go tem apenas 25 palavras-chave. _Muitos_ sistemas podem ser criados a partir da biblioteca padrão e alguns outros pacotes pequenos. O que se espera é que com Go você possa escrever código, voltar a vê-lo 6 meses depois e ele ainda fazer sentido.
-   As ferramentas relacionadas a testes, benchmarking, linting e shipping são incríveis se comparadas à maioria das alternativas.
-   A biblioteca padrão é brilhante.
-   Velocidade de compilação muito rápida para loops de feedback mais frequentes.
-   A famigerada promessa da compatibilidade. Parece que Go vai receber `generics` e outras funcionalidades no futuro, mas os mantenedores prometeram que mesmo o código Go que você escreveram cinco anos atrás ainda vai compilar e funcionar. Eu literalmente passei semanas atualizando um projeto em Scala da versão 2.8 para a 2.10.

Com todas essas propriedades ótimas, ainda podemos acabar criando sistemas terríveis. Por isso, precisamos aplicar lições de engenharia de software que se aplicam independente do quão maravilhosa (ou não) sua linguagem seja.

Em 1974, um engenheiro de software esperto chamado [Manny Lehman](https://pt.wikipedia.org/wiki/Meir_M._Lehman) escreveu as [leis de Lehman para a evolução do software](https://www.baguete.com.br/colunas/jorge-horacio-audy/17/04/2014/as-8-leis-de-lehman-foram-o-manifesto-do-seculo-xx).

> As leis descrevem um equilíbrio entre o desenvolvimento de software em uma ponta e a diminuição do progresso em outra.

É importante entender esses extremos para não acabar em um ciclo infinito de entregar sistemas que se tornam em legado e precisam ser reescritos novamente.

## Lei da Mudança Contínua

> Qualquer software utilizado no mundo real precisa se adaptar ou vai se tornar cada vez mais obsoleto.

Parece óbvio que um software _precisa_ mudar ou acaba se tornando menos útil, mas quantas vezes isso é ignorado?

Muitas equipes são incentivadas a entregar um projeto em uma data específica e passar para o próximo projeto. Se o software tiver "sorte", vai acabar na mão de outro grupo de pessoas para mantê-lo, mas é claro que nenhuma dessas pessoas o escreveu.

As pessoas se preocupam em escolher uma framework que vai ajudá-las a "entregar rapidamente", mas não focam na longevidade do sistema em termos de como precisa ser evoluído.

Mesmo se você for um engenheiro de software incrível, ainda vai cair na armadilha de não saber que futuro aguarda seu software. Já que o negócio muda, o código brilhante que você escreveu já não vai mais ser relevante.

Lehman estava contudo nos anos 70, porque nos deu outra lei para quebrarmos a cabeça.

## Lei da Complexidade Crescente

> Enquanto o softwate evolui, sua complexidade aumenta. A não ser que um esforço seja investido para reduzi-la.

O que ele diz aqui é que não podemos ter equipes de software para funcionar apenas como fábricas de funcionalidades, inserindo mais e mais funcionalidades no software para que ele possa sobreviver a longo prazo.

Nós **temos** que lidar com a complexidade do sistema conforme o conhecimento do nosso domínio muda.

## Refatoração

Existem _diversas_ facetas na engenharia de software que mantêm um software maleável, como:

-   Capacitação do desenvolvimento
-   Em termos gerais, código "bom". Separação sensível de responsabilidades, etc
-   Habilidades de comunicação
-   Arquitetura
-   Observabilidade
-   Implantabilidade
-   Testes automatizados
-   Retornos de feedback

Vou focar na refatoração. Quantas vezes você já ouviu a frase "precisamos refatorar isso"? Provavelmente dita para uma pessoa desenvolvedora em seu primeiro dia de programação sem pensar duas vezes.

De onde essa frase vem? Por que refatorar é diferente de escrever código?

Sei que eu e muitas outras pessoas só _pensaram_ que estavam refatorando, mas estávamos cometendo um erro.

[Martin Fowler descreve como as pessoas entendem a refatoração errada aqui.](https://martinfowler.com/bliki/RefactoringMalapropism.html)

> No entanto, o termo "refatoração" costuma ser utilizado de forma inapropriada. Se alguém fala que um sistema ficará quebrado por alguns dais enquanto está sendo refatorado, pode ter certeza que eles não estão refatorando.

Então o que é refatoração?

### Fatoração

Quando estudava matemática na escola, você provavelmente aprendeu fatoração. Aqui vai um exemplo bem simples:

-   Calcule `1/2 + 1/4`

Para fazer isso você `fatora` os denominadores (você também pode conhecer como MMC, mínimo múltiplo comum), transformando a expressão em `2/4 + 1/4` que então pode se transformar em `3/4`.

Podemos tirar algumas lições importantes disso. Quando `fatoramos a expressão`, `não mudamos o que ela faz`. Ambas as expreções são iguais a `3/4`, mas facilitamos a forma como trabalhamos com esse resultado; trocar `1/2` por `2/4` torna nosso "domínio" mais fácil.

Quando refatora seu código, você tenta encontrar formar de tornar seu código mais fácil de entender e "encaixar" no seu entendimento atual do que o sistema precisa fazer. Mas é extremamente importante que **o comportamento do código não seja alterado**.

#### Exemplo em Go

Aqui está uma função que cumprimenta `nome` em uma `linguagem` específica:

```go
func Ola(nome, linguagem string) string {

  if linguagem == "br" {
     return "Olá, " + nome
  }

  if linguagem == "fr" {
     return "Bonjour, " + nome
  }

  // e mais várias linguagens

  return "Hello, " + nome
}
```

Não é bom ter várias condicionais `if` e temos uma duplicação que concatena um cumprimento específico da linguagem com `,` e o `nome`. Logo, vou refatorar o código.

```go
func Ola(nome, linguagem string) string {
      return fmt.Sprintf(
          "%s, %s",
          cumprimento(linguagem),
          nome,
      )
}

var cumprimentos = map[string]string {
  br: "Olá",
  fr: "Bonjour",
  // etc..
}

func cumprimento(linguagem string) string {
  cumprimento, existe := cumprimentos[linguagem]

  if existe {
     return cumprimento
  }

  return "Hello"
}
```

A natureza dessa refatoração não é tão importante. O que importa é que não mudei o comportamento do código.

Quando estiver refatorando, você pode fazer o que quiser: adicionar interfaces, tipos novos, funções, métodos etc. A única regra é que você não mude o comportamento do software.

### Quando estiver refatorando o código, seu comportamento não deve ser modificado

Isso é muito importante. Se estiver mudando o comportamento enquanto refatora, você vai estar fazendo _duas_ coisas de uma vez. Como engenheiros de software, aprendemos a dividir o sistema em diferentes arquivos/pacotes/funções/etc porque sabemos que tentar entender algo enorme e acoplado é difícil.

Não queremos ter que pensar sobre muitas coisas ao mesmo tempo porque é aí que cometemos erros. Já vi tantos esforços de refatoração falharem pelas pessoas que estavam desenvolvendo darem um passo maior que a perna.

Quando fazia fatorações nas aulas de matemática com papel e caneta, eu precisava verificar manualmente que não havia mudado o significado das expressões na minha cabeça. Como sabemos que não estamos mudando o comportamento quando refatoramos as coisas no código, especialmente em um sistema que não é tão simples?

As pessoas que escolhem não escrever testes vão depender do teste manual. Para quem não trabalha em um projeto pequeno, isso vai ser uma tremenda perda de tempo e não vai escalar a longo prazo.

**Para ter uma refatoração segura, você precisa escrever testes unitários**, porque eles te dão:

-   Confiança de que você pode mudar o código sem se preocupar com mudar seu comportamento
-   Documentação para humanos sobre como o sistema deve se comportar
-   Feedback mais rápido e confiável que o teste manual

#### Exemplo em Go

Um teste unitário para a nossa função `Ola` pode ser feito assim:

```go
func TestOla(t *testing.T) {
  obtido := Ola(“Chris”, br)
  esperado := "Olá, Chris"

  if obtido != esperado {
     t.Errorf("obtido '%s' esperado '%s'", obtido, esperado)
  }
}
```

Na linha de comando, posso executar `go test` e obter feedback imediato se minha refatoração alterou o comportamento da função. Na prática, é melhor aprender aonde fica o botão mágico que vai executar seus testes dentro do seu editor/IDE (ou rodar os testes sempre que salvar o arquivo).

Você deve entrar em uma rotina em que acaba fazendo:

-   Refatorar uma parte pequena
-   Executar testes
-   Repetir

Tudo dentro de um ciclo de feedback contínuo para que você não caia em uma cilada e cometa erros.

Ter um projeto onde os seus principais comportamentos são testados unicamente e te dão feedback em menos de um segundo traz uma relação forte de segurança para refatorar sempre que for necessário. Isso nos ajuda a gerenciar a complexidade crescente que Lehman descreve.

## Se testes unitários são tão bons, por que há resistência em escrevê-los?

De um lado, é possível ver pessoas (como eu) dizendo que testes unitários são importantes para a saúda do seu sistema a longo prazo, porque eles certificam que você possa continuar refatorando com confiança.

Do outro lado, é possível ver pessoas descrevendo experiências com testes unitários que na verdade _dificultaram_ a refatoração.

Se pergunte o seguinte: com qual frequência você precisa mudar seus testes quando refatora? Estive em diversos projetos com boa cobertura de testes e mesmo assim os engenheiros estavam relutantes em refatorar por causa do esforço perceptível de alterar testes.

Esse é o oposto do que prometemos!

### Por que isso acontece?

Imagine que te pediram para desenvolver um quadrado e você chegou à conclusão que seria necessário unir dois triângulos.

![Dois triângulos retângulos formando um quadrado](https://i.imgur.com/ela7SVf.jpg)

Escrevemos nossos testes unitários nos baseando no nosso quadrado para ter certeza de que os lados são iguais e depois escrevemos alguns testes em relação aos nossos triângulos. Queremos ter certeza de que nossos triângulos são renderizados corretamente, então afirmamos que os ângulos somados dos triângulos dão 180 graus, ou verificamos que os dois são criados, etc etc. A cobertura de testes é muito importante e escrever esses testes é bem fácil, então por que não?

Algumas semanas depois, a Lei da Mudança Contínua bate no seu sistema e uma nova pessoa desenvolvedora faz algumas mudanças. Ela acredita que seria melhor se os quadrados fossem formados por dois retângulos ao invés dos dois triângulos.

![Dois retângulos formando um quadrado](https://i.imgur.com/1G6rYqD.jpg)

Ela tenta fazer essa refatoração e percebe que alguns testes falharam. Ela quebrou algum comportamento realmente importante aqui? Agora ela tem que investigar esses testes de triângulo e entender o que está acontecendo.

_Na verdade, não é tão importante que o quadrado seja formado por triângulo_, mas **nossos testes fizeram com que isso parecesse mais importante do que deveria em relação aos detalhes da nossa implementação**

## Favorecer o comportamento do teste ao invés do detalhe da implementação

Quando ouço pessoas reclamando sobre testes unitários, frequentemente o motivo é que eles estão em um nível errado de abstração. Eles testam detalhes da implementação, testando coisas muito específicas ou fazendo muitos mocks.

Acredito que isso deriva de uma falta de entendimento do que testes unitários são e perseguem métricas vaidosas (cobertura de testes).

Se estou apenas testando o comportamento, não deveríamos apenas escrever testes de sistema/caixa preta? Esses tipos de testes geram muito valor em termos de verificar as principais jornadas do usuário, mas costumam ser difíceis de escrever e lentos para rodar. Por esse motivos, eles não são muito úteis para a _refatoração_ porque o ciclo de feedback é lento. Além disso, os testes de caixa preta tendem a não te ajudar muito com as causas de origem comparados aos testes unitários.

Logo, _qual_ é o nível de abstração correto?

## Escrevendo testes unitários de forma efetiva é um problema de design

Deixando testes de lado por um momento, é desejável "unidades" independentes e desacopladas dentro do seu sistema, centradas em torno de conceitos essenciais do seu domínio.

Gosto de imaginar essas unidades tão simples quanto blocos de Lego que têm APIs coerentes e que eu possa combinar com outros blocos para criar sistemas maiores. Por baixo dessas APIs pode haver várisas coisas (tipos, funções etc) colaborando para fazê-las funcionar conforme esperado.

Por exemplo: se estiver escrevendo um banco em Go, você deve ter um pacote "conta". Ele vai te apresentar uma API que não vaza detalhes da implementação e é fácil de ser integrado.

Se tiver essas unidades que seguem essas propriedades, você consegue escrever testes unitários para suas APIs públicas. _Por definição_, esses testes só podem testar os comportamentos importantes. Por baixo dos panos dessas unidades, fico livre para refatorar a implementação o quanto eu precisar e os testes para a maior parte dela não deve me atrapalhar.

### Mas são testes unitários, mesmo?

**SIM**. Testes unitários são feitos para "unidades", como já descrevi. Eles _nunca_ devem ser feitos para uma classe/função/seja lá o que for.

## Conclusão

Falamos sobre

-   Refatoração
-   Testes unitários
-   Desenvolvimento de unidade

O que podemos começar a ver é que essas facetas do desenvolvimento de software reforçam uma à outra.

### Refatoração

-   Nos dá sinais sobre nossos testes unitários. Se precisamos fazer validações manuais, precisamos de mais testes. Se testes estão falhando incorretamente, então nossos testes estão no nível errado de abstração (ou não têm valor e precisam ser deletados).
-   Nos ajuda a lidar com as complexidades dentro e entre nossas unidades.

### Testes unitários

-   Nos dá a garantia para refatoração.
-   Verificam e documentam o comportamento de nossas unidades.

### Unidades (bem definidas)

-   Facilitam a escrita de testes unitários _significativos_.
-   Facilitam a refatoração.

Há um processo que nos ajuda a alcançar um ponto onde podemos refatorar nosso código para lidaar com a complexidade e manter nossos sistemas maleáveis?

## Por que Desenvolvimento Orientado a Testes (TDD)

Algumas pessoas levam as citações de Lehman sobre como o software deve mudar a sério demais e elaboram sistemas complexos demais, gastando muito tempo tentando prever o impossível para criar o sistema extensível "perfeito" e acabam entendendo da forma errada e chegando a lugar nenhum.

Isso vem da época das trevas do software onde um time de analistas costumava perder seis meses escrevendo um documento de requerimentos e a equipe de arquitetura perdia outros seis meses para desenvolvê-lo e alguns anos depois o projeto inteiro falhava.

Eu disse que era uma época das trevas, mas isso ainda acontece!

O movimento ágil nos ensina que precisamos trabalhar de forma iterativa, começando com pouca coisa e evoluindo o software para que tenhamos retorno rápido do design do nosso software e como ele trabalha com usuários reais; o TDD reforça essa abordagem.

O TDD aborda as leis que Lehman fala sobre e outras lições difíceis aprendidas no decorrer da história encorajando uma metodologia de refatoração constante e entrega contínua.

### Etapas pequenas

-   Escrever um teste pequeno para uma unidade do comportamento desejado
-   Verificar que o teste falha com um erro claro (vermelho)
-   Escrever o mínimo de código para fazer o teste passar (verde)
-   Refatorar
-   Repetir

Conforme você pratica, essa mentalidade vai se tornar natural e rápida.

Você vai esperar que esse ciclo de feedback não leve muito tempo e se sentir desconfortável se estiver em um estado em que seu sistema não está "verde" por isso poder indicar que você pode ter deixado algo passar.

Você sempre vai desenvolver de forma a criar funcionalidades pequenas & úteis confortavelmente reforçadas pelo feedback dos seus testes.

## Resumindo

-   O ponto forte do software é que podemos mudá-lo. A _maioria_ dos software requer mudança com o tempo de formas imprevisíveis; não tente pensar muito à frente porque é difícil prever o futuro.
-   Ao invés disso, precisamos criar nosso software de forma que ele possa se manter maleável. Para mudar o software precisamos refatorá-lo conforme ele evolui, ou vai acabar virando uma bagunça.
-   Um bom conjunto de testes pode te ajudar a refatorar mais rápido e de forma menos estressante.
-   Escrever bons testes unitários é um problema de design. Logo, pense em estruturar seu código de forma que ele tenha unidades significativas que possam ser unidas como blocos de Lego.
-   O TDD pode ajudar e te forçar a desenvolver softwares bem fatorados continuamente, reforçados por testes para te ajudar com futuros trabalhos que podem chegar.
