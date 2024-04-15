# Bước 1: Build ứng dụng
FROM golang:1.20.4 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o redis-go .

# Bước 2: Tạo image chạy ứng dụng
FROM nginx:alpine
COPY --from=builder /app/Redis-api-golang /usr/local/bin/redis-nginx
COPY nginx.conf /etc/nginx/nginx.conf
# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# WORKDIR /root/
# COPY --from=builder /app/redis-go .

# EXPOSE 4000 4001 4002 4003 4004
EXPOSE 3000
# CMD ["./redis-go"]
CMD ["sh", "-c", "/usr/local/bin/redis-nginx & nginx -g 'daemon off;'"]