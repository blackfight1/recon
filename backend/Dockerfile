FROM golang:1.21-alpine

WORKDIR /app

# 安装必要工具
RUN apk add --no-cache git docker-cli ca-certificates

# 复制所有文件
COPY . .

# 下载依赖、生成 go.sum 并编译
RUN go mod download && \
    go mod tidy && \
    go build -v -o main . && \
    chmod +x main

# 创建数据目录
RUN mkdir -p /data

# 复制配置文件
RUN cp config.example.yaml config.yaml

EXPOSE 8000

CMD ["./main"]
