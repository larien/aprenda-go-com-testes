# Servidor HTTP

[**Você encontra todo o código-fonte para este capítulo aqui**](https://github.com/quii/learn-go-with-tests/tree/master/http-server)

Você recebeu o desafio de criar um servidor web para que usuários possam acompanhar quantas partidas os jogadores venceram.

* `GET /players/{name}` deve retornar um número indicando o número total de vitórias
* `POST /players/{name}` deve registrar uma vitória para este nome de jogador, incrementando a cada nova chamada `POST`

Vamos seguir com a abordagem do TDD, criando software que funciona o mais rápido possível, e a cada ciclo fazendo pequenas melhorias até uma solução completa. Com essa abordagem, nós

* Mantemos pequeno o escopo do problema em qualquer momento
* Não perdemos o foco por pensar em muito detalhes
* Se ficamos emperrados ou perdidos, voltando para uma versão anterior do código não perdemos muito trabalho.

## Vermelho, verde, refatore

Por todo o livro, enfatizamos o processo TDD de escrever um teste e ver a falha \(vermelho\), escrever a _menor_ quantidade de código para fazer o teste passar/funcionar \(verde\), e então fazemos a reescrita (refatoração).

A disciplina de escrever a menor quantidade de código é importante para garantir a seguraça que o TDD proporciona. Você deve se empenhar em sair do _vermelho_ o quanto antes.

Kent Beck descreve essa prática como:

> Faça o teste passar rapidamente, cometendo quaisquer pecados necessários nesse processo.

E você pode cometer estes pecados porque vai reescrever o código logo depois, com a segurança garantida pelos testes.

### E se você não fizer assim?

Quanto mais alterações você fizer enquanto seu código estiver em _vermelho_, maiores as chances de você adicionar problemas, não cobertos por testes.

A ideia é escrever, iterativamente, código útil em pequenos passos, guiados pelos testes, para que você não perca foco no objetivo principal.

### A galinha e o ovo

Como podemos construir isso de forma incremental? Não podemos obter (`GET`) um nome de jogador sem ter registrado nada anteriormente, e parece complicado saber se o `POST` funcionou sem o endpoint `GET` já implementado.

E é nesse ponto que o _mocking_ vai nos ajudar.

(Nota do tradutor: A expressão _mocking_ significa "zombar", "fazer piada" ou "enganar". Mantemos a expressão original por ser uma expressão comum na literatura em português, por falta de tradução melhor)

* o `GET` precisa de uma _coisa_ `PlayerStore` para obter pontuações de um nome de jogador. Isso deve ser uma interface, para que, ao executar os testes, seja possível criar um código simples de esboço para testar o código sem precisar, neste momento, implementar o código final que será usado para armazenar os dados.
* para o `POST`, podemos _inspecionar_ as chamadas feitas a `PlayerStore` para ter certeza de que os dados são armazenados corretamente. Nossa implementação de gravação dos dados não estará vinculada à busca dos dados.
* para ver código rodando rapidamente vamos fazer uma implementação simples de armazenamento dos dados na memória, e depois podemos criar uma implementação que dá suporte ao mecanismo de armazenamento de preferência.

## Escrevendo o teste primeiro

Podemos escrever um teste e fazer funcionar retornando um valor predeterminado para nos ajudar a começar. Kent Beck se refere a isso como "Fazer de conta". Uma vez que temos um teste funcionando podemos escrever mais testes que nos ajudem a remover este valor predeterminado (constante).

Com este pequeno passo, nós começamos a ter uma estrutura inicial para o projeto funcionando corretamente, sem nos preocuparmos demais com a lógica da aplicação.

Para criar um servidor web (uma aplicação que recebe chamadas via protocolo HTTP) em Go, você vai chamar, normalmente, a função [ListenAndServe](https://golang.org/pkg/net/http/#ListenAndServe).

```go
func ListenAndServe(endereco string, handler Handler) error
```

Isso vai iniciar um servidor web _escutando_ em uma porta, criando uma gorotina para cada requisição e repassando para um [`Handler`](https://golang.org/pkg/net/http/#Handler) (um Handler é um _Tratador_, que recebe a requisição e avalia o que fazer com os dados).

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Esta interface define uma única função que espera dois argumentos, o primeiro que indica onde _escrevemos a resposta_ e o outro com a requisição HTTP que nos foi enviada.

Vamos escrever um teste para a função `PlayerServer` que recebe estes dois argumentos. A requisição enviada serve para obter a pontuação de um Nome de Jogador, que esperamos que seja `"20"`.

```go
func TestGETPlayers(t *testing.T) {
    t.Run("returns Pepper's score", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
        response := httptest.NewRecorder()

        PlayerServer(response, request)

        got := response.Body.String()
        want := "20"

        if got != want {
            t.Errorf("got '%s', want '%s'", got, want)
        }
    })
}
```

Para testar nosso servidor, vamos precisar de um `Request` (_Requisição_) para enviar a requisição ao servidor, e então queremos _inspecionar_ o que o nosso Handler escreve para o `ResponseWriter`.

* Nós usamos o `http.NewRequest` para criar uma requisição. O primeiro argumento é o método da requisição e o segundo é o caminho (_path_) da requisição. O valor `nil` para o segundo argumento corresponde ao corpo (_body_) da requisição, que não precisamos definir para este teste.
* `net/http/httptest` já tem um _inspecionador_ criado para nós, chamado `ResponseRecorder`, então podemos usá-lo. Este possui muitos métodos úteis para inspecionar o que foi escrito como resposta.

## Tente rodar o teste

`./server_test.go:13:2: undefined: PlayerServer`

## Escreva a quantidade mínima de código para o que teste passe e verifique a falha indicada na responta do teste

O compilador está aqui para ajuda, ouça o que ele diz.

Crie a `PlayerServer`

```go
func PlayerServer() {}
```

Tente novamente

```text
./server_test.go:13:14: too many arguments in call to PlayerServer
    have (*httptest.ResponseRecorder, *http.Request)
    want ()
```

Adicione os argumentos à função

```go
import "net/http"

func PlayerServer(w http.ResponseWriter, r *http.Request) {

}
```

Agora o código compila, e o teste falha.

```text
=== RUN   TestGETPlayers/returns_Pepper's_score
    --- FAIL: TestGETPlayers/returns_Pepper's_score (0.00s)
        server_test.go:20: got '', want '20'
```

## Escreva código suficiente para fazer o teste funcionar

Do capítulo sobre injeção de dependências, falamos sobre servidores HTP com a função `Greet`. Aprendemos que a função `ResponseWriter` também implementa a interface `Writer` do pacote io, então podemos usar `fmt.Fprint` para enviar strings como respostas HTTP.

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "20")
}
```

O teste agora deve funcionar.

## Complete a estrutura

Nós queremos converter isso em uma aplicação. Isso é importante porque

* Teremos _software funcionando_, não queremos escrever testes apenas por escrever, e é bom ver código que funciona.
* Conforme refatoramos o código, é provável que vamos mudar a estrutura do programa. Nós queremos garantir que isso é refletido em nossa aplicação também, como parte da abordagem incremental.

Crie um novo arquivo para nossa aplicação, com o código abaixo.

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    handler := http.HandlerFunc(PlayerServer)
    if err := http.ListenAndServe(":5000", handler); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

Até o momento, todo o código de nosso aplicativo está em apenas um arquivo; no entanto, essa não é uma prática recomendada para projetos maiores onde você deseja separar as coisas em arquivos diferentes.

Para executar isso, rode o comando `go build`, que vai pegar todos os arquivos terminados em `.go` neste diretório e construir seu programa. E então você pode executar o programa rodando `./myprogram`.

### `http.HandlerFunc`

Anteriormente, vimos que precisamos implementar a interface `Handler` para criar um servidor. _Normalmente_ fazemos isso criando um `struct` e fazendo com que ele implemente esta interface. No entanto, usualmente utilizamos as _structs_ para armazenar dados, mas como _nesse momento_ não armazenamos um estado, não parece certo criar um _struct_ para isso.

[HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc) nos ajuda a evitar isso.

> o tipo HandlerFunc é um adaptador que permite usar funções comuns como tratadores (_handlers_). Se *f* é uma função com a assinatura adequada, HandlerFunc\(f\) é um Handler que chama *f*.

```go
type HandlerFunc func(ResponseWriter, *Request)
```

Então usamos isso para adaptar a função `PlayerServer` para que ele esteja de acordo com a interface `Handler`.

### `http.ListenAndServe(":5000"...)`

`ListenAndServe` recebe como parâmetro um número de porta para escutar em um `Handler`. Se a porta já estiver sendo usada, será retornado um `error` para que, usando um comando `if`, possamos capturar esse erro e registrar o probema para o usuário.

O que vamos fazer agora é escrever _outro_ teste para nos forçar a fazer uma mudança positiva para tentar nos afastar do valor predefinido.

## Escreva o teste primeiro

Vamos adicionar outro subteste aos nossos testes, que tenta obter a pontuação de um jogador diferente, o que quebrará nossa abordagem que usa um código predefinido.

```go
t.Run("returns Floyd's score", func(t *testing.T) {
    request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
    response := httptest.NewRecorder()

    PlayerServer(response, request)

    got := response.Body.String()
    want := "10"

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
})
```

Você deve estar pensando

> Certamente precisamos de algum tipo de armazenamento para controlar qual jogador recebe qual pontuação. É estranho que os valores pareçam tão predefinidos em nossos testes.

Lembre-se de que estamos apenas tentando dar os menores passos possíveis; e por isso estamos, nesse momento, tentando invalidar o valor da constante.

## Tente rodar o próximo teste

```text
=== RUN   TestGETPlayers/returns_Pepper's_score
    --- PASS: TestGETPlayers/returns_Pepper's_score (0.00s)
=== RUN   TestGETPlayers/returns_Floyd's_score
    --- FAIL: TestGETPlayers/returns_Floyd's_score (0.00s)
        server_test.go:34: got '20', want '10'
```

## Escreva código suficiente para fazer passar

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    if player == "Pepper" {
        fmt.Fprint(w, "20")
        return
    }

    if player == "Floyd" {
        fmt.Fprint(w, "10")
        return
    }
}
```

Este teste nos forçou a olhar para a URL da requisição e tomar uma decisão. Embora ainda estamos pensando em como armazenar os dados do jogador e as interfaces, na verdade o próximo passo a ser dado está relacionado ao _roteamento_ (_routing_).

Se tivéssemos começado com o código de armazenamento dos dados, a quantidade de alterações que precisaríamos fazer seria muito grande. **Este é um pequeno paso em relação ao nosso objetivo final e foi guidado pelos testes**.

Estamos resistindo, nesse momento, à tentação de usar alguma biblioteca de roteamento, e queremos apenas dar o menor passo para fazer nossos testes funcionarem.

`r.URL.Path` retorna o caminho da request, e então usamos a sintaxe de slice para obter a parte final , depois de `/players/`. Não é o recomendado por não ser muito robusto, mas resolve o problema por enquanto.

## Refatorar

Podemos simplificar o `PlayerServer` separando a parte de obtenção da pontuação em uma função.

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    fmt.Fprint(w, GetPlayerScore(player))
}

