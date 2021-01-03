# IO e Ordenação

[**Você pode encontrar todo o código para este capítulo aqui**](https://github.com/larien/aprenda-go-com-testes/tree/master/criando-uma-aplicacao/io)

[No capitulo anterior](../json.md) continuamos interagindo com nossa aplicação pela adição de um novo endpoint `/liga`. Durante o caminho aprendemos como lidar com JSON, tipos embutidos e roteamento.

Nossa dona do produto está de certa forma preocupada, por conta do software perder as pontuações quando o servidor é reiniciado. Ela também não se agradou que nós não interpretamos o endpoint `/liga` que deveria retornar os jogadores ordenados pelo número de vitórias!

## O código até agora

```go
// server.go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// GuardaJogador armazena informações sobre os jogadores
type GuardaJogador interface {
    PegaPontuacaoDoJogador(nome string) int
    SalvaVitoria(nome string)
    PegaLiga() []Jogador
}

// Jogador guarda o nome com o número de vitorias
type Jogador struct {
    Nome string
    Vitorias int
}

// ServidorDoJogador é uma interface HTTP para informações dos jogadores
type ServidorDoJogador struct {
    armazenamento GuardaJogador
    http.Handler
}

const jsonContentType = "application/json"

// NovoServidorDoJogador cria um ServidorDoJogador com roteamento configurado
func NovoServidorDoJogador(armazenamento GuardaJogador) *ServidorDoJogador {
    p := new( ServidorDoJogador)

    p.armazenamento = armazenamento

    roteador := http.NewServeMux()
    roteador.Handle("/liga", http.HandlerFunc(p.ManipulaLiga))
    roteador.Handle("/jogadores/", http.HandlerFunc(p.ManipulaJogador))

    p.Handler = roteador

    return p
}

func (p *ServidorDoJogador) ManipulaLiga(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(p.armazenamento.PegaLiga())
    w.Header().Set("content-type", jsonContentType)
    w.WriteHeader(http.StatusOK)
}

func (p *ServidorDoJogador) ManipulaJogador(w http.ResponseWriter, r *http.Request) {
    jogador := r.URL.Path[len("/jogadores/"):]

    switch r.Method {
    case http.MethodPost:
        p.processaVitoria(w, jogador)
    case http.MethodGet:
        p.mostraPontuacao(w, jogador)
    }
}

func (p *ServidorDoJogador) mostraPontuacao(w http.ResponseWriter, jogador string) {
    pontuacao := p.armazenamento.PegaPontuacaoDoJogador(jogador)

    if pontuacao == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, pontuacao)
}

func (p *ServidorDoJogador) processaVitoria(w http.ResponseWriter, jogador string) {
    p.armazenamento.salvaVitorias(jogador)
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

func (i *ArmazenamentoDeJogadorNaMemoria) PegaLiga() []Jogador {
    var liga []Jogador
    for nome, vitorias := range i.armazenamento {
        liga = append(liga, Jogador{nome, vitorias})
    }
    return liga
}

func (i *ArmazenamentoDeJogadorNaMemoria) SalvaVitoria(nome string) {
    i.armazenamento[nome]++
}

func (i *ArmazenamentoDeJogadorNaMemoria) PegaPontuacaoDoJogador(nome string) int {
    return i.armazenamento[nome]
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
    servidor:= NovoServidorDoJogador(NovoArmazenamentoDeJogadorNaMemoria())

    if err := http.ListenAndServe(":5000", servidor); err != nil {
        log.Fatalf("Não foi possivel ouvir na porta 5000 %v", err)
    }
}
```

Você pode encontrar todos os testes relacionados no link no começo desse capítulo.

## Armazene os dados

Existem diversos bancos de dados que poderíamos usar para isso, mas nós vamos por uma abordagem mais simples. Nós iremos armazenar os dados para essa aplicação em um arquivo como JSON.

Isso mantém os dados bastante manipuláveis e é relativamente simples de implementar.

Não será bem escalável mas, dado que isto é um protótipo, vai funcionar para agora. Se nossas circunstâncias mudarem e isto não for mais apropriado, será simples trocar para algo diferente por conta da abstração de `GuardarJogadores` que nós usamos.

Nós vamos manter o `NovoArmazenamentoDeJogadorNaMemoria` por enquanto para que os testes de integração continuem passando a medida que formos desenvolvendo nossa armazenamento. Quando estivermos confiantes que nossa implementação é suficiente para fazer os testes de integração passarem , nós iremos trocar e apagar `NovoArmazenamentoDeJogadorNaMemoria`

## Escreva os testes primeiro

Por agora você deve estar familiar com as interfaces em torno da biblioteca padrão para leitura de dados \(`io.Reader`\), escrita de dados \(`io.Writer`\) e como nós podemos usar a biblioteca padrão para testar essas funções sem ter que usar arquivos de verdade.

Para esse trabalho ser completo precisamos implementar `GuardaJogador` , então escreveremos testes para nossa armazenamento chamando os métodos que nós precisamos implementar. Começaremos com `PegaLiga`.

```go
func TestSistemaDeArquivoDeArmazenamentoDoJogador(t *testing.T) {

    t.Run("/liga de um leitor", func(t *testing.T) {
        bancoDeDados := strings.NewReader(`[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)

        armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

        recebido := armazenamento.PegaLiga()

        esperado := []Jogador{
            {"Cleo", 10},
            {"Chris", 33},
        }

        defineLiga(t, recebido, esperado)
    })
}
```

Estamos usando `strings.NewReader` que irá nos retornar um `Reader`, que é o que nosso `SistemaDeArquivoDeArmazenamentoDoJogador` irá usar para ler os dados. Em `main` abriremos um arquivo, que também é um `Reader`.

## Tente rodar o teste

```text
# github.com/larien/aprenda-go-com-testes/json-and-io/v7
./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:15:12: undefined: SistemaDeArquivoDeArmazenamentoDoJogador
```

## Escreva código suficiente para fazer o teste rodar e veja o retorno do erro do teste

Vamos definir `SistemaDeArquivoDeArmazenamentoDoJogador` em um novo arquivo

```go
type SistemaDeArquivoDeArmazenamentoDoJogador struct {}
```

Tente de novo

```text
# github.com/larien/aprenda-go-com-testes/json-and-io/v7
./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:15:28: too many values in struct initializer
./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:17:15: armazenamento.PegaLiga undefined (type SistemaDeArquivoDeArmazenamentoDoJogador has no field or method PegaLiga)
```

Está reclamando porque estamos passando para ele um `Reader` mas não está esperando um e não tem `PegaLiga` definida ainda.

```go
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
    bancoDeDados io.Reader
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() []Jogador {
    return nil
}
```

Tente mais uma vez...

```text
=== RUN   TestSistemaDeArquivoDeArmazenamentoDoJogador//league_from_a_reader
    --- FAIL: TestSistemaDeArquivoDeArmazenamentoDoJogador//league_from_a_reader (0.00s)
        SistemaDeArquivoDeArmazenamentoDoJogador_test.go:24: recebido [] esperado [{Cleo 10} {Chris 33}]
