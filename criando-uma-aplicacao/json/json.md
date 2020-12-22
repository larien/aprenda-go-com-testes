# JSON, roteamento e embedding

[**Você pode encontrar todo o código para este capítulo aqui**](https://github.com/quii/learn-go-with-tests/tree/master/criando-uma-aplicacao/json)

[No capítulo anterior](../servidor-http/servidor-http.md) nós criamos um web server para armazenar quantos jogos nossos jogadores venceram.

Nossa gerente de produtos veio com um novo requisito;  criar um novo endpoint chamado `/league` que retorne uma lista contendo todos os jogadores armazenados. Ela gostaria que isto fosse retornado como um JSON. 

## Este é o código que temos até agora

```go
// server.go
package main

import (
    "fmt"
    "net/http"
)

type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
}

type PlayerServer struct {
    store PlayerStore
}

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

```go
// InMemoryPlayerStore.go
package main

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
    return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
    store map[string]int
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
    i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
    return i.store[name]
}
```

```go
// main.go
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

Você pode encontrar os testes correspondentes no link no topo do capítulo.

Nós vamos começar criando o endpoint para a tabela de `league`.

## Escreva os testes primeiro

Ampliaremos a suite de testes existente, pois temos algumas funções de teste úteis e um `PlayerStore` falso para usar.

```go
// server_test.go

func TestLeague(t *testing.T) {
    store := StubPlayerStore{}
    server := &PlayerServer{&store}

    t.Run("it returns 200 on /league", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodGet, "/league", nil)
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusOK)
    })
}
```
Antes de nos preocuparmos sobre as pontuações atuais e o JSON, nós vamos tentar manter as mudanças pequenas com o plano de ir passo a passo rumo ao nosso objetivo. O início mais simples é checar se nós conseguimos consultar `/league` e obter um `OK` de retorno. 

## Tente rodar os testes

```text
=== RUN   TestLeague/it_returns_200_on_/league
panic: runtime error: slice bounds out of range [recovered]
    panic: runtime error: slice bounds out of range

goroutine 6 [running]:
testing.tRunner.func1(0xc42010c3c0)
    /usr/local/Cellar/go/1.10/libexec/src/testing/testing.go:742 +0x29d
panic(0x1274d60, 0x1438240)
    /usr/local/Cellar/go/1.10/libexec/src/runtime/panic.go:505 +0x229
github.com/quii/learn-go-with-tests/json-and-io/v2.(*PlayerServer).ServeHTTP(0xc420048d30, 0x12fc1c0, 0xc420010940, 0xc420116000)
    /Users/quii/go/src/github.com/quii/learn-go-with-tests/json-and-io/v2/server.go:20 +0xec
```

Seu `PlayerServer` deve estar sendo abortado por um panic como acima. Vá para a linha de código que está apontando para `server.go` no stack trace.  

```go
player := r.URL.Path[len("/players/"):]
```

No capítulo anterior, nós mencionamos que isto era uma maneira bastante ingênua de fazer o nosso roteamento. O que está acontecendo é que ele está tentando cortar a string do caminho da URL começando do índice após `/league` e então, isto nos dá um `slice bounds out of range`.

## Escreva somente o código suficiente para fazê-lo passar

