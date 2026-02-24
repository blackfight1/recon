# 开发环境设置

## 本地开发（不使用 Docker）

### 后端开发

1. 安装 PostgreSQL
2. 创建数据库：
```sql
CREATE DATABASE recon;
CREATE USER recon WITH PASSWORD 'recon123456';
GRANT ALL PRIVILEGES ON DATABASE recon TO recon;
```

3. 启动后端：
```bash
cd backend
cp config.example.yaml config.yaml
# 编辑 config.yaml，修改数据库连接为 localhost
go mod download
go run main.go
```

### 前端开发

```bash
cd frontend
npm install
npm run dev
```

访问: http://localhost:5173

## 测试扫描工具

### 测试 Subfinder

```bash
docker run --rm projectdiscovery/subfinder:latest -d example.com
```

### 测试 Httpx

```bash
echo "example.com" | docker run --rm -i projectdiscovery/httpx:latest -json
```

### 测试 Assetfinder

```bash
docker run --rm tomnomnom/assetfinder:latest --subs-only example.com
```

## 调试技巧

### 查看后端日志
```bash
docker-compose logs -f backend
```

### 进入后端容器
```bash
docker-compose exec backend sh
```

### 查看数据库
```bash
docker-compose exec postgres psql -U recon -d recon
```

常用 SQL：
```sql
-- 查看所有目标
SELECT * FROM targets;

-- 查看子域名
SELECT * FROM subdomains ORDER BY last_seen DESC LIMIT 10;

-- 查看变更日志
SELECT * FROM change_logs ORDER BY created_at DESC LIMIT 10;

-- 查看任务
SELECT * FROM scan_tasks ORDER BY created_at DESC LIMIT 10;
```

### 手动触发扫描

```bash
curl -X POST http://localhost:8000/api/targets/1/scan
```

### 查看统计信息

```bash
curl http://localhost:8000/api/stats | jq
```
