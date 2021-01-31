FROM golang:1.15.6-alpine3.13 as builder

COPY . /app

WORKDIR /app

RUN apk -U add librdkafka-dev pkgconf gcc musl-dev

RUN go mod download
RUN CGO_ENABLED=1 go build -tags musl -o bench main.go

FROM alpine:3.13
COPY --from=builder /app/bench /bin/app

USER nobody

ENTRYPOINT ["app"]