Go têm um mecanismo de rotas nativo (built-in) chamado [`ServeMux`](https://golang.org/pkg/net/http/#ServeMux) \(request multiplexer\) que nos permite atracar um `http.Handler` para caminhos de uma requisição em específico.

Vamos cometer alguns pecados e obter os testes passando da maneira mais rápida que pudermos, sabendo que nós podemos refatorar isto com segurança uma vez que nós soubermos que os testes estão passando.

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    router := http.NewServeMux()

    router.Handle("/league", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    router.Handle("/players/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        player := r.URL.Path[len("/players/"):]

        switch r.Method {
        case http.MethodPost:
            p.processWin(w, player)
        case http.MethodGet:
            p.showScore(w, player)
        }
    }))

    router.ServeHTTP(w, r)
}
```

* Quando a requisição começa nós criamos um router e então dizemos para o caminho `x` usar o handler `y`.
* Então para nosso novo endpoint, nós usamos `http.HandlerFunc` e uma _função anônima_ para `w.WriteHeader(http.StatusOK)` quando `/league` é requisitada para fazer nosso novo teste passar.
* Para a rota `/players/` nós somente recortamos e colamos nosso código dentro de outro `http.HandlerFunc`.
* Finalmente, nós lidamos com a requisição que está vindo chamando nosso novo router `ServeHTTP` \(notou como `ServeMux` é _também_ um `http.Handler`?\)

## Refatorando

`ServeHTTP` parece um pouco grande, nós podemos separar as coisas um pouco refatorando nossos handlers em métodos separados.

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    router := http.NewServeMux()
    router.Handle("/league", http.HandlerFunc(p.leagueHandler))
    router.Handle("/players/", http.HandlerFunc(p.playersHandler))

    router.ServeHTTP(w, r)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    switch r.Method {
    case http.MethodPost:
        p.processWin(w, player)
    case http.MethodGet:
        p.showScore(w, player)
    }
}
```

É um pouco estranho \(e ineficiente\) estar configurando um router quando uma requisição chegar e então chamá-lo. O que idealmente queremos fazer é uma função do tipo `NewPlayerServer` que pegará nossas dependências e ao ser chamada, irá fazer a configuração única da criação do router. Desta forma, cada requisição pode usar somente uma instância do nosso router.

```go
type PlayerServer struct {
    store  PlayerStore
    router *http.ServeMux
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
    p := &PlayerServer{
        store,
        http.NewServeMux(),
    }

    p.router.Handle("/league", http.HandlerFunc(p.leagueHandler))
    p.router.Handle("/players/", http.HandlerFunc(p.playersHandler))

    return p
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    p.router.ServeHTTP(w, r)
}
```

* `PlayerServer` agora precisa armazenar um roteador.
* Nós movemos a criação do roteador para fora de `ServeHTTP` e colocamos dentro do nosso `NewPlayerServer`, então isto só será feito uma vez, não por requisição.
* Você vai precisar atualizar todos os testes e código de produção onde nós costumávamos fazer `PlayerServer{&store}` por `NewPlayerServer(&store)`.

### Uma refatoração final

Tente mudar o código para o seguinte:

```go
type PlayerServer struct {
    store  PlayerStore
    http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
    p := new(PlayerServer)

    p.store = store

    router := http.NewServeMux()
    router.Handle("/league", http.HandlerFunc(p.leagueHandler))
    router.Handle("/players/", http.HandlerFunc(p.playersHandler))

    p.Handler = router

    return p
}
```

Finalmente, se certifique de que você **deletou** `func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request)` por não ser mais necessária!

## Incorporando

Nós mudamos a segunda propriedade de `PlayerServer` removendo a propriedade nomeada `router http.ServeMux` e substituindo por `http.Handler`; isto é chamado de _incorporar_. 


> O Go não provê a noção típica de subclasses orientada por tipo, mas tem a habilidade de "emprestar" partes de uma implementação por incorporar tipos dentro de uma struct ou interface.

