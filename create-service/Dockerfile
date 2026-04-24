# --- Stage 1: Base environment and dependencies ---
FROM golang:1.25.4-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# --- Stage 2: Unit Testing (Mandatory) ---
FROM base AS tester
COPY . .
# If tests fail, the build stops here
RUN go test -v ./...

# --- Stage 3: Builder (Inherits from tester) ---
# This ensures builder only runs if tester succeeds
FROM tester AS builder
# Code is already copied from the tester stage
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# --- Stage 4: Final production image ---
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
USER nonroot:nonroot
CMD ["./server"]