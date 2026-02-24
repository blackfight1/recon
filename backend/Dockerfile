FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要工具
RUN apk add --no-cache git docker-cli

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates docker-cli

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config.example.yaml ./config.yaml

# 创建数据目录
RUN mkdir -p /data

EXPOSE 8000

CMD ["./main"]
