FROM golang:1.17.0-alpine3.14 as builder

COPY . /app

WORKDIR /app

RUN apk -U add librdkafka-dev pkgconf gcc musl-dev

RUN go mod tidy
RUN CGO_ENABLED=1 go build -tags musl -o bench main.go

FROM alpine:3.14
COPY --from=builder /app/bench /bin/app

USER nobody

ENTRYPOINT ["app"]
