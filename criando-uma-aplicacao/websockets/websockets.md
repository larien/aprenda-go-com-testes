# Websockets

[**Você pode encontrar todo o código para esse capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/criando-uma-aplicacao/websockets)

## Recapitulando o projeto

Nós temos duas aplicações no nosso código-base de poquer.

* _Aplicação de linha de comando_. Pede ao usuário para que insira o número
de jogadores. A partir daí informa os jogadores o valor da "aposta cega", que
aumenta em função do tempo. A qualquer momento, um usuário pode entrar com
`"{Jogador} ganhou"` para encerrar o jogo e salvar a vitória em um armazenamento.

* _Aplicação Web_. Permite que os usuários salvem os ganhadores e mostrem uma tabela
da liga. Divide o armazenamento com a aplicação de linha de comando.

## Próximos passos

A dona do produto está muito contente com a aplicação por linha de comando, mas
acharia melhor se conseguíssimos levar todas essas funcionalidades para o navegador.
Ela imagina uma página web com uma caixa de texto que permite que o usuário coloque
o número de jogadores e, após submeter esse dado, informe o valor da "aposta cega",
atualizando automaticamente quando for apropriado. Assim como a aplicação por linha
de comando, ela espera que o usuário possa declarar o vencedor e que isso faça com
que as devidas informações sejam salvas no banco de dados.

Descrevendo o projeto dessa forma parece bastante simples, mas sempre precisamos
enfatizar que devemos ter uma abordagem _iterativa_ pra desenvolver os nossos
programas.


