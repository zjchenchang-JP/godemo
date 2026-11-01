package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	// s := make([]int, 5)
	// case1(s)
	// case2(s)
	// PrintSliceStruct(&s)
	fmt.Println("=========================")
	doAppend := func (s []int)  {
		s = append(s, 1)
		printLengthAndCapacity(s)
	}
	s1 := make([]int, 8, 8)
	doAppend(s1[:4])
	// printLengthAndCapacity(s1[:4])
	// fmt.Printf("s1 = %v\n", s1)
	printLengthAndCapacity(s1)
	doAppend(s1)
	printLengthAndCapacity(s1)

}

func printLengthAndCapacity(s []int) {
	fmt.Println(s)
	fmt.Printf("len = %d cap = %d \n", len(s), cap(s))
}

func case2(s []int) {
	s = append(s, 0)
	s[1] = 1
	PrintSliceStruct(&s)
}

func PrintSliceStruct(s *[]int) {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(s))
	// 查看结构
	fmt.Printf("slice struct: %+v, slice is %v\n", ss, s)
}

// 底层数组不变
func case1(s []int) {
	s[1] = 1
	PrintSliceStruct(&s)
}