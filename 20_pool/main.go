package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	f func() error // 任务逻辑
}

func NewTask(funcArg func() error) *Task { // 构造函数
	return &Task{
		f: funcArg,
	}
}

type Pool struct { // 协程池
	RunningWorkers int64
	Capacity       int64
	JobCh          chan *Task
	sync.Mutex
	WorkerSeq int64 // worker编号
}

func NewPool(capacity int64, taskNum int64) *Pool {
	return &Pool{
		Capacity: capacity,
		JobCh:    make(chan *Task, taskNum),
	}
}

func (p *Pool) GetCap() int64 {
	return p.Capacity // p 是 *Pool，Go 对指针访问字段会自动解引用，p.Capacity 等价于 (*p).Capacity。
}

func (p *Pool) incRunning() { // 记录已激活协程数量
	// 等价于 &((*p).RunningWorkers) 原子操作必须直接作用于那块内存做 读 - 写
	atomic.AddInt64(&p.RunningWorkers, 1)
}

func (p *Pool) decRunning() {
	atomic.AddInt64(&p.RunningWorkers, -1)
}

func (p *Pool) GetRunningWorkers() int64 {
	// 并发场景下 读到真实数据 data race
	// 可见性 普通读没有内存屏障，语言层面不保证读到最新值（编译器、CPU 可能重排指令或读缓存里的旧值）。atomic.Load 有同步语义，保证读到的是内存中一致的最新值
	// Java volatile 只保可见性 + 禁止重排序，不保原子性
	return atomic.LoadInt64(&p.RunningWorkers)
}

func (p *Pool) run() {
	p.incRunning()
	workerID := atomic.AddInt64(&p.WorkerSeq, 1)
	go func() {
		// defer 绑定的是 "声明它的那个函数" 的生命周期 ——那个函数返回时才执行;
		// 放协程外面的话 defer 立刻执行 worker 还在干活，计数却显示没人干活
		defer p.decRunning()
		for task := range p.JobCh {
			fmt.Printf("[worker %d] 开始执行\n", workerID) // ← 打印的就是它
			task.f()
		}
	}()
}

// AddTask 往协程池添加任务
func (p *Pool) AddTask(task *Task) {
	p.Lock()
	defer p.Unlock()
	if p.GetRunningWorkers() < p.GetCap() {
		// 创建启动一个任务
		p.run()
	}
	// 任务推入队列 等待消费
	p.JobCh <- task
}

func main() {
	// 创建协程池
	pool := NewPool(3, 10)
	for i := 0; i < 20; i++ {
		// 任务放入池子
		pool.AddTask(NewTask(func() error {
			fmt.Printf("I am Task %d\n", i)
			return nil
		}))
	}
	time.Sleep(time.Second)// 等待执行
}