func GetPlayerScore(name string) string {
    if name == "Pepper" {
        return "20"
    }

    if name == "Floyd" {
        return "10"
    }

    return ""
}
```

E podemos eliminar as repetições de parte do código dos testes montando algumas funções auxiliares("_helpers_")

```go
func TestGETPlayers(t *testing.T) {
    t.Run("returns Pepper's score", func(t *testing.T) {
        request := newGetScoreRequest("Pepper")
        response := httptest.NewRecorder()

        PlayerServer(response, request)

        assertResponseBody(t, response.Body.String(), "20")
    })

    t.Run("returns Floyd's score", func(t *testing.T) {
        request := newGetScoreRequest("Floyd")
        response := httptest.NewRecorder()

        PlayerServer(response, request)

        assertResponseBody(t, response.Body.String(), "10")
    })
}

func newGetScoreRequest(name string) *http.Request {
    req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
    return req
}

func assertResponseBody(t *testing.T, got, want string) {
    t.Helper()
    if got != want {
        t.Errorf("response body is wrong, got '%s' want '%s'", got, want)
    }
}
```

Ainda assim, ainda não estamos felizes. Não parece correto que o servidor sabe as pontuações.

Mas nossa refatoração nos mostra claramente o que fazer.

Nós movemos o cálculo de pontuação pra fora do código principal que trata a requisição (_handler_) para uma função `GetPlayerScore`. Isso parece ser o lugar correto para isolar as responsabilidades usando interfaces.

Vamos alterar a função que refatoramos para ser uma interface

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
}
```

