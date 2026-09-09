name: 功能建议
description: 提出新功能或改进想法
labels: ["enhancement"]
body:
  - type: textarea
    id: problem
    attributes:
      label: 要解决的问题
      description: 描述当前缺失或不便的场景，而不是直接给方案
    validations:
      required: true

  - type: textarea
    id: proposal
    attributes:
      label: 建议方案
      description: 大致思路即可；如涉及接口请说明路径、方法与响应结构
    validations:
      required: true

  - type: textarea
    id: alternatives
    attributes:
      label: 备选方案
      description: 你考虑过的其他做法，以及为什么不选

  - type: checkboxes
    id: scope
    attributes:
      label: 范围确认
      options:
        - label: 我已查看 [roadmap](../docs/roadmap.md)，该能力不在已有规划中
        - label: 我已查看 [CHANGELOG](../CHANGELOG.md)，该能力尚未实现
