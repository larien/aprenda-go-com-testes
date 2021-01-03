# Contexto

[**Você pode encontrar todo o código para esse capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/master/primeiros-passos-com-go/contexto)

Softwares geralmente iniciam processos de longa duração e de uso intensivo de recursos \(muitas vezes em goroutines\). Se a ação que causou isso é cancelada ou falha por algum motivo, você precisa parar esses processos de forma consistente dentro da sua aplicação.

Se você não gerenciar isso, sua aplicação Go tão ágil da qual você tem tanto orgulho pode começar a ter problemas de desempenho difíceis de depurar.

Neste capítulo vamos usar o pacote `context` para nos ajudar a gerenciar processos de longa duração.

Vamos começar com um exemplo clássico de um servidor web que, quando iniciado, abre um processo de longa execução que vai buscar alguns dados para devolver em uma resposta.

Colocaremos em prática um cenário em que um usuário cancela a requisição antes que os dados possam ser recuperados e faremos com que o processo seja instruído a desistir.

Criei um código no caminho feliz para começarmos. Aqui está o código do nosso servidor.

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, store.Fetch())
    }
}
```

A função `Server` (servidor) recebe uma `Store` (armazenamento) e nos retorna um `http.HandlerFunc`. Store está definida como:

```go
type Store interface {
    Fetch() string
}
```

A função retornada chama o método `Fetch` (busca) da `Store` para obter os dados e escrevê-los na resposta.

Nós temos um stub correspondente para `Store` que usamos em um teste.

```go
type StubStore struct {
    response string
}

func (s *StubStore) Fetch() string {
    return s.response
}

func TestHandler(t *testing.T) {
    data := "olá, mundo"
    svr := Server(&StubStore{data})

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`resultado "%s", esperado "%s"`, response.Body.String(), data)
    }
}
```

Agora que temos um caminho feliz, queremos fazer um cenário mais realista onde a `Store` não consiga finalizar o `Fetch` antes que o usuário cancele a requisição.

## Escreva o teste primeiro

Nosso handler precisará de uma maneira de dizer à `Store` para cancelar o trabalho, então atualize a interface.

```go
type Store interface {
    Fetch() string
    Cancel()
}
```

Precisaremos ajustar nosso spy para que leve algum tempo para retornar `data` e uma maneira de saber que foi dito para cancelar. Nós também o renomearemos para `SpyStore`, pois agora vamos observar a forma como ele é chamado. Ele terá que adicionar `Cancel` como um método para implementar a interface `Store`.

```go
type SpyStore struct {
    response string
    cancelled bool
}

func (s *SpyStore) Fetch() string {
    time.Sleep(100 * time.Millisecond)
    return s.response
}

func (s *SpyStore) Cancel() {
    s.cancelled = true
}
```

Vamos adicionar um novo teste onde cancelamos a requisição antes de 100 milissegundos e verificamos a store para ver se ela é cancelada.

```go
t.Run("avisa a store para cancelar o trabalho se a requisição for cancelada", func(t *testing.T) {
      store := &SpyStore{response: data}
      svr := Server(store)

      request := httptest.NewRequest(http.MethodGet, "/", nil)

      cancellingCtx, cancel := context.WithCancel(request.Context())
      time.AfterFunc(5 * time.Millisecond, cancel)
      request = request.WithContext(cancellingCtx)

      response := httptest.NewRecorder()

      svr.ServeHTTP(response, request)

      if !store.cancelled {
          t.Errorf("store não foi avisada para cancelar")
      }
  })
```

Do blog da Google novamente:

> O pacote context fornece funções para derivar novos valores de contexto dos já existentes. Estes valores formam uma árvore: quando um contexto é cancelado, todos os contextos derivados dele também são cancelados.

É importante que você derive seus contextos para que os cancelamentos sejam propagados através da pilha de chamadas ([_call stack_](https://www.ardanlabs.com/blog/2015/01/stack-traces-in-go.html)) para uma determinada requisição.

O que fazemos é derivar um novo `cancellingCtx` da nossa requisição que nos retorna uma função `cancel`. Nós então programamos que a função seja chamada em 5 milissegundos usando `time.AfterFunc`. Por fim, usamos este novo contexto em nossa requisição chamando `request.WithContext`.

## Execute o teste

O teste falha como seria de esperar.

```go
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/avisa_a_store_para_cancelar_o_trabalho_se_a_requisicao_for_cancelada (0.00s)
        context_test.go:62: store no foi avisada para cancelar
