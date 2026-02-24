#!/bin/bash

echo "🚀 启动自动化侦查平台..."

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 检查配置文件
if [ ! -f "backend/config.yaml" ]; then
    echo "📝 创建配置文件..."
    cp backend/config.example.yaml backend/config.yaml
    echo "⚠️  请编辑 backend/config.yaml 配置通知 Webhook"
fi

# 创建数据目录
mkdir -p data

# 启动服务
echo "🐳 启动 Docker 容器..."
docker-compose up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "📊 检查服务状态..."
docker-compose ps

# 检查后端健康状态
echo "🔍 检查后端 API..."
if curl -s http://localhost:8000/health | grep -q "ok"; then
    echo "✅ 后端 API 正常"
else
    echo "⚠️  后端 API 可能未就绪，请稍后再试"
fi

echo ""
echo "✨ 部署完成！"
echo ""
echo "📱 访问地址: http://localhost:8080"
echo "📚 API 文档: http://localhost:8000/api"
echo ""
echo "📋 常用命令:"
echo "  查看日志: docker-compose logs -f"
echo "  停止服务: docker-compose down"
echo "  重启服务: docker-compose restart"
echo ""
