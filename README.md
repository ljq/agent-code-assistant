# Free AI Programming Assistant (Golang)

Agent-code-assistant is a code-assistant (LLM + Tool + Loop)

A **minimalist, auditable, and privatizable** programming assistant CLI implemented using the **Go standard library**, without relying on any third-party Golang frameworks, packages, or Agent SDKs.

⚠️ **Because this involves reading and writing local files, please verify the code before running it to ensure security and avoid any loss of local files! **

---


## Golang version
Golang： ```go 1.25.5```

## ✨ Features

* ✅ **Does not use any third-party Go libraries**

* ✅ **Configuration file driven (config.yaml)**

* ✅ **Supports tool calls (file read/write/directory listing)**

* ✅ **Can connect to large local or remote models (Ollama/DeepSeek/OpenAI style)**

* ✅ **Supports user input length limits (prevents excessively long input)**

* ✅ **Simple code structure, fully understandable, and rewriteable at any time**

---

> **The assistant is essentially just a loop**

The core process consists of only these 5 steps:

``` User input

↓ LLM generates a response (may include tool calls)

↓ Execute the tool locally

↓ Feed back the tool results to the LLM

↓ Output the final answer

```

No Agent magic, no framework black box.

---

## 📁 Project Structure

```text

.
├── main.go # Main program (single file)

├── config.yaml # Model and behavior configuration

└── README.md

```

---

## ⚙️ Configuration Instructions (config.yaml)

```yaml
model: deepseek-r1
endpoint: http://127.0.0.1:11434/api/generate
temperature: 0.2
system_prompt: You are a free programming assistant.
input_limit: 0

```

### Parameter Description

| Parameter | Description |

| --------------- | ------------------- |

| `model` | Model name |

| `endpoint` | Model HTTP interface address |

| `temperature` | Generated temperature (reserved field, currently unused) |

| `system_prompt` | System prompt |

| `input_limit` | Maximum number of characters for user input, `0` means no limit |

> ⚠️ The `config.yaml` parser is a **minimalist implementation**, only supporting

> `key: value` format (one key per line)

---

## 🚀 Quick Start

### 1️⃣ Compile and Run

```bash
go build -o agent-code-assistant main.go

./agent-code-assistant

```

### 2️⃣ Input Example

```text

> List files in the current directory

> Read the contents of main.go

> Create a file hello.txt with the content "Hello World"

```

LLM If it returns the following format:

```text
tool: read_file({"path":"main.go"})

```

The program will automatically execute the tool and continue the conversation.

---

## 🧰 Built-in Tools

Currently, there are 3 built-in tools (can be freely expanded):

| Tool Name | Function |

| ------------ | --------- |

| `read_file` | Read file content |

| `list_files` | List directory files |

| `edit_file` | Create/modify file |

### Tool Call Format (Convention)

```text
tool: tool_name({"arg":"value"})

```

Example:

```text
tool: edit_file({"path":"a.txt","old":"","new":"hello"})

```

---

## 🔐 Input Length Limit

Controlled via `config.yaml`:

```yaml
input_limit: 200

```

Behavior:

* `0` → No limit (default)

* `>0` → Input exceeding the character limit will be rejected.

> Character counting uses `rune`, which is safe for Chinese characters.

---

## 📌 Who is it suitable for?

* For those who want to **understand the essence of AI programming assistants**

* For those who want to **implement the core logic of Copilot/Claude Code** themselves

* For engineers who want to create **private/localized/auditable AI tools**

* For those dissatisfied with the "black box" nature of Agent frameworks

---

## 🧪 Non-goals (Intentionally Avoided)

* ❌ Not pursuing production-grade robustness

* ❌ Not implementing complex prompt templates

* ❌ Not having a built-in permission system/sandbox

* ❌ Not introducing Agent/Workflow frameworks

---

## 🛠 Scalability Directions (Suggestions)

* Add `max_steps` to prevent infinite loops

* Add `exec_cmd` (restricted shell)

* Support multi-tool JSON arrays

* Context persistence (session recovery)

* Transform into an MCP Server (Go)

---

## 📜 License

MIT. You are free to modify, tailor, rewrite, and use it commercially.
