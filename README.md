# 建行对公专业结算综合服务平台 SDK

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

建行对公专业结算综合服务平台 SDK (ccb-hsb-sdk) 是一个用于与中国建设银行对公专业结算综合服务平台进行交互的 Go 语言 SDK。该 SDK 提供了完整的支付、退款、查询等功能,支持 RSA 签名验证,简化了与建行支付平台的集成过程。

## 项目简介

本 SDK 旨在为开发者提供一个简单、易用、安全的接口,用于对接建行对公专业结算综合服务平台。SDK 封装了复杂的签名验证、请求构建、响应解析等逻辑,让开发者可以专注于业务逻辑的实现。

### 主要特性

- **完整的 API 支持**: 支持支付订单创建、支付结果查询、订单退款、退款结果查询等核心功能
- **自动签名验证**: 内置 RSA 签名生成和验证功能,确保交易安全
- **灵活的配置**: 支持多种配置方式,包括代码配置和环境变量配置
- **完善的错误处理**: 提供详细的错误信息,便于问题排查
- **丰富的示例代码**: 提供多种使用场景的示例,帮助快速上手
- **类型安全**: 使用 Go 的类型系统,减少运行时错误
- **线程安全**: 客户端实例可安全地在多个 goroutine 中并发使用

### 支持的功能

- ✅ 支付订单生成 (PlaceOrder)
- ✅ 支付结果查询 (QueryOrder)
- ✅ 订单退款 (RefundOrder)
- ✅ 退款结果查询 (QueryRefund)
- ✅ RSA 签名和验签
- ✅ 多种支付方式支持 (PC端、移动端H5、微信小程序、支付宝小程序等)
- ✅ 子订单和分账支持
- ✅ 消费券支持

## 快速开始

### 安装

使用 go get 命令安装 SDK:

```bash
go get github.com/ljjdev/ccb-hsb-sdk
```

### 基本使用

以下是一个完整的支付订单创建示例:

```go
package main

import (
    "context"
    "crypto/rsa"
    "log"

    "github.com/ljjdev/ccb-hsb-sdk/internal/utils"
    "github.com/ljjdev/ccb-hsb-sdk/pkg/client"
    "github.com/ljjdev/ccb-hsb-sdk/pkg/config"
    "github.com/ljjdev/ccb-hsb-sdk/pkg/model"
)

func main() {
    // 1. 准备 RSA 密钥对
    // 实际使用时,应该从文件加载密钥
    privateKey := loadPrivateKey("path/to/private_key.pem")
    publicKey := loadPublicKey("path/to/public_key.pem")

    // 2. 创建配置
    cfg, err := config.NewConfig(
        config.WithMarketID("12345678901234"),          // 市场编号,由银行提供
        config.WithMerchantID("12345678901234567890"),  // 商家编号,由银行提供
        config.WithGatewayURL("https://marketpay.ccb.com/online/direct"),
        config.WithPrivateKey(privateKey),
        config.WithPublicKey(publicKey),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 3. 创建客户端
    cli, err := client.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // 4. 生成订单编号
    mainOrderNo := utils.GenerateSerialNumber("ORD")

    // 5. 构建支付订单请求
    req := &model.CreateOrderRequest{
        MainOrdrNo:   mainOrderNo,
        PymdCd:       model.PaymentMethodMobileH5,
        PyOrdrTpcd:   model.OrderTypeNormal,
        OrdrTamt:     "100.00",
        TxnTamt:      "100.00",
        Orderlist: []model.SubOrder{
            {
                CmdtyOrdrNo: mainOrderNo + "01",
                OrdrAmt:     "100.00",
                Txnamt:      "100.00",
            },
        },
    }

    // 6. 调用支付订单生成接口
    payURL, err := cli.PlaceOrder(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    // 7. 处理响应
    log.Println("支付URL:", payURL)
}

func loadPrivateKey(path string) *rsa.PrivateKey {
    // 实现从文件加载私钥的逻辑
    return nil
}

func loadPublicKey(path string) *rsa.PublicKey {
    // 实现从文件加载公钥的逻辑
    return nil
}
```

### 查询订单

```go
// 查询支付结果
req := &model.QueryOrderRequest{
    MainOrdrNo: "ORD20240101120000123",
}

resp, err := cli.QueryOrder(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

if resp.IsPaid() {
    log.Println("订单已支付成功")
    log.Printf("支付金额: %s 元\n", resp.Txnamt)
}
```

### 订单退款

```go
// 发起退款
refundReq := &model.RefundOrderRequest{
    MainOrdrNo:   "ORD20240101120000123",
    RefundOrdrNo: utils.GenerateSerialNumber("REF"),
    RefundAmt:    "100.00",
    RefundRsn:    "用户申请退款",
}

refundResp, err := cli.RefundOrder(context.Background(), refundReq)
if err != nil {
    log.Fatal(err)
}

if refundResp.IsSuccess() {
    log.Println("退款成功")
}
```

