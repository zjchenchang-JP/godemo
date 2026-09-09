// ============================================================
// 08 - 结构体与方法（Go 的"面向对象"）
//
// 运行：go run ./08_struct_method
// ============================================================

package main

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// ---------- 1. 结构体定义 ----------
// struct 是一组字段的集合，Go 用它代替"类"
type Person struct {
	Name  string // 字段：首字母大写 = 导出(公开)；小写 = 未导出(私有)
	Age   int
	email string // 私有字段：只在包内可见
}

// ---------- 2. 构造函数（Go 没有类和构造函数，用普通函数约定）----------
// 惯例：命名为 NewXxx，返回指针
func NewPerson(name string, age int, email string) *Person {
	return &Person{Name: name, Age: age, email: email}
}

// ---------- 3. 方法：带"接收者"的函数 ----------
// 值接收者 (p Person)：拿到副本，修改不影响原对象
func (p Person) Greet() string {
	return "你好，我是 " + p.Name
}

// 指针接收者 (p *Person)：拿到地址，可以修改原对象
func (p *Person) Birthday() {
	p.Age++ // 语法糖：自动解引用，等价于 (*p).Age++
}

// 【规则】方法也要遵循导出规则：Greet/Birthday 大写公开，sleep 小写私有
func (p Person) sleep() string {
	return p.Name + " 在睡觉 zzz"
}

// ---------- 4. 值接收者 vs 指针接收者 ----------
// 调用时的地址自动转换：Go 会自动帮 p.Greet() / p.Birthday() 取地址/解引用
func (p Person) RenameValue(name string) { p.Name = name } // 改的是副本，无效
func (p *Person) RenamePtr(name string)  { p.Name = name } // 真正修改

// ---------- 5. 任意自定义类型都能挂方法（不只是 struct）----------
// 相当于给 int 起了个新名字 MyInt 并挂方法
type MyInt int

func (m MyInt) Double() MyInt { return m * 2 }

// ---------- 6. 嵌入(Embedding)：Go 的"继承"其实是组合 ----------
type Animal struct {
	Name string
}

func (a Animal) Describe() string { return "我是动物: " + a.Name }

type Dog struct {
	Animal // 【嵌入字段】：没有字段名，只有类型
	Breed  string
}

// Dog 可以"继承"Animal 的方法（实际是编译器代理：dog.Describe() -> dog.Animal.Describe()）
// 也可以重写（方法遮蔽）
func (d Dog) Describe() string { return "我是狗狗: " + d.Name + "，品种: " + d.Breed }

// ---------- 7. 结构体标签 tag（反射时使用，见第 15 课）----------
type User struct {
	Name string `json:"name" validate:"required"` // 反引号内的元信息
	Age  int    `json:"age"`
}

func main() {
	// ---------- 结构体的创建方式 ----------
	p1 := Person{Name: "Alice", Age: 30, email: "a@x.com"} // 【推荐】指定字段名，顺序无关
	p2 := Person{"Bob", 25, "b@x.com"}                     // 按声明顺序（不推荐）
	p3 := Person{Name: "Charlie"}                          // 部分初始化，其余为零值
	var p4 Person                                          // 零值结构体：{"" 0 ""}
	p5 := new(Person)                                      // 返回指针 *Person

	fmt.Println(p1, p2, p3, p4, p5)
	fmt.Printf("%+v\n", p1) // %+v 打印带字段名: {Name:Alice Age:30 email:a@x.com}
	fmt.Printf("%T", p5) // *main.Person 这是什么类型？？
	fmt.Println("")

	// 结构体是值类型：赋值即拷贝（内含切片/map 字段时是"浅拷贝"，注意！）
	pc := p1
	pc.Name = "AliceCopy"
	fmt.Println(p1.Name, pc.Name) // Alice AliceCopy，互不影响

	// 结构体比较：所有字段都可比较时，可以用 == 直接比较
	fmt.Println(p3 == Person{Name: "Charlie"}) // true
	fmt.Println(p3 == Person{Age: 100})        // false 如果字段不对齐时，或者想自定义比较规则时，如何比较？
	p3.Age = 1
	p3.Name = ""
	fmt.Println("修改后: ", p3 == Person{Age: 1}) // true

	// ---------- 方法调用 ----------
	fmt.Println(p1.Greet())
	p1.Birthday() // 即使 p1 是值变量，也能调用指针方法（Go 自动取地址）
	fmt.Println(p1.Name, p1.Age)

	// 值/指针接收者的区别
	p1.RenameValue("X")
	fmt.Println("RenameValue 后:", p1.Name) // Alice，没变！
	p1.RenamePtr("Alice2")
	fmt.Println("RenamePtr 后:", p1.Name) // Alice2，变了！

	// ---------- 自定义类型的方法 ----------
	m := MyInt(21)
	fmt.Println(m.Double()) // 42

	// ---------- 嵌入 ----------
	d := Dog{Animal: Animal{Name: "旺财"}, Breed: "柴犬"}
	fmt.Println(d.Describe())        // 调用 Dog 自己的（重写）
	fmt.Println(d.Animal.Describe()) // 显式调用被"遮蔽"的父方法

	// 嵌入字段的提升（promotion）：d.Name 其实是 d.Animal.Name
	fmt.Println(d.Name) // 旺财

	// ---------- 方法值与方法表达式（了解）----------
	greet := p1.Greet // 方法值：绑定了 p1 的函数
	fmt.Println(greet())
	fmt.Printf("greet的类型是=%T", greet)
	fmt.Println()

	// ---------- 结构体 tag ----------
	u := User{Name: "Dave", Age: 40}
	fmt.Printf("tag 演示（默认格式）: %v\n", u) // {Dave 40}
	fmt.Println(u)// {Dave 40}
	fmt.Printf("%+v\n",u)

	// 用反射读取 tag（详细原理见第 15 课反射）
	t := reflect.TypeOf(u)
	if f, ok := t.FieldByName("Name"); ok {
		fmt.Printf("Name 字段的 tag: %q\n", f.Tag) // "json:\"name\" validate:\"required\""
		fmt.Println("json 名:", f.Tag.Get("json")) // name
	}

	// ---------- 课后练习演示 ----------
	// 1) Stringer 生效：现在打印 Person 自动走 String()
	//    （注意：文件开头第 87 行 fmt.Println(p1, p2, ...) 的输出从此也全变了）
	fmt.Printf("%v\n", p1) // Person{Name=Alice2, Age=31}
	fmt.Println(&p1)       // 打印指针同样触发

	// 2) Rectangle / Circle 各自的 Area()
	r := Rectangle{W: 3, H: 4}
	c := Circle{R: 5}
	fmt.Printf("矩形面积=%.2f 圆面积=%.2f\n", r.Area(), c.Area())

	// 3) Employee：Person 的字段和方法全部提升
	e := NewEmployee("小明", 30, 15000)
	fmt.Println(e.Greet()) // 复用 Person 的方法
	fmt.Println(e.Name)    // 等价于 e.Person.Name
	e.Birthday()           // 指针方法照样提升调用
	fmt.Println(e)                // Employee 自己的 String()（外层遮蔽）→ 连工资一起打印
	fmt.Println(e.Person.String()) // 显式调用被遮蔽的 Person 版，同 d.Animal.Describe() 的写法
}

