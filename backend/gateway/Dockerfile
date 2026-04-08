FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .
RUN rm -f go.sum && go mod tidy && go build -o bin/gateway main.go

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/gateway .
COPY --from=builder /app/config.yaml .

EXPOSE 8080

CMD ["./gateway"]
