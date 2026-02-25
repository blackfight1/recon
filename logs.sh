#!/bin/bash

# 查看后端日志的快捷脚本

echo "📋 查看后端日志..."
echo "按 Ctrl+C 退出"
echo ""

docker-compose logs -f backend
