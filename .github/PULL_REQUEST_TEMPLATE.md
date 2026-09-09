<!-- 提交前请确认本地自检已通过，参见 CONTRIBUTING.md -->

## 变更说明

<!-- 一段话说清本次变更做了什么、为什么 -->

## 变更类型

- [ ] `feat` 新功能
- [ ] `fix` 缺陷修复
- [ ] `refactor` 重构（不改行为）
- [ ] `test` 测试补充
- [ ] `docs` 文档
- [ ] `ci` / `build` 工程化

## 自检清单

- [ ] `go vet ./...` 通过
- [ ] `go mod tidy -diff` 无输出（module 卫生）
- [ ] 后端测试通过，覆盖率不低于 CI 门槛（当前 8%）
- [ ] `npm test` 与 `npm run build` 通过（涉及前端时）
- [ ] 新增错误响应使用统一 `{"error": "..."}` 结构，500 不泄漏内部细节
- [ ] 新增 Redis key / 业务常量使用命名常量，无裸字符串/裸数字
- [ ] 用户可见变更已同步 [CHANGELOG.md](../CHANGELOG.md)