[Effective Go - Embedding](https://golang.org/doc/effective_go.html#embedding)

O que isto quer dizer é que nosso `PlayerServer` agora tem todos os métodos que `http.Handler` têm, que é somente o `ServeHTTP`.

Para "preencher" o `http.Handler` nós atribuímos ele para o `router` que nós criamos em `NewPlayerServer`. Nós podemos fazer isso porque `http.ServeMux` tem o método `ServeHTTP`.

Isto nos permite remover nosso próprio método `ServeHTTP`, pois nós já estamos expondo um via o tipo incorporado. 

Incorporamento é um recurso muito interessante da linguagem. Você pode usar isto com interfaces para compor novas interfaces.

```go
type Animal interface {
    Eater
    Sleeper
}
```

E você pode usar isto com tipos concretos também, não somente interfaces. Como você pode esperar, se você incorporar um tipo concreto você vai ter acesso a todos os seus métodos e campos públicos. 

### Alguma desvantagem?

Você deve ter cuidado ao incorporar tipos porque você vai expor todos os métodos e campos públicos do tipo que você incorporou. Em nosso caso, está tudo bem porque nós haviamos incorporado apenas a _interface_ que nós queremos expôr \(`http.Handler`\).

Se nós tivéssemos sido "preguiçosos" e incorporado `http.ServeMux` \(o tipo concreto\) por exemplo, também funcionaria _porém_ os usuários de `PlayerServer` seriam capazes de adicionar novas rotas ao nosso servidor porque o método `Handle(path, handler)` seria público.

**Quando incorporamos tipos, realmente devemos pensar sobre qual o impacto que isto terá em nossa API pública**

Isto é um erro _muito_ comum de mau uso de incorporamento, que termina poluindo nossas APIs e expondo os métodos internos dos seus tipos incorporados.

Agora que nós reestruturamos nossa aplicação, nós podemos facilmente adicionar novas rotas e botar para funcionar nosso endpoint `/league`. Agora precisamos fazê-lo retornar algumas informações úteis.

Nós devemos retornar um JSON semelhante a este:

```javascript
[
   {
      "Name":"Bill",
      "Wins":10
   },
   {
      "Name":"Alice",
      "Wins":15
   }
]
```

## Escreva o teste primeiro

Nós vamos começar tentando analizar a resposta dentro de algo mais significativo.

```go
func TestLeague(t *testing.T) {
    store := StubPlayerStore{}
    server := NewPlayerServer(&store)

    t.Run("it returns 200 on /league", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodGet, "/league", nil)
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        var got []Player

        err := json.NewDecoder(response.Body).Decode(&got)

        if err != nil {
            t.Fatalf ("Unable to parse response from server '%s' into slice of Player, '%v'", response.Body, err)
        }

        assertStatus(t, response.Code, http.StatusOK)
    })
}
```
### Por quê não testar o JSON como texto puro?

Você pode argumentar que um simples teste inicial poderia só comparar que o corpo da resposta tem um particular texto em JSON.


Na minha experiência, testes que comparam JSONs de forma literal possuem os seguintes problemas:

* _Fragilidade_. Se você mudar o modelo dos dados seu teste irá falhar.
* _Difícil de debugar_. Pode ser complicado de entender qual é o problema real ao se comparar dois textos JSON.
* _Má intenção_. Embora a saída deva ser JSON, o que é realmente importante é exatamente o que o dado é, ao invés de como ele está codificado.
* _Re-testando a biblioteca padrão_. Não há a necessidade de testar como a biblioteca padrão gera JSON, ela já está testada. Não teste o código de outras pessoas.

Ao invés disso, nós poderíamos analisar o JSON dentro de estruturas de dados que são relevantes para nós e nossos testes.

### Modelagem de dados

Dado o modelo de dados do JSON, parece que nós precisamos de uma lista de `Jogador` com alguns campos, sendo assim nós criaremos um novo tipo para capturarmos isso.

```go
type Jogador struct {
    Nome string
    Vitorias int
}
```
### Decodificação de JSON

```go
var got []Jogador
err := json.NewDecoder(response.Body).Decode(&got)
```

Para analizar o JSON dentro de nosso modelo de dados nós criamos um `Decoder` do pacote `encoding/json` e então chamamos seu método `Decode`. Para criar um `Decoder` é necessário ler de um `io.Reader`, que em nosso caso é nossa própria resposta `Body`.

`Decode` pega o endereço da coisa que nós estamos tentando decodificar, e é por isso que nós declaramos um slice vazio de `Jogador` na linha anterior.

Esse processo de analisar um JSON pode falhar, então `Decode` pode retornar um `error`. Não há ponto de continuidade para o teste se isto acontecer, então nós checamos o erro e paramos o teste com `t.Fatalf`.
Note que nós exibimos o corpo da resposta junto do erro, pois é importante para qualquer outra pessoa que esteja rodando os testes ver que o texto não pôde ser analisado.

## Tente rodar o teste

```text
=== RUN   TestLeague/it_returns_200_on_/league
    --- FAIL: TestLeague/it_returns_200_on_/league (0.00s)
        server_test.go:107: Unable to parse response from server '' into slice of Jogador, 'unexpected end of JSON input'
```

Nosso endpoint atualmente não retorna um corpo, então isso não pode ser analisado como JSON.

## Escreva código suficiente para fazê-lo passar.

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    leagueTable := []Player{
        {"Chris", 20},
    }

    json.NewEncoder(w).Encode(leagueTable)

    w.WriteHeader(http.StatusOK)
}
```
Os testes agora passam.

### Codificando e decodificando
Note a amável simetria na biblioteca padrão.

* Para criar um `Encoder` você precisa de um `io.Writer` que é o que `http.ResponseWriter` implementa.
* Para criar um `Decoder` você precisa de um `io.Reader` que o campo `Body` da nossa resposta implementa.

Ao longo deste livro, nós temos usado `io.Writer`. Isso é uma outra demonstração desta prevalência nas bibliotecas padrões e de como várias bibliotecas facilmente trabalham em conjunto com elas.

## Refatoração

Seria legal introduzir uma separação de conceitos entre nosso handler e o trecho de obter o `leagueTable`. Como sabemos, nós não vamos codificar isso por agora.

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(p.getLeagueTable())
    w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) getLeagueTable() []Player{
    return []Player{
        {"Chris", 20},
    }
}
```

