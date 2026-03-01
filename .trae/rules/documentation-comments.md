---
alwaysApply: false
description: 需要编写golang代码时需要遵守的注释与文档规范
---
GoDoc 兼容:

所有导出的标识符（Type, Function, Interface, Variable）必须有以该标识符开头的文档注释。

格式：// FunctionName 描述了该函数的功能...。

自解释代码: 优先通过良好的命名减少注释。仅在解释“为什么”采用某种复杂算法或业务逻辑时使用行内注释。

包注释: 文件顶部必须有包级别的说明。

TODO 标准: 使用 // TODO(username): 待办事项说明 格式标注待优化项。