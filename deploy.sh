#!/bin/bash

set -e

echo "🚀 部署自动化侦查平台..."
echo ""

# 1. 检查工具
echo "1️⃣ 检查扫描工具..."
TOOLS_PATH="/root/go/bin"
ALL_OK=true

check_tool() {
    local tool_name=$1
    local tool_path="$TOOLS_PATH/$tool_name"
    
    if [ -f "$tool_path" ]; then
        chmod +x "$tool_path" 2>/dev/null || true
        echo "✅ $tool_name"
        return 0
    else
        echo "❌ $tool_name: 未找到"
        ALL_OK=false
        return 1
    fi
}

check_tool "subfinder"
check_tool "samoscout"
check_tool "ksubdomain"
check_tool "httpx"

if [ "$ALL_OK" = false ]; then
    echo ""
    echo "❌ 部分工具未安装，请先安装："
    echo ""
    echo "go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
    echo "go install -v github.com/samogod/samoscout/cmd/samoscout@latest"
    echo "go install -v github.com/boy-hack/ksubdomain/cmd/ksubdomain@latest"
    echo "go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest"
    echo "chmod +x ~/go/bin/*"
    echo ""
    exit 1
fi

echo ""
echo "2️⃣ 停止旧服务..."
docker-compose down

echo ""
echo "3️⃣ 重新构建镜像..."
echo "   构建后端..."
docker-compose build --no-cache backend
echo "   构建前端..."
docker-compose build --no-cache frontend

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
echo "7️⃣ 验证工具挂载..."
docker-compose exec -T backend sh -c "ls -lh /usr/local/bin/subfinder /usr/local/bin/httpx" 2>/dev/null || echo "⚠️  工具挂载检查跳过"

echo ""
echo "8️⃣ 后端日志（最后20行）："
docker-compose logs --tail=20 backend

echo ""
echo "✅ 部署完成！"
echo ""
echo "📍 访问地址："
IP=$(hostname -I | awk '{print $1}' 2>/dev/null || echo "你的IP")
echo "   前端: http://$IP:8080"
echo "   后端: http://$IP:8000"
echo ""
echo "📝 查看日志："
echo "   docker-compose logs -f backend"
echo ""
echo "🧪 测试快速扫描："
echo "   1. 访问前端页面（刷新浏览器 Ctrl+F5）"
echo "   2. 点击'快速扫描'"
echo "   3. 输入域名测试"
