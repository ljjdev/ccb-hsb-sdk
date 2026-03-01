# 模块名称变更计划

## 目标
将 Go 模块名称从 `ccb-hsb-sdk` 更改为 `github.com/ljjdev/ccb-hsb-sdk`，并更新所有代码中的引用。

## 背景
为了符合 Go 模块命名规范并准备将项目发布到 GitHub，需要将本地模块名称更改为完整的 GitHub 路径。

## 实施步骤

### 步骤 1: 更新 go.mod 文件
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/go.mod`
- **操作**: 将 `module ccb-hsb-sdk` 改为 `module github.com/ljjdev/ccb-hsb-sdk`
- **影响**: 这是模块的根声明，所有其他文件都依赖于此

### 步骤 2: 更新示例文件的导入语句
需要更新以下文件中的 import 语句：

#### 2.1 examples/place_order.go
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/examples/place_order.go`
- **操作**: 将所有 `ccb-hsb-sdk/` 开头的导入改为 `github.com/ljjdev/ccb-hsb-sdk/`
- **具体修改**:
  - `"ccb-hsb-sdk/internal/utils"` → `"github.com/ljjdev/ccb-hsb-sdk/internal/utils"`
  - `"ccb-hsb-sdk/pkg/client"` → `"github.com/ljjdev/ccb-hsb-sdk/pkg/client"`
  - `"ccb-hsb-sdk/pkg/config"` → `"github.com/ljjdev/ccb-hsb-sdk/pkg/config"`
  - `"ccb-hsb-sdk/pkg/model"` → `"github.com/ljjdev/ccb-hsb-sdk/pkg/model"`

#### 2.2 examples/query_order.go
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/examples/query_order.go`
- **操作**: 更新所有 `ccb-hsb-sdk/` 导入为 `github.com/ljjdev/ccb-hsb-sdk/`

#### 2.3 examples/refund_order.go
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/examples/refund_order.go`
- **操作**: 更新所有 `ccb-hsb-sdk/` 导入为 `github.com/ljjdev/ccb-hsb-sdk/`

#### 2.4 examples/query_refund.go
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/examples/query_refund.go`
- **操作**: 更新所有 `ccb-hsb-sdk/` 导入为 `github.com/ljjdev/ccb-hsb-sdk/`

### 步骤 3: 更新测试文件的导入语句

#### 3.1 pkg/client/client_test.go
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/pkg/client/client_test.go`
- **操作**: 更新所有 `ccb-hsb-sdk/` 导入为 `github.com/ljjdev/ccb-hsb-sdk/`

#### 3.2 pkg/client/placeorder_test.go
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/pkg/client/placeorder_test.go`
- **操作**: 更新所有 `ccb-hsb-sdk/` 导入为 `github.com/ljjdev/ccb-hsb-sdk/`

#### 3.3 pkg/client/queryorder_test.go
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/pkg/client/queryorder_test.go`
- **操作**: 更新所有 `ccb-hsb-sdk/` 导入为 `github.com/ljjdev/ccb-hsb-sdk/`

#### 3.4 pkg/client/refundorder_test.go
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/pkg/client/refundorder_test.go`
- **操作**: 更新所有 `ccb-hsb-sdk/` 导入为 `github.com/ljjdev/ccb-hsb-sdk/`

#### 3.5 pkg/client/queryrefund_test.go
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/pkg/client/queryrefund_test.go`
- **操作**: 更新所有 `ccb-hsb-sdk/` 导入为 `github.com/ljjdev/ccb-hsb-sdk/`

### 步骤 4: 更新主包文件
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/pkg/client/client.go`
- **操作**: 检查并更新所有 `ccb-hsb-sdk/` 导入为 `github.com/ljjdev/ccb-hsb-sdk/`

### 步骤 5: 更新文档文件

#### 5.1 README.md
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/README.md`
- **操作**: 更新以下内容:
  - 安装命令: `go get ccb-hsb-sdk` → `go get github.com/ljjdev/ccb-hsb-sdk`
  - 所有代码示例中的导入语句
  - 文档中提到的模块名称引用

#### 5.2 docs/integration-guide.md
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/docs/integration-guide.md`
- **操作**: 更新所有代码示例中的导入语句和模块引用

#### 5.3 docs/queryrefund-implementation.md
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/docs/queryrefund-implementation.md`
- **操作**: 更新所有代码示例中的导入语句

#### 5.4 docs/QUERY_ORDER_IMPLEMENTATION.md
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/docs/QUERY_ORDER_IMPLEMENTATION.md`
- **操作**: 更新所有代码示例中的导入语句

### 步骤 6: 更新 Makefile
- **文件**: `/Users/jasonlin/project/ccb-hsb-sdk/Makefile`
- **操作**: 检查并更新所有模块路径引用

### 步骤 7: 清理和验证依赖
- **操作**: 运行 `go mod tidy` 清理依赖
- **目的**: 确保所有依赖关系正确更新

### 步骤 8: 代码格式化
- **操作**: 运行 `go fmt ./...` 格式化所有代码
- **目的**: 确保代码风格一致

### 步骤 9: 静态检查
- **操作**: 运行 `go vet ./...` 进行静态检查
- **目的**: 确保没有编译错误或潜在问题

### 步骤 10: 运行测试
- **操作**: 运行 `go test ./... -v -timeout 30s` 运行所有测试
- **目的**: 验证所有功能正常工作，模块引用正确

## 影响范围

### 需要修改的文件总数: 16 个
- 1 个 go.mod 文件
- 4 个示例文件
- 5 个测试文件
- 1 个主包文件
- 4 个文档文件
- 1 个 Makefile

### 修改类型
- **模块声明**: 1 处 (go.mod)
- **导入语句**: 约 40+ 处 (分布在多个 Go 文件中)
- **文档内容**: 约 20+ 处 (在 README 和文档文件中)

## 验证标准

变更完成后，需要满足以下标准：

1. ✅ `go.mod` 中的模块名称已更新为 `github.com/ljjdev/ccb-hsb-sdk`
2. ✅ 所有 Go 文件中的导入语句已更新
3. ✅ 所有文档文件中的示例代码已更新
4. ✅ `go mod tidy` 运行成功，无错误
5. ✅ `go fmt ./...` 运行成功，无格式问题
6. ✅ `go vet ./...` 运行成功，无警告或错误
7. ✅ `go test ./...` 运行成功，所有测试通过
8. ✅ 项目可以正常编译和运行

## 注意事项

1. **向后兼容性**: 此变更会破坏向后兼容性，使用旧模块路径的代码需要更新导入语句
2. **版本控制**: 建议在 Git 中创建新分支进行此变更，确认无误后再合并
3. **依赖更新**: 如果有其他项目依赖此模块，它们需要更新 go.mod 中的 require 语句
4. **发布**: 更新完成后，需要将代码推送到 GitHub 并打标签发布新版本

## 预期结果

完成后，项目将使用新的模块路径 `github.com/ljjdev/ccb-hsb-sdk`，所有代码引用和文档都已更新，项目可以正常编译、测试和运行。
