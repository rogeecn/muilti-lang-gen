# 多语言文件生成器 (Multilingual File Generator)

一个基于 Go 语言开发的多语言网站生成工具，可以根据模板和语言文件快速生成多语言网站。

## 功能特性

- ✅ 支持自定义 HTML 模板
- ✅ 基于 JSON 格式的语言索引和数据文件
- ✅ 自动生成语言间链接
- ✅ 支持自定义输出文件名格式
- ✅ 命令行界面，使用简单
- ✅ 跨平台支持（Linux、macOS、Windows）

## 快速开始

### 安装

#### 使用 Makefile 构建

```bash
git clone <repository-url>
cd multilang-gen
make build
```

#### 手动构建

```bash
git clone <repository-url>
cd multilang-gen
go build -o multilang-gen
```

### 基本使用

```bash
# 基本用法
./multilang-gen gen template.html ./langs

# 自定义输出文件名
./multilang-gen gen template.html ./langs --output "page-{lang}.html"
```

## 项目结构

```text
multilang-gen/
├── cmd/                   # 命令行相关代码
├── fixtures/              # 测试文件目录
│   ├── templates/         # 测试模板
│   ├── langs/             # 测试语言文件
│   └── output/            # 测试输出目录
├── Makefile              # 构建脚本
├── README.md             # 项目说明
├── requirements.md       # 需求文档
└── main.go               # 主程序入口
```

## 语言文件格式

### JSON 格式（必需）

项目仅支持 JSON 格式的语言文件，提供最佳的数据结构支持和易读性：

#### 语言索引文件 (index.json)

```json
[
  {
    "code": "zh",
    "name": "中文",
    "displayName": "中文",
    "file": "zh.json"
  },
  {
    "code": "en",
    "name": "English",
    "displayName": "English",
    "file": "en.json"
  },
  {
    "code": "fr",
    "name": "Français",
    "displayName": "Français",
    "file": "fr.json"
  }
]
```

#### 语言数据文件

```json
{
  "title": "页面标题",
  "description": "页面描述",
  "welcome": "欢迎信息",
  "content": "页面内容",
  "features": "功能列表标题",
  "feature1": "功能1描述",
  "feature2": "功能2描述",
  "feature3": "功能3描述",
  "switch_language": "切换语言",
  "footer": "页脚信息"
}
```

**文件命名规范**：

- 索引文件：`index.json`
- 语言文件：`{语言代码}.json`（如：`zh.json`, `en.json`, `fr.json`）

**JSON 格式要求**：

- 必须是有效的 JSON 格式
- 键名使用英文，方便模板中引用
- 值为对应语言的翻译文本
- 编码必须为 UTF-8

## 模板语法

模板使用 Go 的 `html/template` 语法，可以访问以下数据：

- `{{.Language}}` - 当前语言标识
- `{{.LangCode}}` - 当前语言代码
- `{{.LangName}}` - 当前语言显示名称
- `{{.Data.key}}` - 语言文件中的键值对数据
- `{__LANG_LINKS__}` - 其他语言链接的占位符（会被自动替换）

### 模板示例

```html
<!DOCTYPE html>
<html lang="{{.Language}}">
  <head>
    <title>{{.Data.title}} - {{.Language}}</title>
  </head>
  <body>
    <h1>{{.Data.title}}</h1>
    <p>{{.Data.description}}</p>

    <div class="lang-links">
      <strong>{{.Data.switch_language}}:</strong> {__LANG_LINKS__}
    </div>

    <main>
      <p>{{.Data.content}}</p>
    </main>
  </body>
</html>
```

## 使用 Makefile

项目提供了 Makefile 来简化常见操作：

```bash
# 构建项目
make build

# 运行功能测试
make test

# 清理生成文件
make clean

# 格式化代码
make fmt

# 代码检查
make vet

# 构建所有平台版本
make build-all

# 运行示例
make example

# 查看所有可用命令
make help
```

## 命令参数

### gen 命令

```text
Usage: multilang-gen gen [template] [language-dir] [flags]

参数:
  template      模板文件路径
  language-dir  语言文件目录路径

选项:
  -o, --output string   输出文件名模式，{lang} 为语言替代符 (默认 "{lang}.html")
  -h, --help           显示帮助信息
```

## 示例

项目在 `fixtures` 目录下包含了一个完整的示例：

1. **模板文件**: `fixtures/templates/test.html`
2. **语言索引**: `fixtures/langs/index.json`
3. **语言文件**: `fixtures/langs/` 目录下的 `zh.json`, `en.json`, `fr.json`
4. **运行命令**:

   ```bash
   make test
   # 或
   ./multilang-gen gen fixtures/templates/test.html fixtures/langs --output "fixtures/output/{lang}.html"
   ```

5. **生成结果**: `fixtures/output/` 目录下的 HTML 文件

## 开发

### 开发环境设置

```bash
# 设置开发环境（安装依赖、格式化、检查）
make dev-setup
```

### 监听文件变化（需要安装 fswatch）

```bash
# 安装 fswatch (macOS)
brew install fswatch

# 监听文件变化并自动构建
make watch
```

## 技术栈

- [Go](https://golang.org/) - 编程语言
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [html/template](https://pkg.go.dev/html/template) - 模板引擎

## 许可证

请查看 LICENSE 文件了解许可证信息。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！

### 开发规范

- 所有测试文件和测试数据必须在 `fixtures` 目录中进行
- 不得在项目根目录创建测试文件，避免污染项目结构
- 提交代码前请运行 `make dev-setup` 进行代码格式化和检查
