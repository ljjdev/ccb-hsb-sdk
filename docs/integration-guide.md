# 对公专业结算综合服务平台 SDK 集成指南

本文档提供了对公专业结算综合服务平台 SDK 的详细集成指南,帮助开发者快速完成 SDK 的集成和部署。

## 目录

- [环境准备](#环境准备)
- [SDK 安装](#sdk-安装)
- [配置说明](#配置说明)
- [快速开始](#快速开始)
- [API 详细说明](#api-详细说明)
- [最佳实践](#最佳实践)
- [错误处理](#错误处理)
- [测试](#测试)
- [部署](#部署)
- [常见问题](#常见问题)

## 环境准备

### 系统要求

- **操作系统**: Linux、macOS、Windows
- **Go 版本**: Go 1.26 或更高版本
- **网络**: 需要能够访问建行支付平台网关

### 前置条件

在开始集成之前,需要准备以下信息:

1. **市场编号 (MarketID)**: 由建行提供,用于标识市场方
2. **商家编号 (MerchantID)**: 由建行提供,用于标识商家
3. **网关地址 (GatewayURL)**: 支付平台的接口网关地址
4. **RSA 密钥对**:
   - 商户私钥: 用于请求签名
   - 银行公钥: 用于响应验签

### 生成 RSA 密钥对

如果建行没有提供密钥对,可以使用 OpenSSL 生成:

```bash
# 生成 2048 位的 RSA 私钥
openssl genrsa -out private_key.pem 2048

# 从私钥提取公钥
openssl rsa -in private_key.pem -pubout -out public_key.pem

# 查看私钥内容(可选)
cat private_key.pem

# 查看公钥内容(可选)
cat public_key.pem
```

**注意**:
- 私钥文件必须妥善保管,不要泄露给第三方
- 公钥需要提供给建行进行配置
- 建议使用 2048 位或更长的密钥长度

## SDK 安装

### 使用 go get 安装

```bash
go get github.com/ljjdev/ccb-hsb-sdk
```

### 验证安装

```bash
go mod tidy
go mod verify
```

## 配置说明

SDK 提供了两种配置方式:

### 1. 代码配置

使用 `config.NewConfig()` 和配置选项函数进行配置:

```go
package main

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "log"
    "os"

    "ccb-hsb-sdk/pkg/config"
)
// loadPrivateKey 将没有头尾的 Base64 字符串解析为 RSA 私钥
func loadPrivateKey(rawStr string) (*rsa.PrivateKey, error) {
   der, err := base64.StdEncoding.DecodeString(rawStr)
   if err != nil {
      return nil, fmt.Errorf("base64 decode error: %v", err)
   }

   // 优先尝试 PKCS#8 格式
   if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
      if priv, ok := key.(*rsa.PrivateKey); ok {
         return priv, nil
      }
   }

   // 备选尝试 PKCS#1 格式
   return x509.ParsePKCS1PrivateKey(der)
}

// loadPublicKey 将没有头尾的 Base64 字符串解析为 RSA 公钥
func loadPublicKey(rawStr string) (*rsa.PublicKey, error) {
   der, err := base64.StdEncoding.DecodeString(rawStr)
   if err != nil {
      return nil, fmt.Errorf("base64 decode error: %v", err)
   }

   // 优先尝试 PKIX 格式 (最常用)
   if pub, err := x509.ParsePKIXPublicKey(der); err == nil {
      if rsaPub, ok := pub.(*rsa.PublicKey); ok {
         return rsaPub, nil
      }
   }

   // 备选尝试 PKCS#1 格式
   return x509.ParsePKCS1PublicKey(der)
}
func main() {
    // 加载密钥
    privateKey, err := loadPrivateKey("1234567")
    if err != nil {
        log.Fatal(err)
    }

    publicKey, err := loadPublicKey("3454545")
    if err != nil {
        log.Fatal(err)
    }

    // 创建配置
    cfg, err := config.NewConfig(
        config.WithMarketID("12345678901234"),
        config.WithMerchantID("12345678901234567890"),
        config.WithGatewayURL("https://marketpay.ccb.com/online/direct"),
        config.WithPrivateKey(privateKey),
        config.WithPublicKey(publicKey),
        config.WithTimeout(30),
        config.WithDebug(false),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 使用配置创建客户端
    // ...
}
```

### 2. 环境变量配置

使用环境变量进行配置,适合容器化部署:

```bash
# 设置环境变量
export CCB_MARKET_ID="12345678901234"
export CCB_MERCHANT_ID="12345678901234567890"
export CCB_GATEWAY_URL="https://marketpay.ccb.com/online/direct"
export CCB_PRIVATE_KEY_PATH="/path/to/private_key.pem"
export CCB_PUBLIC_KEY_PATH="/path/to/public_key.pem"
export CCB_TIMEOUT="30"
export CCB_DEBUG="false"
```

```go
package main

import (
    "log"

    "github.com/ljjdev/ccb-hsb-sdk/pkg/config"
    "github.com/ljjdev/ccb-hsb-sdk/pkg/client"
)

func main() {
    // 从环境变量加载配置
    cfg, err := config.LoadConfigFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // 创建客户端
    cli, err := client.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // 使用客户端
    // ...
}
```

**注意**: `LoadConfigFromEnv()` 目前仅支持从环境变量读取配置项,密钥加载功能需要自行实现。

### 配置项说明

| 配置项 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| MarketID | string | 是 | 市场编号,由银行提供 |
| MerchantID | string | 是 | 商家编号,由银行提供 |
| GatewayURL | string | 是 | 接口网关地址 |
| PrivateKey | *rsa.PrivateKey | 是 | 商户私钥,用于签名 |
| PublicKey | *rsa.PublicKey | 是 | 银行公钥,用于验签 |
| Timeout | time.Duration | 否 | HTTP 请求超时时间,默认 30 秒 |
| Debug | bool | 否 | 是否开启调试模式,默认 false |

## 快速开始

### 完整的支付流程示例

以下是一个完整的支付流程示例,包括创建订单、查询订单、退款等操作:

```go
package main

import (
    "context"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "log"
    "os"
    "time"

    "ccb-hsb-sdk/internal/utils"
    "ccb-hsb-sdk/pkg/client"
    "ccb-hsb-sdk/pkg/config"
    "ccb-hsb-sdk/pkg/model"
)

// 全局客户端实例
var ccbClient *client.Client
// loadPrivateKey 将没有头尾的 Base64 字符串解析为 RSA 私钥
func loadPrivateKey(rawStr string) (*rsa.PrivateKey, error) {
   der, err := base64.StdEncoding.DecodeString(rawStr)
   if err != nil {
      return nil, fmt.Errorf("base64 decode error: %v", err)
   }

   // 优先尝试 PKCS#8 格式
   if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
      if priv, ok := key.(*rsa.PrivateKey); ok {
         return priv, nil
      }
   }

   // 备选尝试 PKCS#1 格式
   return x509.ParsePKCS1PrivateKey(der)
}

// loadPublicKey 将没有头尾的 Base64 字符串解析为 RSA 公钥
func loadPublicKey(rawStr string) (*rsa.PublicKey, error) {
   der, err := base64.StdEncoding.DecodeString(rawStr)
   if err != nil {
      return nil, fmt.Errorf("base64 decode error: %v", err)
   }

   // 优先尝试 PKIX 格式 (最常用)
   if pub, err := x509.ParsePKIXPublicKey(der); err == nil {
      if rsaPub, ok := pub.(*rsa.PublicKey); ok {
         return rsaPub, nil
      }
   }

   // 备选尝试 PKCS#1 格式
   return x509.ParsePKCS1PublicKey(der)
}
func init() {
    // 初始化 SDK 客户端
    privateKey, err := loadPrivateKey("1234567")
    if err != nil {
        log.Fatalf("加载私钥失败: %v", err)
    }

    publicKey, err := loadPublicKey("1234567")
    if err != nil {
        log.Fatalf("加载公钥失败: %v", err)
    }

    cfg, err := config.NewConfig(
        config.WithMarketID("12345678901234"),
        config.WithMerchantID("12345678901234567890"),
        config.WithGatewayURL("https://marketpay.ccb.com/online/direct"),
        config.WithPrivateKey(privateKey),
        config.WithPublicKey(publicKey),
        config.WithTimeout(30*time.Second),
        config.WithDebug(true),
    )
    if err != nil {
        log.Fatalf("创建配置失败: %v", err)
    }

    ccbClient, err = client.NewClient(cfg)
    if err != nil {
        log.Fatalf("创建客户端失败: %v", err)
    }
}

func main() {
    // 1. 创建支付订单
    mainOrderNo := utils.GenerateSerialNumber("ORD")
    payURL, err := createOrder(mainOrderNo, "100.00")
    if err != nil {
        log.Fatalf("创建订单失败: %v", err)
    }

    fmt.Printf("支付URL: %s\n", payURL)

    // 2. 查询订单状态
    // 在实际业务中,应该通过回调或定时任务查询订单状态
    resp, err := queryOrder(mainOrderNo)
    if err != nil {
        log.Fatalf("查询订单失败: %v", err)
    }

    if resp.IsPaid() {
        fmt.Println("订单已支付成功")

        // 3. 发起退款
        refundOrderNo := utils.GenerateSerialNumber("REF")
        err = refundOrder(mainOrderNo, refundOrderNo, "100.00")
        if err != nil {
            log.Fatalf("退款失败: %v", err)
        }

        // 4. 查询退款结果
        err = queryRefund(refundOrderNo)
        if err != nil {
            log.Fatalf("查询退款失败: %v", err)
        }
    }
}

// 创建支付订单
func createOrder(mainOrderNo, amount string) (string, error) {
    req := &model.CreateOrderRequest{
        MainOrdrNo:   mainOrderNo,
        PymdCd:       model.PaymentMethodMobileH5,
        PyOrdrTpcd:   model.OrderTypeNormal,
        OrdrTamt:     amount,
        TxnTamt:      amount,
        PayDsc:       "商品购买",
        OrderTimeOut: "1800",
        Orderlist: []model.SubOrder{
            {
                CmdtyOrdrNo: mainOrderNo + "01",
                OrdrAmt:     amount,
                Txnamt:      amount,
                CmdtyDsc:    "商品",
            },
        },
    }

    payURL, err := ccbClient.PlaceOrder(context.Background(), req)
    if err != nil {
        return "", fmt.Errorf("创建订单失败: %w", err)
    }

    return payURL, nil
}

// 查询订单状态
func queryOrder(mainOrderNo string) (*model.QueryOrderResponse, error) {
    req := &model.QueryOrderRequest{
        MainOrdrNo: mainOrderNo,
    }

    resp, err := ccbClient.QueryOrder(context.Background(), req)
    if err != nil {
        return nil, fmt.Errorf("查询订单失败: %w", err)
    }

    if !resp.IsSuccess() {
        return nil, fmt.Errorf("查询订单失败: %w", resp.GetError())
    }

    return resp, nil
}

// 订单退款
func refundOrder(mainOrderNo, refundOrderNo, refundAmount string) error {
    req := &model.RefundOrderRequest{
        MainOrdrNo:   mainOrderNo,
        RefundOrdrNo: refundOrderNo,
        RefundAmt:    refundAmount,
        RefundRsn:    "用户申请退款",
    }

    resp, err := ccbClient.RefundOrder(context.Background(), req)
    if err != nil {
        return fmt.Errorf("退款失败: %w", err)
    }

    if !resp.IsSuccess() {
        return fmt.Errorf("退款失败: %w", resp.GetError())
    }

    fmt.Printf("退款成功,流水号: %s\n", resp.RefundTrnNo)
    return nil
}

// 查询退款结果
func queryRefund(refundOrderNo string) error {
    req := &model.QueryRefundRequest{
        CustRfndTrcno: refundOrderNo,
    }

    resp, err := ccbClient.QueryRefund(context.Background(), req)
    if err != nil {
        return fmt.Errorf("查询退款失败: %w", err)
    }

    if !resp.IsSuccess() {
        return fmt.Errorf("查询退款失败: %w", resp.GetError())
    }

    fmt.Printf("退款成功,金额: %s 元\n", resp.RfndAmt)
    return nil
}

```

## API 详细说明

### 1. 支付订单生成 (PlaceOrder)

创建支付订单并返回支付 URL。

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| MainOrdrNo | string | 是 | 主订单编号,不允许重复 |
| PymdCd | PaymentMethod | 是 | 支付方式代码 |
| PyOrdrTpcd | OrderType | 是 | 订单类型 |
| OrdrTamt | string | 是 | 订单总金额 |
| TxnTamt | string | 是 | 交易总金额 |
| Orderlist | []SubOrder | 是 | 子订单列表 |
| PayDsc | string | 否 | 支付描述 |
| OrderTimeOut | string | 否 | 订单超时时间(秒) |
| SubAppid | string | 否 | 小程序的APPID |
| SubOpenid | string | 否 | 用户子标识 |

**返回值**:

- `payURL`: 支付 URL,用于引导用户完成支付

**示例**:

```go
req := &model.CreateOrderRequest{
    MainOrdrNo:   utils.GenerateSerialNumber("ORD"),
    PymdCd:       model.PaymentMethodMobileH5,
    PyOrdrTpcd:   model.OrderTypeNormal,
    OrdrTamt:     "100.00",
    TxnTamt:      "100.00",
    Orderlist: []model.SubOrder{
        {
            CmdtyOrdrNo: "SUB20240101120000123",
            OrdrAmt:     "100.00",
            Txnamt:      "100.00",
        },
    },
}

payURL, err := cli.PlaceOrder(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

// 将 payURL 返回给前端,引导用户完成支付
```

### 2. 支付结果查询 (QueryOrder)

查询订单的支付状态。

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| MainOrdrNo | string | 条件必填 | 主订单编号,与 PyTrnNo 必输其一 |
| PyTrnNo | string | 条件必填 | 支付流水号,与 MainOrdrNo 必输其一 |

**返回值**:

- `QueryOrderResponse`: 查询响应,包含订单状态、支付金额等信息

**订单状态**:

| 状态码 | 说明 |
|--------|------|
| 1 | 待支付 |
| 2 | 支付成功 |
| 3 | 支付失败 |
| 9 | 待轮询 |

**示例**:

```go
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

### 3. 订单退款 (RefundOrder)

发起订单退款,支持全额退款和部分退款。

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| MainOrdrNo | string | 是 | 主订单编号 |
| RefundOrdrNo | string | 是 | 退款订单编号,不允许重复 |
| RefundAmt | string | 是 | 退款金额 |
| RefundRsn | string | 否 | 退款原因 |

**返回值**:

- `RefundOrderResponse`: 退款响应,包含退款流水号、退款时间等信息

**示例**:

```go
req := &model.RefundOrderRequest{
    MainOrdrNo:   "ORD20240101120000123",
    RefundOrdrNo: utils.GenerateSerialNumber("REF"),
    RefundAmt:    "100.00",
    RefundRsn:    "用户申请退款",
}

resp, err := cli.RefundOrder(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

if resp.IsSuccess() {
    log.Printf("退款成功,流水号: %s\n", resp.RefundTrnNo)
}
```

### 4. 退款结果查询 (QueryRefund)

查询退款结果。

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| CustRfndTrcno | string | 条件必填 | 客户方退款流水号,与 RfndTrcno 必输其一 |
| RfndTrcno | string | 条件必填 | 退款流水号,与 CustRfndTrcno 必输其一 |

**返回值**:

- `QueryRefundResponse`: 查询响应,包含退款状态、退款金额等信息

**退款状态**:

| 状态码 | 说明 |
|--------|------|
| 00 | 退款成功 |
| 01 | 退款失败 |
| 02 | 退款延迟等待 |
| 03 | 退款结果不确定 |
| 04 | 等待确认(线下订单类型返回) |
| 05 | 没有查询到符合条件的记录 |
| 0a | 已受理(仅异步退款有此状态) |
| 0b | 中断(仅异步退款有此状态) |

**示例**:

```go
req := &model.QueryRefundRequest{
    CustRfndTrcno: "REFUND20240101120000123",
}

resp, err := cli.QueryRefund(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

if resp.IsSuccess() {
    log.Printf("退款成功,金额: %s 元\n", resp.RfndAmt)
}
```

## 最佳实践

### 1. 客户端初始化

建议在应用启动时创建一个全局的 `Client` 实例,并在整个应用生命周期中复用:

```go
package main

import (
    "log"

    "ccb-hsb-sdk/pkg/client"
    "ccb-hsb-sdk/pkg/config"
)

var ccbClient *client.Client

func init() {
    cfg, err := config.NewConfig(...)
    if err != nil {
        log.Fatal(err)
    }

    ccbClient, err = client.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    // 使用全局客户端实例
    // ...
}
```

### 2. 订单编号生成

建议使用 SDK 提供的工具函数生成唯一的订单编号:

```go
import "ccb-hsb-sdk/internal/utils"

// 生成主订单编号
mainOrderNo := utils.GenerateSerialNumber("ORD")

// 生成子订单编号
subOrderNo := mainOrderNo + "01"

// 生成退款订单编号
refundOrderNo := utils.GenerateSerialNumber("REF")
```

### 3. 订单状态轮询

对于需要轮询订单状态的场景,建议使用定时任务或协程:

```go
func pollOrderStatus(mainOrderNo string, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        resp, err := ccbClient.QueryOrder(context.Background(), &model.QueryOrderRequest{
            MainOrdrNo: mainOrderNo,
        })
        if err != nil {
            log.Printf("查询订单失败: %v\n", err)
            time.Sleep(5 * time.Second)
            continue
        }

        if resp.IsPaid() {
            log.Println("订单已支付成功")
            return nil
        }

        if resp.OrdrStcd == model.OrderStatusFailed {
            return fmt.Errorf("订单支付失败")
        }

        log.Printf("订单状态: %s,继续轮询...\n", resp.OrdrStcd)
        time.Sleep(5 * time.Second)
    }

    return fmt.Errorf("订单查询超时")
}
```

### 4. 错误处理

建议对每个 API 调用进行错误处理:

```go
resp, err := cli.PlaceOrder(context.Background(), req)
if err != nil {
    // 处理网络错误、签名错误等
    log.Printf("创建订单失败: %v\n", err)
    return err
}

// 检查业务响应
if !resp.IsSuccess() {
    // 处理业务错误
    log.Printf("订单创建失败: %v\n", resp.GetError())
    return resp.GetError()
}
```

### 5. 日志记录

建议记录所有 API 调用的请求和响应,便于问题排查:

```go
func createOrderWithLogging(req *model.CreateOrderRequest) (string, error) {
    log.Printf("创建订单请求: MainOrdrNo=%s, OrdrTamt=%s\n", req.MainOrdrNo, req.OrdrTamt)

    payURL, err := ccbClient.PlaceOrder(context.Background(), req)
    if err != nil {
        log.Printf("创建订单失败: MainOrdrNo=%s, error=%v\n", req.MainOrdrNo, err)
        return "", err
    }

    log.Printf("创建订单成功: MainOrdrNo=%s, payURL=%s\n", req.MainOrdrNo, payURL)
    return payURL, nil
}
```

### 6. 金额处理

建议使用字符串类型处理金额,避免浮点数精度问题:

```go
// 正确: 使用字符串
amount := "100.00"

// 错误: 使用浮点数
amount := 100.00
```

### 7. 超时设置

根据业务需求设置合理的超时时间:

```go
cfg, err := config.NewConfig(
    config.WithTimeout(30 * time.Second), // 30秒超时
    // ...
)
```

### 8. 调试模式

在开发和测试环境开启调试模式,生产环境关闭:

```go
debug := os.Getenv("ENV") == "development"

cfg, err := config.NewConfig(
    config.WithDebug(debug),
    // ...
)
```

## 错误处理

### 错误类型

SDK 可能返回以下类型的错误:

1. **配置错误**: 私钥/公钥无效、必填参数缺失等
2. **网络错误**: 连接超时、网络不可达等
3. **签名错误**: 签名生成失败、签名验证失败等
4. **业务错误**: 订单创建失败、退款失败等

### 错误处理示例

```go
resp, err := cli.PlaceOrder(context.Background(), req)
if err != nil {
    // 判断错误类型
    if strings.Contains(err.Error(), "timeout") {
        // 处理超时错误
        log.Println("请求超时,请重试")
        return
    }

    if strings.Contains(err.Error(), "signature") {
        // 处理签名错误
        log.Println("签名错误,请检查密钥配置")
        return
    }

    // 其他错误
    log.Printf("创建订单失败: %v\n", err)
    return
}

// 检查业务响应
if !resp.IsSuccess() {
    log.Printf("订单创建失败: %s\n", resp.GetError())
    return
}
```

## 测试

### 单元测试

SDK 提供了完整的单元测试,可以运行以下命令:

```bash
# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./pkg/client

# 运行测试并显示覆盖率
go test -cover ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 集成测试

在集成测试中,可以使用建行提供的测试环境进行测试:

```go
func TestIntegration(t *testing.T) {
    // 使用测试环境的配置
    cfg, err := config.NewConfig(
        config.WithMarketID("test_market_id"),
        config.WithMerchantID("test_merchant_id"),
        config.WithGatewayURL("https://test.marketpay.ccb.com/online/direct"),
        // ...
    )
    if err != nil {
        t.Fatal(err)
    }

    cli, err := client.NewClient(cfg)
    if err != nil {
        t.Fatal(err)
    }

    // 执行测试
    // ...
}
```

## 部署

### 环境变量配置

在部署时,建议使用环境变量配置 SDK:

```bash
# 生产环境
export CCB_MARKET_ID="prod_market_id"
export CCB_MERCHANT_ID="prod_merchant_id"
export CCB_GATEWAY_URL="https://marketpay.ccb.com/online/direct"
export CCB_PRIVATE_KEY_PATH="/path/to/prod_private_key.pem"
export CCB_PUBLIC_KEY_PATH="/path/to/prod_public_key.pem"
export CCB_TIMEOUT="30"
export CCB_DEBUG="false"
```

### Docker 部署

使用 Docker 部署时,可以通过环境变量传递配置:

```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY private_key.pem .
COPY public_key.pem .
ENV CCB_MARKET_ID="your_market_id"
ENV CCB_MERCHANT_ID="your_merchant_id"
ENV CCB_GATEWAY_URL="https://marketpay.ccb.com/online/direct"
ENV CCB_PRIVATE_KEY_PATH="/app/private_key.pem"
ENV CCB_PUBLIC_KEY_PATH="/app/public_key.pem"
ENV CCB_TIMEOUT="30"
ENV CCB_DEBUG="false"
CMD ["./main"]
```

### Kubernetes 部署

在 Kubernetes 中部署时,可以使用 ConfigMap 和 Secret 管理配置:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ccb-sdk-config
data:
  CCB_MARKET_ID: "your_market_id"
  CCB_MERCHANT_ID: "your_merchant_id"
  CCB_GATEWAY_URL: "https://marketpay.ccb.com/online/direct"
  CCB_TIMEOUT: "30"
  CCB_DEBUG: "false"
---
apiVersion: v1
kind: Secret
metadata:
  name: ccb-sdk-secret
type: Opaque
stringData:
  private_key.pem: |
    -----BEGIN RSA PRIVATE KEY-----
    ...
    -----END RSA PRIVATE KEY-----
  public_key.pem: |
    -----BEGIN PUBLIC KEY-----
    ...
    -----END PUBLIC KEY-----
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ccb-sdk-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ccb-sdk
  template:
    metadata:
      labels:
        app: ccb-sdk
    spec:
      containers:
      - name: app
        image: your-image:latest
        envFrom:
        - configMapRef:
            name: ccb-sdk-config
        volumeMounts:
        - name: keys
          mountPath: /app/keys
          readOnly: true
      volumes:
      - name: keys
        secret:
          secretName: ccb-sdk-secret
```

## 常见问题

### Q1: 签名验证失败怎么办?

A: 请检查以下几点:
1. 确认使用的私钥和公钥是否正确
2. 确认公钥是否已提供给建行进行配置
3. 开启调试模式,查看详细的签名信息
4. 检查请求参数是否符合建行规范

### Q2: 如何处理网络超时?

A: 可以通过以下方式处理:
1. 增加超时时间: `config.WithTimeout(60 * time.Second)`
2. 实现重试机制
3. 检查网络连接是否正常

### Q3: 如何处理订单号重复?

A: 建议使用 SDK 提供的工具函数生成唯一的订单编号:

```go
mainOrderNo := utils.GenerateSerialNumber("ORD")
```

该函数生成的订单编号格式为: 前缀 + 时间戳 + 随机数,可以保证唯一性。

### Q4: 如何处理部分退款?

A: 可以多次调用退款接口进行部分退款,但累计退款金额不能超过订单总金额。建议在退款前查询订单的已退款金额。

### Q5: 如何实现订单回调?

A: 建行支付平台支持订单回调,需要在创建订单时设置回调地址。SDK 目前不直接支持回调处理,需要自行实现回调接口。

### Q6: 如何处理并发请求?

A: SDK 的 `Client` 实例是线程安全的,可以在多个 goroutine 中并发使用。建议创建一个全局的 `Client` 实例,并在整个应用生命周期中复用。

### Q7: 如何查看详细的请求和响应信息?

A: 开启调试模式:

```go
cfg, err := config.NewConfig(
    config.WithDebug(true),
    // ...
)
```

开启调试模式后,SDK 会输出请求和响应的详细信息。

### Q8: 如何处理金额精度问题?

A: 建议使用字符串类型处理金额,避免浮点数精度问题:

```go
// 正确: 使用字符串
amount := "100.00"

// 错误: 使用浮点数
amount := 100.00
```

### Q9: 如何处理订单超时?

A: 在创建订单时,可以通过 `OrderTimeOut` 字段设置订单超时时间(单位:秒):

```go
req := &model.CreateOrderRequest{
    OrderTimeOut: "1800", // 30分钟
    // ...
}
```

### Q10: 如何联系技术支持?

A: 如果遇到问题,可以通过以下方式获取帮助:
- 查看 [README.md](../README.md) 获取快速开始指南
- 查看 examples 目录下的示例代码
- 联系建行客服或客户经理

## 附录

### A. 支付方式代码

| 代码 | 说明 |
|------|------|
| 01 | PC端收银台 |
| 02 | 线下支付 |
| 03 | 移动端H5页面 |
| 05 | 微信小程序 |
| 06 | 对私网银 |
| 07 | 聚合二维码 |
| 08 | 龙支付 |
| 09 | 被扫 |
| 11 | 数字电子钱包 |
| 12 | 无感支付 |
| 13 | 共享钱包 |
| 14 | 支付宝小程序 |
| 15 | 免密支付 |

### B. 订单类型代码

| 代码 | 说明 |
|------|------|
| 02 | 消费券购买订单 |
| 03 | 在途订单 |
| 04 | 普通订单 |

### C. 订单状态代码

| 代码 | 说明 |
|------|------|
| 1 | 待支付 |
| 2 | 支付成功 |
| 3 | 支付失败 |
| 9 | 待轮询 |

### D. 退款状态代码

| 代码 | 说明 |
|------|------|
| 00 | 退款成功 |
| 01 | 退款失败 |
| 02 | 退款延迟等待 |
| 03 | 退款结果不确定 |
| 04 | 等待确认(线下订单类型返回) |
| 05 | 没有查询到符合条件的记录 |
| 0a | 已受理(仅异步退款有此状态) |
| 0b | 中断(仅异步退款有此状态) |

### E. 时间格式

SDK 使用以下时间格式:

- `yyyyMMddHHmmssfff`: 完整时间戳,包含毫秒
- `yyyyMMdd`: 短日期格式
- `yyyyMMddHHmmss`: 长时间格式,不包含毫秒

示例:
- `20240101120000123`: 2024年1月1日12:00:00.123
- `20240101`: 2024年1月1日
- `20240101120000`: 2024年1月1日12:00:00

### F. 金额格式

金额使用字符串类型,保留两位小数:

- 正确: `"100.00"`
- 错误: `"100"`
- 错误: `"100.000"`

### G. 参考资源

- 建行对公专业结算综合服务平台技术对接接口文档
- [Go 官方文档](https://golang.org/doc/)
- [Go 标准库文档](https://pkg.go.dev/std)
