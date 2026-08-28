// ============================================================
// 02 - 变量与常量
//
// 运行：go run ./02_variables
// ============================================================

package main

import "fmt"

// 【包级变量】在函数外声明，整个包内可见（只能用 var，不能用 :=）
var globalCount = 100 // 声明时可以不给类型，由初始值自动推断为 int

// var 也可以分组声明（Go 惯用写法）
var (
	serverName = "godemo"
	port       = 8080 // 同类型可省略前面变量的类型
	debug      bool   // 没写初始值就必须显式给类型，默认值为 false
)

func main() {
	// ---------- 1. 声明变量的三种方式 ----------

	// 方式一：var 关键字，类型在后（注意：Go 的类型都写在变量名后面！）
	var a int = 10 // 完整形式：var 名字 类型 = 值
	var b = 20     // 省略类型：由初始值推断为 int
	var c string   // 省略初始值：自动为"零值"（string 的零值是 ""）
	fmt.Println(a, b, c == "")

	// 方式二：短变量声明 :=（只能在函数内使用，最常用！）
	x := 42        // 自动推断类型
	y := "hello"   // string
	z := 3.14      // float64
	fmt.Println(x, y, z)

	// 方式三：多重赋值
	i, j, k := 1, 2, "三" // 多个变量一次声明，类型可以各不相同
	i, j = j, i           // 【经典】一行交换两个变量的值（无需临时变量）
	fmt.Println(i, j, k)

	// 【重要】声明了却未使用的变量会导致编译错误！
	// 取消下面这行的注释试试：unused := 1

	// ---------- 2. 零值（zero value）----------
	// Go 中变量声明后一定有值，没赋值就是类型的"零值"：
	//   数值类型 -> 0     布尔 -> false
	//   字符串  -> ""     指针/切片/map/接口/函数/channel -> nil
	var (
		zi    int
		zf    float64
		zb    bool
		zs    string
		zp    *int // 指针，第 07 课详讲
		zslic []int
	)
	fmt.Println("零值：", zi, zf, zb, zs, zp == nil, zslic == nil)

	// ---------- 3. 空白标识符 _ ----------
	// _ 用来丢弃不需要的值，"我声明了但不使用它"就不会报错。
	_, keep := divide(10, 3) // 只要第二个返回值（函数在下方）
	fmt.Println("10/3 的余数 =", keep)

	// ---------- 4. 常量 const ----------
	// 常量在编译期确定，不可修改。类型可以是任意基本类型。
	const Pi = 3.14159
	const (
		StatusOK  = 200 // 无类型常量，使用时根据上下文自动适配类型
		Language  = "Go"
		MaxSize   = 1 << 20 // 1MB，常量可以做编译期运算
	)
	const Big float64 = 1e20 // 也可以显式指定类型
	fmt.Println(Pi, StatusOK, Language, MaxSize, Big)

	// ---------- 5. iota：枚举生成器 ----------
	// iota 在 const 块中从 0 开始，每新增一行自动 +1
	const (
		Sunday    = iota // 0
		Monday           // 1（自动重复上一行的表达式 iota）
		Tuesday          // 2
		Wednesday        // 3
	)
	fmt.Println("Sunday =", Sunday, "Wednesday =", Wednesday)

	// iota 的常见技巧：位运算枚举（权限位）
	const (
		Read    = 1 << iota // 1 (001)
		Write               // 2 (010)
		Execute             // 4 (100)
	)
	perm := Read | Write // 3：同时拥有读和写权限
	fmt.Printf("权限=%d 可读=%v 可执行=%v\n", perm, perm&Read != 0, perm&Execute != 0)

	// iota 技巧二：跳值（_ 占位跳过一行）
	const (
		_  = iota // 跳过 0
		KB = 1 << (10 * iota) // 1 << 10
		MB                    // 1 << 20
		GB                    // 1 << 30
	)
	fmt.Println("KB =", KB, "MB =", MB, "GB =", GB)

	// ---------- 6. 变量作用域 ----------
	if n := 10; n > 5 { // if 语句中声明的变量只在 if/else 块内有效
		fmt.Println("块内变量 n =", n)
	}
	// fmt.Println(n) // 取消注释会报错：n 在此处未定义
	fmt.Println("包级变量：", globalCount, serverName, port, debug)
}

// divide 返回商和余数（多返回值，第 05 课详讲）
func divide(a, b int) (int, int) {
	return a / b, a % b
}

// ============================================================
// 课后练习：
// 1. 声明三个变量并交换它们的值。
// 2. 用 iota 生成一个季节枚举（Spring=1..Winter=4，提示：_ 跳过 0）。
// 3. 打印 1<<62 的值，体会 Go 整型不溢出的"常量魔法"（变量会溢出）。
// ============================================================
