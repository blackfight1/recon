# 🔍 自动化侦查平台

一个专为漏洞赏金猎人设计的资产收集与监控平台，支持子域名自动发现、存活验证、变更检测和实时通知。

## ✨ 核心功能

- **资产监控** - 添加目标域名，自动定时扫描（每6小时）
- **快速扫描** - 输入域名立即扫描，实时查看进度和日志
- **子域名收集** - 集成 Subfinder、Samoscout
- **DNS 验证** - 使用 Ksubdomain 极速验证 DNS 存活，过滤泛解析
- **存活验证** - 使用 Httpx 验证子域名存活状态，识别状态码、标题、技术栈
- **变更检测** - 自动对比历史数据，发现新增/失效的子域名
- **通知推送** - 支持企业微信和钉钉通知
- **暗色主题** - 现代化暗色 UI，护眼舒适

## 🔧 扫描工具

项目使用本地安装的扫描工具（方案二架构）：

- **Subfinder** - 子域名收集
- **Samoscout** - 子域名收集（补充）
- **Ksubdomain** - DNS 存活验证，过滤泛解析
- **Httpx** - HTTP 存活验证，识别状态码、标题、技术栈

工具安装位置：`/root/go/bin/`

## 🚀 快速开始

### 前置要求

1. **安装 Docker 和 Docker Compose**
2. **安装扫描工具到 VPS**（必须）

```bash
# 安装 Go（如果未安装）
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:~/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装扫描工具
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install -v github.com/samogod/samoscout/cmd/samoscout@latest
go install -v github.com/boy-hack/ksubdomain/cmd/ksubdomain@latest
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest

# 验证安装
subfinder -version
httpx -version
```

### 1. 克隆项目

```bash
git clone <your-repo-url>
cd recon
```

### 2. 检查工具

```bash
chmod +x check-tools.sh
./check-tools.sh
```

### 3. 一键部署

```bash
chmod +x deploy.sh
./deploy.sh
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
# 检查工具
chmod +x check-tools.sh
./check-tools.sh

# 一键部署
chmod +x deploy.sh
./deploy.sh

# 启动项目
chmod +x start.sh
./start.sh

# 查看日志
docker-compose logs -f backend

# 查看服务状态
docker-compose ps

# 重新构建
chmod +x rebuild.sh
./rebuild.sh

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 卸载（保留镜像）
chmod +x uninstall.sh
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
- **前端**: Vue 3 + Element Plus + Vite（暗色主题）
- **部署**: Docker + Docker Compose
- **扫描工具**: Subfinder, Samoscout, Ksubdomain, Httpx（本地安装）

## 🏗️ 架构说明

项目采用**方案二架构**：扫描工具安装在 VPS 主机上，通过 Docker 卷挂载到容器内使用。

**优势：**
- 工具可独立更新，无需重新构建镜像
- 容器镜像体积小
- 调试方便，可在主机上直接测试工具
- 支持任何类型的扫描工具（Go、Python、Rust 等）

## 📊 扫描流程

1. **子域名收集** - Subfinder + Samoscout
2. **DNS 验证** - Ksubdomain 过滤泛解析
3. **HTTP 验证** - Httpx 获取状态码、标题、技术栈
4. **数据对比** - 发现新增、恢复、失效的子域名
5. **通知推送** - 企业微信 / 钉钉

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
