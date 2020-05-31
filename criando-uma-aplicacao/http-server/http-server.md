# Servidor HTTP

[**Você encontra todo o código-fonte para este capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/http-server)

Você recebeu o desafio de criar um servidor web para que usuários possam acompanhar quantas partidas os jogadores (players) venceram.

* `GET /players/{name}` deve retornar um número indicando o número total de vitórias
* `POST /players/{name}` deve registrar uma vitória para este nome de jogador, incrementando a cada nova chamada `POST`

Vamos seguir com a abordagem do TDD, criando software que funciona o mais rápido possível, e a cada ciclo fazendo pequenas melhorias até uma solução completa. Com essa abordagem, nós

* Mantemos pequeno o escopo do problema em qualquer momento
* Não perdemos o foco por pensar em muito detalhes
* Se ficamos emperrados ou perdidos, podemos voltar para uma versão anterior do código sem perder muito trabalho.

## Vermelho, verde, refatore

Por todo o livro, enfatizamos o processo TDD de escrever um teste e ver a falha \(vermelho\), escrever a _menor_ quantidade de código para fazer o teste passar/funcionar \(verde\), e então fazemos a reescrita (refatoração).

A disciplina de escrever a menor quantidade de código é importante para garantir a segurança que o TDD proporciona. Você deve se empenhar em sair do _vermelho_ o quanto antes.

Kent Beck descreve essa prática como:

> Faça o teste passar rapidamente, cometendo quaisquer pecados necessários nesse processo.

E você pode cometer estes pecados porque vamos reescrever o código logo depois, com a segurança garantida pelos testes.

### E se você não fizer assim?

Quanto mais alterações você fizer enquanto seu código estiver em _vermelho_, maiores as chances de você adicionar problemas, não cobertos por testes.

A ideia é escrever iterativamente código útil em pequenos passos, guiados pelos testes, para que você não perca foco no objetivo principal.

### A galinha e o ovo

Como podemos construir isso de forma incremental? Não podemos obter um jogador (`GET`) sem tê-lo registrado nada anteriormente, e parece complicado saber se o `POST` funcionou sem o endpoint `GET` já implementado.

E é nesse ponto que o _mocking_ vai nos ajudar.

(Nota do tradutor: A expressão _mocking_ significa "zombar", "fazer piada" ou "enganar". Mas em programação, _mocking_ significa criar _algo_, como uma classe ou função, que retorna os valores esperados de forma predefinida.  Mantemos a expressão original por ser uma expressão comum na literatura em português, por falta de tradução melhor)

* o `GET` precisa de uma _coisa_ `PlayerStore` para obter pontuações de um nome de jogador. Isso deve ser uma interface, para que, ao executar os testes, seja possível criar um código simples de esboço para testar o código sem precisar, neste momento, implementar o código final que será usado para armazenar os dados.
* para o `POST`, podemos _inspecionar_ as chamadas feitas a `PlayerStore` para ter certeza de que os dados são armazenados corretamente. Nossa implementação de gravação dos dados não estará vinculada à busca dos dados.
* para ver código rodando rapidamente vamos fazer uma implementação simples de armazenamento dos dados na memória, e depois podemos criar uma implementação que dá suporte ao mecanismo de armazenamento de preferência.

## Escrevendo o teste primeiro

Podemos escrever um teste e fazer funcionar retornando um valor predeterminado para nos ajudar a começar. Kent Beck se refere a isso como "Fazer de conta". Uma vez que temos um teste funcionando podemos escrever mais testes que nos ajudem a remover este valor predeterminado (constante).

Com este pequeno passo, nós começamos a ter uma estrutura inicial para o projeto funcionando corretamente, sem nos preocuparmos demais com a lógica da aplicação.