Para que o `PlayerServer` consiga usar o `PlayerStore`, é necessário ter uma referência a ele. Agora nos parece o momento certo para alterar nossa arquitetura, e nosso `PlayerServer` agora é uma `struct`.

```go
type PlayerServer struct {
    store PlayerStore
}
```

E então, vamos implementar a interface do `Handler` adicionando um método à nossa nova struct e adicionado neste método o código existente.


```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]
    fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

Outra alteração a fazer: agora usamos a `store.GetPlayerStore` para obter a pontuação, ao invés da função local definida anteriormente \(e que podemos remover\).

Abaixo, a listagem completa do servidor

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
}

type PlayerServer struct {
    store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]
    fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

### Ajustar os problemas

Fizemos muitas mudanças, e sabemos que nossos testes não irão funcionar e a compilação deixou de funcionar nesse momento; mas relaxa, e deixa o compilador fazer o trabalho.

`./main.go:9:58: type PlayerServer is not an expression`

Precisamos mudar os nossos testes, que agora devem criar uma nova instânia de `PlayerServer` e então chamar o método `ServeHTTP`.

```go
func TestGETPlayers(t *testing.T) {
    server := &PlayerServer{}

    t.Run("returns Pepper's score", func(t *testing.T) {
        request := newGetScoreRequest("Pepper")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertResponseBody(t, response.Body.String(), "20")
    })

    t.Run("returns Floyd's score", func(t *testing.T) {
        request := newGetScoreRequest("Floyd")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertResponseBody(t, response.Body.String(), "10")
    })
}
```

Perceba que ainda não nos preocupamos, _por enquanto_, com o armazenamento dos dados, nós apenas queremos a compilação funcionando o quanto antes.

Você deve ter o hábito de priorizar, sempre, código que compila antes de ter código que passa nos testes.

Adicionando mais funcionalidades \(como códigos esboço de armazenamento\) a um código que não ainda não compila, nos arriscamos a ter, potencialmente, _mais_ problemas de compilação.

Agora `main.go`não vai compilar pelas mesmas razões.

```go
func main() {
    server := &PlayerServer{}

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

Agora tudo compila, mas os testes falham.

```text
=== RUN   TestGETPlayers/returns_the_Pepper's_score
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
    panic: runtime error: invalid memory address or nil pointer dereference
```

Isso porque não passamos um `PlayerStore` em nossos testes. Precisamos fazer um código de esboço para nos ajudar.

```go
type StubPlayerStore struct {
    scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
    score := s.scores[name]
    return score
}
```

Um `map` é um jeito simples e rápido de fazer um armazenamento chave/valor de esboço para os nossos testes. Agora vamos criar um desses armazenamentos para os nosso testes e inserir em nosso `PlayerServer`.

```go
func TestGETPlayers(t *testing.T) {
    store := StubPlayerStore{
        map[string]int{
            "Pepper": 20,
            "Floyd":  10,
        },
    }
    server := &PlayerServer{&store}

    t.Run("returns Pepper's score", func(t *testing.T) {
        request := newGetScoreRequest("Pepper")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertResponseBody(t, response.Body.String(), "20")
    })

    t.Run("returns Floyd's score", func(t *testing.T) {
        request := newGetScoreRequest("Floyd")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertResponseBody(t, response.Body.String(), "10")
    })
}
```

Nossos testes agora passam, e parecem melhores. Agora a _intenção_ do nosso código é clara, por conta da adição do armazenamento. Estamos dizendo ao leitor que, por termos _este dado em um `PlayerStore`, quando você o usar com um  `PlayerServer` você deve obter as respostas definidas.

### Rodar a aplicação

Agora que nossos testes estão passando, a última coisa que precisamos fazer para completar a refatoração é verificar se a aplicação está funcionando. O programa deve iniciar, mas você vai receber uma mensagem horrível se tentar acessar o servidor em `http://localhost:5000/players/Pepper`.

E a razão pra isso é: não passamos um `PlayerStore`.

Precisamos fazer uma implementação de um, mas isso é difícil no momento, já que não estamos armazenando nenhum dado significativo, por isso precisará ser um valor predefinido por enquanto.

```go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
    return 123
}

func main() {
    server := &PlayerServer{&InMemoryPlayerStore{}}

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

Se você rodao novamente o `go build` e acessar a mesma URL você deve receber `"123"`. Não é fantástico, mas até armazenarmos os dados, é o melhor que podemos fazer.

Temos algumas opções para decidir o que fazer agora

* Tratar o cenário onde o jogador não existe
* Tratar o cenário de `POST /players/{name}`
* Não foi exatamente bom perceber que nossa aplicação principal iniciou mas não funcionou. Tivemos que testar manualmente para ver o problema.

Enquanto o cenário do `POST` nos deixa mais perto do "caminho ideal", eu sinto que vai ser mais fácil atacar o cenário de "jogador não existente" antes, já que estamos neste assunto. Veremos os outros itens posteriormente.

## Escreva o teste primeiro.

Adicione o cenário de um jogador inexistente aos nossos testes


```go
t.Run("returns 404 on missing players", func(t *testing.T) {
    request := newGetScoreRequest("Apollo")
    response := httptest.NewRecorder()

    server.ServeHTTP(response, request)

    got := response.Code
    want := http.StatusNotFound

    if got != want {
        t.Errorf("got status %d want %d", got, want)
    }
})
```

## Tente rodar o teste

```text
=== RUN   TestGETPlayers/returns_404_on_missing_players
    --- FAIL: TestGETPlayers/returns_404_on_missing_players (0.00s)
        server_test.go:56: got status 200 want 404
```

## Escreva código necessário para que o teste funcione

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    w.WriteHeader(http.StatusNotFound)

    fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

Às vezes eu me incomodo quando os defensores do TDD dizem "tenha certeza de você escreveu apenas a mínima quantidade de código para fazer o teste funcionar", porque isso me parece muito pedante.

Mas este cenário ilustra muito bem o que querem dizer. Eu fiz o mínimo \(sabendo que não era a implementação correta\), que foi retornar um `StatusNotFound`em **todas as respostas**, mas todos os nossos testes estão passando!

**Implementando o mínimo par que os testes passem pode evidenciar lacunas nos testes**. Em nosso caso, nós não estamos validando que devemos receber um `StatusOK` quando jogadores _existem_ em nosso armazenamento.

Atualize os dois outros testes para validr o retorno e ajuste o código.

Eis os novos testes

```go
func TestGETPlayers(t *testing.T) {
    store := StubPlayerStore{
        map[string]int{
            "Pepper": 20,
            "Floyd":  10,
        },
    }
    server := &PlayerServer{&store}

    t.Run("returns Pepper's score", func(t *testing.T) {
        request := newGetScoreRequest("Pepper")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusOK)
        assertResponseBody(t, response.Body.String(), "20")
    })

    t.Run("returns Floyd's score", func(t *testing.T) {
        request := newGetScoreRequest("Floyd")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusOK)
        assertResponseBody(t, response.Body.String(), "10")
    })

    t.Run("returns 404 on missing players", func(t *testing.T) {
        request := newGetScoreRequest("Apollo")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusNotFound)
    })
}

