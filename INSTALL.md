# 安装部署指南

## 环境准备

### 1. 安装 Docker 和 Docker Compose

Ubuntu/Debian:
```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh

# 启动 Docker 服务
sudo systemctl start docker
sudo systemctl enable docker

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker-compose --version
```

### 2. 拉取扫描工具镜像（可选，首次运行会自动拉取）

```bash
docker pull projectdiscovery/subfinder:latest
docker pull projectdiscovery/httpx:latest
docker pull tomnomnom/assetfinder:latest
```

## 部署步骤

### 1. 克隆项目

```bash
git clone <your-repo-url>
cd recon-platform
```

### 2. 配置后端

```bash
cd backend
cp config.example.yaml config.yaml
```

编辑 `config.yaml`，配置通知 Webhook：

```yaml
notification:
  # 企业微信机器人
  wecom:
    enabled: true
    webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY"
  
  # 钉钉机器人
  dingtalk:
    enabled: true
    webhook: "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN"
    secret: "YOUR_SECRET"  # 如果启用了加签
```

### 3. 启动服务

```bash
# 返回项目根目录
cd ..

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 4. 访问系统

打开浏览器访问: http://your-server-ip:8080

## 验证安装

### 1. 检查服务状态

```bash
docker-compose ps
```

应该看到三个服务都在运行：
- recon-db (PostgreSQL)
- recon-api (Go 后端)
- recon-web (Vue 前端)

### 2. 检查 API

```bash
curl http://localhost:8000/health
```

应该返回: `{"status":"ok"}`

### 3. 测试扫描

1. 在 Web 界面添加一个测试目标（如 example.com）
2. 点击"扫描"按钮
3. 在"任务列表"中查看扫描进度
4. 扫描完成后在"子域名"和"变更中心"查看结果

## 常见问题

### 1. 后端无法连接数据库

检查 PostgreSQL 是否正常启动：
```bash
docker-compose logs postgres
```

### 2. 扫描工具无法运行

确保 Docker socket 已挂载：
```bash
ls -la /var/run/docker.sock
```

如果权限不足：
```bash
sudo chmod 666 /var/run/docker.sock
```

### 3. 前端无法访问后端 API

检查网络配置：
```bash
docker network ls
docker network inspect recon-platform_recon-network
```

### 4. 扫描结果为空

1. 检查目标域名是否正确
2. 查看后端日志：`docker-compose logs backend`
3. 确保服务器可以访问外网

## 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose down
docker-compose up -d --build

# 查看日志
docker-compose logs -f
```

## 卸载

```bash
# 停止并删除所有容器
docker-compose down

# 删除数据卷（注意：会删除所有数据）
docker-compose down -v

# 删除镜像
docker rmi recon-platform_backend recon-platform_frontend
```

## 生产环境建议

1. 修改数据库密码（docker-compose.yml 和 config.yaml）
2. 配置 Nginx 反向代理并启用 HTTPS
3. 定期备份 PostgreSQL 数据
4. 监控磁盘空间（扫描数据会持续增长）
5. 配置防火墙规则，只开放必要端口

## 性能优化

1. 调整扫描并发数（config.yaml 中的 concurrency）
2. 根据服务器性能调整扫描间隔
3. 定期清理旧的扫描数据
4. 使用 SSD 存储提升数据库性能
