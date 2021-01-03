# JSON, roteamento e embedding

[**Você pode encontrar todo o código para este capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/master/criando-uma-aplicacao/json)

[No capítulo anterior](../servidor-http/servidor-http.md) nós criamos um servidor web para armazenar quantos jogos nossos jogadores venceram.

Nossa gerente de produtos veio com um novo requisito;  criar um novo endpoint chamado `/liga` que retorne uma lista contendo todos os jogadores armazenados. Ela gostaria que isto fosse retornado como um JSON. 

## Este é o código que temos até agora

```go
// servidor.go
package main

import (
    "fmt"
    "net/http"
)

type ArmazenamentoJogador interface {
    ObtemPontuacaoDoJogador(nome string) int
    GravarVitoria(nome string)
}

type ServidorJogador struct {
    armazenamento ArmazenamentoJogador
}

func (s *ServidorJogador) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    jogador := r.URL.Path[len("/jogadores/"):]

    switch r.Method {
    case http.MethodPost:
        s.processarVitoria(w, jogador)
    case http.MethodGet:
        s.mostrarPontuacao(w, jogador)
    }
}

func (s *ServidorJogador) mostrarPontuacao(w http.ResponseWriter, jogador string) {
    pontuação := s.armazenamento.ObtemPontuacaoDoJogador(jogador)

    if pontuação == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, pontuação)
}

func (s *ServidorJogador) processarVitoria(w http.ResponseWriter, jogador string) {
    s.armazenamento.GravarVitoria(jogador)
    w.WriteHeader(http.StatusAccepted)
}
```

```go
// ArmazenamentoDeJogadorNaMemoria.go
package main

func NovoArmazenamentoDeJogadorNaMemoria() *ArmazenamentoDeJogadorNaMemoria {
    return &ArmazenamentoDeJogadorNaMemoria{map[string]int{}}
}

type ArmazenamentoDeJogadorNaMemoria struct {
    armazenamento map[string]int
}

func (a *ArmazenamentoDeJogadorNaMemoria) GravarVitoria(nome string) {
    a.armazenamento[nome]++
}

func (a *ArmazenamentoDeJogadorNaMemoria) ObtemPontuacaoDoJogador(nome string) int {
    return a.armazenamento[nome]
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
    servidor := &ServidorJogador{NovoArmazenamentoDeJogadorNaMemoria()}

    if err := http.ListenAndServe(":5000", servidor); err != nil {
        log.Fatalf("não foi possível ouvir na porta 5000 %v", err)
    }
}
```

Você pode encontrar os testes correspondentes no endereço no topo do capítulo.

Nós vamos começar criando o endpoint para a tabela de `liga`.

## Escreva os testes primeiro

Ampliaremos a suite de testes existente, pois temos algumas funções de teste úteis e um `ArmazenamentoJogador` falso para usar.

```go
// server_test.go

func TestLiga(t *testing.T) {
    armazenamento := EsbocoArmazenamentoJogador{}
    servidor := &ServidorJogador{&armazenamento}

    t.Run("retorna 200 em /liga", func(t *testing.T) {
        requisicao, _ := http.NewRequest(http.MethodGet, "/liga", nil)
        resposta := httptest.NewRecorder()

        servidor.ServeHTTP(resposta, requisicao)

        verificaStatus(t, resposta.Code, http.StatusOK)
    })
}
```
Antes de nos preocuparmos sobre as pontuações atuais e o JSON, nós vamos tentar manter as mudanças pequenas com o plano de ir passo a passo rumo ao nosso objetivo. O início mais simples é checar se nós conseguimos consultar `/liga` e obter um `OK` de retorno. 

## Tente rodar os testes