```

## Escreva código suficiente para fazer passar

Nós lemos JSON de um leitor antes

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() []Jogador {
    var liga []Jogador
    json.NewDecoder(f.bancoDeDados).Decode(&liga)
    return liga
}
```

O teste deve passar.

## Refatore

_Fizemos_ isso antes! Nosso código de teste para o servidor tinha que decodificar o JSON da resposta.

Vamos tentar DRYando isso em uma função.

Crie um novo arquivo chamado `liga.go` e coloque isso nele.

```go
func NovaLiga(rdr io.Reader) ([]Jogador, error) {
    var liga []Jogador
    err := json.NewDecoder(rdr).Decode(&liga)
    if err != nil {
        err = fmt.Errorf("Problema parseando a liga, %v", err)
    }

    return liga, err
}
```

Chame isso em nossa implementação e em nosso teste helper `obterLigaDaResposta` in `serv_test.go`

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() []Jogador {
    liga, _ := NovaLiga(f.bancoDeDados)
    return liga
}
```

Ainda não temos a estratégia para lidar com a análise de erros mas vamos continuar.

### Procurando problemas

Existe um problema na nossa implementação. Primeiramente, vamos relembrar como `io.Reader` é definida.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

Com nosso arquivo, você consegue imagina-lo lendo byte por byte até o fim. O que acontece se você tentar e `ler` uma segunda vez?

Adicione o seguinte no final do seu teste atual.

```go
// read again
recebido = armazenamento.PegaLiga()
defineLiga(t, recebido, esperado)
```

Queremos que passe, mas se você rodar o teste ele não passa.

O problema é nosso `Reader` chegou no final, então não tem mais nada para ser lido. Precisamos de um jeito de avisar para voltar ao inicio.

[ReadSeeker](https://golang.org/pkg/io/#ReadSeeker) é outra interface na biblioteca padrão que pode ajudar.

```go
type ReadSeeker interface {
    Reader
    Seeker
}
```

Lembra-se do incorporamento? Esta é uma interface composta de `Reader` e [`Seeker`](https://golang.org/pkg/io/#Seeker)

```go
type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}
```

Parece bom, podemos mudar `SistemaDeArquivoDeArmazenamentoDoJogador` para pegar essa interface no lugar?

```go
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
    bancoDeDados io.ReadSeeker
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() []Jogador {
    f.bancoDeDados.Seek(0, 0)
    liga, _ := NovaLiga(f.bancoDeDados)
    return liga
}
```

Tente rodar o teste,agora passa! Ainda bem que `string.NewReader` que nós usamos em nosso teste também implementa `ReadSeeker` então não precisamos mudar nada.

A seguir vamos implementar `PegarPontuacaooDoJogador`.

## Escreva o teste primeiro

```go
t.Run("pegar pontuação do jogador", func(t *testing.T) {
    bancoDeDados := strings.NewReader(`[
        {"Nome": "Cleo", "Vitorias": 10},
        {"Nome": "Chris", "Vitorias": 33}]`)

    armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

    recebido := armazenamento.("Chris")

    esperado := 33

    if recebido != esperado {
        t.Errorf("recebido %d esperado %d", recebido, esperado)
    }
})
```

## Tente rodar o teste

`./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:38:15: armazenamento. undefined (type SistemaDeArquivoDeArmazenamentoDoJogador has no field or method )`

## Escreva código suficiente para fazer o teste rodar e veja o retorno do erro do teste

Precisamos adicionar o método para o novo tipo para fazer o teste compilar.

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) (nome string) int {
    return 0
}
```

Agora compila e o teste falha

```text
=== RUN   TestSistemaDeArquivoDeArmazenamentoDoJogador/get_player_score
    --- FAIL: TestSistemaDeArquivoDeArmazenamentoDoJogador//get_player_score (0.00s)
        SistemaDeArquivoDeArmazenamentoDoJogador_test.go:43: recebido 0 esperado 33
```