Para criar um servidor web (uma aplicação que recebe chamadas via protocolo HTTP) em Go, você vai chamar, usualmente, a função [ListenAndServe](https://golang.org/pkg/net/http/#ListenAndServe).

```go
func ListenAndServe(addr string, handler Handler) error
```

Isso vai iniciar um servidor web _escutando_ em uma porta, criando uma gorotina para cada requisição, e repassando para um [`Handler`](https://golang.org/pkg/net/http/#Handler) (um Handler é um _Tratador_, que recebe a requisição e avalia o que fazer com os dados).

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

* Nós usamos o `http.NewRequest` para criar uma requisição. O primeiro argumento é o _método_ da requisição e o segundo é o caminho (_path_) da requisição. O valor `nil` para o segundo argumento corresponde ao corpo (_body_) da requisição, que não precisamos definir para este teste.
* `net/http/httptest` já tem um _inspecionador_ criado para nós, chamado `ResponseRecorder`, então podemos usá-lo. Este possui muitos métodos úteis para inspecionar o que foi escrito como resposta.

## Tente rodar o teste

`./server_test.go:13:2: undefined: PlayerServer`

## Escreva a quantidade mínima de código para o que teste passe e verifique a falha indicada na responta do teste

O compilador está aqui para ajudar, ouça o que ele diz.

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

Do capítulo sobre injeção de dependências, falamos sobre servidores HTTP com a função `Greet`. Aprendemos que a função `ResponseWriter` também implementa a interface `Writer` do pacote io, então podemos usar `fmt.Fprint` para enviar strings como respostas HTTP.

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "20")
}
```

O teste agora deve funcionar.

## Complete a estrutura

Nós queremos converter isso em uma aplicação. Isso é importante porque

* Teremos _software funcionando_; não queremos escrever testes apenas por escrever, e é bom ver código que funciona.
* Conforme refatoramos o código, é provável mudaremos a estrutura do programa. Nós queremos garantir que isso é refletido em nossa aplicação também, como parte da abordagem incremental.

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

Para executar, execute o comando `go build`, que vai pegar todos os arquivos terminados em `.go` neste diretório e construir seu programa. E então você pode executar o programa rodando `./myprogram`.

### `http.HandlerFunc`

Anteriormente, vimos que precisamos implementar a interface `Handler` para criar um servidor. _Normalmente_ fazemos isso criando um `struct` e fazendo com que ele implemente esta interface. No entanto, mesmo que o comum seja utilizar as _structs_ para armazenar dados, _nesse momento_ não armazenamos um estado, então não parece certo criar um _struct_ para isso.

[HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc) nos ajuda a evitar isso.

> O tipo HandlerFunc é um adaptador que permite usar funções comuns como tratadores (_handlers_). Se *f* é uma função com a assinatura adequada, HandlerFunc\(f\) é um Handler que chama *f*.

```go
type HandlerFunc func(ResponseWriter, *Request)
```

Então usamos essa construção para adaptar a função `PlayerServer`, fazendo com que esteja de acordo com a interface `Handler`.

### `http.ListenAndServe(":5000"...)`

`ListenAndServe` recebe como parâmetro um número de porta para escutar em um `Handler`. Se a porta já estiver sendo usada, será retornado um `error` para que, usando um comando `if`, possamos capturar esse erro e registrar o problema para o usuário.

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

Se tivéssemos começado com o código de armazenamento dos dados, a quantidade de alterações que precisaríamos fazer seria muito grande. **Este é um pequeno passo em relação ao nosso objetivo final e foi guiado pelos testes**.

Estamos resistindo, nesse momento, à tentação de usar alguma biblioteca de roteamento, e queremos apenas dar o menor passo para fazer nossos testes funcionarem.

`r.URL.Path` retorna o caminho da request, e então usamos a sintaxe de slice para obter a parte final, depois de `/players/`. Não é o recomendado por não ser muito robusto, mas resolve o problema por enquanto.

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

Ainda assim, ainda não estamos felizes. Não parece correto que o servidor saiba as pontuações.

Mas nossa refatoração nos mostra claramente o que fazer.

Nós movemos o cálculo de pontuação para fora do código principal que trata a requisição (_handler_) para uma função `GetPlayerScore`. Isso parece ser o lugar correto para isolar as responsabilidades usando interfaces.

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

E agora, vamos implementar a interface do `Handler` adicionando um método à nossa nova struct e adicionado neste método o código existente.

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

Fizemos muitas mudanças, e sabemos que nossos testes não irão funcionar e a compilação deixou de funcionar nesse momento; mas relaxe, e deixe o compilador fazer o trabalho.

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

Você deve ter o hábito de priorizar, _sempre_, código que compila antes de ter código que passa nos testes.

Adicionando mais funcionalidades \(como códigos de esboço - _stub_ - de armazenamento\) a um código que não ainda não compila, nos arriscamos a ter, potencialmente, _mais_ problemas de compilação.

Agora `main.go` não vai compilar pelas mesmas razões.

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

Isso porque não passamos um `PlayerStore` em nossos testes. Precisamos fazer um código de esboço \(stub\) para nos ajudar.

```go
type StubPlayerStore struct {
    scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
    score := s.scores[name]
    return score
}
```

Um `map` é um jeito simples e rápido de fazer um armazenamento chave/valor de  para os nossos testes. Agora vamos criar um desses armazenamentos para os nosso testes e inserir em nosso `PlayerServer`.

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

Nossos testes agora passam, e parecem melhores. Agora a _intenção_ do nosso código é clara, por conta da adição do armazenamento. Estamos dizendo a quem lê o código que, por termos _este dado em um `PlayerStore`_, quando você o usar com um  `PlayerServer` você deve obter as respostas definidas.

### Rodar a aplicação

Agora que nossos testes estão passando, a última coisa que precisamos fazer para completar a refatoração é verificar se a aplicação está funcionando. O programa deve iniciar, mas você vai receber uma mensagem horrível se tentar acessar o servidor em `http://localhost:5000/players/Pepper`.

