FROM golang:1.26.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY .env /app/.env

EXPOSE 8080

ENV APP_ENV=production

ENTRYPOINT ["/app/server"]
