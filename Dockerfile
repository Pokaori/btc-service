FROM golang:1.13-alpine AS build_base

RUN apk add --no-cache 

WORKDIR ~/bitcoin-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /bitcoin-service


FROM alpine:latest
RUN apk add --no-cache 

COPY --from=build_base /bitcoin-service /app/bitcoin-service
COPY .env /.env

EXPOSE 8000

CMD ["/app/bitcoin-service"]