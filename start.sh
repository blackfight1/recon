#!/bin/bash

echo "🚀 启动自动化侦查平台..."
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装"
    exit 1
fi

# 检查配置文件
if [ ! -f "backend/config.yaml" ]; then
    echo "📝 创建配置文件..."
    cp backend/config.example.yaml backend/config.yaml
fi

# 创建数据目录
mkdir -p data

# 启动服务
echo "🐳 启动 Docker 容器..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 15

# 检查服务状态
echo "📊 检查服务状态..."
docker-compose ps

# 检查后端健康状态
echo ""
echo "🔍 检查后端 API..."
if curl -s http://localhost:8000/health 2>/dev/null | grep -q "ok"; then
    echo "✅ 后端 API 正常"
else
    echo "⚠️  后端 API 可能未就绪，查看日志："
    docker-compose logs --tail=20 backend
fi

echo ""
echo "✨ 部署完成！"
echo ""
echo "📱 访问地址: http://localhost:8080"
echo ""
echo "📋 常用命令:"
echo "  查看日志: docker-compose logs -f"
echo "  停止服务: docker-compose down"
echo "  重启服务: docker-compose restart"
echo "  卸载系统: ./uninstall.sh"
echo ""
