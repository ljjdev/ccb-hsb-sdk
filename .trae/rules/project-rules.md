---
alwaysApply: false
description: 需要编写代码或创建文件以及文件夹时需要遵守的项目规范
---
目录与项目结构规范 (Project Structure)
遵循 Standard Go Layout: 必须采用社区公认的目录结构。

/cmd: 包含主应用程序入口（如 main.go），每个二进制文件一个子目录。

/internal: 存放不希望被外部项目引用的私有代码，确保封装性。

/pkg: 存放可以被外部项目安全使用的库代码。

/api: 存放 API 定义文件（如 Proto, OpenAPI/Swagger）。

扁平化原则: 避免过度嵌套。包名应简短且具有描述性（如 user 而非 user_management_service）。

禁止循环依赖: 设计包时必须确保单向依赖流。