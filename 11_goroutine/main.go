// ============================================================
// 11 - Goroutine 与同步原语（WaitGroup / Mutex / Once / atomic）
//
// 运行：go run ./11_goroutine
// 检测数据竞争：go run -race ./11_goroutine
// ============================================================

package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 【什么是 goroutine】
// goroutine 是 Go 运行时管理的"轻量级线程"：
//   - 初始栈仅 2KB（线程约 1MB），可以轻松开几十万个
//   - 由 Go 调度器调度，内核感知不到（用户态调度）
//   - 通过 go 关键字启动，函数调用前加 go 即可

func main() {
	// ---------- 1. 启动 goroutine ----------
	go say("世界", 3) // 新开一个协程并发执行
	say("你好", 1)     // 当前主 goroutine 执行

	// 【经典新手问题】main 函数退出时，所有 goroutine 会被直接杀掉。
	// 不做同步的话，上面的输出可能还没打印程序就结束了。
	// （time.Sleep 是"等一下"的土办法，仅用于演示，别在生产代码用！）
	time.Sleep(100 * time.Millisecond)

	// ---------- 2. WaitGroup：等待一组 goroutine 完成（标准做法）----------
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1) // 计数器 +1（必须在 go 之前调用）
		go func(id int) {
			defer wg.Done() // 完成时计数器 -1
			fmt.Printf("  worker %d 开始工作\n", id)
			time.Sleep(time.Duration(rand(id)) * time.Millisecond) // 模拟耗时
			fmt.Printf("  worker %d 完成\n", id)
		}(i) // 【注意】把 i 作为参数传入，而不是直接闭包引用！
	}
	wg.Wait() // 阻塞直到计数器归零
	fmt.Println("所有 worker 完成")

	// ---------- 3. 数据竞争（race condition）演示 ----------
	// 多个 goroutine 同时写同一个变量 => 结果不可预测！
	counter := 0
	var wg2 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			counter++ // 这行不是原子操作（读-改-写三步）
		}()
	}
	wg2.Wait()
	fmt.Println("无锁计数（可能小于1000）:", counter) // 用 -race 跑会直接报告竞争

	// ---------- 4. Mutex：互斥锁，保护共享数据 ----------
	var (
		mu      sync.Mutex
		safeCnt int
		wg3     sync.WaitGroup
	)
	for i := 0; i < 1000; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			mu.Lock()         // 加锁：其他 goroutine 会阻塞在这里
			safeCnt++         // 临界区：同一时刻只有一个 goroutine 在执行
			mu.Unlock()       // 解锁（习惯上和 Lock 成对，或用 defer）
		}()
	}
	wg3.Wait()
	fmt.Println("加锁计数（一定是1000）:", safeCnt)

	// RWMutex：读多写少的场景，读锁可并行（RLock），写锁独占（Lock）
	// var rw sync.RWMutex; rw.RLock(); ...; rw.RUnlock()

	// ---------- 5. atomic：简单计数的最快方案 ----------
	var atomicCnt atomic.Int64 // Go 1.19+ 的类型化原子变量
	var wg4 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg4.Add(1)
		go func() {
			defer wg4.Done()
			atomicCnt.Add(1) // 原子加
		}()
	}
	wg4.Wait()
	fmt.Println("原子计数:", atomicCnt.Load())

	// ---------- 6. sync.Once：只执行一次（单例初始化等）----------
	var once sync.Once
	for i := 0; i < 3; i++ {
		once.Do(func() { // 只有第一次调用会真正执行函数体
			fmt.Println("这行只会打印一次")
		})
	}

	// ---------- 7. 主 goroutine 与调度 ----------
	// GOMAXPROCS：默认等于 CPU 核数，控制并行度
	// runtime.Gosched()：主动让出 CPU（很少用）
	fmt.Println("CPU 核数:", runtime.GOMAXPROCS(0), "当前 goroutine 数:", runtime.NumGoroutine())

	// 【哲学】Go 名言：
	//   "不要通过共享内存来通信，而要通过通信来共享内存"
	// 下一课的 channel 就是这种"以通信为核心"的并发风格。
}

func say(msg string, times int) {
	for i := 0; i < times; i++ {
		fmt.Println(msg)
	}
}

// rand 简单伪随机（避免引入完整示例的复杂性）
func rand(n int) int {
	return (n*7919 + 13) % 50
}

// ============================================================
// 课后练习：
// 1. 起 10 个 goroutine 并发打印 1~10，用 WaitGroup 等待全部完成。
// 2. 用 -race 模式分别跑"无锁计数"和"加锁计数"，观察检测报告。
// 3. 用 sync.Once 实现一个懒加载的单例配置对象。
// ============================================================
