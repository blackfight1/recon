#!/bin/bash

echo "🔍 检查扫描工具是否安装..."
echo ""

TOOLS_PATH="/root/go/bin"
ALL_OK=true

# 检查工具函数
check_tool() {
    local tool_name=$1
    local tool_path="$TOOLS_PATH/$tool_name"
    
    if [ -f "$tool_path" ] && [ -x "$tool_path" ]; then
        echo "✅ $tool_name: $tool_path"
        # 尝试获取版本
        $tool_path -version 2>&1 | head -n 1 || echo "   (已安装)"
        return 0
    else
        echo "❌ $tool_name: 未找到或不可执行"
        ALL_OK=false
        return 1
    fi
}

echo "检查路径: $TOOLS_PATH"
echo ""

# 检查所有工具
check_tool "subfinder"
echo ""

check_tool "samoscout"
echo ""

check_tool "ksubdomain"
echo ""

check_tool "httpx"
echo ""

# 总结
if [ "$ALL_OK" = true ]; then
    echo "✅ 所有工具已安装并可执行！"
    echo ""
    echo "下一步："
    echo "  1. 修改 docker-compose.yml 挂载工具路径"
    echo "  2. 修改 scanner.go 使用本地工具"
    echo "  3. 运行 ./rebuild.sh"
    exit 0
else
    echo ""
    echo "❌ 部分工具未安装"
    echo ""
    echo "安装命令："
    echo "  go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
    echo "  go install -v github.com/samogod/samoscout/cmd/samoscout@latest"
    echo "  go install -v github.com/boy-hack/ksubdomain/cmd/ksubdomain@latest"
    echo "  go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest"
    exit 1
fi
