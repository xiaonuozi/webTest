package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 创建上线列表
	var m sync.Map
	m.Store("1", "2")
	go f1(&m)
	go f2(&m)
	go f3(&m)
	go f4(&m)

	// 等待所有协程执行完毕
	// 可以使用其他方式来同步协程，这里使用Sleep暂停一段时间简化示例
	time.Sleep(time.Second)
}

func f1(m *sync.Map) {
	m.Store("2", "3")
	fmt.Println("f1: Stored key 2")
}

func f2(m *sync.Map) {
	m.Store("3", "4")
	fmt.Println("f2: Stored key 3")
}

func f3(m *sync.Map) {
	m.Store("4", "5")
	fmt.Println("f3: Stored key 4")
}

func f4(m *sync.Map) {
	fmt.Println("f4: Printing map contents")
	m.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}
