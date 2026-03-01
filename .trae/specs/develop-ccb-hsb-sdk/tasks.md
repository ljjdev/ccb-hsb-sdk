# Tasks

- [x] Task 1: 创建项目基础结构
  - [x] 创建pkg目录结构(pkg/client, pkg/signature, pkg/model, pkg/config)
  - [x] 创建internal目录结构(internal/utils)
  - [x] 创建examples目录用于示例代码
  - [x] 更新go.mod文件,添加必要的依赖
  - [x] 创建Makefile,包含build、test、fmt、vet、lint等目标

- [x] Task 2: 实现配置管理模块
  - [x] 定义Config结构体,包含平台公钥、客户方私钥、市场编号、商户编号等配置项
  - [x] 实现配置加载和验证功能
  - [x] 提供默认配置和自定义配置选项

- [x] Task 3: 实现数据模型
  - [x] 定义CcbPaymentInfo结构体(支付订单主信息)
  - [x] 定义CcbSubOrder结构体(子订单信息)
  - [x] 定义CcbLedgerAccount结构体(分账方信息)
  - [x] 定义CcbCoupon结构体(消费券信息)
  - [x] 定义支付请求和响应结构体(PlaceOrderRequest, PlaceOrderResponse)
  - [x] 定义退款请求和响应结构体(RefundOrderRequest, RefundOrderResponse)
  - [x] 定义查询支付结果请求和响应结构体(QueryOrderRequest, QueryOrderResponse)
  - [x] 定义查询退款结果请求和响应结构体(QueryRefundRequest, QueryRefundResponse)
  - [x] 实现JSON序列化和反序列化标签

- [x] Task 4: 实现RSA签名与验签模块
  - [x] 实现RSA密钥加载功能(支持PEM格式)
  - [x] 实现SHA256withRSA签名功能
  - [x] 实现SHA256withRSA验签功能
  - [x] 添加单元测试

- [x] Task 5: 实现JSON排序拼接模块
  - [x] 实现JSON对象递归排序功能
  - [x] 实现字段过滤功能(排除SIGN_INF、Svc_Rsp_St、Svc_Rsp_Cd、Rsp_Inf)
  - [x] 实现签名字符串生成功能
  - [x] 添加单元测试

- [x] Task 6: 实现HTTP客户端模块
  - [x] 定义HTTPClient结构体
  - [x] 实现POST请求发送功能
  - [x] 实现请求头设置(Content-Type等)
  - [x] 实现超时控制
  - [x] 实现错误处理

- [x] Task 7: 实现支付订单生成接口
  - [x] 定义PlaceOrder方法
  - [x] 实现请求JSON构建(包含主订单、子订单、分账方、消费券)
  - [x] 实现请求签名
  - [x] 实现请求发送
  - [x] 实现响应验签
  - [x] 实现支付URL提取和URL解码
  - [x] 添加单元测试

- [x] Task 8: 实现查询支付结果接口
  - [x] 定义QueryOrder方法
  - [x] 实现请求JSON构建
  - [x] 实现请求签名
  - [x] 实现请求发送
  - [x] 实现响应验签
  - [x] 添加单元测试

- [x] Task 9: 实现订单退款接口
  - [x] 定义RefundOrder方法
  - [x] 实现请求JSON构建(支持全额退款和部分退款)
  - [x] 实现请求签名
  - [x] 实现请求发送
  - [x] 实现响应验签
  - [x] 添加单元测试

- [x] Task 10: 实现查询退款结果接口
  - [x] 定义QueryRefund方法
  - [x] 实现请求JSON构建
  - [x] 实现请求签名
  - [x] 实现请求发送
  - [x] 实现响应验签
  - [x] 添加单元测试

- [x] Task 11: 实现SDK客户端
  - [x] 定义Client结构体
  - [x] 实现NewClient构造函数
  - [x] 集成配置管理、签名、HTTP客户端等模块
  - [x] 提供便捷的API方法(PlaceOrder, QueryOrder, RefundOrder, QueryRefund)

- [x] Task 12: 编写示例代码
  - [x] 创建支付订单示例(atherPlaceorder)
  - [x] 创建查询支付结果示例(paymentQuery)
  - [x] 创建订单退款示例(refund)
  - [x] 创建查询退款结果示例(refundQuery)
  - [x] 展示配置初始化
  - [x] 展示错误处理
  - [x] 添加详细注释

- [x] Task 13: 编写文档
  - [x] 编写包级别文档
  - [x] 编写API文档(GoDoc格式)
  - [x] 编写README.md(使用说明、快速开始、API参考)
  - [x] 编写集成指南

- [x] Task 14: 添加单元测试
  - [x] 为签名模块添加测试
  - [x] 为JSON拼接模块添加测试
  - [x] 为HTTP客户端添加测试
  - [x] 为支付订单接口添加测试(使用mock)
  - [x] 为查询支付结果接口添加测试(使用mock)
  - [x] 为订单退款接口添加测试(使用mock)
  - [x] 为查询退款结果接口添加测试(使用mock)

- [x] Task 15: 代码质量检查
  - [x] 运行go fmt格式化代码
  - [x] 运行go vet检查代码
  - [x] 运行golint检查代码风格
  - [x] 确保所有导出函数都有文档注释

# Task Dependencies
- [Task 2] depends on [Task 1]
- [Task 3] depends on [Task 1]
- [Task 4] depends on [Task 1]
- [Task 5] depends on [Task 1]
- [Task 6] depends on [Task 1]
- [Task 7] depends on [Task 2, Task 3, Task 4, Task 5, Task 6]
- [Task 8] depends on [Task 2, Task 3, Task 4, Task 5, Task 6]
- [Task 9] depends on [Task 2, Task 3, Task 4, Task 5, Task 6]
- [Task 10] depends on [Task 2, Task 3, Task 4, Task 5, Task 6]
- [Task 11] depends on [Task 2, Task 4, Task 5, Task 6]
- [Task 12] depends on [Task 11]
- [Task 13] depends on [Task 11]
- [Task 14] depends on [Task 4, Task 5, Task 7, Task 8, Task 9, Task 10]
- [Task 15] depends on [Task 7, Task 8, Task 9, Task 10, Task 11]
