# 多语言文件生成器需求文档

## 项目概述

本项目用于实现一个多语言文件生成器，能够根据指定模板和语言文件生成对应的多语言页面。

## 功能需求

### 子命令: `gen`

#### 命令描述

根据指定模板生成多语言文件

#### 参数说明

- **参数 1**: 模板名称 (必需)
- **参数 2**: 语言列表目录 (必需)

#### 配置选项

- **输出文件名**: 支持自定义生成文件的名称
  - 默认格式: `{lang}.html`
  - `{lang}` 为语言替代符

## 实现流程

### 1. 语言文件扫描

遍历语言目录，获取所有语言文件

### 2. 模板解析

使用 `template/html` 解析模板文件

### 3. 文件生成

遍历语言列表，生成多语言文件

### 4. 链接替换

在生成的各个语言文件中，将 `{__LANG_LINKS__}` 替换为其它各个 `{lang}.html` 的链接

## 技术要求

- 使用 Go 语言的 `template/html` 包进行模板解析
- 支持动态语言链接生成
- 文件名支持自定义模板格式
- 仅支持 JSON 格式的语言文件
- 必须提供 `index.json` 语言索引文件

## 语言文件格式要求

### 索引文件 (index.json)

语言索引文件必须为 JSON 数组格式，包含所有支持的语言配置：

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
  }
]
```

### 语言数据文件

每个语言的数据文件必须为 JSON 格式，包含模板中使用的所有键值对：

```json
{
  "title": "页面标题",
  "description": "页面描述",
  "welcome": "欢迎信息",
  "content": "页面内容",
  "switch_language": "切换语言"
}
```

## 开发规范

### 测试目录规范

- 所有测试文件和测试数据必须在 `fixtures` 目录中进行
- 不得在项目根目录创建测试文件，避免污染项目结构
- 测试生成的 HTML 文件应输出到 `fixtures/output` 目录
- 测试用的模板文件应放在 `fixtures/templates` 目录
- 测试用的语言文件应放在 `fixtures/langs` 目录

### 目录结构规范

```text
fixtures/
├── templates/            # 测试模板文件
│   └── test.html        # 测试用模板
├── langs/               # 测试语言文件
│   ├── index.json       # 语言索引文件
│   ├── zh.json          # 中文语言文件
│   ├── en.json          # 英文语言文件
│   └── fr.json          # 法文语言文件
└── output/              # 测试输出目录
    ├── zh.html          # 生成的中文页面
    ├── en.html          # 生成的英文页面
    └── fr.html          # 生成的法文页面
```

## 构建和测试

### 使用 Makefile

项目提供了 Makefile 来简化常见操作：

```bash
# 构建项目
make build

# 运行功能测试
make test

# 清理生成文件
make clean

# 构建所有平台版本
make build-all

# 查看所有可用命令
make help
```

### 手动构建

```bash
# 构建
go build -o multilang-gen

# 运行
./multilang-gen gen template.html ./langs --output "{lang}.html"
```
