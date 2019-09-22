# Instalação do Go: defina seu ambiente para produtividade

As instruções oficiais de instalação do Go estão disponíveis [aqui](http://www.golangbr.org/doc/instalacao).

Esse guia vai presumir que você está usando um gerenciador de pacotes como [Homebrew](https://brew.sh), [Chocolatey](https://chocolatey.org), [Apt](https://help.ubuntu.com/community/AptGet/Howto) ou [yum](https://access.redhat.com/solutions/9934).

Para propósitos de demonstração, vamos te mostrar o procedimento de instalação para o OSX usando Homebrew.

## Instalação

### Mac OSX

O processo de instalação é bem simples. Primeiro, o que você precisa fazer é executar o comando abaixo pra instalar o homebrew (brew). O Brew depende do Xcode, então você deve se certificar de instalá-lo primeiro.

```bash
xcode-select --install
```

Depois, execute o comando a seguir para instalar o homebrew:

```bash
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

Agora você consegue instalar o Go:

```bash
brew install go
```

_Siga todas as instruções recomendadas pelo seu gerenciador de pacotes. **Nota** cada grupo de instruções varia de sistema operacional para sistema operacional._

Você pode verificar a instalação com:

```bash
$ go version
go version go1.10 darwin/amd64
```

### Linux

O processo de instalação é bem simples. Primeiro você precisa escolher e baixar a versão do Go que você deseja instalar. Para isso [acesse o site oficial](https://golang.org/) da linguagem e copie o [link](https://golang.org/dl/) da versão desejada (recomendamos instalar sempre a versão mais atual).

Para baixá-lo execute o seguinte comando no seu terminal.

```bash
#escolha a vsersão go Go que você deseja instalar, no nosso exemplo estamos utilizando a versão 1.10
VERSAO_GO=1.10
cd ~
curl -O "https://dl.google.com/go/go${VERSAO_GO}.linux-amd64.tar.gz"
```

Agora descompacte os arquivos com o seguinte comando.

```bash
tar xvf "go${VERSAO_GO}.linux-amd64.tar.gz"
```

E em seguida, mova os arquivos para o diretório de binário do seu usuário.

```bash
sudo mv go /usr/local
```

Agora teste a sua instalação.

```bash
go version
go version go1.13 linux/amd64
```

Nos próximos passos vamos configurar o ambiente Go. As instruções abaixo valem tanto para sistema operacional OSX quanto para o Linux.

## O Ambiente Go

O Go divide opiniões.

Por convenção, todo o código Go é colocado dentro de apenas um workspace (pasta). Esse workspace pode estar em qualquer lugar da sua máquina. Se você não especificar, o Go vai definir o \$HOME/go como workspace padrão. Ele é identificado (e modificado) pela variável de ambiente [GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable).

Você precisa definir a variável de ambiente para que possa utilizar futuramente em scripts, shells etc.

Atualize seu .bash_profile para conter os seguintes `exports`:

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

_Nota_ você deve abrir um novo terminal para definir essas variáveis de ambiente.

O Go presume que seu workspace contenha uma estrutura de diretórios específica.

Ele coloca seus arquivos em três diretórios: todo o código-fonte fica em `src`, os objetos dos pacotes ficam em `pkg` e os programas compilados são colocados em `bin`. É possível criar esses diretórios com o comando a seguir:

```bash
mkdir -p $GOPATH/src $GOPATH/pkg $GOPATH/bin
```

Agora você é capaz de usar o _go get_ para que o `src/package/bin` seja instalado corretamente no diretório \$GOPATH/xxx apropriado.

## Editor Go

A escolha de editor é bem pessoal. Você pode já ter um de sua preferência que tem suporte a Go. Se não tiver, leve em consideração um Editor como o [Visual Studio Code](https://code.visualstudio.com), que tem um suporte exceptional à linguagem.

Você pode instalá-lo com o comando a seguir:

```bash
brew cask install visual-studio-code
```

Confirme que o VS Code foi instalado corretamente executando o seguinte comando:

```bash
code .
```

O VS Code é lançado com poucos softwares habilidados. Você pode habilitar novos softwares instalando extensões. Para adicionar o suporte a Go, você deve instalar uma extensão. Existem várias disponíveis para o VS Code, mas uma excepcional é a do [Luke Hoban](https://github.com/Microsoft/vscode-go). Instale-a da forma a seguir:

```bash
code --install-extension ms-vscode.go
```

Quando abrir um arquivo Go pela primeira vez no VS Code, ele vai indicar que ferramentas de análises estão faltando. Clique no botão para instalá-las. A lista de ferramentas que são instaladas (e usadas) pelo VS Code estão disponíveis [aqui](https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on).

## Debugger do Go

Uma boa opção para debugar seus programas em Go (que é integrado com o VS Code) é o Delve. Ele pode ser instalado da seguinte maneira usando `go get`:

```bash
go get -u github.com/go-delve/delve/cmd/dlv
```

## Linter do Go

Uma melhoria sob o linter padrão pode ser configurada usando o [GolangCI-Lint](https://github.com/golangci/golangci-lint).

Que pode ser instalada da seguinte forma:

```bash
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
```

## Refatoração e suas ferramentas

Uma grande ênfase nesse livro é dada na importância da refatoração.

Suas ferramentas podem te ajudar a fazer uma refatoração com maior confiança.

Você deve ter familiaridade o suficiente com seu editor para performar as ações a seguir com uma simples combinação de teclas:

-   **Extrair/alinhar variável**. Ser capaz de pegar valores mágicos e dar um nome a eles vai simplificar seu código rapidamente.
-   **Extrair método/função**. É crucial ser capaz de tirar uma seção do código e extrair funções/métodos.
-   **Renomear**. Você deve se sentir capaz de renomear símbolos no decorrer dos arquivos com confiança.
-   **go fmt**. O Go tem um formatador nativo chamado `go fmt`. Seu editor deve executar esse comando a cada vez que salvar o arquivo.
-   **Executar testes**. Não precisa nem dizer que você deve ser capaz de fazer todos os pontos acima e então re-executar seus testes rapidamente para certificar que sua refatoração não quebrou nada.

Além disso, para te ajudar a trabalhar com seu código, você deve ser capaz de:

-   **Verificar a assinatura da função**. Nunca tenha dúvida sobre a forma de chamar uma função em Go. Sua IDE deve descrever uma função em termos de sua documentação, seus parâmetros e o que ela retorna.
-   **Ver a definição da função**. Se não tiver certeza sobre como uma função funciona, você deve ser capaz de ir para o código fonte de descobrir por si facilmente.
-   **Encontrar usos de um símbolo**. Ser capaz de ver o contexto de uma função sendo chamada pode te ajudar com o processo de refatoração.

Dominar suas ferramentas vai te ajudar a concentrar no código e reduzir a troca de contexto.

## Resumindo

Nesse ponto você já deve ter o Go instalado, um editor disponível e algumas ferramentas básicas configuradas. O Go tem um ecossistema enorme de produtos feitos por outras pessoas. Identificamos alguns componentes úteis aqui, mas você pode encontrar uma lista mais completa no [Awesome Go](https://awesome-go.com).