func assertStatus(t *testing.T, got, want int) {
    t.Helper()
    if got != want {
        t.Errorf("did not get correct status, got %d, want %d", got, want)
    }
}

func newGetScoreRequest(name string) *http.Request {
    req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
    return req
}

func assertResponseBody(t *testing.T, got, want string) {
    t.Helper()
    if got != want {
        t.Errorf("response body is wrong, got '%s' want '%s'", got, want)
    }
}
```

Estamos verificando o `status` (código HTTP de retorno) em todos os nossos testes, por isso existe a função auxiliar `assertStatus` para ajudar com isso.

Agora os primeiros dois testes falham porque o `status` recebido é 404, ao invés do esperado 200. Então vamos corrigir o `PlayerServer` para que retorne *não encontrado* (HTTP status 404) se a pontuação for 0.

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    score := p.store.GetPlayerScore(player)

    if score == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, score)
}
```

### Armazenando pontuações

Agora que podemos obter pontuações de um armazenamento, também podemos armazenar novas pontuações.

## Escreva os testes primeiro

```go
func TestStoreWins(t *testing.T) {
    store := StubPlayerStore{
        map[string]int{},
    }
    server := &PlayerServer{&store}

    t.Run("it returns accepted on POST", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusAccepted)
    })
}
```

