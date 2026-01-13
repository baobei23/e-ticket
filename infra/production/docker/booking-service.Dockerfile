FROM golang:1.25.5 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/services/booking-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o booking-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/booking-service/booking-service .
ENTRYPOINT ["./booking-service"]