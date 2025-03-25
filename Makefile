# GOfficer - Markdown to DOCX 转换工具 Makefile

# 设置变量
BINARY_NAME=goffice
BUILD_DIR=build
MAIN_DIR=cmd/goffice

# Go命令
GO=go
GOBUILD=$(GO) build
GOTEST=$(GO) test
GOMOD=$(GO) mod
GOGET=$(GO) get
GOCLEAN=$(GO) clean
GOVET=$(GO) vet
GOFMT=$(GO) fmt

# 默认目标
.PHONY: all
all: clean test build

# 构建应用
.PHONY: build
build:
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_DIR)/main.go

# 运行测试
.PHONY: test
test:
	@echo "运行所有测试..."
	$(GOTEST) -v ./...

# 运行特定包的测试
.PHONY: test-pkg
test-pkg:
	@echo "运行 $(PKG) 包的测试..."
	$(GOTEST) -v ./$(PKG)/...

# 代码格式化
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	$(GOFMT) ./...

# 代码检查
.PHONY: vet
vet:
	@echo "代码静态检查..."
	$(GOVET) ./...

# 依赖管理
.PHONY: deps
deps:
	@echo "更新依赖..."
	$(GOMOD) tidy

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)

# 安装应用
.PHONY: install
install: build
	@echo "安装 $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# 示例：将示例Markdown转换为DOCX
.PHONY: example
example: build
	@echo "运行示例转换..."
	$(BUILD_DIR)/$(BINARY_NAME) example.md example.docx

# 帮助信息
.PHONY: help
help:
	@echo "可用的目标:"
	@echo "  all        - 清理、测试并构建应用"
	@echo "  build      - 构建应用"
	@echo "  test       - 运行所有测试"
	@echo "  test-pkg   - 运行特定包的测试 (使用: make test-pkg PKG=internal/parser)"
	@echo "  fmt        - 格式化代码"
	@echo "  vet        - 代码静态检查"
	@echo "  deps       - 更新依赖"
	@echo "  clean      - 清理构建文件"
	@echo "  install    - 安装应用到系统" 
	@echo "  example    - 运行示例转换" 