```text
=== RUN   TestLiga/retorna_200_em_/liga
panic: runtime error: slice bounds out of range [recovered]
    panic: runtime error: slice bounds out of range

goroutine 6 [running]:
testing.tRunner.func1(0xc42010c3c0)
    /usr/local/Cellar/go/1.10/libexec/src/testing/testing.go:742 +0x29d
panic(0x1274d60, 0x1438240)
    /usr/local/Cellar/go/1.10/libexec/src/runtime/panic.go:505 +0x229
github.com/larien/aprenda-go-com-testes/json-and-io/v2.(*ServidorJogador).ServeHTTP(0xc420048d30, 0x12fc1c0, 0xc420010940, 0xc420116000)
    /Users/larien/go/src/github.com/larien/aprenda-go-com-testes/json-and-io/v2/servidor.go:20 +0xec
```

Seu `ServidorJogador` deve estar sendo abortado por um panic como acima. Vá para a linha de código que está apontando para `servidor.go` no stack trace.

```go
jogador := r.URL.Path[len("/jogadores/"):]
```

No capítulo anterior, nós mencionamos que esta era uma maneira bastante ingênua de fazer o nosso roteamento. O que está acontecendo é que ele está tentando cortar a string do caminho da URL começando do índice após `/liga` e então, isto nos dá um `slice bounds out of range`.

## Escreva somente o código suficiente para fazê-lo passar