Mais adiante, nós vamos querer estender nossos testes para então podermos controlar exatamente qual dado nós queremos receber de volta.

## Escreva o teste primeiro

Nós podemos atualizar o teste para afirmar que a tabela das ligas contem alguns jogadores que nós vamos pôr em nossa loja.

Atualize `StubPlayerStore` para permitir que ele armazene uma liga, que é apenas um slice de `Jogador`. Nós vamos armazenar nossos dados esperados lá.

```go
type StubPlayerStore struct {
    scores   map[string]int
    winCalls []string
    league []Player
}
```
Adiante, atualize nossos testes colocando alguns jogadores na propriedade da liga, para então afirmar que eles foram retornados do nosso servidor.

```go
func TestLeague(t *testing.T) {

    t.Run("it returns the league table as JSON", func(t *testing.T) {
        wantedLeague := []Player{
            {"Cleo", 32},
            {"Chris", 20},
            {"Tiest", 14},
        }

        store := StubPlayerStore{nil, nil, wantedLeague}
        server := NewPlayerServer(&store)

        request, _ := http.NewRequest(http.MethodGet, "/league", nil)
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        var got []Player

        err := json.NewDecoder(response.Body).Decode(&got)

        if err != nil {
            t.Fatalf("Unable to parse response from server '%s' into slice of Player, '%v'", response.Body, err)
        }

        assertStatus(t, response.Code, http.StatusOK)

        if !reflect.DeepEqual(got, wantedLeague) {
            t.Errorf("got %v want %v", got, wantedLeague)
        }
    })
}
```

## Tente rodar o teste

```text
./server_test.go:33:3: too few values in struct initializer
./server_test.go:70:3: too few values in struct initializer
```
## Escreva o minimo de código para que o teste rode e cheque as falhas na saída dele.

Você vai precisar atualizar os outros testes, assim como nós temos um novo campo em `StubPlayerStore`; ponha-o como nulo para os outros testes.

Tente executar os testes novamente e você deverá ter:

```text
=== RUN   TestLeague/it_returns_the_league_table_as_JSON
    --- FAIL: TestLeague/it_returns_the_league_table_as_JSON (0.00s)
        server_test.go:124: got [{Chris 20}] want [{Cleo 32} {Chris 20} {Tiest 14}]
```

## Escreva código suficiente para fazê-lo passar

Nós sabemos que o dado está em nosso `StubPlayerStore` e nós abstraímos esses dados para uma interface `PlayerStore`. Nós precisamos atualizar isto então qualquer um passando-nos um `PlayerStore` pode prover-nos com dados para as ligas.

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
    GetLeague() []Player
}
```

Agora nós podemos atualizar o código do nosso handler para chamar isto ao invés de retornar uma lista manualmente escrita. Delete nosso método `getLeagueTable()` e então atualize `leagueHandler` para chamar `GetLeague()`.

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(p.store.GetLeague())
    w.WriteHeader(http.StatusOK)
}
```

