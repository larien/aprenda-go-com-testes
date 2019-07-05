# OS Exec

[**Você pode encontrar todo o código aqui**](https://github.com/larien/learn-go-with-tests/tree/master/q-and-a/os-exec)

[keith6014](https://www.reddit.com/user/keith6014) perguntou no [reddit](https://www.reddit.com/r/golang/comments/aaz8ji/testdata_and_function_setup_help/)

> Estou executando um comando usando os/exec.Command\(\) que gera dados XML. O comando será executado em uma função chamada GetData\(\). 
>
> Para testar GetData\(\), tenho alguns dados que criei para testar. 
>
> No meu \_test.go tenho o TestGetData que executa GetData\(\) mas isso vai usar os.exec, e ao invés disso eu gostaria de usar os meus dados para teste.
>
> Qual seria uma boa forma de conseguir isso? Quando chamo GetData, eu deveria ter uma flag "test" para que assim ela leia um arquivo - como GetData\(modo string\)?

Algumas considerações:

* Quando alguma coisa é difícil de se testar, geralmente é porque a separação de conceitos não foi feita muito bem.
* Não coloque "modos de teste" dentro do seu código. Ao invés disso, use [Injeção de Dependência](../primeiros-passos-com-go/dependency-injection.md) para que então você possa modelar suas dependências e separar os conceitos.

Eu tomei a liberdade de supor como o código deveria ser:

```go
type Payload struct {
    Message string `xml:"message"`
}

func GetData() string {
    cmd := exec.Command("cat", "msg.xml")

    out, _ := cmd.StdoutPipe()
    var payload Payload
    decoder := xml.NewDecoder(out)

    // esses 3 podem retornar erros, mas estou ignorando para ser mais direto
    cmd.Start()
    decoder.Decode(&payload)
    cmd.Wait()

    return strings.ToUpper(payload.Message)
}
```

* Uso `exec.Command` que te permite executar um comando externo ao processo.
* We capture the output in `cmd.StdoutPipe` which returns us a `io.ReadCloser` \(this will become important\)
* Capturamos a saída em `cmd.StdoutPipe` que retorna um `io.ReadCloser` \(isso será importante\).
* O resto do código é mais ou menos a cópia da [excelente documentação](https://golang.org/pkg/os/exec/#example_Cmd_StdoutPipe).
  * Capturamos qualquer saída de stdout em um `io.ReadCloser`e então rodamos o comando `Start`, e esperamos até todos os dados serem lidos executando `Wait`. Entre essas duas chamadas usamos o `Decode` na nossa struct `Payload`.

Esse é o conteúdo de `msg.xml`:

```markup
<payload>
    <message>Feliz Ano Novo!</message>
</payload>
```

Escrevi um teste simples para mostrar isso na prática:

```go
func TestGetData(t *testing.T) {
    got := GetData()
    want := "FELIZ ANO NOVO!"

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
}
```

## Código testável

Código testável é código desacoplado e com um propósito único. Na minha opinião, há duas preocupações com esse código:

1. Obtendo o dado cru do XML
2. Decodificar o XML e aplicá-lo na nossa regra de negócio \(nesse caso, `strings.ToUpper` no valor de `<message>`\).

A primeira parte é só uma cópia do exemplo da lib padrão.

A segunda parte é onde temos nossa regra de negócio e olhando para o código podemos ver onde a "costura" na nossa lógica começa; é onde pegamos nosso `io.ReadCloser`. Podemos usar essa abstração existente para dividir os conceitos e tornar nosso código testável.

**O problema com GetData é que a regra de negócio está acoplada com a parte de pegar o XML. Para fazer o design do nosso código melhor, precisamos separar essas partes.**

Nosso `TestGetData` pode agir como nosso teste de integração entre as duas responsabilidades, então vamos mantê-lo para garantir que o código continue funcionando.

Abaixo é como o recém dividido código fica:

```go
type Payload struct {
    Message string `xml:"message"`
}

func GetData(data io.Reader) string {
    var payload Payload
    xml.NewDecoder(data).Decode(&payload)
    return strings.ToUpper(payload.Message)
}

func getXMLFromCommand() io.Reader {
    cmd := exec.Command("cat", "msg.xml")
    out, _ := cmd.StdoutPipe()

    cmd.Start()
    data, _ := ioutil.ReadAll(out)
    cmd.Wait()

    return bytes.NewReader(data)
}

func TestGetDataIntegration(t *testing.T) {
    got := GetData(getXMLFromCommand())
    want := "HAPPY NEW YEAR!"

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
}
```

Agora que `GetData` tem na sua entrada somente um `io.Reader` nós o deixamos testável e não há mais a preocupação sobre como os dados são obtidos; e todo mundo pode reusar a função com qualquer coisa que retorne um `io.Reader` \(o que é bem comum\). Por exemplo, podemos começar a pegar o XML de uma URL ao invés da linha de comando.

```go
func TestGetData(t *testing.T) {
    input := strings.NewReader(`
<payload>
    <message>Gatos são os melhores animais</message>
</payload>`)

    got := GetData(input)
    want := "GATOS SÃO OS MELHORES ANIMAIS"

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
}
```

Esse é um exemplo de um teste unitário para `GetData`.

Separando os conceitos e usado as abstrações existentes dentro do Go, testar nossa preciosa regra de negócio é moleza.

