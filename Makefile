# 多语言文件生成器 Makefile

# 变量定义
BINARY_NAME=multilang-gen
BUILD_DIR=build
FIXTURES_DIR=fixtures
GO_FILES=$(shell find . -name "*.go" -type f)
INSTALL_PATH=/usr/local/bin

# 默认目标
.PHONY: all
all: build

# 构建
.PHONY: build
build:
	@echo "构建 $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f $(FIXTURES_DIR)/output/*.html

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "代码检查..."
	go vet ./...

# 运行测试
.PHONY: test
test: build
	@echo "运行功能测试..."
	@mkdir -p $(FIXTURES_DIR)/output
	./$(BINARY_NAME) gen $(FIXTURES_DIR)/templates/test.html $(FIXTURES_DIR)/langs --output "$(FIXTURES_DIR)/output/{lang}.html"
	@echo "测试完成，生成的文件："
	@ls -la $(FIXTURES_DIR)/output/*.html

# 安装依赖
.PHONY: deps
deps:
	@echo "安装依赖..."
	go mod tidy
	go mod download

# 开发环境设置
.PHONY: dev-setup
dev-setup: deps fmt vet

# 发布构建（跨平台）
.PHONY: build-all
build-all: clean
	@echo "构建所有平台版本..."
	@mkdir -p $(BUILD_DIR)

	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .

	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

	# macOS ARM64 (M1/M2)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

	@echo "构建完成，文件在 $(BUILD_DIR) 目录："
	@ls -la $(BUILD_DIR)/

# 示例运行
.PHONY: example
example: build
	@echo "运行示例..."
	@mkdir -p $(FIXTURES_DIR)/output
	./$(BINARY_NAME) gen $(FIXTURES_DIR)/templates/test.html $(FIXTURES_DIR)/langs --output "$(FIXTURES_DIR)/output/example-{lang}.html"
	@echo "示例完成，查看生成的文件："
	@ls -la $(FIXTURES_DIR)/output/example-*.html


# 检查安装状态
.PHONY: install-check
install-check:
	@echo "检查 $(BINARY_NAME) 安装状态..."
	@if command -v $(BINARY_NAME) >/dev/null 2>&1; then \
		echo "✅ $(BINARY_NAME) 已安装"; \
		echo "安装位置: $$(which $(BINARY_NAME))"; \
		echo "版本信息:"; \
		$(BINARY_NAME) --help | head -3; \
	else \
		echo "❌ $(BINARY_NAME) 未安装"; \
		echo "运行 'make install' 进行安装"; \
	fi

# 帮助信息
.PHONY: help
help:
	@echo "可用的 Make 目标："
	@echo "  build         - 构建项目"
	@echo "  clean         - 清理构建文件和测试输出"
	@echo "  fmt           - 格式化 Go 代码"
	@echo "  vet           - 运行 Go 代码检查"
	@echo "  test          - 运行功能测试"
	@echo "  deps          - 安装和整理依赖"
	@echo "  dev-setup     - 设置开发环境"
	@echo "  build-all     - 构建所有平台版本"
	@echo "  example       - 运行示例"
	@echo "  install       - 安装到系统 ($(INSTALL_PATH))"
	@echo "  uninstall     - 从系统卸载"
	@echo "  install-check - 检查安装状态"
	@echo "  help          - 显示此帮助信息"

# 监听文件变化并自动构建（需要安装 fswatch）
.PHONY: watch
watch:
	@echo "监听文件变化并自动构建..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make build; \
	else \
		echo "请先安装 fswatch: brew install fswatch"; \
	fi
