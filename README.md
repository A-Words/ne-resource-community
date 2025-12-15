# 网络主题资源共享站 (NE Resource Community)

面向网络工程垂直领域的资源社区，聚合工具、配置模板、文档资料与学习资源。支持本地文件存储、基于 PostgreSQL 内置全文检索的多条件搜索，以及用户上传、评分与下载统计。

## 技术栈
- **前端**：Vue 3 + Vite + TypeScript + Element Plus + Pinia + Vue Router
- **后端**：Go (Gin) + GORM + PostgreSQL（内置 FTS `tsvector`/`websearch_to_tsquery`）
- **安全**：ClamAV (病毒扫描) + JWT (鉴权)
- **文件存储**：本地磁盘 `server/uploads`

## 核心功能
- **结构化资源库**：四大核心分类（网络工具、配置模板、文档资料、学习资源），支持针对性筛选。
- **资源求助**：用户可发布资源求助（Request），社区互助解决资源缺失问题。
- **灵活存储**：支持本地文件上传与外部直链（External Link）两种资源形态，节省存储空间。
- **智能检索**：多条件筛选 + PostgreSQL 全文检索 + 相似资源推荐。
- **贡献与互动**：用户上传、评分与评论、下载统计、收藏功能。
- **质量控制**：
    - **病毒扫描**：集成 ClamAV 对上传文件进行实时病毒检测。
    - **去重机制**：支持文件哈希去重。
    - **审核流程**：管理员审核流程（Pending/Approved/Rejected）。
- **个人空间**：收藏列表、下载历史、我的上传管理。
- **管理后台**：资源审核、举报处理、用户管理、系统审计日志。

## 目录结构
- `server/`
    - `cmd/api/`：API 服务入口
    - `cmd/promote_admin/`：管理员提权 CLI 工具
    - `internal/`：核心业务逻辑（Handlers, Models, Middleware, Scanner 等）
    - `uploads/`：本地上传文件目录
- `web/`：前端应用（基于 Vue 3 + TypeScript）

## 环境依赖
- Go 1.23+
- Node.js 18+/npm 9+
- PostgreSQL 14+（启用默认英文全文检索字典即可）
- ClamAV (可选，用于文件病毒扫描)

## 环境变量（后端）
| 变量名 | 说明 | 默认值 |
| --- | --- | --- |
| `SERVER_ADDR` | 监听地址 | `:8080` |
| `DATABASE_URL` | PostgreSQL 连接串 | `postgres://postgres:postgres@localhost:5432/ne_resource?sslmode=disable` |
| `JWT_SECRET` | JWT 签名密钥 | `dev-secret-change-me` |
| `UPLOAD_DIR` | 本地上传目录 | `uploads` |
| `ENV` | 运行环境标记 | `dev` |
| `CLAMAV_ADDR` | ClamAV 服务地址 | `tcp://localhost:3310` |

## 快速开始

### 1. 数据库准备
确保 PostgreSQL 中已创建数据库 `ne_resource`。

### 2. 后端启动
```bash
cd server
# 安装依赖
go mod tidy
# 运行服务
go run cmd/api/main.go
```

### 3. 创建管理员
注册一个普通用户后，使用 CLI 工具将其提升为管理员：
```bash
# 在 server 目录下
go run cmd/promote_admin/main.go -email <your_email>
```

### 4. 前端启动
```bash
cd web
# 安装依赖
npm install
# 开发模式运行
npm run dev
```
访问 `http://localhost:5173` 即可体验。
