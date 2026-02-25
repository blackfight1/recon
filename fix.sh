#!/bin/bash

echo "🔧 修复并重新部署..."
echo ""

# 给所有脚本添加执行权限
echo "1️⃣ 添加脚本执行权限..."
chmod +x *.sh

# 停止所有服务
echo "2️⃣ 停止所有服务..."
docker-compose down

# 清理旧的构建缓存
echo "3️⃣ 清理构建缓存..."
docker-compose build --no-cache backend

# 启动服务
echo "4️⃣ 启动服务..."
docker-compose up -d

# 等待启动
echo "5️⃣ 等待服务启动（20秒）..."
sleep 20

# 检查状态
echo "6️⃣ 检查服务状态..."
docker-compose ps

echo ""
echo "7️⃣ 后端日志（最后30行）："
docker-compose logs --tail=30 backend

echo ""
echo "✅ 修复完成！"
echo ""
echo "如果后端仍在重启，请运行："
echo "  ./debug.sh"
echo ""
echo "实时查看日志："
echo "  ./logs.sh"
