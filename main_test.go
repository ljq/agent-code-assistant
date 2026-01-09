package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

/*
	==============================
	配置加载测试
	==============================
*/

func TestLoadConfig(t *testing.T) {
	content := `
model: test-model
endpoint: http://localhost
input_limit: 128
# this is a comment
`

	f, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString(content)
	f.Close()

	cfg := loadConfig(f.Name())

	if cfg["model"] != "test-model" {
		t.Fatal("model parse failed")
	}
	if cfg["endpoint"] != "http://localhost" {
		t.Fatal("endpoint parse failed")
	}
	if cfg["input_limit"] != "128" {
		t.Fatal("input_limit parse failed")
	}
}

/*
	==============================
	Tool 调用解析测试
	==============================
*/

func TestParseTools(t *testing.T) {
	text := `
hello
tool: read_file({"path":"a.txt"})
tool: edit_file({"path":"b.txt","old":"","new":"hi"})
`

	calls := parseTools(text)
	if len(calls) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(calls))
	}

	if calls[0].Name != "read_file" {
		t.Fatal("first tool name wrong")
	}
	if calls[0].Args["path"] != "a.txt" {
		t.Fatal("read_file arg error")
	}
}

/*
	==============================
	文件工具测试
	==============================
*/

func TestFileTools(t *testing.T) {
	dir := t.TempDir()

	// 创建文件
	r1 := editFile(dir+"/a.txt", "", "hello")
	if r1["status"] != "ok" {
		t.Fatal("edit_file failed")
	}

	// 读取文件
	r2 := readFile(dir + "/a.txt")
	if r2["content"] != "hello" {
		t.Fatal("read_file content mismatch")
	}

	// 列出目录
	r3 := listFiles(dir)
	if r3["files"] == "" {
		t.Fatal("list_files empty result")
	}
}

/*
	==============================
	输入长度限制测试
	==============================
*/

func TestInputLimit(t *testing.T) {
	limit := 3

	allow := func(s string) bool {
		if limit > 0 && len([]rune(s)) > limit {
			return false
		}
		return true
	}

	// 中文字符测试
	if !allow("你好") {
		t.Fatal("valid chinese input rejected")
	}
	if allow("你好啊") {
		t.Fatal("input limit not enforced")
	}
}

/*
	==============================
	LLM 调用（Mock HTTP）
	==============================
*/

func TestCallLLM(t *testing.T) {
	// 模拟模型服务
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"response":"mock-ok"}`))
		}),
	)
	defer server.Close()

	cfg := Config{
		"model":    "mock-model",
		"endpoint": server.URL,
	}

	out := callLLM(cfg, []string{"user: hello"})
	if out != "mock-ok" {
		t.Fatalf("unexpected llm output: %s", out)
	}
}