Tente executar os testes:

```text
# github.com/quii/learn-go-with-tests/json-and-io/v4
./main.go:9:50: cannot use NewInMemoryPlayerStore() (type *InMemoryPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *InMemoryPlayerStore does not implement PlayerStore (missing GetLeague method)
./server_integration_test.go:11:27: cannot use store (type *InMemoryPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *InMemoryPlayerStore does not implement PlayerStore (missing GetLeague method)
./server_test.go:36:28: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *StubPlayerStore does not implement PlayerStore (missing GetLeague method)
./server_test.go:74:28: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *StubPlayerStore does not implement PlayerStore (missing GetLeague method)
./server_test.go:106:29: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *StubPlayerStore does not implement PlayerStore (missing GetLeague method)
```

O compilador está reclamando porque `InMemoryPlayerStore` e `StubPlayerStore` não tem os novos métodos que nós adicionamos em nossa interface.

Para `StubPlayerStore` isto é bem fácil, apenas retorne o campo `league` que nós adicionamos anteriormente.

```go
func (s *StubPlayerStore) GetLeague() []Player {
    return s.league
}
```
Aqui está uma lembrança de como `InMemoryStore` é implementado:

```go
type InMemoryPlayerStore struct {
    store map[string]int
}
```
Embora seja bastante simples para implementar `GetLeague` "propriamente", iterando sobre o map, lembre que nós estamos apenas tentando _escrever o mínimo de código para fazer os testes passarem_.

Então vamos apenas deixar o compilador feliz por enquanto e viver com o desconfortável sentimento de uma implementação incompleta em nosso `InMemoryStore`.

```go
func (i *InMemoryPlayerStore) GetLeague() []Player {
    return nil
}
```

O que isto está realmente nos dizendo é que _depois_ nós vamos querer testar isto, porém vamos estacionar isto por hora.

Tente executar os testes, o compilador deve passar e os testes deverão estar passando!

## Refatoração

O código de teste não transmite suas intenções muito bem e possui vários trechos que podem ser refatorados.

```go
t.Run("it returns the league table as JSON", func(t *testing.T) {
    wantedLeague := []Player{
        {"Cleo", 32},
        {"Chris", 20},
        {"Tiest", 14},
    }

    store := StubPlayerStore{nil, nil, wantedLeague}
    server := NewPlayerServer(&store)

    request := newLeagueRequest()
    response := httptest.NewRecorder()

    server.ServeHTTP(response, request)

    got := getLeagueFromResponse(t, response.Body)
    assertStatus(t, response.Code, http.StatusOK)
    assertLeague(t, got, wantedLeague)
})
```

Aqui estão os novos helpers:

```go
func getLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
    t.Helper()
    err := json.NewDecoder(body).Decode(&league)

    if err != nil {
        t.Fatalf("Unable to parse response from server '%s' into slice of Player, '%v'", body, err)
    }

    return
}

func assertLeague(t *testing.T, got, want []Player) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}

func newLeagueRequest() *http.Request {
    req, _ := http.NewRequest(http.MethodGet, "/league", nil)
    return req
}
```

Uma última coisa que nós precisamos fazer para nosso servidor funcionar é ter certeza de que nós retornamos um `content-type` correto na resposta, então as máquinas podem reconhecer que nós estamos retornando um `JSON`.

## Escreva os testes primeiro

Adicione essa afirmação no teste existente

```go
if response.Result().Header.Get("content-type") != "application/json" {
    t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
}
```

## Tente rodar o teste

```text
=== RUN   TestLeague/it_returns_the_league_table_as_JSON
    --- FAIL: TestLeague/it_returns_the_league_table_as_JSON (0.00s)
        server_test.go:124: response did not have content-type of application/json, got map[Content-Type:[text/plain; charset=utf-8]]
```

## Escreva código suficiente para fazê-lo passar