Inicialmente vamos verificar se obtemos o status HTTP correto ao fazer a requisição em uma rota específica usando POST. Isso nos permite preparar o caminho da funcionalidade que aceita um tipo diferente de requisição e tratar de forma diferente a requisição para `GET /players/{name}`. Uma vez que isso funciona como esperaodo, então podemos começar a testar a interação do nosso tratador com o armazenamento.

## Tente rodar o teste

```text
=== RUN   TestStoreWins/it_returns_accepted_on_POST
    --- FAIL: TestStoreWins/it_returns_accepted_on_POST (0.00s)
        server_test.go:70: did not get correct status, got 404, want 202
```

## Escreva código suficiente pra fazer passar

Lembre-se que estamos comentendo pecados deliberadamente, então um comando `if` para identificar o método da requisi~ao vai resolver o problema.

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodPost {
        w.WriteHeader(http.StatusAccepted)
        return
    }

    player := r.URL.Path[len("/players/"):]

    score := p.store.GetPlayerScore(player)

    if score == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, score)
}
```

## Refatorar

O tratador parece um pouco bagunçado agora. Vamos separar o códido para ficar simples de entender e isolar as diferentes funcionalidade em novas funções.

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    switch r.Method {
    case http.MethodPost:
        p.processWin(w)
    case http.MethodGet:
        p.showScore(w, r)
    }

}

func (p *PlayerServer) showScore(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    score := p.store.GetPlayerScore(player)

    if score == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter) {
    w.WriteHeader(http.StatusAccepted)
}
```

