# Servidor HTTP

[**Você encontra todo o código-fonte para este capítulo aqui**](https://github.com/larien/learn-go-with-tests/tree/master/http-server)

Você recebeu o desafio de criar um servidor web para que usuários possam acompanhar quantas partidas os jogadores venceram.

* `GET /jogadores/{name}` deve retornar um número indicando o número total de vitórias
* `POST /jogadores/{name}` deve registrar uma vitória para este nome de jogador, incrementando a cada nova chamada de submissão de dados (método HTTP `POST`). 

Vamos seguir com a abordagem do Desenvolvimento Orientado a Testes, criando software que funciona o mais rápido possível, e a cada ciclo fazendo pequenas melhorias até uma solução completa. Com essa abordagem, nós

* Mantemos pequeno o escopo do problema em qualquer momento
* Não perdemos o foco por pensar em muito detalhes
* Se ficamos emperrados ou perdidos, podemos voltar para uma versão anterior do código sem perder muito trabalho.

## Vermelho, verde, refatore

Por todo o livro, enfatizamos o processo Desenvolvimento Orientado a Testes de escrever um teste e ver a falha \(vermelho\), escrever a _menor_ quantidade de código para fazer o teste passar/funcionar \(verde\), e então fazemos a reescrita (refatoração).

A disciplina de escrever a menor quantidade de código é importante para garantir a segurança que o Desenvolvimento Orientado a Testes proporciona. Você deve se empenhar em sair do _vermelho_ o quanto antes.

Kent Beck descreve essa prática como:

> Faça o teste passar rapidamente, cometendo quaisquer pecados necessários nesse processo.

E você pode cometer estes pecados porque vamos reescrever o código logo depois, com a segurança garantida pelos testes.

### E se você não fizer assim?

Quanto mais alterações você fizer enquanto seu código estiver em _vermelho_, maiores as chances de você adicionar problemas, não cobertos por testes.

A ideia é escrever iterativamente código útil em pequenos passos, guiados pelos testes, para que você não perca foco no objetivo principal.

### A galinha e o ovo

Como podemos construir isso de forma incremental? Não podemos obter um jogador (método HTTP `GET`) sem tê-lo registrado anteriormente, e parece complicado saber se a chamada ao método HTTP `POST` funcionou sem a chamada ao método HTTP `GET` já implementado.

E é nesse ponto que o uso de _classes com valores predefinidos_ vai nos ajudar. 

(Nota do tradutor: No original, é usada a expressão _mocking_, que significa "zombar", "fazer piada" ou "enganar". Em programação, _mocking_ significa criar _algo_, como uma classe ou função, que retorna os valores esperados de forma predefinida.)

* a implementação que responde ao método HTTP `GET` precisa de uma _coisa_ `JogadorArmazenamento` para obter pontuações de um nome de jogador. Isso deve ser uma interface, para que, ao executar os testes, seja possível criar um código simples de esboço para testar o código sem precisar, neste momento, implementar o código final que será usado para armazenar os dados.
* para o método HTTP `POST`, podemos _inspecionar_ as chamadas feitas a `JogadorArmazenamento` para ter certeza de que os dados são armazenados corretamente. Nossa implementação de gravação dos dados não estará vinculada à busca dos dados.
* para ver código rodando rapidamente vamos fazer uma implementação simples de armazenamento dos dados na memória, e depois podemos criar uma implementação que dá suporte ao mecanismo de armazenamento de preferência.

## Escrevendo o teste primeiro

Podemos escrever um teste e fazer funcionar retornando um valor predeterminado para nos ajudar a começar. Kent Beck se refere a isso como "Fazer de conta". Uma vez que temos um teste funcionando podemos escrever mais testes que nos ajudem a remover este valor predeterminado (constante).

Com este pequeno passo, nós começamos a ter uma estrutura inicial para o projeto funcionando corretamente, sem nos preocuparmos demais com a lógica da aplicação.

Para criar um servidor web (uma aplicação que recebe chamadas via protocolo HTTP) em Go, você vai chamar, usualmente, a função [ListenAndServe](https://golang.org/pkg/net/http/#ListenAndServe).

```go
func ListenAndServe(endereco string, tratador Handler) error
```

Isso vai iniciar um servidor web _escutando_ em uma porta, criando uma gorotina para cada requisição, e repassando para um _Tratador_, que é representado pela interface [`Handler`](https://golang.org/pkg/net/http/#Handler), usada para receber a requisição e avaliar o que fazer com os dados recebidos.

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Esta interface define uma única função que espera dois argumentos, o primeiro que indica onde _escrevemos a resposta_ e o outro com a requisição HTTP que nos foi enviada.

Vamos criar o primeiro arquivo, `servidor_test.go` e escrever um teste para a função `JogadorServidor` que recebe estes dois argumentos. A requisição enviada serve para obter a pontuação de um Nome de Jogador, que esperamos que seja `"20"`.

```go
func TestObterJogadores(t *testing.T) {
    t.Run("retornar resultado de Maria", func(t *testing.T) {
        requisicao, _ := http.NewRequest(http.MethodGet, "/jogadores/Maria", nil)
        resposta := httptest.NewRecorder()

        JogadorServidor(resposta, requisicao)

        recebido := resposta.Body.String()
        esperado := "20"

        if recebido != esperado {
            t.Errorf("recebido '%s', esperado '%s'", recebido, esperado)
        }
    })
}
```

Para testar nosso servidor, vamos precisar de uma _Requisição_ (`Request`) para enviar a requisição ao servidor, e então queremos _inspecionar_ o que o nosso Tratador escreve para o `ResponseWriter`.

* Nós usamos o `http.NewRequest` para criar uma requisição. O primeiro argumento é o _método_ da requisição e o segundo é o caminho da requisição. O valor `nil` para o segundo argumento corresponde ao corpo (_body_) da requisição, que não precisamos definir para este teste.
* `net/http/httptest` já tem um _inspecionador_ criado para nós, chamado `ResponseRecorder` (_GravadorDeResposta_), então podemos usá-lo. Este possui muitos métodos úteis para inspecionar o que foi escrito como resposta.

## Tente rodar o teste

`./servidor_test.go:13:2: undefined: JogadorServidor`

## Escreva a quantidade mínima de código para o que teste passe e verifique a falha indicada na responta do teste

O compilador está aqui para ajudar, ouça o que ele diz.

Crie o arquivo `servidor.go`, e nele a função `JogadorServidor`

```go
func JogadorServidor() {}
```

Tente novamente

```text
./servidor_test.go:13:14: too many arguments in call to JogadorServidor
    have (*httptest.ResponseRecorder, *http.Request)
    want ()
```

Adicione os argumentos à função

```go
import "net/http"

func JogadorServidor(w http.ResponseWriter, r *http.Request) {

}
```

Agora o código compila, e o teste falha.

```text
--- FAIL: TestObterJogadores (0.00s)
    --- FAIL: TestObterJogadores/retornar_resultado_de_Maria (0.00s)
        servidor_test.go:20: recebido '', esperado '20'
```

## Escreva código suficiente para fazer o teste funcionar

Do capítulo sobre injeção de dependências, falamos sobre servidores HTTP com a função `Greet`. Aprendemos que a função `ResponseWriter` também implementa a interface `Writer` do pacote io, então podemos usar `fmt.Fprint` para enviar strings como respostas HTTP.

```go
func JogadorServidor(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "20")
}
```

O teste agora deve funcionar.

## Complete a estrutura

Nós queremos converter o código acima em uma aplicação. Isso é importante porque

* Teremos _software funcionando_; não queremos escrever testes apenas por escrever, e é bom ver código que funciona.
* Conforme refatoramos o código, é provável mudaremos a estrutura do programa. Nós queremos garantir que isso é refletido em nossa aplicação também, como parte da abordagem incremental.

Crie um novo arquivo `main.go` para nossa aplicação, com o código abaixo.

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    tratador := http.HandlerFunc(JogadorServidor)
    if err := http.ListenAndServe(":5000", tratador); err != nil {
        log.Fatalf("não foi possível escutar na porta 5000 %v", err)
    }
}
```

Para executar, execute o comando `go build`, que vai pegar todos os arquivos terminados em `.go` neste diretório e construir seu programa. E então você pode executar o programa rodando `./programa`.

### `http.HandlerFunc`

Anteriormente, vimos que precisamos implementar a interface `Handler` para criar um servidor. _Normalmente_ fazemos isso criando um estrutura (`struct`) e fazendo com que esta implemente a interface. No entanto, mesmo que o comum seja utilizar as _estruturas_ para armazenar dados, _nesse momento_ não armazenamos um estado, então não parece certo criar uma _estrutura_ para isso.

Usar a função [HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc) nos ajuda a resolver este problema.

> O tipo HandlerFunc é um adaptador que permite usar funções comuns como tratadores (_handlers_). Se *f* é uma função com a assinatura adequada, HandlerFunc\(f\) é um _Handler_ que chama *f*.

```go
type HandlerFunc func(ResponseWriter, *Request)
```

Então usamos essa construção para adaptar a função `JogadorServidor`, fazendo com que esteja de acordo com a interface `Handler`.

### `http.ListenAndServe(":5000"...)`

`ListenAndServe` recebe como parâmetro um número de porta para escutar em um _tratador_ (`Handler`). Se a porta já estiver sendo usada, será retornado um _erro_ (`error`) para que, usando um comando `if`, possamos capturar esse erro e informar o problema para o usuário.

O que vamos fazer agora é escrever _outro_ teste para nos forçar a fazer uma mudança para tentar nos afastar do valor predefinido.

## Escreva o teste primeiro

Vamos adicionar outro subteste aos nossos testes, que tenta obter a pontuação de um jogador diferente, o que causará um problema em nossa implementação que usa um código predefinido.

```go
t.Run("retornar resultado de Pedro", func(t *testing.T) {
    requisicao, _ := http.NewRequest(http.MethodGet, "/jogadores/Pedro", nil)
    resposta := httptest.NewRecorder()

    JogadorServidor(resposta, requisicao)

    recebido := resposta.Body.String()
    esperado := "10"

    if recebido != esperado {
        t.Errorf("recebido '%s', esperado '%s'", recebido, esperado)
    }
})
```

Você deve estar pensando

> Certamente precisamos de algum tipo de armazenamento para controlar qual jogador recebe qual pontuação. É estranho que os valores sejam predefinidos em nossos testes.

Lembre-se de que estamos apenas tentando dar os menores passos possíveis; e por isso estamos, nesse momento, tentando invalidar o valor da constante.

## Tente rodar o próximo teste

```text
=== RUN   TestObterJogadores
=== RUN   TestObterJogadores/retornar_resultado_de_Maria
=== RUN   TestObterJogadores/retornar_resultado_de_Pedro
    TestObterJogadores/retornar_resultado_de_Pedro: servidor_test.go:34: recebido '20', esperado '10'
--- FAIL: TestObterJogadores (0.00s)
    --- PASS: TestObterJogadores/retornar_resultado_de_Maria (0.00s)
    --- FAIL: TestObterJogadores/retornar_resultado_de_Pedro (0.00s)
```

## Escreva código suficiente para fazer passar

```go
func JogadorServidor(w http.ResponseWriter, r *http.Request) {
    jogador := r.URL.Path[len("/jogadores/"):]

    if jogador == "Maria" {
        fmt.Fprint(w, "20")
        return
    }

    if jogador == "Pedro" {
        fmt.Fprint(w, "10")
        return
    }
}
```

Este teste nos forçou a olhar para a URL da requisição e tomar uma decisão. Embora ainda estamos pensando em como armazenar os dados do jogador e as interfaces, na verdade o próximo passo a ser dado está relacionado ao _roteamento_ (_routing_).

Se tivéssemos começado com o código de armazenamento dos dados, a quantidade de alterações que precisaríamos fazer seria muito grande. **Este é um pequeno passo em relação ao nosso objetivo final e foi guiado pelos testes**.

Estamos resistindo, nesse momento, à tentação de usar alguma biblioteca de roteamento, e queremos apenas dar o menor passo para fazer nossos testes funcionarem.

`r.URL.Path` retorna o caminho da request, e então usamos a sintaxe de slice para obter a parte final, depois de `/jogadores/`. Não é o recomendado por não ser muito robusto, mas resolve o problema por enquanto.

## Refatorar

Podemos simplificar a `JogadorServidor` separando a parte de obtenção da pontuação em uma função.

```go
func JogadorServidor(w http.ResponseWriter, r *http.Request) {
    jogador := r.URL.Path[len("/jogadores/"):]

    fmt.Fprint(w, ObterPontuacaoJogador(jogador))
}

func ObterPontuacaoJogador(nome string) string {
    if nome == "Maria" {
        return "20"
    }

    if nome == "Pedro" {
        return "10"
    }

    return ""
}
```

E podemos eliminar as repetições de parte do código dos testes montando algumas funções auxiliares("_helpers_")

```go
func TestObterJogadores(t *testing.T) {
    t.Run("retornar resultado de Maria", func(t *testing.T) {
        requisicao := novaRequisicaoObterPontuacao("Maria")
        resposta := httptest.NewRecorder()

        JogadorServidor(resposta, requisicao)

        verificarCorpoRequisicao(t, resposta.Body.String(), "20")
    })

    t.Run("returns Pedro's score", func(t *testing.T) {
        requisicao := novaRequisicaoObterPontuacao("Pedro")
        resposta := httptest.NewRecorder()

        JogadorServidor(resposta, requisicao)

        verificarCorpoRequisicao(t, resposta.Body.String(), "10")
    })
}

func novaRequisicaoObterPontuacao(nome string) *http.Request {
    requisicao, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
    return requisicao
}

func verificarCorpoRequisicao(t *testing.T, recebido, esperado string) {
    t.Helper()
    if recebido != esperado {
        t.Errorf("corpo da requisição é inválido, obtive '%s' esperava '%s'", recebido, esperado)
    }
}
```

Ainda assim, ainda não estamos felizes. Não parece correto que o servidor saiba as pontuações.

Mas nossa refatoração nos mostra claramente o que fazer.

Nós movemos o cálculo de pontuação para fora do código principal que trata a requisição (_handler_) para uma função `ObterPontuacaoJogador`. Isso parece ser o lugar correto para isolar as responsabilidades usando interfaces.

Vamos alterar, em `servidor.go`, a função que refatoramos para ser uma interface

```go
type JogadorArmazenamento interface {
    ObterPontuacaoJogador(nome string) int
}
```

Para que o `JogadorServidor` consiga usar o `JogadorArmazenamento`, é necessário ter uma referência a ele. Agora nos parece o momento certo para alterar nossa arquitetura, e nosso `JogadorServidor` agora se torna uma estrutura (`struct`).

```go
type JogadorServidor struct {
    armazenamento JogadorArmazenamento
}
```

E agora, vamos implementar a interface do _tratador_ (`Handler`) adicionando um método à nossa nova estrutura `JogadorServidor` e adicionado neste método o código existente.

```go
func (js *JogadorServidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/jogadores/"):]
    fmt.Fprint(w, js.armazenamento.ObterPontuacaoJogador(jogador))
}
```

Outra alteração a fazer: agora usamos a `armazenamento.ObterPontuacaoJogador` para obter a pontuação, ao invés da função local definida anteriormente \(e que podemos remover\).

Abaixo, a listagem completa do servidor (arquivo `servidor.go`)

```go
type JogadorArmazenamento interface {
	ObterPontuacaoJogador(nome string) int
}

type JogadorServidor struct {
	armazenamento JogadorArmazenamento
}

func (js *JogadorServidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]
	fmt.Fprint(w, js.armazenamento.ObterPontuacaoJogador(jogador))
}
```

### Ajustar os problemas

Fizemos muitas mudanças, e sabemos que nossos testes não irão funcionar e a compilação deixou de funcionar nesse momento; mas relaxe, e deixe o compilador fazer o trabalho.

`./main.go:9:58: type JogadorServidor is not an expression`

Precisamos mudar os nossos testes, que agora devem criar uma nova instância de `JogadorServidor` e então chamar o método `ServeHTTP`.

```go
func TestObterJogadores(t *testing.T) {
    servidor := &JogadorServidor{}

    t.Run("returns Maria's score", func(t *testing.T) {
        request := novaRequisicaoObterPontuacao("Maria")
        response := httptest.NewRecorder()

        servidor.ServeHTTP(response, request)

        verificarCorpoRequisicao(t, response.Body.String(), "20")
    })

    t.Run("returns Pedro's score", func(t *testing.T) {
        request := novaRequisicaoObterPontuacao("Pedro")
        response := httptest.NewRecorder()

        servidor.ServeHTTP(response, request)

        verificarCorpoRequisicao(t, response.Body.String(), "10")
    })
}
```

Perceba que ainda não nos preocupamos, _por enquanto_, com o armazenamento dos dados, nós apenas queremos a compilação funcionando o quanto antes.

Você deve ter o hábito de priorizar, _sempre_, código que compila antes de ter código que passa nos testes.

Adicionando mais funcionalidades \(como códigos de esboço de armazenamento\) a um código que não ainda não compila, nos arriscamos a ter, potencialmente, _mais_ problemas de compilação.

Agora `main.go` não vai compilar pelas mesmas razões.

```go
func main() {
	servidor := &JogadorServidor{}

	if err := http.ListenAndServe(":5000", servidor); err != nil {
		log.Fatalf("não foi possível escutar na porta 5000 %v", err)
	}
}
```

Agora tudo compila, mas os testes falham.

```text
--- FAIL: TestObterJogadores (0.00s)
    --- FAIL: TestObterJogadores/retorna_pontucao_de_Maria (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
        panic: runtime error: invalid memory address or nil pointer dereference
```

Isso porque não passamos um `JogadorArmazenamento` em nossos testes. Precisamos fazer, no arquivo `servidor_test.go` um código de esboço para nos ajudar.

```go
type EsbocoJogadorArmazenamento struct {
	pontuacoes map[string]int
}

func (e *EsbocoJogadorArmazenamento) ObterPontuacaoJogador(name string) int {
	pontuacao := e.pontuacoes[name]
	return pontuacao
}
```

Um _mapa_ (`map`) é um jeito simples e rápido de fazer um armazenamento chave/valor para os nossos testes. Agora vamos criar um desses armazenamentos para os nosso testes e inserir em nosso `JogadorServidor`.

```go
func TestObterJogadores(t *testing.T) {
	armazenamento := EsbocoJogadorArmazenamento{
		map[string]int{
			"Maria": 20,
			"Pedro": 10,
		},
	}
	servidor := &JogadorServidor{&armazenamento}

	t.Run("retorna pontuacao de Maria", func(t *testing.T) {
		request := novaRequisicaoObterPontuacao("Maria")
		response := httptest.NewRecorder()

		servidor.ServeHTTP(response, request)

		verificarCorpoRequisicao(t, response.Body.String(), "20")
	})

	t.Run("retorna pontuacao de Pedro", func(t *testing.T) {
		request := novaRequisicaoObterPontuacao("Pedro")
		response := httptest.NewRecorder()

		servidor.ServeHTTP(response, request)

		verificarCorpoRequisicao(t, response.Body.String(), "10")
	})
}
```

Nossos testes agora passam, e parecem melhores. Agora a _intenção_ do nosso código é clara, por conta da adição do armazenamento. Estamos dizendo a quem lê o código que, por termos _este dado em um `JogadorArmazenamento`_, quando você o usar com um  `JogadorServidor` você deve obter as respostas definidas.

### Rodar a aplicação

Agora que nossos testes estão passando, a última coisa que precisamos fazer para completar a refatoração é verificar se a aplicação está funcionando. O programa deve iniciar, mas você vai receber uma mensagem horrível se tentar acessar o servidor em `http://localhost:5000/jogadores/Maria`.

E a razão pra isso é: não informamos um `JogadorArmazenamento`.

Precisamos fazer uma implementação de um, mas isso é difícil no momento, já que não estamos armazenando nenhum dado significativo, por isso precisará ser, por enquanto, um valor predefinido. Vamos alterar na `main.go`:

```go
type JogadorArmazenamentoNaMemoria struct{}

func (i *JogadorArmazenamentoNaMemoria) ObterPontuacaoJogador(name string) int {
    return 123
}

func main() {
    server := &JogadorServidor{&JogadorArmazenamentoNaMemoria{}}

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("não foi possível escutar na porta 5000 %v", err)
    }
}
```

Se você rodar novamente o `go build` e acessar a mesma URL você deve receber `"123"`. Não é fantástico, mas até armazenarmos os dados, é o melhor que podemos fazer.

Temos algumas opções para decidir o que fazer agora

* Tratar o cenário onde o jogador não existe
* Tratar o cenário da chamado ao método HTTP `POST` em `/jogadores/{nome}`
* Não foi particularmente interessante perceber que nossa aplicação principal iniciou mas não funcionou. Tivemos que testar manualmente para ver o problema.

Enquanto o cenário do tratamento ao método HTTP `POST` nos deixa mais perto do "caminho ideal", eu sinto que vai ser mais fácil atacar o cenário de "jogador não existente" antes, já que estamos neste assunto. Veremos os outros itens posteriormente.

## Escreva o teste primeiro

Adicione o cenário de um jogador inexistente aos nossos testes

```go
t.Run("retorna 404 para jogador não encontrado", func(t *testing.T) {
    requisicao := novaRequisicaoObterPontuacao("Jorge")
    resposta := httptest.NewRecorder()

    server.ServeHTTP(resposta, requisicao)

    recebido := resposta.Code
    esperado := http.StatusNotFound

    if recebido != esperado {
        t.Errorf("recebido status %d esperado %d", recebido, esperado)
    }
})
```

## Tente rodar o teste

```text
--- FAIL: TestObterJogadores (0.00s)
    --- FAIL: TestObterJogadores/retorna_404_para_jogador_não_encontrado (0.00s)
        servidor_test.go:56: recebido status 200 esperado 404
```

## Escreva código necessário para que o teste funcione

```go
func (p *JogadorServidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/jogadores/"):]

    w.WriteHeader(http.StatusNotFound)

    fmt.Fprint(w, p.store.ObterPontuacaoJogador(player))
}
```

Às vezes eu me incomodo quando os defensores do Desenvolvimento Orientado a Testes dizem "tenha certeza de você escreveu apenas a mínima quantidade de código para fazer o teste funcionar", porque isso me parece muito pedante.

Mas este cenário ilustra muito bem o que querem dizer. Eu fiz o mínimo \(sabendo que não era a implementação correta\), que foi retornar um `StatusNotFound` em **todas as respostas**, mas todos os nossos testes estão passando!

**Implementando o mínimo para que os testes passem vai evidenciar as lacunas nos testes**. Em nosso caso, nós não estamos validando que devemos receber um `StatusOK` quando jogadores _existem_ em nosso armazenamento.

Atualize os outros dois testes para validar o retorno e corrija o código.

Eis os novos testes

```go
func TestObterJogadores(t *testing.T) {
	armazenamento := EsbocoJogadorArmazenamento{
		map[string]int{
			"Maria": 20,
			"Pedro": 10,
		},
	}
	servidor := &JogadorServidor{&armazenamento}

	t.Run("retorna pontuacao de Maria", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Maria")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		verificarCorpoRequisicao(t, resposta.Body.String(), "20")
	})

	t.Run("retorna pontuacao de Pedro", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Pedro")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		verificarCorpoRequisicao(t, resposta.Body.String(), "10")
	})

	t.Run("retorna 404 para jogador não encontrado", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Jorge")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		recebido := resposta.Code
		esperado := http.StatusNotFound

		if recebido != esperado {
			t.Errorf("recebido status %d esperado %d", recebido, esperado)
		}
	})
}

func novaRequisicaoObterPontuacao(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func verificarCorpoRequisicao(t *testing.T, recebido, esperado string) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("corpo da requisição é inválido, recebido '%s' esperado '%s'", recebido, esperado)
	}
}

func verificarRespostaCodigoStatus(t *testing.T, recebido, esperado int) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("não recebeu código de status HTTP esperado, recebido %d, esperado %d", recebido, esperado)
	}
}
```

Estamos verificando o `status` (código de retorno HTTP) em todos os nossos testes, por isso existe a função auxiliar `verificarRespostaCodigoStatus` para ajudar com isso.

Agora os primeiros dois testes falham porque o código de status recebido é 404, ao invés do esperado 200. Então vamos corrigir o `JogadorServidor` para que retorne *não encontrado* (HTTP status 404) se a pontuação for 0.

```go
func (js *JogadorServidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]

	pontuacao := js.armazenamento.ObterPontuacaoJogador(jogador)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}
```

### Armazenando pontuações

Agora que podemos obter pontuações de um armazenamento, também podemos armazenar novas pontuações.

## Escreva os testes primeiro

```go
func TestArmazenamentoVitorias(t *testing.T) {
	armazenamento := EsbocoJogadorArmazenamento{
		map[string]int{},
	}
	servidor := &JogadorServidor{&armazenamento}

	t.Run("retorna status 'aceito' para chamadas ao método POST", func(t *testing.T) {
		requisicao, _ := http.NewRequest(http.MethodPost, "/jogadores/Maria", nil)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusAccepted)
	})
}
```

Inicialmente vamos verificar se obtemos o código de status HTTP correto ao fazer a requisição em uma rota específica usando o método POST. Isso nos permite preparar o caminho da funcionalidade que aceita um tipo diferente de requisição, e tratar de forma diferente a requisição para `GET /jogadores/{nome}`. Uma vez que isso funcione como esperado, podemos começar a testar a interação do nosso tratador (_handler_) com o armazenamento.

## Tente rodar o teste

```text
--- FAIL: TestArmazenamentoVitorias (0.00s)
    --- FAIL: TestArmazenamentoVitorias/retorna_status_'aceito'_para_chamadas_ao_método_POST (0.00s)
        servidor_test.go:75: não recebeu código de status HTTP esperado, recebido 404, esperado 202
```

## Escreva código suficiente pra fazer passar

Lembre-se que estamos cometendo pecados deliberadamente, então um comando `if` para identificar o método da requisição vai resolver o problema.

```go
func (js *JogadorServidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	jogador := r.URL.Path[len("/jogadores/"):]

	pontuacao := js.armazenamento.ObterPontuacaoJogador(jogador)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}
```

## Refatorar

O tratador parece um pouco bagunçado agora. Vamos separar o código para ficar simples de entender e isolar as diferentes funcionalidades em novas funções.

```go
func (js *JogadorServidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		js.registrarVitoria(w)
	case http.MethodGet:
		js.mostrarPontuacao(w, r)
	}
}

func (p *JogadorServidor) mostrarPontuacao(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/jogadores/"):]

	score := p.armazenamento.ObterPontuacaoJogador(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *JogadorServidor) registrarVitoria(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}
```

Isso faz com que a responsabilidade de roteamento do `ServeHTTP` esteja mais clara; e também permite que, em nossas próximas iterações, o código para armazenamento possa estar dentro de `registrarVitoria`.

Agora, queremos verificar que, quando fazemos a chamada `POST` a `/jogadores/{nome}`, nosso `JogadorArmazenamento` registra a vitória.

## Escreva primeiro o teste

Vamos implementar isso estendendo o `EsbocoJogadorArmazenamento` com um novo método `RecordWin` e então inspecionar as chamadas.

```go
type EsbocoJogadorArmazenamento struct {
	pontuacoes        map[string]int
	registrosVitorias []string
}

func (e *EsbocoJogadorArmazenamento) ObterPontuacaoJogador(nome string) int {
	pontuacao := e.pontuacoes[nome]
	return pontuacao
}

func (e *EsbocoJogadorArmazenamento) RegistrarVitoria(nome string) {
	e.registrosVitorias = append(e.registrosVitorias, nome)
}
```

Agora, para começar, estendemos o teste para verificar a quantidade de chamadas

```go
func TestArmazenamentoVitorias(t *testing.T) {
	armazenamento := EsbocoJogadorArmazenamento{
		map[string]int{},
		nil,
	}
	servidor := &JogadorServidor{&armazenamento}

	t.Run("registra vitorias na chamada ao método HTTP POST", func(t *testing.T) {
		requisicao := novaRequisicaoRegistrarVitoriaPost("Maria")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusAccepted)

		if len(armazenamento.registrosVitorias) != 1 {
			t.Errorf("verifiquei %d chamadas a RegistrarVitoria, esperava %d", len(armazenamento.registrosVitorias), 1)
		}
	})
}

func novaRequisicaoRegistrarVitoriaPost(nome string) *http.Request {
	requisicao, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/jogadores/%s", nome), nil)
	return requisicao
}
```

## Tente rodar o teste

```text
./servidor_test.go:26:17: too few values in EsbocoJogadorArmazenamento literal
./servidor_test.go:70:17: too few values in EsbocoJogadorArmazenamento literal
```

## Escreva a mínima quantidade de código para a execução do teste e verifique a falha indicada no retorno

Como adicionamos um campo, precisamos atualizar o código onde criamos o `EsbocoJogadorArmazenamento`

```go
armazenamento := EsbocoJogadorArmazenamento{
    map[string]int{},
    nil,
}
```

```text
--- FAIL: TestArmazenamentoVitorias (0.00s)
    --- FAIL: TestArmazenamentoVitorias/#00 (0.00s)
        servidor_test.go:85: verifiquei 0 chamadas a RegistrarVitoria, esperava 1
```

## Escreva código suficiente para o teste passar

Como estamos apenas verificando o número de chamadas, e não seus valores específicos, nossa iteração inicial é um pouco menor.

Para conseguir invocar a `RegistrarVitoria`, precisamos atualizar a definição de `JogadorArmazenamento` para que o `JogadorServidor` funcione como esperado.

```go
type JogadorArmazenamento interface {
    ObterPontuacaoJogador(name string) int
    RegistrarVitoria(name string)
}
```

E, ao fazer isso, `main` não compila mais

```text
./main.go:15:29: cannot use &JogadorArmazenamentoNaMemoria literal (type *JogadorArmazenamentoNaMemoria) as type JogadorArmazenamento in field value:
        *JogadorArmazenamentoNaMemoria does not implement JogadorArmazenamento (missing RegistrarVitoria method)
```

O compilador nos informa o que está errado. Vamos alterar `InMemoryJogadorArmazenamento`, adicionando esse método.

```go
type JogadorArmazenamentoNaMemoria struct{}

func (ja *JogadorArmazenamentoNaMemoria) RegistrarVitoria(nome string) {}
```

Com essa alteração, o código volta a compilar - mas os testes ainda falham.

Agora que `JogadorArmazenamento` tem o método `RecordWin`, podemos chamar de dentro do nosso `JogadorServidor`

```go
func (p *JogadorServidor) processWin(w http.ResponseWriter) {
    p.store.RecordWin("Marcela")
    w.WriteHeader(http.StatusAccepted)
}
```

Rode os testes e deve estar funcionando sem erros! Claro, `"Marcela"` não é bem o que queremos enviar para `RegistrarVitoria`, então vamos ajustar os testes.

## Escreva os testes primeiro

```go
t.Run("registra vitorias na chamada ao método HTTP POST", func(t *testing.T) {
    jogador := "Maria"

    requisicao := novaRequisicaoRegistrarVitoriaPost(jogador)
    resposta := httptest.NewRecorder()

    servidor.ServeHTTP(resposta, requisicao)

    verificarRespostaCodigoStatus(t, resposta.Code, http.StatusAccepted)

    if len(armazenamento.registrosVitorias) != 1 {
        t.Errorf("verifiquei %d chamadas a RegistrarVitoria, esperava %d", len(armazenamento.registrosVitorias), 1)
    }

    if armazenamento.registrosVitorias[0] != jogador {
        t.Errorf("não registrou o vencedor corretamente, recebi '%s', esperava '%s'", armazenamento.registrosVitorias[0], jogador)
    }
})
```

Agora sabemos que existe um elemento no slice `registrosVitorias`, e então podemos acessar, sem erros, o primeiro elemento e verificar se é igual a `jogador`.

## Tente rodar o teste

```text
--- FAIL: TestArmazenamentoVitorias (0.00s)
    --- FAIL: TestArmazenamentoVitorias/registra_vitorias_na_chamada_ao_método_HTTP_POST (0.00s)
        servidor_test.go:91: não registrou o vencedor corretamente, recebi 'Marcela', esperava 'Maria'
```

## Escreva código suficiente para o teste passar

```go
func (js *JogadorServidor) registrarVitoria(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]
	js.armazenamento.RegistrarVitoria(jogador)
	w.WriteHeader(http.StatusAccepted)
}
```

Mudamos `registrarVitoria` para obter a `http.Request`, e assim conseguir extrair o nome do jogador da URL. Com o nome, podemos chamar o `armazenamento` com o valor correto para fazer os testes passarem.

## Refatorar

Podemos eliminar repetições no código, porque estamos obtendo o nome do "player" do mesmo jeito em dois lugares diferentes.

```go
func (p *JogadorServidor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/jogadores/"):]

    switch r.Method {
    case http.MethodPost:
        p.processWin(w, player)
    case http.MethodGet:
        p.showScore(w, player)
    }
}

func (p *JogadorServidor) showScore(w http.ResponseWriter, player string) {
    score := p.store.ObterPontuacaoJogador(player)

    if score == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, score)
}

func (p *JogadorServidor) processWin(w http.ResponseWriter, player string) {
    p.store.RecordWin(player)
    w.WriteHeader(http.StatusAccepted)
}
```

Mesmo com os testes passando, não temos código funcionando de forma ideal. Se executar a `main` e usar o programa como planejado, não vai funcionar porque ainda não nos dedicamos a implementar corretamente `JogadorArmazenamento`. Mas isso não é um problema; como focamos no tratamento da requisição, identificamos a interface necessária, ao invés de tentar definir antecipadamente.

_Poderíamos_ começar a escrever alguns testes para a `InMemoryJogadorArmazenamento`, mas ela é apenas uma solução temporária até a implementação de um modo mais robusto de registrar as pontuações \(por exemplo, em um banco de dados\).

O que vamos fazer agora é escrever um _teste de integração_ entre `JogadorServidor` e `InMemoryJogadorArmazenamento` para terminar a funcionalidade. Isso vai permitir confiar que a aplicação está funcionando, sem ter que testar diretamente `InMemoryJogadorArmazenamento`. E não apenas isso, mas quando implementarmos `JogadorArmazenamento` com um banco de dados, usaremos esse mesmo teste para verificar se a implementação funciona como esperado.

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
    store := InMemoryJogadorArmazenamento{}
    server := JogadorServidor{&store}
    player := "Maria"

    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

    response := httptest.NewRecorder()
    server.ServeHTTP(response, newGetScoreRequest(player))
    assertStatus(t, response.Code, http.StatusOK)

    assertResponseBody(t, response.Body.String(), "3")
}
```

* Estamos criando os dois componentes que queremos integrar: `InMemoryJogadorArmazenamento` e `JogadorServidor`.
* Então fazemos 3 requisições para registrar 3 vitórias para `player`. Não nos preocupamos com os códigos de retorno no teste, porque isso não é relevante para verificar se a integração funciona como esperado.
* Registramos a próxima resposta \(por isso guardamos o valor em `response`\) porque vamos obter a pontuação do `player`.

## Tente rodar o teste

```text
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:24: response body is wrong, got '123' want '3'
```

## Escreva código suficiente para passar

Abaixo, há mais código do que o esperado para se escrever sem ter os testes correspondentes.

_Isso é permitido_! Ainda existem testes verificando se as coisas estão funcionando como esperado, mas não focando na parte específica em que estamos trabalhando \(`InMemoryJogadorArmazenamento`\).

Se houvesse algum problema para continuarmos, era só reverter as alterações para antes do teste que falhou e então escrever mais testes unitários específicos para `InMemoryJogadorArmazenamento`, que nos ajudariam a encontrar a solução.

```go
func NewInMemoryJogadorArmazenamento() *InMemoryJogadorArmazenamento {
    return &InMemoryJogadorArmazenamento{map[string]int{}}
}

type InMemoryJogadorArmazenamento struct{
    store map[string]int
}

func (i *InMemoryJogadorArmazenamento) RecordWin(name string) {
    i.store[name]++
}

func (i *InMemoryJogadorArmazenamento) ObterPontuacaoJogador(name string) int {
    return i.store[name]
}
```

* Para armazenar os dados, adicionamos um `map[string]int` na struct `InMemoryJogadorArmazenamento`
* Para ajudar nos testes, criamos a `NewInMemoryJogadorArmazenamento` para inicializar o armazenamento, e o código do teste de integração foi atualizado para usar esta função \(`store := NewInMemoryJogadorArmazenamento()`\).
* O resto do código é apenas para fazer o `map` funcionar.

Nosso teste de integração passa, e agora só é preciso mudar o `main` para usar o `NewInMemoryJogadorArmazenamento()`

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    server := &JogadorServidor{NewInMemoryJogadorArmazenamento()}

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

Após compilar e rodar, use o `curl` para testar.

* Execute o comando a seguir algumas vezes, mude o nome do jogador se quiser `curl -X POST http://localhost:5000/jogadores/Maria`
* Verifique a pontuação, rodando `curl http://localhost:5000/jogadores/Maria`

Ótimo! Criamos um serviço de acordo com os padrões REST! Se quiser continuar, você pode escolher um armazenamento de dados com maior persistência, que não vai perder os dados quando o programa terminar.

* Escolher uma tecnologia de armazenamento \(Bolt? Mongo? Postgres? Sistema de arquivos?\)
* Fazer `PostgresJogadorArmazenamento` implementar `JogadorArmazenamento`
* Desenvolver a funcionalidade usando Desenvolvimento Orientado a Testes para ter certeza de que funciona
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
* o Desenvolvimento Orientado a Testes nos ajudou a definir as interfaces necessárias

### Cometa pecados, e daí refatore \(e então registre no controle de versão\)

* Você precisa tratar falhas na compilação ou nos testes como uma situação urgente, a qual precisa resolver o mais rápido possível.
* Escreva apenas o código necessário para resolver o problema. _Logo depois_ refatore e faça um código melhor
* Ao tentar fazer muitas alterações enquanto o código não está compilando ou os testes estão falhando, corremos o risco de acumular e agravar os problemas.
* Nos manter fiéis à essa abordagem nos obriga a escrever pequenos testes, o que significa pequenas alterações, o que nos ajuda a continuar trabalhando em sistemas complexos de forma gerenciável.