## Escreva código sufience para fazer passar

Podemos iterar sobre a liga para encontrar o jogador e retornar a pontuação dele.

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) (nome string) int {

    var vitorias int

    for _, jogador := range f.PegaLiga() {
        if jogador.Nome == nome {
            vitorias = jogador.Vitorias
            break
        }
    }

    return vitorias
}
```

## Refatore

Você terá visto vários refatoramentos de teste helper, então deixarei este para você fazer funcionar

```go
t.Run("/pega pontuacao do  jogador", func(t *testing.T) {
    bancoDeDados := strings.NewReader(`[
        {"Nome": "Cleo", "Vitorias": 10},
        {"Nome": "Chris", "Vitorias": 33}]`)

    armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

    recebido := armazenamento.("Chris")
    esperado := 33
    definePontuacaoIgual(t, recebido, esperado)
})
```

Finalmente, precisamos começar a salvar pontuações com `SalvaVitoria`.

## Escreva o teste primeiro

Nossa abordagem é um pouco ruim para escritas. Não podemos \(facilmente\) apenas atualizar uma "linha" de JSON em um arquivo. Precisaremos armazenar a _inteira_ nova representação de nosso banco de dados em cada escrita.

Como escrevemos? Normalmente usaríamos um `Writer`, mas já temos nosso `ReadSeeker`. Potencialmente podemos ter duas dependências, mas a biblioteca padrão já tem uma interface para nós: o `ReadWriteSeeker`, que permite fazermos tudo que precisamos com um arquivo.

Vamos atualizar nosso tipo:

```go
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
    bancoDeDados io.ReadWriteSeeker
}
```

Veja se compila:

```go
./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:15:34: cannot use bancoDeDados (type *strings.Reader) as type io.ReadWriteSeeker in field value:
    *strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:36:34: cannot use bancoDeDados (type *strings.Reader) as type io.ReadWriteSeeker in field value:
    *strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
