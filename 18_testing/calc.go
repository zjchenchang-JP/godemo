// Package calc 是被测试的目标代码。
//
// 注意：这个包不是 main，所以不能 go run，只能被导入或被测试。
// Go 的规定：测试文件必须与被测代码在【同一个包】，
// 因此 calc_test.go 也声明为 package calc。

package calc

import "fmt"

// Add 返回两数之和。
func Add(a, b int) int {
	return a + b
}

// Divide 安全的除法：除数为零时返回错误而不是 panic。
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零: %g / %g", a, b)
	}
	return a / b, nil
}

// IsPrime 判断 n 是否为质数。
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
