FROM golang:1.21-alpine

WORKDIR /app

# 安装必要工具（这层会被缓存）
RUN apk add --no-cache git docker-cli ca-certificates

# 先复制 go.mod 和 go.sum（如果存在），利用 Docker 缓存
COPY go.mod ./

# 下载依赖（这层会被缓存，除非 go.mod 改变）
RUN go mod download

# 复制源代码（只有代码改变时才重新构建）
COPY . .

# 生成 go.sum 并编译
RUN go mod tidy && \
    go build -v -o main . && \
    chmod +x main

# 创建数据目录
RUN mkdir -p /data

# 复制配置文件
RUN cp config.example.yaml config.yaml

EXPOSE 8000

CMD ["./main"]