```

Não é tão surpreendente que `strings.Reader` não implementa `ReadWriteSeeker`, então o que vamos fazer?

Temos duas opções:

-   Criar um arquivo temporário para cada teste. `*os.File` implementa `ReadWriteSeeker`. O pró disso é que isso se torna mais um teste de integração, mas nós realmente estamos lendo e escrevendo de um sistema de arquivos então isso nos dará um alto nível de confiança. Os contras são que preferimos testes unitários porque são mais rápidos e normalmente mais simples. Também precisaremos trabalhar mais criando arquivos temporários e então ter certeza que serão removidos após o teste.
-   Poderíamos usar uma biblioteca externa. [Mattetti](https://github.com/mattetti) escreveu uma biblioteca [filebuffer](https://github.com/mattetti/filebuffer) que implementa a interface que precisamos e assim não precisariamos modificar o sistema de arquivos.

Não acredito que exista uma resposta especialmente errada aqui, mas ao escolher usar uma biblioteca externa eu teria que explicar o gerenciamento de dependências! Então usaremos os arquivos.

Antes de adicionarmos nosso teste precisamos fazer nossos outros testes compilarem substituindo o `strings.Reader` com um `os.File`.

Vamos criar uma função auxiliar que irá criar um arquivo temporário com alguns dados dentro dele

```go
func criaArquivoTemporario(t *testing.T, dadoInicial string) (io.ReadWriteSeeker, func()) {
    t.Helper()

   arquivotmp, err := ioutil.TempFile("", "db")

    if err != nil {
        t.Fatalf("não foi possivel escrever o arquivo temporário %v", err)
    }

    arquivotmp.Write([]byte(dadoInicial))

    removeArquivo := func() {
        arquivotmp.Close()
        os.Remove(arquivotmp.Name())
    }

    return arquivotmp, removeArquivo
}
```

[TempFile](https://golang.org/pkg/io/ioutil/#TempDir) cria um arquivo temporário para usarmos. O valor `"db"` que passamos é um prefixo colocado em um arquivo de nome aleatório que vai criar. Isto é para garantir que não vai dar conflito acidental com outros arquivos.

Você irá notar que não estamos retornando apenas nosso `ReadWriteSeeker` \(o arquivo\) mas também uma função. Precisamos garantir que o arquivo é removido uma vez que o teste é finalizado. Não queremos que dados sejam vazados dos arquivos no teste como é possível acontecer e desinteressante para o leitor. Ao retornar uma função `removeArquivo` , cuidamos dos detalhes no nosso auxiliar e tudo que a chamada precisa fazer é executar `defer limpaBancoDeDados()`.

```go
func TestaArmazenamentoDeSistemaDeArquivo(t *testing.T) {

    t.Run("liga de um leitor", func(t *testing.T) {
        bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)
        defer limpaBancoDeDados()

        armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

        recebido := armazenamento.PegaLiga()

        esperado := []Jogador{
            {"Cleo", 10},
            {"Chris", 33},
        }

        defineLiga(t, recebido, esperado)

        // ler novamente
        recebido = armazenamento.PegaLiga()
        defineLiga(t, recebido, esperado)
    })

    t.Run("retorna pontuação do jogador", func(t *testing.T) {
        bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)
        defer limpaBancoDeDados()

        armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

        recebido := armazenamento.("Chris")
        esperado := 33
        definePontuacaoIgual(t, recebido, esperado)
    })
}
```

Rode os testes e eles devem estar passando! Teve uma quantidade razoável de mudanças mas agora parece que nossa definição de interface completa e deve ser muito fáci adicionar novos testes de agora em diante.

Vamos pegar a primeira iteração de gravar uma vitória de um jogador existente

```go
t.Run("armazena vitórias de um jogador existente", func(t *testing.T) {
    bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
        {"Nome": "Cleo", "Vitorias": 10},
        {"Nome": "Chris", "Vitorias": 33}]`)
    defer limpaBancoDeDados()

    armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

    armazenamento.SalvaVitoria("Chris")

    recebido := armazenamento.("Chris")
    esperado := 34
    definePontuacaoIgual(t, recebido, esperado)
})
```

## Tente rodar o teste

`./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:67:8: armazenamento.SalvaVitoria undefined (type SistemaDeArquivoDeArmazenamentoDoJogador has no field or method SalvaVitoria)`

## Escreva código suficiente para fazer o teste rodar e veja o retorno do erro do teste

Adicione um novo método

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {

}
```

```text
=== RUN   TestSistemaDeArquivoDeArmazenamentoDoJogador/store_wins_for_existing_players
    --- FAIL: TestSistemaDeArquivoDeArmazenamentoDoJogador/store_wins_for_existing_players (0.00s)
        SistemaDeArquivoDeArmazenamentoDoJogador_test.go:71: recebido 33 esperado 34
```

Nossa implementação está vazia então a pontuação anterior está sendo retornada.

## Escreva código sufience para fazer passar

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {
    liga := f.PegaLiga()

    for i, jogador := range liga {
        if jogador.Nome == nome {
            liga[i].Vitorias++
        }
    }

    f.bancoDeDados.Seek(0,0)
    json.NewEncoder(f.bancoDeDados).Encode(liga)
}
```

Você deve está se perguntando por que estou fazendo `liga[i].Vitorias++` invés de `jogador.Vitorias++`.

Quando você `percorre` sobre um pedaço é retornado o índice atual do laço \(no nosso caso `i`\) e uma _cópia_ do elemento naquele índice. Mudando o valor `Vitorias` não irá afetar no pedaço `liga` que iteramos sobre. Por este motivo, precisamos pegar a referência do valor atual fazendo `liga[i]` e então mudando este valor.

Se rodar os testes, eles devem estar passando.

## Refatore

Em `PegaPontuacaoDoJogador` e `SalvaVitoria`, estamos iterando sobre `[]Jogador` para encontrar um jogador pelo nome.

Poderíamos refatorar esse código comum nos internos de `SistemaDeArquivoDeArmazenamentoDoJogador` mas para mim, parece que talvez seja um código util então poderíamos colocar em um novo tipo. Trabalhando com uma "Liga" até agora tem sido com `[]Jogador` mas podemos criar um novo tipo chamado `Liga`. Será mais fácil para outros desenvolvedores entenderem e assim podemos anexar métodos utéis dentro desse tipo para usarmos.

Dentro de `liga.go` adicionamos o seguinte

```go
type Liga []Jogador
func (l Liga) Find(nome string) *Jogador {
    for i, p := range l {
        if p.Nome==nome {
            return &l[i]
        }
    }
    return nil
}
```

Agora se qualquer um tiver uma `Liga` facilmente será encontrado um dado jogador.

Mude nossa interface `GuardaJogador` para retornar `Liga` invés de `[]Jogador`. Tente e rode novamente os teste, você terá um problema de compilação por termos modificado a interface mas é fácil de resolver; apenas modifique o tipo de retorno de `[]Jogador` to `Liga`.

Isso nos permite simplificar os métodos em `SistemaDeArquivoDeArmazenamentoDoJogador`.

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) (nome string) int {

    jogador := f.PegaLiga().Find(nome)

    if  jogador != nil {
        return  jogador.Vitorias
    }

    return 0
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {
    liga := f.PegaLiga()
    jogador :=liga.Find(nome)

    if  jogador != nil {
        jogador.Vitorias++
    }

    f.bancoDeDados.Seek(0, 0)
    json.NewEncoder(f.bancoDeDados).Encode(liga)
}
```

Isto parece bem melhor and podemos ver como talvez possamos encontrar como outras funcionalidades úteis em torno de `Liga` podem ser refatoradas.

Agora precisamos tratar o cenário de salvar vitórias de novos jogadores.

## Escreva o teste primeiro

```go
t.Run("armazena vitorias de novos jogadores", func(t *testing.T) {
    bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
        {"Nome": "Cleo", "Vitorias": 10},
        {"Nome": "Chris", "Vitorias": 33}]`)
    defer limpaBancoDeDados()

    armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

    armazenamento.SalvaVitoria("Pepper")

    recebido := armazenamento.("Pepper")
    esperado := 1
    definePontuacaoIgual(t, recebido, esperado)
})
```

## Tente rodar o teste

```text
=== RUN   TestSistemaDeArquivoDeArmazenamentoDoJogador/store_wins_for_new_players#01
    --- FAIL: TestSistemaDeArquivoDeArmazenamentoDoJogador/store_wins_for_new_players#01 (0.00s)
        SistemaDeArquivoDeArmazenamentoDoJogador_test.go:86: recebido 0 esperado 1
```

## Escreva código suficiente para fazer passar

Apenas precisamos tratar o caso onde `Find` returna `nil` por não ter conseguido encontrar o jogador.

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {
    liga := f.PegaLiga()
    jogador := liga.Find(nome)

    if jogador != nil {
        jogador.Wins++
    } else {
        liga = append(liga, Jogador{nome, 1})
    }

    f.bancoDeDados.Seek(0, 0)
    json.NewEncoder(f.bancoDeDados).Encode(liga)
}
```

O caminho feliz parece bom então agora vamos tentar usar nossa nova `armazenamento` no teste de integração. Isto nos dará mais confiança que o software funciona e então podemos deletar o redundante `NovoArmazenamentoDeJogadorNaMemoria`.

Em `TestRecordingWinsAndRetrievingThem` substitui a velha armazenamento.

```go
bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, "")
defer limpaBancoDeDados()
armazenamento := &SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}
```

Se você rodar o teste ele deve passar e agora podemos deletar `NovoArmazenamentoDeJogadorNaMemoria`. `main.go` terá problemas de compilação que nos motivará para agora usar nossa nova armazenamento no código "real".

```go
package main

