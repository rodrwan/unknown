# Etapa de construcción
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Etapa final
FROM alpine:latest

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/main .
# Copiar archivos estáticos y templates
COPY --from=builder /app/static ./static

# Exponer el puerto
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./main"]