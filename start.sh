#!/bin/bash

set -e

echo "🚀 启动自动化侦查平台..."
echo ""

# 检查 Docker
echo "1️⃣ 检查 Docker 环境..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

echo "✅ Docker 环境正常"
echo ""

# 检查配置文件
echo "2️⃣ 检查配置文件..."
if [ ! -f "backend/config.yaml" ]; then
    echo "⚠️  未找到配置文件，复制默认配置..."
    cp backend/config.example.yaml backend/config.yaml
    echo "✅ 已创建 backend/config.yaml，请根据需要修改"
fi
echo ""

# 启动服务
echo "3️⃣ 启动所有服务..."
docker-compose up -d

echo ""
echo "4️⃣ 等待服务启动（20秒）..."
sleep 20

# 检查服务状态
echo ""
echo "5️⃣ 检查服务状态..."
docker-compose ps

echo ""
echo "✅ 启动完成！"
echo ""
echo "📍 访问地址："
echo "   前端: http://localhost:8080"
echo "   后端: http://localhost:8000"
echo "   健康检查: http://localhost:8000/health"
echo ""
echo "📝 常用命令："
echo "   查看日志: docker-compose logs -f"
echo "   停止服务: docker-compose down"
echo "   重启服务: docker-compose restart"
echo "   重新构建: ./rebuild.sh"
echo ""