import (
    "log"
    "net/http"
    "os"
)

const dbFileName = "game.db.json"

func main() {
    db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        log.Fatalf("problema abrindo %s %v", dbFileName, err)
    }

    armazenamento := &SistemaDeArquivoDeArmazenamentoDoJogador{db}
    server := NovoServidorDoJogador(armazenamento)

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("não foi possivel escutar na porta 5000 %v", err)
    }
}
```

-   Nós criamos um arquivo para nosso banco de dados.
-   O 2º argumento para `os.OpenFile` permite definir as permissões para abrir um arquivo, no nosso caso `O_RDWR` significa que queremos ler e escrever _e_ `os.O_CREATE` significa criar um arquivo se ele não existe.
-   O 3º argumento significa definir as permissões para o arquivo, no nosso caso, todos os usuários podem ler e escrever o arquivo. [\(Veja superuser.com para uma explicação mais detalhada\)](https://superuser.com/questions/295591/what-is-the-meaning-of-chmod-666).

Rodando o programa agora os dados permanecem em um arquivo entre reinicializações, uhu!

## Mais refatoramento e preocupações com performance

Toda vez que alguém chama `PegaLiga()` ou `()` estamos lendo o arquivo do ínicio, e transformando ele em JSON. Não deveríamos ter que fazer isso porque `SistemaDeArquivoDeArmazenamentoDoJogador` é inteiramente responsável pelo estado da liga; apenas queremos usar o arquivo para pegar o estado atual e atualiza-lo quando os dados mudarem.

Podemos criar um construtor que pode fazer parte dessa inicialização para nós e armazena a liga como um valor em nosso `SistemaDeArquivoDeArmazenamentoDoJogador` para ser usado nas leitura então.

```go
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
    bancoDeDados io.ReadWriteSeeker
    liga Liga
}

func NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados io.ReadWriteSeeker) *SistemaDeArquivoDeArmazenamentoDoJogador {
    bancoDeDados.Seek(0, 0)
    liga, _ := NovaLiga(bancoDeDados)
    return &SistemaDeArquivoDeArmazenamentoDoJogador{
        bancoDeDados:bancoDeDados,
        liga:liga,
    }
}
```

Desta maneira precisamos ler do disco apenas uma vez . Podemos agora substituir todas as nossas chamadas anteriores para pegar a liga do disco e apenas usar `f.liga` no lugar.

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() Liga {
    return f.liga
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) (nome string) int {

    jogador := f.liga.Find(nome)

    if jogador != nil {
        return jogador.Vitorias
    }

    return 0
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {
    jogador := f.liga.Find(nome)

    if jogador != nil {
        jogador.Vitorias++
    } else {
        f.liga = append(f.liga, Jogador{nome, 1})
    }

    f.bancoDeDados.Seek(0, 0)
    json.NewEncoder(f.bancoDeDados).Encode(f.liga)
}
```

Se você tentar e rodar os testes eles agora vão reclamar sobre inicializar `SistemaDeArquivoDeArmazenamentoDoJogador` então fixe-o chamando nosso construtor.

### Outro problema

Existe mais alguma ingenuidade na maneira como estamos lidando com arquivos que _poderiamos_ criar um erro bem bobo futuramente.

Quando nós chamamos `SalvaVitoria` nós `procuramos` no ínicio do arquivo e então escrevemos o novo dado mas e se o novo dado for menor que o que estava lá antes?

Na nossa situação atual, isso é impossível. Nunca editamos ou apagamos pontuações, então os dados apenas podem aumentar, mas seria irresponsabilidade nossa deixar o código desse jeito, não é inimaginável que um cenário de apagamento poderia aparecer.

Como iremos testar isso então? O que precisamos fazer primeiro é refatorar nosso código, então separamos nossa preocupação do _tipo de dados que escrevemos, da escrita_. Podemos então testar isso separadamente para verificar se funciona como esperamos.

Agora iremos criar um novo tipo para encapsular nossa funcionalidade "quando escrevemos, vamos para o começo". Vou chama-la de `Fita`. Criamos um novo arquivo com o seguinte

```go
package main

import "io"

type fita struct {
    arquivo io.ReadWriteSeeker
}

func (t *fita) Write(p []byte) (n int, err error) {
    t.arquivo.Seek(0, 0)
    return t.arquivo.Write(p)
}
```

Note que apenas implementamos `Write` agora, já que encapsula a parte de `Procura` . Isso que dizer que `SistemaDeArquivoDeArmazenamentoDoJogador` pode ter uma referência a `Writer` invés disso.

```go
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
    bancoDeDados io.Writer
    liga   Liga
}
```

Atualize o construtor para usar `fita`

