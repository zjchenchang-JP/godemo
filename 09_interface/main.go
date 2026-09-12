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
	// 只有 Area、没有 Perimeter 的类型现在也能传（改造前会编译报错：没实现 Shape）
	// Triangle does not implement Shape (missing method Perimeter)
	printShapeInfo(Triangle{Base: 3, H: 4})

	// ---------- 4. 类型断言：从接口中取回具体类型 ----------
	var v any = "hello" // any 是 interface{} 的别名（Go 1.18+）

	str := v.(string)         // 断言 v 里装的是 string
	fmt.Println("断言结果:", str, "长度:", len(str))

	// 【安全写法】comma-ok：断言失败不会 panic，ok 为 false
	n, ok := v.(int)
	fmt.Println("v 是 int 吗?", ok, "值:", n) // false 0
	var v2 any = 3
	n1, ok := v2.(string)
	fmt.Println("值是：",n1) // ""

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
	var p *Circle // p 是 nil 指针  静态类型是*Circle
	var si Shape = p
	fmt.Println("si == nil ?", si == nil) // false！（类型信息 Circle 存在）
	// 所以：函数返回 error 接口时，不要返回具体类型的 nil 指针
	fmt.Printf("si类型是=%T\n", si)
	fmt.Println("si值是=",si)
	fmt.Printf("p类型是=%T\n", p)
	fmt.Println("p值是=",p)

	// ---------- 练习1演示：Animal 多态 ----------
	animals := []Animal{
		Cat{name: "咪咪", Action: "抓老鼠"},
		Dog{Cat{name: "旺财", Action: "看家"}},
	}
	for _, a := range animals {
		fmt.Printf("我是 %s: ", a.Name()) // Dog 的 Name() 是从 Cat 提升来的
		a.Sound()                         // 同一句代码，不同叫声 = 多态
	}
	
}
// ---------- 9. 面向接口编程的原则 ----------
// "接口由使用方定义，而不是实现方"（Go 惯例：在消费接口的地方声明小接口）
// 标准库典范：io.Reader / io.Writer / fmt.Stringer / sort.Interface ...

// ---------- "接口由使用方定义"的活例子 ----------
// 【原版】要求整个 Shape（Area+Perimeter 两个方法），但函数体只用 Area()
// 大接口会把"只有 Area 的类型"挡在门外 —— 接口越大，抽象越弱
// func printShapeInfo(s Shape) {
// 	fmt.Printf("  类型=%T 面积=%.2f\n", s, s.Area())
// }

// “接口由使用方定义”的完整心智模型是：接口不是类型的标签，而是函数开出的需求清单
// 在消费现场开、按最小需求开、开完自动生效
// 【方式一】使用方现场声明小接口：只列自己真正用到的方法（惯例 -er 结尾）
type Arear interface {
	Area() float64
}
// “接受接口，返回具体类型”是一体两面：
// 参数收接口：只索取我需要的最小能力（Arear 而不是 Shape）
// 返回具体类型：给调用方完整的字段和方法（比如返回 *Employee 而不是某个接口，调用方想用啥都能拿到）
func printShapeInfo(s Arear) {
	fmt.Printf("  类型=%T 面积=%.2f\n", s, s.Area())
}

// 【方式二】极致版：匿名接口直接写在参数里，连名字都不取
// 【证据】Triangle 只有 Area、没有 Perimeter：
// 改造前传给 printShapeInfo 会编译报错（没实现 Shape），现在畅通无阻
type Triangle struct{ Base, H float64 }

func (t Triangle) Area() float64 { return t.Base * t.H / 2 }

// describeAny 用 type switch 分别处理不同类型
// 接口定义“必须会什么”，断言探测“额外会什么”——一个管契约，一个管彩蛋
// v.(string) 和 r.(io.Closer) 是同一个语法——T 写具体类型就是“身份断言”，T 写接口类型就是“能力断言”
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

// 练习1：Animal 接口 + Cat/Dog 多态
type Animal interface {
	Sound()
	Name() string // 返回名字（打印的事交给调用方）
}

// 【坑】字段和方法共用同一个命名空间：既然有 Name() 方法，
//       就不能再有 Name 字段（c.Name 分不清是字段还是方法）→ 编译报错
// 【惯例】字段小写 name（私有），方法大写 Name() 当 getter —— Go 不写 GetName()
type Cat struct {
	name   string
	Action string
}

func (c Cat) Sound() { 
	fmt.Println("喵~ (" + c.Action + ")") 
}
func (c Cat) Name() string { 
	return c.name 
}

// Dog 嵌入 Cat：Name() 直接提升复用，不用重写（复习 08 课的方法提升/遮蔽）
// 只重写 Sound() 遮蔽 Cat 的版本
type Dog struct {
	Cat
}

func (d Dog) Sound() { fmt.Println("汪汪! (" + d.Action + ")") }