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
quando eles acessarem `/jogo`.

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
código `200` quando acessamos o `GET /jogo`.

```go
func TestJogo(t *testing.T) {
    t.Run("GET /jogo retorna 200", func(t *testing.T) {
        servidor := NovoServidorJogador(&EsbocoDeArmazenamentoJogador{})

        requisicao, _ := http.NewRequest(http.MethodGet, "/jogo", nil)
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
roteador.Handle("/jogo", http.HandlerFunc(p.jogo))
```

E então escreva o método `jogo`:

```go
func (p *ServidorJogador) jogo(w http.ResponseWriter, r *http.Request) {
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
    t.Run("GET /jogo retorna 200", func(t *testing.T) {
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
<section id="jogo">
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

### Como testamos que retornamos a marcação correta?

Existem algumas formas. Como foi enfatizado no decorrer do livro, é importante que os testes que você escreve têm valor o suficiente para justificar o custo.

1. Escreva um teste baseado no navegador, usando algo como Selenium. Esses testes são os mais "realistas" de todas as abordagens porque começam um navegador web de verdade e simula um usuário interagindo com ele. Esses testes podem te dar muita confiança de que seu sistma funciona, mas são mais difíceis e escrever que os testes unitários e muito mais lentos de serem executados. Para os propósitos do nosso produto, isso é exagero.
2. Fazer uma comparação exata de textos. Isso _pode_ funcionar, mas esses tipos de testes acabam sendo muito frágeis. No momento que alguém muda a marcação, você vai ter um teste falhando quando na prática nada está _de fato_ falhando.
3. Verificar que chamamos o template correto. Vamos usar uma biblioteca de template da biblioteca padrão para servir o HTML (que falamos brevemente) e podemos injetar na _coisa_ que gera o HTML e espionar suas chamadas para verificar que estamos fazendo tudo corretamente. Isso teria um impacto no design do nosso código, mas na realidade isso não estaríamos testando algo tão crítico além de verificar se estamos chamando o arquivo de template correto. Dito isso, só vamos ter um template no nosso projeto e a chance de falha aqui parece pequena.

Então, pela primeira vez no livro "Aprenda Go com Testes", não vamos escrever nenhum teste.

Coloque a marcação em um arquivo chamado `jogo.html`.

Na próxima mudança do endoint, vamos apenas escrever o seguinte:

```go
func (p *ServidorJogador) jogo(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("jogo.html")

    if err != nil {
        http.Error(w, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
        return
    }

    tmpl.Execute(w, nil)
}
```

O [`html/template`](https://golang.org/pkg/html/template/) é um pacote do Go para criar HTML. No nosso caso, chamamos `template.ParseFiles` enviando o caminho do nosso arquivo HTML. Presumindo que não há nenhum erro, chamamos a função `Execute` para "executar" o template, que o escreve para um `ìo.Writer`. No nosso caso, esperamos que o template seja escrito na internet, então enviamos o nosso `http.ResponseWriter`.

Já que não escrevemos um teste, seria prudente testar nosso servidor web manualmente só para ter certeza de que as coisas estão funcionamos como esperamos. Vá para `cmd/webserver` e execute o arquivo `main.go`. Visite `http://localhost:5000/jogo`.

Você _deve_ ter obtido um erro sobre não ser capaz de encontrar o template. Você pode ou mudar o caminho para ser relativo à sua pasta, ou pode ter uma cópia de ``jogo.html` no diretório `cmd/webserver`. Eu escolho criar um symlink \(`ln -s ../../jogo.html jogo.html`\) para o arquivo dentro da raiz do projeto para caso eu faça alterações, elas reflitam quando o servidor estiver sendo executado.

Se fizer essa alteração e rodar novamente, deve conseguir ver a interface.

Agora precisamos testar que, quando obtemos uma string sob uma conexão WebSocket para o nosso servidor, declaramos a pessoa como vencedora de um jogo.

## Escreva o teste primeiro

Pela primeia vez, vamos usar uma biblioteca externa para trabalhar com WebSockets.

Rode `go get github.com/gorilla/websocket`.

Isso vai obter o código para a excelente biblioteca [Gorilla WebSocket](https://github.com/gorilla/websocket). Agora podemos atualizar nossos testes para nosso novo requerimento.

```go
t.Run("quando recebemos uma mensagem de um websocket que é vencedor da jogo", func(t *testing.T) {
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
=== RUN   TestJogo/quando_recebemos_uma_mensagem_via_websocket_que_ha_um_vencedor_de_uma_jogo
    --- FAIL: TestJogo/quando_recebemos_uma_mensagem_via_websocket_que_ha_um_vencedor_de_uma_jogo (0.00s)
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
=== RUN   TestJogo/quando_recebemos_uma_mensagem_via_websocket_que_ha_um_vencedor_de_uma_jogo
    --- FAIL: TestJogo/quando_recebemos_uma_mensagem_via_websocket_que_ha_um_vencedor_de_uma_jogo (0.00s)
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

Nossa chamada para `template.ParseFiles("jogo.html")` vai ser executada a cada `GET /jogo`, o que significa que vamos usar o sistema de arquivo a cada requisição apesar de não ser necessário parsear o template novamente. Vamos refatorar o código para que possamos fazer o parse do template uma vez em `NovoServidorJogador` ao invés disso. Vamos ter que fazer isso para que nossa função possa retornar um erro caso tenhamos problema ao obter o template do disco ou fazer parse dele.

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
    roteador.Handle("/jogo", http.HandlerFunc(p.jogo))
    roteador.Handle("/ws", http.HandlerFunc(p.webSocket))

    p.Handler = roteador

    return p, nil
}

func (p *ServidorJogador) jogo(w http.ResponseWriter, r *http.Request) {
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

Agora que os testes estão passando, tent executar o servidor e declarar alguns vencedores em `/jogo`. Devemos vê-los gravados em `/liga`. Lembre-se que sempre que tivermos um vencedor, vamos _fechar a conexão_, e você vai precisar atualizar a página para abrir a conexão novamente.

Fizemos um formulário simples da web que permite que usuários gravem o vencedor de uma jogo. Vamos iterar nele para fazer com que o usuário possa começar uma jogo inserindo o número de jogadores e o servidor vai mostrar mensagens para o cliente informando-o qual é o valor do blind conforme o tempo passa.

Primeiramente, atualize o `jogo.html` para atualizar o código do lado do cliente para os novos requerimentos:

```markup
<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <title>Vamos jogar pôquer</title>
</head>
<corpo>
<section id="jogo">
    <div id="jogo-start">
        <label for="jogador-count">Número de jogadores</label>
        <input type="number" id="jogador-count"/>
        <button id="start-jogo">Começar</button>
    </div>

    <div id="declare-vencedor">
        <label for="vencedor">Vencedor</label>
        <input type="text" id="vencedor"/>
        <button id="vencedor-button">Declare vencedor</button>
    </div>

    <div id="blind-value"/>
</section>

<section id="jogo-end">
    <h1>Outra ótima jogo de pôquer, pessoal!!</h1>
    <p><a href="/liga">Verifique a tabela da liga</a></p>
</section>

</corpo>
<script type="application/javascript">
    const startGame = document.getElementById('jogo-start')

    const declareWinner = document.getElementById('declare-vencedor')
    const submitWinnerButton = document.getElementById('vencedor-button')
    const entradaVencedor = document.getElementById('vencedor')

    const blindContainer = document.getElementById('blind-value')

    const gameContainer = document.getElementById('jogo')
    const gameEndContainer = document.getElementById('jogo-end')

    declareWinner.hidden = true
    gameEndContainer.hidden = true

    document.getElementById('start-jogo').addEventListener('click', event => {
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

As principais alterações envolvem inserir uma seção para definir o número de jogadores e uma seção para mostrar o valor do blind. Temos um pouco de lógica para mostrar/esconder a interface do usuário dependendo da etapa da jogo.

Para qualquer mensagem que recebermos via `conexão.onmessage`, presumimos ser alertas de blind e então definimos o `blindContainer.innerText` de acordo.

Como fazemos para enviar os alertas de blind?No capítulo anterior, mostramos a ideia de `Jogo` para que nosso código CLI possa chamar um `Jogo` e todo o restante se responsabilizaria por agendar os alertas de blind. Isso acabou até sendo uma boa separação de responsabilidades.

```go
type Jogo interface {
    Começar(numeroDeJogadores int)
    Terminar(vencedor string)
}
```

Quando o usuário era requisitado pela CLI pelo número de jogadores, ele precisava `Começar` a jogo, o que ativaria os alertas de blind, e quando o usuario declarava o vencedor, isso iria `Terminar`. Esses sã os mesmos requerimentos que temos agora, só que a obtenção das entradas era diferente; logo, só precisamos reutilizar esse conceito aonde possível.

Nossa implementação "real" de `Jogo` é `TexasHoldem`:

```go
type TexasHoldem struct {
    alertador AlertadorDeBlind
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

Isso funciona no CLI porque estamos _sempre esperando para enviar os alertas para `os.Stdout`_, mas isso não vai funcionar no nosso servidor web. Para cada requisição, obtemos um novo `http.ResponseWriter` que então melhoramos para uma `*websocket.Conn`. Logo, não odemos saber quando construímos nossas dependências para onde nossos alertas precisam ir.

Por esse motivo, precisamos mudar o `AlertadorDeBlind.AgendarAlertaPara` para que ele receba um destino paara os alertas para que possamos reutiliza-lo no nosso servidor web.

Abra o `AlertadorDeBlind.go` e adicione o parâmetro para io.Writer`:

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

Não faz sentido nenhum que o `TexasHoldem` saiba para onde enviar os alertas de blind. Agora, vamos atualizar o `Jogo` para que quando você começa uma jogo, declare _para onde_ os alertas devem ir.

```go
type Jogo interface {
    Começar(numeroDeJogadores int, destinoDosAlertas io.Writer)
    Terminar(vencedor string)
}
```

Deixe o compilador te dizer o que precisa ser corrigido. As alterações não são tão ruins:

* Atualize o `TexasHoldem` para que implemente `Jogo` corretamente
* No `CLI`, quando começamos a jogo, preciamos passar nosssa propriedade `saida` \(`cli.jogo.Começar(numeroDeJogadores, cli.saida`\)
* No teste do `TexasHoldem`, precisamos usar `jogo.Começar(5, ioutil.Discard)` para corrigir o problema de compilação e configurar a saída do alerta para ser descartada 

Se tiver feito tudo certo, todos os testes devem passar! Agora podemos usar `Jogo` dentro do `Servidor`.

## Escreva os testes primeiro

Os requerimentos de `CLI` e `Servidor` são os mesmos! É apenas o mecanismo de entrega que é diferente.

Vamos dar uma olhada no nosso teste do `CLI` para inspiração.

```go
t.Run("começa jogo com 3 jogadores e termina jogo com 'Chris' como vencedor", func(t *testing.T) {
    jogo := &JogoEspiao{}

    saida := &bytes.Buffer{}
    in := usuarioEnvia("3", "Chris venceu")

    poquer.NovaCLI(in, saida, jogo).JogarPoquer()

    verificaMensagensEnviadasParaUsuario(t, saida, poquer.PromptJogador)
    verificaJogoComeçadoCom(t, jogo, 3)
    verificaTerminosChamadosCom(t, jogo, "Chris")
})
```

Parece que devemos ser capazes de testar um resultado semelhante usando `JogoEspiao`.

Substitua o antigo teste de websocket com o seguinte:

```go
t.Run("começa uma jogo com 3 jogadores e declara Ruth vencedora", func(t *testing.T) {
    jogo := &poquer.JogoEspiao{}
    vencedor := "Ruth"
    servidor := httptest.NewServer(deveFazerServidorJogador(t, ArmazenamentoJogadorTosco, jogo))
    ws := deveConectarAoWebSocket(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

    defer servidor.Close()
    defer ws.Close()

    escreverMensagemNoWebsocket(t, ws, "3")
    escreverMensagemNoWebsocket(t, ws, vencedor)

    time.Sleep(10 * time.Millisecond)
    verificaJogoComeçadoCom(t, jogo, 3)
    verificaTerminosChamadosCom(t, jogo, vencedor)
})
```

* Conforme discutidos, criamos um espião de `Jogo` e passamos para o ``deveFazerServidorJogador` \(certifique-se de atualizar a função auxiliar para isso\).
* Depois, enviamos mensagens no web socket para uma jogo.
* Por mim, verificamos que a jogo começou e finalizamos com o que esperamos.

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
func NovoServidorJogador(armazenamento ArmazenamentoJogador, jogo Jogo) (*ServidorJogador, error) {
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
    jogo Jogo
}
```

\(Já temos um método chamado `jogo`, então é só renomeá-lo para `jogarJogo`\)

A seguir, vamos atribui-lo no nosso construtor:

```go
func NovoServidorJogador(armazenamento ArmazenamentoJogador, jogo Jogo) (*ServidorJogador, error) {
    p := new(ServidorJogador)

    tmpl, err := template.ParseFiles("jogo.html")

    if err != nil {
        return nil, fmt.Errorf("problema ao abrir %s %v", caminhoTemplateHTML, err)
    }

    p.jogo = jogo

    // etc
```

Agora podemos usar nosso `Jogo` dentro de `webSocket`.

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    conexão, _ := atualizadorDeWebsocket.Upgrade(w, r, nil)

    _, mensagemNumeroDeJogadores, _ := conexão.ReadMessage()
    numeroDeJogadores, _ := strconv.Atoi(string(mensagemNumeroDeJogadores))
    p.jogo.Começar(numeroDeJogadores, ioutil.Discard) //todo: Não descartar as mensagens de blind!

    _, vencedor, _ := conexão.ReadMessage()
    p.jogo.Terminar(string(vencedor))
}
```

Uhul! Os testes estão passando.

Não vamos enviar as mensagens de blind para nenhum lugar _por enquanto_ já que precisamos de um tempo para pensar nisso. Quando chamamos `jogo.Começar`, enviamos os dados para `ioutil.Discard` que vai apenar descartar qualquer mensagem escrita nele.

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

    jogo := poquer.NovoTexasHoldem(poquer.AlertadorDeBlindFunc(poquer.Alertador), armazenamento)

    servidor, err := poquer.NovoServidorJogador(armazenamento, jogo)

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

A forma que estamos usando WebSocker é bem básica e a mnnipulação de erro é bem fraca, então gostaria de encapsular isso em um tipo só para remover essa bagunça do código do servidor. Precisaremos revisitar isso depois, mas por enqaunto isso vai melhorar um pouco as coisas.

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

Agora o código do servidor fica um pouco mais simples:

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    ws := novoWebsocketServidorJogador(w, r)

    mensagemNumeroDeJogadores := ws.EsperarPelaMensagem()
    numeroDeJogadores, _ := strconv.Atoi(mensagemNumeroDeJogadores)
    p.jogo.Começar(numeroDeJogadores, ioutil.Discard) //todo: Não descartar as mensagens de blind!

    vencedor := ws.EsperarPelaMensagem()
    p.jogo.Terminar(vencedor)
}
```

Quando descobrirmos como não descartar as mensagens de blind teremos terminado essa etapa.

### _Não_ vamos escrever um teste!

Às vezes, quando não temos certeza de como vamos fazer algo, é melhor apenas brincar e testar coisas diferentes! Tenha certeza de que seu trabalho está salvo primeiro porque quando descobrirmos o que fazer, vamos implementá-lo junto de um teste.

A linha problemática do código que temos é:

```go
p.jogo.Começar(numeroDeJogadores, ioutil.Discard) //todo: Não descartar as mensagens de blind!
```

Precisamos passar um `io.Writer` para a jogo para ter aonde escrever os alertas be blind.

Não seria legal se apenas precisássemos passar o nosso `websocketServidorJogador` de antes? É o nosso wrapper em torno do nosso WebSocket, então _parece_ que devemos ser capazes de enviá-lo para que nosso `Jogo` seja capaz de enviar mensagens para ele.

Vamos tentar:

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    ws := novoWebsocketServidorJogador(w, r)

    mensagemNumeroDeJogadores := ws.EsperarPelaMensagem()
    numeroDeJogadores, _ := strconv.Atoi(mensagemNumeroDeJogadores)
    p.jogo.Começar(numeroDeJogadores, ws)
    //etc...
```

O compilador reclama:

```text
./servidor.go:71:14: cannot use ws (type *websocketServidorJogador) as type io.Writer in argument para p.jogo.Começar:
    *websocketServidorJogador does not implement io.Writer (missing Write method)
```

Parece que a coisa óbvia a se fazer é fazer com que o `websocketServidorJogador` _implementa_ o `io.Writer`. Para fazer isso, precisamos usar do `*websocket.Conn` para ussar a escrita de mensagem `WriteMessage` para enviar a mensagem para o websocket.

```go
func (w *websocketServidorJogador) Write(p []byte) (n int, err error) {
    err = w.WriteMessage(1, p)

    if err != nil {
        return 0, err
    }

    return len(p), nil
}
```

Isso parece fácil demais! Execute a aplicação para ver se funciona.

Mas antes edite o `TexasHoldem` para que o tempo de incremento do blind seja mais curto para que você possa ver as coisas em ação:

```go
incrementoDeBlind := time.Duration(5+numeroDeJogadores) * time.Second // (ao invés de um minuto)
```

As coisas devem estar funcionando! A quantidade do blind é incrementada no computador como se fosse mágica.

Agora vamos reverter o código e pensar como testá-lo. Para _implementar_ isso tudo o que precisamos fazer foi passar o `websocketServidorJogador` para `ComeçarJogo` no lugar do `ioutil.Discard`, então isso faz parecer que tenhamos que espionar a chamada para verificar se ela funciona.

Espionar é ótimo e nos ajuda a verificar os detalhes de implementação, mas sempre devemos favorecer o teste do comportamento _real_ se possível, porque caso seja necessário refatorar isso os testes espiões são os primeiros a começar a falhar por geralmente verificarem os detalhes de implementação que estamos tentando alterar.

Nosso teste atualmente abre uma conexão websocket para nosso servidor em execução e envia mensagens para fazê-lo efetuar ações. De forma semelhante, devemos ser capazes de testar as mensagens que o nosso servidor envia de volta para a conexão de websocket.

## Escreva o teste primeiro

Vamos editar nosso teste existente.

Atualmente, nosso `JogoEspiao` não envia nenhum dado para a `saida` quando você chama `Começar`. Devemos alterar isso para que possamos configurá-lo para enviar uma mensagem e então verificar se a mensagem é enviada para o websocket. Isso deve nos dar confiança que configuramos as coisas corretamente enquanto ainda exercitamos o comportamento real do que esperamos.

```go
type JogoEspiao struct {
    ComecouASerChamado     bool
    ComecouASerChamadoCom int
    AlertaDeBlind      []byte

    TerminouDeSerChamado   bool
    TerminouDeSerChamadoCom string
}
```

Adicione o campo de `AlertaDeBlind`.

Atualize o `Começar` do `JogoEspiao` para enviar a mensagem para a `saída`.

```go
func (j *JogoEspiao) Começar(numeroDeJogadores int, saida io.Writer) {
    j.ComecouASerChamado = true
    j.ComecouASerChamadoCom = numeroDeJogadores
    saida.Write(j.AlertaDeBlind)
}
```

Agora isso significa que quando usarmos o `ServidorJogador`, quando ele tentar `Começar` o jogo, deve acabar enviando mensagens pelo websocket se as coisas estiverem funcionando direito.

Finalmente podemos atualizar o teste:

```go
t.Run("começa uma artida com  3 jogadores, envia alguns alertas de blind no websocket e declara Ruth como vencedora", func(t *testing.T) {
    alertaDeBlindEsperado := "Blind é 100"
    vencedor := "Ruth"

    jogo := &JogoEspiao{AlertaDeBlind: []byte(alertaDeBlindEsperado)}
    servidor := httptest.NewServer(deveFazerServidorJogador(t, ArmazenamentoJogadorTosco, jogo))
    ws := deveConectarAoWebSocket(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

    defer servidor.Close()
    defer ws.Close()

    escreverMensagemNoWebsocket(t, ws, "3")
    escreverMensagemNoWebsocket(t, ws, vencedor)

    time.Sleep(10 * time.Millisecond)
    verificaJogoComeçadoCom(t, jogo, 3)
    verificaTerminosChamadosCom(t, jogo, vencedor)

    _, alertaDeBlindObtido, _ := ws.ReadMessage()

    if string(alertaDeBlindObtido) != alertaDeBlindEsperado {
        t.Errorf("alerta de blind obtido '%s', esperado '%s'", string(alertaDeBlindObtido), alertaDeBlindEsperado)
    }
})
```

* Adicionamos um `alertaDeBlindEsperado` e configuramos nosso `JogoEspiao` para enviá-lo para a `saida` se `Começar` for chamado.
* Esperamos que ela seja enviada na conexão do websocket, então adicionamos uma chamada para `ws.ReadMessage()` para esperar por uma mensagem ser enviada e então verificamos se é aquela que esperamos.

## Execute o teste

Talvez você pense que o teste demora demais. Isso acontece porque o ``ws.ReadMessage()` vai bloqueá-lo até obter a mensagem, que nunca vai chegar.
## Escreva o mínimo de código necessário para o teste ser executado e verifique a saída do teste falhando

Nunca devemos ter testes que demoram, então vamos apresentar uma nova forma de lidar com coigo que esperamos com um timeout.

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

O que o `within` faz é pegar uma função `assert` como argumento e então o executa dentro de uma goroutine. Se/Quando a função termina, ela avisa que terminou através do canal `done`.

Enquanto isso acontece, usamos uma declaração `select` que nos permite esperar por um canal para enviar uma mensagem. A partir daí é uma corrida entre a função de `assert` e o `time.After` que vai enviar um sinal qunado a duração chega ao fim.

Por mim, fiz uma função auxiliar para a nossa verificação so para melhorar um pouco as coisas:

```go
func verificaSeWebSocketObteveMensagem(t *testing.T, ws *websocket.Conn, esperado string) {
    _, msg, _ := ws.ReadMessage()
    if string(msg) != esperado {
        t.Errorf(`obtido "%s", esperado "%s"`, string(msg), esperado)
    }
}
```

É assim que o teste fica agora:

```go
t.Run("começa uma artida com  3 jogadores, envia alguns alertas de blind no websocket e declara Ruth como vencedora", func(t *testing.T) {
    alertaDeBlindEsperado := "Blind é 100"
    vencedor := "Ruth"

    jogo := &JogoEspiao{AlertaDeBlind: []byte(alertaDeBlindEsperado)}
    servidor := httptest.NewServer(deveFazerServidorJogador(t, ArmazenamentoJogadorTosco, jogo))
    ws := deveConectarAoWebSocket(t, "ws"+strings.TrimPrefix(servidor.URL, "http")+"/ws")

    defer servidor.Close()
    defer ws.Close()

    escreverMensagemNoWebsocket(t, ws, "3")
    escreverMensagemNoWebsocket(t, ws, vencedor)

    time.Sleep(tenMS)

    verificaJogoComeçadoCom(t, jogo, 3)
    verificaTerminosChamadosCom(t, jogo, vencedor)
    within(t, tenMS, func() { verificaSeWebSocketObteveMensagem(t, ws, alertaDeBlindEsperado) })
})
```

Agora se você rodar o teste...

```text
=== RUN   TestJogo
=== RUN   TestJogo/começa_um_jogo_com_3_jogadores,envia_alguns_alertas_de_blind_para_o_websocket_e_declara_Ruth_como_vencedora
--- FAIL: TestJogo (0.02s)
    --- FAIL: TestJogo/começa_um_jogo_com_3_jogadores,envia_alguns_alertas_de_blind_para_o_websocket_e_declara_Ruth_como_vencedora (0.02s)
        server_test.go:143: timed out
        server_test.go:150: obtido "", esperado "Blind é 100"
```

## Escreva código suficiente para fazer o teste passar

Finalmente podemos alterar o código do nosso servidor para que ele envie a mensagem para nossa conexão com o WebSocket para a jogo quando ela começa:

```go
func (p *ServidorJogador) webSocket(w http.ResponseWriter, r *http.Request) {
    ws := novoWebsocketServidorJogador(w, r)

    mensagemNumeroDeJogadores := ws.EsperarPelaMensagem()
    numeroDeJogadores, _ := strconv.Atoi(mensagemNumeroDeJogadores)
    p.jogo.Começar(numeroDeJogadores, ws)

    vencedor := ws.EsperarPelaMensagem()
    p.jogo.Terminar(vencedor)
}
```

## Refatorar

O código do servidor sofreu uma mudança bem pequena, então não tem muito o
que mudar aqui, mas o código de teste ainda tem uma chamada `time.Sleep`
porque temos que esperar até que o nosso servidor termina sua tarefa assíncronamente.

Podemos refatorar nossas funções auxiliares `verificaJogoComeçadoCom` e `verificaTerminosChamadosCom` para que possam tentar as verificações novamente logo após falharem.

Abaixo esta como fazer isso com o `verificaTerminosChamadosCom` e você pode usar a mesma abordagem para a outra função auxiliar.

```go
func verificaTerminosChamadosCom(t *testing.T, jogo *JogoEspiao, vencedor string) {
    t.Helper()

    passou := tentarNovamenteAte(500*time.Millisecond, func() bool {
        return jogo.TerminouDeSerChamadoCom == vencedor
    })

    if !passou {
        t.Errorf("esperava chamada de término com '%s' mas obteve '%s' ", vencedor, jogo.TerminouDeSerChamadoCom)
    }
}
```

Aqui está como  `tentarNovamenteAte` está definida:

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