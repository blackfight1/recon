# 🔍 自动化侦查平台 (Recon Platform)

> 专为漏洞赏金猎人打造的资产监控与变更追踪系统

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-4FC08D?logo=vue.js)](https://vuejs.org)

## ✨ 功能特性

### MVP 版本 (当前)

- ✅ **子域名自动收集** - 集成 Subfinder、Assetfinder、cert.sh
- ✅ **存活验证** - 使用 Httpx 批量验证并获取详细信息
- ✅ **智能差量对比** - 自动识别新增、恢复、失效的资产
- ✅ **实时通知** - 支持企业微信和钉钉机器人
- ✅ **Web 管理界面** - 现代化的 Vue 3 界面
- ✅ **定时自动扫描** - 可配置的扫描间隔
- ✅ **变更历史追踪** - 完整的资产变更时间轴

### 即将推出

- 🔜 端口扫描 (Naabu + Nmap)
- 🔜 网页截图 (Gowitness)
- 🔜 漏洞扫描 (Nuclei)
- 🔜 指纹识别增强

## 🎯 核心优势

- **变更优先** - 专注于资产变化，第一时间发现新资产
- **容器化部署** - 扫描工具按需启动，节省资源
- **历史回溯** - 完整的资产变更历史记录
- **个人优化** - 专为个人使用场景设计，轻量高效

## 🚀 快速开始

### 一键部署

**Linux/Mac:**
```bash
chmod +x rebuild.sh
./rebuild.sh
```

**Windows:**
```bash
start.bat
```

### 访问系统

- 本地: http://localhost:8080
- VPS: http://你的IP:8080

### 卸载

```bash
chmod +x uninstall.sh
./uninstall.sh
```

详细部署指南请查看 [VPS_DEPLOYMENT.md](VPS_DEPLOYMENT.md)

## 📖 文档

- [部署清单](DEPLOY.md) - 快速部署步骤 ⭐
- [VPS 部署](VPS_DEPLOYMENT.md) - VPS 服务器部署和故障排查
- [安装指南](INSTALL.md) - 详细安装说明
- [使用指南](USAGE.md) - 功能使用说明
- [API 文档](API.md) - API 接口文档

## 🏗️ 技术架构

```
┌─────────────┐
│  Vue 3 前端  │  Element Plus UI
└──────┬──────┘
       │ REST API
┌──────▼──────┐
│  Gin 后端   │  Go + GORM
└──────┬──────┘
       │
┌──────▼──────┐
│ PostgreSQL  │  数据存储
└─────────────┘
       │
┌──────▼──────────────────┐
│  Docker 容器化扫描工具   │
│  Subfinder | Httpx      │
│  Assetfinder | ...      │
└─────────────────────────┘
```

### 技术栈

**后端:**
- Go 1.21+ - 高性能并发处理
- Gin - 轻量级 Web 框架
- GORM - ORM 框架
- PostgreSQL - 关系型数据库
- robfig/cron - 定时任务调度

**前端:**
- Vue 3 - 渐进式框架
- Vite - 极速构建工具
- Element Plus - 企业级 UI 组件
- Axios - HTTP 客户端

**扫描工具:**
- Subfinder - 子域名收集
- Assetfinder - 子域名收集
- Httpx - HTTP 探测
- cert.sh - 证书透明度查询

**部署:**
- Docker - 容器化
- Docker Compose - 服务编排
- Nginx - 前端服务器

## 📊 系统截图

### 仪表盘
展示资产总览和关键指标

### 目标管理
添加和管理监控目标

### 子域名列表
查看所有发现的子域名及详细信息

### 变更中心
时间轴展示所有资产变更

## 📋 常用命令

```bash
# 首次部署或代码更新
./rebuild.sh                  # 重新构建（利用缓存）

# 快速重启（不重新构建）
./restart.sh                  # 或 docker-compose restart

# 启动服务
./start.sh                    # 或 docker-compose up -d

# 查看日志
docker-compose logs -f

# 查看状态
docker-compose ps

# 停止服务
docker-compose down

# 卸载系统（保留基础镜像）
./uninstall.sh
```

## 💡 缓存优化

项目已优化 Docker 构建缓存：

- **基础镜像**（golang、postgres、nginx）会被保留
- **Go 依赖包**会被缓存，只有 go.mod 改变时才重新下载
- **扫描工具镜像**默认保留，避免重复下载
- 只有代码改变时才重新编译

这样每次更新代码后重新构建会快很多！

## 📝 使用示例

### 1. 添加监控目标

```bash
curl -X POST http://localhost:8000/api/targets \
  -H "Content-Type: application/json" \
  -d '{"domain":"example.com","description":"测试目标"}'
```

### 2. 手动触发扫描

```bash
curl -X POST http://localhost:8000/api/targets/1/scan
```

### 3. 查看变更

```bash
curl http://localhost:8000/api/changes | jq
```

## 🛠️ 开发

### 后端开发

```bash
cd backend
go mod download
go run main.go
```

### 前端开发

```bash
cd frontend
npm install
npm run dev
```

### 构建

```bash
# 构建后端
cd backend
go build -o main .

# 构建前端
cd frontend
npm run build
```

## 🐛 故障排查

### 常见问题

1. **扫描无结果**
   - 检查目标域名格式
   - 查看后端日志: `docker-compose logs backend`
   - 确认网络连接正常

2. **通知未收到**
   - 验证 Webhook URL
   - 检查通知配置是否启用
   - 查看后端日志确认发送状态

3. **服务无法访问**
   - 检查容器状态: `docker-compose ps`
   - 确认端口未被占用
   - 查看防火墙设置

更多问题请查看 [INSTALL.md](INSTALL.md)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## ⚠️ 免责声明

本工具仅供授权的安全测试使用。使用者应遵守相关法律法规，仅对有权限的目标进行测试。作者不对任何非法使用承担责任。

## 🙏 致谢

感谢以下开源项目：
- [ProjectDiscovery](https://github.com/projectdiscovery) - Subfinder, Httpx, Naabu, Nuclei
- [tomnomnom](https://github.com/tomnomnom) - Assetfinder
- [Gin](https://github.com/gin-gonic/gin)
- [Vue.js](https://github.com/vuejs/core)
- [Element Plus](https://github.com/element-plus/element-plus)

## 📮 联系方式

- 提交 Issue: [GitHub Issues](https://github.com/yourusername/recon-platform/issues)
- 讨论: [GitHub Discussions](https://github.com/yourusername/recon-platform/discussions)

---

⭐ 如果这个项目对你有帮助，请给个 Star！
