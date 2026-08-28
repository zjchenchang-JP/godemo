// ============================================================
// 08 - 结构体与方法（Go 的"面向对象"）
//
// 运行：go run ./08_struct_method
// ============================================================

package main

import (
	"fmt"
	"reflect"
)

// ---------- 1. 结构体定义 ----------
// struct 是一组字段的集合，Go 用它代替"类"
type Person struct {
	Name string // 字段：首字母大写 = 导出(公开)；小写 = 未导出(私有)
	Age  int
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
	p2 := Person{"Bob", 25, "b@x.com"}                      // 按声明顺序（不推荐）
	p3 := Person{Name: "Charlie"}                          // 部分初始化，其余为零值
	var p4 Person                                          // 零值结构体：{"" 0 ""}
	p5 := new(Person)                                      // 返回指针 *Person

	fmt.Println(p1, p2, p3, p4, p5)
	fmt.Printf("%+v\n", p1) // %+v 打印带字段名: {Name:Alice Age:30 email:a@x.com}

	// 结构体是值类型：赋值即拷贝（内含切片/map 字段时是"浅拷贝"，注意！）
	pc := p1
	pc.Name = "AliceCopy"
	fmt.Println(p1.Name, pc.Name) // Alice AliceCopy，互不影响

	// 结构体比较：所有字段都可比较时，可以用 == 直接比较
	fmt.Println(p3 == Person{Name: "Charlie"}) // true

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

	// ---------- 结构体 tag ----------
	u := User{Name: "Dave", Age: 40}
	fmt.Printf("tag 演示: %v\n", u)

	// 用反射读取 tag（详细原理见第 15 课反射）
	t := reflect.TypeOf(u)
	if f, ok := t.FieldByName("Name"); ok {
		fmt.Printf("Name 字段的 tag: %q\n", f.Tag)
		fmt.Println("json 名:", f.Tag.Get("json")) // name
	}
}

// ============================================================
// 课后练习：
// 1. 给 Person 加方法 String() string，让 Printf 用 %v 打印时更友好。
//    （提示：这叫实现 fmt.Stringer 接口，第 09 课会讲）
// 2. 设计 Rectangle 和 Circle，都挂 Area() 方法，体会为什么
//    大结构体适合用指针接收者。
// 3. 用嵌入实现 Employee 包含 Person，再加 Salary 字段。
// ============================================================
