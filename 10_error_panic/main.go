// ============================================================
// 10 - 错误处理与 panic/recover
//
// 运行：go run ./10_error_panic
// ============================================================

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// ---------- 1. error 只是一个接口 ----------
// 标准库定义：type error interface { Error() string }
// 所以任何实现了 Error() string 的类型都是 error

// ---------- 2. 哨兵错误（sentinel error）：预定义的错误值 ----------
var ErrNotFound = errors.New("记录不存在")

// ---------- 3. 自定义错误类型（结构体实现 error 接口）----------
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("字段校验失败 [%s]: %s", e.Field, e.Msg)
}

// ---------- 4. 演示用的业务函数 ----------
func findUser(id int) (string, error) {
	users := map[int]string{1: "Alice", 2: "Bob"}
	name, ok := users[id]
	if !ok {
		return "", ErrNotFound // 返回哨兵错误
	}
	return name, nil
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Msg: "不能为负数"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Msg: "太大了"}
	}
	return nil
}

// ---------- 5. 错误包装（wrap）：给底层错误加上下文 ----------
func loadConfig(path string) error {
	// 模拟底层读取失败
	_, err := strconv.Atoi("abc")
	if err != nil {
		// %w 动词：包装错误，保留原始错误链（区别于 %v 只拼接文字）
		return fmt.Errorf("加载配置文件 %s 失败: %w", path, err)
	}
	return nil
}

func main() {
	// ---------- 标准错误处理范式 ----------
	name, err := findUser(1)
	if err != nil {
		fmt.Println("出错了:", err)
	} else {
		fmt.Println("找到用户:", name)
	}

	// 忽略错误时用 _ 显式声明（表示"我知道有错误但我故意忽略"）
	n, _ := strconv.Atoi("42")
	fmt.Println("n =", n)

	// ---------- errors.Is：判断错误链中是否包含指定错误 ----------
	if _, err := findUser(99); err != nil {
		if errors.Is(err, ErrNotFound) { // 类似 ==，但会沿着包装链查找
			fmt.Println("用 errors.Is 匹配到 ErrNotFound")
		}
	}

	// ---------- errors.As：从错误链中提取自定义错误类型 ----------
	if err := validateAge(-5); err != nil {
		var ve *ValidationError
		if errors.As(err, &ve) { // 取出具体类型，能访问字段
			fmt.Println("校验错误字段:", ve.Field, "原因:", ve.Msg)
		}
	}

	// ---------- 错误包装链 ----------
	if err := loadConfig("app.yaml"); err != nil {
		fmt.Println("外层:", err)
		fmt.Println("用 errors.Unwrap 取底层:", errors.Unwrap(err))
	}

	// ---------- defer + recover：捕获 panic ----------
	// panic：不可恢复的严重错误（数组越界、空指针等），会导致程序崩溃。
	// recover：只能在 defer 中调用，能"接住" panic 让程序继续运行。
	result := safeDivide(10, 0)
	fmt.Println("safeDivide 结果:", result)

	// panic 的常见合法用途：
	// 1) 程序启动时配置缺失/环境错误（启动即失败）
	// 2) "不可能到达" 的分支（防御性编程）
	// 3) 标准库/框架中的编程错误提示
	// 【原则】业务错误永远用 error 返回，不要用 panic！

	// ---------- Go 1.13+ 的错误处理要点总结 ----------
	// 1. fmt.Errorf("...: %w", err)   包装错误
	// 2. errors.Is(err, target)       判断是否是某个错误（含包装链）
	// 3. errors.As(err, &target)      提取某类型错误
	// 4. errors.Join(err1, err2)      合并多个错误（Go 1.20+）
	joined := errors.Join(errors.New("错误A"), errors.New("错误B"))
	fmt.Println("合并错误:", joined)

	// 练习1
	card, err := ParseCard("♥13")
	if err != nil {
		var ice *InvalidCardError
		if errors.As(err, &ice) {
			fmt.Println("输入出错：", ice.Input, "原因：", ice.Msg)
		}
		// return
	} else {
		fmt.Println(card)
	}

	// 练习2. 用 errors.Is 处理 os.Open 的 os.ErrNotExist。
	f, err := os.Open("config.xxx")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("文件不存在， 改用默认配置")
			// 常规降级业务分支
			return
		}
		fmt.Println("打开失败：", err)
		return
	}
	defer f.Close()
}

// safeDivide 用 recover 捕获除零 panic，转为 error 返回
func safeDivide(a, b int) (result int) {
	// recover 必须写在 defer 的匿名函数里才能生效
	err := func() (err error) {
		defer func() {
			if r := recover(); r != nil { // r 是 panic 的值
				err = fmt.Errorf("捕获到 panic: %v", r)
			}
		}()
		result = a / b // b=0 时这里会 panic: integer divide by zero
		return nil
	}() // IIFE（Immediately Invoked Function Expression，立即调用的函数表达式）
	if err != nil {
		fmt.Println("处理错误:", err)
		return 0
	}
	return result
}

