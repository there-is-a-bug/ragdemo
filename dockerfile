# 构建阶段 go1.26.8
FROM golang:1.26.8-alpine AS builder

WORKDIR /app

# 安装git依赖
RUN apk add --no-cache git

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 编译静态二进制，关闭CGO
RUN CGO_ENABLED=0 GOOS=linux go build -o rag .

# 运行阶段，轻量alpine
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/rag ./rag
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./rag"]
