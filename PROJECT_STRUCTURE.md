# 项目结构说明

## 目录结构

```
recon-platform/
├── backend/                    # Go 后端
│   ├── config/                # 配置管理
│   │   └── config.go         # 配置结构和加载
│   ├── controllers/           # 控制器层
│   │   └── target.go         # API 处理函数
│   ├── database/              # 数据库层
│   │   └── database.go       # 数据库初始化和连接
│   ├── models/                # 数据模型
│   │   └── models.go         # 数据库表结构
│   ├── notifier/              # 通知模块
│   │   └── notifier.go       # 企业微信/钉钉通知
│   ├── router/                # 路由层
│   │   └── router.go         # API 路由定义
│   ├── scanner/               # 扫描引擎
│   │   └── scanner.go        # 子域名扫描和差量对比
│   ├── scheduler/             # 任务调度
│   │   └── scheduler.go      # 定时任务管理
│   ├── config.example.yaml    # 配置文件模板
│   ├── Dockerfile            # 后端 Docker 镜像
│   ├── go.mod                # Go 依赖管理
│   └── main.go               # 程序入口
│
├── frontend/                  # Vue 前端
│   ├── src/
│   │   ├── api/              # API 接口封装
│   │   │   └── index.js
│   │   ├── router/           # 路由配置
│   │   │   └── index.js
│   │   ├── views/            # 页面组件
│   │   │   ├── Dashboard.vue    # 仪表盘
│   │   │   ├── Targets.vue      # 目标管理
│   │   │   ├── Subdomains.vue   # 子域名列表
│   │   │   ├── Changes.vue      # 变更中心
│   │   │   └── Tasks.vue        # 任务列表
│   │   ├── App.vue           # 根组件
│   │   └── main.js           # 入口文件
│   ├── Dockerfile            # 前端 Docker 镜像
│   ├── nginx.conf            # Nginx 配置
│   ├── index.html            # HTML 模板
│   ├── package.json          # npm 依赖
│   └── vite.config.js        # Vite 配置
│
├── docker-compose.yml         # Docker Compose 配置
├── .gitignore                # Git 忽略文件
├── README.md                 # 项目说明
├── INSTALL.md                # 安装指南
├── USAGE.md                  # 使用指南
├── API.md                    # API 文档
├── dev-setup.md              # 开发环境设置
├── PROJECT_STRUCTURE.md      # 本文件
└── start.sh                  # 快速启动脚本
```

## 核心模块说明

### 后端模块

#### 1. Config (配置管理)
- 使用 Viper 加载 YAML 配置
- 支持服务器、数据库、扫描器、通知等配置
- 配置文件：`backend/config.yaml`

#### 2. Database (数据库层)
- 使用 GORM 作为 ORM
- PostgreSQL 数据库
- 自动迁移表结构
- 四个核心表：
  - targets: 监控目标
  - subdomains: 子域名记录
  - scan_tasks: 扫描任务
  - change_logs: 变更日志

#### 3. Models (数据模型)
- Target: 监控目标模型
- Subdomain: 子域名模型
- ScanTask: 扫描任务模型
- ChangeLog: 变更日志模型

#### 4. Scanner (扫描引擎)
核心功能模块，负责：
- 子域名收集（Subfinder + Assetfinder + cert.sh）
- 存活验证（Httpx）
- 差量对比（对比历史数据）
- 变更检测（新增/恢复/失效）

工作流程：
```
collectSubdomains() -> verifyAlive() -> compareAndUpdate()
```

#### 5. Scheduler (任务调度)
- 使用 robfig/cron 实现定时任务
- 默认每 6 小时扫描一次
- 支持启动时立即执行一次
- 并发扫描多个目标

#### 6. Notifier (通知模块)
- 支持企业微信机器人
- 支持钉钉机器人（含加签）
- 变更通知格式化
- 自动限制通知内容长度

#### 7. Router (路由层)
- 使用 Gin 框架
- RESTful API 设计
- CORS 支持
- 路由分组管理

#### 8. Controllers (控制器层)
处理所有 API 请求：
- 目标 CRUD
- 子域名查询
- 变更日志查询
- 任务管理
- 统计信息

