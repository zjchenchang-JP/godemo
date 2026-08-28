// ============================================================
// 09 - 接口：Go 面向对象的核心（隐式实现 / 多态 / 类型断言）
//
// 运行：go run ./09_interface
// ============================================================

package main

import "fmt"

// ---------- 1. 接口定义 ----------
// 接口 = 方法签名的集合。"谁实现了这些方法，谁就是这个接口"。
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 【核心】Go 的接口是【隐式实现】的：
// 不需要像 Java/C# 那样写 "class Circle implements Shape"，
// 只要类型实现了接口要求的全部方法，就自动满足该接口 —— 这叫"鸭子类型"。

// Circle 用结构体实现 Shape
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64      { return 3.14159 * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * 3.14159 * c.Radius }

// Rect 也实现了 Shape
type Rect struct {
	W, H float64
}

func (r Rect) Area() float64      { return r.W * r.H }
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }

// ---------- 2. 接口组合（嵌入其他接口）----------
type Describable interface {
	Describe() string
}
type BetterShape interface {
	Shape          // 嵌入 Shape 接口：拥有 Area/Perimeter
	Describe() string // 再加上自己的方法
}

// ---------- 3. 常见的标准库接口 ----------
// fmt.Stringer：实现了 String() string 的类型，用 Printf 打印时会自动调用
type Temperature float64

func (t Temperature) String() string {
	return fmt.Sprintf("%.1f°C", float64(t))
}

func main() {
	// ---------- 多态：一个接口变量装不同类型 ----------
	var s Shape // 接口变量（此时 s == nil）

	s = Circle{Radius: 5}
	fmt.Printf("圆  面积=%.2f 周长=%.2f\n", s.Area(), s.Perimeter())

	s = Rect{W: 3, H: 4}
	fmt.Printf("矩形 面积=%.2f 周长=%.2f\n", s.Area(), s.Perimeter())

	// 接口可以装"任何实现了它的类型"，用切片收集
	shapes := []Shape{Circle{1}, Rect{2, 3}, Circle{10}}
	total := 0.0
	for _, sh := range shapes { // 同一段代码处理不同类型 = 多态
		total += sh.Area()
	}
	fmt.Printf("总面积 = %.2f\n", total)

	// 接口作为函数参数（最常见用法）
	printShapeInfo(Circle{2})
	printShapeInfo(Rect{5, 6})

	// ---------- 4. 类型断言：从接口中取回具体类型 ----------
	var v any = "hello" // any 是 interface{} 的别名（Go 1.18+）

	str := v.(string)         // 断言 v 里装的是 string
	fmt.Println("断言结果:", str, "长度:", len(str))

	// 【安全写法】comma-ok：断言失败不会 panic，ok 为 false
	n, ok := v.(int)
	fmt.Println("v 是 int 吗?", ok, "值:", n) // false 0

	// i.(T) 单值形式失败会直接 panic
	// bad := v.(int) // 取消注释会 panic: interface conversion

	// ---------- 5. 类型开关（type switch）：对类型分支处理 ----------
	for _, item := range []any{42, "Go", 3.14, true, []int{1, 2}} {
		describeAny(item)
	}

	// ---------- 6. 空接口 any/interface{} ----------
	// 不含任何方法的接口，任何类型都实现了它 => 能装任何值
	var data any = 100
	data = "字符串也行"
	data = []float64{1.1, 2.2}
	fmt.Println("any 可以装任何值:", data)

	// ---------- 7. fmt.Stringer 接口 ----------
	t := Temperature(36.6)
	fmt.Println("体温:", t) // 打印时自动调用 t.String()

	// ---------- 8. 接口的底层（理解 nil 陷阱）----------
	// 接口值 = (动态类型, 动态值) 二元组。
	// 【经典坑】类型不为 nil 但值为 nil 的接口变量 != nil！
	var p *Circle // p 是 nil 指针
	var si Shape = p
	fmt.Println("si == nil ?", si == nil) // false！（类型信息 Circle 存在）
	// 所以：函数返回 error 接口时，不要返回具体类型的 nil 指针

	// ---------- 9. 面向接口编程的原则 ----------
	// "接口由使用方定义，而不是实现方"（Go 惯例：在消费接口的地方声明小接口）
	// 标准库典范：io.Reader / io.Writer / fmt.Stringer / sort.Interface ...
}

// printShapeInfo 接收接口参数：能传入任何 Shape
func printShapeInfo(s Shape) {
	fmt.Printf("  类型=%T 面积=%.2f\n", s, s.Area())
}

// describeAny 用 type switch 分别处理不同类型
func describeAny(v any) {
	switch x := v.(type) { // 注意：x 在每个 case 中是对应的具体类型
	case int:
		fmt.Println("整数:", x)
	case string:
		fmt.Println("字符串:", x, "长度:", len(x))
	case float64:
		fmt.Println("浮点数:", x)
	case bool:
		fmt.Println("布尔值:", x)
	case []int:
		fmt.Println("int切片:", x, "长度:", len(x))
	default:
		fmt.Printf("未知类型: %T\n", x)
	}
}

// ============================================================
// 课后练习：
// 1. 定义 Animal 接口(Sound()/Name())，让 Cat、Dog 实现，
//    用 []Animal 遍历打印（多态练习）。
// 2. 给 08 课的 Person 实现 String() 方法，体验 fmt.Stringer。
// 3. 思考：为什么 Go 推荐"接受接口，返回具体类型"？
// ============================================================
