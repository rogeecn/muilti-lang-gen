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
# 在当前目录生成多语言文件
./multilang-gen gen .

# 在指定目录生成多语言文件
./multilang-gen gen /path/to/project

# 自定义输出文件名格式
./multilang-gen gen . --output "page-{lang}.html"

# 只生成指定语言的文件
./multilang-gen gen . --lang zh
./multilang-gen gen . --lang zh,en
./multilang-gen gen . --lang zh --lang en

# 初始化项目模板
./multilang-gen init ./my-project
```

### 项目目录结构

使用 gen 命令前，项目目录应包含以下结构：

```text
project/
├── index.tmpl              # HTML 模板文件（必需）
├── manifest.json           # 站点配置（可选）
├── langs/                  # 语言文件目录（必需）
│   ├── index.json          # 语言索引（必需）
│   ├── zh-CN.json          # 中文语言包
│   └── en-US.json          # 英文语言包
└── outputs/                # 输出目录（自动创建）
    ├── zh.html             # 中文页面
    └── en.html             # 英文页面
```

│ ├── langs/ # 测试语言文件
│ └── output/ # 测试输出目录
├── Makefile # 构建脚本
├── README.md # 项目说明
├── requirements.md # 需求文档
└── main.go # 主程序入口

````

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
````

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

模板使用 Go 的 `html/template` 语法，可以访问以下数据结构：

### 可用变量

- `{{.Lang}}` - 当前语言信息对象
  - `{{.Lang.Code}}` - 语言代码（如：zh, en, fr）
  - `{{.Lang.Name}}` - 语言名称（如：中文, English）
  - `{{.Lang.DisplayName}}` - 显示名称（用于链接文本）
  - `{{.Lang.URL}}` - 当前语言的文件 URL
- `{{.LangLinks}}` - 其他语言链接列表
- `{{.I18N}}` - 当前语言的国际化数据内容（从 JSON 文件加载）

### 语言链接渲染

在模板中使用 `range` 循环来渲染语言切换链接：

```html
<div class="lang-links">
  <strong>{{.I18N.switch_language}}:</strong>
  {{range .LangLinks}}
  <a href="{{.URL}}">{{.DisplayName}}</a>
  {{end}}
</div>
```

### 高级语言链接渲染

你可以根据需要自定义链接的渲染样式：

```html
<!-- 简单链接列表 -->
<div class="languages">
  {{range .LangLinks}}
  <a href="{{.URL}}" class="lang-link">{{.DisplayName}}</a>
  {{end}}
</div>

<!-- 下拉菜单样式 -->
<select onchange="location = this.value;">
  <option value="{{.Lang.URL}}">{{.Lang.DisplayName}} (当前)</option>
  {{range .LangLinks}}
  <option value="{{.URL}}">{{.DisplayName}}</option>
  {{end}}
</select>

<!-- 带图标的链接 -->
<ul class="lang-menu">
  {{range .LangLinks}}
  <li>
    <a href="{{.URL}}" title="切换到{{.DisplayName}}">
      <span class="flag flag-{{.Code}}"></span>
      {{.DisplayName}}
    </a>
  </li>
  {{end}}
</ul>
```

### 条件渲染

可以使用条件语句来控制显示逻辑：

```html
<!-- 只有多种语言时才显示切换链接 -->
{{if .LangLinks}}
<div class="lang-switcher">
  <span>{{.I18N.switch_language}}:</span>
  {{range .LangLinks}}
  <a href="{{.URL}}">{{.DisplayName}}</a>
  {{end}}
</div>
{{end}}

<!-- 根据语言代码显示不同内容 -->
{{if eq .Lang.Code "zh"}}
<p>这是中文特有的内容</p>
{{else if eq .Lang.Code "en"}}
<p>This is English-specific content</p>
{{end}}
```

### 模板示例

```html
<!DOCTYPE html>
<html lang="{{.Lang.Code}}">
  <head>
    <title>{{.I18N.title}} - {{.Lang.DisplayName}}</title>
  </head>
  <body>
    <h1>{{.I18N.title}}</h1>
    <p>{{.I18N.description}}</p>

    <div class="lang-links">
      <strong>{{.I18N.switch_language}}:</strong>
      {{range .LangLinks}}
      <a href="{{.URL}}">{{.DisplayName}}</a>
      {{end}}
    </div>

    <main>
      <p>{{.I18N.content}}</p>
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
