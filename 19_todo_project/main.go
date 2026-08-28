// ============================================================
// 19 - 综合项目：命令行 Todo 应用
//
// 综合运用：结构体+方法、切片、map、错误处理、JSON 持久化、指针接收者
//
// 用法（在项目根目录执行）：
//   go run ./19_todo_project add "学习 Go 基础"
//   go run ./19_todo_project add "写一个并发程序"
//   go run ./19_todo_project list
//   go run ./19_todo_project done 1
//   go run ./19_todo_project rm 2
//   go run ./19_todo_project help
// ============================================================

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// 数据文件路径（与运行时的当前目录同级）
const dataFile = "todos.json"

// Todo 一条待办事项
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoList 待办清单：持有所有数据和方法
type TodoList struct {
	Todos  []Todo `json:"todos"`
	nextID int    // 私有字段不参与序列化
}

// ---------- 面向对象式的业务方法 ----------

// Add 新增待办并返回它的编号
func (l *TodoList) Add(title string) int {
	l.nextID++
	l.Todos = append(l.Todos, Todo{
		ID:        l.nextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	})
	return l.nextID
}

// Find 按编号查找待办（指针接收者：可能修改 Done 字段）
func (l *TodoList) Find(id int) *Todo {
	for i := range l.Todos {
		if l.Todos[i].ID == id {
			return &l.Todos[i] // 返回切片元素的指针，可直接修改
		}
	}
	return nil
}

// MarkDone 把指定待办标记为完成
func (l *TodoList) MarkDone(id int) bool {
	t := l.Find(id)
	if t == nil {
		return false
	}
	t.Done = true
	return true
}

// Remove 删除指定待办（切片删除技巧：见第 06 课）
func (l *TodoList) Remove(id int) bool {
	for i, t := range l.Todos {
		if t.ID == id {
			l.Todos = append(l.Todos[:i], l.Todos[i+1:]...)
			return true
		}
	}
	return false
}

// Print 打印全部待办
func (l *TodoList) Print() {
	if len(l.Todos) == 0 {
		fmt.Println("（空空如也，用 add 命令添加一条吧）")
		return
	}
	fmt.Printf("%-4s %-3s %-24s %s\n", "ID", "状态", "标题", "创建时间")
	for _, t := range l.Todos {
		status := "[ ]"
		if t.Done {
			status = "[✓]"
		}
		fmt.Printf("%-4d %-3s %-24s %s\n",
			t.ID, status, t.Title, t.CreatedAt.Format("2006-01-02 15:04"))
	}
	// 统计（闭包的简单应用）
	doneCount := func() int {
		n := 0
		for _, t := range l.Todos {
			if t.Done {
				n++
			}
		}
		return n
	}()
	fmt.Printf("共 %d 条，已完成 %d 条\n", len(l.Todos), doneCount)
}

// ---------- JSON 持久化（错误处理的综合练习）----------

// Load 从文件读取清单；文件不存在时返回空清单（不算错误）
func Load(path string) (*TodoList, error) {
	l := &TodoList{}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) { // 首次运行没有文件是正常情况
		return l, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取 %s 失败: %w", path, err)
	}
	if err := json.Unmarshal(data, l); err != nil {
		return nil, fmt.Errorf("解析 %s 失败: %w", path, err)
	}
	// 恢复自增 ID：找到历史最大编号
	for _, t := range l.Todos {
		if t.ID > l.nextID {
			l.nextID = t.ID
		}
	}
	return l, nil
}

// Save 把清单写回文件
func (l *TodoList) Save(path string) error {
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// ---------- 命令行入口 ----------

func usage() {
	fmt.Println(`命令行 Todo 应用
用法：
  todo add <标题>    添加待办
  todo list          查看全部
  todo done <编号>   标记完成
  todo rm <编号>     删除待办`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	cmd := os.Args[1]

	list, err := Load(dataFile)
	if err != nil {
		fmt.Println("加载数据失败:", err)
		os.Exit(1)
	}

	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("请提供标题：todo add <标题>")
			return
		}
		id := list.Add(os.Args[2])
		fmt.Printf("已添加 #%d: %s\n", id, os.Args[2])

	case "list":
		list.Print()
		// list 不修改数据，直接退出
		return

	case "done":
		id := parseID(os.Args)
		if id == -1 {
			return
		}
		if !list.MarkDone(id) {
			fmt.Printf("找不到编号为 %d 的待办\n", id)
		} else {
			fmt.Printf("已完成 #%d\n", id)
		}

	case "rm":
		id := parseID(os.Args)
		if id == -1 {
			return
		}
		if !list.Remove(id) {
			fmt.Printf("找不到编号为 %d 的待办\n", id)
		} else {
			fmt.Printf("已删除 #%d\n", id)
		}

	case "help":
		usage()
		return

	default:
		fmt.Println("未知命令:", cmd)
		usage()
		return
	}

	// 只有会修改数据的命令才落盘
	if err := list.Save(dataFile); err != nil {
		fmt.Println("保存失败:", err)
		os.Exit(1)
	}
}

// parseID 从命令行参数解析编号，失败返回 -1
func parseID(args []string) int {
	if len(args) < 3 {
		fmt.Println("请提供编号")
		return -1
	}
	id, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("无效的编号 %q（必须是整数）\n", args[2])
		return -1
	}
	return id
}
