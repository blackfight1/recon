# 快速开始指南

## 5 分钟快速部署

### 第一步：环境检查

确保已安装 Docker 和 Docker Compose：

```bash
docker --version
docker-compose --version
```

如果未安装，请参考 [INSTALL.md](INSTALL.md)

### 第二步：启动服务

**Linux/Mac:**
```bash
chmod +x start.sh
./start.sh
```

**Windows:**
双击运行 `start.bat`

或手动执行：
```bash
docker-compose up -d
```

### 第三步：访问系统

打开浏览器访问: http://localhost:8080

### 第四步：添加第一个目标

1. 点击顶部菜单"目标管理"
2. 点击"添加目标"按钮
3. 输入域名（例如：example.com）
4. 点击"确定"
5. 点击"扫描"按钮触发首次扫描

### 第五步：查看结果

- 在"任务列表"查看扫描进度
- 扫描完成后在"子域名"查看发现的资产
- 在"变更中心"查看资产变更

## 配置通知（可选）

### 企业微信通知

1. 在企业微信群中添加机器人，获取 Webhook URL
2. 编辑 `backend/config.yaml`：

```yaml
notification:
  wecom:
    enabled: true
    webhook: "你的Webhook URL"
```

3. 重启后端服务：
```bash
docker-compose restart backend
```

### 钉钉通知

1. 在钉钉群中添加自定义机器人，获取 Webhook URL
2. 如果启用了加签，记录密钥
3. 编辑 `backend/config.yaml`：

```yaml
notification:
  dingtalk:
    enabled: true
    webhook: "你的Webhook URL"
    secret: "你的密钥"  # 如果启用了加签
```

4. 重启后端服务：
```bash
docker-compose restart backend
```

## 验证安装

### 检查服务状态

```bash
docker-compose ps
```

应该看到三个服务都在运行：
- recon-db (PostgreSQL)
- recon-api (Go 后端)
- recon-web (Vue 前端)

### 检查后端 API

```bash
curl http://localhost:8000/health
```

应该返回: `{"status":"ok"}`

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 只查看后端日志
docker-compose logs -f backend

# 只查看前端日志
docker-compose logs -f frontend
```

## 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 进入后端容器
docker-compose exec backend sh

# 进入数据库
docker-compose exec postgres psql -U recon -d recon

# 重新构建并启动
docker-compose up -d --build

# 停止并删除所有数据
docker-compose down -v
```

## 测试扫描

### 使用 API 测试

1. 添加目标：
```bash
curl -X POST http://localhost:8000/api/targets \
  -H "Content-Type: application/json" \
  -d '{"domain":"example.com","description":"测试目标"}'
```

2. 触发扫描：
```bash
curl -X POST http://localhost:8000/api/targets/1/scan
```

3. 查看任务状态：
```bash
curl http://localhost:8000/api/tasks | jq
```

4. 查看子域名：
```bash
curl http://localhost:8000/api/subdomains | jq
```

5. 查看变更：
```bash
curl http://localhost:8000/api/changes | jq
```

## 故障排查

### 问题 1: 端口被占用

错误信息：`Bind for 0.0.0.0:8080 failed: port is already allocated`

解决方法：
```bash
# 查看占用端口的进程
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Linux/Mac

# 修改 docker-compose.yml 中的端口映射
# 例如将 8080:80 改为 8081:80
```

### 问题 2: Docker 权限不足

错误信息：`permission denied while trying to connect to the Docker daemon socket`

解决方法：
```bash
# Linux
sudo usermod -aG docker $USER
# 重新登录或执行
newgrp docker

# 或使用 sudo
sudo docker-compose up -d
```

### 问题 3: 扫描无结果

可能原因：
1. 目标域名不存在或无子域名
2. 网络连接问题
3. 扫描工具镜像未下载

解决方法：
```bash
# 查看后端日志
docker-compose logs backend

# 手动拉取扫描工具镜像
docker pull projectdiscovery/subfinder:latest
docker pull projectdiscovery/httpx:latest
docker pull tomnomnom/assetfinder:latest

# 测试网络连接
curl https://crt.sh
```

### 问题 4: 前端无法访问后端

可能原因：
1. 后端服务未启动
2. 网络配置问题

解决方法：
```bash
# 检查后端状态
docker-compose ps backend

# 检查后端日志
docker-compose logs backend

# 测试后端 API
curl http://localhost:8000/health

# 重启服务
docker-compose restart
```

## 下一步

- 阅读 [USAGE.md](USAGE.md) 了解详细使用方法
- 查看 [API.md](API.md) 了解 API 接口
- 参考 [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) 了解项目结构

## 获取帮助

- 查看文档：[README.md](README.md)
- 提交问题：GitHub Issues
- 查看日志：`docker-compose logs -f`

## 卸载

如果需要完全卸载：

```bash
# 停止并删除容器
docker-compose down

# 删除数据卷（会删除所有数据）
docker-compose down -v

# 删除镜像
docker rmi recon-platform_backend recon-platform_frontend

# 删除扫描工具镜像（可选）
docker rmi projectdiscovery/subfinder:latest
docker rmi projectdiscovery/httpx:latest
docker rmi tomnomnom/assetfinder:latest
```

---

🎉 恭喜！你已经成功部署了自动化侦查平台！
