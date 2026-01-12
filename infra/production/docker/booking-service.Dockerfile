FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/services/booking-service
RUN CGO_ENABLED=0 GOOS=linux go build -o booking-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/booking-service/booking-service .
ENTRYPOINT ["./booking-service"]