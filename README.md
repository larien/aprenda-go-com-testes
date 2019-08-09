# Introdução

![](.gitbook/assets/red-green-blue-gophers-smaller.png)

[Arte por Denise](https://twitter.com/deniseyu21)

![Build Status](https://travis-ci.org/larien/learn-go-with-tests.svg?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/larien/learn-go-with-tests)](https://goreportcard.com/report/github.com/quii/learn-go-with-tests)

-   Formatos: [Gitbook](https://larien.gitbook.io/aprenda-go-com-testes), [EPUB or PDF](https://github.com/larien/learn-go-with-tests/releases)
-   Versão original: [English](https://quii.gitbook.io/learn-go-with-tests/)
-   Traduções: [中文](https://studygolang.gitbook.io/learn-go-with-tests)

## Motivação

-   Explore a linguagem Go escrevendo testes
-   **Tenha uma base com TDD**. O Go é uma boa linguagem para aprender TDD por ser simples de aprender e ter testes nativamente
-   Tenha confiança de que você será capaz de escrever sistemas robustos e bem testados em Go
-   [Assista a um vídeo ou leia sobre o motivo pelo qual testes unitários e TDD são importantes](meta/motivacao.md)

## Índice

### Primeiros Passos com Go

1. [Instalação do Go](primeiros-passos-com-go/instalacao-do-go.md) - Prepare o ambiente para produtividade.
2. [Olá, mundo](primeiros-passos-com-go/hello-world.md) - Declarando variáveis, constantes, declarações `if`/`else`, switch, escreva seu primeiro programa em Go e seu primeiro teste. Sintaxe de subteste e closures.
3. [Inteiros](primeiros-passos-com-go/integers.md) - Mais conteúdo sobre sintaxe de declaração de função e aprenda novas formas de melhorar a documentação do seu código.
4. [Iteração](primeiros-passos-com-go/iteracao.md) - Aprenda sobre `for` e benchmarking.
5. [Arrays e slices](primeiros-passos-com-go/arrays-e-slices.md) - Aprenda sobre arrays, slices, `len`, variáveis recebidas como argumentos, `range` e cobertura de testes.
6. [Estruturas, métodos e interfaces](primeiros-passos-com-go/structs-methods-and-interfaces.md) - Aprenda sobre `structs`, métodos, `interface` e testes orientados a tabela \(table driven tests\).
7. [Ponteiros e erros](primeiros-passos-com-go/pointers-and-errors.md) - Aprenda sobre ponteiros e erros.
8. [Maps](primeiros-passos-com-go/maps.md) - Aprenda sobre armazenamento de valores na estrutura de dados `map`.
9. [Injeção de dependência](primeiros-passos-com-go/dependency-injection.md) - Aprenda sobre injeção de dependência, qual sua relação com interfaces e uma introdução a I/O.
10. [Mocking](primeiros-passos-com-go/mocks.md) - Use injeção de dependência com mocking para testar um código sem nenhum teste.
11. [Concorrência](primeiros-passos-com-go/concurrency.md) - Aprenda como escrever código concorrente para tornar seu software mais rápido.
12. [Select](primeiros-passos-com-go/select.md) - Aprenda a sincronizar processos assíncronos de forma elegante.
13. [Reflection](primeiros-passos-com-go/reflection.md) - Aprenda sobre reflection.
14. [Sync](primeiros-passos-com-go/sync.md) - Conheça algumas funcionalidades do pacote `sync`, como `WaitGroup` e `Mutex`.
15. [Context](primeiros-passos-com-go/context.md) - Use o pacote `context` para gerenciar e cancelar processos de longa duração.

### Crie uma aplicação

Agora que você já deu seus _Primeiros Passos com Go_, esperamos que você tenha uma base sólida das principais funcionalidades da linguagem e como TDD funciona.

Essa seção envolve a criação de uma aplicação.

Cada capítulo é uma continuação do anterior, expandindo as funcionalidades da aplicação conforme nosso "Product Owner" dita.

Novos conceitos serão apresentados para ajudar a escrever código de qualidade, mas a maior parte do material novo terá relação com o que pode ser feito com a biblioteca padrão do Go.

No final desse capítulo, você deverá ter uma boa ideia de como escrever uma aplicação em Go testada.

-   [Servidor HTTP](criando-uma-aplicacao/http-server.md) - Vamos criar uma aplicação que espera por requisições HTTP e as responde.
-   [JSON, routing e embedding](criando-uma-aplicacao/json.md) - Vamos fazer nossos endpoints retornarem JSON e explorar como trabalhar com rotas.
-   [IO e classificação](criando-uma-aplicacao/io.md) - Vamos persistir e ler nossos dados do disco e falar sobre classificação de dados.
-   [Linha de comando e estrutura do projeto](criando-uma-aplicacao/command-line.md) - Suportar diversas aplicações em uma base de código e ler entradas da linha de comando.
-   [Tempo](criando-uma-aplicacao/time.md) - Usar o pacote `time` para programar atividades.
-   [Websockets](criando-uma-aplicacao/websockets.md) - Aprender a escrever e testar um servidor que usa websockets.

### Dúvidas e respostas

Costumo ver perguntas nas Interwebs como:

> Como testo minha função incrível que faz x, y e z?

Se tiver esse tipo de dúvida, crie uma Issue no GitHub e vou tentar achar tempo para escrever um pequeno capítulo para resolver o problema. Acho que conteúdo como esse é valioso, já que está resolvendo problemas `reais` envolvendo testes que as pessoas têm.

-   [OS exec](perguntas-e-respostas/os-exec.md) - Um exemplo de como podemos usar o sistema operacional para executar comandos para buscar dados e manter nossa lógica de negócio testável.
-   [Tipos de erro](perguntas-e-respostas/error-types.md) - Exemplo de como criar seus próprios tipos de erro para melhorar seus testes e tornar seu código mais fácil de se trabalhar.

## Contribuição

-   _Esse projeto está em desenvolvimento_, tanto seu conteúdo original quanto sua tradução. Se tiver interesse em contribuir, por favor entre em contato.
-   Leia [contribuindo.md](meta/contribuindo.md) para algumas diretrizes.
-   Quer ajudar com a tradução para o português? Leia [traduzindo.md](meta/traduzindo.md) e entenda como o processo de tradução está organizado.
-   Tem ideias? Crie uma issue!

## Explicação

Tenho experiência em apresentar Go a equipes de desenvolvimento e tenho testado abordagens diferentes sobre como evoluir um grupo de pessoas que têm curiosidade sobre Go para criadores extremamente eficazes de sistemas em Go.

### O que não funcionou

#### Ler _o_ livro

Uma abordagem que tentamos foi pegar [o livro azul](https://www.amazon.com.br/Linguagem-Programa%C3%A7%C3%A3o-Go-Alan-Donovan/dp/8575225464) e toda semana discutir um capítulo junto de exercícios.

Amo esse livro, mas ele exige muito comprometimento. O livro é bem detalhado na explicação de conceitos, o que obviamente é ótimo, mas significa que o progresso é lento e uniforme - não é para todo mundo.

Descobri que apenas um pequeno número de pessoas pegaria o capítulo X para ler e faria os exercícios, enquanto que a maioria não.

#### Resolver alguns problemas

Katas são divertidos, mas geralmente se limitam ao escopo de aprender uma linguagem; é improvável que você use goroutines para resolver um kata.

Outro problema é quando você tem níveis diferentes de entusiasmo. Algumas pessoas aprendem mais da linguagem que outras e, quando demonstram o que já fizeram, confundem essas pessoas apresentando funcionalidades que as outras ainda não conhecem.

Isso acaba tornando o aprendizado bem _desestruturado_ e _específico_.

### O que funcionou

De longe, a forma mais eficaz foi apresentar os conceitos da linguagem aos poucos lendo o [go by example](https://gobyexample.com/), explorando-o com exemplos e discutindo-o como um grupo. Essa abordagem foi bem mais interativa do que "leia o capítulo X como lição de casa".

Com o tempo, a equipe ganhou uma base sólida da _gramática_ da linguagem para que conseguíssemos começar a desenvolver sistemas.

Para mim, é semelhante à ideia de praticar escalas quando se tenta aprender a tocar violão.

Não importa quão artístico você seja; é improvável que você crie músicas boas sem entender os fundamentos e praticando os mecanismos.

### O que funcionou para mim

Quando _eu_ aprendo uma nova linguagem de programação, costumo começar brincando em um REPL, mas hora ou outra preciso de mais estrutura.

O que eu gosto de fazer é explorar conceitos e então solidificar as ideias com testes. Testes certificam de que o código que escrevi está correto e documentam a funcionalidade que aprendi.

Usando minha experiência de aprendizado em grupo e a minha própria, vou tentar criar algo que seja útil para outras equipes. Aprender os conceitos escrevendo testes pequenos para que você possa usar suas habilidades de desenvolvimento de software e entregar sistemas ótimos.

## Para quem isso foi feito

-   Pessoas que se interessam em aprender Go.
-   Pessoas que já sabem Go, mas querem explorar testes com TDD.

## O que vamos precisar

-   Um computador!
-   [Go instalado](https://golang.org/)
-   Um editor de texto
-   Experiência com programação. Entendimento de conceitos como `if`, variáveis, funções etc.
-   Se sentir confortável com o terminal

## Traduzido com <3 por

-   Davi Marcondes Moreira

[github](http://github.com/devdrops) [twitter](https://twitter.com/devdrops)

-   Diego Nascimento

[github](https://github.com/diegonvs) [twitter](https://twitter.com/diegonvs97) [linkedin](https://www.linkedin.com/in/dnvs97/)

-   Edmilton Neves

[site](http://edmilton.com.br/) [github](https://github.com/edmilton) [twitter](https://twitter.com/ed_neves) [linkedin](https://www.linkedin.com/in/edmilton/)

-   Jéssica Paz

[site](https://jessicapaz.me) [github](https://github.com/jessicapaz) [twitter](https://twitter.com/jessicamorim42) [linkedin](https://www.linkedin.com/in/jessica-paz/)

-   Lauren Ferreira

[site](https://larien.dev) [github](https://github.com/larien) [twitter](https://twitter.com/larienmf) [linkedin](https://www.linkedin.com/in/lauren-ferreira/)

-   Rafael Acioly

[github](https://github.com/rafa-acioly) [twitter](https://twitter.com/R_acioly) [linkedin](https://www.linkedin.com/in/rafaelacioly/)

## Feedback

-   Crie issues/submita PRs [aqui](https://github.com/quii/learn-go-with-tests) ou [me envie um tweet em @quii](https://twitter.com/quii).
-   Para a versão em português, submita um PR [aqui](https://github.com/larien/learn-go-with-tests) ou entre em contato comigo pelo [meu site](https://larien.dev).

[MIT license](https://github.com/larien/learn-go-with-tests/tree/09aafaeebaef4443e80a6216cc46fa3d7bfdabbb/LICENSE.md)

[Logo criado por egonelbre](https://github.com/egonelbre) Que estrela!
