#!/bin/bash

echo "🔍 检查扫描工具..."
echo ""

TOOLS_PATH="/root/go/bin"
ALL_OK=true

check_tool() {
    local tool_name=$1
    local tool_path="$TOOLS_PATH/$tool_name"
    
    if [ -f "$tool_path" ]; then
        # 确保有执行权限
        chmod +x "$tool_path" 2>/dev/null || true
        
        if [ -x "$tool_path" ]; then
            echo "✅ $tool_name: $tool_path"
            # 尝试获取版本
            $tool_path -version 2>&1 | head -n 1 || $tool_path -h 2>&1 | head -n 1 || echo "   (已安装)"
            return 0
        else
            echo "❌ $tool_name: 无执行权限"
            ALL_OK=false
            return 1
        fi
    else
        echo "❌ $tool_name: 未找到"
        ALL_OK=false
        return 1
    fi
}

echo "检查路径: $TOOLS_PATH"
echo ""

check_tool "subfinder"
echo ""

check_tool "samoscout"
echo ""

check_tool "ksubdomain"
echo ""

check_tool "httpx"
echo ""

if [ "$ALL_OK" = true ]; then
    echo "✅ 所有工具已安装并可执行！"
    echo ""
    echo "下一步："
    echo "  ./deploy.sh"
    exit 0
else
    echo "❌ 部分工具未安装或无执行权限"
    echo ""
    echo "安装命令："
    echo "  go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
    echo "  go install -v github.com/samogod/samoscout/cmd/samoscout@latest"
    echo "  go install -v github.com/boy-hack/ksubdomain/cmd/ksubdomain@latest"
    echo "  go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest"
    echo ""
    echo "添加执行权限："
    echo "  chmod +x ~/go/bin/*"
    exit 1
fi
