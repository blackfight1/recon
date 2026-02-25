#!/bin/bash

echo "🔄 快速重启服务..."
echo ""

# 停止容器
echo "1️⃣ 停止容器..."
docker-compose down

# 启动容器（使用现有镜像）
echo "2️⃣ 启动容器..."
docker-compose up -d

# 等待服务启动
echo "3️⃣ 等待服务启动（15秒）..."
sleep 15

# 检查服务状态
echo ""
echo "📊 容器状态："
docker-compose ps

echo ""
echo "🔍 检查服务..."

# 检查后端
if curl -s http://localhost:8000/health 2>/dev/null | grep -q "ok"; then
    echo "✅ 后端 API 正常"
else
    echo "⚠️  后端 API 可能未就绪"
    docker-compose logs --tail=10 backend
fi

# 检查前端
if curl -s -I http://localhost:8080 2>/dev/null | grep -q "200"; then
    echo "✅ 前端正常"
else
    echo "⚠️  前端可能未就绪"
fi

echo ""
echo "✨ 重启完成！"
echo ""
echo "📱 访问地址: http://localhost:8080"
echo "📋 查看日志: docker-compose logs -f"
echo ""