Go tem um mecanismo de rotas nativo (built-in) chamado [`ServeMux`](https://golang.org/pkg/net/http/#ServeMux) \(requisição multiplexadora\) que nos permite atracar um `http.Handler` para caminhos de uma requisição em específico.

Vamos cometer alguns pecados e obter os testes passando da maneira mais rápida que pudermos, sabendo que nós podemos refatorar isto com segurança uma vez que nós soubermos que os testes estão passando.

```go
func (s *ServidorJogador) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    roteador := http.NewServeMux()

    roteador.Handle("/liga", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    roteador.Handle("/jogadores/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        jogador := r.URL.Path[len("/jogadores/"):]

        switch r.Method {
        case http.MethodPost:
            s.processarVitoria(w, jogador)
        case http.MethodGet:
            s.mostrarPontuacao(w, jogador)
        }
    }))

    roteador.ServeHTTP(w, r)
}
```

* Quando a requisição começa nós criamos um roteador e então dizemos para o caminho `x` usar o handler `y`.
* Então para nosso novo endpoint, nós usamos `http.HandlerFunc` e uma _função anônima_ para `w.WriteHeader(http.StatusOK)` quando `/liga` é requisitada para fazer nosso novo teste passar.
* Para a rota `/jogadores/` nós somente recortamos e colamos nosso código dentro de outro `http.HandlerFunc`.
* Finalmente, nós lidamos com a requisição que está vindo chamando nosso novo roteador `ServeHTTP` \(notou como `ServeMux` é _também_ um `http.Handler`?\)

## Refatorando

`ServeHTTP` parece um pouco grande, nós podemos separar as coisas um pouco refatorando nossos handlers em métodos separados.

```go
func (s *ServidorJogador) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    roteador := http.NewServeMux()
    roteador.Handle("/liga", http.HandlerFunc(s.manipulaLiga))
    roteador.Handle("/jogadores/", http.HandlerFunc(s.manipulaJogadores))

    roteador.ServeHTTP(w, r)
}

func (s *ServidorJogador) manipulaLiga(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func (s *ServidorJogador) manipulaJogadores(w http.ResponseWriter, r *http.Request) {
    jogador := r.URL.Path[len("/jogadores/"):]

    switch r.Method {
    case http.MethodPost:
        s.processarVitoria(w, jogador)
    case http.MethodGet:
        s.mostrarPontuacao(w, jogador)
    }
}
```

É um pouco estranho \(e ineficiente\) estar configurando um roteador quando uma requisição chegar e então chamá-lo. O que idealmente queremos fazer é uma função do tipo `NovoServidorJogador` que pegará nossas dependências e ao ser chamada, irá fazer a configuração única da criação do roteador. Desta forma, cada requisição pode usar somente uma instância do nosso roteador.

```go
type ServidorJogador struct {
    armazenamento  ArmazenamentoJogador
    roteador *http.ServeMux
}

func NovoServidorJogador(armazenamento ArmazenamentoJogador) *ServidorJogador {
    p := &ServidorJogador{
        armazenamento,
        http.NewServeMux(),
    }

    s.roteador.Handle("/liga", http.HandlerFunc(s.manipulaLiga))
    s.roteador.Handle("/jogadores/", http.HandlerFunc(s.manipulaJogadores))

    return s
}

func (s *ServidorJogador) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.roteador.ServeHTTP(w, r)
}
```

* `ServidorJogador` agora precisa armazenar um roteador.
* Nós movemos a criação do roteador para fora de `ServeHTTP` e colocamos dentro do nosso `NovoServidorJogador`, então isto só será feito uma vez, não por requisição.
* Você vai precisar atualizar todos os testes e código de produção onde nós costumávamos fazer `ServidorJogador{&armazenamento}` por `NovoServidorJogador(&armazenamento)`.

### Uma refatoração final

Tente mudar o código para o seguinte:

```go
type ServidorJogador struct {
    armazenamento  ArmazenamentoJogador
    http.Handler
}

func NovoServidorJogador(armazenamento ArmazenamentoJogador) *ServidorJogador {
    s := new(ServidorJogador)

    s.armazenamento = armazenamento

    roteador := http.NewServeMux()
    roteador.Handle("/liga", http.HandlerFunc(s.manipulaLiga))
    roteador.Handle("/jogadores/", http.HandlerFunc(s.manipulaJogadores))

    s.Handler = roteador

    return s
}
```

Finalmente, se certifique de que você **deletou** `func (s *ServidorJogador) ServeHTTP(w http.ResponseWriter, r *http.Request)` por não ser mais necessária!

## Incorporando

Nós mudamos a segunda propriedade de `ServidorJogador` removendo a propriedade nomeada `roteador http.ServeMux` e substituindo por `http.Handler`; isto é chamado de _incorporar_. 


> O Go não provê a noção típica de subclasses orientada por tipo, mas tem a habilidade de "emprestar" partes de uma implementação por incorporar tipos dentro de uma struct ou interface.

[Effective Go - Embedding](https://golang.org/doc/effective_go.html#embedding)

O que isto quer dizer é que nosso `ServidorJogador` agora tem todos os métodos que `http.Handler` têm, que é somente o `ServeHTTP`.

Para "preencher" o `http.Handler` nós atribuímos ele para o `roteador` que nós criamos em `NovoServidorJogador`. Nós podemos fazer isso porque `http.ServeMux` tem o método `ServeHTTP`.

Isto nos permite remover nosso próprio método `ServeHTTP`, pois nós já estamos expondo um via o tipo incorporado. 

Incorporamento é um recurso muito interessante da linguagem. Você pode usar isto com interfaces para compor novas interfaces.

```go
type Animal interface {
    Comedor
    Dormente
}
```

E você pode usar isto com tipos concretos também, não somente interfaces. Como você pode esperar, se você incorporar um tipo concreto você vai ter acesso a todos os seus métodos e campos públicos. 

### Alguma desvantagem?

Você deve ter cuidado ao incorporar tipos porque você vai expor todos os métodos e campos públicos do tipo que você incorporou. Em nosso caso, está tudo bem porque nós haviamos incorporado apenas a _interface_ que nós queremos expôr \(`http.Handler`\).

Se nós tivéssemos sido "preguiçosos" e incorporado `http.ServeMux` \(o tipo concreto\) por exemplo, também funcionaria _porém_ os usuários de `ServidorJogador` seriam capazes de adicionar novas rotas ao nosso servidor porque o método `Handle(path, handler)` seria público.

**Quando incorporamos tipos, realmente devemos pensar sobre qual o impacto que isto terá em nossa API pública**

Isto é um erro _muito_ comum de mau uso de incorporamento, que termina poluindo nossas APIs e expondo os métodos internos dos seus tipos incorporados.

Agora que nós reestruturamos nossa aplicação, nós podemos facilmente adicionar novas rotas e botar para funcionar nosso endpoint `/liga`. Agora precisamos fazê-lo retornar algumas informações úteis.

Nós devemos retornar um JSON semelhante a este:

```javascript
[
   {
      "Nome":"Bill",
      "Vitórias":10
   },
   {
      "Nome":"Alice",
      "Vitórias":15
   }
]
```

## Escreva o teste primeiro

Nós vamos começar tentando analizar a resposta dentro de algo mais significativo.

```go
func TestLiga(t *testing.T) {
    armazenamento := EsbocoArmazenamentoJogador{}
    servidor := NovoServidorJogador(&armazenamento)

    t.Run("retorna 200 em /liga", func(t *testing.T) {
        requisicao, _ := http.NewRequest(http.MethodGet, "/liga", nil)
        resposta := httptest.NewRecorder()

        servidor.ServeHTTP(resposta, requisicao)

        var obtido []Jogador

        err := json.NewDecoder(resposta.Body).Decode(&obtido)

        if err != nil {
            t.Fatalf ("Não foi possível fazer parse da resposta do servidor '%s' no slice de Jogador, '%v'", resposta.Body, err)
        }

        verificaStatus(t, resposta.Code, http.StatusOK)
    })
}
```
### Por que não testar o JSON como texto puro?

Você pode argumentar que um simples teste inicial poderia só comparar que o não foi possível ouvir na porta 5000 tem um particular texto em JSON.


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
var obtido []Jogador
err := json.NewDecoder(resposta.Body).Decode(&obtido)
```

Para analizar o JSON dentro de nosso modelo de dados nós criamos um `Decoder` do pacote `encoding/json` e então chamamos seu método `Decode`. Para criar um `Decoder` é necessário ler de um `io.Reader`, que em nosso caso é nossa própria resposta `Body`.

`Decode` pega o endereço da coisa que nós estamos tentando decodificar, e é por isso que nós declaramos um slice vazio de `Jogador` na linha anterior.

Esse processo de analisar um JSON pode falhar, então `Decode` pode retornar um `error`. Não há ponto de continuidade para o teste se isto acontecer, então nós checamos o erro e paramos o teste com `t.Fatalf`.
Note que nós exibimos o não foi possível ouvir na porta 5000 junto do erro, pois é importante para qualquer outra pessoa que esteja rodando os testes ver que o texto não pôde ser analisado.

## Tente rodar o teste

```text
=== RUN   TestLiga/retorna_200_em_/liga
    --- FAIL: TestLiga/retorna_200_em_/liga (0.00s)
        server_test.go:107: Não foi possível fazer parse da resposta do servidor '' no slice de Jogador, 'unexpected end of JSON input'
```

Nosso endpoint atualmente não retorna um corpo, então isso não pode ser analisado como JSON.

## Escreva código suficiente para fazê-lo passar

```go
func (s *ServidorJogador) manipulaLiga(w http.ResponseWriter, r *http.Request) {
    tabelaDaLiga := []Jogador{
        {"Chris", 20},
    }

    json.NewEncoder(w).Encode(tabelaDaLiga)

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

Seria legal introduzir uma separação de conceitos entre nosso handler e o trecho de obter o `tabelaDaLiga`. Como sabemos, nós não vamos codificar isso por agora.

```go
func (s *ServidorJogador) manipulaLiga(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(s.obterTabelaDaLiga())
    w.WriteHeader(http.StatusOK)
}

func (s *ServidorJogador) obterTabelaDaLiga() []Jogador{
    return []Jogador{
        {"Chris", 20},
    }
}
```

Mais adiante, nós vamos querer estender nossos testes para então podermos controlar exatamente qual dado nós queremos receber de volta.

## Escreva o teste primeiro

Nós podemos atualizar o teste para afirmar que a tabela das ligas contem alguns jogadores que nós vamos pôr em nossa loja.

Atualize `EsbocoArmazenamentoJogador` para permitir que ele armazene uma liga, que é apenas um slice de `Jogador`. Nós vamos armazenar nossos dados esperados lá.

```go
type EsbocoArmazenamentoJogador struct {
    pontuações   map[string]int
    chamadasDeVitoria []string
    liga []Jogador
}
```
Adiante, atualize nossos testes colocando alguns jogadores na propriedade da liga, para então afirmar que eles foram retornados do nosso servidor.

```go
func TestLiga(t *testing.T) {

    t.Run("retorna a tabela da Liga como JSON", func(t *testing.T) {
        ligaEsperada := []Jogador{
            {"Cleo", 32},
            {"Chris", 20},
            {"Tiest", 14},
        }

        armazenamento := EsbocoArmazenamentoJogador{nil, nil, ligaEsperada}
        servidor := NovoServidorJogador(&armazenamento)

        requisicao, _ := http.NewRequest(http.MethodGet, "/liga", nil)
        resposta := httptest.NewRecorder()

        servidor.ServeHTTP(resposta, requisicao)

        var obtido []Jogador

        err := json.NewDecoder(resposta.Body).Decode(&obtido)

        if err != nil {
            t.Fatalf("Não foi possível fazer parse da resposta do servidor '%s' no slice de Jogador, '%v'", resposta.Body, err)
        }

        verificaStatus(t, resposta.Code, http.StatusOK)

        if !reflect.DeepEqual(obtido, ligaEsperada) {
            t.Errorf("obtido %v esperado %v", obtido, ligaEsperada)
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

Você vai precisar atualizar os outros testes, assim como nós temos um novo campo em `EsbocoArmazenamentoJogador`; ponha-o como nulo para os outros testes.

Tente executar os testes novamente e você deverá ter:

```text
=== RUN   TestLiga/retorna_a_tabela_da_liga_como_JSON
    --- FAIL: TestLiga/retorna_a_tabela_da_liga_como_JSON (0.00s)
        server_test.go:124: obtido [{Chris 20}] esperado [{Cleo 32} {Chris 20} {Tiest 14}]
```

## Escreva código suficiente para fazê-lo passar

Nós sabemos que o dado está em nosso `EsbocoArmazenamentoJogador` e nós abstraímos esses dados para uma interface `ArmazenamentoJogador`. Nós precisamos atualizar isto então qualquer um passando-nos um `ArmazenamentoJogador` pode prover-nos com dados para as ligas.

```go
type ArmazenamentoJogador interface {
    ObtemPontuacaoDoJogador(nome string) int
    GravarVitoria(nome string)
    ObterLiga() []Jogador
}
```

Agora nós podemos atualizar o código do nosso handler para chamar isto ao invés de retornar uma lista manualmente escrita. Delete nosso método `obterTabelaDaLiga()` e então atualize `manipulaLiga` para chamar `ObterLiga()`.

```go
func (s *ServidorJogador) manipulaLiga(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(s.armazenamento.ObterLiga())
    w.WriteHeader(http.StatusOK)
}
```

Tente executar os testes:

```text
# github.com/larien/aprenda-go-com-testes/json-and-io/v4
./main.go:9:50: cannot use NovoArmazenamentoDeJogadorNaMemoria() (type *ArmazenamentoDeJogadorNaMemoria) as type ArmazenamentoJogador in argument to NovoServidorJogador:
    *ArmazenamentoDeJogadorNaMemoria does not implement ArmazenamentoJogador (missing ObterLiga method)
./servidor_integration_test.go:11:27: cannot use armazenamento (type *ArmazenamentoDeJogadorNaMemoria) as type ArmazenamentoJogador in argument to NovoServidorJogador:
    *ArmazenamentoDeJogadorNaMemoria does not implement ArmazenamentoJogador (missing ObterLiga method)
./server_test.go:36:28: cannot use &armazenamento (type *EsbocoArmazenamentoJogador) as type ArmazenamentoJogador in argument to NovoServidorJogador:
    *EsbocoArmazenamentoJogador does not implement ArmazenamentoJogador (missing ObterLiga method)
./server_test.go:74:28: cannot use &armazenamento (type *EsbocoArmazenamentoJogador) as type ArmazenamentoJogador in argument to NovoServidorJogador:
    *EsbocoArmazenamentoJogador does not implement ArmazenamentoJogador (missing ObterLiga method)
./server_test.go:106:29: cannot use &armazenamento (type *EsbocoArmazenamentoJogador) as type ArmazenamentoJogador in argument to NovoServidorJogador:
    *EsbocoArmazenamentoJogador does not implement ArmazenamentoJogador (missing ObterLiga method)
```

O compilador está reclamando porque `ArmazenamentoDeJogadorNaMemoria` e `EsbocoArmazenamentoJogador` não tem os novos métodos que nós adicionamos em nossa interface.

Para `EsbocoArmazenamentoJogador` isto é bem fácil, apenas retorne o campo `liga` que nós adicionamos anteriormente.

```go
func (s *EsbocoArmazenamentoJogador) ObterLiga() []Jogador {
    return s.liga
}
```
Aqui está uma lembrança de como `InMemoryStore` é implementado:

```go
type ArmazenamentoDeJogadorNaMemoria struct {
    armazenamento map[string]int
}
```
Embora seja bastante simples para implementar `ObterLiga` "propriamente", iterando sobre o map, lembre que nós estamos apenas tentando _escrever o mínimo de código para fazer os testes passarem_.

Então vamos apenas deixar o compilador feliz por enquanto e viver com o desconfortável sentimento de uma implementação incompleta em nosso `InMemoryStore`.

```go
func (a *ArmazenamentoDeJogadorNaMemoria) ObterLiga() []Jogador {
    return nil
}
```

O que isto está realmente nos dizendo é que _depois_ nós vamos querer testar isto, porém vamos estacionar isto por hora.

Tente executar os testes, o compilador deve passar e os testes deverão estar passando!

## Refatoração

O código de teste não transmite suas intenções muito bem e possui vários trechos que podem ser refatorados.

```go
t.Run("retorna a tabela da Liga como JSON", func(t *testing.T) {
    ligaEsperada := []Jogador{
        {"Cleo", 32},
        {"Chris", 20},
        {"Tiest", 14},
    }

    armazenamento := EsbocoArmazenamentoJogador{nil, nil, ligaEsperada}
    servidor := NovoServidorJogador(&armazenamento)

    requisicao := novaRequisicaoDeLiga()
    resposta := httptest.NewRecorder()

    servidor.ServeHTTP(resposta, requisicao)

    obtido := obterLigaDaResposta(t, resposta.Body)
    verificaStatus(t, resposta.Code, http.StatusOK)
    verificaLiga(t, obtido, ligaEsperada)
})
```

Aqui estão os novos helpers:

```go
func obterLigaDaResposta(t *testing.T, body io.Reader) (liga []Jogador) {
    t.Helper()
    err := json.NewDecoder(body).Decode(&liga)

    if err != nil {
        t.Fatalf("Não foi possível fazer parse da resposta do servidor '%s' no slice de Jogador, '%v'", body, err)
    }

    return
}

func verificaLiga(t *testing.T, obtido, esperado []Jogador) {
    t.Helper()
    if !reflect.DeepEqual(obtido, esperado) {
        t.Errorf("obtido %v esperado %v", obtido, esperado)
    }
}

func novaRequisicaoDeLiga() *http.Request {
    req, _ := http.NewRequest(http.MethodGet, "/liga", nil)
    return req
}
```

Uma última coisa que nós precisamos fazer para nosso servidor funcionar é ter certeza de que nós retornamos um `content-type` correto na resposta, então as máquinas podem reconhecer que nós estamos retornando um `JSON`.

## Escreva os testes primeiro

Adicione essa afirmação no teste existente

```go
if resposta.Result().Header.Get("content-type") != "application/json" {
    t.Errorf("resposta não tinha o tipo de conteúdo de application/json, obtido %v", resposta.Result().Header)
}
```

## Tente rodar o teste

```text
=== RUN   TestLiga/retorna_a_tabela_da_liga_como_JSON
    --- FAIL: TestLiga/retorna_a_tabela_da_liga_como_JSON (0.00s)
        server_test.go:124: resposta não tinha o tipo de conteúdo de application/json, obtido map[Content-Type:[text/plain; charset=utf-8]]
```

## Escreva código suficiente para fazê-lo passar

Atualize `manipulaLiga`

```go
func (s *ServidorJogador) manipulaLiga(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    json.NewEncoder(w).Encode(s.armazenamento.ObterLiga())
}
```

O teste deve passar.

## Refatoração

Adicione um helper para `verificaTipoDoConteudo`.

```go
const tipoDoConteudoJSON = "application/json"

func verificaTipoDoConteudo(t *testing.T, resposta *httptest.ResponseRecorder, esperado string) {
    t.Helper()
    if resposta.Result().Header.Get("content-type") != esperado {
        t.Errorf("resposta não obteve content-type de %s, obtido %v", esperado, resposta.Result().Header)
    }
}
```

Use isso no teste.

```go
verificaTipoDoConteudo(t, resposta, tipoDoConteudoJSON)
```

Agora que nós resolvemos `ServidorJogador`, por agora podemos mudar nossa atenção para `ArmazenamentoDeJogadorNaMemoria` porque no momento se nós tentarmos demonstrá-lo para o gerente de produto, `/liga` não vai funcionar.

A forma mais rápida de nós termos alguma confiança é adicionar a nosso teste de integração, nós podemos bater no novo endpoint e checar se nós recebemos a resposta correta de `/liga`.

## Escreva o teste primeiro

Nós podemos usar `t.Run` para parar este teste um pouco e então reusar os helpers dos testes do nosso servidor - novamente mostrando a importância de refatoração dos testes.

```go
func TestGravaVitoriasEAsRetorna(t *testing.T) {
    armazenamento := NovoArmazenamentoDeJogadorNaMemoria()
    servidor := NovoServidorJogador(armazenamento)
    jogador := "Pepper"

    servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
    servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))
    servidor.ServeHTTP(httptest.NewRecorder(), novaRequisiçãoPostDeVitoria(jogador))

    t.Run("obter pontuação", func(t *testing.T) {
        resposta := httptest.NewRecorder()
        servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
        verificaStatus(t, resposta.Code, http.StatusOK)

        verificaCorpoDaResposta(t, resposta.Body.String(), "3")
    })

    t.Run("obter liga", func(t *testing.T) {
        resposta := httptest.NewRecorder()
        servidor.ServeHTTP(resposta, novaRequisicaoDeLiga())
        verificaStatus(t, resposta.Code, http.StatusOK)

        obtido := obterLigaDaResposta(t, resposta.Body)
        esperado := []Jogador{
            {"Pepper", 3},
        }
        verificaLiga(t, obtido, esperado)
    })
}
```

## Tente rodar o teste

```text
=== RUN   TestGravaVitoriasEAsRetorna/obter_liga
    --- FAIL: TestGravaVitoriasEAsRetorna/obter_liga (0.00s)
        servidor_integration_test.go:35: obtido [] esperado [{Pepper 3}]
```

## Escreva código suficiente para fazê-lo passar

`ArmazenamentoDeJogadorNaMemoria` is returning `nil` when you call `ObterLiga()` so we'll need to fix that.

```go
func (a *ArmazenamentoDeJogadorNaMemoria) ObterLiga() []Jogador {
    var liga []Jogador
    for nome, vitórias := range a.armazenamento {
        liga = append(liga, Jogador{nome, vitórias})
    }
    return liga
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
