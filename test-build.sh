#!/bin/bash

echo "🔍 测试后端构建..."
echo ""

cd backend

echo "1️⃣ 检查 Go 环境..."
docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine go version

echo ""
echo "2️⃣ 检查文件..."
ls -lh *.go

echo ""
echo "3️⃣ 测试 go mod download..."
docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine sh -c "go mod download && echo '✅ 依赖下载成功'"

echo ""
echo "4️⃣ 测试 go mod tidy..."
docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine sh -c "go mod tidy && echo '✅ go mod tidy 成功'"

echo ""
echo "5️⃣ 测试编译..."
docker run --rm -v $(pwd):/app -w /app golang:1.21-alpine sh -c "CGO_ENABLED=0 GOOS=linux go build -o main . && ls -lh main && echo '✅ 编译成功'"

echo ""
echo "6️⃣ 清理测试文件..."
rm -f main

cd ..

echo ""
echo "✅ 测试完成！如果所有步骤都成功，说明代码没问题。"
echo "   如果有失败，请查看错误信息。"
