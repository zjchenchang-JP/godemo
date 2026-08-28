// ============================================================
// 17 - 网络编程：HTTP 服务器与客户端
//
// 运行：go run ./17_net_http
// （服务器在本机随机端口启动，客户端请求后自动关闭，程序会正常退出）
// ============================================================

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// Note API 的数据模型
type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func main() {
	// ==================== 一、构建 HTTP 服务器 ====================

	// ServeMux = 路由器（Go 1.22+ 支持方法匹配和路径参数）
	mux := http.NewServeMux()

	// 处理函数签名固定：func(w http.ResponseWriter, r *http.Request)
	//   w：往里写响应   r：读取请求信息
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Go HTTP server!") // 最简单的文本响应
	})

	// 返回 JSON 的接口
	mux.HandleFunc("GET /api/notes", func(w http.ResponseWriter, r *http.Request) {
		notes := []Note{{1, "学习 Go"}, {2, "写 HTTP 服务"}}
		w.Header().Set("Content-Type", "application/json")
		// NewEncoder 一步完成"编码 + 写入"（也常配合 Marshal + w.Write）
		json.NewEncoder(w).Encode(notes)
	})

	// 接收 JSON 的 POST 接口
	mux.HandleFunc("POST /api/echo", func(w http.ResponseWriter, r *http.Request) {
		var in map[string]any
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			// 返回错误状态码的正确姿势
			http.Error(w, "无效的 JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(in)
	})

	// 路径参数（Go 1.22+）：/api/notes/42 中取出 id
	mux.HandleFunc("GET /api/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "你查询的笔记 ID 是 %s", id)
	})

	// ==================== 二、启动服务器（随机端口）====================
	// Listen 到 127.0.0.1:0 表示让系统分配空闲端口，演示更安全
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("监听失败:", err)
		return
	}
	baseURL := "http://" + ln.Addr().String()

	srv := &http.Server{Handler: mux}
	go func() {
		// Serve 会一直阻塞，直到 Shutdown 被调用
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			fmt.Println("服务器错误:", err)
		}
	}()
	fmt.Println("服务器已启动:", baseURL)

	// ==================== 三、HTTP 客户端 ====================

	// 自定义客户端（带超时！默认 Client 没有超时，生产必配）
	client := &http.Client{Timeout: 3 * time.Second}

	// --- GET 文本 ---
	resp, err := client.Get(baseURL + "/hello")
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	body, _ := io.ReadAll(resp.Body) // 读响应体
	resp.Body.Close()                // 【必须】响应体要关闭，否则连接泄漏
	fmt.Println("GET /hello      ->", resp.Status, "|", strings.TrimSpace(string(body)))

	// --- GET JSON，解码到结构体 ---
	resp2, _ := client.Get(baseURL + "/api/notes")
	var notes []Note
	json.NewDecoder(resp2.Body).Decode(&notes) // 流式解码，无需手动 ReadAll
	resp2.Body.Close()
	fmt.Println("GET /api/notes  ->", resp2.Status, "|", len(notes), "条笔记:", notes[1].Text)

	// --- POST JSON ---
	resp3, _ := client.Post(baseURL+"/api/echo", "application/json",
		strings.NewReader(`{"msg":"你好服务端"}`))
	var echoed map[string]any
	json.NewDecoder(resp3.Body).Decode(&echoed)
	resp3.Body.Close()
	fmt.Println("POST /api/echo  ->", resp3.Status, "|", echoed)

	// --- 路径参数 ---
	resp4, _ := client.Get(baseURL + "/api/notes/42")
	b4, _ := io.ReadAll(resp4.Body)
	resp4.Body.Close()
	fmt.Println("GET /api/notes/42 ->", resp4.Status, "|", strings.TrimSpace(string(b4)))

	// ==================== 四、优雅关闭 ====================
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	srv.Shutdown(ctx) // 等待处理中的请求完成再退出
	fmt.Println("服务器已优雅关闭")

	// ==================== 补充知识点 ====================
	// 1) 最简服务器：http.ListenAndServe(":8080", mux) 一行搞定
	// 2) net/http 已内置优雅并发：每个请求一个 goroutine，无需自己开
	// 3) 中间件 = 包装 HandlerFunc 的函数（func(h http.Handler) http.Handler），
	//    日志/鉴权/限流都靠它，可以试试自己写一个
	// 4) 生产框架推荐：gin / echo / chi（net/http 完全够用时优先标准库）
}

// ============================================================
// 课后练习：
// 1. 加一个中间件 loggingMiddleware，打印每个请求的方法、路径和耗时。
// 2. 用 HTML 模板（html/template）渲染一个简单页面。
// 3. 把这个例子改成真正的 RESTful TODO API（配合第 16 课的文件持久化）。
// ============================================================
