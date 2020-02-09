package quelazyexe

// package main

import (
	"fmt"
	"sync"
	"time"
)

type Executor struct {
	lock     sync.Mutex
	IsOpen   bool
	Continue bool
	Exec     func()
}

// 执行器执行过程
func (e *Executor) taskExecutor() {
	z.Lock()
	s++
	z.Unlock()
	for {
		e.lock.Lock()
		if e.Continue {
			e.Continue = false
			e.lock.Unlock()
		} else {
			e.IsOpen = false
			e.lock.Unlock()
			break
		}

		e.Exec()

		time.Sleep(time.Millisecond) // 每次延时 1 ms
	}
}

// 执行器唤醒过程
func (e *Executor) TaskSignIn() {
	e.lock.Lock()
	e.Continue = true
	if !e.IsOpen {
		e.IsOpen = true
		go e.taskExecutor()
	}
	e.lock.Unlock()
}

// 生成 执行器 封装
func NewExecutor(f func()) *Executor {
	return &Executor{IsOpen: false, Continue: false, Exec: f}
}

// ======== 测试部分 ========
var target = 10000

var p = 0 // 将被累加一 target 次
var c = 0 // 将在懒执行器中被同步为 p
var s = 0 // 执行器被启动的次数
// 每次执行进行一次同步操作
func example() {
	y.Lock()
	c = p
	// 如果达到指定值，报告, 报告最多 2 次(多做 1 次同步)
	if c == target {
		t := time.Now()
		fmt.Println("done")
		fmt.Println(t.Sub(tStart))
	}
	y.Unlock()
}

// 测试用锁
var y sync.Mutex     // 显式累加 p 的操作锁
var z sync.Mutex     // 执行器执行次数 s 的操作锁
var tStart time.Time // 程序开始时间

// 测试
func main() {
	e := NewExecutor(example)

	tStart = time.Now()
	for range make([]bool, target) {
		go func() {
			y.Lock()
			p++
			y.Unlock()
			e.TaskSignIn()
		}()
	}

	time.Sleep(time.Second * 1)
	fmt.Println("同步结果：", c, "执行器启动次数：", s)
}
