// ============================================================
// 15 - 反射（reflect）：运行时检查和操作类型
//
// 运行：go run ./15_reflection
// ============================================================

package main

import (
	"fmt"
	"reflect"
)

// 反射的三大用途：
//  1. 读写任意结构体的字段和 tag（encoding/json、ORM、依赖注入）
//  2. 动态调用方法
//  3. 泛型出现前的通用容器（现在多数场景可用泛型替代）
//
// 缺点：慢（比直接调用慢一个量级）、绕过编译期检查、代码难读。
// 原则：能不用就不用；需要"框架级通用性"时才用。

// User 演示用结构体：带 json tag
type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"` // omitempty：零值时省略
	phone string // 私有字段：反射只能读不能改
}

// Hello 演示用方法
func (u *User) Hello(who string) string {
	return "你好 " + who + "，我是 " + u.Name
}

func main() {
	u := User{Name: "Alice", Age: 20, phone: "123"}

	// ---------- 1. TypeOf / ValueOf：反射的两大入口 ----------
	t := reflect.TypeOf(u)  // 类型信息（reflect.Type）
	v := reflect.ValueOf(u) // 值信息（reflect.Value）

	fmt.Println("类型名:", t.Name())     // User
	fmt.Println("种类 Kind:", t.Kind()) // struct
	fmt.Println("字段数:", t.NumField())
	fmt.Println("值种类:", v.Kind(), "可寻址?", v.CanAddr()) // 传值拷贝，不可寻址

	// ---------- 2. 遍历结构体字段（反射最经典的用法）----------
	fmt.Println("\n遍历 User 字段：")
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)        // 字段的类型信息 StructField
		fv := v.Field(i)       // 字段的值
		fmt.Printf("  %-6s 类型=%-7v json标签=%-18q 可导出=%-6v 值=%v\n",
			f.Name, f.Type, f.Tag.Get("json"), f.IsExported(), fv)
	}

	// ---------- 3. 修改值：必须传指针 ----------
	// reflect.ValueOf(u) 拿到的是副本，改不了；
	// 要修改必须传地址，再用 Elem() 解一层指针。
	pv := reflect.ValueOf(&u)
	fmt.Println("\n指针种类:", pv.Kind()) // ptr
	ev := pv.Elem()                      // 等价于 *u，可寻址、可修改
	fmt.Println("Elem 后可寻址?", ev.CanAddr(), "可设置?", ev.CanSet())

	ev.FieldByName("Age").SetInt(21) // 把 u.Age 改成 21
	fmt.Println("反射修改后 Age =", u.Age)

	// 【注意】修改未导出字段会 panic：
	// ev.FieldByName("phone").SetString("456") // panic: reflect: reflect.Value.SetString using value obtained using unexported field

	// ---------- 4. 动态创建与判断类型 ----------
	inspect(3.14)
	inspect("hello")
	inspect([]int{1, 2})
	inspect(u)
	inspect(func() {})

	// ---------- 5. 动态调用方法 ----------
	mv := reflect.ValueOf(&u).MethodByName("Hello") // 指针才能拿到指针方法集
	if mv.IsValid() {
		// Call 的参数和返回值都是 []reflect.Value
		results := mv.Call([]reflect.Value{reflect.ValueOf("Go")})
		fmt.Println("反射调用 Hello:", results[0].String())
	}

	// ---------- 6. StructTag 的语法 ----------
	// tag 是一个字符串，格式：`key1:"value1" key2:"value2"`
	// f.Tag.Get("json") 取单个 key；Lookup 能判断是否存在
	if f, ok := reflect.TypeOf(u).FieldByName("Email"); ok {
		jsonName, exists := f.Tag.Lookup("json")
		fmt.Println("Email tag:", jsonName, exists)
	}

	// ---------- 7. 慎用反射的检查清单 ----------
	// - 性能敏感的热路径？ -> 不要用
	// - 编译期就知道类型？ -> 用类型断言/泛型代替
	// - 真要做通用序列化/ORM？ -> 参考 encoding/json 的实现
}

// inspect 用 Kind 判断任意值的实际类别
func inspect(x any) {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Float64:
		fmt.Printf("float64: %.2f\n", v.Float())
	case reflect.String:
		fmt.Printf("string: %q 长度=%d\n", v.String(), v.Len())
	case reflect.Slice:
		fmt.Printf("slice: 长度=%d 元素类型=%v\n", v.Len(), v.Type().Elem())
	case reflect.Struct:
		fmt.Printf("struct: %s 有 %d 个字段\n", v.Type().Name(), v.NumField())
	case reflect.Func:
		fmt.Println("func: 函数类型不能直接调用（需 Call）")
	default:
		fmt.Printf("其他: %v (%v)\n", v.Kind(), v.Type())
	}
}

// ============================================================
// 课后练习：
// 1. 写函数 PrintJSONNames(v any)：用反射打印结构体所有字段的
//    json 标签名（这是 ORM/序列化库的核心逻辑）。
// 2. 对比 reflect.TypeOf(x).Kind() 和 %T 的区别。
// 3. 思考：为什么 reflect.ValueOf(x) 不能修改 x？（提示：值拷贝）
// ============================================================
