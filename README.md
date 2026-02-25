# 🔍 自动化侦查平台

一个专为漏洞赏金猎人设计的资产收集与监控平台，支持子域名自动发现、存活验证、变更检测和实时通知。

## ✨ 核心功能

- **资产监控** - 添加目标域名，自动定时扫描（每6小时）
- **快速扫描** - 输入域名立即扫描，实时查看进度和日志
- **子域名收集** - 集成 Subfinder、Assetfinder、cert.sh
- **存活验证** - 使用 Httpx 验证子域名存活状态
- **变更检测** - 自动对比历史数据，发现新增/失效的子域名
- **通知推送** - 支持企业微信和钉钉通知

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone <your-repo-url>
cd recon
```

### 2. 配置（可选）

```bash
cp backend/config.example.yaml backend/config.yaml
# 编辑配置文件，修改通知 webhook 等
nano backend/config.yaml
```

### 3. 一键启动

```bash
chmod +x *.sh
./start.sh
```

### 4. 访问系统

- 前端界面: `http://你的IP:8080`
- 后端 API: `http://你的IP:8000`
- 健康检查: `http://你的IP:8000/health`

## 📖 使用说明

### 资产监控模式

1. 进入"目标管理"页面
2. 点击"添加目标"，输入域名（如 example.com）
3. 系统自动每6小时扫描一次
4. 在"变更中心"查看新发现的子域名

### 快速扫描模式

1. 进入"快速扫描"页面
2. 输入域名，点击"开始扫描"
3. 实时查看扫描进度（0-100%）和详细日志
4. 扫描完成后查看结果

## 🛠️ 常用命令

```bash
# 启动项目
./start.sh

# 查看后端日志
./logs.sh

# 诊断问题
./debug.sh

# 查看服务状态
docker-compose ps

# 重新构建
./rebuild.sh

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 卸载（保留镜像）
./uninstall.sh

# 完全卸载
./uninstall.sh --full
```

## 📁 项目结构

```
recon/
├── backend/              # Go 后端
│   ├── config/          # 配置管理
│   ├── controllers/     # API 控制器
│   ├── database/        # 数据库连接
│   ├── models/          # 数据模型
│   ├── notifier/        # 通知模块
│   ├── router/          # 路由配置
│   ├── scanner/         # 扫描引擎（核心）
│   └── scheduler/       # 定时任务
├── frontend/            # Vue 3 前端
│   └── src/
│       ├── api/         # API 封装
│       ├── router/      # 路由配置
│       └── views/       # 页面组件
├── docker-compose.yml   # Docker 编排
├── start.sh            # 一键启动
├── rebuild.sh          # 重新构建
└── uninstall.sh        # 卸载脚本
```

## 🔧 技术栈

- **后端**: Go 1.21 + Gin + GORM + PostgreSQL
- **前端**: Vue 3 + Element Plus + Vite
- **部署**: Docker + Docker Compose
- **扫描工具**: Subfinder, Assetfinder, Httpx, cert.sh

## 🐛 故障排查

### 后端容器不断重启

```bash
# 1. 查看后端日志
docker-compose logs -f backend

# 2. 运行诊断脚本
chmod +x debug.sh
./debug.sh

# 3. 检查常见问题
# - 数据库是否启动成功
# - 配置文件是否存在
# - 环境变量是否正确
```

### 服务无法启动

```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend
docker-compose logs -f frontend

# 重新构建
./rebuild.sh

# 诊断问题
./debug.sh
```

### 常用命令

```bash
# 查看后端日志
./logs.sh
# 或
docker-compose logs -f backend

# 重启单个服务
docker-compose restart backend

# 进入容器调试
docker-compose exec backend sh
docker-compose exec postgres psql -U recon -d recon
```

### 扫描工具无法运行

```bash
# 检查 Docker 网络
docker network ls

# 手动测试工具
docker run --rm projectdiscovery/subfinder:latest -version
docker run --rm projectdiscovery/httpx:latest -version
```

### 数据库连接失败

```bash
# 检查数据库容器
docker-compose ps postgres

# 进入数据库
docker-compose exec postgres psql -U recon -d recon
```

## 📝 API 文档

详见 [API.md](API.md)

主要端点：
- `GET /api/stats` - 统计信息
- `GET /api/targets` - 获取目标列表
- `POST /api/targets` - 添加目标
- `POST /api/targets/:id/scan` - 手动触发扫描
- `POST /api/quick-scan` - 快速扫描
- `GET /api/tasks/:id/logs` - 获取任务日志
- `GET /api/tasks/:id/progress` - 获取任务进度
- `GET /api/subdomains` - 获取子域名列表
- `GET /api/changes` - 获取变更日志

## 🎯 下一步计划

- [ ] 添加端口扫描（Naabu）
- [ ] 添加指纹识别（Wappalyzer）
- [ ] 添加截图功能（Gowitness）
- [ ] 添加漏洞扫描（Nuclei）
- [ ] 优化通知模板
- [ ] 添加 Web 界面配置

## ⚠️ 注意事项

1. 首次启动会下载 Docker 镜像，需要一些时间
2. 扫描工具镜像会被保留，避免重复下载
3. 数据存储在 Docker 卷中，卸载时不会丢失（除非使用 --full）
4. 建议在 VPS 上运行，本地开发需要修改配置
5. 快速扫描不会添加到监控列表，适合临时测试

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**Made with ❤️ for Bug Bounty Hunters**
