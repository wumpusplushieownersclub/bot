FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY src src

RUN go build -o wumpus-bot ./src

FROM alpine

#RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/wumpus-bot .

ENTRYPOINT [ "./wumpus-bot" ]