E a razão pra isso é: não informamos um `PlayerStore`.

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

Se você rodar novamente o `go build` e acessar a mesma URL você deve receber `"123"`. Não é fantástico, mas até armazenarmos os dados, é o melhor que podemos fazer.

Temos algumas opções para decidir o que fazer agora

* Tratar o cenário onde o jogador não existe
* Tratar o cenário de `POST /players/{name}`
* Não foi exatamente bom perceber que nossa aplicação principal iniciou mas não funcionou. Tivemos que testar manualmente para ver o problema.

Enquanto o cenário do `POST` nos deixa mais perto do "caminho ideal", eu sinto que vai ser mais fácil atacar o cenário de "jogador não existente" antes, já que estamos neste assunto. Veremos os outros itens posteriormente.

## Escreva o teste primeiro

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

Mas este cenário ilustra muito bem o que querem dizer. Eu fiz o mínimo \(sabendo que não era a implementação correta\), que foi retornar um `StatusNotFound` em **todas as respostas**, mas todos os nossos testes estão passando!

**Implementando o mínimo par que os testes passem pode evidenciar lacunas nos testes**. Em nosso caso, nós não estamos validando que devemos receber um `StatusOK` quando jogadores _existem_ em nosso armazenamento.

Atualize os outros dois testes para validar o retorno e corrija o código.

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

Inicialmente vamos verificar se obtemos o status HTTP correto ao fazer a requisição em uma rota específica usando POST. Isso nos permite preparar o caminho da funcionalidade que aceita um tipo diferente de requisição e tratar de forma diferente a requisição para `GET /players/{name}`. Uma vez que isso funciona como esperado, então podemos começar a testar a interação do nosso _handler_ com o armazenamento.

## Tente rodar o teste

```text
=== RUN   TestStoreWins/it_returns_accepted_on_POST
    --- FAIL: TestStoreWins/it_returns_accepted_on_POST (0.00s)
        server_test.go:70: did not get correct status, got 404, want 202
```

## Escreva código suficiente pra fazer passar

Lembre-se que estamos cometendo pecados deliberadamente, então um comando `if` para identificar o método da requisição vai resolver o problema.

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

O _handler_ parece um pouco bagunçado agora. Vamos separar o código para ficar simples de entender e isolar as diferentes funcionalidades em novas funções.

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

Isso faz com que a função de roteamento do `ServeHTTP` esteja mais clara; e também permite que, em nossas próximas iterações, o código para armazenamento possa estar dentro de `processWin`.

Agora, queremos verificar que, quando fazemos a chamada `POST` a `/players/{name}`, nosso `PlayerStore` registra a vitória.

