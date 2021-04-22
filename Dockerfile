FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git tzdata ca-certificates
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/diosteama

FROM scratch
COPY --from=builder /app/diosteama /diosteama
COPY --from=builder /app/resources/ /resources/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/diosteama"]
