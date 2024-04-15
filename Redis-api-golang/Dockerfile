# Bước 1: Build ứng dụng
FROM golang:1.20.4 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-redis .

# Bước 2: Tạo image chạy ứng dụng
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/go-redis .

EXPOSE 3000
CMD ["./go-redis"]