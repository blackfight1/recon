# 部署验证指南

本文档用于验证项目是否正确部署并正常运行。

## 📋 验证清单

### 1. 文件结构验证 ✅

#### 根目录文件
- [x] docker-compose.yml
- [x] README.md
- [x] INSTALL.md
- [x] USAGE.md
- [x] API.md
- [x] QUICKSTART.md
- [x] PROJECT_STRUCTURE.md
- [x] PROJECT_SUMMARY.md
- [x] CHECKLIST.md
- [x] LICENSE
- [x] .gitignore
- [x] start.sh
- [x] start.bat
- [x] check-setup.sh
- [x] test-deployment.sh

#### 后端文件
- [x] backend/main.go
- [x] backend/go.mod
- [x] backend/config.example.yaml
- [x] backend/Dockerfile
- [x] backend/.dockerignore
- [x] backend/config/config.go
- [x] backend/controllers/target.go
- [x] backend/database/database.go
- [x] backend/models/models.go
- [x] backend/notifier/notifier.go
- [x] backend/router/router.go
- [x] backend/scanner/scanner.go
- [x] backend/scheduler/scheduler.go

#### 前端文件
- [x] frontend/package.json
- [x] frontend/vite.config.js
- [x] frontend/index.html
- [x] frontend/Dockerfile
- [x] frontend/nginx.conf
- [x] frontend/.dockerignore
- [x] frontend/src/main.js
- [x] frontend/src/App.vue
- [x] frontend/src/api/index.js
- [x] frontend/src/router/index.js
- [x] frontend/src/views/Dashboard.vue
- [x] frontend/src/views/Targets.vue
- [x] frontend/src/views/Subdomains.vue
- [x] frontend/src/views/Changes.vue
- [x] frontend/src/views/Tasks.vue

### 2. 环境验证

#### 检查 Docker
```bash
docker --version
# 预期输出: Docker version 20.x.x 或更高
```

#### 检查 Docker Compose
```bash
docker-compose --version
# 预期输出: docker-compose version 1.29.x 或更高
```

### 3. 配置验证

#### 创建配置文件
```bash
cp backend/config.example.yaml backend/config.yaml
```

#### 验证配置文件
```bash
cat backend/config.yaml
# 确认配置文件存在且格式正确
```

### 4. 启动验证

#### 启动服务
```bash
docker-compose up -d
```

#### 检查容器状态
```bash
docker-compose ps
```

预期输出：
```
Name                 Command               State           Ports
------------------------------------------------------------------------
recon-api      ./main                       Up      0.0.0.0:8000->8000/tcp
recon-db       docker-entrypoint.sh postgres Up      0.0.0.0:5432->5432/tcp
recon-web      nginx -g daemon off;         Up      0.0.0.0:8080->80/tcp
```

所有服务的 State 应该是 "Up"。

### 5. 服务验证

#### 验证后端 API
```bash
curl http://localhost:8000/health
```

预期输出：
```json
{"status":"ok"}
```

#### 验证统计接口
```bash
curl http://localhost:8000/api/stats
```

预期输出：
```json
{
  "data": {
    "total_targets": 0,
    "enabled_targets": 0,
    "total_subdomains": 0,
    "alive_subdomains": 0,
    "new_subdomains": 0,
    "recent_changes": 0
  }
}
```

#### 验证目标列表接口
```bash
curl http://localhost:8000/api/targets
```

预期输出：
```json
{"data":[]}
```

#### 验证前端访问
```bash
curl -I http://localhost:8080
```

预期输出应包含：
```
HTTP/1.1 200 OK
```

### 6. 数据库验证

#### 连接数据库
```bash
docker-compose exec postgres psql -U recon -d recon
```

#### 检查表结构
```sql
\dt
```

预期输出应包含：
- targets
- subdomains
- scan_tasks
- change_logs

#### 退出数据库
```sql
\q
```

### 7. 日志验证

#### 查看后端日志
```bash
docker-compose logs backend
```

应该看到类似输出：
```
Database connected successfully
Scheduler started successfully
Server starting on port 8000...
```

#### 查看前端日志
```bash
docker-compose logs frontend
```

应该看到 Nginx 启动日志。

#### 查看数据库日志
```bash
docker-compose logs postgres
```

应该看到 PostgreSQL 启动成功的日志。

### 8. 功能验证

#### 8.1 添加目标

使用 Web 界面：
1. 访问 http://localhost:8080
2. 点击"目标管理"
3. 点击"添加目标"
4. 输入域名：example.com
5. 输入描述：测试目标
6. 点击"确定"

