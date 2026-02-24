#!/bin/bash

echo "🔍 检查项目设置..."
echo ""

# 检查必要文件
echo "📁 检查文件结构..."
files=(
    "docker-compose.yml"
    "backend/main.go"
    "backend/go.mod"
    "backend/config.example.yaml"
    "backend/Dockerfile"
    "frontend/package.json"
    "frontend/Dockerfile"
    "frontend/src/main.js"
)

missing_files=0
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✅ $file"
    else
        echo "  ❌ $file (缺失)"
        missing_files=$((missing_files + 1))
    fi
done

echo ""

# 检查目录结构
echo "📂 检查目录结构..."
dirs=(
    "backend/config"
    "backend/controllers"
    "backend/database"
    "backend/models"
    "backend/notifier"
    "backend/router"
    "backend/scanner"
    "backend/scheduler"
    "frontend/src/api"
    "frontend/src/router"
    "frontend/src/views"
)

missing_dirs=0
for dir in "${dirs[@]}"; do
    if [ -d "$dir" ]; then
        echo "  ✅ $dir"
    else
        echo "  ❌ $dir (缺失)"
        missing_dirs=$((missing_dirs + 1))
    fi
done

echo ""

# 检查 Docker
echo "🐳 检查 Docker..."
if command -v docker &> /dev/null; then
    echo "  ✅ Docker 已安装: $(docker --version)"
else
    echo "  ❌ Docker 未安装"
fi

if command -v docker-compose &> /dev/null; then
    echo "  ✅ Docker Compose 已安装: $(docker-compose --version)"
else
    echo "  ❌ Docker Compose 未安装"
fi

echo ""

# 检查配置文件
echo "⚙️  检查配置..."
if [ -f "backend/config.yaml" ]; then
    echo "  ✅ backend/config.yaml 已存在"
else
    echo "  ⚠️  backend/config.yaml 不存在（首次运行会自动创建）"
fi

echo ""

# 总结
echo "📊 检查总结:"
echo "  缺失文件: $missing_files"
echo "  缺失目录: $missing_dirs"

if [ $missing_files -eq 0 ] && [ $missing_dirs -eq 0 ]; then
    echo ""
    echo "✨ 项目结构完整！可以开始部署了。"
    echo ""
    echo "下一步："
    echo "  1. 编辑 backend/config.yaml 配置通知"
    echo "  2. 运行 ./start.sh 启动服务"
    echo "  3. 访问 http://localhost:8080"
else
    echo ""
    echo "⚠️  项目结构不完整，请检查缺失的文件和目录。"
fi

echo ""
