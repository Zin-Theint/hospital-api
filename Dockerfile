#build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o hospital-api ./cmd/server

#run
FROM alpine:3.19
RUN adduser -D api
WORKDIR /home/api
COPY --from=builder /app/hospital-api .
USER api
EXPOSE 8080
ENTRYPOINT ["./hospital-api"]
