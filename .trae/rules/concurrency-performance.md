---
alwaysApply: false
description: 需要编写golang代码时需要遵守的并发与性能规范
---
Goroutine 管理: 启动 Goroutine 时必须明确其生命周期，优先通过 context.Context 控制退出。

Channel 使用: 仅在需要通信或同步时使用 Channel；简单的状态保护优先使用 sync.Mutex。

零分配意识: 在高频调用路径中，避免不必要的 interface{} 转换和内存逃逸。