```

## Escreva código o suficiente para fazer o teste passar

Lembre-se de ser disciplinado com o TDD. Escreva a quantidade _mínima_ de código para fazer nosso teste passar.

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        store.Cancel()
        fmt.Fprint(w, store.Fetch())
    }
}
```

Isto faz com que este teste passe, mas não parece tão bom! Certamente não deveríamos estar cancelando a `Store` antes do fetch em _cada requisição_.

Ao ser disciplinado ele destacou uma falha em nossos testes, isso é uma coisa boa!

Vamos precisar atualizar nosso teste de caminho feliz para verificar que ele não será cancelado.

```go
t.Run("retorna dados da store", func(t *testing.T) {
    store := SpyStore{response: data}
    svr := Server(&store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`resultado "%s", esperado "%s"`, response.Body.String(), data)
    }

    if store.cancelled {
        t.Error("não deveria ter cancelado a store")
    }
})
```

Execute ambos os testes. O teste do caminho feliz deve agora estar falhando e somos forçados a fazer uma implementação mais sensata.

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        data := make(chan string, 1)

        go func() {
            data <- store.Fetch()
        }()

        select {
        case d := <-data:
            fmt.Fprint(w, d)
        case <-ctx.Done():
            store.Cancel()
        }
    }
}
```

O que fizemos aqui?

`context` tem um método `Done()` que retorna um canal que recebe um sinal quando o context estiver "done" (finalizado) ou "cancelled" (cancelado). Queremos ouvir esse sinal e chamar `store.Cancel` se o obtivermos, mas queremos ignorá-lo se a nossa `Store` conseguir finalizar o `Fetch` antes dele.

Para gerenciar isto, executamos o `Fetch` em uma goroutine e ele irá escrever o resultado em um novo channel `data`. Nós então usamos `select` para efetivamente correr para os dois processos assíncronos e então escrevemos uma resposta ou cancelamos com `Cancel`.

## Refatoração

Podemos refatorar um pouco o nosso código de teste fazendo métodos de verificação no nosso spy.

```go
func (s *SpyStore) assertWasCancelled() {
    s.t.Helper()
    if !s.cancelled {
        s.t.Errorf("store não foi avisada para cancelar")
    }
}

func (s *SpyStore) assertWasNotCancelled() {
    s.t.Helper()
    if s.cancelled {
        s.t.Errorf("store foi avisada para cancelar")
    }
}
```

Lembre-se de passar o `*testing.T` ao criar o spy.

```go
func TestServer(t *testing.T) {
    data := "olá, mundo"

    t.Run("retorna dados da store", func(t *testing.T) {
        store := &SpyStore{response: data, t: t}
        svr := Server(store)

        request := httptest.NewRequest(http.MethodGet, "/", nil)
        response := httptest.NewRecorder()

        svr.ServeHTTP(response, request)

        if response.Body.String() != data {
            t.Errorf(`recebi "%s", quero "%s"`, response.Body.String(), data)
        }

        store.assertWasNotCancelled()
    })

    t.Run("avisa a store para cancelar o trabalho se a requisição for cancelada", func(t *testing.T) {
        store := &SpyStore{response: data, t: t}
        svr := Server(store)

        request := httptest.NewRequest(http.MethodGet, "/", nil)

        cancellingCtx, cancel := context.WithCancel(request.Context())
        time.AfterFunc(5*time.Millisecond, cancel)
        request = request.WithContext(cancellingCtx)

        response := httptest.NewRecorder()

        svr.ServeHTTP(response, request)

        store.assertWasCancelled()
    })
}
```

Esta abordagem é boa, mas é idiomática?

Faz sentido para o nosso servidor web estar preocupado com o cancelamento manual da `Store`? E se a `Store` também depender de outros processos de execução lenta? Nós teremos que ter certeza que a `Store.Cancel` propagará corretamente o cancelamento para todos os seus dependentes.

Um dos pontos principais do `context` é que é uma maneira consistente de oferecer cancelamento.

[Da documentação do go](https://golang.org/pkg/context/)

> As requisições de entrada para um servidor devem criar um Context e as chamadas de saída para servidores devem aceitar um Context. A cadeia de chamadas de função entre eles deve propagar o Context, substituindo-o opcionalmente por um Context derivado criado usando WithCancel, WithDeadline, WithTimeout ou WithValue. Quando um Context é cancelado, todos os Contexts derivados dele também são cancelados.

Do blog da Google novamente:

> Na Google, exigimos que os programadores Go passem um parâmetro Context como o primeiro argumento para cada função no caminho de chamada entre requisições de entrada e saída. Isto permite que o código Go desenvolvido por muitas equipes diferentes interopere bem. Ele fornece um controle simples sobre timeouts e cancelamentos e garante que valores críticos, como credenciais de segurança, transitem corretamente pelos programas Go.

\(Pare por um momento e pense nas ramificações de cada função tendo que enviar um context e a ergonomia disso.\)

Se sentindo um pouco desconfortável? Bom. Vamos tentar seguir essa abordagem e, em vez disso, passar o `context` para nossa `Store` e deixá-la ser responsável. Dessa maneira, ela também pode passar o `context` para os seus dependentes e eles também podem ser responsáveis por se pararem.

## Escreva o teste primeiro

Teremos de alterar os nossos testes existentes, uma vez que as suas responsabilidades estão mudando. As únicas coisas que nosso handler é responsável agora é certificar-se que emite um contexto à `Store` em cascata (downstream) e que trata o erro que virá da `Store` quando é cancelada.

Vamos atualizar nossa interface `Store` para mostrar as novas responsabilidades.

```go
type Store interface {
    Fetch(ctx context.Context) (string, error)
}
```

Apague o código dentro do nosso handler por enquanto:

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    }
}
```

