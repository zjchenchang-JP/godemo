// ============================================================
// 13 - 泛型（Go 1.18+）：类型参数、约束、泛型容器
//
// 运行：go run ./13_generics
// ============================================================

package main

import (
	"cmp" // Go 1.21+ 的比较工具包（提供 cmp.Ordered 约束）
	"fmt"
)

// ---------- 1. 泛型解决什么问题？----------
// 没有泛型时，想写"通用求和"只能：为 int 写一份、float 写一份，
// 或用 any + 类型断言（丢失编译期类型检查）。
// 泛型 = 写一份代码，编译期保证类型安全。

// ---------- 2. 自定义约束 ----------
// [T Number] 声明类型参数 T，必须是满足 Number 约束的类型。
// ~int 表示"底层类型是 int"的所有类型（包括 type MyInt int）。
type Number interface {
	~int | ~int64 | ~float64 // 联合类型：三者之一即可
}

// Sum 对任意数值切片求和
func Sum[T Number](nums []T) T {
	var total T // 泛型的零值写法：只声明不赋值
	for _, n := range nums {
		total += n
	}
	return total
}

// ---------- 3. any 约束：不做任何限制 ----------
// Map 把切片的每个元素经 f 转换后返回新切片（输入 T 输出 U，两个类型参数）
func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// Filter 过滤切片元素
func Filter[T any](s []T, keep func(T) bool) []T {
	var out []T
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// ---------- 4. comparable：可比较约束 ----------
// map 的键类型要求可比较，用 K comparable 约束
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// ---------- 5. cmp.Ordered：标准库现成的"可排序"约束 ----------
// 支持 ~int|~float|~string 及其变体，比自定义 Number 更全面
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// ---------- 6. 泛型类型（容器）----------
// Stack[T any]：元素类型为 T 的栈
type Stack[T any] struct {
	items []T
}

// 注意：泛型类型的方法也要带上 [T any]
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// Pop 弹出栈顶元素；空栈返回零值和 false
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T // 零值
		return zero, false
	}
	n := len(s.items) - 1
	v := s.items[n]
	s.items = s.items[:n]
	return v, true
}

func (s *Stack[T]) Len() int { return len(s.items) }

// MyInt 验证 ~int 约束对"底层类型"生效
type MyInt int

func main() {
	// ---------- 泛型函数调用 ----------
	ints := []int{1, 2, 3, 4}
	floats := []float64{1.5, 2.5}
	myInts := []MyInt{10, 20} // MyInt 底层是 int，满足 ~int

	fmt.Println(Sum(ints))    // 类型推断：无需写 Sum[int](ints)
	fmt.Println(Sum(floats))
	fmt.Println(Sum(myInts))  // 没有 ~ 就匹配不了它！
	fmt.Println(Sum[int](ints)) // 显式指定类型参数（推断不出来时用）

	// Map：int 切片 -> string 切片
	words := Map(ints, func(n int) string { return fmt.Sprintf("#%d", n) })
	fmt.Println(words)

	// Filter：保留偶数
	evens := Filter(ints, func(n int) bool { return n%2 == 0 })
	fmt.Println("偶数:", evens)

	// Keys
	m := map[string]int{"a": 1, "b": 2}
	fmt.Println("keys:", Keys(m))

	// Min：同一函数支持数字和字符串
	fmt.Println(Min(3, 7), Min(2.5, 1.5), Min("apple", "banana"))

	// ---------- 泛型容器 ----------
	s1 := &Stack[int]{} // 泛型类型实例化必须显式写类型参数
	s1.Push(1)
	s1.Push(2)
	s1.Push(3)
	v, _ := s1.Pop()
	fmt.Println("int 栈弹出:", v, "剩余:", s1.Len())

	s2 := &Stack[string]{} // 同一个类型，装字符串
	s2.Push("hello")
	sv, _ := s2.Pop()
	fmt.Println("string 栈弹出:", sv)

	// ---------- 何时用泛型？----------
	// 适合：容器/工具函数（栈、队列、映射、过滤）、算法库
	// 不适合：业务代码里滥用（接口通常更简单）。
	// 原则：先用普通函数/接口，真正需要"类型安全的复用"时再上泛型。
}

// ============================================================
// 课后练习：
// 1. 实现泛型函数 Contains[T comparable](s []T, v T) bool。
// 2. 实现泛型队列 Queue[T]（Enqueue/Dequeue）。
// 3. 用 cmp.Slice 排序泛型切片（提示：slices.Sort 也是泛型的）。
// ============================================================
