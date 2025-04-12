FROM --platform=linux/amd64 golang:alpine AS builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o sync_users ./cmd/sync_users

FROM alpine

WORKDIR /app

COPY --from=builder /app/sync_users /sync_users
RUN chmod +x /sync_users

COPY --from=builder /app/database/migrations /app/database/migrations

CMD ["/sync_users"] 