# 网络主题资源共享站 (NE Resource Community)

面向网络工程垂直领域的资源社区，聚合工具、配置模板、文档资料与学习资源。支持本地文件存储、基于 PostgreSQL 内置全文检索的多条件搜索，以及用户上传、评分与下载统计。

## 技术栈
- 前端：Vue 3 + Vite + TypeScript + Element Plus + Pinia + Vue Router
- 后端：Go (Gin) + GORM + PostgreSQL（内置 FTS `tsvector`/`websearch_to_tsquery`）
- 文件存储：本地磁盘 `server/uploads`（无病毒扫描，按需求可后续接入）

## 核心功能
- 结构化资源库：工具 / 配置模板 / 文档资料 / 学习资源分类
- 智能检索：多条件筛选 + PostgreSQL 全文检索 + 相似资源推荐
- 贡献与互动：用户上传、评分与评论、下载统计、积分占位
- 质量控制占位：支持去重/审核流程扩展（当前未启用病毒扫描）
- 个人空间占位：收藏、下载历史、进度/贡献统计可扩展

## 目录结构
- `server/`：Gin API、GORM 模型、JWT 鉴权、文件上传与静态服务
- `web/`：前端应用（页面：资源广场、上传、登录/注册、资源详情、个人空间）
- `server/uploads/`：本地上传文件目录（运行时自动创建，已在 `.gitignore`）

## 环境依赖
- Go 1.23+
- Node.js 18+/npm 9+
- PostgreSQL 14+（启用默认英文全文检索字典即可）

## 环境变量（后端）
| 变量名 | 说明 | 默认值 |
| --- | --- | --- |
| `SERVER_ADDR` | 监听地址 | `:8080` |
| `DATABASE_URL` | PostgreSQL 连接串 | `postgres://postgres:postgres@localhost:5432/ne_resource?sslmode=disable` |
| `JWT_SECRET` | JWT 签名密钥 | `dev-secret-change-me` |
| `UPLOAD_DIR` | 本地上传目录 | `uploads` |
| `ENV` | 运行环境标记 | `dev` |

## 快速开始
### 后端
```bash
cd server
# 安装依赖（已包含 go.sum）
go mod tidy
# 运行
 go run cmd/api/main.go
```
确保 PostgreSQL 中已创建数据库 `ne_resource`（或将 `DATABASE_URL` 指向已有库）。

### 前端
```bash
cd web
npm install
npm run dev  # 默认 5173，已代理 /api 与 /uploads 到 8080
```
生产构建：`npm run build`

## API 概览（简要）
- `GET /api/health` 健康检查
- `POST /api/auth/register` 注册；`POST /api/auth/login` 登录（返回 JWT）
- `GET /api/resources` 列表/筛选/全文搜索
- `POST /api/resources` 上传资源（需 Bearer Token，multipart）
- `GET /api/resources/:id` 详情
- `GET /api/resources/:id/recommendations` 相似推荐
- `GET /api/resources/:id/download` 下载并计数
- `POST /api/resources/:id/reviews` 评分/评价（需登录）

## 搜索与索引
`resources.search_vector` 为生成列，基于 `title/description/tags`，建有 GIN 索引用于 `websearch_to_tsquery`；无需额外扩展插件。

## 后续可扩展点
- 加入文件类型/大小限制、MD5 去重、审核流与管理员角色
- 收藏/下载历史接口与个人空间展示
- 更细的推荐策略（标签/协议/厂商综合）
- 对接对象存储或 CDN，接入病毒扫描（如 ClamAV）
