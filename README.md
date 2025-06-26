# 屎山代码鉴定器 (fuck-u-code)

一个专为挖掘项目"屎坑"设计的代码质量分析工具，能无情揭露代码的丑陋真相，并用毫不留情的幽默语言告诉你：你的代码到底有多烂。

## 特性

- **多语言支持**: 全面分析 Go、JavaScript/TypeScript、Python、Java、C/C++ 等多种编程语言
- **屎山指数评分**: 0~100 分的质量评分系统（分数越高越臭）
- **全面质量检测**: 七大维度评估代码质量
  - 循环复杂度 - 评估代码逻辑复杂程度
  - 函数长度 - 分析函数大小和状态管理
  - 注释覆盖率 - 检查代码文档完整性
  - 错误处理 - 评估异常处理的健壮性
  - 命名规范 - 判断标识符命名质量
  - 代码重复度 - 检测冗余代码片段
  - 代码结构 - 分析嵌套深度和组织结构
- **幽默精辟点评**: 用扎心又好笑的方式提供改进建议
- **高性能分析**: 多核并行处理，支持大型项目分析
- **智能排除**: 自动排除常见依赖目录和构建文件
- **可配置报告**: 支持详细/摘要模式，多语言输出
- **优化的代码结构**: 精简无冗余，高可维护性设计
- **风格统一**: 符合 Go 最佳实践的代码风格
- **高效分析算法**: 改进的通配符匹配和文件查找算法

## 安装

### 1. 从源码安装

```bash
go install github.com/Done-0/fuck-u-code@latest
```

### 2. 从源码构建

```bash
git clone https://github.com/Done-0/fuck-u-code.git
cd fuck-u-code
go build -o fuck-u-code ./cmd/fuck-u-code
```

## 使用方法

### 基本分析

```bash
fuck-u-code analyze /path/to/your/project
# 或者 fuck-u-code /path/to/your/project
```

不指定路径时，默认分析当前目录:

```bash
fuck-u-code analyze
```

### 命令行选项

| 选项         | 简写   | 描述                               |
| ------------ | ------ | ---------------------------------- |
| `--verbose`  | `-v`   | 显示详细分析报告                   |
| `--top N`    | `-t N` | 显示问题最多的前 N 个文件 (默认 5) |
| `--issues N` | `-i N` | 每个文件显示 N 个问题 (默认 5)     |
| `--summary`  | `-s`   | 只显示总结结论，不看过程           |
| `--lang`     | `-l`   | 指定输出语言 (zh-CN, en-US)        |
| `--exclude`  | `-e`   | 排除特定文件/目录模式 (可多次使用) |

### 使用示例

```bash
# 分析并显示详细报告
fuck-u-code analyze --verbose

# 只查看最糟糕的3个文件
fuck-u-code analyze --top 3

# 英文报告
fuck-u-code analyze --lang en-US

# 只查看总结信息
fuck-u-code analyze --summary

# 排除特定文件夹
fuck-u-code analyze --exclude "**/test/**" --exclude "**/legacy/**"
```

## 高级用法

### 分析前端项目

前端项目通常包含大量依赖和生成文件，工具默认已排除以下路径：

- node_modules、bower_components
- dist、build、.next、out、.cache、.nuxt、.output
- 压缩文件 (_.min.js, _.bundle.js, \_.chunk.js)
- 静态资源文件夹 (public/assets, static/js, static/css)

### 分析后端项目

分析后端项目时，工具默认会排除以下内容：

- vendor、bin、target、obj
- 临时文件夹 (tmp, temp, logs)
- 生成文件 (generated, migrations)
- 测试数据 (testdata, test-results)

## 技术架构

- **分析器引擎**: 基于 AST 的代码分析，具有语言特定优化
- **度量指标系统**: 可扩展的代码质量评估框架
- **国际化支持**: 中英文完整支持
- **高性能设计**: 多协程并行文件处理
- **模块化结构**:
  - `analyzer`: 代码分析核心逻辑
  - `metrics`: 各项代码质量指标实现
  - `parser`: 多语言代码解析器
  - `common`: 通用工具和类型
  - `i18n`: 国际化支持
  - `report`: 分析报告生成

## 代码优化特点

- **精简高效**: 消除冗余逻辑和不必要的条件判断
- **统一命名**: 符合 Go 语言规范的命名约定
- **优化流程控制**: 使用 switch 替代复杂的 if-else 链
- **提升性能**: 改进文件查找和解析算法
- **增强可维护性**: 清晰的模块边界和职责划分
- **代码风格统一**: 遵循行业最佳实践和规范

## 许可证

本项目采用 MIT 许可证

## 安利一下

另一个开源项目 [Jank](https://github.com/Done-0/Jank)
