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

Nosso servidor tem um roteador definido, então deve ser relativamente fácil corrigir isso.

Adicione o seguinte no nosso roteador:

```go
roteador.Handle("/partida", http.HandlerFunc(p.partida))
```

E então escreva o método `partida`:

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

Você também vai notar que mudei o `verificaStatus` para aceitar `resposta` ao invés de `resposta.Code` já que parece combinar melhor.

Agora precisamos que o endpoint retorne um pouco de HTML, e aqui está ele:

```markup
<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <title>Vamos jogar pôquer</title>
</head>
<corpo>
<section id="partida">
    <div id="declare-vencedor">
        <label for="vencedor">Vencedor</label>
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

Put the markup in a arquivo called `jogo.html`

Next change the endpoint we just wrote para the following

```go
func (p *ServidorJogador) partida(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("jogo.html")

    if err != nil {
        http.Error(w, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, nil)
}
```

[`html/template`](https://golang.org/pkg/html/template/) is a Go package for creating HTML. In our case we call `template.ParseFiles`, giving the path of our html arquivo. Assuming there is no error you can then `Execute` the template, which writes it para an `io.Writer`. In our case we esperado it para `Write` para the internet, so we give it our `http.ResponseWriter`.

As we have not written a test, it would be prudent para manually test our web servidor just para make sure things are working as we'd hope. Go para `cmd/webserver` and run the `main.go` arquivo. Visit `http://localhost:5000/partida`.

You _should_ have obtido an error about not being able para find the template. You can either change the path para be relative para your folder, or you can have a copy of the `jogo.html` in the `cmd/webserver` directory. I chose para create a symlink \(`ln -s ../../jogo.html jogo.html`\) para the arquivo inside the root of the project so if I make changes they are reflected when running the servidor.

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

Certifique-se que tenha importado o pacote `websocket`. Minha IDE fez isso automaticamente para mim e a sua deve fazer o mesmo.

Para testar o que acontece do navegador, temos que abrir nossa própria conexão WebSocket e escrever nela.

Nossos testes anteriores do servidor apenas chamavam métodos no nosso servidor, mas agora precisamos ter uma conexão persistente nele. Para fazer isso, usamos o `httptest.NewServer`, que recebe um `http.Handler` que vai esperar conexões.

Ao usar `websocket.DefaultDialer.Dial`, tentamos conectar no nosso servidor para então enviar uma mensagem com nosso `vencedor`.

Por fim, verificamos o armazenamento do jogador para certificar que o vencedor foi gravado.

## Execute o teste

```text
=== RUN   TestJogo/quando_recebemos_uma_mensagem_via_websocket_que_ha_um_vencedor_de_uma_partida
    --- FAIL: TestJogo/quando_recebemos_uma_mensagem_via_websocket_que_ha_um_vencedor_de_uma_partida (0.00s)
        server_test.go:124: não foi possível abrir uma conexão de websocket em ws://127.0.0.1:55838/ws websocket: bad handshake
```

Não mudamos nosso servidor para aceitar conexões WebSocket em `/ws`, então ainda não estamos [apertando as mãos](https://pt.wikipedia.org/wiki/Handshake).

## Escreva código suficiente para fazer o teste passar

Adicione outra linha no nosso roteador:

```go
roteador.Handle("/ws", http.HandlerFunc(p.webSocket))
```

E adicione nosso novo manipulador `webSocket`:

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }
    upgrader.Upgrade(w, r, nil)
}
```

Para aceitar uma conexão WebSocker, precisamos de um método `Upgrade` para atualizar a requisição. Agora, se você executar o teste novamente, o próximo erro deve aparecer.

```text
=== RUN   TestJogo/quando_recebemos_uma_mensagem_via_websocket_que_ha_um_vencedor_de_uma_partida
    --- FAIL: TestJogo/quando_recebemos_uma_mensagem_via_websocket_que_ha_um_vencedor_de_uma_partida (0.00s)
        server_test.go:132: obtido 0 chamadas paraGravarVitoria esperado 1