This makes the routing aspect of `ServeHTTP` a bit clearer and means our next iterations on storing can just be inside `processWin`.

Next, we want to check that when we do our `POST /players/{name}` that our `PlayerStore` is told to record the win.

## Write the test first

We can accomplish this by extending our `StubPlayerStore` with a new `RecordWin` method and then spy on its invocations.

```go
type StubPlayerStore struct {
    scores   map[string]int
    winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
    score := s.scores[name]
    return score
}

func (s *StubPlayerStore) RecordWin(name string) {
    s.winCalls = append(s.winCalls, name)
}
```

Now extend our test to check the number of invocations for a start

```go
func TestStoreWins(t *testing.T) {
    store := StubPlayerStore{
        map[string]int{},
    }
    server := &PlayerServer{&store}

    t.Run("it records wins when POST", func(t *testing.T) {
        request := newPostWinRequest("Pepper")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusAccepted)

        if len(store.winCalls) != 1 {
            t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
        }
    })
}

func newPostWinRequest(name string) *http.Request {
    req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
    return req
}
```

## Try to run the test

```text
./server_test.go:26:20: too few values in struct initializer
./server_test.go:65:20: too few values in struct initializer
```

## Write the minimal amount of code for the test to run and check the failing test output

We need to update our code where we create a `StubPlayerStore` as we've added a new field

