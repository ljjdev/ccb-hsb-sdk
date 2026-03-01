# Makefile for ccb-hsb-sdk
# 提供构建、测试、代码检查等功能

# 变量定义
BINARY_NAME=ccb-hsb-sdk
BUILD_DIR=build
GO=go
GOFLAGS=-v
LDFLAGS=-s -w

# 颜色定义
COLOR_RESET=\033[0m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m
COLOR_RED=\033[31m

# 默认目标
.PHONY: all
all: fmt vet lint test build

# 帮助信息
.PHONY: help
help:
	@echo "$(COLOR_BLUE)可用命令:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make build$(COLOR_RESET)       - 编译项目"
	@echo "  $(COLOR_GREEN)make test$(COLOR_RESET)        - 运行测试"
	@echo "  $(COLOR_GREEN)make test-coverage$(COLOR_RESET) - 运行测试并生成覆盖率报告"
	@echo "  $(COLOR_GREEN)make fmt$(COLOR_RESET)         - 格式化代码"
	@echo "  $(COLOR_GREEN)make vet$(COLOR_RESET)         - 运行 go vet 检查"
	@echo "  $(COLOR_GREEN)make lint$(COLOR_RESET)        - 运行 golangci-lint 检查"
	@echo "  $(COLOR_GREEN)make clean$(COLOR_RESET)       - 清理构建文件"
	@echo "  $(COLOR_GREEN)make deps$(COLOR_RESET)        - 下载依赖"
	@echo "  $(COLOR_GREEN)make tidy$(COLOR_RESET)        - 整理依赖"
	@echo "  $(COLOR_GREEN)make all$(COLOR_RESET)         - 执行所有检查并构建"

# 编译项目
.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	@echo "$(COLOR_BLUE)开始编译项目...$(COLOR_RESET)"
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(COLOR_GREEN)编译完成!$(COLOR_RESET)"

# 运行测试
.PHONY: test
test:
	@echo "$(COLOR_BLUE)运行测试...$(COLOR_RESET)"
	$(GO) test $(GOFLAGS) -race -coverprofile=coverage.out ./...
	@echo "$(COLOR_GREEN)测试完成!$(COLOR_RESET)"

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "$(COLOR_BLUE)运行测试并生成覆盖率报告...$(COLOR_RESET)"
	$(GO) test $(GOFLAGS) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)覆盖率报告已生成: coverage.html$(COLOR_RESET)"

# 格式化代码
.PHONY: fmt
fmt:
	@echo "$(COLOR_BLUE)格式化代码...$(COLOR_RESET)"
	$(GO) fmt ./...
	@echo "$(COLOR_GREEN)代码格式化完成!$(COLOR_RESET)"

# 运行 go vet 检查
.PHONY: vet
vet:
	@echo "$(COLOR_BLUE)运行 go vet 检查...$(COLOR_RESET)"
	$(GO) vet ./...
	@echo "$(COLOR_GREEN)go vet 检查通过!$(COLOR_RESET)"

# 运行 golangci-lint 检查
.PHONY: lint
lint:
	@echo "$(COLOR_BLUE)运行 golangci-lint 检查...$(COLOR_RESET)"
	@which golangci-lint > /dev/null || (echo "$(COLOR_YELLOW)golangci-lint 未安装,正在安装...$(COLOR_RESET)" && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin)
	golangci-lint run --timeout=5m ./...
	@echo "$(COLOR_GREEN)golangci-lint 检查通过!$(COLOR_RESET)"

# 清理构建文件
.PHONY: clean
clean:
	@echo "$(COLOR_BLUE)清理构建文件...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)清理完成!$(COLOR_RESET)"

# 下载依赖
.PHONY: deps
deps:
	@echo "$(COLOR_BLUE)下载依赖...$(COLOR_RESET)"
	$(GO) mod download
	@echo "$(COLOR_GREEN)依赖下载完成!$(COLOR_RESET)"

# 整理依赖
.PHONY: tidy
tidy:
	@echo "$(COLOR_BLUE)整理依赖...$(COLOR_RESET)"
	$(GO) mod tidy
	@echo "$(COLOR_GREEN)依赖整理完成!$(COLOR_RESET)"

# 检查代码风格
.PHONY: check-style
check-style: fmt vet lint
	@echo "$(COLOR_GREEN)代码风格检查全部通过!$(COLOR_RESET)"

# 运行示例
.PHONY: run-example
run-example:
	@echo "$(COLOR_BLUE)运行示例...$(COLOR_RESET)"
	$(GO) run examples/*.go

# 生成文档
.PHONY: docs
docs:
	@echo "$(COLOR_BLUE)生成文档...$(COLOR_RESET)"
	$(GO) doc -all ./...

# 安装到本地
.PHONY: install
install:
	@echo "$(COLOR_BLUE)安装到本地...$(COLOR_RESET)"
	$(GO) install ./...
	@echo "$(COLOR_GREEN)安装完成!$(COLOR_RESET)"