```

Agora que temos uma conexão aberta, vamos esperar por uma pensagem e então gravá-la como vencedor.

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

\(Sim, estamos ignorando vários erros nesse momento!\)

`conexão.ReadMessage()`  bloqueia a espera por uma mensagem na conexão. Quando obtivermos uma, vamos usá-la para `GravarVitoria`. Isso finalmente fecharia a conexão WebSocket.

Se tentar executar o teste, ele ainda vai falhar.

O problema está no tempo. Há um atraso entre nossa conexão WebSocket ler a mensagem e gravar a vitória e nosso teste termina sua execução antes disso acontecer. Você pode testar isso colocando um `time.Sleep` curto antes da verificação final.

Vamos continuar com isso por enquanto, mas saiba que colocar sleeps arbitrários em testes **é uma prática muito ruim**.

```go
time.Sleep(10 * time.Millisecond)
VerificaVitoriaDoVencedor(t, armazenamento, vencedor)
```

## Refatore

Cometemos vários pecados para fazer esse teste funcionar tanto no código do servidor quanto no código do teste, mas lembre-se que essa é a forma mais fácil para fazer as coisas funcionarem.

Temos um software horrível e cheio de gambiarras _funcionando_ apoiado por testes, então agora temos a liberdade para torná-lo elegante sabendo que não vamos quebrar nada por acidente.

Então, vamos começar com o código do servidor.

Podemos mover o `upgrader` para um valor privado dentro do nosso pacote porque não precisamos redeclará-lo em toda requisição na conexão com o WebSocket.

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

Nossa chamada para `template.ParseFiles("jogo.html")` vai ser executada a cada `GET /partida`, o que significa que vamos usar o sistema de arquivo a cada requisição apesar de não ser necessário parsear o template novamente. Vamos refatorar o código para que possamos fazer o parse do template uma vez em `NovoServidorJogador` ao invés disso. Vamos ter que fazer isso para que nossa função possa retornar um erro caso tenhamos problema ao obter o template do disco ou fazer parse dele.

Agora vamos às mudanças relevantes do `ServidorJogador`:

```go
type ServidorJogador struct {
    armazenamento ArmazenamentoJogador
    http.Handler
    template *template.Template
}

const caminhoTemplateHTML = "jogo.html"

