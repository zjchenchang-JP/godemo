// ============================================================
// 05 - 函数：多返回值、可变参数、闭包、defer、init
//
// 运行：go run ./05_functions
// ============================================================

package main

import (
	"errors"
	"fmt"
	"strings"
)

// ---------- 1. 函数定义的基本形式 ----------
// func 函数名(参数名 类型, 参数名 类型) 返回类型 { ... }
// 连续同类型参数可以只写最后一个类型
func add(a, b int) int {
	return a + b
}

// ---------- 2. 多返回值（Go 的招牌特性）----------
// 返回值写在括号里，通常最后一个返回 error（Go 的错误处理约定）
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil // nil 表示"没有错误"
}

// ---------- 3. 命名返回值 ----------
// 返回值预先命名，函数体内可直接使用，return 可以"裸返回"
// （小函数方便，大函数不建议用，会影响可读性）
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return // 裸返回：自动返回 x, y
}

// ---------- 4. 可变参数 ----------
// nums 本质上是一个切片 []int，必须放在参数最后
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// ---------- 5. 函数是一等公民：可作为参数和返回值 ----------
// 参数 op 是一个函数类型：接收两个 int，返回 int
func calc(a, b int, op func(int, int) int) int {
	return op(a, b)
}

// 返回函数的函数（通常用于生成"配置好"的函数）
func makeGreeter(prefix string) func(string) string {
	return func(name string) string {
		return prefix + " " + name
	}
}

// ---------- 6. init 函数：特殊的初始化函数 ----------
// 执行顺序：包级变量初始化 -> init() -> main()
// 每个文件都可以有自己的 init()，甚至一个文件里可以写多个（但不推荐）
func init() {
	fmt.Println("[init] 包级变量初始化完成后、main 之前自动执行")
}

// 包级变量的初始化先于 init()
var pkgValue = computePkgValue()

func computePkgValue() int {
	fmt.Println("[包级变量] 初始化执行")
	return 42
}

func main() {
	// ---------- 普通调用 ----------
	fmt.Println(add(1, 2))

	// ---------- 多返回值的标准用法 ----------
	result, err := divide(10, 3)
	if err != nil { // Go 的错误处理范式：先检查 err
		fmt.Println("出错了:", err)
		return
	}
	fmt.Println("10/3 =", result)

	// 忽略某个返回值用 _
	_, err2 := divide(10, 0)
	if err2 != nil {
		fmt.Println("捕获错误:", err2)
	}

	// ---------- 命名返回值 ----------
	fmt.Println(split(17)) // 7 10

	// ---------- 可变参数 ----------
	fmt.Println(sum(1, 2, 3))        // 传任意个参数
	fmt.Println(sum())               // 0 个也行
	nums := []int{4, 5, 6}
	fmt.Println(sum(nums...))        // 切片展开传入（注意三个点）

	// ---------- 函数作为参数 ----------
	fmt.Println("加法:", calc(3, 4, func(a, b int) int { return a + b })) // 匿名函数直接传
	fmt.Println("乘法:", calc(3, 4, mul))

	// ---------- 闭包 ----------
	// 闭包 = 匿名函数 + 它引用的外部变量。函数"记住"了外部环境
	counter := makeCounter() // counter 捕获了自己的 count 变量
	fmt.Println(counter(), counter(), counter()) // 1 2 3
	counter2 := makeCounter()
	fmt.Println(counter2()) // 1：每个闭包有独立的状态

	greet := makeGreeter("Hello")
	fmt.Println(greet("Go"), makeGreeter("你好")("世界"))

	// ---------- defer 延迟调用 ----------
	// defer 把函数推迟到"当前函数返回前"执行，常用于释放资源（关文件/解锁）
	// 多个 defer 按【后进先出 LIFO】顺序执行（像栈一样）
	defer fmt.Println("defer 1：最先声明，最后执行")
	defer fmt.Println("defer 2：后声明，先执行")

	// 经典用途：成对操作（打开/关闭、加锁/解锁）
	readFileDemo()

	// 【重要】defer 的参数在 defer 语句处【立即求值】，函数体延迟执行
	x := 10
	defer fmt.Println("defer 捕获的 x =", x) // 打印 10
	x = 20
	fmt.Println("当前 x =", x)

	// 【坑】在循环里 defer 不会立即执行，会堆积到函数结束
	for i := 0; i < 3; i++ {
		defer fmt.Print(i, " ") // 全部堆到 main 结束才打印：2 1 0
	}
}

func mul(a, b int) int {
	return a * b
}

// makeCounter 返回一个闭包：每次调用计数 +1
func makeCounter() func() int {
	count := 0 // 这个变量被闭包捕获，函数返回后依然存活
	return func() int {
		count++
		return count
	}
}

// readFileDemo 演示 defer 关闭资源的惯用写法
func readFileDemo() {
	// 用 strings.Reader 模拟一个数据源（真实文件操作见第 16 课）
	r := strings.NewReader("假装这是一个文件的内容")
	defer func() {
		fmt.Println("defer: 释放资源（关闭文件/连接）")
		r.Len() // 这里仅为演示引用 r
	}()
	buf := make([]byte, 8)
	n, _ := r.Read(buf)
	fmt.Printf("读取到 %d 字节: %q\n", n, buf[:n])
}

// ============================================================
// 课后练习：
// 1. 写一个 minMax(nums ...int) (int, int) 函数，返回最小值和最大值。
// 2. 用闭包实现一个斐波那契生成器 fibonacci() func() int。
// 3. 思考：如果 defer 后面跟的是闭包 defer func(){ print(x) }()，
//    x=20 的修改会影响结果吗？（提示：闭包捕获的是变量引用）
// ============================================================
