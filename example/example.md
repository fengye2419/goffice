# GOfficer 示例文档

这是一个用于演示 GOfficer 功能的 Markdown 文档。GOfficer 可以将 Markdown 文件转换为 DOCX 格式。

## 支持的功能

### 基本格式

Markdown 支持**粗体**格式的文本。

### 数学公式

支持行内数学公式，例如：$E=mc^2$ 或 $a^2 + b^2 = c^2$。

也支持块级数学公式：

```math
\int_0^{\infty} e^{-x} dx = 1
```

以及复杂的数学表达式：

```math
\frac{d}{dx}(x^n) = nx^{n-1}
```

### 麦克斯韦方程组

麦克斯韦方程组是电磁学的基本方程组：

```math
\nabla \times \vec{E} = -\frac{\partial \vec{B}}{\partial t}
```

```math
\nabla \times \vec{B} = \mu_0 \vec{J} + \mu_0 \epsilon_0 \frac{\partial \vec{E}}{\partial t}
```

```math
\nabla \cdot \vec{E} = \frac{\rho}{\epsilon_0}
```

```math
\nabla \cdot \vec{B} = 0
```

## 使用方法

使用以下命令将此示例文件转换为 DOCX：

```
make example
```

或者直接运行：

```
goffice example.md example.docx
``` 