```go
store := StubPlayerStore{
    map[string]int{},
    nil,
}
```

```text
--- FAIL: TestStoreWins (0.00s)
    --- FAIL: TestStoreWins/it_records_wins_when_POST (0.00s)
        server_test.go:80: got 0 calls to RecordWin want 1
```

## Write enough code to make it pass

As we're only asserting the number of calls rather than the specific values it makes our initial iteration a little smaller.

We need to update `PlayerServer`'s idea of what a `PlayerStore` is by changing the interface if we're going to be able to call `RecordWin`.

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
}
```

By doing this `main` no longer compiles

```text
./main.go:17:46: cannot use InMemoryPlayerStore literal (type *InMemoryPlayerStore) as type PlayerStore in field value:
    *InMemoryPlayerStore does not implement PlayerStore (missing RecordWin method)
```

The compiler tells us what's wrong. Let's update `InMemoryPlayerStore` to have that method.

```go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) RecordWin(name string) {}
```

Try and run the tests and we should be back to compiling code - but the test is still failing.

Now that `PlayerStore` has `RecordWin` we can call it within our `PlayerServer`

```go
func (p *PlayerServer) processWin(w http.ResponseWriter) {
    p.store.RecordWin("Bob")
    w.WriteHeader(http.StatusAccepted)
}
```

Run the tests and it should be passing! Obviously `"Bob"` isn't exactly what we want to send to `RecordWin`, so let's further refine the test.

## Write the test first

```go
t.Run("it records wins on POST", func(t *testing.T) {
    player := "Pepper"

    request := newPostWinRequest(player)
    response := httptest.NewRecorder()

    server.ServeHTTP(response, request)

    assertStatus(t, response.Code, http.StatusAccepted)

    if len(store.winCalls) != 1 {
        t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
    }

    if store.winCalls[0] != player {
        t.Errorf("did not store correct winner got '%s' want '%s'", store.winCalls[0], player)
    }
})
```

Now that we know there is one element in our `winCalls` slice we can safely reference the first one and check it is equal to `player`.

## Try to run the test

```text
=== RUN   TestStoreWins/it_records_wins_on_POST
    --- FAIL: TestStoreWins/it_records_wins_on_POST (0.00s)
        server_test.go:86: did not store correct winner got 'Bob' want 'Pepper'
```

## Write enough code to make it pass

```go
func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]
    p.store.RecordWin(player)
    w.WriteHeader(http.StatusAccepted)
}
```

We changed `processWin` to take `http.Request` so we can look at the URL to extract the player's name. Once we have that we can call our `store` with the correct value to make the test pass.

## Refactor

We can DRY up this code a bit as we're extracting the player name the same way in two places

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    switch r.Method {
    case http.MethodPost:
        p.processWin(w, player)
    case http.MethodGet:
        p.showScore(w, player)
    }
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
    score := p.store.GetPlayerScore(player)

    if score == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
    p.store.RecordWin(player)
    w.WriteHeader(http.StatusAccepted)
}
```

Even though our tests are passing we don't really have working software. If you try and run `main` and use the software as intended it doesn't work because we haven't got round to implementing `PlayerStore` correctly. This is fine though; by focusing on our handler we have identified the interface that we need, rather than trying to design it up-front.

We _could_ start writing some tests around our `InMemoryPlayerStore` but it's only here temporarily until we implement a more robust way of persisting player scores \(i.e. a database\).

