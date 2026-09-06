# Changelog

本项目遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/) 格式，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### Added
- **测试深化**：feed 与 interaction 领域层单元测试（场景归一化、热榜权重、评论校验、状态判断），Redis 缓存页/卡片命中与回源路径测试
- **接口测试**：feed handler 请求-响应断言（真实 service 链路 + 仓储 stub），集成测试错误响应统一契约断言（`{"error": ...}`）
- **覆盖率门槛**：CI 低于阈值（当前 8%，基线 9.1%）即失败，防止覆盖率随代码增长静默退化
- **错误日志**：三个异步 worker（互动落库/fanout 分发/embedding 生成）失败时记录带事件上下文的日志，便于排查
- `2026-08-25` 2ef1fe6 worker 失败日志（user/video/author 上下文）
- `2026-08-30` 843aef2 feed 领域单测（17 用例）
- `2026-08-31` 32ddbaa interaction 领域单测（约 30 用例）
- `2026-09-01` ed0b793 页/卡片缓存命中路径测试
- `2026-09-03` b8d8915 feed handler 请求-响应测试
- `2026-09-04` c9b2793 集成测试断言强化（错误契约断言）
- `2026-09-05` 51a5697 CI 覆盖率防退化门槛

### Changed
- **错误响应统一**：全部 HTTP 错误返回一致的 `{"error": "..."}` 结构；鉴权中间件 `message` 字段统一为 `error`；500 错误不再泄漏内部细节（详情进日志）
- **敏感信息脱敏**：GORM SQL 日志限制为 Error 级，避免慢查询打印含绑定参数（如密码哈希）的语句
- **常量化**：Redis key 前缀（8 个）与分页上限抽取为命名常量，路由分组补 RESTful 语义注释
- `2026-08-23` 991171e 统一错误响应结构
- `2026-08-26` 91cf4cc GORM 日志脱敏
- `2026-08-27` 5a734d8 Redis key 前缀常量
- `2026-08-28` 794ac77 视频分页上限常量
- `2026-08-29` eca1e23 路由分组注释

## [0.1.0] - 2026-08-13

首个可运行版本：短视频 Feed 后端，覆盖 Feed 流（时间线/热榜/关注）、互动（点赞/收藏/评论，异步落库）、账号（JWT 鉴权）、发布（推拉混合分发）与监控（Prometheus 指标）。

## [0.2.0] - 2026-08-24

### Added
- **CI 工作流**（2026-08-15）：Go vet + Go test + 前端构建，push/PR 触发
- **CI 依赖缓存**（2026-08-16）：Go module 与 npm 缓存，加速二次构建
- **测试覆盖率报告**（2026-08-17）：CI 输出覆盖率摘要
- **根目录 Makefile**（2026-08-18）：`make dev / test / build / lint`
- **前端测试**（2026-08-19）：vitest + 组件冒烟测试
- **配置必填校验**（2026-08-21）：port/DSN/Redis 地址缺失时启动即报错
- **分页参数校验**（2026-08-22）：limit 1-100、offset >= 0，非法返回 400
- **JWT TTL 配置核对**（2026-08-24）：锁定 `jwt_access_ttl` 从配置生效的行为

### Changed
- **Go module 卫生检查**（2026-08-20）：CI 中 `go mod tidy -diff` + `go mod verify`
- **LF 行尾强制**（2026-08-20）：Go module 文件统一行尾

### Fixed
- **互动并发冲突错误码**（2026-08-21）：高并发点赞下 Redis 乐观锁冲突（`TxFailedErr`）从误报 500 修正为 409 Conflict（压测发现，三层修复 + 单测兜底）

[Unreleased]: https://github.com/rickererer/PulseFeed/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/rickererer/PulseFeed/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/rickererer/PulseFeed/releases/tag/v0.1.0
