# 修改 BuildSignatureString 函数的参数过滤逻辑

## 任务描述
修改 `/Users/jasonlin/project/ccb-hsb-sdk/pkg/signature/rsa.go` 文件中的 `BuildSignatureString` 函数，将参数过滤逻辑从忽略 "sign" 参数改为忽略 "Sign_Inf" 参数。

## 修改内容

### 文件路径
`/Users/jasonlin/project/ccb-hsb-sdk/pkg/signature/rsa.go`

### 具体修改步骤

1. **修改注释（第 115 行）**
   - 原注释：`// 4. 忽略 sign 参数`
   - 修改为：`// 4. 忽略 Sign_Inf 参数`

2. **修改条件判断（第 124 行）**
   - 原代码：`if k != "sign" && params[k] != "" {`
   - 修改为：`if k != "Sign_Inf" && params[k] != "" {`

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

4. 如果以上检查都通过，运行全量测试确保没有破坏其他功能：
   ```bash
   go test ./... -v -timeout 30s
   ```

## 预期结果

- 代码格式化通过
- 静态检查通过
- 单元测试通过
- BuildSignatureString 函数现在会忽略 "Sign_Inf" 参数而不是 "sign" 参数
