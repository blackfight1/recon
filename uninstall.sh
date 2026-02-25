#!/bin/bash

echo "🗑️  自动化侦查平台 - 卸载脚本"
echo "================================"
echo ""
echo "⚠️  警告：此操作将："
echo "   1. 停止所有容器"
echo "   2. 删除所有容器"
echo "   3. 删除项目镜像（保留基础镜像）"
echo "   4. 删除所有数据卷（包括数据库数据）"
echo "   5. 删除网络"
echo ""
read -p "确定要继续吗？(yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ 取消卸载"
    exit 0
fi

echo ""
echo "开始卸载..."
echo ""

# 1. 停止并删除容器
echo "1️⃣ 停止并删除容器..."
docker-compose down

# 2. 删除数据卷
echo "2️⃣ 删除数据卷（包括所有数据）..."
docker-compose down -v

# 3. 删除项目镜像（保留基础镜像）
echo "3️⃣ 删除项目镜像（保留基础镜像 golang、postgres、node、nginx）..."
docker rmi recon_backend recon_frontend 2>/dev/null || true
docker rmi recon-backend recon-frontend 2>/dev/null || true

# 4. 询问是否删除扫描工具镜像
echo ""
echo "📦 扫描工具镜像（建议保留）："
echo "   - projectdiscovery/subfinder:latest"
echo "   - projectdiscovery/httpx:latest"
echo "   - tomnomnom/assetfinder:latest"
echo ""
read -p "是否删除扫描工具镜像？(yes/no，默认 no): " remove_tools

if [ "$remove_tools" = "yes" ]; then
    echo "4️⃣ 删除扫描工具镜像..."
    docker rmi projectdiscovery/subfinder:latest 2>/dev/null || true
    docker rmi projectdiscovery/httpx:latest 2>/dev/null || true
    docker rmi tomnomnom/assetfinder:latest 2>/dev/null || true
else
    echo "4️⃣ 保留扫描工具镜像"
fi

# 5. 询问是否删除基础镜像
echo ""
echo "📦 基础镜像（强烈建议保留）："
echo "   - golang:1.21-alpine"
echo "   - postgres:15-alpine"
echo "   - node:18-alpine"
echo "   - nginx:alpine"
echo ""
read -p "是否删除基础镜像？(yes/no，默认 no): " remove_base

if [ "$remove_base" = "yes" ]; then
    echo "5️⃣ 删除基础镜像..."
    docker rmi golang:1.21-alpine 2>/dev/null || true
    docker rmi postgres:15-alpine 2>/dev/null || true
    docker rmi node:18-alpine 2>/dev/null || true
    docker rmi nginx:alpine 2>/dev/null || true
else
    echo "5️⃣ 保留基础镜像"
fi

# 6. 清理网络
echo "6️⃣ 清理网络..."
docker network rm recon_recon-network 2>/dev/null || true

# 7. 清理数据目录
echo "7️⃣ 清理数据目录..."
if [ -d "data" ]; then
    rm -rf data
fi

# 8. 清理配置文件
echo ""
read -p "是否删除配置文件 backend/config.yaml？(yes/no，默认 no): " remove_config

if [ "$remove_config" = "yes" ]; then
    echo "8️⃣ 删除配置文件..."
    rm -f backend/config.yaml
else
    echo "8️⃣ 保留配置文件"
fi

# 9. 显示剩余内容
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 卸载完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "已删除："
echo "  ✓ Docker 容器"
echo "  ✓ 项目镜像"
echo "  ✓ 数据卷"
echo "  ✓ 网络"
echo "  ✓ 数据目录"
if [ "$remove_config" = "yes" ]; then
    echo "  ✓ 配置文件"
fi
echo ""
echo "保留的内容："
echo "  • 源代码文件"
echo "  • 文档文件"
if [ "$remove_config" != "yes" ]; then
    echo "  • backend/config.yaml"
fi
if [ "$remove_tools" != "yes" ]; then
    echo "  • 扫描工具镜像"
fi
if [ "$remove_base" != "yes" ]; then
    echo "  • 基础镜像（golang、postgres、node、nginx）"
fi
echo ""
echo "💡 查看保留的镜像："
echo "   docker images"
echo ""
echo "💡 如需完全删除项目，请手动删除项目目录："
echo "   cd .. && rm -rf recon"
echo ""
