#!/bin/bash

echo "🔄 重新构建项目..."

# 停止服务
echo "1️⃣ 停止服务..."
docker-compose down

# 重新构建
echo "2️⃣ 重新构建镜像..."
docker-compose build --no-cache

# 启动服务
echo "3️⃣ 启动服务..."
docker-compose up -d

# 等待服务启动
echo "4️⃣ 等待服务启动（10秒）..."
sleep 10

# 检查服务状态
echo "5️⃣ 检查服务状态..."
docker-compose ps

echo ""
echo "✅ 重新构建完成！"
echo "   前端: http://localhost:8080"
echo "   后端: http://localhost:8000"
