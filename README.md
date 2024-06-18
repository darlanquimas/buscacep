# BuscaCEP

Este é um projeto em Go que demonstra como consultar informações de CEP utilizando um banco de dados SQLite local e a API do ViaCep.

## Pré-requisitos

Antes de executar o projeto, certifique-se de ter o seguinte instalado em seu ambiente de desenvolvimento:

- Go (versão 1.16 ou superior)
- SQLite (para armazenamento local dos dados)
- Conexão com a internet (para consultar dados via API)

## Instalação

1. Clone o repositório:

   ```bash
   git clone https://github.com/seu-usuario/seu-repositorio.git

   ```

2. Navegue até o diretório do projeto:

   ```bash
   cd buscacep\cmd\viacepapp

   ```

3. Execute o comando para compilar e instalar as dependências do Go:

   ```bash
   go build

   .\viacepapp.exe <CEP1> <CEP2>
   ```

# Configuração

Antes de executar o projeto pela primeira vez, é necessário configurar o banco de dados SQLite e a inicialização das tabelas. Certifique-se de que o SQLite esteja configurado corretamente e atualize o caminho do banco de dados conforme necessário no código.

````go
// Exemplo de inicialização do banco de dados
db, err := InitDB("caminho/para/seu/banco-de-dados.db")
if err != nil {
    fmt.Fprintf(os.Stderr, "Erro ao inicializar banco de dados: %v\n", err)
    os.Exit(1)
}
defer db.Close()


# Uso
Para executar o projeto, use o comando abaixo:

    ```bash
    go run main.go <CEP1> <CEP2> ...
````
