FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/services/event-service
RUN CGO_ENABLED=0 GOOS=linux go build -o event-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/event-service/event-service .
ENTRYPOINT ["./event-service"]