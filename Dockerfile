FROM golang:1.25-alpine

# Diretório de trabalho dentro do container
WORKDIR /app

# Copiar ficheiros de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o resto do código
COPY . .

# Compilar a aplicação
RUN go build -o stand_api .

# Comando para arrancar a aplicação
CMD ["./stand_api"]