Atualize nosso `SpyStore`:

```go
type SpyStore struct {
    response string
    t        *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
    data := make(chan string, 1)

    go func() {
        var result string
        for _, c := range s.response {
            select {
            case <-ctx.Done():
                s.t.Log("spy store foi cancelado")
                return
            default:
                time.Sleep(10 * time.Millisecond)
                result += string(c)
            }
        }
        data <- result
    }()

    select {
    case <-ctx.Done():
        return "", ctx.Err()
    case res := <-data:
        return res, nil
    }
}
```

Temos que fazer nosso spy agir como um método real que funciona com o `context`.

Estamos simulando um processo lento onde construímos o resultado lentamente adicionando a string, caractere por caractere em uma goroutine. Quando a goroutine termina seu trabalho, ela escreve a string no channel `data`. A goroutine escuta o `ctx.Done` e irá parar o trabalho se um sinal for enviado nesse channel.

Finalmente o código usa outro `select` para esperar que a goroutine termine seu trabalho ou que o cancelamento ocorra.

É semelhante à nossa abordagem de antes onde usamos as primitivas de concorrência do Go para fazerem dois processos assíncronos disputarem um contra o outro para determinar o que retornamos.

Você usará uma abordagem similar ao escrever suas próprias funções e métodos que aceitam um `context`, por isso certifique-se de que está entendendo o que está acontecendo.

Nós removemos a referência ao `ctx` dos campos do `SpyStore` porque não é mais interessante para nós. Estamos estritamente testando o comportamento agora, que preferimos em comparação aos detalhes da implementação dos testes, como "você passou um determinado valor para a função `foo`".

Finalmente podemos atualizar nossos testes. Comente nosso teste de cancelamento para que possamos corrigir o teste do caminho feliz primeiro.

```go
t.Run("retorna dados da store", func(t *testing.T) {
    store := &SpyStore{response: data, t: t}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`resultado "%s", esperado "%s"`, response.Body.String(), data)
    }
})
```

## Execute o teste

```text
=== RUN   TestServer/retorna_dados_da_store
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/retorna_dados_da_store (0.00s)
        context_test.go:22: resultado "", esperado "olá, mundo"
```

## Escreva código o suficiente para fazer o teste passar

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data, _ := store.Fetch(r.Context())
        fmt.Fprint(w, data)
    }
}
```

O nosso caminho feliz deve estar... feliz. Agora podemos corrigir o outro teste.

## Escreva o teste primeiro

Precisamos testar que não escrevemos qualquer tipo de resposta no caso de erro. Infelizmente o `httptest.ResponseRecorder` não tem uma maneira de descobrir isso, então teremos que usar nosso próprio spy para testar.

```go
type SpyResponseWriter struct {
    written bool
}

