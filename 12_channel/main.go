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
		fmt.Println(ch)
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
	// 【原则】由发送方关闭，不要在接收方关闭！向已关闭的 channel 发送会 panic，所以只能由发送方关闭
	naturals := make(chan int)
	go func() {
		for i := 1; i <= 5; i++ {
			naturals <- i
		}
		// 注释下行代码 发生异常all goroutines are asleep - deadlock!
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
		// select 每执行一次只完成一次通信 所以循环
		for {
			select {
			case <-done: // 收到退出信号
				fmt.Println("后台任务退出")
				return
			default:
				time.Sleep(10 * time.Millisecond) // 模拟工作
				fmt.Println("工作耗时10ms")
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
		jobs <- j // 如果没有close(job) worker 永远卡在 <-jobs ** goroutine 泄漏**
	}
	// close(jobs) // 告诉 worker：没有新任务了
	for a := 1; a <= 9; a++ {
		fmt.Print(<-results, " ") // 收集 9 个结果
	}
	fmt.Println("\nworker pool 完成")

	// ---------- 11. 课后练习答案演示（实现见文件末尾）----------
	exercise1()
	exercise2()
	exercise3()

	// ---------- channel 要点总结 ----------
	// 1) nil channel：收发都会永久阻塞（可用于禁用 select 的某个分支）
	// 2) 向已关闭的 channel 发送会 panic；关闭 nil channel 也会 panic
	// 3) 无缓冲=同步(接力棒)；有缓冲=异步(传送带)
	// 4) goroutine 泄漏：channel 没人收/没人发，goroutine 永远卡住

	// 课后练习（自己的尝试——有 bug，先注释掉让程序能跑完，原因见文件末尾参考答案前的说明）
	// fibo := make(chan int,10)
	// for n := range fibo { // bug：fibo 没有任何发送方，range 在这里永久阻塞
	// 	go fibonaci(n)      //      → 最后所有 goroutine 都睡着，deadlock panic
	// }
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
	go func() { time.Sleep(2 * time.Millisecond); c2 <- "two" }()

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
	// range 的终止条件就是“channel 已close且已读空”——否则 range 永远不会结束
	for j := range jobs {
		time.Sleep(10 * time.Millisecond) // 模拟处理耗时
		results <- j * j
		fmt.Printf("worker %d 处理了任务 %d\n", id, j)
	}
	fmt.Println("worker", id, "退出") // 外部注释close(),永不执行
}

// ============================================================
// 课后练习：
// 1. 用 goroutine + channel 生成斐波那契数列的前 10 个数。
// 2. 实现并发求和：把 1~1000 分给 4 个 goroutine，channel 汇总结果。
// 3. 给 worker pool 加一个 ctx 参数，支持随时取消所有 worker。
// ============================================================
func fibonaci(n int) int{
	if n <= 2 {
		fmt.Println(1)
		return 1
	} else {
		j := fibonaci(n - 1) + fibonaci(n - 2)
		fmt.Println(j)
		return j
	}
}

// ============================================================
// 课后练习参考答案
//
// 为什么上面自己的尝试跑不起来：
//
//	fibo := make(chan int,10)
//	for n := range fibo { go fibonaci(n) }
//
// 1) fibo 从头到尾没有任何发送方（没人执行 fibo <- x）。channel 是"管道"，
//    两头都得有人，只有一头时 for range fibo 永久阻塞 → 所有 goroutine 都睡
//    着 → deadlock panic。和第 3 节注释掉 close(naturals) 是同一类问题：
//    消费方在等一个永远不会来的信号。
// 2) 思路反了：不是"往 channel 放序号、每个序号起 goroutine 去算"，
//    而是"一个 goroutine 生产数列，main 消费数列"（生产者-消费者模型）。
// 3) 递归版 fibonaci(n) 求第 n 项会从头重算，且结果 return 掉了没进 channel。
// ============================================================

// ======================= 练习 1 =======================
// 生产者-消费者模型：fibGen 生产，main 消费，无缓冲 channel 锁住两边节奏
//
// 【考点】
// 1) 无缓冲的 ch <- x 是"当面交接"：数据没有暂存处，main 没执行 <-ch 之前
//    发送方一直卡在这行（阻塞 = 挂起等待，不是报错、不烧 CPU）
// 2) close(ch) 是消费方 range 能结束的前提——不 close，range 永远等下一个数
// 3) 生产一个 → 交一个 → 消费一个 → 再生产下一个，
//    生产者永远跑不到消费者前面（这就是"无缓冲=同步(接力棒)"）
func fibGen(n int, ch chan<- int) {
	x, y := 1, 1
	for i := 0; i < n; i++ {
		ch <- x       // 无缓冲：main 没来收之前，这里会阻塞等待
		x, y = y, x+y // 经典递推：x 变成 y，y 变成 x+y
	}
	close(ch) // 发送方负责关闭，告诉消费方"没有更多了"
}

func exercise1() {
	ch := make(chan int)
	go fibGen(10, ch)
	for v := range ch { // range 自动读到 close 为止
		fmt.Print(v, " ")
	}
	fmt.Println() // 输出：0 1 1 2 3 5 8 13 21 34
}

// ======================= 练习 2 =======================
// 1~1000 切成 4 段（1~250 / 251~500 / 501~750 / 751~1000），
// 4 个 goroutine 各算一段，各自把小结果发进 channel，main 收 4 次汇总
//
// 【考点】
// 1) partSum 缓冲设为 worker 数：发送方发完就能退出，不和 main 的收取节奏耦合
// 2) start 通过参数传入而非闭包直接捕获循环变量 w——Go 1.22 之前的经典坑：
//    闭包捕获的是同一个变量，4 个 goroutine 可能全看到 w=4
//    （新版 Go 已修复，传参仍是防误读的好习惯）
// 3) 结果验证：1~1000 总和 = 1000×1001/2 = 500500
// 4) 工程里也可用 sync.WaitGroup 等 4 个都完成再汇总，本题按要求用 channel
func exercise2() {
	const (
		workers  = 4
		total    = 1000
		partSize = total / workers // 每人算 250 个
	)

	partSum := make(chan int, workers) // 缓冲 = worker 数（对照实验：改成 make(chan int) 再跑一遍试试）

	for w := 0; w < workers; w++ {
		go func(start int) {
			sum := 0
			for i := start; i < start+partSize; i++ {
				sum += i
			}
			partSum <- sum // 发送：缓冲够 → 立即放行；缓冲 0 → 卡在这等 main 来收
			// 【验证】观察这行何时打出来：
			// 缓冲 4：main 睡觉的 50ms 里 4 行全部出现（发送方不用等接收方）
			// 缓冲 0：main 睡觉的 50ms 里一行都没有（发送方全卡在上一行）
			fmt.Println("段", start, "发送完成:", sum)
		}(w*partSize + 1) // w=0 → 1~250，w=1 → 251~500，依次类推
	}

	fmt.Println("--- main 故意等 50ms 再去收，这期间观察上面有没有输出 ---")
	time.Sleep(50 * time.Millisecond)

	sum := 0
	for i := 0; i < workers; i++ {
		v := <-partSum
		fmt.Println("main 收到:", v)
		sum += v
	}
	fmt.Println("1~1000 并发总和:", sum) // 500500
}

// ======================= 练习 3 =======================
// 可取消的 worker pool：ctx 一 cancel，3 个 worker 全部立刻退出
//
// 【考点】
// 1) for range jobs 没法同时监听 ctx，必须拆成 for + select，
//    再用 comma-ok（j, ok := <-jobs）补上原来 range 自动做的"关闭检测"
// 2) 模拟耗时的 Sleep 也要用 select 包一层：Sleep 睡死了没法打断，
//    select + ctx.Done() 能在取消信号到来时立刻醒来退出
// 3) 取消是双方的事：worker 尊重 ctx 还不够，main 收结果的循环
//    若还死等 9 个会永远阻塞——收发两端都要放进 select
// 4) 「ctx 取消」和「jobs 已关闭」经常同时就绪，select 随机挑一个分支，
//    但两条路都通向 return，逻辑闭环，不用担心随机性
func ctxWorker(ctx context.Context, id int, jobs <-chan int, results chan<- int) {
	for {
		select {
		case <-ctx.Done(): // 收到取消信号 → 立即退出
			fmt.Println("worker", id, "收到取消，退出")
			return
		case j, ok := <-jobs:
			if !ok { // jobs 已关闭且读空 → 活干完了，正常退出
				fmt.Println("worker", id, "任务全部完成，退出")
				return
			}
			// 模拟耗时：不用 time.Sleep 而是 select 包住，取消时能立刻中断
			select {
			case <-ctx.Done():
				fmt.Println("worker", id, "处理任务", j, "时被取消，丢弃")
				return
			case <-time.After(10 * time.Millisecond):
			}
			results <- j * j
		}
	}
}

func exercise3() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 养成习惯：拿到 cancel 就 defer，防止 context 泄漏

	for w := 1; w <= 3; w++ {
		go ctxWorker(ctx, w, jobs, results)
	}
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	// 15ms 后取消：此时 3 个 worker 刚各做完 1 个任务（10ms/个），正在做第二轮
	go func() {
		time.Sleep(15 * time.Millisecond)
		fmt.Println("---- main 发出取消信号 ----")
		cancel()
	}()

	// 收结果的一端也必须监听 ctx：
	// 取消后 worker 不再投结果，死等 9 个会永远阻塞（死锁）
	collected := 0
	for a := 1; a <= 9; a++ {
		select {
		case r := <-results:
			fmt.Print(r, " ")
			collected++
		case <-ctx.Done():
			fmt.Printf("\n提前取消，只收到 %d/9 个结果\n", collected)
			time.Sleep(30 * time.Millisecond) // 给 worker 一点时间打印退出日志
			return
		}
	}
	fmt.Println("\nworker pool（可取消版）正常完成")
}
