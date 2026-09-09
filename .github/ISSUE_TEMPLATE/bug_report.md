name: Bug 报告
description: 报告功能异常或行为与预期不符
labels: ["bug"]
body:
  - type: markdown
    attributes:
      value: |
        感谢反馈。请尽量填写完整信息，缺少复现步骤的问题可能被关闭。

  - type: textarea
    id: repro-steps
    attributes:
      label: 复现步骤
      description: 从哪个操作开始，按什么顺序执行
      placeholder: |
        1. 启动服务 `cd apps && docker compose up -d --build`
        2. 调用 `PUT /api/videos/1/like` ...
    validations:
      required: true

  - type: textarea
    id: expected
    attributes:
      label: 期望行为
    validations:
      required: true

  - type: textarea
    id: actual
    attributes:
      label: 实际行为
      description: 请附响应体 / 日志片段
    validations:
      required: true

  - type: textarea
    id: request-response
    attributes:
      label: 请求/响应示例
      description: 涉及接口行为时请提供（含状态码与响应体）
      placeholder: |
        > PUT /api/videos/1/like
        HTTP 500
        {"error": "internal server error"}

  - type: dropdown
    id: component
    attributes:
      label: 影响组件
      options:
        - api（后端）
        - web（前端）
        - worker（异步任务）
        - monitoring（指标/面板）
        - 部署/环境
    validations:
      required: true

  - type: input
    id: env
    attributes:
      label: 运行环境
      description: 操作系统 / Docker 版本 / Go 或 Node 版本
    validations:
      required: true