Atualize `leagueHandler`

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    json.NewEncoder(w).Encode(p.store.GetLeague())
}
```

O teste deve passar.

## Refatoração

Adicione um helper para `assertContentType`.

```go
const jsonContentType = "application/json"

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
    t.Helper()
    if response.Result().Header.Get("content-type") != want {
        t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
    }
}
```

Use isso no teste.

```go
assertContentType(t, response, jsonContentType)
```

Agora que nós resolvemos `PlayerServer`, por agora podemos mudar nossa atenção para `InMemoryPlayerStore` porque no momento se nós tentarmos demonstrá-lo para o gerente de produto, `/league` não vai funcionar.

A forma mais rápida de nós termos alguma confiança é adicionar a nosso teste de integração, nós podemos bater no novo endpoint e checar se nós recebemos a resposta correta de `/league`.

## Escreva o teste primeiro

Nós podemos usar `t.Run` para parar este teste um pouco e então reusar os helpers dos testes do nosso servidor - novamente mostrando a importância de refatoração dos testes.

```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
    store := NewInMemoryPlayerStore()
    server := NewPlayerServer(store)
    player := "Pepper"

    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

    t.Run("get score", func(t *testing.T) {
        response := httptest.NewRecorder()
        server.ServeHTTP(response, newGetScoreRequest(player))
        assertStatus(t, response.Code, http.StatusOK)

        assertResponseBody(t, response.Body.String(), "3")
    })

    t.Run("get league", func(t *testing.T) {
        response := httptest.NewRecorder()
        server.ServeHTTP(response, newLeagueRequest())
        assertStatus(t, response.Code, http.StatusOK)

        got := getLeagueFromResponse(t, response.Body)
        want := []Player{
            {"Pepper", 3},
        }
        assertLeague(t, got, want)
    })
}
```

## Tente rodar o teste

```text
=== RUN   TestRecordingWinsAndRetrievingThem/get_league
    --- FAIL: TestRecordingWinsAndRetrievingThem/get_league (0.00s)
        server_integration_test.go:35: got [] want [{Pepper 3}]
```

## Escreva código suficiente para fazê-lo passar

`InMemoryPlayerStore` is returning `nil` when you call `GetLeague()` so we'll need to fix that.

```go
func (i *InMemoryPlayerStore) GetLeague() []Player {
    var league []Player
    for name, wins := range i.store {
        league = append(league, Player{name, wins})
    }
    return league
}
```

Tudo que nós precisamos fazer é iterar através do map e converter cada chave/valor para um `Jogador`

O teste deve passar agora.

## Concluindo

Nós temos continuado a seguramente iterar no nosso programa usando TDD, fazendo ele suportar novos endpoints de uma forma manutenível com um roteador e isso pode agora retornar JSON para nossos consumidores. No próximo capítulo, nós vamos cobrir persistência de dados e ordenação de nossas ligas.

O que nós cobrimos:

* **Roteamento**. A biblioteca padrão oferece uma fácil forma de usar tipos para fazer roteamento. Ela abraça completamente a interface `http.Handler` nela, tanto que você pode atribuir rotas para `Handler`s e a rota em si também é um `Handler`. Ela não tem alguns recursos que você pode esperar, como caminhos para variáveis \(ex. `/users/{id}`\). Você pode facilmente analisar esta informação por si mesmo porém você pode querer considerar olhar para outras bibliotecas de roteamento se isso se tornar um fardo. Muitas das mais populares seguem a filosofia das bibliotecas padrões e também implementam `http.Handler`.

* **Composição**. Nós tocamos um pouco nesta técnica porém você pode [ler mais sobre isso de Effective Go](https://golang.org/doc/effective_go.html#embedding). Se há uma coisa que você deve tirar disso é que composições podem ser extremamente úteis, porém _sempre pensando na sua API pública, só exponha o que é apropriado_.
* **Serialização e Desserialização de JSON**. A biblioteca padrão faz isto de forma bastante trivial ao serializar e desserializar nosso dado. Isto também abre para configurações e você pode customizar como esta transformação de dados funciona se necessário.
