# 修改 BuildSignatureString 函数排除更多参数

## 任务描述

修改 `/Users/jasonlin/project/ccb-hsb-sdk/pkg/signature/rsa.go` 文件中的 `BuildSignatureString` 函数，在现有排除 `Sign_Inf` 参数的基础上，增加排除 `Svc_Rsp_St`、`Svc_Rsp_Cd`、`Rsp_Inf` 这三个参数。

## 修改内容

### 文件路径
`/Users/jasonlin/project/ccb-hsb-sdk/pkg/signature/rsa.go`

### 具体修改步骤

1. **修改注释（第 115-116 行）**
   - 原注释：
   ```go
   // 3. 忽略值为空的参数
   // 4. 忽略 Sign_Inf 参数
   ```
   - 修改为：
   ```go
   // 3. 忽略值为空的参数
   // 4. 忽略 Sign_Inf、Svc_Rsp_St、Svc_Rsp_Cd、Rsp_Inf 参数
   ```

2. **修改条件判断（第 125 行）**
   - 原代码：
   ```go
   if k != "Sign_Inf" && params[k] != "" {
   ```
   - 修改为：
   ```go
   if k != "Sign_Inf" && k != "Svc_Rsp_St" && k != "Svc_Rsp_Cd" && k != "Rsp_Inf" && params[k] != "" {
   ```

## 修改后的完整代码片段

```go
// 3. 忽略值为空的参数
// 4. 忽略 Sign_Inf、Svc_Rsp_St、Svc_Rsp_Cd、Rsp_Inf 参数
func BuildSignatureString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	// 提取并排序参数键
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "Sign_Inf" && k != "Svc_Rsp_St" && k != "Svc_Rsp_Cd" && k != "Rsp_Inf" && params[k] != "" {
			keys = append(keys, k)
		}
	}

	// 按字典序排序
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	// 拼接参数
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(k)
		builder.WriteByte('=')
		builder.WriteString(params[k])
	}

	return builder.String()
}
```

## 验证步骤

修改完成后，执行以下验证步骤：

1. 运行代码格式化：
   ```bash
   go fmt ./pkg/signature/...
   ```

2. 运行静态检查：
   ```bash
   go vet ./pkg/signature/...
   ```

3. 运行单元测试：
   ```bash
   go test ./pkg/signature/... -v
   ```

4. 运行 client 包测试（包含签名验证测试）：
   ```bash
   go test ./pkg/client/... -v
   ```

5. 如果以上检查都通过，运行全量测试确保没有破坏其他功能：
   ```bash
   go test ./... -v -timeout 30s
   ```

## 预期结果

- 代码格式化通过
- 静态检查通过
- 单元测试通过
- BuildSignatureString 函数现在会排除以下四个参数：
  - `Sign_Inf`（签名信息）
  - `Svc_Rsp_St`（服务响应状态）
  - `Svc_Rsp_Cd`（服务响应码）
  - `Rsp_Inf`（响应信息）
- 签名字符串中不再包含这四个参数
- 签名验证成功

## 相关文件

- `/Users/jasonlin/project/ccb-hsb-sdk/pkg/signature/rsa.go` - 需要修改的文件
- `/Users/jasonlin/project/ccb-hsb-sdk/pkg/signature/rsa_test.go` - 可能需要更新的测试文件
