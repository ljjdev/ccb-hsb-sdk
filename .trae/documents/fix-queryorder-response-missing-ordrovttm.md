# 修复 QueryOrderResponse 缺少 OrdrOvtmTm 字段的问题

## 问题描述

在 `/Users/jasonlin/project/ccb-hsb-sdk/pkg/client/client.go` 的 `QueryOrder` 方法中，调用 `BuildSignatureString` 函数构建签名字符串时，丢失了 `Ordr_Ovtm_Tm` 字段。

## 根本原因

经过分析，发现问题出在 `QueryOrderResponse` 结构体定义中：

1. **响应示例包含该字段**：
   - 在 `/Users/jasonlin/project/ccb-hsb-sdk/doc/响应示例.md` 第 33 行，查询支付结果接口的响应包含：
   ```json
   "Ordr_Ovtm_Tm": "20260228154541"
   ```

2. **结构体缺少该字段**：
   - 在 `/Users/jasonlin/project/ccb-hsb-sdk/pkg/model/order.go` 第 514-553 行的 `QueryOrderResponse` 结构体中，**没有定义** `OrdrOvtmTm` 字段
   - 现有字段包括：`OrdrGenTm`（订单生成时间），但缺少 `OrdrOvtmTm`（订单超时时间）

3. **影响链路**：
   - JSON 反序列化时，由于结构体没有 `OrdrOvtmTm` 字段，该字段被忽略
   - 调用 `resp.ToMap()` 时，map 中不包含 `Ordr_Ovtm_Tm` 键
   - `BuildSignatureString` 构建签名字符串时，自然也就不包含该字段
   - 最终导致签名验证失败

## 解决方案

在 `QueryOrderResponse` 结构体中添加 `OrdrOvtmTm` 字段。

### 文件路径
`/Users/jasonlin/project/ccb-hsb-sdk/pkg/model/order.go`

### 具体修改步骤

1. **在 `QueryOrderResponse` 结构体中添加 `OrdrOvtmTm` 字段**
   - 位置：在 `OrdrGenTm` 字段之后（第 525 行之后）
   - 添加内容：
   ```go
   // OrdrOvtmTm 订单超时时间
   OrdrOvtmTm string `json:"Ordr_Ovtm_Tm,omitempty"`
   ```

2. **修改后的结构体顺序**（参考 `CreateOrderResponse` 的字段顺序）：
   ```go
   type QueryOrderResponse struct {
       // MainOrdrNo 主订单编号
       MainOrdrNo string `json:"Main_Ordr_No,omitempty"`

       // PyTrnNo 支付流水号
       PyTrnNo string `json:"Py_Trn_No"`

       // Txnamt 交易金额
       Txnamt float64 `json:"Txnamt,omitempty"`

       // OrdrGenTm 订单生成时间
       OrdrGenTm string `json:"Ordr_Gen_Tm,omitempty"`

       // OrdrOvtmTm 订单超时时间
       OrdrOvtmTm string `json:"Ordr_Ovtm_Tm,omitempty"`

       // OrdrStcd 订单状态代码
       OrdrStcd OrderStatus `json:"Ordr_Stcd"`

       // ... 其他字段
   }
   ```

## 验证步骤

修改完成后，执行以下验证步骤：

1. 运行代码格式化：
   ```bash
   go fmt ./pkg/model/...
   ```

2. 运行静态检查：
   ```bash
   go vet ./pkg/model/...
   ```

3. 运行单元测试：
   ```bash
   go test ./pkg/model/... -v
   ```

4. 运行 client 包测试（包含签名验证测试）：
   ```bash
   go test ./pkg/client/... -v
   ```

5. 运行全量测试确保没有破坏其他功能：
   ```bash
   go test ./... -v -timeout 30s
   ```

## 预期结果

- 代码格式化通过
- 静态检查通过
- 单元测试通过
- `QueryOrderResponse` 结构体包含 `OrdrOvtmTm` 字段
- JSON 反序列化时能正确解析 `Ordr_Ovtm_Tm` 字段
- `BuildSignatureString` 构建的签名字符串包含 `Ordr_Ovtm_Tm` 参数
- 签名验证成功

## 相关文件

- `/Users/jasonlin/project/ccb-hsb-sdk/pkg/model/order.go` - 需要修改的结构体定义文件
- `/Users/jasonlin/project/ccb-hsb-sdk/pkg/client/client.go` - 调用签名验证的文件（无需修改）
- `/Users/jasonlin/project/ccb-hsb-sdk/pkg/signature/rsa.go` - 签名字符串构建函数（无需修改）
- `/Users/jasonlin/project/ccb-hsb-sdk/doc/响应示例.md` - 响应示例参考文档
