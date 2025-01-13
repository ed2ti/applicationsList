FROM golang:1.23

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos go.mod e go.sum para o container
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o restante do código fonte para o container
COPY . .

# Compila o aplicativo
RUN GOOS=linux GOARCH=amd64 go build -o app

# Define o comando padrão
CMD ["./app"]