```go
func NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados io.ReadWriteSeeker) *SistemaDeArquivoDeArmazenamentoDoJogador {
    bancoDeDados.Seek(0, 0)
    liga, _ := NovaLiga(bancoDeDados)

    return &SistemaDeArquivoDeArmazenamentoDoJogador{
        bancoDeDados: &fita{bancoDeDados},
        liga:   liga,
    }
}
```

Finalmente, podemos ter o incrível beneficio que queríamos removendo `Procura` de `SalvaVitoria`. Sim, não parece muito, mas pelo menos isso significa que, se fizermos qualquer outro tipo de escritas, podemos confiar no nosso `Write` para se comportar como precisamos. Além disso, agora podemos testar o potencial código problemático separadamente e corrigi-lo.

Agora vamos escrever o teste onde atualizamos todo o conteúdo de um arquivo com algo menor que o conteúdo original . Em `fita_test.go`:

## Escreva o teste primeiro

Vamos apenas criar um arquivo, tentar e escrever nele usando nossa fita, ler todo novamente e visualizar o que está no arquivo

```go
func TestaFita_Escrita(t *testing.T) {
    arquivo, limpa := criaArquivoTemporario(t, "12345")
    defer limpa()

    fita := &fita{arquivo}

    fita.Write([]byte("abc"))

    arquivo.Seek(0, 0)
    novoConteudoDoArquivo, _ := ioutil.ReadAll(arquivo)

    recebido := string(novoConteudoDoArquivo)
    esperado := "abc"

    if recebido != esperado {
        t.Errorf("recebido '%s' esperado '%s'", recebido, esperado)
    }
}
```

## Tente rodar o teste

```text
=== RUN   TestaFita_Escrita
--- FAIL: TestaFita_Escrita (0.00s)
    fita_test.go:23: recebido 'abc45' esperado 'abc'
```

Como pensamos! Ele apenas escreve os dados que queremos, deixando todo o resto.

## Escreva código suficiente para fazer passar

`os.File` tem uma função truncada que vai permitir que o arquivo seja esvaziado eficientemente. Devemos ser capazes de apenas chama-la para conseguir o que queremos.

Mude `fita` para o seguinte

```go
type fita struct {
    file *os.File
}

func (t *fita) Write(p []byte) (n int, err error) {
    t.file.Truncate(0)
    t.file.Seek(0, 0)
    return t.file.Write(p)
}
```

O compilador irá falhar em alguns lugares quando esperamos um `io.ReadWriteSeeker` mas estamos mandando um `*os.File`. Você deve ser capaz de corrigir esses problemas por conta própria, mas se ficar preso basta checar o código fonte.

Uma vez que você tenha refatorado nosso teste `TestaFita_Escrita` deve estar passando!

### Uma outra pequena refatoração

Em `SalvaVitoria` temos uma linha`json.NewEncoder(f.bancoDeDados).Encode(f.league)`.

Não precisamos criar um novo codificador toda vez que escrevemos, podemos inicializar um em nosso construtor e usa-lo.

Armazena uma referência para um `Encoder` para nosso tipo.

```go
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
    bancoDeDados *json.Encoder
    liga   Liga
}
```

Inicialize no construtor

```go
func NovoSistemaDeArquivoDeArmazenamentoDoJogador(arquivo *os.File) *SistemaDeArquivoDeArmazenamentoDoJogador {
    arquivo.Seek(0, 0)
    liga, _ := NovaLiga(arquivo)

    return &SistemaDeArquivoDeArmazenamentoDoJogador{
        bancoDeDados: json.NewEncoder(&fita{arquivo}),
        liga:   liga,
    }
}
```

Use em `SalvaVitoria`.

## Não quebramos algumas regras ali? Testando coisas privadas? Sem interfaces?

### Testando tipos privados

É verdade que _no geral_ deve ser favorecido não testar coisas privadas, uma vez que isso, as vezes, leva a testar coisas bastante acopladas para a implementação; que pode impedir refatoramento no futuro.

Entretanto,não devemos esquecer que testes nos dá _confiança_.

Não estamos confiantes que nossa implementação funcionaria se tivéssemos adicionado algum tipo de funcionalidade para editar ou deletar. Não queremos deixar o código assim, especialmente se isso foi trabalhado por mais de uma pessoa que talvez não estivesse ciente dos defeitos da nossa abordagem.

Finalmente, é apenas um teste! Se decidirmos mudar a maneira como funciona não será um desastre deletar o teste, mas teremos que ter pego o requisito para futuro mantenedores.

### Interfaces

Começamos o código usando `io.Reader` como o caminho mais fácil para testar de forma unitária nosso novo `GuardaJogador`. A medida que desenvolvemos nosso código, movemos para `io.ReadWriter` e então para `io.ReadWriteSeeker`. Descobrimos então que não tinha nada na biblioteca padrão que implementasse isso além de `*os.File`. Poderiamos ter decidido escrever o nosso ou usar um de código aberto, mas isso pareceu pragmático apenas para fazer arquivos temporários para os testes.

Finalmente, precisamos de `Truncate` que também está no `*os.File`. Isso seria uma opção para criar nossa própria interface pegando esses requisitos.

```go
type ReadWriteSeekTruncate interface {
    io.ReadWriteSeeker
    Truncate(size int64) error
}
```

Mas o que isso está realmente nos dando? Lembre-se que _não estamos mockando_ e isso é irrealista para um armazenamento de **sistema de arquivos** receber outro tipo além que um `*os.File` então não precisamos do polimorfismo que interface nos dá.

