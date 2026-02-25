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

| 工具 | 用途 | 命令示例 |
|------|------|----------|
| **Subfinder** | 子域名收集 | `subfinder -d domain.com -all -silent -o subs.txt` |
| **Samoscout** | 子域名收集（补充） | `samoscout -d domain.com -silent -o subs.txt` |
| **Ksubdomain** | DNS 存活验证 | `ksubdomain verify -f subs.txt --silent -o output.txt` |
| **Httpx** | HTTP 存活验证 | `httpx -l subs.txt -sc -title -td -json -o httpx.json` |

工具安装位置：`/root/go/bin/`

## 🚀 快速开始

### 前置要求

1. **Docker 和 Docker Compose**
2. **扫描工具**（必须安装）

### 安装扫描工具

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

# 确保工具有执行权限
chmod +x ~/go/bin/*

# 验证安装
subfinder -version
httpx -version
```

### 部署项目

```bash
# 1. 克隆项目
git clone <your-repo-url>
cd recon

# 2. 一键部署
chmod +x deploy.sh
./deploy.sh
```

### 访问系统

- **前端界面**: `http://你的IP:8080`
- **后端 API**: `http://你的IP:8000`
- **健康检查**: `http://你的IP:8000/health`

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
# 一键部署（推荐）
./deploy.sh

# 启动项目
./start.sh

# 查看日志
docker-compose logs -f backend

# 查看服务状态
docker-compose ps

# 重新构建
./rebuild.sh

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 卸载（保留镜像和数据）
./uninstall.sh

# 完全卸载（删除所有）
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
├── frontend/            # Vue 3 前端（暗色主题）
│   └── src/
│       ├── api/         # API 封装
│       ├── router/      # 路由配置
│       └── views/       # 页面组件
├── check-tools.sh       # 检查工具
├── deploy.sh            # 一键部署
├── docker-compose.yml   # Docker 编排
├── rebuild.sh           # 重新构建
├── start.sh             # 启动项目
└── uninstall.sh         # 卸载脚本
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

```
1. Subfinder 收集子域名
2. Samoscout 补充收集
3. Ksubdomain DNS 验证（过滤泛解析）
4. Httpx HTTP 验证（状态码、标题、技术栈）
5. 数据对比分析（新增、恢复、失效）
6. 通知推送（企业微信/钉钉）
```

## 🐛 故障排查

### 工具找不到错误

```
错误: fork/exec /usr/local/bin/subfinder: no such file or directory
```

**解决方案：**

```bash
# 1. 检查工具是否安装
ls -lh /root/go/bin/

# 2. 确保有执行权限
chmod +x /root/go/bin/*

# 3. 重新部署
./deploy.sh

# 4. 验证容器内工具
docker-compose exec backend ls -lh /usr/local/bin/subfinder
```

### 前端没有暗色主题

前端需要重新构建才能应用新的样式：

```bash
# 重新构建前端
docker-compose build --no-cache frontend
docker-compose up -d frontend

# 或者完整重新部署
./deploy.sh
```

### 后端容器不断重启

```bash
# 查看日志
docker-compose logs -f backend

# 检查数据库连接
docker-compose ps postgres

# 重新构建
./rebuild.sh
```

### 扫描无结果

```bash
# 1. 检查工具是否正常
docker-compose exec backend subfinder -version
docker-compose exec backend httpx -version

# 2. 手动测试工具
docker-compose exec backend subfinder -d example.com -silent

# 3. 查看扫描日志
docker-compose logs -f backend | grep -A 10 "扫描"
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
- [ ] 支持多域名批量扫描
- [ ] 优化通知模板
- [ ] 添加 Web 界面配置

## ⚠️ 注意事项

1. **工具路径固定** - 必须安装在 `/root/go/bin/`
2. **权限要求** - 工具必须有执行权限（`chmod +x`）
3. **环境依赖** - 换服务器需要重新安装工具
4. **挂载配置** - docker-compose.yml 中的挂载路径必须正确
5. **首次启动** - 会下载 Docker 镜像，需要一些时间
6. **数据持久化** - 数据存储在 Docker 卷中，卸载时不会丢失（除非使用 --full）

## 📄 更新日志

### v2.0.0 - 2026-02-25

- ✅ 架构重构：从容器内 Docker 调用改为使用 VPS 本地工具
- ✅ 新增 Samoscout 和 Ksubdomain 工具
- ✅ 移除 Assetfinder 和 cert.sh
- ✅ 全新暗色主题 UI
- ✅ 优化扫描流程和日志显示
- ✅ 精简脚本，保留 5 个核心脚本

### v1.0.0 - 2026-02-24

- 初始版本发布

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**Made with ❤️ for Bug Bounty Hunters**
