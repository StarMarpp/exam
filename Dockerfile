FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /appsource 

COPY go.mod go.sum ./

COPY ./bin ./bin
COPY ./internal ./internal
COPY ./migrations ./migrations
COPY ./frontend ./frontend

RUN go mod tidy

RUN go build -o app ./internal
RUN chmod +x ./app

FROM alpine:3.21

WORKDIR /myapp
COPY --from=builder /appsource/app .

EXPOSE 8080

CMD ["./app"]