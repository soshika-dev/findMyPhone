# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux

# Ensure module files are up to date before building
RUN go mod tidy

RUN go build -o server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/server ./

EXPOSE 8080

ENTRYPOINT ["/app/server"]
