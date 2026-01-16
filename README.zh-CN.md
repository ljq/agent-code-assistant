
# Free AI Programming Assistant (Golang)

[English](README.md)

agent-code-assistant is a code-assistant（LLM + Tool + Loop）

一个**极简、可审计、可私有化**的编程助手CLI
使用 **Go 标准库** 实现，不依赖任何第三方Golang框架或者package或 Agent SDK。

⚠️ **因为涉及本地文件的读写操作，请在运行前，校验代码确保安全性避免对本地文件造成损失！**

---

## Golang version
Golang： ```go 1.25.5```

## ✨ 特性

* ✅ **不使用任何 Go 第三方库**
* ✅ **配置文件驱动（config.yaml）**
* ✅ **支持工具调用（文件读 / 写 / 列目录）**
* ✅ **可接本地或远程大模型（Ollama / DeepSeek / OpenAI 风格）**
* ✅ **支持用户输入长度限制（防止超长输入）**
* ✅ **代码结构简单、可完全理解、可随时重写**

---

> **assistant本质上只是一个循环**

核心流程只有这 5 步：

```
用户输入
   ↓
LLM 生成回复（可能包含 tool 调用）
   ↓
本地执行工具
   ↓
将工具结果反馈给 LLM
   ↓
输出最终回答
```

没有 Agent 魔法，没有框架黑箱。

---

## 📁 项目结构

```text
.
├── main.go        # 主程序（单文件）
├── config.yaml    # 模型与行为配置
├── main_test.go   # 单元测试
└── README.md
```

---

## ⚙️ 配置说明（config.yaml）

```yaml
model: deepseek-r1
endpoint: http://127.0.0.1:11434/api/generate
temperature: 0.2
system_prompt: You are a free programming assistant.
input_limit: 0
```

### 参数说明

| 参数              | 说明                  |
| --------------- | ------------------- |
| `model`         | 模型名称                |
| `endpoint`      | 模型 HTTP 接口地址        |
| `temperature`   | 生成温度（保留字段，当前未使用）    |
| `system_prompt` | 系统提示词               |
| `input_limit`   | 用户输入最大字符数，`0` 表示不限制 |

> ⚠️ `config.yaml` 解析器为**极简实现**，仅支持
> `key: value` 形式（一行一个键）

---

## 🚀 快速开始

### 1️⃣ 编译运行

```bash
go build -o agent-code-assistant main.go
./agent-code-assistant
```

### 2️⃣ 输入示例

```text
> 列出当前目录文件
> 读取 main.go 的内容
> 创建一个 hello.txt 文件，内容是 Hello World
```

LLM 如果返回如下格式：

```text
tool: read_file({"path":"main.go"})
```

程序将自动执行工具并继续对话。

---

## 🧰 内置工具

当前内置 3 个工具（可自由扩展）：

| 工具名称          | 函数       |
| ------------ | --------- |
| `read_file`  | 读取文件内容    |
| `list_files` | 列出目录文件    |
| `edit_file`  | 创建 / 修改文件 |

### 工具调用格式（约定）

```text
tool: tool_name({"arg":"value"})
```

示例：

```text
tool: edit_file({"path":"a.txt","old":"","new":"hello"})
```

---

## 🔐 输入长度限制

通过 `config.yaml` 控制：

```yaml
input_limit: 200
```

行为：

* `0` → 不限制（默认）
* `>0` → 超过字符数直接拒绝输入

> 字符统计使用 `rune`，对中文安全

---

## 📌 适合谁？

* 想**理解 AI 编程助手本质**的人
* 想**自己实现 Copilot / Claude Code 核心逻辑**的人
* 想做 **私有化 / 本地化 / 可审计 AI 工具** 的工程师
* 对 Agent 框架“黑箱”不满意的人

---

## 🧪 非目标（刻意不做）

* ❌ 不追求生产级健壮性
* ❌ 不实现复杂 prompt 模板
* ❌ 不内置权限系统 / 沙箱
* ❌ 不引入 Agent / Workflow 框架

---

## 🛠 可扩展方向（建议）

* 添加 `max_steps` 防止死循环
* 增加 `exec_cmd`（受限 shell）
* 支持多工具 JSON 数组
* Context 持久化（会话恢复）
* 改造为 MCP Server（Go）

---

## 📜 License
MIT 可以自由修改、裁剪、重写、商用。
