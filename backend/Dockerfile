FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN rm -f /run/flags/setup.flag

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release-backend ./cmd/backend

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/release-backend .

EXPOSE $PORT

# Should be changed when an update is made
ENV VERSION=1.0.0

CMD ["./release-backend"]