func (s *SpyResponseWriter) Header() http.Header {
    s.written = true
    return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
    s.written = true
    return 0, errors.New("não implementado")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
    s.written = true
}
```

Nosso `SpyResponseWriter` implementa `http.ResponseWriter` para que possamos usá-lo no teste.

```go
t.Run("avisa a store para cancelar o trabalho se a requisição for cancelada", func(t *testing.T) {
    store := &SpyStore{response: data, t: t}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)

    cancellingCtx, cancel := context.WithCancel(request.Context())
    time.AfterFunc(5*time.Millisecond, cancel)
    request = request.WithContext(cancellingCtx)

    response := &SpyResponseWriter{}

    svr.ServeHTTP(response, request)

    if response.written {
        t.Error("uma resposta não deveria ter sido escrita")
    }
})
```

## Execute o teste

```text
=== RUN   TestServer
=== RUN   TestServer/avisa_a_store_para_cancelar_o_trabalho_se_a_requisicao_for_cancelada
--- FAIL: TestServer (0.01s)
    --- FAIL: TestServer/avisa_a_store_para_cancelar_o_trabalho_se_a_requisicao_for_cancelada (0.01s)
        context_test.go:47: uma resposta não deveria ter sido escrita
```

## Escreva código o suficiente para fazer o teste passar

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data, err := store.Fetch(r.Context())

        if err != nil {
            return // todo: registre o erro como você quiser
        }

        fmt.Fprint(w, data)
    }
}
```

Podemos ver depois disso que o código do servidor se tornou simplificado, pois não é mais explicitamente responsável pelo cancelamento. Ele simplesmente passa o `context` e confia nas funções em cascata (downstream) para respeitar qualquer cancelamento que possa ocorrer.

## Resumo

### Sobre o que falamos

* Como testar um handler HTTP que teve a requisição cancelada pelo cliente.
* Como usar o contexto para gerenciar o cancelamento.
* Como escrever uma função que aceita `context` e o usa para se cancelar usando goroutines, `select` e canais.
* Seguir as diretrizes da Google a respeito de como controlar o cancelamento propagando o contexto escopado da requisição através da sua pilha de chamadas (_call stack_).
* Como levar seu próprio spy para `http.ResponseWriter` se você precisar dele.

### E quanto ao context.Value?

[Michal Štrba](https://faiface.github.io/post/context-should-go-away-go2/) e eu temos uma opinião semelhante.

> Se você usar o ctx.Value na minha empresa \(inexistente\), você está demitido

Alguns engenheiros têm defendido a passagem de valores através do `context` porque _parece conveniente_.

A conveniência é muitas vezes a causa do código ruim.

O problema com `context.Values` é que ele é apenas um mapa não tipado para que você não tenha nenhum tipo de segurança e você tem que lidar com ele não realmente contendo seu valor. Você tem que criar um acoplamento de chaves de mapa de um módulo para outro e se alguém muda alguma coisa começar a quebrar.

Resumindo, **se uma função necessita de alguns valores, coloque-os como parâmetros tipados em vez de tentar obtê-los a partir de `context.Value`**. Isto torna-o estaticamente verificado e documentado para que todos o vejam.

#### Mas...

Por outro lado, pode ser útil incluir informações que sejam ortogonais a uma requisição em um contexto, como um identificador único. Potencialmente esta informação não seria necessária para todas as funções da sua pilha de chamadas (_call stack_) e tornaria as suas assinaturas funcionais muito confusas.

[Jack Lindamood diz que **Context.Value deve informar, não controlar**](https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39)

> O conteúdo do context.Value é para os mantenedores e não para os usuários. Ele nunca deve ser uma entrada necessária para resultados documentados ou esperados.

### Material adicional

* Gostei muito de ler [Context should go away for Go 2 por Michal Štrba](https://faiface.github.io/post/context-should-go-away-go2/). Seu argumento é que ter que passar o `context` em toda parte é um indicador que está apontando a uma deficiência na linguagem a respeito do cancelamento. Ele diz que seria melhor se isso fosse resolvido de alguma forma no nível de linguagem, em vez de em um nível de biblioteca. Até que isso aconteça, você precisará do `context` se quiser gerenciar processos de longa duração.
* O [blog do Go descreve ainda mais a motivação para trabalhar com `context` e tem alguns exemplos](https://blog.golang.org/context)

