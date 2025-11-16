FROM golang:1.25.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8083
CMD ["./server"]
