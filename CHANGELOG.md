# 更新日志

## v2.0.0 - 2026-02-25

### 🎉 重大更新

#### 架构重构
- ✅ 从容器内 Docker 调用改为使用 VPS 本地工具（方案二）
- ✅ 移除 Docker socket 挂载，提升安全性
- ✅ 工具通过卷挂载到容器，独立更新无需重建镜像

#### 扫描工具优化
- ✅ 新增 **Samoscout** - 子域名收集工具
- ✅ 新增 **Ksubdomain** - DNS 存活验证，过滤泛解析
- ✅ 移除 **Assetfinder** - 由 Samoscout 替代
- ✅ 移除 **cert.sh** - 简化扫描流程
- ✅ 优化 **Httpx** - 完整支持 JSON 输出，提取技术栈

#### 前端美化
- ✅ 全新暗色主题 UI
- ✅ 现代化设计，护眼舒适
- ✅ 优化快速扫描页面，实时动画效果
- ✅ 改进日志显示，图标化状态
- ✅ 响应式布局优化

#### 代码优化
- ✅ 重构 `scanner.go`，使用本地命令
- ✅ 扩展 `HttpxResult` 结构体，支持完整字段
- ✅ 新增 `readSubdomainsFromFile` 辅助函数
- ✅ 优化错误处理和日志记录

#### 脚本精简
- ✅ 保留 5 个核心脚本
- ❌ 删除 `fix.sh`, `logs.sh`, `debug.sh`, `quick-fix.sh`
- ✅ 新增 `check-tools.sh` - 工具检查脚本
- ✅ 新增 `deploy.sh` - 一键部署脚本

### 📋 扫描流程

```
1. Subfinder 收集子域名
2. Samoscout 补充收集
3. Ksubdomain DNS 验证
4. Httpx HTTP 验证
5. 数据对比分析
6. 通知推送
```

### 🔧 工具要求

必须在 VPS 上安装以下工具到 `/root/go/bin/`：

```bash
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install -v github.com/samogod/samoscout/cmd/samoscout@latest
go install -v github.com/boy-hack/ksubdomain/cmd/ksubdomain@latest
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
```

### 📦 部署

```bash
# 1. 检查工具
./check-tools.sh

# 2. 一键部署
./deploy.sh
```

### 🎨 UI 预览

- 暗色主题配色
- 实时进度条动画
- 彩色日志图标
- 现代化卡片设计

### ⚠️ 破坏性变更

- 不再支持容器内 Docker 调用
- 必须在主机上安装扫描工具
- 移除 cert.sh 支持

---

## v1.0.0 - 2026-02-24

### 初始版本

- 基础资产监控功能
- 快速扫描功能
- 子域名收集（Subfinder, Assetfinder, cert.sh）
- 存活验证（Httpx）
- 变更检测
- 通知推送（企业微信、钉钉）
- Docker 部署
