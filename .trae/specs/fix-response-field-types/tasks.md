# Tasks

- [x] Task 1: 修正 QueryOrderResponse 的 Txnamt 字段类型
  - [x] SubTask 1.1: 将 Txnamt 字段类型从 string 改为 float64
  - [x] SubTask 1.2: 更新字段注释

- [x] Task 2: 修正 QueryRefundResponse 的 RfndAmt 字段类型
  - [x] SubTask 2.1: 将 RfndAmt 字段类型从 string 改为 *string
  - [x] SubTask 2.2: 更新字段注释

- [x] Task 3: 验证字段类型修正
  - [x] SubTask 3.1: 运行测试确保 JSON 解析正确
  - [x] SubTask 3.2: 检查 go vet 无错误

# Task Dependencies
无依赖关系，所有任务可以并行执行
