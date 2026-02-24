# 部署清单

## ✅ 部署步骤

### 1. 准备环境

```bash
# 检查 Docker
docker --version
docker-compose --version
```

### 2. 部署系统

```bash
# 给脚本执行权限
chmod +x rebuild.sh start.sh uninstall.sh

# 执行部署
./rebuild.sh
```

### 3. 验证部署

```bash
# 检查容器状态（应该都是 Up）
docker-compose ps

# 检查后端（应该返回 {"status":"ok"}）
curl http://localhost:8000/health

# 检查前端（应该返回 200）
curl -I http://localhost:8080
```

### 4. 访问系统

```bash
# 获取服务器 IP
curl ifconfig.me

# 在浏览器访问
http://你的IP:8080
```

## 🐛 故障排查

### 问题 1: 构建失败

```bash
# 查看详细日志
docker-compose logs backend

# 完全重建
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```

### 问题 2: 无法访问

```bash
# 检查端口监听
netstat -tlnp | grep 8080

# 检查容器状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 问题 3: 数据库连接失败

```bash
# 重启数据库
docker-compose restart postgres

# 等待 10 秒后重启后端
sleep 10
docker-compose restart backend
```

## 📝 配置通知（可选）

编辑 `backend/config.yaml`:

```yaml
notification:
  wecom:
    enabled: true
    webhook: "你的企业微信Webhook"
  dingtalk:
    enabled: true
    webhook: "你的钉钉Webhook"
    secret: "你的密钥"
```

重启后端：
```bash
docker-compose restart backend
```

## 🗑️ 卸载

```bash
./uninstall.sh
```

## 📞 获取帮助

- 查看日志: `docker-compose logs -f`
- 查看状态: `docker-compose ps`
- 详细文档: `cat VPS_DEPLOYMENT.md`
