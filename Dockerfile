FROM golang:alpine AS builder

RUN apk update

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux  go build -o bin/career-path cmd/app/main.go

FROM alpine AS app

RUN apk update && apk add --no-cache ca-certificates make && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/career-path .
COPY --from=builder /app/.env .env
COPY --from=builder /app/Makefile ./Makefile
COPY --from=builder /app/database/migrations ./database/migrations
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

RUN chmod +x /usr/local/bin/migrate

ENTRYPOINT ["./career-path"]