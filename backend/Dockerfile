FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要工具
RUN apk add --no-cache git docker-cli curl

# 复制 go.mod 和 go.sum（如果存在）
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 生成 go.sum 并编译
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main . && \
    chmod +x main && \
    ls -lh main

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates docker-cli curl

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config.example.yaml ./config.yaml

# 验证文件存在
RUN ls -lh /app/main && chmod +x /app/main

# 创建数据目录
RUN mkdir -p /data

EXPOSE 8000

CMD ["./main"]
