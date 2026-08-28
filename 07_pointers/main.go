// ============================================================
// 07 - 指针：& 取地址，* 解引用
//
// 运行：go run ./07_pointers
// ============================================================

package main

import "fmt"

func main() {
	// ---------- 1. 指针基础 ----------
	// 指针 = 存放"变量内存地址"的变量
	// &x  ：取变量 x 的地址
	// *p  ：解引用，访问指针 p 指向的变量
	a := 10
	p := &a          // p 的类型是 *int（指向 int 的指针）
	fmt.Println("a 的地址:", p)
	fmt.Println("p 指向的值:", *p) // 10

	*p = 20          // 通过指针修改 a 的值
	fmt.Println("修改后 a =", a)  // 20

	fmt.Printf("p 的类型: %T\n", p) // *int

	// ---------- 2. 指针的零值是 nil ----------
	var np *int
	fmt.Println("空指针:", np == nil)
	// fmt.Println(*np) // 取消注释会 panic：nil pointer dereference

	// ---------- 3. 指针作为函数参数（实现"引用传递"效果）----------
	// Go 函数参数都是【值传递】（拷贝一份）。
	// 想在函数内修改外部变量，就传它的地址。
	x, y := 1, 2
	failSwap(x, y) // 传值：交换了个寂寞
	fmt.Println("failSwap 后:", x, y) // 1 2
	realSwap(&x, &y) // 传地址：真的交换了
	fmt.Println("realSwap 后:", x, y) // 2 1

	// ---------- 4. 指针与结构体（最常用的场景）----------
	user := User{Name: "小明", Age: 18}
	up := &user // 指向结构体的指针

	// 【语法糖】通过指针访问字段不需要写 (*up).Name，直接 up.Name
	up.Age = 19
	fmt.Println(user.Age) // 19

	// ---------- 5. new 与 & 的对比 ----------
	// new(T) 分配一个 T 类型的零值内存，返回指针（*T）
	n1 := new(int)    // *int，指向 0
	*n1 = 100
	n2 := new(User)   // *User，零值结构体
	n2.Name = "小红"

	// 更常见的方式：直接对字面量取地址
	n3 := &User{Name: "小刚", Age: 20}

	fmt.Println(*n1, n2, n3)

	// ---------- 6. Go 与 C 指针的区别 ----------
	// 1) 没有 +1/-1 这类【指针运算】：p++ 是编译错误
	// 2) 有垃圾回收（GC），不需要手动 free/delayed delete
	// 3) 函数不能返回局部变量的地址吗？——可以！GC 会自动延长其生命周期
	addr := returnLocalPtr()
	fmt.Println("返回局部变量指针:", *addr) // 42，安全！

	// ---------- 7. 指针的指针（了解即可）----------
	pp := &p // **int，指向 p（p 又指向 a）
	fmt.Println("**pp =", **pp) // 20，两层解引用

	// ---------- 8. 什么时候用指针？----------
	// 1) 需要修改调用方的变量/结构体
	// 2) 结构体较大，拷贝成本高
	// 3) 需要表达"可能为空"（nil）
	// 注意：切片、map 本身就带"引用"性质，传值即可修改内容（见第 06 课）

	s := []int{1, 2, 3}
	modifySlice(s) // 切片传值也能改（共享底层数组）
	fmt.Println("切片被函数修改:", s)
}

// User 演示用结构体
type User struct {
	Name string
	Age  int
}

// failSwap 传值：函数内拿到的是副本，改了也没用
func failSwap(a, b int) {
	a, b = b, a
}

// realSwap 传指针：通过地址直接修改外部变量
func realSwap(a, b *int) {
	*a, *b = *b, *a
}

// returnLocalPtr 返回局部变量的地址是安全的（GC 会追踪引用）
func returnLocalPtr() *int {
	local := 42
	return &local
}

// modifySlice 切片参数：即使传值，底层数组仍是共享的
func modifySlice(s []int) {
	s[0] = 999
}

// ============================================================
// 课后练习：
// 1. 写函数 setAge(u *User, age int) 修改用户年龄。
// 2. 思考：为什么 failSwap 失败而 realSwap 成功？
// 3. 思考：modifySlice 没传指针为什么能改成功？如果 append 呢？
//    （提示：append 扩容后指向新数组，修改不会影响原切片）
// ============================================================
