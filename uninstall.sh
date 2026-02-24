#!/bin/bash

echo "🗑️  自动化侦查平台 - 卸载脚本"
echo "================================"
echo ""
echo "⚠️  警告：此操作将："
echo "   1. 停止所有容器"
echo "   2. 删除所有容器"
echo "   3. 删除所有镜像"
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

# 3. 删除镜像
echo "3️⃣ 删除项目镜像..."
docker rmi recon_backend recon_frontend 2>/dev/null || true
docker rmi recon-backend recon-frontend 2>/dev/null || true

# 4. 删除扫描工具镜像（可选）
echo ""
read -p "是否删除扫描工具镜像？(yes/no): " remove_tools

if [ "$remove_tools" = "yes" ]; then
    echo "4️⃣ 删除扫描工具镜像..."
    docker rmi projectdiscovery/subfinder:latest 2>/dev/null || true
    docker rmi projectdiscovery/httpx:latest 2>/dev/null || true
    docker rmi tomnomnom/assetfinder:latest 2>/dev/null || true
    docker rmi postgres:15-alpine 2>/dev/null || true
fi

# 5. 清理网络
echo "5️⃣ 清理网络..."
docker network rm recon_recon-network 2>/dev/null || true

# 6. 清理数据目录
echo "6️⃣ 清理数据目录..."
if [ -d "data" ]; then
    rm -rf data
fi

# 7. 清理配置文件
echo ""
read -p "是否删除配置文件 backend/config.yaml？(yes/no): " remove_config

if [ "$remove_config" = "yes" ]; then
    echo "7️⃣ 删除配置文件..."
    rm -f backend/config.yaml
fi

# 8. 显示剩余内容
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 卸载完成！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "已删除："
echo "  ✓ Docker 容器"
echo "  ✓ Docker 镜像"
echo "  ✓ 数据卷"
echo "  ✓ 网络"
echo "  ✓ 数据目录"
if [ "$remove_config" = "yes" ]; then
    echo "  ✓ 配置文件"
fi
echo ""
echo "保留的文件："
echo "  • 源代码文件"
echo "  • 文档文件"
if [ "$remove_config" != "yes" ]; then
    echo "  • backend/config.yaml"
fi
echo ""
echo "💡 如需完全删除项目，请手动删除项目目录："
echo "   cd .. && rm -rf recon-platform"
echo ""
