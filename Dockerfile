FROM golang:alpine AS builder

RUN apk update

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux  go build -o bin/career-path cmd/app/main.go

FROM surnet/alpine-wkhtmltopdf:3.20.2-0.12.6-full as wkhtmltopdf
FROM alpine AS app

RUN apk update && apk add --no-cache ca-certificates make && update-ca-certificates
RUN apk add --no-cache \
    libstdc++ \
    libx11 \
    libxrender \
    libxext \
    libssl3 \
    ca-certificates \
    fontconfig \
    freetype \
    ttf-dejavu \
    ttf-droid \
    ttf-freefont \
    ttf-liberation \
    # more fonts
  && apk add --no-cache --virtual .build-deps \
    msttcorefonts-installer \
  # Install microsoft fonts
  && update-ms-fonts \
  && fc-cache -f \
  # Clean up when done
  && rm -rf /tmp/* \
  && apk del .build-deps

# Copy wkhtmltopdf files from docker-wkhtmltopdf image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=wkhtmltopdf /bin/wkhtmltoimage /bin/wkhtmltoimage
COPY --from=wkhtmltopdf /lib/libwkhtmltox* /lib/

WORKDIR /app

COPY --from=builder /app/bin/career-path .
COPY --from=builder /app/.env .env
COPY --from=builder /app/Makefile ./Makefile
COPY --from=builder /app/database/migrations ./database/migrations
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/resume_template.gohtml ./resume_template.gohtml

RUN chmod +x /usr/local/bin/migrate

ENTRYPOINT ["./career-path"]
