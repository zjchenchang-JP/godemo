# Go 语言学习示例集

一套可以直接运行的 Go 语言学习代码，每个目录一个主题，代码内含详细的中文注释。
所有示例已在 **go 1.25** 环境下验证可编译运行。

## 如何运行

在项目根目录（`godemo/`）下执行：

```bash
go run ./01_hello_world      # 运行单个示例
go build ./...               # 编译全部示例，验证无误
go vet ./...                 # 静态检查
```

测试类示例（18_testing）用：

```bash
go test ./18_testing -v              # 运行单元测试
go test ./18_testing -bench=. -v     # 运行基准测试
```

## 学习顺序

### 阶段一：入门基础（先跑起来，熟悉语法）
| 顺序 | 目录 | 主题 | 重点 |
|---|---|---|---|
| 1 | [01_hello_world](01_hello_world/) | 程序结构 | package、import、main 函数、fmt 输出、注释 |
| 2 | [02_variables](02_variables/) | 变量与常量 | var、:=、零值、const、iota 枚举 |
| 3 | [03_basic_types](03_basic_types/) | 基本类型 | 整型/浮点/字符串/rune、类型转换、Printf 动词 |
| 4 | [04_control_flow](04_control_flow/) | 流程控制 | if、for（Go 只有 for）、switch、break/continue/label、goto |

### 阶段二：核心语法（Go 的日常写法）
| 顺序 | 目录 | 主题 | 重点 |
|---|---|---|---|
| 5 | [05_functions](05_functions/) | 函数 | 多返回值、可变参数、闭包、defer、init |
| 6 | [06_array_slice_map](06_array_slice_map/) | 集合类型 | 数组 vs 切片、append/copy、切片底层与陷阱、map |
| 7 | [07_pointers](07_pointers/) | 指针 | & 和 *、指针传参、new、nil 指针 |
| 8 | [08_struct_method](08_struct_method/) | 结构体与方法 | struct、值/指针接收者、嵌入(组合)、构造函数 |

### 阶段三：接口与错误（Go 的"面向对象"哲学）
| 顺序 | 目录 | 主题 | 重点 |
|---|---|---|---|
| 9 | [09_interface](09_interface/) | 接口 | 隐式实现、多态、类型断言/类型开关、any、Stringer |
| 10 | [10_error_panic](10_error_panic/) | 错误处理 | error、%w 包装、errors.Is/As、自定义错误、panic/recover |

### 阶段四：并发（Go 的杀手锏）
| 顺序 | 目录 | 主题 | 重点 |
|---|---|---|---|
| 11 | [11_goroutine](11_goroutine/) | goroutine | go 关键字、WaitGroup、Mutex、Once、atomic、竞态 |
| 12 | [12_channel](12_channel/) | channel 通道 | 无缓冲/有缓冲、close/range、select、超时、worker pool |

### 阶段五：现代特性与工程化
| 顺序 | 目录 | 主题 | 重点 |
|---|---|---|---|
| 13 | [13_generics](13_generics/) | 泛型 | 类型参数、约束(any/comparable/自定义)、泛型容器 |
| 14 | [14_package_module](14_package_module/) | 包与模块 | 包的可见性、init 执行顺序、go.mod、第三方依赖 |
| 15 | [15_reflection](15_reflection/) | 反射 | TypeOf/ValueOf、遍历字段、读写值、struct tag |

### 阶段六：标准库与实战
| 顺序 | 目录 | 主题 | 重点 |
|---|---|---|---|
| 16 | [16_json_io](16_json_io/) | IO 与 JSON | 文件读写、bufio.Scanner、encoding/json |
| 17 | [17_net_http](17_net_http/) | 网络编程 | http.Server、路由处理、JSON API、客户端请求 |
| 18 | [18_testing](18_testing/) | 测试 | 单元测试、表驱动测试、基准测试、Example |
| 19 | [19_todo_project](19_todo_project/) | 综合项目 | 命令行 Todo：结构体+方法+JSON 持久化+错误处理 |

## 学习建议

1. **按顺序学**：每个文件都是独立可运行的 `main` 程序，先看注释，再改代码跑一跑。
2. **动手改**：每个文件末尾都有"课后练习"，做完再进入下一课。
3. **善用工具**：
   - `go fmt ./...` 格式化代码（Go 有官方统一风格）
   - `go doc fmt.Println` 查看任何函数的文档
   - 在 [pkg.go.dev](https://pkg.go.dev) 查标准库文档
4. **并发篇重点学**：goroutine + channel 是 Go 区别于其他语言的核心，值得多花时间。

## 延伸资源

- [A Tour of Go（中文版交互教程）](https://tour.go-zh.org/)
- [Go by Example（示例驱动学习）](https://gobyexample.com/)
- [Go 官方文档](https://go.dev/doc/)
- 《The Go Programming Language》（"圣书"，进阶必读）
- [Effective Go（官方代码风格指南）](https://go.dev/doc/effective_go)
