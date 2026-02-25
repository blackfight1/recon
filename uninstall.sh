#!/bin/bash

echo "🗑️  卸载自动化侦查平台..."
echo ""

# 检查参数
FULL_UNINSTALL=false
if [ "$1" == "--full" ]; then
    FULL_UNINSTALL=true
    echo "⚠️  完全卸载模式：将删除所有数据和镜像"
else
    echo "ℹ️  标准卸载模式：保留数据和基础镜像"
    echo "   使用 --full 参数进行完全卸载"
fi
echo ""

# 停止并删除容器
echo "1️⃣ 停止并删除容器..."
docker-compose down

# 删除数据卷（仅完全卸载）
if [ "$FULL_UNINSTALL" = true ]; then
    echo "2️⃣ 删除数据卷..."
    docker volume rm recon_postgres_data recon_scan_data 2>/dev/null || true
    
    echo "3️⃣ 删除镜像..."
    docker rmi recon-backend recon-frontend 2>/dev/null || true
    docker rmi postgres:15-alpine 2>/dev/null || true
    docker rmi projectdiscovery/subfinder:latest 2>/dev/null || true
    docker rmi projectdiscovery/httpx:latest 2>/dev/null || true
    docker rmi tomnomnom/assetfinder:latest 2>/dev/null || true
else
    echo "2️⃣ 保留数据卷和基础镜像"
fi

echo ""
echo "✅ 卸载完成！"

if [ "$FULL_UNINSTALL" = false ]; then
    echo ""
    echo "ℹ️  数据已保留，下次启动时会恢复"
    echo "   如需完全删除，请运行: ./uninstall.sh --full"
fi