What we'll do for now is write an _integration test_ between our `PlayerServer` and `InMemoryPlayerStore` to finish off the functionality. This will let us get to our goal of being confident our application is working, without having to directly test `InMemoryPlayerStore`. Not only that, but when we get around to implementing `PlayerStore` with a database, we can test that implementation with the same integration test.

### Integration tests

Integration tests can be useful for testing that larger areas of your system work but you must bear in mind:

* They are harder to write
* When they fail, it can be difficult to know why \(usually it's a bug within a component of the integration test\) and so can be harder to fix
* They are sometimes slower to run \(as they often are used with "real" components, like a database\)

For that reason, it is recommended that you research _The Test Pyramid_.

## Write the test first

In the interest of brevity, I am going to show you the final refactored integration test.

```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
    store := InMemoryPlayerStore{}
    server := PlayerServer{&store}
    player := "Pepper"

    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

    response := httptest.NewRecorder()
    server.ServeHTTP(response, newGetScoreRequest(player))
    assertStatus(t, response.Code, http.StatusOK)

    assertResponseBody(t, response.Body.String(), "3")
}
```

* We are creating our two components we are trying to integrate with: `InMemoryPlayerStore` and `PlayerServer`.
* We then fire off 3 requests to record 3 wins for `player`. We're not too concerned about the status codes in this test as it's not relevant to whether they are integrating well.
* The next response we do care about \(so we store a variable `response`\) because we are going to try and get the `player`'s score.

## Try to run the test

```text
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:24: response body is wrong, got '123' want '3'
```

## Write enough code to make it pass

I am going to take some liberties here and write more code than you may be comfortable with without writing a test.

_This is allowed_! We still have a test checking things should be working correctly but it is not around the specific unit we're working with \(`InMemoryPlayerStore`\).

If I were to get stuck in this scenario, I would revert my changes back to the failing test and then write more specific unit tests around `InMemoryPlayerStore` to help me drive out a solution.

```go
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
    return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct{
    store map[string]int
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
    i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
    return i.store[name]
}
```

* We need to store the data so I've added a `map[string]int` to the `InMemoryPlayerStore` struct
* For convenience I've made `NewInMemoryPlayerStore` to initialise the store, and updated the integration test to use it \(`store := NewInMemoryPlayerStore()`\)
* The rest of the code is just wrapping around the `map`

The integration test passes, now we just need to change `main` to use `NewInMemoryPlayerStore()`

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    server := &PlayerServer{NewInMemoryPlayerStore()}

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

Build it, run it and then use `curl` to test it out.

* Run this a few times, change the player names if you like `curl -X POST http://localhost:5000/players/Pepper`
* Check scores with `curl http://localhost:5000/players/Pepper`

Great! You've made a REST-ish service. To take this forward you'd want to pick a data store to persist the scores longer than the length of time the program runs.

* Pick a store \(Bolt? Mongo? Postgres? File system?\)
* Make `PostgresPlayerStore` implement `PlayerStore`
* TDD the functionality so you're sure it works
* Plug it into the integration test, check it's still ok
* Finally plug it into `main`

## Wrapping up

### `http.Handler`

* Implement this interface to create web servers
* Use `http.HandlerFunc` to turn ordinary functions into `http.Handler`s
* Use `httptest.NewRecorder` to pass in as a `ResponseWriter` to let you spy on the responses your handler sends
* Use `http.NewRequest` to construct the requests you expect to come in to your system

### Interfaces, Mocking and DI

* Lets you iteratively build the system up in smaller chunks
* Allows you to develop a handler that needs a storage without needing actual storage
* TDD to drive out the interfaces you need

### Commit sins, then refactor \(and then commit to source control\)

* You need to treat having failing compilation or failing tests as a red situation that you need to get out of as soon as you can.
* Write just the necessary code to get there. _Then_ refactor and make the code nice.
* By trying to do too many changes whilst the code isn't compiling or the tests are failing puts you at risk of compounding the problems.
* Sticking to this approach forces you to write small tests, which means small changes, which helps keep working on complex systems manageable.