Em primeiro lugar, vamos precisar apresentar um HTML. Até agora, todos os
nossos _endpoints_ HTTP retornaram texto puro ou JSON. Nós _poderíamos_ usar
as mesmas técnicas que conhecemos \(porque, no fim, tanto o texto puro quanto
o JSON são strings\), mas nós também podemos usar o pacote
[html/template](https://golang.org/pkg/html/template/) para uma solução mais
limpa.

Nós também temos que ser capazes de enviar mensagens assíncronas para o usuário
dizendo `A aposta blind é *y*` sem ter que recarregar o navegador. Para facilitar
isso, podemos usar [WebSockets](https://pt.wikipedia.org/wiki/WebSocket).

> WebSocket é uma tecnologia que permite a comunicação bidirecional por canais full-duplex
sobre um único socket TCP (Transmission Control Protocol)

Como estamos adotando várias técnicas, é ainda mais importante que façamos o menor
trabalho possível primeiro e só então iteramos.

Por causa disso, a primeira coisa que faremos é criar uma página web com um formulário
para o usuário salvar um vencedor. Em vez de usar um formulário simples, vamos usar os
WebSockets para enviar os dados para o nosso servidor o salvar.

Depois disso, iremos trabalhaor nos alertas cegos, uma vez que já teremos algum
código de infraestrutura pronto.

### E os testes para o JavaScript?

Haverá algum JavaScript escrito para cumprir nossa tarefa, mas não vamos
escrever testes para ele.

É claro que é possível, mas, em nome da breviedade, não incluíremos quaisquer
explicações para isso.

Desculpem, amigos. Peçam para a O'Reilly me pagar para fazer um "Aprenda
JavaScript com testes".

## Escreva o teste primeiro

A primeira coisa que precisamos fazer é montar algum HTML para os usuários
quando eles acessarem `/partida`.

Aqui está um lembrete do código no nosso servidor web:

```go
type ServidorJogador struct {
    armazenamento ArmazenamentoJogador
    http.Handler
}

const tipoConteudoJSON = "application/json"

func NovoServidorJogador(armazenamento ArmazenamentoJogador) *ServidorJogador {
    p := new(ServidorJogador)

    p.armazenamento = armazenamento

    roteador := http.NewServeMux()
    roteador.Handle("/liga", http.HandlerFunc(p.manipulaLiga))
    roteador.Handle("/jogadores/", http.HandlerFunc(p.manipulaJogadores))

    p.Handler = roteador

    return p
}
```

A maneira _mais fácil_ que podemos fazer por agora é checar que recebemos um
código `200` quando acessamos o `GET /partida`.

```go
func TestJogo(t *testing.T) {
    t.Run("GET /partida retorna 200", func(t *testing.T) {
        servidor := NovoServidorJogador(&EsbocoDeArmazenamentoJogador{})

        requisicao, _ := http.NewRequest(http.MethodGet, "/partida", nil)
        resposta := httptest.NewRecorder()

        servidor.ServeHTTP(resposta, requisicao)

        verificaStatus(t, resposta.Code, http.StatusOK)
    })
}
```

## Tente rodar o teste

```text
--- FAIL: TestJogo (0.00s)
=== RUN   TestJogo/GET_/game_returns_200
    --- FAIL: TestJogo/GET_/game_returns_200 (0.00s)
        server_test.go:109: não obteve o status correto, obtido 404, esperado 200
```

## Escreva código suficiente para fazer o teste passar

Our servidor has a roteador setup so it's relatively easy para fix.

To our roteador add

```go
roteador.Handle("/partida", http.HandlerFunc(p.partida))
```

E então escreva o método `partida`

```go
func (p *ServidorJogador) partida(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}
```

## Refatore

O servidor já está bem graças às inserções que fizemos no código já bem refatorado.

Podemos ajeitar ainda mais o teste um pouco ao adicionarmos uma função auxiliar
`novaRequisicaoDeJogo` para fazer a requisição para `/jogo`. Tente escrever essa
função você mesmo.

```go
func TestJogo(t *testing.T) {
    t.Run("GET /partida retorna 200", func(t *testing.T) {
        servidor := NovoServidorJogador(&EsbocoDeArmazenamentoJogador{})

        requisicao :=  novaRequisicaoJogo()
        resposta := httptest.NewRecorder()

        servidor.ServeHTTP(resposta, requisicao)

        verificaStatus(t, resposta, http.StatusOK)
    })
}
```

You'll also notice I changed `verificaStatus` para accept `resposta` rather than `resposta.Code` as I feel it reads better.

Now we need para make the endpoint return some HTML, here it is

```markup
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Vamos jogar pôquer</title>
</head>
<corpo>
<section id="partida">
    <div id="declare-vencedor">
        <label for="vencedor">Winner</label>
        <input type="text" id="vencedor"/>
        <button id="vencedor-button">Declare vencedor</button>
    </div>
</section>
</corpo>
<script type="application/javascript">

    const submitWinnerButton = document.getElementById('vencedor-button')
    const entradaVencedor = document.getElementById('vencedor')

    if (window['WebSocket']) {
        const conexão = new WebSocket('ws://' + document.location.host + '/ws')

        submitWinnerButton.onclick = event => {
            conexão.send(entradaVencedor.value)
        }
    }
</script>
</html>
```

Temos uma página web bem simples:

* Uma entrada de texto para a pessoa inserir a vitória
* Um botão onde pode-se clicar para declarar quem venceu
* Um pouco de JavaScript para abrir uma conexão WebSocket para nosso servidor e
assim gerenciar o envio dos dados ao pressionar o botão

`WebSocket` é integrado na maioria dos navegadores modernos, logo não precisamos
nos preocupar em instalar bibliotecas. A página web não vai funcionar em
navegadores antigos, mas para nosso caso tá tudo bem.

### How do we test we return the correct markup?

There are a few ways. As has been emphasised throughout the book, it is important that the tests you write have sufficient value para justify the cost.

1. Write a browser based test, using something like Selenium. These tests are the most "realistic" of all approaches because they start an actual web browser of some kind and simulates a user interacting with it. These tests can give you a lot of confidence your system works but are more difficult para write than unit tests and much slower para run. For the purposes of our product this is overkill.
2. Do an exact string match. This _can_ be ok but these kind of tests end up being very brittle. The moment someone changes the markup you will have a test failing when in practice nothing has _actually broken_.
3. Check we call the correct template. We will be using a templating library from the standard lib para serve the HTML \(discussed shortly\) and we could inject in the _thing_ para generate the HTML and spy on its call para check we're doing it right. This would have an impact on our code's design but doesn't actually test a great deal; other than we're calling it with the correct template arquivo. Given we will only have the one template in our project the chance of failure here seems low.

So in the book "Learn Go with Tests" for the first time, we're not going para write a test.

Put the markup in a arquivo called `partida.html`

Next change the endpoint we just wrote para the following

```go
func (p *ServidorJogador) partida(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("partida.html")

    if err != nil {
        http.Error(w, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, nil)
}
```

[`html/template`](https://golang.org/pkg/html/template/) is a Go package for creating HTML. In our case we call `template.ParseFiles`, giving the path of our html arquivo. Assuming there is no error you can then `Execute` the template, which writes it para an `io.Writer`. In our case we esperado it para `Write` para the internet, so we give it our `http.ResponseWriter`.

As we have not written a test, it would be prudent para manually test our web servidor just para make sure things are working as we'd hope. Go para `cmd/webserver` and run the `main.go` arquivo. Visit `http://localhost:5000/partida`.

You _should_ have obtido an error about not being able para find the template. You can either change the path para be relative para your folder, or you can have a copy of the `partida.html` in the `cmd/webserver` directory. I chose para create a symlink \(`ln -s ../../partida.html partida.html`\) para the arquivo inside the root of the project so if I make changes they are reflected when running the servidor.

If you make this change and run again you should see our UI.

Now we need para test that when we obtera string over a WebSocket connection para our servidor that we declare it as a vencedor of a partida.

## Write the test first

For the first time we are going para use an external library so that we can work with WebSockets.

Run `go obtergithub.com/gorilla/websocket`

This will fetch the code for the excellent [Gorilla WebSocket](https://github.com/gorilla/websocket) library. Now we can update our tests for our new requirement.

```go
t.Run("quando recebemos uma mensagem de um websocket que é vencedor da partida", func(t *testing.T) {
    armazenamento := &EsbocoDeArmazenamentoJogador{}
    vencedor := "Ruth"
    servidor := httptest.NewServer(NovoServidorJogador(armazenamento))
    defer servidor.Close()

    wsURL := "ws" + strings.TrimPrefix(servidor.URL, "http") + "/ws"

    ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("não foi possível abrir uma conexão de websocket em %s %v", wsURL, err)
    }
    defer ws.Close()

    if err := ws.WriteMessage(websocket.TextMessage, []byte(vencedor)); err != nil {
        t.Fatalf("não foi possível enviar mensagem na conexão websocket %v", err)
    }

    VerificaVitoriaDoVencedor(t, armazenamento, vencedor)
})
```

Make sure that you have an import for the `websocket` library. My IDE automatically did it for me, so should yours.

To test what happens from the browser we have para open up our own WebSocket connection and write para it.

Our previous tests around our servidor just called methods on our servidor but now we need para have a persistent connection para our servidor. To do that we use `httptest.NewServer` which takes a `http.Handler` and will spin it up and listen for connections.

Using `websocket.DefaultDialer.Dial` we try para dial in para our servidor and then we'll try and send a mensagem with our `vencedor`.

Finally we assert on the jogador armazenamento para check the vencedor was recorded.

## Try para run the test

```text
=== RUN   TestJogo/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game
    --- FAIL: TestJogo/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game (0.00s)
        server_test.go:124: não foi possível abrir uma conexão de websocket em ws://127.0.0.1:55838/ws websocket: bad handshake
```

We have not changed our servidor para accept WebSocket connections on `/ws` so we're not shaking hands yet.

## Write enough code para make it pass

Add another listing para our roteador

```go
roteador.Handle("/ws", http.HandlerFunc(p.webSocket))
```

Then add our new `webSocket` handler

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }
    upgrader.Upgrade(w, r, nil)
}
```

To accept a WebSocket connection we `Upgrade` the requisicao. If you now re-run the test you should move on para the next error.

```text
=== RUN   TestJogo/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game
    --- FAIL: TestJogo/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game (0.00s)
        server_test.go:132: obtido 0 chamadas paraGravarVitoria esperado 1
```

Now that we have a connection opened, we'll esperado para listen for a mensagem and then record it as the vencedor.

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }
    conexão, _ := upgrader.Upgrade(w, r, nil)
    _, winnerMsg, _ := conexão.ReadMessage()
    p.armazenamento.GravarVitoria(string(winnerMsg))
}
```

\(Yes, we're ignoring a lot of errors right now!\)

`conexão.ReadMessage()` blocks on waiting for a mensagem on the connection. Once we obterone we use it para `GravarVitoria`. This would finally close the WebSocket connection.

If you try and run the test, it's still failing.

The issue is timing. There is a delay between our WebSocket connection reading the mensagem and recording the win and our test finishes before it happens. You can test this by putting a short `time.Sleep` before the final assertion.

Let's go with that for now but acknowledge that putting in arbitrary sleeps into tests **is very bad practice**.

```go
time.Sleep(10 * time.Millisecond)
VerificaVitoriaDoVencedor(t, armazenamento, vencedor)
```

## Refactor

We committed many sins para make this test work both in the servidor code and the test code but remember this is the easiest way for us para work.

We have nasty, horrible, _working_ software backed by a test, so now we are free para make it nice and know we wont break anything accidentally.

Let's start with the servidor code.

We can move the `upgrader` para a private value inside our package because we don't need para redeclare it on every WebSocket connection requisicao

```go
var atualizadorDeWebsocket = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    conexão, _ := atualizadorDeWebsocket.Upgrade(w, r, nil)
    _, winnerMsg, _ := conexão.ReadMessage()
    p.armazenamento.GravarVitoria(string(winnerMsg))
}
```

Our call para `template.ParseFiles("partida.html")` will run on every `GET /partida` which means we'll go para the arquivo system on every requisicao even though we have no need para re-parse the template. Let's refactor our code so that we parse the template once in `NovoServidorJogador` instead. We'll have para make it so this function can now return an error in case we have problems fetching the template from disk or parsing it.

Here's the relevant changes para `ServidorJogador`

```go
type ServidorJogador struct {
    armazenamento ArmazenamentoJogador
    http.Handler
    template *template.Template
}

const htmlTemplatePath = "partida.html"

func NovoServidorJogador(armazenamento ArmazenamentoJogador) (*ServidorJogador, error) {
    p := new(ServidorJogador)

    tmpl, err := template.ParseFiles("partida.html")

    if err != nil {
        return nil, fmt.Errorf("problema ao abrir %s %v", htmlTemplatePath, err)
    }

    p.template = tmpl
    p.armazenamento = armazenamento

    roteador := http.NewServeMux()
    roteador.Handle("/liga", http.HandlerFunc(p.manipulaLiga))
    roteador.Handle("/jogadores/", http.HandlerFunc(p.manipulaJogadores))
    roteador.Handle("/partida", http.HandlerFunc(p.partida))
    roteador.Handle("/ws", http.HandlerFunc(p.webSocket))

    p.Handler = roteador

    return p, nil
}

func (p *ServidorJogador) partida(w http.ResponseWriter, r *http.Request) {
    p.template.Execute(w, nil)
}
```

By changing the signature of `NovoServidorJogador` we now have compilation problems. Try and fix them yourself or refer para the source code if you struggle.

For the test code I made a helper called `deveFazerServidorJogador(t *testing.T, armazenamento ArmazenamentoJogador) *ServidorJogador` so that I could hide the error noise away from the tests.

```go
func deveFazerServidorJogador(t *testing.T, armazenamento ArmazenamentoJogador) *ServidorJogador {
    servidor, err := NovoServidorJogador(armazenamento)
    if err != nil {
        t.Fatal("problema ao criar o servidor do jogador", err)
    }
    return servidor
}
```

Similarly I created another helper `mustDialWS` so that I could hide nasty error noise when creating the WebSocket connection.

```go
func mustDialWS(t *testing.T, url string) *websocket.Conn {
    ws, _, err := websocket.DefaultDialer.Dial(url, nil)

    if err != nil {
        t.Fatalf("não foi possível abrir uma conexão de websocket em %s %v", url, err)
    }

    return ws
}
```

Finally in our test code we can create a helper para tidy up sending mensagens

```go
func escreverMensagemNoWebsocket(t *testing.T, conexão *websocket.Conn, mensagem string) {
    t.Helper()
    if err := conexão.WriteMessage(websocket.TextMessage, []byte(mensagem)); err != nil {
        t.Fatalf("não foi possível enviar mensagem na conexão websocket %v", err)
    }
}
```

Now the tests are passing try running the servidor and declare some winners in `/partida`. You should see them recorded in `/liga`. Remember that every time we obtera vencedor we _close the connection_, you will need para refresh the page para open the connection again.

We've made a trivial web form that lets users record the vencedor of a partida. Let's iterate on it para make it so the user can start a partida by providing a number of jogadores and the servidor will push mensagens para the client informing them of what the blind value is as time passes.

First of all update `partida.html` para update our client side code for the new requirements

```markup
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Lets play poquer</title>
</head>
<corpo>
<section id="partida">
    <div id="partida-start">
        <label for="jogador-count">Number of jogadores</label>
        <input type="number" id="jogador-count"/>
        <button id="start-partida">Começar</button>
    </div>

    <div id="declare-vencedor">
        <label for="vencedor">Winner</label>
        <input type="text" id="vencedor"/>
        <button id="vencedor-button">Declare vencedor</button>
    </div>

    <div id="blind-value"/>
</section>

<section id="partida-end">
    <h1>Another great partida of poquer everyone!</h1>
    <p><a href="/liga">Go check the liga table</a></p>
</section>

</corpo>
<script type="application/javascript">
    const startGame = document.getElementById('partida-start')

    const declareWinner = document.getElementById('declare-vencedor')
    const submitWinnerButton = document.getElementById('vencedor-button')
    const entradaVencedor = document.getElementById('vencedor')

    const blindContainer = document.getElementById('blind-value')

    const gameContainer = document.getElementById('partida')
    const gameEndContainer = document.getElementById('partida-end')

    declareWinner.hidden = true
    gameEndContainer.hidden = true

    document.getElementById('start-partida').addEventListener('click', event => {
        startGame.hidden = true
        declareWinner.hidden = false

        const numeroDeJogadores = document.getElementById('jogador-count').value

        if (window['WebSocket']) {
            const conexão = new WebSocket('ws://' + document.location.host + '/ws')

            submitWinnerButton.onclick = event => {
                conexão.send(entradaVencedor.value)
                gameEndContainer.hidden = false
                gameContainer.hidden = true
            }

            conexão.onclose = evt => {
                blindContainer.innerText = 'Connection closed'
            }

            conexão.onmessage = evt => {
                blindContainer.innerText = evt.data
            }

            conexão.onopen = function () {
                conexão.send(numeroDeJogadores)
            }
        }
    })
</script>
</html>
```

The main changes is bringing in a section para enter the number of jogadores and a section para display the blind value. We have a little logic para show/hide the user interface depending on the stage of the partida.

Any mensagem we receive via `conexão.onmessage` we assume para be blind alerts and so we set the `blindContainer.innerText` accordingly.

How do we go about sending the blind alerts? In the previous chapter we introduced the idea of `Jogo` so our CLI code could call a `Jogo` and everything else would be taken care of including scheduling blind alerts. This turned out para be a good separation of concern.

```go
type Jogo interface {
    Começar(numeroDeJogadores int)
    Terminar(vencedor string)
}
```

When the user was prompted in the CLI for number of jogadores it would `Começar` the partida which would kick off the blind alerts and when the user declared the vencedor they would `Terminar`. This is the same requirements we have now, just a different way of getting the inputs; so we should look para re-use this concept if we can.

Our "real" implementation of `Jogo` is `TexasHoldem`

```go
type TexasHoldem struct {
    alerter AlertadorDeBlind
    armazenamento   ArmazenamentoJogador
}
```

By sending in a `AlertadorDeBlind` `TexasHoldem` can schedule blind alerts para be sent para _wherever_

```go
type AlertadorDeBlind interface {
    AgendarAlertaPara(duracao time.Duration, quantia int)
}
```

And as a reminder, here is our implementation of the `AlertadorDeBlind` we use in the CLI.

```go
func SaidaAlertador(duracao time.Duration, quantia int) {
    time.AfterFunc(duracao, func() {
        fmt.Fprintf(os.Stdout, "Blind agora é %d\n", quantia)
    })
}
```

This works in CLI because we _always esperado para send the alerts para `os.Stdout`_ but this wont work for our web servidor. For every requisicao we obtera new `http.ResponseWriter` which we then upgrade para `*websocket.Conn`. So we cant know when constructing our dependencies where our alerts need para go.

For that reason we need para change `AlertadorDeBlind.AgendarAlertaPara` so that it takes a destination for the alerts so that we can re-use it in our webserver.

Open AlertadorDeBlind.go and add the parameter `para io.Writer`

```go
type AlertadorDeBlind interface {
    AgendarAlertaPara(duracao time.Duration, quantia int, para io.Writer)
}

type AlertadorDeBlindFunc func(duracao time.Duration, quantia int, para io.Writer)

func (a AlertadorDeBlindFunc) AgendarAlertaPara(duracao time.Duration, quantia int, para io.Writer) {
    a(duracao, quantia, para)
}
```

The idea of a `StdoutAlerter` doesn't fit our new model so just rename it para `Alertador`

```go
func Alertador(duracao time.Duration, quantia int, para io.Writer) {
    time.AfterFunc(duracao, func() {
        fmt.Fprintf(para, "Blind agora é %d\n", quantia)
    })
}
```

If you try and compile, it will fail in `TexasHoldem` because it is calling `AgendarAlertaPara` without a destination, para obterthings compiling again _for now_ hard-code it para `os.Stdout`.

Try and run the tests and they will fail because `AlertadorDeBlindEspiao` no longer implementa `AlertadorDeBlind`, fix this by updating the signature of `AgendarAlertaPara`, run the tests and we should still be green.

It doesn't make any sense for `TexasHoldem` para know where para send blind alerts. Let's now update `Jogo` so that when you start a partida you declare _where_ the alerts should go.

```go
type Jogo interface {
    Começar(numeroDeJogadores int, destinoDosAlertas io.Writer)
    Terminar(vencedor string)
}
```

Let the compiler tell you what you need para fix. The change isn't so bad:

* Update `TexasHoldem` so it properly implementa `Jogo`
* In `CLI` when we start the partida, pass in our `out` property \(`cli.partida.Começar(numeroDeJogadores, cli.out)`\)
* In `TexasHoldem`'s test i use `partida.Começar(5, ioutil.Discard)` para fix the compilation problem and configure the alert output para be discarded

If you've obtido everything right, everything should be green! Now we can try and use `Jogo` within `Server`.

## Write the test first

The requirements of `CLI` and `Server` are the same! It's just the delivery mechanism is different.

Let's take a look at our `CLI` test for inspiration.

```go
t.Run("começa partida com 3 jogadores e termina partida com 'Chris' como vencedor", func(t *testing.T) {
    partida := &JogoEspiao{}

    out := &bytes.Buffer{}
    in := usuarioEnvia("3", "Chris venceu")

    poquer.NovaCLI(in, out, partida).JogarPoquer()

    verificaMensagensEnviadasParaUsuario(t, out, poquer.PromptJogador)
    verificaJogoComeçadoCom(t, partida, 3)
    verificaTerminosChamadosCom(t, partida, "Chris")
})
```

It looks like we should be able para test drive out a similar outcome using `JogoEspiao`

Replace the old websocket test with the following

```go
t.Run("start a partida with 3 jogadores and declare Ruth the vencedor", func(t *testing.T) {
    partida := &poquer.JogoEspiao{}
    vencedor := "Ruth"
    servidor := httptest.NewServer(deveFazerServidorJogador(t, ArmazenamentoJogadorTosco, partida))
    ws := mustDialWS(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

    defer servidor.Close()
    defer ws.Close()

    escreverMensagemNoWebsocket(t, ws, "3")
    escreverMensagemNoWebsocket(t, ws, vencedor)

    time.Sleep(10 * time.Millisecond)
    verificaJogoComeçadoCom(t, partida, 3)
    verificaTerminosChamadosCom(t, partida, vencedor)
})
```

* As discussed we create a spy `Jogo` and pass it into `deveFazerServidorJogador` \(be sure para update the helper para support this\).
* We then send the web socket mensagens for a partida.
* Finally we assert that the partida is started and finished with what we expect.

## Try para run the test

You'll have a number of compilation errors around `deveFazerServidorJogador` in other tests. Introduce an unexported variable `dummyGame` and use it through all the tests that aren't compiling

```go
var (
    dummyGame = &JogoEspiao{}
)
```

The final error is where we are trying para pass in `Jogo` para `NovoServidorJogador` but it doesn't support it yet

```text
./server_test.go:21:38: too many arguments in call para "github.com/larien/learn-go-with-tests/WebSockets/v2".NovoServidorJogador
    have ("github.com/larien/learn-go-with-tests/WebSockets/v2".ArmazenamentoJogador, "github.com/larien/learn-go-with-tests/WebSockets/v2".Jogo)
    esperado ("github.com/larien/learn-go-with-tests/WebSockets/v2".ArmazenamentoJogador)
```

## Write the minimal quantia of code for the test para run and check the failing test output

Just add it as an argument for now just para obterthe test running

```go
func NovoServidorJogador(armazenamento ArmazenamentoJogador, partida Jogo) (*ServidorJogador, error) {
```

Finally!

```text
=== RUN   TestJogo/start_a_game_with_3_players_and_declare_Ruth_the_winner
--- FAIL: TestJogo (0.01s)
    --- FAIL: TestJogo/start_a_game_with_3_players_and_declare_Ruth_the_winner (0.01s)
        server_test.go:146: wanted Começar called with 3 but obtido 0
        server_test.go:147: expected finish called with 'Ruth' but obtido ''
FAIL
```

## Write enough code para make it pass

We need para add `Jogo` as a field para `ServidorJogador` so that it can use it when it gets requests.

```go
type ServidorJogador struct {
    armazenamento ArmazenamentoJogador
    http.Handler
    template *template.Template
    partida Jogo
}
```

\(We already have a method called `partida` so rename that para `jogarJogo`\)

Next lets assign it in our constructor

```go
func NovoServidorJogador(armazenamento ArmazenamentoJogador, partida Jogo) (*ServidorJogador, error) {
    p := new(ServidorJogador)

    tmpl, err := template.ParseFiles("partida.html")

    if err != nil {
        return nil, fmt.Errorf("problema ao abrir %s %v", htmlTemplatePath, err)
    }

    p.partida = partida

    // etc
```

Now we can use our `Jogo` within `webSocket`.

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    conexão, _ := atualizadorDeWebsocket.Upgrade(w, r, nil)

    _, mensagemNumeroDeJogadores, _ := conexão.ReadMessage()
    numeroDeJogadores, _ := strconv.Atoi(string(mensagemNumeroDeJogadores))
    p.partida.Começar(numeroDeJogadores, ioutil.Discard) //todo: Dont discard the blinds mensagens!

    _, vencedor, _ := conexão.ReadMessage()
    p.partida.Terminar(string(vencedor))
}
```

Hooray! The tests pass.

We are not going para send the blind mensagens anywhere _just yet_ as we need para have a think about that. When we call `partida.Começar` we send in `ioutil.Discard` which will just discard any mensagens written para it.

For now start the web servidor up. You'll need para update the `main.go` para pass a `Jogo` para the `ServidorJogador`

```go
func main() {
    db, err := os.OpenFile(nomeArquivoBaseDeDados, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        log.Fatalf("problema ao abrir %s %v", nomeArquivoBaseDeDados, err)
    }

    armazenamento, err := poquer.NovoSistemaArquivoArmazenamentoJogador(db)

    if err != nil {
        log.Fatalf("problema ao criar sistema de arquivo de armazenamento do jogador, %v ", err)
    }

    partida := poquer.NovoTexasHoldem(poquer.AlertadorDeBlindFunc(poquer.Alertador), armazenamento)

    servidor, err := poquer.NovoServidorJogador(armazenamento, partida)

    if err != nil {
        log.Fatalf("problema ao criar o servidor do jogador %v", err)
    }

    if err := http.ListenAndServe(":5000", servidor); err != nil {
        log.Fatalf("não foi possível ouvir na porta 5000 %v", err)
    }
}
```

Discounting the fact we're not getting blind alerts yet, the app does work! We've managed para re-use `Jogo` with `ServidorJogador` and it has taken care of all the details. Once we figure out how para send our blind alerts through para the web sockets rather than discarding them it _should_ all work.

Before that though, let's tidy up some code.

## Refactor

The way we're using WebSockets is fairly basic and the error handling is fairly naive, so I wanted para encapsulate that in a type just para remove that messyness from the servidor code. We may wish para revisit it later but for now this'll tidy things up a bit

```go
type websocketServidorJogador struct {
    *websocket.Conn
}

func novoWebsocketServidorJogador(w http.ResponseWriter, r *http.Request) *websocketServidorJogador {
    conexão, err := atualizadorDeWebsocket.Upgrade(w, r, nil)

    if err != nil {
        log.Printf("houve um problema ao atualizar a conexão para WebSockets %v\n", err)
    }

    return &websocketServidorJogador{conexão}
}

func (w *websocketServidorJogador) EsperarPelaMensagem() string {
    _, msg, err := w.ReadMessage()
    if err != nil {
        log.Printf("erro ao ler do websocket %v\n", err)
    }
    return string(msg)
}
```

Now the servidor code is a bit simplified

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    ws := novoWebsocketServidorJogador(w, r)

    mensagemNumeroDeJogadores := ws.EsperarPelaMensagem()
    numeroDeJogadores, _ := strconv.Atoi(mensagemNumeroDeJogadores)
    p.partida.Começar(numeroDeJogadores, ioutil.Discard) //todo: Dont discard the blinds mensagens!

    vencedor := ws.EsperarPelaMensagem()
    p.partida.Terminar(vencedor)
}
```

Once we figure out how para not discard the blind mensagens we're done.

### Let's _not_ write a test!

Sometimes when we're not sure how para do something, it's best just para play around and try things out! Make sure your work is committed first because once we've figured out a way we should drive it through a test.

The problematic line of code we have is

```go
p.partida.Começar(numeroDeJogadores, ioutil.Discard) //todo: Dont discard the blinds mensagens!
```

We need para pass in an `io.Writer` for the partida para write the blind alerts para.

Wouldn't it be nice if we could pass in our `websocketServidorJogador` from before? It's our wrapper around our WebSocket so it _feels_ like we should be able para send that para our `Jogo` para send mensagens para.

Give it a go:

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    ws := novoWebsocketServidorJogador(w, r)

    mensagemNumeroDeJogadores := ws.EsperarPelaMensagem()
    numeroDeJogadores, _ := strconv.Atoi(mensagemNumeroDeJogadores)
    p.partida.Começar(numeroDeJogadores, ws)
    //etc...
```

The compiler complains

```text
./servidor.go:71:14: cannot use ws (type *websocketServidorJogador) as type io.Writer in argument para p.partida.Começar:
    *websocketServidorJogador does not implement io.Writer (missing Write method)
```

It seems the obvious thing para do, would be para make it so `websocketServidorJogador` _does_ implement `io.Writer`. To do so we use the underlying `*websocket.Conn` para use `WriteMessage` para send the mensagem down the websocket

```go
func (w *websocketServidorJogador) Write(p []byte) (n int, err error) {
    err = w.WriteMessage(1, p)

    if err != nil {
        return 0, err
    }

    return len(p), nil
}
```

This seems too easy! Try and run the application and see if it works.

Beforehand edit `TexasHoldem` so that the blind increment time is shorter so you can see it in action

```go
blindIncrement := time.Duration(5+numeroDeJogadores) * time.Second // (rather than a minute)
```

You should see it working! The blind quantia increments in the browser as if by magic.

Now let's revert the code and think how para test it. In order para _implement_ it all we did was pass through para `StartGame` was `websocketServidorJogador` rather than `ioutil.Discard` so that might make you think we should perhaps spy on the call para verify it works.

Spying is great and helps us check implementation details but we should always try and favour testing the _real_ behaviour if we can because when you decide para refactor it's often spy tests that start failing because they are usually checking implementation details that you're trying para change.

Our test currently opens a websocket connection para our running servidor and sends mensagens para make it do things. Equally we should be able para test the mensagens our servidor sends back over the websocket connection.

## Write the test first

We'll edit our existing test.

Currently our `JogoEspiao` does not send any data para `out` when you call `Começar`. We should change it so we can configure it para send a canned mensagem and then we can check that mensagem gets sent para the websocket. This should give us confidence that we have configured things correctly whilst still exercising the real behaviour we esperado.

```go
type JogoEspiao struct {
    ComecouASerChamado     bool
    ComecouASerChamadoCom int
    AlertaDeBlind      []byte

    TerminouDeSerChamado   bool
    TerminouDeSerChamadoCom string
}
```

Add `AlertaDeBlind` field.

Update `JogoEspiao` `Começar` para send the canned mensagem para `out`.

```go
func (j *JogoEspiao) Começar(numeroDeJogadores int, out io.Writer) {
    j.ComecouASerChamado = true
    j.ComecouASerChamadoCom = numeroDeJogadores
    out.Write(j.AlertaDeBlind)
}
```

This now means when we exercise `ServidorJogador` when it tries para `Começar` the partida it should end up sending mensagens through the websocket if things are working right.

Finally we can update the test

```go
t.Run("start a partida with 3 jogadores, send some blind alerts down WS and declare Ruth the vencedor", func(t *testing.T) {
    wantedBlindAlert := "Blind is 100"
    vencedor := "Ruth"

    partida := &JogoEspiao{AlertaDeBlind: []byte(wantedBlindAlert)}
    servidor := httptest.NewServer(deveFazerServidorJogador(t, ArmazenamentoJogadorTosco, partida))
    ws := mustDialWS(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

    defer servidor.Close()
    defer ws.Close()

    escreverMensagemNoWebsocket(t, ws, "3")
    escreverMensagemNoWebsocket(t, ws, vencedor)

    time.Sleep(10 * time.Millisecond)
    verificaJogoComeçadoCom(t, partida, 3)
    verificaTerminosChamadosCom(t, partida, vencedor)

    _, gotBlindAlert, _ := ws.ReadMessage()

    if string(gotBlindAlert) != wantedBlindAlert {
        t.Errorf("obtido blind alert '%s', esperado '%s'", string(gotBlindAlert), wantedBlindAlert)
    }
})
```

* We've added a `wantedBlindAlert` and configured our `JogoEspiao` para send it para `out` if `Começar` is called.
* We hope it gets sent in the websocket connection so we've added a call para `ws.ReadMessage()` para wait for a mensagem para be sent and then check it's the one we expected.

## Try para run the test

You should find the test hangs forever. This is because `ws.ReadMessage()` will block until it gets a mensagem, which it never will.

## Write the minimal quantia of code for the test para run and check the failing test output

We should never have tests that hang so let's introduce a way of handling code that we esperado para timeout.

```go
func within(t *testing.T, d time.Duration, assert func()) {
    t.Helper()

    done := make(chan struct{}, 1)

    go func() {
        assert()
        done <- struct{}{}
    }()

    select {
    case <-time.After(d):
        t.Error("timed out")
    case <-done:
    }
}
```

What `within` does is take a function `assert` as an argument and then runs it in a go routine. If/When the function finishes it will signal it is done via the `done` channel.

While that happens we use a `select` statement which lets us wait for a channel para send a mensagem. From here it is a race between the `assert` function and `time.After` which will send a signal when the duracao has occurred.

Finally I made a helper function for our assertion just para make things a bit neater

```go
func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, esperado string) {
    _, msg, _ := ws.ReadMessage()
    if string(msg) != esperado {
        t.Errorf(`obtido "%s", esperado "%s"`, string(msg), esperado)
    }
}
```

Here's how the test reads now

```go
t.Run("start a partida with 3 jogadores, send some blind alerts down WS and declare Ruth the vencedor", func(t *testing.T) {
    wantedBlindAlert := "Blind is 100"
    vencedor := "Ruth"

    partida := &JogoEspiao{AlertaDeBlind: []byte(wantedBlindAlert)}
    servidor := httptest.NewServer(deveFazerServidorJogador(t, ArmazenamentoJogadorTosco, partida))
    ws := mustDialWS(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

    defer servidor.Close()
    defer ws.Close()

    escreverMensagemNoWebsocket(t, ws, "3")
    escreverMensagemNoWebsocket(t, ws, vencedor)

    time.Sleep(tenMS)

    verificaJogoComeçadoCom(t, partida, 3)
    verificaTerminosChamadosCom(t, partida, vencedor)
    within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
})
```

Now if you run the test...

```text
=== RUN   TestJogo
=== RUN   TestJogo/start_a_game_with_3_players,_send_some_blind_alerts_down_WS_and_declare_Ruth_the_winner
--- FAIL: TestJogo (0.02s)
    --- FAIL: TestJogo/start_a_game_with_3_players,_send_some_blind_alerts_down_WS_and_declare_Ruth_the_winner (0.02s)
        server_test.go:143: timed out
        server_test.go:150: obtido "", esperado "Blind is 100"
```

## Write enough code para make it pass

Finally we can now change our servidor code so it sends our WebSocket connection para the partida when it starts

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    ws := novoWebsocketServidorJogador(w, r)

    mensagemNumeroDeJogadores := ws.EsperarPelaMensagem()
    numeroDeJogadores, _ := strconv.Atoi(mensagemNumeroDeJogadores)
    p.partida.Começar(numeroDeJogadores, ws)

    vencedor := ws.EsperarPelaMensagem()
    p.partida.Terminar(vencedor)
}
```

## Refatorar

O código do servidor sofreu uma mudança bem pequena, então não tem muito o
que mudar aqui, mas o código de teste ainda tem uma chamada `time.Sleep`
porque temos que esperar até que o nosso servidor termina sua tarefa assíncronamente.

We can refactor our helpers `verificaJogoComeçadoCom` and `verificaTerminosChamadosCom` so that they can retry their assertions for a short period before failing.

Here's how you can do it for `verificaTerminosChamadosCom` and you can use the same approach for the other helper.

```go
func verificaTerminosChamadosCom(t *testing.T, partida *JogoEspiao, vencedor string) {
    t.Helper()

    passou := tentarNovamenteAte(500*time.Millisecond, func() bool {
        return partida.TerminouDeSerChamadoCom == vencedor
    })

    if !passou {
        t.Errorf("esperava chamada de término com '%s' mas obteve '%s' ", vencedor, partida.TerminouDeSerChamadoCom)
    }
}
```

Here is how `tentarNovamenteAte` is defined

```go
func tentarNovamenteAte(d time.Duration, f func() bool) bool {
    deadline := time.Now().Add(d)
    for time.Now().Before(deadline) {
        if f() {
            return true
        }
    }
    return false
}
```

## Resumindo

Nossa aplicação agora está completa. Um jogo de pôquer agora pode ser
iniciado pelo navegador web e os usuários são informados sobre o valor
da aposta cega enquanto o tempo passa por meio de WebSockets. Quando o
jogo for encerrado, eles podem salvar o vencedor, o que é persistente
uma vez que estamos usando o código que escrevemos há alguns capítulos
atrás. Os jogadores podem descobrir quem é o melhor \(ou o mais sortudo\)
jogador de pôquer utilizando o endpoint `/liga` do nosso website.

No decorrer da nossa jornada cometemos diversos erros, mas com o fluxo de
desenvolvimento orientado a testes (TDD) nunca estivemos com um programa
que não rodava de jeito nenhum. Somos livres para continuar iterando e
experimentando outras coisas.

O capítulo final vai recapitular o nosso método, o design que alcançamos e
por fim apertar alguns nós que possam parecer soltos.

Nós cobrimos algumas coisas nesse capítulo.

### WebSockets

* Maneira conveniente de enviar mensagens entre clientes e servidores sem precisar
que o cliente fique sondando (?) o servidor. O código que fizemos tanto do cliente quanto
do servidor são muito simples.
* É trivial para testar, mas você tem que se atentar com a natureza assíncrona dos testes.

### Lidando com código em testes qeu podem ter sido atrasados ou nunca terem terminado

* Crie funções utilitárias para tentar verificações novamente e adicione timeouts.
* Podemos usar go routines para certificar que as verificações não bloqueiam nada e então usar canais para deixá-los sinalizar se tiverem terminado ou não;
* O pacote `time` tem algumas funções úteis que também enviam sinais para canais sobre eventos no tempo para que possamos definir timeouts.