## API 参考

### 核心包

#### pkg/client

客户端包,提供与建行支付平台交互的核心功能。

**主要类型:**

- `Client`: SDK 客户端,用于调用各种 API

**主要方法:**

- `NewClient(cfg *config.Config) (*Client, error)`: 创建新的客户端实例
- `PlaceOrder(ctx context.Context, req *model.CreateOrderRequest) (string, error)`: 创建支付订单并返回支付 URL
- `CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error)`: 创建支付订单并返回完整响应
- `QueryOrder(ctx context.Context, req *model.QueryOrderRequest) (*model.QueryOrderResponse, error)`: 查询支付结果
- `RefundOrder(ctx context.Context, req *model.RefundOrderRequest) (*model.RefundOrderResponse, error)`: 订单退款
- `QueryRefund(ctx context.Context, req *model.QueryRefundRequest) (*model.QueryRefundResponse, error)`: 查询退款结果

#### pkg/config

配置包,提供 SDK 的配置管理功能。

**主要类型:**

- `Config`: SDK 配置
- `Option`: 配置选项函数类型

**主要方法:**

- `NewConfig(opts ...Option) (*Config, error)`: 创建新的配置实例
- `WithMarketID(marketID string) Option`: 设置市场编号
- `WithMerchantID(merchantID string) Option`: 设置商家编号
- `WithGatewayURL(gatewayURL string) Option`: 设置网关地址
- `WithPrivateKey(privateKey *rsa.PrivateKey) Option`: 设置私钥
- `WithPublicKey(publicKey *rsa.PublicKey) Option`: 设置公钥
- `WithTimeout(timeout time.Duration) Option`: 设置超时时间
- `WithDebug(debug bool) Option`: 设置调试模式
- `LoadConfigFromEnv() (*Config, error)`: 从环境变量加载配置

#### pkg/model

模型包,定义了请求和响应的数据结构。

**主要类型:**

- `CreateOrderRequest`: 创建订单请求
- `CreateOrderResponse`: 创建订单响应
- `QueryOrderRequest`: 查询订单请求
- `QueryOrderResponse`: 查询订单响应
- `RefundOrderRequest`: 退款订单请求
- `RefundOrderResponse`: 退款订单响应
- `QueryRefundRequest`: 查询退款请求
- `QueryRefundResponse`: 查询退款响应
- `SubOrder`: 子订单
- `Participant`: 分账方
- `Coupon`: 消费券

**常量:**

- `OrderStatus`: 订单状态
  - `OrderStatusPending`: 待支付
  - `OrderStatusSuccess`: 支付成功
  - `OrderStatusFailed`: 支付失败
  - `OrderStatusPolling`: 待轮询

- `PaymentMethod`: 支付方式
  - `PaymentMethodPC`: PC端收银台
  - `PaymentMethodOffline`: 线下支付
  - `PaymentMethodMobileH5`: 移动端H5页面
  - `PaymentMethodWechatMini`: 微信小程序
  - `PaymentMethodAlipayMini`: 支付宝小程序
  - 等等...

- `OrderType`: 订单类型
  - `OrderTypeCoupon`: 消费券购买订单
  - `OrderTypeTransit`: 在途订单
  - `OrderTypeNormal`: 普通订单

- `RefundStatus`: 退款状态
  - `RefundStatusSuccess`: 退款成功
  - `RefundStatusFailed`: 退款失败
  - `RefundStatusDelayed`: 退款延迟等待
  - 等等...

#### pkg/signature

签名包,提供 RSA 签名和验签功能。

**主要类型:**

- `RSAService`: RSA 签名服务
- `Signer`: 签名器接口
- `Verifier`: 验签器接口

**主要方法:**

- `NewRSAService(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *RSAService`: 创建 RSA 服务实例
- `Sign(data string) (string, error)`: 对数据进行签名
- `Verify(data string, signature string) error`: 验证签名
- `BuildSignatureString(params map[string]string) string`: 构建待签名字符串
- `BuildSignatureStringFromJSON(jsonStr string) (string, error)`: 从 JSON 字符串构建待签名字符串

#### internal/utils

工具包,提供内部使用的工具函数。

**主要方法:**

- `CurrentTimestamp() string`: 获取当前时间戳
- `GenerateSerialNumber(prefix string) string`: 生成流水号
- `FormatAmount(amount float64) string`: 格式化金额
- `ParseAmount(amount string) (float64, error)`: 解析金额
- 等等...

## 集成指南

详细的集成指南请参考 [docs/integration-guide.md](docs/integration-guide.md)。

### 环境要求

- Go 1.26 或更高版本
- 建行对公专业结算综合服务平台商户账号
- RSA 密钥对(由银行提供或自行生成)

### 配置说明

SDK 支持两种配置方式:

1. **代码配置**: 使用 `config.NewConfig()` 和配置选项函数
2. **环境变量配置**: 使用 `config.LoadConfigFromEnv()`

