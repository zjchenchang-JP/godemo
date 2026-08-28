// ============================================================
// 06 - 数组、切片(Slice)、映射(Map)
//
// 运行：go run ./06_array_slice_map
// ============================================================

package main

import "fmt"

func main() {
	// =========================================================
	// 一、数组（Array）：固定长度，值类型（了解即可，实际很少直接用）
	// =========================================================
	var arr1 [3]int            // 声明：长度是类型的一部分！[3]int 和 [4]int 是不同类型
	arr1[0] = 1                // 通过下标赋值
	arr2 := [3]int{1, 2, 3}    // 字面量初始化
	arr3 := [...]int{1, 2, 3, 4} // ... 让编译器自动数长度（仍是数组）
	fmt.Println(arr1, arr2, arr3, len(arr3))

	// 【重点】数组是值类型：赋值/传参会完整拷贝！
	copyArr := arr2
	copyArr[0] = 99
	fmt.Println(arr2[0], copyArr[0]) // 1 99：原数组不受影响

	// =========================================================
	// 二、切片（Slice）：动态长度，引用类型（Go 最常用的数据结构！）
	// =========================================================

	// 创建方式
	s1 := []int{1, 2, 3}                 // 字面量
	s2 := make([]int, 5)                 // make(类型, 长度)：长度5，容量5，全为0
	s3 := make([]int, 0, 10)             // make(类型, 长度, 容量)：预留10个空间
	var s4 []int                         // nil 切片（没分配底层数组）
	fmt.Println(s1, s2, s3, s4, s4 == nil, len(s2), cap(s3))

	// len() = 当前元素个数；cap() = 底层数组容量（不够时 append 会自动扩容）
	fmt.Printf("s1 len=%d cap=%d\n", len(s1), cap(s1))

	// append 追加元素（可能触发扩容，返回【新的】切片头）
	s1 = append(s1, 4)
	s1 = append(s1, 5, 6)          // 一次追加多个
	s1 = append(s1, []int{7, 8}...) // 追加另一个切片
	fmt.Println("append 后:", s1, "cap =", cap(s1))

	// 切片表达式 s[low:high]：左闭右开 [low, high)
	sub := s1[1:4] // 取下标 1,2,3
	fmt.Println("s1[1:4] =", sub)
	fmt.Println("s1[:2] =", s1[:2], "s1[6:] =", s1[6:], "s1[:] =", s1[:])

	// copy 复制（真正的数据拷贝，目标要有足够长度）
	dst := make([]int, 3)
	copy(dst, s1) // 拷贝前 3 个
	fmt.Println("copy 后 dst =", dst)

	// 删除下标 i 的元素（Go 没有内置 remove，惯用 append 技巧）：
	i := 2
	s5 := append(s1[:i], s1[i+1:]...)
	fmt.Println("删除下标2后:", s5)

	// 【重要】切片是"引用类型"：共享底层数组，修改会互相影响！
	base := []int{1, 2, 3, 4, 5}
	view := base[1:3]
	view[0] = 99
	fmt.Println(base) // [1 99 3 4 5]：base 也被改了！

	// append 的经典陷阱：扩容后会指向新数组，不再影响原切片
	a := []int{1, 2, 3}
	b := append(a, 4) // cap 不够，b 指向新数组
	b[0] = 100
	fmt.Println(a, b) // a 不受影响

	// 切片判断空：用 len(s) == 0，不要用 s == nil
	empty := []int{}
	fmt.Println(len(empty) == 0, empty == nil) // true false

	// =========================================================
	// 三、映射（Map）：键值对，引用类型，无序
	// =========================================================

	// 创建方式
	m1 := map[string]int{"apple": 5, "banana": 3} // 字面量
	m2 := make(map[string]int)                    // make
	var m3 map[string]int                         // nil map（只读，写入会 panic！）
	fmt.Println(m1, m2, m3 == nil)

	// 增、改、查、删
	m2["go"] = 1        // 增/改：存在则覆盖
	m2["python"] = 2
	delete(m2, "python") // 删：delete(map, key)，key 不存在也不报错

	// 查：comma-ok 语法（第二个返回值表示 key 是否存在）
	count, ok := m1["apple"]
	fmt.Println("apple:", count, "存在?", ok)
	missing, ok2 := m1["cherry"]
	fmt.Println("cherry:", missing, "存在?", ok2) // 不存在时，值是零值 0

	// 【经典坑】直接取不存在的 key 得到零值，不会报错
	fmt.Println(m1["notexist"]) // 0

	// 遍历（map 是无序的！每次遍历顺序可能不同）
	for key, value := range m1 {
		fmt.Printf("  %s=%d ", key, value)
	}
	fmt.Println()
	for key := range m1 { // 只要 key
		fmt.Print(key, " ")
	}
	fmt.Println()

	// map 的值可以是任意类型：切片、map、结构体...
	graph := map[string][]string{ // 值为切片
		"a": {"b", "c"},
	}
	graph["a"] = append(graph["a"], "d")

	nested := map[string]map[string]int{ // 嵌套 map：内层需要单独 make
		"sc": {"chengdu": 1},
	}
	nested["gd"] = make(map[string]int) // 必须先初始化内层 map
	nested["gd"]["gz"] = 2
	fmt.Println(graph, nested)

	// 【重要】map 不是并发安全的！并发读写要用 sync.Map 或加锁（第 11 课）

	// 求长度
	fmt.Println(len(m1), len(graph))
}

// ============================================================
// 课后练习：
// 1. 写函数 reverse(s []int) []int 原地反转切片。
// 2. 统计一句话中每个单词出现的次数（strings.Fields + map）。
// 3. 找出切片中第二大的数（不用排序）。
// ============================================================
