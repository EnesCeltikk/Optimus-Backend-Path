# --- Build aşaması ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Bağımlılıkları yükle
COPY go.mod go.sum ./
RUN go mod download

# Tüm kaynak kodunu kopyala
COPY . .

# Binary adı argüman olarak alınabilir
ARG BINARY_NAME=app

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -o /$BINARY_NAME ./main.go

# --- Çalıştırma aşaması ---
FROM alpine:latest
WORKDIR /root/

# Gerekli paketler
RUN apk --no-cache add ca-certificates

# Binary'i kopyala
ARG BINARY_NAME=app
COPY --from=builder /$BINARY_NAME ./

# .env dosyasını çalışma dizinine kopyalamak için (isteğe bağlı)
COPY .env .

# Portu expose et
EXPOSE 8080

# Uygulamayı başlat
CMD ["./app"] 