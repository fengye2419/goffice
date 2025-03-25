# GOfficer - Markdown 转 DOCX 工具

这是一个将 Markdown 文件转换为 DOCX 文档的命令行工具，支持标题、段落、加粗文本和数学公式等元素。

## 功能特性

- 支持 Markdown 基本语法
- 支持数学公式（LaTeX 格式）
- 生成标准 DOCX 文件

## 使用方法

```bash
go run cmd/goffice/main.go 输入文件.md 输出文件.docx
```

或者构建后使用：

```bash
go build -o goffice cmd/goffice/main.go
./goffice 输入文件.md 输出文件.docx
```

## 项目结构

```
goffice/
├── cmd/            # 命令行入口
├── internal/       # 内部包
│   ├── models/     # 数据模型
│   ├── parser/     # Markdown 解析器
│   └── docx/       # DOCX 生成器
└── pkg/            # 公共包
    ├── latex/      # LaTeX 处理
    └── utils/      # 通用工具
```

## 开发

```bash
# 获取依赖
go mod tidy

# 运行测试
go test ./...
```