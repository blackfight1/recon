#!/bin/bash

echo "🔍 诊断后端问题..."
echo ""

echo "1️⃣ 容器状态："
docker-compose ps
echo ""

echo "2️⃣ 后端最近日志（最后50行）："
docker-compose logs --tail=50 backend
echo ""

echo "3️⃣ 数据库连接测试："
docker-compose exec postgres psql -U recon -d recon -c "SELECT version();" 2>&1 || echo "❌ 数据库连接失败"
echo ""

echo "4️⃣ 检查配置文件："
docker-compose exec backend ls -la /app/config.yaml 2>&1 || echo "❌ 配置文件不存在"
echo ""

echo "5️⃣ 检查环境变量："
docker-compose exec backend env | grep DB_ 2>&1 || echo "❌ 无法获取环境变量"
echo ""

echo "📝 如需实时查看日志，运行："
echo "   docker-compose logs -f backend"
