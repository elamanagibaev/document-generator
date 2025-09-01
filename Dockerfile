
FROM golang:1.24 AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o document-generator ./cmd/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/document-generator /app/document-generator
COPY templates/ /app/templates/


ENTRYPOINT ["/app/document-generator"]
