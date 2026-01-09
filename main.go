package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

 /**************** Configuration Loading (配置加载) ****************/

// Config uses the simplest map to store configuration
// Does not import yaml library, only supports key: value
// (Config 使用最简单的 map 保存配置，不引入 yaml 库，仅支持 key: value)
type Config map[string]string

// loadConfig reads config.yaml
// Supports:
// - Empty lines
// - # comments
// - key: value
// (loadConfig 读取 config.yaml，支持：空行，# 注释，key: value)
func loadConfig(path string) Config {
	cfg := Config{}
	data, _ := os.ReadFile(path)

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		// Ignore empty lines and comments
		// (忽略空行和注释)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split only by the first colon
		// (仅按第一个冒号分割)
		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		cfg[key] = val
	}
	return cfg
}

 /**************** Tool Definitions (工具定义) ****************/

// ToolCall represents a tool call
// Example: tool: read_file({"path":"main.go"})
// (ToolCall 表示一次工具调用，例如：tool: read_file({"path":"main.go"}))
type ToolCall struct {
	Name string
	Args map[string]string
}

// readFile reads file content
// (读取文件内容)
func readFile(path string) map[string]string {
	b, _ := os.ReadFile(path)
	return map[string]string{
		"content": string(b),
	}
}

// listFiles lists directory files
// (列出目录文件)
func listFiles(path string) map[string]string {
	es, _ := os.ReadDir(path)
	var names []string
	for _, e := range es {
		names = append(names, e.Name())
	}
	return map[string]string{
		"files": strings.Join(names, " "),
	}
}

// editFile / create file
// If old is empty, directly write new file
// (编辑 / 创建文件，old 为空表示直接写新文件)
func editFile(path, old, nw string) map[string]string {
	var content string

	if b, err := os.ReadFile(path); err == nil {
		content = string(b)
	}

	if old != "" {
		content = strings.Replace(content, old, nw, 1)
	} else {
		content = nw
	}

	os.WriteFile(path, []byte(content), 0644)
	return map[string]string{"status": "ok"}
}

// Tool registry (LLM can only call capabilities here)
// (工具注册表（LLM 只能调用这里的能力）)
var tools = map[string]func(map[string]string) map[string]string{
	"read_file": func(a map[string]string) map[string]string {
		return readFile(a["path"])
	},
	"list_files": func(a map[string]string) map[string]string {
		return listFiles(a["path"])
	},
	"edit_file": func(a map[string]string) map[string]string {
		return editFile(a["path"], a["old"], a["new"])
	},
}

 /**************** LLM Invocation (LLM 调用) ****************/

// callLLM sends request to model endpoint
// Here using Ollama / DeepSeek style as example
// (callLLM 向模型 endpoint 发送请求，这里以 Ollama / DeepSeek 风格为例)
func callLLM(cfg Config, ctx []string) string {
	reqBody := map[string]any{
		"model":  cfg["model"],
		"prompt": strings.Join(ctx, "\n"),
		"stream": false,
	}

	j, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(
		"POST",
		cfg["endpoint"],
		bytes.NewReader(j),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "LLM request error"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Compatible with Ollama response structure
	// (兼容 Ollama 返回结构)
	var r map[string]any
	json.Unmarshal(body, &r)
	if v, ok := r["response"]; ok {
		return fmt.Sprint(v)
	}

	return string(body)
}

 /**************** Tool Parsing (工具解析) ****************/

// parseTools extracts tool calls from LLM output
// Format:
// tool: edit_file({"path":"a.txt","old":"","new":"hi"})
// (parseTools 从 LLM 输出中提取 tool 调用，格式：tool: edit_file({"path":"a.txt","old":"","new":"hi"}))
func parseTools(txt string) []ToolCall {
	var calls []ToolCall

	for _, line := range strings.Split(txt, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "tool:") {
			continue
		}

		// tool:name(args)
		// (tool:name(args))
		p := strings.SplitN(line[5:], "(", 2)
		if len(p) != 2 {
			continue
		}

		args := map[string]string{}
		json.Unmarshal(
			[]byte(strings.TrimSuffix(p[1], ")")),
			&args,
		)

		calls = append(calls, ToolCall{
			Name: strings.TrimSpace(p[0]),
			Args: args,
		})
	}
	return calls
}

 /**************** Main Loop (主循环) ****************/

func main() {
	cfg := loadConfig("config.yaml")

	// Read input length limit
	// (读取输入长度限制)
	limit, _ := strconv.Atoi(cfg["input_limit"])

	reader := bufio.NewReader(os.Stdin)

	// Dialogue context (the "state machine's memory")
	// (对话上下文（就是"状态机的内存"）)
	ctx := []string{
		"system: " + cfg["system_prompt"],
	}

	for {
		fmt.Print("> ")
		user, _ := reader.ReadString('\n')
		user = strings.TrimSpace(user)

		// Input length limit logic
		// (输入长度限制逻辑)
		if limit > 0 && len([]rune(user)) > limit {
			fmt.Printf(
				"输入过长（最大 %d 字符）\n",
				limit,
			)
			continue
		}

		ctx = append(ctx, "user: "+user)

		for {
			// LLM → decision
			// (LLM → 决策)
			out := callLLM(cfg, ctx)

			// Whether tool invocation is needed
			// (是否需要调用工具)
			calls := parseTools(out)
			if len(calls) == 0 {
				fmt.Println(out)
				ctx = append(ctx, "assistant: "+out)
				break
			}

			// Execute tools and feed results back to model
			// (执行工具，并把结果喂回模型)
			for _, c := range calls {
				if fn, ok := tools[c.Name]; ok {
					res := fn(c.Args)
					j, _ := json.Marshal(res)
					ctx = append(ctx, "tool: "+string(j))
				}
			}
		}
	}
}