或使用 API：
```bash
curl -X POST http://localhost:8000/api/targets \
  -H "Content-Type: application/json" \
  -d '{"domain":"example.com","description":"测试目标"}'
```

预期输出：
```json
{
  "data": {
    "id": 1,
    "domain": "example.com",
    "description": "测试目标",
    "enabled": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 8.2 触发扫描

使用 Web 界面：
1. 在目标列表中点击"扫描"按钮

或使用 API：
```bash
curl -X POST http://localhost:8000/api/targets/1/scan
```

预期输出：
```json
{"message":"Scan triggered successfully"}
```

#### 8.3 查看任务状态

使用 Web 界面：
1. 点击"任务列表"
2. 查看扫描任务状态

或使用 API：
```bash
curl http://localhost:8000/api/tasks
```

#### 8.4 查看扫描结果

等待扫描完成后（通常 5-10 分钟）：

使用 Web 界面：
1. 点击"子域名"查看发现的子域名
2. 点击"变更中心"查看变更记录

或使用 API：
```bash
# 查看子域名
curl http://localhost:8000/api/subdomains

# 查看变更
curl http://localhost:8000/api/changes
```

### 9. 通知验证（可选）

如果配置了通知：

#### 9.1 配置企业微信
编辑 `backend/config.yaml`：
```yaml
notification:
  wecom:
    enabled: true
    webhook: "你的Webhook URL"
```

#### 9.2 配置钉钉
编辑 `backend/config.yaml`：
```yaml
notification:
  dingtalk:
    enabled: true
    webhook: "你的Webhook URL"
    secret: "你的密钥"
```

#### 9.3 重启后端
```bash
docker-compose restart backend
```

#### 9.4 触发扫描
再次触发扫描，如果有变更应该会收到通知。

### 10. 性能验证

#### 检查资源使用
```bash
docker stats
```

正常情况下：
- recon-api: ~100MB 内存
- recon-web: ~50MB 内存
- recon-db: ~200MB 内存

#### 检查磁盘使用
```bash
docker system df
```

### 11. 清理测试

#### 停止服务
```bash
docker-compose down
```

#### 删除数据（可选）
```bash
docker-compose down -v
```

## ✅ 验证结果

如果以上所有步骤都成功，说明项目部署正确且运行正常！

## 🐛 常见问题

### 问题 1: 容器无法启动

**症状**: `docker-compose ps` 显示容器状态为 Exit

**解决方法**:
```bash
# 查看详细日志
docker-compose logs

# 重新构建
docker-compose up -d --build
```

### 问题 2: API 返回 502

**症状**: `curl http://localhost:8000/health` 返回 502

**解决方法**:
```bash
# 检查后端日志
docker-compose logs backend

# 重启后端
docker-compose restart backend
```

### 问题 3: 前端无法访问

**症状**: 浏览器无法打开 http://localhost:8080

**解决方法**:
```bash
# 检查前端日志
docker-compose logs frontend

# 检查端口占用
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Linux/Mac

# 重启前端
docker-compose restart frontend
```

### 问题 4: 扫描无结果

**症状**: 扫描完成但没有发现子域名

**可能原因**:
1. 目标域名确实没有子域名
2. 网络连接问题
3. 扫描工具镜像未下载

**解决方法**:
```bash
# 查看后端日志
docker-compose logs backend

# 手动拉取镜像
docker pull projectdiscovery/subfinder:latest
docker pull projectdiscovery/httpx:latest
docker pull tomnomnom/assetfinder:latest

# 测试网络
curl https://crt.sh
```

### 问题 5: 数据库连接失败

**症状**: 后端日志显示数据库连接错误

**解决方法**:
```bash
# 检查数据库状态
docker-compose ps postgres

# 检查数据库日志
docker-compose logs postgres

# 重启数据库
docker-compose restart postgres

# 等待数据库就绪后重启后端
docker-compose restart backend
```

## 📞 获取帮助

如果遇到问题：

1. 查看日志：`docker-compose logs -f`
2. 检查文档：README.md, INSTALL.md, USAGE.md
3. 运行测试脚本：`./test-deployment.sh`
4. 提交 Issue：GitHub Issues

## 🎉 验证完成

恭喜！如果所有验证都通过，你的自动化侦查平台已经成功部署并可以使用了！

下一步：
1. 添加真实的监控目标
2. 配置通知渠道
3. 等待自动扫描或手动触发
4. 查看资产变更

祝你挖洞愉快！🚀
