# ---------- BUILD STAGE ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

# зависимости (важно для alpine)
RUN apk add --no-cache git

# копируем модули
COPY go.mod go.sum ./
RUN go mod download

# копируем весь код
COPY . .

# собираем бинарь
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/eurovision-voting/main.go

# ---------- RUN STAGE ----------
FROM alpine:latest

WORKDIR /app

# сертификаты (нужно для https/db connections)
RUN apk --no-cache add ca-certificates

# копируем бинарь из builder
COPY --from=builder /app/app .

# порт (Render/локал)
ENV PORT=8080

EXPOSE 8080

# запуск
CMD ["./app"]