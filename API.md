# API 文档

Base URL: `http://localhost:8000/api`

## 目标管理

### 获取所有目标
```
GET /targets
```

响应：
```json
{
  "data": [
    {
      "id": 1,
      "domain": "example.com",
      "description": "测试目标",
      "enabled": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 创建目标
```
POST /targets
Content-Type: application/json

{
  "domain": "example.com",
  "description": "测试目标"
}
```

### 更新目标
```
PUT /targets/:id
Content-Type: application/json

{
  "domain": "example.com",
  "description": "更新描述",
  "enabled": false
}
```

### 删除目标
```
DELETE /targets/:id
```

### 触发扫描
```
POST /targets/:id/scan
```

## 子域名

### 获取所有子域名
```
GET /subdomains?page=1&page_size=50&status=alive
```

参数：
- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 50）
- `status`: 状态筛选（new/alive/dead）

响应：
```json
{
  "data": [
    {
      "id": 1,
      "target_id": 1,
      "subdomain": "www.example.com",
      "status": "alive",
      "status_code": 200,
      "title": "Example Domain",
      "server": "nginx",
      "ip": "93.184.216.34",
      "first_seen": "2024-01-01T00:00:00Z",
      "last_seen": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "page": 1
}
```

### 获取指定目标的子域名
```
GET /subdomains/target/:target_id
```

## 变更日志

### 获取所有变更
```
GET /changes?page=1&page_size=50
```

响应：
```json
{
  "data": [
    {
      "id": 1,
      "target_id": 1,
      "change_type": "subdomain_new",
      "content": "api.example.com (状态码: 200, 标题: API)",
      "notified": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 50,
  "page": 1
}
```

变更类型：
- `subdomain_new`: 新增子域名
- `subdomain_alive`: 子域名恢复存活
- `subdomain_dead`: 子域名失效

### 获取指定目标的变更
```
GET /changes/target/:target_id
```

## 任务管理

### 获取所有任务
```
GET /tasks?page=1&page_size=20
```

响应：
```json
{
  "data": [
    {
      "id": 1,
      "target_id": 1,
      "task_type": "subdomain",
      "status": "completed",
      "result": "Found 10 alive subdomains, 2 changes",
      "error_msg": "",
      "started_at": "2024-01-01T00:00:00Z",
      "completed_at": "2024-01-01T00:05:00Z",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 20,
  "page": 1
}
```

任务类型：
- `subdomain`: 子域名扫描
- `port`: 端口扫描
- `screenshot`: 截图
- `vuln`: 漏洞扫描

任务状态：
- `pending`: 等待中
- `running`: 运行中
- `completed`: 已完成
- `failed`: 失败

### 获取单个任务
```
GET /tasks/:id
```

## 统计信息

### 获取统计数据
```
GET /stats
```

响应：
```json
{
  "data": {
    "total_targets": 5,
    "enabled_targets": 4,
    "total_subdomains": 150,
    "alive_subdomains": 120,
    "new_subdomains": 10,
    "recent_changes": 25
  }
}
```

## 健康检查

```
GET /health
```

响应：
```json
{
  "status": "ok"
}
```

## 错误响应

所有错误响应格式：
```json
{
  "error": "错误描述"
}
```

HTTP 状态码：
- `200`: 成功
- `201`: 创建成功
- `400`: 请求参数错误
- `404`: 资源不存在
- `500`: 服务器内部错误
