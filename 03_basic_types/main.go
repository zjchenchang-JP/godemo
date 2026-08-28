// ============================================================
// 03 - 基本数据类型与类型转换
//
// 运行：go run ./03_basic_types
// ============================================================

package main

import (
	"fmt"
	"strconv" // 字符串与其他类型的转换
	"unicode/utf8"
)

func main() {
	// ---------- 1. 整型 ----------
	// 有符号：int8/16/32/64 + int（平台相关，64位系统上是64位）
	// 无符号：uint8/16/32/64 + uint
	// 两个别名：byte = uint8（字节）、rune = int32（Unicode 码点）
	var i8 int8 = 127       // int8 范围：-128 ~ 127
	ui8 := uint8(255)       // uint8 范围：0 ~ 255
	big := 9_223_372_036_854_775_807 // int64 最大值，可用下划线分隔提高可读性
	hex, oct, bin := 0xFF, 0o17, 0b1010 // 十六进制、八进制、二进制字面量
	fmt.Println(i8, ui8, big, hex, oct, bin)

	// 溢出会"环绕"（不会报错，要小心！）
	overflow := int8(127)
	overflow++
	fmt.Println("int8(127)+1 =", overflow) // -128，环绕到了最小值

	// ---------- 2. 浮点型 ----------
	f32 := 3.14            // 不写类型时，浮点字面量默认 float64
	var f64 float64 = 2.718281828
	fmt.Printf("float32=%.4f float64=%.10f\n", f32, f64)

	// 【经典坑】浮点数有精度误差，不能直接用 == 比较
	a := 0.1 + 0.2
	fmt.Println(0.1+0.2 == 0.3)         // false！
	fmt.Printf("%.17f\n", a)            // 0.30000000000000004
	fmt.Println(abs(a-0.3) < 1e-9)      // 正确做法：差值小于某个极小值

	// ---------- 3. 布尔型 ----------
	ok := true
	fmt.Println(!ok, ok && !ok, ok || !ok) // 非、与、或
	// 注意：Go 不允许把整数当布尔用，如 if 1 {} 是编译错误

	// ---------- 4. 字符串 ----------
	s := "Hello, 世界"
	fmt.Println("长度(字节):", len(s))                   // 13：中文每字占 3 字节(UTF-8)
	fmt.Println("字符数:", utf8.RuneCountInString(s))    // 9

	// 字符串是不可变的 byte 切片
	// s[0] = 'h' // 编译错误！不能修改
	fmt.Printf("第1个字节: %c, 类型: %T\n", s[0], s[0]) // 'H', uint8

	// 遍历字符串：range 按"字符(rune)"遍历，而不是字节
	for i, r := range s {
		fmt.Printf("  下标%d: %q\n", i, r)
	}

	// 原始字符串（反引号）：不转义、可换行，常用于正则和 JSON 模板
	raw := `C:\Windows\n换行也保留
第二行`
	fmt.Println(raw)

	// ---------- 5. rune 与 byte ----------
	var ch rune = '中'      // 单引号 = 字符字面量，类型是 rune(int32)
	var bt byte = 'A'       // 也可以是 byte(uint8)
	fmt.Printf("%c=%d, %c=%d\n", ch, ch, bt, bt)

	// ---------- 6. 类型转换 ----------
	// 【重点】Go 没有隐式类型转换！任何类型之间都必须显式转换。
	// var f float64 = 1     // 整数字面量可以（无类型常量）
	// var f2 float64 = intVar // 但 int 变量不行，必须写 float64(intVar)
	n := 42
	f := float64(n)     // int -> float64
	back := int(f)      // float64 -> int（直接截断小数，不四舍五入）
	ui := uint8(n)      // 范围内转换 OK
	fmt.Println(f, back, ui)

	// 注意：常量直接转 int 会在编译期报错（int(3.99) 不允许），要经过变量
	pi := 3.99
	fmt.Println(int(pi)) // 3：截断，不是 4！

	// 字符串 <-> 数字 要用 strconv 包（不是类型转换！）
	numStr := "123"
	num, err := strconv.Atoi(numStr) // 字符串转 int
	fmt.Println(num+1, err)          // 124 nil
	str := strconv.Itoa(456)         // int 转字符串
	fmt.Println(str + "!")

	fl, _ := strconv.ParseFloat("3.14", 64) // 字符串转 float64
	fmt.Println(fl * 2)

	// int -> string 的【经典大坑】：string(整数) 是按 Unicode 码点转成"一个字符"，
	// 不是把数字变成字符串！比如 string(65) 得到 "A" 而不是 "65"。
	// （新版 go vet 会直接警告这种写法，必须先显式转成 rune）
	code := 65
	fmt.Println(string(rune(code)))  // "A"（不是 "65"！）
	fmt.Println(string(rune(20013))) // "中"
	fmt.Println(strconv.Itoa(code))  // "65"（把数字变成字符串的正确做法）

	// 字符串 <-> []byte / []rune
	bs := []byte(s)  // 用于 IO 操作、网络传输
	rs := []rune(s)  // 用于按字符处理（如取第 N 个字符）
	fmt.Println(string(bs) == s, string(rs[7])) // true 世（第8个字符）

	// ---------- 7. 复数与类型别名（了解）----------
	var cplx complex128 = complex(3, 4) // 复数：3+4i
	fmt.Println(cplx, real(cplx), imag(cplx))
	// 类型别名：byte≈uint8、rune≈int32、any=interface{}（Go 1.18 起）
	var anyVal any = "任意类型，第 09 课详讲"
	fmt.Println(anyVal)
}

// abs 求绝对值
func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// ============================================================
// 课后练习：
// 1. 打印 "Hello, 世界" 的前 5 个字符（提示：[]rune 切片）。
// 2. 用 strconv 把 "3.14" 转成 float64 再乘以 2。
// 3. 思考：string(300) 会输出什么？为什么？（提示：UTF-8 编码）
// ============================================================
