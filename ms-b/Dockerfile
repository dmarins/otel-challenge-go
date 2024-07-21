FROM golang:1.22 AS builder
WORKDIR /app
RUN apt-get install -y ca-certificates
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ms-b

FROM scratch AS runner
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/ms-b /ms-b
ENTRYPOINT ["/ms-b"]
