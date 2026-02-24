#!/bin/bash

echo "🔄 完全重建项目..."
echo ""

# 停止所有容器
echo "1️⃣ 停止所有容器..."
docker-compose down

# 删除旧镜像
echo "2️⃣ 删除旧镜像..."
docker-compose rm -f
docker rmi recon_backend recon_frontend 2>/dev/null || true

# 重新构建（不使用缓存）
echo "3️⃣ 重新构建镜像（这可能需要几分钟）..."
docker-compose build --no-cache

# 启动服务
echo "4️⃣ 启动所有服务..."
docker-compose up -d

# 等待服务启动
echo "5️⃣ 等待服务启动（20秒）..."
sleep 20

# 检查容器状态
echo ""
echo "📊 容器状态："
docker-compose ps

echo ""
echo "🔍 检查服务..."

# 检查后端
echo ""
echo "后端健康检查："
if curl -s http://localhost:8000/health 2>/dev/null | grep -q "ok"; then
    echo "✅ 后端 API 正常运行"
else
    echo "❌ 后端 API 异常"
    echo "后端日志（最后 20 行）："
    docker-compose logs --tail=20 backend
fi

# 检查前端
echo ""
echo "前端检查："
if curl -s -I http://localhost:8080 2>/dev/null | grep -q "200"; then
    echo "✅ 前端正常运行"
else
    echo "❌ 前端异常"
    echo "前端日志（最后 20 行）："
    docker-compose logs --tail=20 frontend
fi

# 检查数据库
echo ""
echo "数据库检查："
if docker-compose exec -T postgres psql -U recon -d recon -c "SELECT 1" >/dev/null 2>&1; then
    echo "✅ 数据库正常运行"
else
    echo "❌ 数据库异常"
    docker-compose logs --tail=20 postgres
fi

# 显示访问信息
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✨ 部署完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 获取服务器 IP
SERVER_IP=$(curl -s ifconfig.me 2>/dev/null || curl -s ipinfo.io/ip 2>/dev/null || echo "YOUR_SERVER_IP")

echo "📱 访问地址："
echo "   本地访问: http://localhost:8080"
echo "   外网访问: http://$SERVER_IP:8080"
echo ""
echo "🔧 API 接口："
echo "   本地: http://localhost:8000/api"
echo "   外网: http://$SERVER_IP:8000/api"
echo ""
echo "💡 健康检查："
echo "   curl http://localhost:8000/health"
echo ""
echo "📋 常用命令："
echo "   查看日志: docker-compose logs -f"
echo "   查看状态: docker-compose ps"
echo "   停止服务: docker-compose down"
echo "   重启服务: docker-compose restart"
echo ""

# 检查端口监听
echo "🔌 端口监听状态："
netstat -tlnp 2>/dev/null | grep -E ':(8080|8000|5432)' || ss -tlnp | grep -E ':(8080|8000|5432)'

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