## Escreva primeiro o teste

Vamos implementar isso estendendo o `StubPlayerStore` com um novo método `RecordWin` e então inspecionar as chamadas.

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

Agora, para começar, estendemos o teste para verificar a quantidade de chamadas

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

## Tente rodar o teste

```text
./server_test.go:26:20: too few values in struct initializer
./server_test.go:65:20: too few values in struct initializer
```

## Escreva a mínima quantidade de código para a execução do teste e verifique a falha indicada no retorno

Como adicionamos um campo, precisamos atualizar o código onde criamos o `StubPlayerStore`

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

## Escreva código suficiente para o teste passar

Como estamos apenas verificando o número de chamadas, e não seus valores específicos, nossa iteração inicial é um pouco menor.

Para conseguir invocar a `RecordWin`, precisamos atualizar a definição de `PlayerStore` para que o `PlayerServer` funcione como esperado.

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
}
```

E, ao fazer isso, `main` não compila mais

```text
./main.go:17:46: cannot use InMemoryPlayerStore literal (type *InMemoryPlayerStore) as type PlayerStore in field value:
    *InMemoryPlayerStore does not implement PlayerStore (missing RecordWin method)
```

O compilador nos informa o que está errado. Vamos alterar `InMemoryPlayerStore`, adicionando esse método.

```go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) RecordWin(name string) {}
```

Com essa alteração, o código volta a compilar - mas os testes ainda falham.

Agora que `PlayerStore` tem o método `RecordWin`, podemos chamar de dentro do nosso `PlayerServer`

```go
func (p *PlayerServer) processWin(w http.ResponseWriter) {
    p.store.RecordWin("Bob")
    w.WriteHeader(http.StatusAccepted)
}
```

Rode os testes e deve estar funcionando sem erros! Claro, `"Bob"` não é bem o que queremos enviar para `RecordWin`, então vamos ajustar os testes.

## Escreva os testes primeiro

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

Agora sabemos que existe um elemento no slice `winCalls`, e então podemos acessar, sem erros, o primeiro elemento e verificar se é igual a `player`.

## Tente rodar o teste

```text
=== RUN   TestStoreWins/it_records_wins_on_POST
    --- FAIL: TestStoreWins/it_records_wins_on_POST (0.00s)
        server_test.go:86: did not store correct winner got 'Bob' want 'Pepper'
