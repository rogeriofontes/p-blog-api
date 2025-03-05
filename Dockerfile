# Usar a imagem oficial do Golang
FROM golang:1.23-alpine

# Definir diretório de trabalho dentro do container
WORKDIR /app

# Copiar os arquivos do projeto para dentro do container
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante do código-fonte
COPY . .

# Compilar a aplicação
RUN go build -o main .

# Expor a porta da API
EXPOSE 8484

# Comando para rodar a API
CMD ["/app/main"]
