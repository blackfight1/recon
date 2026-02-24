#!/bin/bash

echo "🧪 测试部署..."
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数
PASSED=0
FAILED=0

# 测试函数
test_command() {
    local name=$1
    local command=$2
    
    echo -n "测试 $name... "
    if eval $command > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 通过${NC}"
        PASSED=$((PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ 失败${NC}"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# 1. 检查 Docker
echo "=== 环境检查 ==="
test_command "Docker" "docker --version"
test_command "Docker Compose" "docker-compose --version"
echo ""

# 2. 检查文件结构
echo "=== 文件结构检查 ==="
test_command "docker-compose.yml" "test -f docker-compose.yml"
test_command "backend/main.go" "test -f backend/main.go"
test_command "backend/go.mod" "test -f backend/go.mod"
test_command "backend/Dockerfile" "test -f backend/Dockerfile"
test_command "frontend/package.json" "test -f frontend/package.json"
test_command "frontend/Dockerfile" "test -f frontend/Dockerfile"
echo ""

# 3. 检查配置
echo "=== 配置检查 ==="
if [ ! -f "backend/config.yaml" ]; then
    echo -e "${YELLOW}⚠ backend/config.yaml 不存在，从模板创建...${NC}"
    cp backend/config.example.yaml backend/config.yaml
    echo -e "${GREEN}✓ 配置文件已创建${NC}"
else
    echo -e "${GREEN}✓ backend/config.yaml 已存在${NC}"
fi
echo ""

# 4. 启动服务
echo "=== 启动服务 ==="
echo "正在启动 Docker 容器..."
docker-compose up -d

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 服务启动成功${NC}"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗ 服务启动失败${NC}"
    FAILED=$((FAILED + 1))
    exit 1
fi
echo ""

# 5. 等待服务就绪
echo "=== 等待服务就绪 ==="
echo "等待 15 秒..."
sleep 15
echo ""

# 6. 检查容器状态
echo "=== 容器状态检查 ==="
test_command "PostgreSQL 容器" "docker-compose ps | grep recon-db | grep Up"
test_command "后端容器" "docker-compose ps | grep recon-api | grep Up"
test_command "前端容器" "docker-compose ps | grep recon-web | grep Up"
echo ""

# 7. 测试 API
echo "=== API 测试 ==="
test_command "健康检查" "curl -s http://localhost:8000/health | grep ok"
test_command "统计接口" "curl -s http://localhost:8000/api/stats"
test_command "目标列表" "curl -s http://localhost:8000/api/targets"
echo ""

# 8. 测试前端
echo "=== 前端测试 ==="
test_command "前端访问" "curl -s http://localhost:8080 | grep -i html"
echo ""

# 9. 测试数据库连接
echo "=== 数据库测试 ==="
test_command "数据库连接" "docker-compose exec -T postgres psql -U recon -d recon -c 'SELECT 1'"
test_command "表结构检查" "docker-compose exec -T postgres psql -U recon -d recon -c '\dt' | grep targets"
echo ""

# 10. 显示日志（最后 20 行）
echo "=== 服务日志 ==="
echo "后端日志（最后 10 行）："
docker-compose logs --tail=10 backend
echo ""

# 总结
echo "=== 测试总结 ==="
echo -e "通过: ${GREEN}$PASSED${NC}"
echo -e "失败: ${RED}$FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！系统运行正常。${NC}"
    echo ""
    echo "访问地址："
    echo "  - Web 界面: http://localhost:8080"
    echo "  - API 接口: http://localhost:8000/api"
    echo ""
    echo "下一步："
    echo "  1. 在 Web 界面添加监控目标"
    echo "  2. 触发扫描测试"
    echo "  3. 查看扫描结果"
    echo ""
    echo "查看日志: docker-compose logs -f"
    echo "停止服务: docker-compose down"
    exit 0
else
    echo -e "${RED}❌ 部分测试失败，请检查错误信息。${NC}"
    echo ""
    echo "故障排查："
    echo "  1. 查看完整日志: docker-compose logs"
    echo "  2. 检查容器状态: docker-compose ps"
    echo "  3. 重启服务: docker-compose restart"
    echo "  4. 查看文档: INSTALL.md"
    exit 1
fi
