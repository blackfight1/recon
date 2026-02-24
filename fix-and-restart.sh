#!/bin/bash

echo "🔧 修复并重启服务..."
echo ""

# 停止服务
echo "1️⃣ 停止现有服务..."
docker-compose down

# 清理
echo "2️⃣ 清理旧容器和镜像..."
docker-compose rm -f

# 重新构建
echo "3️⃣ 重新构建镜像..."
docker-compose build --no-cache

# 启动服务
echo "4️⃣ 启动服务..."
docker-compose up -d

# 等待服务启动
echo "5️⃣ 等待服务启动（15秒）..."
sleep 15

# 检查状态
echo "6️⃣ 检查服务状态..."
docker-compose ps

echo ""
echo "7️⃣ 检查后端健康..."
if curl -s http://localhost:8000/health | grep -q "ok"; then
    echo "✅ 后端 API 正常"
else
    echo "❌ 后端 API 异常，查看日志："
    docker-compose logs --tail=20 backend
fi

echo ""
echo "8️⃣ 检查前端..."
if curl -s -I http://localhost:8080 | grep -q "200"; then
    echo "✅ 前端正常"
else
    echo "❌ 前端异常，查看日志："
    docker-compose logs --tail=20 frontend
fi

echo ""
echo "✨ 修复完成！"
echo ""
echo "📱 本地访问: http://localhost:8080"
echo "📱 VPS 访问: http://YOUR_VPS_IP:8080"
echo ""
echo "⚠️  如果 VPS 无法访问，请检查："
echo "   1. 防火墙: sudo ufw allow 8080/tcp"
echo "   2. 云服务商安全组配置"
echo "   3. 查看完整指南: cat VPS_DEPLOYMENT.md"
echo ""
echo "📋 查看日志: docker-compose logs -f"