### 前端模块

#### 1. Router (路由)
- Vue Router 4
- 五个主要页面路由
- History 模式

#### 2. API (接口封装)
- Axios 封装
- 统一请求/响应拦截
- 错误处理

#### 3. Views (页面组件)

**Dashboard (仪表盘)**
- 统计卡片展示
- 快速操作入口
- 自动刷新数据

**Targets (目标管理)**
- 目标列表展示
- 添加/编辑/删除目标
- 手动触发扫描
- 启用/禁用目标

**Subdomains (子域名)**
- 子域名列表
- 状态筛选
- 分页展示
- 详细信息展示

**Changes (变更中心)**
- 时间轴展示
- 变更类型图标
- 分页加载

**Tasks (任务列表)**
- 任务状态展示
- 执行结果查看
- 自动刷新

## 数据流

### 扫描流程数据流

```
用户添加目标
    ↓
Scheduler 定时触发
    ↓
Scanner.ScanTarget()
    ↓
collectSubdomains() → Docker 容器执行工具
    ↓
verifyAlive() → Docker 容器执行 Httpx
    ↓
compareAndUpdate() → 对比数据库
    ↓
生成 ChangeLog
    ↓
Notifier 发送通知
    ↓
更新 ScanTask 状态
```

### API 请求流程

```
前端发起请求
    ↓
Axios 拦截器
    ↓
后端 Router
    ↓
Controller 处理
    ↓
Database 查询
    ↓
返回 JSON 响应
    ↓
前端更新界面
```

## 技术栈

### 后端
- 语言：Go 1.21+
- 框架：Gin (HTTP)
- ORM：GORM
- 数据库：PostgreSQL 15
- 定时任务：robfig/cron
- 配置：Viper
- 容器：Docker SDK

### 前端
- 框架：Vue 3
- 构建工具：Vite
- UI 库：Element Plus
- 路由：Vue Router 4
- HTTP：Axios
- 图标：Element Plus Icons

### 扫描工具
- Subfinder (子域名收集)
- Assetfinder (子域名收集)
- Httpx (存活验证)
- cert.sh API (证书透明度)

### 部署
- Docker & Docker Compose
- Nginx (前端服务器)
- PostgreSQL (数据库)

## 扩展点

### 1. 添加新的扫描工具

在 `backend/scanner/scanner.go` 中添加新的收集函数：

```go
func collectFromNewTool(domain string) ([]string, error) {
    // 实现新工具的调用
}
```

然后在 `collectSubdomains()` 中调用。

### 2. 添加新的通知渠道

在 `backend/notifier/notifier.go` 中添加新的发送函数：

```go
func sendToNewChannel(message string) {
    // 实现新通知渠道
}
```

在配置文件中添加相应配置。

### 3. 添加端口扫描

1. 在 models 中添加 Port 模型
2. 在 scanner 中实现端口扫描函数
3. 在前端添加端口展示页面

### 4. 添加截图功能

1. 集成 Gowitness
2. 存储截图路径到数据库
3. 在前端展示截图墙

### 5. 添加漏洞扫描

1. 集成 Nuclei
2. 添加 Vulnerability 模型
3. 在前端展示漏洞列表

## 性能考虑

### 1. 数据库优化
- 为常用查询字段添加索引
- 使用连接池
- 定期清理历史数据

### 2. 并发控制
- 限制同时扫描的目标数量
- 使用 goroutine 池
- 控制 Docker 容器数量

### 3. 缓存策略
- 统计数据缓存
- API 响应缓存
- 前端数据缓存

### 4. 资源管理
- Docker 容器资源限制
- 扫描数据定期清理
- 日志轮转

## 安全考虑

### 1. 输入验证
- 域名格式验证
- SQL 注入防护（GORM 自动处理）
- XSS 防护（Vue 自动转义）

### 2. 访问控制
- 后续可添加用户认证
- API 访问限流
- CORS 配置

### 3. 数据安全
- 数据库密码加密存储
- Webhook URL 保密
- 敏感信息脱敏

### 4. 容器安全
- 使用官方镜像
- 定期更新镜像
- 最小权限原则
