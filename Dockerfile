# Etapa 1: build
FROM golang:1.25 AS builder

WORKDIR /app

# Copiar go.mod e go.sum primeiro (para cache eficiente)
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante
COPY . .

# Compilar binário
RUN go build -o bot .

# Etapa 2: imagem final
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

# Copiar binário do builder
COPY --from=builder /app/bot .

COPY .env ./

# Variável de ambiente para token
# ENV TELEGRAM_BOT_TOKEN=""

CMD ["./bot"]
