FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

RUN apk update && apk upgrade && apk add --no-cache bash curl

ADD https://github.com/pressly/goose/releases/download/v3.24.1/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

RUN curl -fsSL -o sqlc.tar.gz https://github.com/sqlc-dev/sqlc/releases/download/v1.28.0/sqlc_1.28.0_linux_amd64.tar.gz && \
    tar -xzf sqlc.tar.gz -C /bin && \
    rm sqlc.tar.gz && \
    chmod +x /bin/sqlc && \
    /bin/sqlc version

EXPOSE 8080

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]