func NovoServidorJogador(armazenamento ArmazenamentoJogador) (*ServidorJogador, error) {
    p := new(ServidorJogador)

    tmpl, err := template.ParseFiles("jogo.html")

    if err != nil {
        return nil, fmt.Errorf("problema ao abrir %s %v", caminhoTemplateHTML, err)
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

Ao alterar a assinatura de `NovoServidorJogador`, agora temos problemas de compilação. Tente corrigir por si só ou olhe para o código fonte caso para ver a solução.

Para o código de teste, fiz uma função auxiliar chamada `deveFazerServidorJogador(t *testing.T, armazenamento ArmazenamentoJogador) *ServidorJogador` para que eu possa esconder o erro dos testes.

```go
func deveFazerServidorJogador(t *testing.T, armazenamento ArmazenamentoJogador) *ServidorJogador {
    servidor, err := NovoServidorJogador(armazenamento)
    if err != nil {
        t.Fatal("problema ao criar o servidor do jogador", err)
    }
    return servidor
}
```

Da mesma forma, criei outra função auxiliar `deveConectarAoWebSocket` para que eu possa esconder um erro ao criar uma conexão de WebSocket.

```go
func deveConectarAoWebSocket(t *testing.T, url string) *websocket.Conn {
    ws, _, err := websocket.DefaultDialer.Dial(url, nil)

    if err != nil {
        t.Fatalf("não foi possível abrir uma conexão de websocket em %s %v", url, err)
    }

    return ws
}
```

Finalmente, podemos criar uma função auxiliar no nosso código de teste para enviar mensagens:

```go
func escreverMensagemNoWebsocket(t *testing.T, conexão *websocket.Conn, mensagem string) {
    t.Helper()
    if err := conexão.WriteMessage(websocket.TextMessage, []byte(mensagem)); err != nil {
        t.Fatalf("não foi possível enviar mensagem na conexão websocket %v", err)
    }
}
```

Agora que os testes estão passando, tent executar o servidor e declarar alguns vencedores em `/partida`. Devemos vê-los gravados em `/liga`. Lembre-se que sempre que tivermos um vencedor, vamos _fechar a conexão_, e você vai precisar atualizar a página para abrir a conexão novamente.

Fizemos um formulário simples da web que permite que usuários gravem o vencedor de uma partida. Vamos iterar nele para fazer com que o usuário possa começar uma partida inserindo o número de jogadores e o servidor vai mostrar mensagens para o cliente informando-o qual é o valor do blind conforme o tempo passa.

Primeiramente, atualize o `jogo.html` para atualizar o código do lado do cliente para os novos requerimentos:

```markup
<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <title>Vamos jogar pôquer</title>
</head>
<corpo>
<section id="partida">
    <div id="partida-start">
        <label for="jogador-count">Número de jogadores</label>
        <input type="number" id="jogador-count"/>
        <button id="start-partida">Começar</button>
    </div>

    <div id="declare-vencedor">
        <label for="vencedor">Vencedor</label>
        <input type="text" id="vencedor"/>
        <button id="vencedor-button">Declare vencedor</button>
    </div>

    <div id="blind-value"/>
</section>

<section id="partida-end">
    <h1>Outra ótima partida de pôquer, pessoal!!</h1>
    <p><a href="/liga">Verifique a tabela da liga</a></p>
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

As principais alterações envolvem inserir uma seção para definir o número de jogadores e uma seção para mostrar o valor do blind. Temos um pouco de lógica para mostrar/esconder a interface do usuário dependendo da etapa da partida.

Para qualquer mensagem que recebermos via `conexão.onmessage`, presumimos ser alertas de blind e então definimos o `blindContainer.innerText` de acordo.

Como fazemos para enviar os alertas de blind?No capítulo anterior, mostramos a ideia de `Jogo` para que nosso código CLI possa chamar um `Jogo` e todo o restante se responsabilizaria por agendar os alertas de blind. Isso acabou até sendo uma boa separação de responsabilidades.

```go
type Jogo interface {
    Começar(numeroDeJogadores int)
    Terminar(vencedor string)
}
```

Quando o usuário era requisitado pela CLI pelo número de jogadores, ele precisava `Começar` a partida, o que ativaria os alertas de blind, e quando o usuario declarava o vencedor, isso iria `Terminar`. Esses sã os mesmos requerimentos que temos agora, só que a obtenção das entradas era diferente; logo, só precisamos reutilizar esse conceito aonde possível.

Nossa implementação "real" de `Jogo` é `TexasHoldem`:

```go
type TexasHoldem struct {
    alerter AlertadorDeBlind
    armazenamento   ArmazenamentoJogador
}
```

Ao enviar um `AlertadorDeBlind`, o `TexasHoldem` pode agendar alertas de blind para enviar para _qualquer lugar_.

```go
type AlertadorDeBlind interface {
    AgendarAlertaPara(duracao time.Duration, quantia int)
}
```

E só para lembrar, aqui está nossa implementação do `AlertadorDeBlind` que usamos na CLI.

```go
func SaidaAlertador(duracao time.Duration, quantia int) {
    time.AfterFunc(duracao, func() {
        fmt.Fprintf(os.Stdout, "Blind agora é %d\n", quantia)
    })
}
```

Isso funciona no CLI porque estamos _sempre esperando para enviar os alertas ara `os.Stdout`_, mas isso não vai funcionar no nosso servidor web. Para cada requisição, obtemos um novo `http.ResponseWriter` que então melhoramos para uma `*websocket.Conn`. Logo, não odemos saber quando construímos nossas dependências para onde nossos alertas precisam ir.

Por esse motivo, precisamos mudar o `AlertadorDeBlind.AgendarAlertaPara` para que ele receba um destino paara os alertas para que possamos reutiliza-lo no nosso servidor web.

Abra o AlertadorDeBlind.go e adicione o parâmetro para io.Writer`:

```go
type AlertadorDeBlind interface {
    AgendarAlertaPara(duracao time.Duration, quantia int, para io.Writer)
}

type AlertadorDeBlindFunc func(duracao time.Duration, quantia int, para io.Writer)

func (a AlertadorDeBlindFunc) AgendarAlertaPara(duracao time.Duration, quantia int, para io.Writer) {
    a(duracao, quantia, para)
}
```

A ideia de um `SaidaAlertador` não encaixa bem no nosso modelo, então vamos apenas renomeá-lo para `Alertador`:

```go
func Alertador(duracao time.Duration, quantia int, para io.Writer) {
    time.AfterFunc(duracao, func() {
        fmt.Fprintf(para, "Blind agora é %d\n", quantia)
    })
}
```

Se tentar compilar, haverá uma falha em `TexasHoldem` porque estamos chamando `AgendarAlertaPara` sem uma descrição. Só para deixar tudo compilando novamente, vamos escrevê-lo para `os.Stdout`.

Execute os testes e eles vão falhar porque o `AlertadorDeBlindEspiao` não implementa mais o `AlertadorDeBlind`. Corrija isso atualizando a assinatura de `AgendarAlertaPara`, execute os testes e todos devem estar passando.

Não faz sentido nenhum que o `TexasHoldem` saiba para onde enviar os alertas de blind. Agora, vamos atualizar o `Jogo` para que quando você começa uma partida, declare _para onde_ os alertas devem ir.

```go
type Jogo interface {
    Começar(numeroDeJogadores int, destinoDosAlertas io.Writer)
    Terminar(vencedor string)
}
```

Deixe o compilador te dizer o que precisa ser corrigido. As alterações não são tão ruins:

* Atualize o `TexasHoldem` para que implemente `Jogo` corretamente
* No `CLI`, quando começamos a partida, preciamos passar nosssa propriedade `saida` \(`cli.partida.Começar(numeroDeJogadores, cli.saida`\)
* No teste do `TexasHoldem`, precisamos usar `partida.Começar(5, ioutil.Discard)` para corrigir o problema de compilação e configurar a saída do alerta para ser descartada 

Se tiver feito tudo certo, todos os testes devem passar! Agora podemos usar `Jogo` dentro do `Servidor`.

## Escreva os testes primeiro

Os requerimentos de `CLI` e `Servidor` são os mesmos! É apenas o mecanismo de entrega que é diferente.

Vamos dar uma olhada no nosso teste do `CLI` para inspiração.

```go
t.Run("começa partida com 3 jogadores e termina partida com 'Chris' como vencedor", func(t *testing.T) {
    partida := &JogoEspiao{}

    saida := &bytes.Buffer{}
    in := usuarioEnvia("3", "Chris venceu")

    poquer.NovaCLI(in, saida, partida).JogarPoquer()

    verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador)
    verificaJogoComeçadoCom(t, partida, 3)
    verificaTerminosChamadosCom(t, partida, "Chris")
})
```

Parece que devemos ser capazes de testar um resultado semelhante usando `JogoEspiao`.

Substitua o antigo teste de websocket com o seguinte:

```go
t.Run("começa uma partida com 3 jogadores e declara Ruth vencedora", func(t *testing.T) {
    partida := &poquer.JogoEspiao{}
    vencedor := "Ruth"
    servidor := httptest.NewServer(deveFazerServidorJogador(t, ArmazenamentoJogadorTosco, partida))
    ws := deveConectarAoWebSocket(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

    defer servidor.Close()
    defer ws.Close()

    escreverMensagemNoWebsocket(t, ws, "3")
    escreverMensagemNoWebsocket(t, ws, vencedor)

    time.Sleep(10 * time.Millisecond)
    verificaJogoComeçadoCom(t, partida, 3)
    verificaTerminosChamadosCom(t, partida, vencedor)
})
```

* Conforme discutidos, criamos um espião de `Jogo` e passamos para o ``deveFazerServidorJogador` \(certifique-se de atualizar a função auxiliar para isso\).
* Depois, enviamos mensagens no web socket para uma partida.
* Por mim, verificamos que a partida começou e finalizamos com o que esperamos.

## Execute o teste

Você terá vários erros de compilação envolvendo `deveFazerServidorJogador` em outros testes. Crie uma variável não exportada `jogoTosco` e use-a em todos os testes que não estão compilando:

```go
var (
    jogoTosco = &JogoEspiao{}
)
```

O erro final se encontra onde estamos tentando passar em `Jogo`, pois `NovoServidorJogador` ainda não o suporta:

```text
./server_test.go:21:38: too many arguments in call para "github.com/larien/learn-go-with-tests/WebSockets/v2".NovoServidorJogador
    have ("github.com/larien/learn-go-with-tests/WebSockets/v2".ArmazenamentoJogador, "github.com/larien/learn-go-with-tests/WebSockets/v2".Jogo)
    esperado ("github.com/larien/learn-go-with-tests/WebSockets/v2".ArmazenamentoJogador)
```

## Escreva o mínimo de código possível para o teste funcionar e verifique a saída do teste falhado

Basta adicionar um argumento por enquanto para fazer o teste funcionar:

```go
func NovoServidorJogador(armazenamento ArmazenamentoJogador, partida Jogo) (*ServidorJogador, error) {
```

Finalmente!

```text
=== RUN   TestJogo/começa_um_jogo_com_3_jogadores_e_declara_Ruth_a_vencedora
--- FAIL: TestJogo (0.01s)
    --- FAIL: TestJogo/começa_um_jogo_com_3_jogadores_e_declara_Ruth_a_vencedora (0.01s)
        server_test.go:146: esperava Começar chamado com 3 mas obteve 0
        server_test.go:147: esperava Terminar chamado com 'Ruth' mas obteve ''
FAIL
```

## Escreva código suficiente para fazer o teste passar

Precisamos adicionar `Jogo` como campo para `ServidorJogador` para que possamos usá-lo quando ele obtiver requisições.

```go
type ServidorJogador struct {
    armazenamento ArmazenamentoJogador
    http.Handler
    template *template.Template
    partida Jogo
}
```

\(Já temos um método chamado `partida`, então é só renomeá-lo para `jogarJogo`\)

A seguir, vamos atribui-lo no nosso construtor:

```go
func NovoServidorJogador(armazenamento ArmazenamentoJogador, partida Jogo) (*ServidorJogador, error) {
    p := new(ServidorJogador)

    tmpl, err := template.ParseFiles("jogo.html")

    if err != nil {
        return nil, fmt.Errorf("problema ao abrir %s %v", caminhoTemplateHTML, err)
    }

    p.partida = partida

    // etc
```

Agora podemos usar nosso `Jogo` dentro de `webSocket`.

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    conexão, _ := atualizadorDeWebsocket.Upgrade(w, r, nil)

    _, mensagemNumeroDeJogadores, _ := conexão.ReadMessage()
    numeroDeJogadores, _ := strconv.Atoi(string(mensagemNumeroDeJogadores))
    p.partida.Começar(numeroDeJogadores, ioutil.Discard) //todo: Não descartar as mensagens de blind!

    _, vencedor, _ := conexão.ReadMessage()
    p.partida.Terminar(string(vencedor))
}
```

Uhul! Os testes estão passando.

Não vamos enviar as mensagens de blind para nenhum lugar _por enquanto_ já que precisamos de um tempo para pensar nisso. Quando chamamos `partida.Começar`, enviamos os dados para `ioutil.Discard` que vai apenar descartar qualquer mensagem escrita nele.

Por enquanto, vamos iniciar o servidor. Você vai precisar atualizar a ``main.go` para passar um `Jogo` para o `ServidorJogador`:

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

Tirando o fato de que não temos alertas de blind por enquanto, a aplicação funciona! Conseguimos reutilizar `Jogo` com `ServidorJogador` e ele toma conta dos detalhes. Quando descobrirmos como enviar mensagens de blind atraves de web sockets ao invés de descartá-las, tudo _deve_ ficar pronto.

Antes disso, vamos mexer um pouco no código.

## Refatore

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
    p.partida.Começar(numeroDeJogadores, ioutil.Discard) //todo: Não descartar as mensagens de blind!

    vencedor := ws.EsperarPelaMensagem()
    p.partida.Terminar(vencedor)
}
```

Once we figure out how para not discard the blind mensagens we're done.

### Let's _not_ write a test!

Sometimes when we're not sure how para do something, it's best just para play around and try things out! Make sure your work is committed first because once we've figured out a way we should drive it through a test.

The problematic line of code we have is

```go
p.partida.Começar(numeroDeJogadores, ioutil.Discard) //todo: Não descartar as mensagens de blind!
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
    ws := deveConectarAoWebSocket(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

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
    ws := deveConectarAoWebSocket(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

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

## Escreva código suficiente para fazer o teste passar

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