// // 更简洁
// func safeDivide(a, b int) (result int, err error) {
//     defer func() {
//         if r := recover(); r != nil {
//             err = fmt.Errorf("捕获到 panic: %v", r)
//         }
//     }()
//     return a / b, nil
// }

// ============================================================
// 课后练习：
// 1. 写 ParseCard(cardStr string) (Card, error)，非法输入时返回
//    自定义的 InvalidCardError（包含原始字符串字段）。
// 2. 用 errors.Is 处理 os.Open 的 os.ErrNotExist。
// 3. 思考：为什么 recover 必须在 defer 函数中直接调用？
// ============================================================

// 1. 写 ParseCard(cardStr string) (Card, error)
type Card struct {
	Suit string // 花色：♠ ♥ ♦ ♣
	Rank string // 点数：A 2~10 J Q K
}

// 自定义错误：关键是保留"原始输入"这个上下文
type InvalidCardError struct {
	Input string // 哪个字符串出的错
	Msg   string // 为什么错
}

func (e *InvalidCardError) Error() string {
	return fmt.Sprintf("非法卡牌 %q: %s", e.Input, e.Msg)
}

var validRanks = map[string]bool{
	"A": true, "2": true, "3": true, "4": true, "5": true, "6": true,
	"7": true, "8": true, "9": true, "10": true, "J": true, "Q": true, "K": true,
}

var validSuits = map[string]bool{"♠": true, "♥": true, "♦": true, "♣": true}

func ParseCard(cardStr string) (Card, error) {
	runes := []rune(cardStr) // 转成 rune 切片，花色符号是多字节字符
	if len(runes) < 2 {
		return Card{}, &InvalidCardError{Input: cardStr, Msg: "长度不足"}
	}
	rank := string(runes[:len(runes)-1]) // 除最后一个字符外都是点数
	suit := string(runes[len(runes)-1])  // 最后一个字符是花色

	if !validRanks[rank] {
		return Card{}, &InvalidCardError{Input: cardStr, Msg: "点数不合法: " + rank}
	}
	if !validSuits[suit] {
		return Card{}, &InvalidCardError{Input: cardStr, Msg: "花色不合法: " + suit}
	}
	return Card{Suit: suit, Rank: rank}, nil
}

/**
练习 3：为什么 recover 必须在 defer 里直接调用？
分两层理解：

（1）为什么必须在 defer 里 —— 时机问题

panic 的执行模型是「栈展开」：

panic 一发生，当前函数立刻中断，后面的代码一行都不执行
runtime 开始沿调用栈回退，只执行沿途函数注册过的 defer
defer 跑完还没人 recover → 整个程序崩溃
所以 panic 发生后，唯一还在执行的代码就是 defer 函数——recover 写在别处根本没有运行的机会。写在 panic 之前？那时还没 panic，返回 nil；写在函数外？栈已经展开过去了，轮不到它执行。defer 是 panic 世界里唯一的“逃生舱口”。

（2）为什么必须“直接”调用 —— 识别机制问题

recover 生效的条件是：它的调用者必须恰好是那个正在被栈展开执行的 defer 函数。runtime 靠这一点判断“该在哪一层停止展开、让哪个函数正常返回”。隔着函数包一层，recover 就认不出自己身处 defer 上下文，返回 nil。

对照着看，前两个是经典坑：


defer recover()                  // ✗ 无效：被编译器直接当作 no-op
defer fmt.Println(recover())     // ✗ 无效：参数在注册 defer 时就求值了，那时还没 panic

defer func() {                   // ✓ 标准写法：匿名函数本体直接调 recover
	if r := recover(); r != nil {
		fmt.Println("接住了:", r)
	}
}()

func myRecover() {               // ✓ 也有效：myRecover 本身就是被 defer 的函数
	if r := recover(); r != nil {
		fmt.Println("接住了:", r)
	}
}
defer myRecover()

defer func() { helper() }()      // ✗ 无效：调 recover 的是 helper，
                                 //    不是被 defer 的那个函数（隔了一层）
（3）语义串起来

recover() 的含义是“取消当前的栈展开，让我所属的这个函数正常返回”。它天然是绑定到某一层函数的，所以只能写在那层函数的 defer 里——这也解释了 main.go:121 里 safeDivide 的结构：匿名函数“拥有”那个 defer，panic 被 recover 后它带着命名返回值 err 正常返回，外层代码才得以继续跑。
*/
