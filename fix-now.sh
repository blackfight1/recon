#!/bin/bash

echo "🔧 快速修复当前问题..."
echo ""

# 1. 确保工具有执行权限
echo "1️⃣ 添加工具执行权限..."
chmod +x /root/go/bin/* 2>/dev/null || true
ls -lh /root/go/bin/

echo ""
echo "2️⃣ 停止服务..."
docker-compose down

echo ""
echo "3️⃣ 重新构建（包含前端暗色主题）..."
docker-compose build --no-cache

echo ""
echo "4️⃣ 启动服务..."
docker-compose up -d

echo ""
echo "5️⃣ 等待启动（15秒）..."
sleep 15

echo ""
echo "6️⃣ 验证工具挂载..."
docker-compose exec -T backend sh -c "ls -lh /usr/local/bin/subfinder /usr/local/bin/samoscout /usr/local/bin/ksubdomain /usr/local/bin/httpx"

echo ""
echo "7️⃣ 测试工具..."
docker-compose exec -T backend subfinder -version || echo "⚠️  Subfinder 无法执行"
docker-compose exec -T backend httpx -version || echo "⚠️  Httpx 无法执行"

echo ""
echo "8️⃣ 查看服务状态..."
docker-compose ps

echo ""
echo "✅ 修复完成！"
echo ""
echo "访问前端查看暗色主题："
echo "  http://$(hostname -I | awk '{print $1}'):8080"
echo ""
echo "如果还有问题，请运行："
echo "  docker-compose logs -f backend"
