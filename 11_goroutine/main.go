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

	// 【经典新手问题】main 函数退出时，所有 goroutine 会被直接杀掉。[类比守护线程]
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

	var wg5 sync.WaitGroup
	go10Print(&wg5)
	wg5.Wait()

	loopVarTrapDemo()

	// 课后练习 2 见运行说明（-race 对比），练习 3 在这里演示：
	singletonDemo()
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

// 1. 起 10 个 goroutine 并发打印 1~10，用 WaitGroup 等待全部完成。
func go10Print(wg5 *sync.WaitGroup) {
	for i := 1; i <= 10; i++ {
		wg5.Add(1)
		go func(id int) {
			defer wg5.Done()
			fmt.Printf("  goroutine %d\n", id) // 【勘误】用参数 id，别用闭包里的 i！
		}(i)
	}
	// 注意：1~10 的"数字集合"是确定的，但打印"顺序"完全随机 —— 这就是并发。
}

// 【为什么用参数 id，而不是闭包里的 i？】
// 闭包捕获的是"变量本身"（引用），不是 go 语句执行那一刻的值快照！
//
// Go 1.22 之前：整个 for 循环只有一个共享的 i。goroutine 真正被调度执行时，
//   循环往往早已跑完，大家读到的都是 i 的终值（本例是 11）=> 打印十个 11（经典面试题）。
// Go 1.22 起（要求 go.mod >= 1.22，本项目是 1.25）：每轮迭代都会创建新的 i，
//   直接闭包引用 i 也"碰巧"正确了 —— 但这依赖语言版本，老代码/被复制的代码常翻车。
// 把 i 作为参数传入：go 执行的一瞬间值就被"复制"进参数，每个 goroutine 各持一份，
//   与 i 后来的变化彻底无关 —— 任何 Go 版本都对，意图也更明显。
func loopVarTrapDemo() {
	var wg sync.WaitGroup

	// 人为复刻老语义：i 声明在循环【外】=> 全程只有一个共享的 i（任何版本都如此）
	var i int
	for i = 1; i <= 10; i++ {
		wg.Add(1)
		go func() { // 闭包引用共享的 i
			defer wg.Done()
			fmt.Print(i, " ")
		}()
	}
	wg.Wait()
	fmt.Println("<- 闭包直接引用共享变量：全是 11（循环结束时的终值）！")

	// 对照组：同样在循环外声明，但 go 时把值复制进参数
	k := 0
	for k = 1; k <= 10; k++ {
		wg.Add(1)
		go func(n int) { // go 执行瞬间，k 的当前值被复制进 n
			defer wg.Done()
			fmt.Print(n, " ")
		}(k)
	}
	wg.Wait()
	fmt.Println("<- 传参复制值：1~10（顺序随机）")
}

// ------------------------------------------------------------
// 3. 用 sync.Once 实现一个懒加载的单例配置对象。
//
// 【懒加载】第一次用到时才初始化，而不是程序一启动就构造（省内存/启动快）。
// 【单例】无论多少 goroutine 同时调用 GetConfig，拿到的都是同一个 *Config。
// sync.Once 内部用原子操作 + 互斥锁实现，Do 是并发安全的：
//   - 第一个调用 Do 的 goroutine 执行 f，其余调用者阻塞等待它完成
//   - 之后再调用 Do 直接返回，f 不会执行第二次
// ------------------------------------------------------------
type Config struct {
	Addr    string
	Timeout time.Duration
}

var (
	configOnce sync.Once // 零值即可用，不需要初始化（和 Mutex 一样）
	instance   *Config   // 在 once.Do 的函数体里赋值，写之前无人能读到它
)

// GetConfig 所有地方都通过它获取配置（不要直接碰 instance 变量）
func GetConfig() *Config {
	configOnce.Do(func() {
		fmt.Println("  [加载配置] 这行只应出现一次...")
		time.Sleep(50 * time.Millisecond) // 模拟读文件/连数据库等昂贵初始化
		instance = &Config{Addr: "127.0.0.1:8080", Timeout: 3 * time.Second}
	})
	return instance
}

// singletonDemo 并发调用 GetConfig，验证"只初始化一次 + 全员同一个实例"
func singletonDemo() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := GetConfig()
			// %p 打印指针地址：5 行的地址应完全相同 => 同一个对象
			fmt.Printf("  goroutine %d 拿到配置 %p = %+v\n", id, cfg, *cfg)
		}(i)
	}
	wg.Wait()

	// 主 goroutine 再取一次：这次 Do 里的函数不会再执行（上面没打印第二次）
	fmt.Println("  再次获取:", *GetConfig())
}

