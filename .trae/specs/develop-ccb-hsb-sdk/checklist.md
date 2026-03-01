# Checklist

- [ ] 项目基础结构符合Standard Go Layout规范
- [ ] pkg目录结构正确创建(client, signature, model, config)
- [ ] internal目录结构正确创建
- [ ] examples目录正确创建
- [ ] go.mod文件包含必要的依赖
- [ ] Makefile创建成功,包含build、test、fmt、vet、lint等目标

- [ ] Config结构体包含所有必需的配置项(平台公钥、客户方私钥、市场编号、商户编号等)
- [ ] 配置加载功能正常工作
- [ ] 配置验证功能正常工作
- [ ] 支持默认配置和自定义配置

- [ ] CcbPaymentInfo结构体定义正确,包含所有必需字段
- [ ] CcbSubOrder结构体定义正确,包含所有必需字段
- [ ] CcbLedgerAccount结构体定义正确,包含所有必需字段
- [ ] CcbCoupon结构体定义正确,包含所有必需字段
- [ ] PlaceOrderRequest和PlaceOrderResponse结构体定义正确
- [ ] RefundOrderRequest和RefundOrderResponse结构体定义正确
- [ ] QueryOrderRequest和QueryOrderResponse结构体定义正确
- [ ] QueryRefundRequest和QueryRefundResponse结构体定义正确
- [ ] JSON序列化和反序列化标签正确设置

- [ ] RSA密钥加载功能支持PEM格式
- [ ] SHA256withRSA签名功能正常工作
- [ ] SHA256withRSA验签功能正常工作
- [ ] 签名模块单元测试通过

- [ ] JSON对象递归排序功能正常工作
- [ ] 字段过滤功能正确排除SIGN_INF、Svc_Rsp_St、Svc_Rsp_Cd、Rsp_Inf
- [ ] 签名字符串生成功能正常工作
- [ ] JSON拼接模块单元测试通过

- [ ] HTTPClient结构体定义正确
- [ ] POST请求发送功能正常工作
- [ ] 请求头正确设置
- [ ] 超时控制正常工作
- [ ] 错误处理完善

- [ ] PlaceOrder方法实现正确
- [ ] 请求JSON构建正确(包含主订单、子订单、分账方、消费券)
- [ ] 请求签名正确
- [ ] 请求发送正常
- [ ] 响应验签正确
- [ ] 支付URL提取和URL解码正确
- [ ] 支付订单接口单元测试通过

- [ ] QueryOrder方法实现正确
- [ ] 请求JSON构建正确
- [ ] 请求签名正确
- [ ] 请求发送正常
- [ ] 响应验签正确
- [ ] 查询支付结果接口单元测试通过

- [ ] RefundOrder方法实现正确
- [ ] 请求JSON构建正确(支持全额退款和部分退款)
- [ ] 请求签名正确
- [ ] 请求发送正常
- [ ] 响应验签正确
- [ ] 订单退款接口单元测试通过

- [ ] QueryRefund方法实现正确
- [ ] 请求JSON构建正确
- [ ] 请求签名正确
- [ ] 请求发送正常
- [ ] 响应验签正确
- [ ] 查询退款结果接口单元测试通过

- [ ] Client结构体定义正确
- [ ] NewClient构造函数正常工作
- [ ] 各模块集成正确
- [ ] API方法便捷易用(PlaceOrder, QueryOrder, RefundOrder, QueryRefund)

- [ ] 支付订单示例可以正常运行
- [ ] 查询支付结果示例可以正常运行
- [ ] 订单退款示例可以正常运行
- [ ] 查询退款结果示例可以正常运行
- [ ] 配置初始化示例清晰
- [ ] 错误处理示例完整
- [ ] 示例代码注释详细

- [ ] 包级别文档完整
- [ ] API文档符合GoDoc格式
- [ ] README.md包含使用说明、快速开始、API参考
- [ ] 集成指南完整

- [ ] 签名模块单元测试通过
- [ ] JSON拼接模块单元测试通过
- [ ] HTTP客户端单元测试通过
- [ ] 支付订单接口单元测试通过
- [ ] 查询支付结果接口单元测试通过
- [ ] 订单退款接口单元测试通过
- [ ] 查询退款结果接口单元测试通过

- [ ] 代码通过go fmt格式化
- [ ] 代码通过go vet检查
- [ ] 代码通过golint检查
- [ ] 所有导出函数都有文档注释
- [ ] 代码符合项目规范(coding-style.md, concurrency-performance.md, documentation-comments.md)
