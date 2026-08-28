// ============================================================
// 12 - Channel 通道：goroutine 之间的通信管道 + select
//
// 运行：go run ./12_channel
// ============================================================

package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// ---------- 1. 无缓冲 channel：同步通信 ----------
	// 创建：ch := make(chan T)。发送和接收必须双方同时就位（ rendezvous ），
	// 否则阻塞 —— 天然的同步机制。
	ch := make(chan string)
	go func() {
		ch <- "来自 goroutine 的消息" // 发送（在主协程接收前会阻塞）
	}()
	msg := <-ch // 接收（阻塞直到有数据）
	fmt.Println("收到:", msg)

	// ---------- 2. 有缓冲 channel：异步通信 ----------
	bufCh := make(chan int, 3) // 缓冲区大小 3
	bufCh <- 1                 // 不阻塞（缓冲区没满）
	bufCh <- 2
	bufCh <- 3
	// bufCh <- 4 // 取消注释会死锁：缓冲区已满，没有接收方
	fmt.Println("缓冲区:", len(bufCh), "/", cap(bufCh))
	fmt.Println(<-bufCh, <-bufCh, <-bufCh) // 取出 1 2 3

	// ---------- 3. close 与 range ----------
	// 发送方 close 后：接收方能读完剩余数据，之后得到零值/ok=false。
	// 【原则】由发送方关闭，不要在接收方关闭！
	naturals := make(chan int)
	go func() {
		for i := 1; i <= 5; i++ {
			naturals <- i
		}
		close(naturals) // 发送完毕后关闭
	}()
	for n := range naturals { // range 会自动读到 channel 关闭为止
		fmt.Print(n, " ")
	}
	fmt.Println()

	// comma-ok 语法判断 channel 是否关闭
	v, ok := <-naturals
	fmt.Println("关闭后读取:", v, ok) // 0 false

	// ---------- 4. 单向 channel：限制方向，增强安全 ----------
	// 类型 chan<- T 只能发送，<-chan T 只能接收。
	// 常用于函数参数：明确"这个函数只会发/只会收"，防止误用。
	counter := counterGen() // counterGen 内部持有发送端，外部只拿到接收端
	for i := 0; i < 3; i++ {
		fmt.Println("计数:", counter())
	}

	// ---------- 5. select：同时等待多个 channel ----------
	// 哪个 channel 就绪就执行哪个 case（多个就绪则随机选）
	selectDemo()

	// ---------- 6. 超时控制：select + time.After ----------
	slowCh := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond) // 模拟慢操作
		slowCh <- "慢响应"
	}()
	select {
	case r := <-slowCh:
		fmt.Println("收到:", r)
	case <-time.After(50 * time.Millisecond): // 50ms 超时通道
		fmt.Println("超时了！不等了")
	}

	// ---------- 7. select + default：非阻塞收发 ----------
	ch2 := make(chan int, 1)
	select {
	case v := <-ch2:
		fmt.Println("读到:", v)
	default: // 没有 case 就绪时走 default（不阻塞）
		fmt.Println("暂时没数据，先干别的")
	}

	// ---------- 8. done channel：通知 goroutine 退出 ----------
	done := make(chan struct{}) // 空结构体：不占内存的信号 channel
	go func() {
		for {
			select {
			case <-done: // 收到退出信号
				fmt.Println("后台任务退出")
				return
			default:
				time.Sleep(10 * time.Millisecond) // 模拟工作
			}
		}
	}()
	time.Sleep(50 * time.Millisecond)
	close(done) // 广播退出信号（所有接收方都能收到）
	time.Sleep(20 * time.Millisecond)

	// ---------- 9. context：更优雅的取消机制（工程标准做法）----------
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("工作完成")
	case <-ctx.Done(): // 超时或手动 cancel 都会触发
		fmt.Println("context 取消:", ctx.Err()) // context deadline exceeded
	}

	// ---------- 10. worker pool 模式（并发编程经典范式）----------
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// 3 个 worker 消费同一个 jobs 队列
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs) // 告诉 worker：没有新任务了
	for a := 1; a <= 9; a++ {
		fmt.Print(<-results, " ") // 收集 9 个结果
	}
	fmt.Println("\nworker pool 完成")

	// ---------- channel 要点总结 ----------
	// 1) nil channel：收发都会永久阻塞（可用于禁用 select 的某个分支）
	// 2) 向已关闭的 channel 发送会 panic；关闭 nil channel 也会 panic
	// 3) 无缓冲=同步(接力棒)；有缓冲=异步(传送带)
	// 4) goroutine 泄漏：channel 没人收/没人发，goroutine 永远卡住
}

// counterGen 返回一个"只接收"channel 的闭包，产生递增序列
func counterGen() func() int {
	ch := make(chan int)
	go func() {
		for i := 1; ; i++ {
			ch <- i // 每次调用 counter() 才发送下一个
		}
	}()
	return func() int { return <-ch }
}

// selectDemo 演示 select 的多路复用
func selectDemo() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() { time.Sleep(10 * time.Millisecond); c1 <- "one" }()
	go func() { time.Sleep(20 * time.Millisecond); c2 <- "two" }()

	// 两次 select 分别接收先到的消息
	for i := 0; i < 2; i++ {
		select {
		case m1 := <-c1:
			fmt.Println("select 收到:", m1)
		case m2 := <-c2:
			fmt.Println("select 收到:", m2)
		}
	}
}

// worker 从 jobs 取任务，计算平方后放入 results
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		time.Sleep(10 * time.Millisecond) // 模拟处理耗时
		results <- j * j
	}
}

// ============================================================
// 课后练习：
// 1. 用 goroutine + channel 生成斐波那契数列的前 10 个数。
// 2. 实现并发求和：把 1~1000 分给 4 个 goroutine，channel 汇总结果。
// 3. 给 worker pool 加一个 ctx 参数，支持随时取消所有 worker。
// ============================================================
