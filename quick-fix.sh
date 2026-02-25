#!/bin/bash

echo "⚡ 快速修复并重启后端..."
echo ""

# 停止后端
echo "1️⃣ 停止后端容器..."
docker-compose stop backend

# 重新构建后端
echo "2️⃣ 重新构建后端..."
docker-compose build --no-cache backend

# 启动后端
echo "3️⃣ 启动后端..."
docker-compose up -d backend

# 等待启动
echo "4️⃣ 等待后端启动（10秒）..."
sleep 10

# 检查状态
echo "5️⃣ 检查后端状态..."
docker-compose ps backend

echo ""
echo "6️⃣ 后端日志（最后30行）："
docker-compose logs --tail=30 backend

echo ""
echo "✅ 修复完成！"
echo ""
echo "测试快速扫描："
echo "  访问 http://你的IP:8080"
echo "  点击'快速扫描'菜单"
echo ""
echo "查看实时日志："
echo "  ./logs.sh"
