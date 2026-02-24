@echo off
chcp 65001 >nul
echo 🚀 启动自动化侦查平台...
echo.

REM 检查 Docker 是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker 未安装，请先安装 Docker Desktop
    pause
    exit /b 1
)

REM 检查 Docker Compose 是否安装
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker Compose 未安装
    pause
    exit /b 1
)

REM 检查配置文件
if not exist "backend\config.yaml" (
    echo 📝 创建配置文件...
    copy backend\config.example.yaml backend\config.yaml
    echo ⚠️  请编辑 backend\config.yaml 配置通知 Webhook
)

REM 创建数据目录
if not exist "data" mkdir data

REM 启动服务
echo 🐳 启动 Docker 容器...
docker-compose up -d

REM 等待服务启动
echo ⏳ 等待服务启动...
timeout /t 10 /nobreak >nul

REM 检查服务状态
echo 📊 检查服务状态...
docker-compose ps

echo.
echo ✨ 部署完成！
echo.
echo 📱 访问地址: http://localhost:8080
echo 📚 API 文档: http://localhost:8000/api
echo.
echo 📋 常用命令:
echo   查看日志: docker-compose logs -f
echo   停止服务: docker-compose down
echo   重启服务: docker-compose restart
echo.
pause