Não tenha medo de cortar e mudar tipos e experimentar como temos aqui. O bom de usar uma linguagem tipada estaticamente é o compilador que ajudará você com toda mudança.

## Tratamento de erros

Antes de começarmos no ordenamento, devemos ter certeza que estamos contentes com nosso código atual e remover qualquer débito técnico que ainda resta. É um principio importante para trabalhar com software o mais rápido possível \(mantenha-se fora do estado vermelho\) mas isso não quer dizer que devemos ignorar os casos de erro!

Se voltarmos para `SistemaDeArquivoDeArmazenamentoDoJogador.go` temos `liga, _ := NovaLiga(f.bancoDeDados)` no nosso construtor.

`NovaLiga`pode retornar um erro se é instável passar a liga do `io.Reader` que fornecemos.

Era pragmático ignorar isso naquela hora como já tinhamos testes falhando. Se tivemos tentado lidar com isso ao mesmo tempo estamos lidando com duas coisas de uma vez.

Vamos fazer com que nosso construtor seja capaz de retornar um erro.

```go
func NovoSistemaDeArquivoDeArmazenamentoDoJogador(arquivo *os.File) (*SistemaDeArquivoDeArmazenamentoDoJogador, error) {
    arquivo.Seek(0, 0)
    liga, err := NovaLiga(arquivo)

    if err != nil {
        return nil, fmt.Errorf("problema carregando o armazenamento do jogador  de arquivo %s, %v", arquivo.Nome(), err)
    }

    return &SistemaDeArquivoDeArmazenamentoDoJogador{
        bancoDeDados: json.NewEncoder(&fita{arquivo}),
        liga:   liga,
    }, nil
}
```

Lembre-se que é importante retornar mensagens de erro úteis \(assim como nossos testes\). As pessoas na internet dizem que a maioria dos códigos em Go é

```go
if err != nil {
    return err
}
```

**Isso é 100% não idiomático.** Adicionando informação contextual \(i.e o que você estava fazendo que causou o erro\\) para suas mensagens de erro facilita manipular o software.

Se você tentar e compilar, vai ver alguns erros.

```text
./main.go:18:35: multiple-value NovoSistemaDeArquivoDeArmazenamentoDoJogador() in single-value context
./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:35:36: multiple-value NovoSistemaDeArquivoDeArmazenamentoDoJogador() in single-value context
./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:57:36: multiple-value NovoSistemaDeArquivoDeArmazenamentoDoJogador() in single-value context
./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:70:36: multiple-value NovoSistemaDeArquivoDeArmazenamentoDoJogador() in single-value context
./SistemaDeArquivoDeArmazenamentoDoJogador_test.go:85:36: multiple-value NovoSistemaDeArquivoDeArmazenamentoDoJogador() in single-value context
./server_integration_test.go:12:35: multiple-value NovoSistemaDeArquivoDeArmazenamentoDoJogador() in single-value context
```

Em main vamos querer sair do programa, imprimindo o erro.

```go
armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(db)

if err != nil {
    log.Fatalf("problema criando o sistema de arquivo do armazenamento do jogador, %v ", err)
}
```

Nos nossos testes podemos garantir que não exista erro . Podemos fazer uma função auxiliar para ajudar com isto.

```go
func defineSemErro(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("não esperava um erro mas obteve um, %v", err)
    }
}
```

Trabalhe nos outros problemas de compilação usando essa auxiliar. Finalmente, você deve ter um teste falhando

```text
=== RUN   TestRecordingWinsAndRetrievingThem
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:14: não esperava um erro mas obteve um, problema carregando o armazenamento do jogador  de arquivo /var/folders/nj/r_ccbj5d7flds0sf63yy4vb80000gn/T/db841037437, problem parsing league, EOF
```

Não podemos analisar a liga porque o arquivo está vazio.Não estávamos obtendo erros antes porque sempre os ignoramos.

Vamos corrigir nosso grande teste de integração colocando algum JSON válido nele e então podemos escrever um teste específico para este cenário.

```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
    bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[]`)
    //etc...
```

Agora todos os testes estão passando, precisamos então lidar com o cenário onde o arquivo está vazio.

## Escreva o teste primeiro

```go
t.Run("funciona com um arquivo vazio", func(t *testing.T) {
    bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, "")
    defer limpaBancoDeDados()

    _, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados)

    defineSemErro(t, err)
})
```

## Tente rodar o teste

```text
=== RUN   TestSistemaDeArquivoDeArmazenamentoDoJogador/works_with_an_empty_file
    --- FAIL: TestSistemaDeArquivoDeArmazenamentoDoJogador/works_with_an_empty_file (0.00s)
        SistemaDeArquivoDeArmazenamentoDoJogador_test.go:108: não esperava um erro mas obteve um, problema carregando o armazenamento do jogador  de arquivo /var/folders/nj/r_ccbj5d7flds0sf63yy4vb80000gn/T/db019548018, problem parsing league, EOF
