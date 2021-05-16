# build stage
FROM golang:1.16-alpine3.13 AS builder
RUN apk --update add ca-certificates
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main server.go

# Run stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080
CMD [ "/app/main" ]