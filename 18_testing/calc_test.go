// ============================================================
// 18 - 测试：单元测试 / 表驱动测试 / 基准测试 / Example
//
// 运行本文件（在项目根目录执行）：
//   go test ./18_testing -v              # 运行全部测试，显示细节
//   go test ./18_testing -run TestAdd    # 只跑指定测试
//   go test ./18_testing -bench=. -v     # 跑基准测试
//   go test ./18_testing -cover          # 查看覆盖率
//
// 测试文件的规则：
//   文件名必须以 _test.go 结尾（不会被编译进正式程序）
//   测试函数必须以 Test 开头，且参数是 *testing.T
//   基准函数以 Benchmark 开头，参数是 *testing.B
//   示例函数以 Example 开头，没有参数
// ============================================================

package calc // 【重点】与被测代码同包，可以直接测试私有函数

import (
	"fmt"
	"testing"
)

// ---------- 1. 最基本的单元测试 ----------
func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		// t.Errorf：报告错误但继续执行后续断言
		// t.Fatalf：报告错误并立即终止本测试函数
		t.Errorf("Add(2, 3) = %d; 期望 %d", got, want)
	}
}

// ---------- 2. 表驱动测试（Go 社区最推崇的风格！）----------
// 把用例组织成"输入+期望值"的切片，新增用例只需加一行。
func TestIsPrime(t *testing.T) {
	tests := []struct {
		name string // 子测试名（失败时能精确定位）
		in   int
		want bool
	}{
		{"负数", -7, false},
		{"零", 0, false},
		{"一", 1, false},
		{"二_最小质数", 2, true},
		{"九_合数", 9, false},
		{"十三", 13, true},
		{"九十七", 97, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // t.Run 创建子测试
			if got := IsPrime(tt.in); got != tt.want {
				t.Errorf("IsPrime(%d) = %v; 期望 %v", tt.in, got, tt.want)
			}
		})
	}
}

// ---------- 3. 测试错误分支 ----------
func TestDivide(t *testing.T) {
	// 正常路径
	got, err := Divide(10, 4)
	if err != nil || got != 2.5 {
		t.Errorf("Divide(10, 4) = (%v, %v); 期望 (2.5, nil)", got, err)
	}

	// 错误路径：除零必须返回 error 而不是 panic
	if _, err := Divide(1, 0); err == nil {
		t.Error("Divide(1, 0) 应当返回 error，却返回了 nil")
	}
}

// ---------- 4. Example：可执行文档 ----------
// 注释里的 "// Output:" 是关键：go test 会真实运行函数，
// 拿实际输出和 Output 后的内容比对。它同时会出现在 godoc 文档里。
func ExampleAdd() {
	fmt.Println(Add(20, 22))
	// Output: 42
}

// ---------- 5. 基准测试：测量性能 ----------
// 运行：go test -bench=. ./18_testing
// b.N 由框架动态调整，直到获得稳定的测量结果
func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(7919)
	}
}

// ---------- 6. 其他常用测试工具（本文件未展开，了解即可）----------
// go test -race                检测数据竞争（强烈建议在 CI 里开启）
// go test -cover -coverprofile=c.out && go tool cover -html=c.out  覆盖率报告
// testing.T.Helper()           标记辅助函数，报错时定位到调用行
// t.Skip("跳过原因")           跳过某个测试（如环境不满足）
// testify 库                   社区最流行的断言库（assert/require）
// go test -fuzz=. ./xxx        模糊测试（Go 1.18+）