```

## Escreva código suficiente para o teste passar

```go
func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]
    p.store.RecordWin(player)
    w.WriteHeader(http.StatusAccepted)
}
```

Mudamos `processWin` para obter a `http.Request`, para conseguir extrair o nome do jogador da URL. Com o nome, podemos chamar o `store` com o valor correto para fazer os testes passarem.

## Refatorar

Podemos eliminar repetições no código, porque estamos obtendo o nome do "player" do mesmo jeito em dois lugares diferentes.

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

Mesmo com os testes passando, não temos código funcionando de forma ideal. Se executar a `main` e usar o programa como planejado, não vai funcionar porque ainda não nos dedicamos a implementar corretamente `PlayerStore`. Mas isso não é um problema; como focamos no tratamento da requisição, identificamos a interface necessária, ao invés de tentar definir antecipadamente.

_Poderíamos_ começar a escrever alguns testes para a `InMemoryPlayerStore`, mas ela é apenas uma solução temporária até a implementação de um modo mais robusto de registrar as pontuações \(por exemplo, em um banco de dados\).

O que vamos fazer agora é escrever um _teste de integração_ entre `PlayerServer` e `InMemoryPlayerStore` para terminar a funcionalidade. Isso vai permitir confiar que a aplicação está funcionando, sem ter que testar diretamente `InMemoryPlayerStore`. E não apenas isso, mas quando implementarmos `PlayerStore` com um banco de dados, usaremos esse mesmo teste para verificar se a implementação funciona como esperado.

### Testes de integração

Testes de integração podem ser úteis para testar partes maiores do sistema, mas saiba que:

* São mais difíceis de escrever
* Quando falham, é difícil saber o porquê \(normalmente é um problema dentro de um componente do teste de integração\) e pode ser difícil de corrigir
* Às vezes são mais lentos para rodar \(porque são usados com componentes "reais", como um banco de dados\)

Por isso, é recomendado que pesquise sobre _Pirâmide de Testes_.

## Escreva os testes primeiro

Para ser mais breve, vou te mostrar o teste de integração, já refatorado.

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

* Estamos criando os dois componentes que queremos integrar: `InMemoryPlayerStore` e `PlayerServer`.
* Então fazemos 3 requisições para registrar 3 vitórias para `player`. Não nos preocupamos com os códigos de retorno no teste, porque isso não é relevante para verificar se a integração funciona como esperado.
* Registramos a próxima resposta \(por isso guardamos o valor em `response`\) porque vamos obter a pontuação do `player`.

## Tente rodar o teste

```text
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:24: response body is wrong, got '123' want '3'
```

## Escreva código suficiente para passar

Abaixo, há mais código do que o esperado para se escrever sem ter os testes correspondentes.

_Isso é permitido_! Ainda existem testes verificando se as coisas estão funcionando como esperado, mas não focando na parte específica em que estamos trabalhando \(`InMemoryPlayerStore`\).

Se houvesse algum problema para continuarmos, era só reverter as alterações para antes do teste que falhou e então escrever mais testes unitários específicos para `InMemoryPlayerStore`, que nos ajudariam a encontrar a solução.

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

* Para armazenar os dados, adicionamos um `map[string]int` na struct `InMemoryPlayerStore`
* Para ajudar nos testes, criamos a `NewInMemoryPlayerStore` para inicializar o armazenamento, e o código do teste de integração foi atualizado para usar esta função \(`store := NewInMemoryPlayerStore()`\).
* O resto do código é apenas para fazer o `map` funcionar.

Nosso teste de integração passa, e agora só é preciso mudar o `main` para usar o `NewInMemoryPlayerStore()`

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

Após compilar e rodar, use o `curl` para testar.

* Execute o comando a seguir algumas vezes, mude o nome do jogador se quiser `curl -X POST http://localhost:5000/players/Pepper`
* Verifique a pontuação, rodando `curl http://localhost:5000/players/Pepper`

Ótimo! Criamos um serviço de acordo com os padrões REST! Se quiser continuar, você pode escolher um armazenamento de dados com maior persistência, que não vai perder os dados quando o programa terminar.

* Escolher uma tecnologia de armazenamento \(Bolt? Mongo? Postgres? Sistema de arquivos?\)
* Fazer `PostgresPlayerStore` implementar `PlayerStore`
* Desenvolver a funcionalidade usando TDD para ter certeza de que funciona
* Conectar nos testes de integração, verificar se tudo funciona
* E, finalmente, integrar dentro de `main`.

## Finalizando

### `http.Handler`

* Implemente essa interface para criar servidores web
* Use `http.HandlerFunc` para transformar funções simples em `http.Handler`s
* Use `httptest.NewRecorder` para informar um `ResponseWriter` que permite inspecionar as respostas que a função tratadora envia
* Use `http.NewRequest` para construir as requisições que você espera que seu sistema receba

### Interfaces, _Mocking_ e Injeção de Dependência

* Permitem que você construa a sua aplicação de forma iterativa, um pedaço de cada vez
* Te permite desenvolver uma funcionalidade de tratamento de requisições que precisa de um armazenamento sem precisar exatamente de uma estrutura de armazenamento
* o TDD nos ajudou a definir as interfaces necessárias

### Cometa pecados, e daí refatore \(e então registre no controle de versão\)

* Você precisa tratar falhas na compilação ou nos testes como uma situação urgente, a qual precisa resolver o mais rápido possível.
* Escreva apenas o código necessário para resolver o problema. _Logo depois_ refatore e faça um código melhor
* Ao tentar fazer muitas alterações enquanto o código não está compilando ou os testes estão falhando, corremos o risco de acumular e agravar os problemas.
* Nos manter fiéis à essa abordagem nos obriga a escrever pequenos testes, o que significa pequenas alterações, o que nos ajuda a continuar trabalhando em sistemas complexos de forma gerenciável.
