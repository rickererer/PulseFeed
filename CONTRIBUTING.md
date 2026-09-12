# 贡献指南

感谢关注 PulseFeed。本文档说明如何搭建开发环境、跑通验证命令以及提交变更。

## 环境要求

| 工具 | 版本 | 用途 |
|---|---|---|
| Go | 1.26+ | 后端（`apps/api`） |
| Node.js | 22+ | 前端（`apps/web`） |
| Docker Compose | v2 | 一键启动 MySQL/Redis/RabbitMQ/API/Web/Grafana |

## 快速开始

```bash
# 一键启动全部服务（首次会构建镜像）
cd apps && docker compose up -d --build

# Web:  http://127.0.0.1:5173
# Grafana: http://127.0.0.1:3000（面板 PulseFeed Overview）
```

不使用容器时，也可以只起依赖，本地跑 API：

```bash
cd apps/api && go run ./cmd/feed
```

## 配置

后端配置来自 `apps/api/configs/config.yaml`。配置由 YAML 文件直接读取（**不读取环境变量**），因此每次调整都改文件本身：

```bash
cp apps/api/configs/config.example.yaml apps/api/configs/config.yaml
```

随后按本机环境填写。注意事项：

- `config.example.yaml` 是字段模板，敏感项一律为占位符；**真实凭据不要提交**，本地覆盖可放在 `config.local.yaml`（已在 `.gitignore`）
- 容器内运行使用 `config.docker.yaml`，其中 host 指向 compose 服务名（`mysql` / `redis` / `rabbitmq`）
- 字段缺失或非法会在启动阶段由 `Config.Validate` 直接报错，不会静默使用默认值
- RabbitMQ 为可选项：`rabbitmq.url` 留空即禁用异步队列（相关 consumer 与 publish 路径不启动）

## 常用命令

根目录 Makefile 提供统一入口：

| 命令 | 作用 |
|---|---|
| `make dev` | 本地开发模式启动 |
| `make test` | 运行后端与前端测试 |
| `make build` | 构建后端与前端产物 |
| `make lint` | 代码静态检查（go vet） |
| `make clean` | 清理构建产物 |

后端与前端也可以分别执行：

```bash
cd apps/api && go test ./... -covermode=count
cd apps/web && npm test && npm run build
```

## 提交变更

### 分支与提交

1. 从 `main` 拉出特性分支，命名如 `feat/xxx`、`fix/xxx`、`docs/xxx`。
2. 提交信息使用 Conventional Commits 风格（`feat:`、`fix:`、`docs:`、`test:`、`refactor:`、`ci:`）。
3. 一次 PR 聚焦一件事；拒绝空提交与注水提交，小步真实改进。

### 提交前自检

本地全部通过后再推送：

```bash
cd apps/api
go vet ./...
go mod tidy -diff && go mod verify   # module 卫生（CI 同样检查）
go test ./... -covermode=count        # 覆盖率不得低于 CI 门槛
cd ../web && npm test && npm run build
```

### CI 检查项

PR 会触发 GitHub Actions，全部通过才可合并：

- `go vet`
- Go module 卫生（tidy diff 为空 + verify 通过）
- 后端测试 + 覆盖率门槛（当前 8%，低于即失败）
- 前端 vitest + 构建

### 错误与风格约定

- 所有 HTTP 错误响应使用统一的 `{"error": "..."}` 结构（见 `CHANGELOG` Unreleased 段）。
- 500 错误返回固定文案，内部细节写日志而不是响应体。
- Redis key 前缀、分页上限等业务常量从 `internal/infra/cache/feed_cache.go` 与领域层常量引用，不要新写裸字符串/裸数字。
- 新增异步 worker 的失败分支需要带事件上下文的日志。

### 报告问题

提交 Issue 时请说明：复现步骤、期望行为、实际行为、相关日志片段。涉及接口行为的请附请求/响应示例。

## 版本与变更记录

合并到 `main` 的用户可见变更请同步更新 [CHANGELOG.md](CHANGELOG.md)（Keep a Changelog 格式）。