环境变量列表:

- `CCB_MARKET_ID`: 市场编号
- `CCB_MERCHANT_ID`: 商家编号
- `CCB_GATEWAY_URL`: 网关地址
- `CCB_PRIVATE_KEY_PATH`: 私钥文件路径
- `CCB_PUBLIC_KEY_PATH`: 公钥文件路径
- `CCB_TIMEOUT`: 超时时间(秒)
- `CCB_DEBUG`: 调试模式(true/false)

### 签名说明

SDK 自动处理所有请求的签名和响应的验签,无需手动处理。签名算法为 SHA256withRSA。

### 错误处理

SDK 的所有方法都返回 error 类型,建议在调用时进行错误处理:

```go
resp, err := cli.PlaceOrder(context.Background(), req)
if err != nil {
    // 处理错误
    log.Printf("创建订单失败: %v\n", err)
    return
}

if !resp.IsSuccess() {
    // 处理业务错误
    log.Printf("订单创建失败: %v\n", resp.GetError())
    return
}
```

### 测试

运行测试:

```bash
# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./pkg/client

# 运行测试并显示覆盖率
go test -cover ./...
```

## 常见问题

### Q1: 如何获取市场编号和商家编号?

A: 市场编号和商家编号由建行提供,需要联系建行客服或客户经理申请。

### Q2: 如何生成 RSA 密钥对?

A: 可以使用 OpenSSL 生成 RSA 密钥对:

```bash
# 生成私钥
openssl genrsa -out private_key.pem 2048

# 生成公钥
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

### Q3: 支持哪些支付方式?

A: SDK 支持多种支付方式,包括:
- PC端收银台 (PaymentMethodPC)
- 移动端H5页面 (PaymentMethodMobileH5)
- 微信小程序 (PaymentMethodWechatMini)
- 支付宝小程序 (PaymentMethodAlipayMini)
- 龙支付 (PaymentMethodDragonPay)
- 聚合二维码 (PaymentMethodQRCode)
- 等等...

### Q4: 如何处理订单超时?

A: 在创建订单时,可以通过 `OrderTimeOut` 字段设置订单超时时间(单位:秒):

```go
req := &model.CreateOrderRequest{
    OrderTimeOut: "1800", // 30分钟
    // ...
}
```

### Q5: 如何实现订单状态轮询?

A: 可以使用定时任务或协程定期查询订单状态:

```go
for i := 0; i < maxRetries; i++ {
    resp, err := cli.QueryOrder(context.Background(), req)
    if err != nil {
        continue
    }

    if resp.IsPaid() {
        // 订单已支付
        break
    }

    if resp.OrdrStcd == model.OrderStatusFailed {
        // 订单支付失败
        break
    }

    // 等待一段时间后重试
    time.Sleep(5 * time.Second)
}
```

### Q6: 如何开启调试模式?

A: 在创建配置时设置 `WithDebug(true)`:

```go
cfg, err := config.NewConfig(
    config.WithDebug(true),
    // ...
)
```

开启调试模式后,SDK 会输出请求和响应的详细信息,便于问题排查。

### Q7: 如何处理部分退款?

A: 可以多次调用退款接口进行部分退款,但累计退款金额不能超过订单总金额:

```go
// 第一次退款 50 元
refundReq := &model.RefundOrderRequest{
    MainOrdrNo:   "ORD20240101120000123",
    RefundOrdrNo: utils.GenerateSerialNumber("REF"),
    RefundAmt:    "50.00",
    RefundRsn:    "部分退款",
}

cli.RefundOrder(context.Background(), refundReq)

// 第二次退款 30 元
refundReq2 := &model.RefundOrderRequest{
    MainOrdrNo:   "ORD20240101120000123",
    RefundOrdrNo: utils.GenerateSerialNumber("REF"),
    RefundAmt:    "30.00",
    RefundRsn:    "部分退款",
}

cli.RefundOrder(context.Background(), refundReq2)
```

### Q8: SDK 是否线程安全?

A: 是的,`Client` 实例可以在多个 goroutine 中安全地并发使用。建议在应用启动时创建一个全局的 `Client` 实例,并在整个应用生命周期中复用。

### Q9: 如何验证签名是否正确?

A: SDK 会自动验证所有响应的签名。如果签名验证失败,会返回错误。可以通过开启调试模式查看详细的签名信息。

### Q10: 如何联系技术支持?

A: 如果遇到问题,可以通过以下方式获取帮助:
- 查看 [docs/integration-guide.md](docs/integration-guide.md) 获取详细的集成指南
- 查看 examples 目录下的示例代码
- 联系建行客服或客户经理

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request!

## 更新日志

### v1.0.0 (2026-03-01)

- 初始版本发布
- 支持支付订单创建
- 支持支付结果查询
- 支持订单退款
- 支持退款结果查询
- 支持 RSA 签名和验签
- 提供完整的示例代码