// ============================================================
// 课后练习：
// 1. 给 Person 加方法 String() string，让 Printf 用 %v 打印时更友好。
//    （提示：这叫实现 fmt.Stringer 接口，第 09 课会讲）
// 2. 设计 Rectangle 和 Circle，都挂 Area() 方法，体会为什么
//    大结构体适合用指针接收者。
// 3. 用嵌入实现 Employee 包含 Person，再加 Salary 字段。
// ============================================================
// ---------- 练习 1：Person 实现 fmt.Stringer ----------
// fmt 在 %v / %s / Println 打印时，会优先调用 String()
// 【坑1】String() 只负责"返回描述"，千万别改字段——它会被任何打印语句
//        随时调用，改字段 = 打印一次污染一次数据
// 【坑2】接收者要用值 (p Person)：用指针的话 fmt.Println(p1) 打印【值】不触发，
//        只有 fmt.Println(&p1) 打印指针才触发
// 【坑3】方法内部不能再用 %v 打印自己（fmt.Sprintf("%v", p)）→ 无限递归栈溢出
func (p Person) String() string {
	return "Person{Name=" + p.Name + ", Age=" + strconv.Itoa(p.Age) + "}"
}

// ---------- 练习 2：Rectangle / Circle 都挂 Area() ----------
type Rectangle struct{ W, H float64 }
type Circle struct{ R float64 }

func (r Rectangle) Area() float64 { return r.W * r.H }
func (c Circle) Area() float64    { return math.Pi * c.R * c.R }

// 为什么"大结构体适合用指针接收者"？
// 值接收者 = 每次调用方法都完整【拷贝一份】结构体。
// Rectangle 只有 16 字节，拷贝无所谓；下面这位每次调用要拷 1MB：
type BigImage struct {
	Pixels [1024 * 1024]byte // 1MB 的字段
}

// 指针接收者：只传 8 字节地址，不拷贝 1MB，还能就地修改
func (b *BigImage) SetAll(v byte) {
	for i := range b.Pixels {
		b.Pixels[i] = v
	}
}

// ---------- 练习 3：嵌入实现 Employee ----------
// Person 的字段和方法全部"提升"：e.Name / e.Greet() / e.Birthday() 都能直接用
type Employee struct {
	Person
	Salary int32
}

func NewEmployee(name string, age int, salary int32) *Employee {
	return &Employee{Person: Person{Name: name, Age: age}, Salary: salary}
}

// 【细节】Person.String() 也会被提升 → 直接打印 Employee 只显示 Person 部分（Salary 消失）
// 原因：提升来的 String() 其实是 e.Person.String()，它的接收者是 Person 那半个，
//       代码里只看得见 Name/Age，根本不知道有 Salary 这个字段
// 解法：Employee 自己定义 String() 遮蔽它（原理同上面 Dog.Describe 遮蔽 Animal.Describe）
func (e Employee) String() string {
	return fmt.Sprintf("Employee{名字:%s, 年龄:%d, 工资:%d}", e.Name, e.Age, e.Salary)
}