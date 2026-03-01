---
alwaysApply: false
description: 需要编写golang代码时需要遵守的代码风格规范
---
代码风格规范 (Coding Style)
命名公约:

变量与函数使用 mixedCaps（小驼峰），导出名使用 MixedCaps（大驼峰）。

缩写词应全大写（如 URL, HTTP, ID），不要写成 Url 或 Http。

接口命名通常以 er 结尾（如 Reader, Writer, Storer）。

Receiver 规范:

统一使用指针接收器 (s *Service) 以保证一致性，除非是极小的基础类型。

接收器命名通常取结构体首字母的小写（如 u *User），严禁使用 this 或 self。

显式错误处理:

绝对禁止忽略错误（禁止 _ = DoSomething()）。

错误信息应为小写且不以标点结尾（如 fmt.Errorf("user not found: %w", err)）。

优先使用 errors.Is 和 errors.As 进行错误判定。