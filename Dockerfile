# build stage
FROM golang:1.16-alpine3.13 AS builder
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main server.go

# Run stage
FROM alpine:3.13
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD [ "/app/main" ]