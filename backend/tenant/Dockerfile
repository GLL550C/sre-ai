FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .
RUN rm -f go.sum && go mod tidy && go build -o bin/tenant main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/tenant .
COPY --from=builder /app/config.yaml .

EXPOSE 8083

CMD ["./tenant"]
