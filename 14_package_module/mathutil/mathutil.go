// Package mathutil 演示"库包"的写法（包注释写在 package 上方）。
//
// 这是第 14 课的子包：注意目录名是 mathutil，包名也是 mathutil。
// Go 的惯例：目录名 = 包名（导入路径按目录，使用时按包名）。
package mathutil

import "fmt"

// Version 导出的包级变量（大写开头 = 公开）。
var Version = "1.0.0"

// 内部计数：小写开头 = 未导出（包外不可见）
var callCount int

// init 是特殊函数：包被导入时自动执行，先于使用者的 main()。
// 用途：初始化数据库连接、注册驱动、校验配置等。
// 执行顺序：导入包的变量初始化 -> 导入包的 init -> 本包变量初始化 -> 本包 init -> main
func init() {
	fmt.Println("  [mathutil] init() 执行：包被加载了")
}

// Add 返回两数之和（文档注释：以 Add 开头，go doc 会展示）。
func Add(a, b int) int {
	callCount++
	return a + b
}

// Square 返回 n 的平方。
func Square(n int) int {
	callCount++
	return n * n
}

// SumAll 对任意个整数求和。
func SumAll(nums ...int) int {
	total := 0
	for _, n := range nums {
		total = addOne(total, n) // 调用未导出的私有函数：包内随便用
	}
	return total
}

// addOne 是私有函数：只能被本包内的代码调用。
// 在 14_package_module/main.go 里调用 mathutil.addOne 会编译报错！
func addOne(a, b int) int {
	return a + b
}

// Calls 返回本包函数被调用的次数（演示封装私有状态）。
func Calls() int {
	return callCount
}
