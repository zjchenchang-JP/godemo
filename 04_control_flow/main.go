// ============================================================
// 04 - 流程控制：if / for / switch / break / continue / goto
//
// 运行：go run ./04_control_flow
// ============================================================

package main

import "fmt"

func main() {
	// ==================== if 语句 ====================
	// 特点：1.条件不加小括号  2.花括号必须有  3.左花括号必须和 if 同行
	age := 20
	if age >= 18 {
		fmt.Println("成年人")
	} else if age >= 12 {
		fmt.Println("青少年")
	} else {
		fmt.Println("儿童")
	}

	// 【Go 特色】if 可以带一个初始化语句（分号分隔），变量作用域仅限 if 块
	if n := age * 2; n > 30 {
		fmt.Println("n > 30, n =", n)
	} else {
		fmt.Println("n <= 30, n =", n) // 这里也能用 n
	}
	// fmt.Println(n) // 错误：n 已超出作用域

	// ==================== for 循环 ====================
	// Go 【只有 for】一种循环关键字，但有 3 种形态：

	// 形态一：经典 C 风格（init; condition; post）
	for i := 0; i < 3; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// 形态二：只有条件（相当于其他语言的 while）
	count := 0
	for count < 3 {
		fmt.Print(count, " ")
		count++
	}
	fmt.Println()

	// 形态三：无限循环（死循环，配合 break 退出）
	n := 0
	for {
		n++
		if n > 2 {
			break // break 跳出循环
		}
	}
	fmt.Println("无限循环执行了", n, "次")

	// continue：跳过本次循环的剩余部分，进入下一次
	for i := 0; i < 6; i++ {
		if i%2 == 0 {
			continue // 跳过偶数
		}
		fmt.Print(i, " ") // 1 3 5
	}
	fmt.Println()

	// 【最常用】for range：遍历切片、数组、map、字符串、channel
	nums := []int{10, 20, 30}
	for idx, val := range nums { // 第一个值是下标，第二个是元素
		fmt.Printf("nums[%d]=%d ", idx, val)
	}
	fmt.Println()

	for _, val := range nums { // 只要值，用 _ 丢弃下标
		fmt.Print(val, " ")
	}
	fmt.Println()

	for idx := range nums { // 只要下标（range 一个变量时是下标）
		fmt.Print(idx, " ")
	}
	fmt.Println()

	// 遍历字符串：按字符(rune)
	for i, r := range "Go语言" {
		fmt.Printf("[%d:%c]", i, r) // 注意中文的下标会跳 3
	}
	fmt.Println()

	// ==================== switch 语句 ====================
	// Go 的 switch 默认每个 case 自动 break！不会穿透到下一个 case
	lang := "go"
	switch lang {
	case "python":
		fmt.Println("蟒蛇")
	case "go", "golang": // 一个 case 可以匹配多个值
		fmt.Println("Go!")
	case "java":
		fmt.Println("爪哇")
	default:
		fmt.Println("未知语言")
	}

	// 无表达式的 switch：替代 if-else 链，更清晰（推荐！）
	score := 85
	switch {
	case score >= 90:
		fmt.Println("优秀")
	case score >= 80:
		fmt.Println("良好")
	case score >= 60:
		fmt.Println("及格")
	default:
		fmt.Println("不及格")
	}

	// fallthrough：强制执行下一个 case（不判断条件，很少用）
	switch x := 1; x {
	case 1:
		fmt.Println("一")
		fallthrough // 穿透到 case 2
	case 2:
		fmt.Println("二（被 fallthrough 带进来的）")
	}

	// switch 带初始化语句
	switch os := "windows"; os {
	case "linux", "darwin":
		fmt.Println("Unix 系")
	default:
		fmt.Println("其他系统:", os)
	}

	// 类型开关（type switch）：判断接口变量的动态类型，第 09 课再深入
	var val any = 3.14
	switch v := val.(type) {
	case int:
		fmt.Println("整数:", v)
	case float64:
		fmt.Println("浮点数:", v)
	case string:
		fmt.Println("字符串:", v)
	default:
		fmt.Println("其他类型")
	}

	// ==================== 标签 label ====================
	// break/continue 只能作用于最内层循环；
	// 想跳出外层循环时，给外层循环贴个"标签"
outer:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if i*j > 4 {
				fmt.Println("i*j>4, 跳出全部循环")
				break outer // 直接跳出外层循环
			}
			fmt.Printf("i=%d j=%d | ", i, j)
		}
	}
	fmt.Println()

	// continue + 标签：结束外层本次循环，直接进入外层的下一次
outer2:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if j == 2 {
				continue outer2 // 相当于"跳过外层循环体剩余部分"，i++ 继续
			}
			fmt.Printf("%d-%d ", i, j) // 只会打印 1-1 2-1 3-1
		}
	}
	fmt.Println()

	// ==================== goto（了解，不推荐）====================
	// goto 可以无条件跳转到同函数内的标签，通常只在生成代码里出现
	m := 0
loop:
	m++
	if m < 3 {
		goto loop
	}
	fmt.Println("goto 循环结束, m =", m)
}

// ============================================================
// 课后练习：
// 1. 用 for 打印九九乘法表（嵌套循环 + Printf 对齐 %-2d）。
// 2. 用无表达式 switch 判断某年份是否为闰年。
// 3. 用 break+label 查找二维切片中第一个负数的位置。
// ============================================================
