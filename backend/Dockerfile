FROM golang:1.21

WORKDIR /app

# 安装必要工具
RUN apk add --no-cache git ca-certificates

# 先复制 go.mod
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
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
