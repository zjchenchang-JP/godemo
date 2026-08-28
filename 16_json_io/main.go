// ============================================================
// 16 - 文件 IO 与 JSON 处理
//
// 运行：go run ./16_json_io
// ============================================================

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Task 演示 JSON 序列化的结构体
type Task struct {
	ID    int    `json:"id"`              // 字段名小写输出
	Title string `json:"title"`           //
	Done  bool   `json:"done"`            //
	Owner string `json:"owner,omitempty"` // omitempty：空字符串时该字段不输出
	// 反序列化时 tag 名要与 JSON 里的 key 对应；
	// JSON 多余的字段会被忽略，缺失的字段保持零值。
}

func main() {
	// ==================== 一、文件读写 ====================

	// 方式一（最简单）：os.ReadFile / os.WriteFile —— 一次读/写整个文件
	content := []byte("第一行：Go 文件写入演示\n第二行：bufio 更适合大文件\n")
	if err := os.WriteFile("demo.txt", content, 0644); err != nil {
		// 0644 是 Unix 权限位（Windows 下会被忽略）
		fmt.Println("写入失败:", err)
		return
	}
	data, err := os.ReadFile("demo.txt")
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}
	fmt.Printf("ReadFile 读到 %d 字节:\n%s", len(data), data)

	// 方式二：os.Open + bufio.Scanner —— 逐行处理，适合大文件
	f, err := os.Open("demo.txt")
	if err != nil {
		fmt.Println("打开失败:", err)
		return
	}
	defer f.Close() // 【惯用法】打开后立即 defer 关闭

	scanner := bufio.NewScanner(f) // Scanner 默认按行分割
	for scanner.Scan() {           // Scan 返回 false 表示读完（或出错）
		line := scanner.Text() // 去掉换行符的一行
		fmt.Println("  扫描行:", line)
	}
	if err := scanner.Err(); err != nil { // 别忘了检查扫描错误
		fmt.Println("扫描出错:", err)
	}

	// 写大文件：os.Create + bufio.Writer（缓冲写入，减少系统调用）
	fw, _ := os.Create("demo_buf.txt")
	w := bufio.NewWriter(fw)
	for i := 1; i <= 3; i++ {
		fmt.Fprintf(w, "bufio 写入第 %d 行\n", i)
	}
	w.Flush()  // 【必须】把缓冲区刷到文件
	fw.Close()
	os.Remove("demo_buf.txt") // 演示完删掉

	// ==================== 二、JSON 序列化 ====================

	tasks := []Task{
		{ID: 1, Title: "学习 Go 基础", Done: true, Owner: "me"},
		{ID: 2, Title: "写一个项目", Done: false}, // Owner 为空，会被 omitempty 省略
	}

	// Marshal：结构体 -> JSON 字节切片
	b, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("序列化失败:", err)
		return
	}
	fmt.Println("\nJSON:", string(b))

	// MarshalIndent：带缩进的美化输出
	pretty, _ := json.MarshalIndent(tasks, "", "  ")
	fmt.Println("美化版:")
	fmt.Println(string(pretty))

	// ==================== 三、JSON 反序列化 ====================

	jsonStr := `[
		{"id": 10, "title": "读入的任务", "done": false, "extra": "多余字段会被忽略"},
		{"id": 11, "title": "没有 done 字段", "done": false}
	]`
	var loaded []Task
	if err := json.Unmarshal([]byte(jsonStr), &loaded); err != nil { // 注意传指针！
		fmt.Println("反序列化失败:", err)
		return
	}
	fmt.Println("解析到", len(loaded), "个任务:", loaded[0].Title)

	// 解析未知结构的 JSON：map[string]any
	var anyData map[string]any
	json.Unmarshal([]byte(`{"name":"Go","year":2009,"tags":["简单","并发"]}`), &anyData)
	// 【注意】数字会被解析成 float64！
	fmt.Printf("name=%v year=%v(%T)\n", anyData["name"], anyData["year"], anyData["year"])

	// ==================== 四、流式编解码 ====================

	// json.NewDecoder：从 io.Reader 读（常用于 HTTP 响应、网络流）
	var fromReader []Task
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.Decode(&fromReader)
	fmt.Println("Decoder 读到:", len(fromReader), "个任务")

	// json.NewEncoder：写到 io.Writer（一个函数搞定编码+写入）
	fmt.Print("Encoder 输出: ")
	json.NewEncoder(os.Stdout).Encode(tasks[0]) // 直接编码到标准输出

	// 清理演示文件
	os.Remove("demo.txt")
	fmt.Println("\n演示文件已清理")
}

// ============================================================
// 课后练习：
// 1. 把 tasks 序列化写入 tasks.json，再读回来反序列化打印。
// 2. 定义嵌套结构体（用户含地址），练习嵌套 JSON 的 tag。
// 3. 用 bufio.Scanner 实现一个简单的 "wc -l"（统计行数）程序。
// ============================================================
