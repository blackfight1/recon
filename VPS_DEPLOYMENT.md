# VPS 部署指南

## 问题修复

### 1. 构建错误修复

已修复 `go.sum` 文件缺失的问题。现在可以正常构建。

### 2. VPS 访问配置

在 VPS 上部署后无法访问的常见原因和解决方法：

## 快速修复步骤

### 步骤 1: 清理并重新构建

```bash
# 停止并清理
docker-compose down -v

# 重新构建并启动
docker-compose up -d --build
```

### 步骤 2: 检查服务状态

```bash
# 查看容器状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 步骤 3: 开放防火墙端口

#### Ubuntu/Debian (使用 ufw)

```bash
# 检查防火墙状态
sudo ufw status

# 开放 8080 端口
sudo ufw allow 8080/tcp

# 开放 8000 端口（API）
sudo ufw allow 8000/tcp

# 重新加载防火墙
sudo ufw reload
```

#### CentOS/RHEL (使用 firewalld)

```bash
# 检查防火墙状态
sudo firewall-cmd --state

# 开放端口
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --permanent --add-port=8000/tcp

# 重新加载
sudo firewall-cmd --reload
```

#### 使用 iptables

```bash
# 开放端口
sudo iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8000 -j ACCEPT

# 保存规则
sudo iptables-save > /etc/iptables/rules.v4
```

### 步骤 4: 检查云服务商安全组

如果使用阿里云、腾讯云、AWS 等，需要在控制台配置安全组：

1. 登录云服务商控制台
2. 找到你的 VPS 实例
3. 进入"安全组"或"防火墙"设置
4. 添加入站规则：
   - 端口：8080
   - 协议：TCP
   - 来源：0.0.0.0/0（或你的 IP）
5. 添加入站规则：
   - 端口：8000
   - 协议：TCP
   - 来源：0.0.0.0/0（或你的 IP）

### 步骤 5: 验证端口监听

```bash
# 检查端口是否在监听
netstat -tlnp | grep 8080
netstat -tlnp | grep 8000

# 或使用 ss
ss -tlnp | grep 8080
ss -tlnp | grep 8000
```

### 步骤 6: 测试本地访问

```bash
# 在 VPS 上测试
curl http://localhost:8080
curl http://localhost:8000/health
```

如果本地可以访问但外网不行，说明是防火墙问题。

### 步骤 7: 测试外网访问

从你的本地电脑测试：

```bash
# 替换 YOUR_VPS_IP 为你的 VPS IP
curl http://YOUR_VPS_IP:8080
curl http://YOUR_VPS_IP:8000/health
```

## 完整部署流程

### 1. 连接到 VPS

```bash
ssh root@YOUR_VPS_IP
```

### 2. 安装 Docker（如果未安装）

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh

# 启动 Docker
systemctl start docker
systemctl enable docker

# 安装 Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
```

### 3. 克隆或上传项目

```bash
# 创建目录
mkdir -p ~/recon-platform
cd ~/recon-platform

# 上传项目文件（从本地）
# 或使用 git clone
```

### 4. 配置项目

```bash
# 复制配置文件
cp backend/config.example.yaml backend/config.yaml

# 编辑配置（可选）
nano backend/config.yaml
```

### 5. 启动服务

```bash
# 构建并启动
docker-compose up -d --build

# 查看日志
docker-compose logs -f
```

### 6. 配置防火墙

```bash
# Ubuntu/Debian
sudo ufw allow 8080/tcp
sudo ufw allow 8000/tcp
```

### 7. 访问系统

在浏览器中访问：
```
http://YOUR_VPS_IP:8080
```

## 故障排查

### 问题 1: 构建失败

**错误**: `go.sum: not found`

**解决**:
```bash
# 已修复，重新构建
docker-compose down
docker-compose up -d --build
```

### 问题 2: 容器无法启动

**检查日志**:
```bash
docker-compose logs backend
docker-compose logs frontend
docker-compose logs postgres
```

**常见原因**:
- 端口被占用
- 内存不足
- 磁盘空间不足

**解决**:
```bash
# 检查端口占用
netstat -tlnp | grep 8080

# 检查磁盘空间
df -h

# 检查内存
free -h
```

### 问题 3: 无法外网访问

**检查顺序**:

1. 容器是否运行
```bash
docker-compose ps
```

2. 本地是否可访问
```bash
curl http://localhost:8080
```

3. 防火墙是否开放
```bash
sudo ufw status
```

4. 云服务商安全组是否配置

5. 端口是否监听在 0.0.0.0
```bash
netstat -tlnp | grep 8080
```

### 问题 4: 后端 API 无法连接数据库

**检查数据库**:
```bash
docker-compose logs postgres
docker-compose exec postgres psql -U recon -d recon -c "SELECT 1"
```

**重启数据库**:
```bash
docker-compose restart postgres
docker-compose restart backend
```

## 生产环境建议

### 1. 使用 Nginx 反向代理

创建 `nginx.conf`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. 配置 HTTPS

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com
```

### 3. 修改数据库密码

编辑 `docker-compose.yml` 和 `backend/config.yaml`，修改数据库密码。

### 4. 定期备份

```bash
# 备份数据库
docker-compose exec postgres pg_dump -U recon recon > backup.sql

# 恢复数据库
docker-compose exec -T postgres psql -U recon recon < backup.sql
```

### 5. 监控资源

```bash
# 查看容器资源使用
docker stats

# 设置资源限制（在 docker-compose.yml 中）
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
```

## 快速命令参考

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f

# 查看状态
docker-compose ps

# 重新构建
docker-compose up -d --build

# 进入容器
docker-compose exec backend sh

# 查看资源使用
docker stats

# 清理所有数据
docker-compose down -v
```

## 访问地址

- Web 界面: http://YOUR_VPS_IP:8080
- API 接口: http://YOUR_VPS_IP:8000/api
- 健康检查: http://YOUR_VPS_IP:8000/health

## 安全建议

1. 修改默认数据库密码
2. 配置防火墙只允许必要的 IP 访问
3. 使用 HTTPS
4. 定期更新 Docker 镜像
5. 定期备份数据
6. 监控日志和资源使用

---

如果还有问题，请查看日志：`docker-compose logs -f`
