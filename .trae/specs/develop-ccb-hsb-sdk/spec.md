# 建设银行对公专业结算综合服务平台SDK开发规范

## Why
为第三方提供便捷、安全、可靠的Go语言SDK,用于对接建设银行对公专业结算综合服务平台,简化支付、退款、分账等业务功能的集成过程。

## What Changes
- 创建完整的Go SDK项目结构
- 实现RSA签名与验签功能(SHA256withRSA算法)
- 实现JSON数据排序拼接功能
- 实现HTTP客户端封装
- 实现支付订单生成接口(gatherPlaceorder)
- 实现查询支付结果接口(gatherEnquireOrder)
- 实现订单退款接口(refundOrder)
- 实现查询退款结果接口(enquireRefundOrder)
- 实现分账功能支持
- 实现消费券功能支持
- 提供配置管理(平台公钥、客户方公钥、客户方私钥、市场编号、商户编号等)
- 提供完整的错误处理机制
- 提供单元测试和示例代码

## Impact
- Affected specs: REST API集成、RSA签名验签、JSON数据处理
- Affected code: pkg/client, pkg/signature, pkg/model, pkg/config

## ADDED Requirements

### Requirement: SDK配置管理
SDK SHALL提供配置管理功能,支持配置以下参数:
- 平台公钥(用于验签银行响应)
- 客户方私钥(用于签名请求)
- 客户方公钥
- 市场编号(Mkt_Id)
- 商户编号(MarketMerchantId)
- 接口基础URL
- 超时设置

#### Scenario: 初始化SDK配置
- **WHEN** 开发者创建SDK客户端实例
- **THEN** SDK应成功加载配置并初始化客户端

### Requirement: RSA签名与验签
SDK SHALL提供RSA签名和验签功能,使用SHA256withRSA算法。

#### Scenario: 签名请求
- **WHEN** 开发者调用签名方法
- **THEN** SDK应使用客户方私钥对数据进行SHA256withRSA签名

#### Scenario: 验签响应
- **WHEN** 开发者调用验签方法
- **THEN** SDK应使用平台公钥验证银行响应的签名

### Requirement: JSON数据排序拼接
SDK SHALL提供JSON数据排序拼接功能,按照银行要求的规则生成签名字符串。

#### Scenario: 生成签名字符串
- **WHEN** 开发者传入JSON数据
- **THEN** SDK应按照键名排序、排除特定字段(SIGN_INF、Svc_Rsp_St、Svc_Rsp_Cd、Rsp_Inf)的规则生成签名字符串

### Requirement: 支付订单生成
SDK SHALL提供支付订单生成接口,支持创建主订单和子订单。

#### Scenario: 创建支付订单
- **WHEN** 开发者调用创建支付订单方法
- **THEN** SDK应构建请求JSON、签名、发送请求、验签响应并返回支付URL或支付参数

### Requirement: 查询支付结果
SDK SHALL提供查询支付结果接口。

#### Scenario: 查询支付结果
- **WHEN** 开发者调用查询支付结果方法
- **THEN** SDK应构建请求JSON、签名、发送请求、验签响应并返回支付结果

### Requirement: 订单退款
SDK SHALL提供订单退款接口,支持全额退款和部分退款。

#### Scenario: 订单退款
- **WHEN** 开发者调用订单退款方法
- **THEN** SDK应构建请求JSON、签名、发送请求、验签响应并返回退款结果

### Requirement: 查询退款结果
SDK SHALL提供查询退款结果接口。

#### Scenario: 查询退款结果
- **WHEN** 开发者调用查询退款结果方法
- **THEN** SDK应构建请求JSON、签名、发送请求、验签响应并返回退款结果

### Requirement: 分账功能
SDK SHALL支持分账功能,允许配置分账规则和分账方列表。

#### Scenario: 配置分账
- **WHEN** 开发者创建订单时配置分账信息
- **THEN** SDK应将分账信息正确包含在请求中

### Requirement: 消费券功能
SDK SHALL支持消费券功能,允许在订单中使用消费券。

#### Scenario: 使用消费券
- **WHEN** 开发者创建订单时配置消费券信息
- **THEN** SDK应将消费券信息正确包含在请求中

### Requirement: HTTP客户端封装
SDK SHALL提供HTTP客户端封装,处理请求发送和响应接收。

#### Scenario: 发送请求
- **WHEN** SDK发送HTTP请求
- **THEN** SDK应正确设置请求头、处理超时、处理错误响应

### Requirement: 错误处理
SDK SHALL提供完善的错误处理机制,包括网络错误、签名错误、业务错误等。

#### Scenario: 处理错误
- **WHEN** 请求失败或验签失败
- **THEN** SDK应返回明确的错误信息

### Requirement: Makefile构建工具
SDK SHALL提供Makefile,用于简化项目的构建、测试、格式化等操作。

#### Scenario: 使用Makefile
- **WHEN** 开发者使用make命令
- **THEN** SDK应支持常用的make目标,如build、test、fmt、vet、lint等

## MODIFIED Requirements
无

## REMOVED Requirements
无
