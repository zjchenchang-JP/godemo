// ============================================================
// 14 - 包(Package)与模块(Module)：代码组织、可见性、依赖管理
//
// 运行：go run ./14_package_module
// ============================================================

package main

import (
	"fmt"

	"godemo/14_package_module/mathutil" // 导入本项目内的包：模块名/路径/包名
)

// 【一个模块(module)的结构】
//
//	godemo/                  <- 模块根目录（有 go.mod）
//	├── go.mod               <- 模块定义文件：module godemo + go 版本
//	├── 01_hello_world/      <- package main（可执行）
//	├── ...
//	└── 14_package_module/
//	    ├── main.go          <- package main
//	    └── mathutil/        <- package mathutil（库包）
//
// 【import 路径 = 模块名 + 目录路径】，与包名（package xxx）可以不同，
// 但惯例上保持一致。导入后用"包名.成员"访问：mathutil.Add(...)。

func main() {
	fmt.Println("== 导入自定义包 mathutil ==")

	// 大写开头 = 导出（公开）：包外可用
	fmt.Println("Add(3, 5) =", mathutil.Add(3, 5))
	fmt.Println("Square(4) =", mathutil.Square(4))
	fmt.Println("SumAll(1,2,3,4,5) =", mathutil.SumAll(1, 2, 3, 4, 5))
	fmt.Println("Version =", mathutil.Version, "累计调用 =", mathutil.Calls())

	// 小写开头 = 未导出（私有）：包外不可见，取消注释会编译错误：
	// fmt.Println(mathutil.addOne(1, 2)) // undefined: mathutil.addOne

	fmt.Println()
	fmt.Println("== 可见性规则总结 ==")
	fmt.Println("标识符（变量/函数/类型/字段/方法）首字母：")
	fmt.Println("  大写 -> 导出，任何导入者都能访问")
	fmt.Println("  小写 -> 未导出，仅本包内可见")

	fmt.Println()
	fmt.Println("== go.mod 与依赖管理（模块时代，Go 1.11+）==")
	fmt.Println("初始化模块：  go mod init 模块名        （如 go mod init github.com/you/godemo）")
	fmt.Println("添加第三方包： go get github.com/xxx/yyy  （自动写入 go.mod，校验存入 go.sum）")
	fmt.Println("整理依赖：    go mod tidy                 （移除没用到的依赖）")
	fmt.Println("下载依赖：    go mod download")
	fmt.Println("查看依赖：    go list -m all")

	fmt.Println()
	fmt.Println("== internal 包（强制私有）==")
	fmt.Println("路径中含 internal/ 的包，只有其父目录下的代码可以导入，")
	fmt.Println("外部模块强行导入会直接编译报错 —— 用来放内部实现。")

	fmt.Println()
	fmt.Println("== 包的初始化顺序（面试常问）==")
	fmt.Println("1. 按 import 依赖图，先初始化被导入的包")
	fmt.Println("2. 每个包内：包级变量初始化 -> init()（可多个，按声明顺序）")
	fmt.Println("3. 最后执行 package main 的 main()")

	fmt.Println()
	fmt.Println("== 标准布局（社区惯例 go-project-layout，了解即可）==")
	fmt.Println("cmd/xxx/     各个可执行程序入口")
	fmt.Println("internal/    私有代码")
	fmt.Println("pkg/         可被外部引用的库")
	fmt.Println("小型项目不需要严格遵守，扁平结构完全没问题（就像本仓库）")
}

// ============================================================
// 课后练习：
// 1. 在本目录新建 stringutil 包，写 Reverse(s string) string。
// 2. 故意调用 mathutil 的私有函数，观察编译错误信息。
// 3. 执行 go list -m all，看看当前模块的依赖情况。
// ============================================================