```

## Escreva código sufience para fazer passar

Mude nosso construtor para o seguinte

```go
func NovoSistemaDeArquivoDeArmazenamentoDoJogador(arquivo *os.File) (*SistemaDeArquivoDeArmazenamentoDoJogador, error) {

    arquivo.Seek(0, 0)

    info, err := arquivo.Stat()

    if err != nil {
        return nil, fmt.Errorf("problema ao usar o arquivo  %s, %v", arquivo.Nome(), err)
    }

    if info.Size() == 0 {
        file.Write([]byte("[]"))
        file.Seek(0, 0)
    }

    liga, err := NovaLiga(file)

    if err != nil {
        return nil, fmt.Errorf("problema carregando armazenamento de jogador do aquivo %s, %v", arquivo.Nome(), err)
    }

    return &SistemaDeArquivoDeArmazenamentoDoJogador{
        bancoDeDados: json.NewEncoder(&fita{file}),
        liga:   liga,
    }, nil
}
```

`Arquivo.Stat` retorna estatísticas do nosso arquivo. Isto nos permite checar o tamanho do arquivo, se está vazio podemos `Escrever` um array JSON vazio e `Busca` de volta para o ínicio, pronto para o resto do arquivo.

## Refatore

Nosso construtor está um pouco bagunçado, podemos extrair o código de inicialização em uma função

```go
func iniciaArquivoBDDeJogador(arquivo *os.File) error {
    arquivo.Seek(0, 0)

    info, err := arquivo.Stat()

    if err != nil {
        return fmt.Errorf("problema ao usar arquivo %s, %v", file.Name(), err)
    }

    if info.Size()==0 {
        arquivo.Write([]byte("[]"))
        arquivo.Seek(0, 0)
    }

    return nil
}
```

```go
func NovoSistemaDeArquivoDeArmazenamentoDoJogador(arquivo *os.File) (*SistemaDeArquivoDeArmazenamentoDoJogador, error) {

    err := iniciaArquivoBDDeJogador(file)

    if err != nil {
        return nil, fmt.Errorf("problema inicializando arquivo do jogador, %v", err)
    }

    liga, err := Nova(liga)

    if err != nil {
        return nil, fmt.Errorf("problema carregando armazenamento de jogador do arquivo %s, %v", arquivo.Nome(), err)
    }

    return &SistemaDeArquivoDeArmazenamentoDoJogador{
        bancoDeDados: json.NewEncoder(&fita{file}),
        liga:   liga,
    }, nil
}
```

## Ordenação

Nossa dona do produto quer que `/liga` retorne os jogadores ordenados pela pontuação.

A principal decisão a ser feita é onde isso deve acontecer no software. Se estamos usando um "verdadeiro" banco de dados usariamos coisas como `ORDER BY` , então o ordenamento é super rápido por esse motivo parece que a implementção de `GuardaJogador` deve ser responsável.

## Escreva o teste primeiro

Podemos atualizar a inserção no nosso primeiro teste em `TestaArmazenamentoDeSistemaDeArquivo`

```go
t.Run("liga ordernada", func(t *testing.T) {
    bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
        {"Nome": "Cleo", "Vitorias": 10},
        {"Nome": "Chris", "Vitorias": 33}]`)
    defer limpaBancoDeDados()

    armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{bancoDeDados}

   recebido := armazenamento.PegaLiga()

   esperado:= []Jogador{
        {"Chris", 33},
        {"Cleo", 10},
    }

    defineLiga(t, recebido, esperado)

    // read again
    recebido = armazenamento.PegaLiga()
    defineLiga(t, recebido, esperado)
})
```

A ordem que está sendo recebida do JSON está errada e nosso `esperado` vai checar que é retornado para o chamador na ordem correta.

## Tente rodar o teste

```text
=== RUN   TestSistemaDeArquivoDeArmazenamentoDoJogador/league_from_a_reader,_sorted
    --- FAIL: TestSistemaDeArquivoDeArmazenamentoDoJogador/league_from_a_reader,_sorted (0.00s)
        SistemaDeArquivoDeArmazenamentoDoJogador_test.go:46: recebido [{Cleo 10} {Chris 33}] esperado [{Chris 33} {Cleo 10}]
        SistemaDeArquivoDeArmazenamentoDoJogador_test.go:51: recebido [{Cleo 10} {Chris 33}] esperado [{Chris 33} {Cleo 10}]
```

## Escreva código sufience para fazer passar

```go
func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() League {
    sort.Slice(f.liga, func(i, j int) bool {
        return f.liga[i].Vitorias > f.liga[j].Vitorias
    })
    return f.liga
}
```

[`sort.Slice`](https://golang.org/pkg/sort/#Slice)

> Slice ordena a parte fornecida dada a menor função fornecida

Moleza!

## Finalizando

### O que cobrimos

-   A interface `Seeker` e sua relação com `Reader` e `Writer`.
-   Trabalhando com arquivos.
-   Criando uma auxiliar fácil de usar para testes com arquivos que escondem todas as bagunças.
-   `sort.Slice` para ordenar partes.
-   Usando o compilador para nos ajudar a fazer mudanças estruturais de forma segura na aplicação.

### Quebrando regras

-   Maior partes das regras em engenharia de software não são realmente regras, apenas boas práticas que funcionam 80% do tempo.
-   Descobrimos um cenário onde nos "regras" anteriores de não testar funções internas não foi útil, então quebramos essa regra.
-   É importante entender o que estamos perdendo e ganhado ao quebrar as regras . No nosso caso, não tinha problema porque era apenas um teste e seria muito difícil exercitar o cenário contrário.
-   Para poder quebrar as regras, **você deve entende-las**. Uma analogia é com aprender a tocar violão. Não importa quão criativo você seja, você deve entender e praticar os fundamentos.

### Onde nosso software está

-   Temos uma API HTTP onde é possível criar jogadores e aumentar a pontuação deles..
-   Podemos retornar uma liga das pontuações de todos como JSON.
-   O dado é mantindo com um arquivo JSON.
