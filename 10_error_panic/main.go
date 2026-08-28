// ============================================================
// 10 - 错误处理与 panic/recover
//
// 运行：go run ./10_error_panic
// ============================================================

package main

import (
	"errors"
	"fmt"
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
	}()
	if err != nil {
		fmt.Println("处理错误:", err)
		return 0
	}
	return result
}

// ============================================================
// 课后练习：
// 1. 写 ParseCard(cardStr string) (Card, error)，非法输入时返回
//    自定义的 InvalidCardError（包含原始字符串字段）。
// 2. 用 errors.Is 处理 os.Open 的 os.ErrNotExist。
// 3. 思考：为什么 recover 必须在 defer 函数中直接调用？
// ============================================================
