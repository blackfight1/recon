#!/bin/bash

echo "🚀 部署重构后的项目..."
echo ""

# 1. 检查工具
echo "1️⃣ 检查扫描工具..."
chmod +x check-tools.sh
./check-tools.sh

if [ $? -ne 0 ]; then
    echo ""
    echo "❌ 工具检查失败，请先安装所需工具"
    exit 1
fi

echo ""
echo "2️⃣ 停止旧服务..."
docker-compose down

echo ""
echo "3️⃣ 重新构建镜像..."
docker-compose build --no-cache backend

echo ""
echo "4️⃣ 启动所有服务..."
docker-compose up -d

echo ""
echo "5️⃣ 等待服务启动（20秒）..."
sleep 20

echo ""
echo "6️⃣ 检查服务状态..."
docker-compose ps

echo ""
echo "7️⃣ 后端日志（最后30行）："
docker-compose logs --tail=30 backend

echo ""
echo "✅ 部署完成！"
echo ""
echo "访问地址："
echo "  前端: http://你的IP:8080"
echo "  后端: http://你的IP:8000"
echo ""
echo "测试快速扫描："
echo "  1. 访问前端页面"
echo "  2. 点击'快速扫描'"
echo "  3. 输入域名测试"
echo ""
echo "查看实时日志："
echo "  docker-compose logs -f backend"
