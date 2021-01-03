# Como contribuir

Contribuições são mais que bem vindas. Espero que esse se torne um ótimo lar para guias sobre como aprender Go escrevendo testes. Você pode submeter uma PR ou criar uma issue [aqui](https://github.com/larien/aprenda-go-com-testes).

## O que estamos buscando

-   Ensinar funcionalidades de Go (conceitos como `if`, `select`, estruturas, métodos etc).
-   Demonstrar as funcionalidades interessantes dentro da biblioteca padrão e o quão fácil é utilizar o TDD para criar um servidor HTTP, por exemplo.
-   Mostrar como o ferramental do Go, incluindo benchmarking, race detectors, etc pode te ajudar a construir um ótimo software.

Se não se sentir confiante em enviar seu próprio guia, criar uma issue para algo que queira aprender também é uma contribuição válida.

## Estilo a ser seguido

-   Sempre reforce o ciclo TDD. Dê uma olhada no [Exemplo de Capítulo](exemplo.md).
-   Dê ênfase em iterar sobre funcionalidades orientando-se por testes. O exemplo [Olá, mundo](../primeiros-passos-com-go/ola-mundo/ola-mundo.md) funciona bem porque aos poucos tornamos o código mais sofisticado e aprendemos novas técnicas _orientadas por testes_. Por exemplo:
    -   `Hello()` &lt;- como escrever funções e retornar tipos.
    -   `Hello(name string)` &lt;- argumentos, constantes.
    -   `Hello(name string)` &lt;- padrão para "mundo" usando `if`.
    -   `Hello(name, language string)` &lt;- `switch`.
-   Tente diminuir a barreira de conhecimento com explicações claras e simples.
    -   É importante pensar em exemplos que demonstrem o que você está tentando ensinar sem confundir a leitura com outras funcionalidades.
    -   Por exemplo: você pode aprender `structs` sem entender ponteiros.
    -   Seja breve.
-   Siga o [guia de estilo para Comentários de Revisão de Código](https://github.com/golang/go/wiki/CodeReviewComments). É importante ter um estilo consistente em todas as seções.
-   Sua seção deve ter uma aplicação executável no final (como um `package main` com uma função `main`) para que as pessoas possam vê-la em ação e brincar com ela.
-   Todos os testes devem passar.
-   Execute o `/build.sh` antes de subir uma PR.
