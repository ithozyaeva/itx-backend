FROM --platform=linux/amd64 golang:alpine AS builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the application
RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/backend

FROM alpine

WORKDIR /app

COPY --from=builder /app/backend /backend
RUN chmod +x /backend

# Copy migrations and other necessary files
COPY --from=builder /app/database/migrations /app/database/migrations
COPY --from=builder /app/config /app/config

CMD ["/backend"]