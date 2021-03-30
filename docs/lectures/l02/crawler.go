package main

import (
	"fmt"
	"sync"
)

//
// Serial crawler
// By 1uvu
// url 是当前待获取的 url
// fetched 保存 url 的获取状态
//
func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	if fetched[url] {
		return
	}
	fetched[url] = true // 更新当前 url 为已被获取
	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	// 用递归实现 DFS 获取 page 包含的每个 url
	for _, u := range urls {
		Serial(u, fetcher, fetched)
	}
	return
}

//
// Concurrent crawler with shared state and Mutex
// 定义获取状态
// mu 原语为修改状态上锁
// fetched 记录每个 url string 是否已被获取
//
type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

func ConcurrentMutex(url string, fetcher Fetcher, f *fetchState) {
	// 使用或修改状态时上锁
	f.mu.Lock()
	already := f.fetched[url]
	f.fetched[url] = true
	f.mu.Unlock()

	if already {
		return
	}

	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	// 并发获取当前 page 的 urls
	// done 是一个 WaitGroup, 用来存放处于阻塞状态的 goroutine 线程
	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1) // 一个 Add 需要一个 Done 来平衡
		//u2 := u
		//go func() {
		// defer done.Done()
		// ConcurrentMutex(u2, fetcher, f)
		//}()
		go func(u string) {
			defer done.Done()
			ConcurrentMutex(u, fetcher, f)
		}(u) // 带有参数的匿名函数, 启动 goroutine 来执行
	}
	done.Wait()
	return
}

func makeState() *fetchState {
	f := &fetchState{}
	f.fetched = make(map[string]bool)
	return f
}

//
// Concurrent crawler with channels
//
func worker(url string, ch chan []string, fetcher Fetcher) {
	urls, err := fetcher.Fetch(url)
	// "<-" 符号表示数据流的方向,
	// 前面是消息的接收者, 后面代表发生者
	if err != nil {
		ch <- []string{} // channel 中已无 urls 待获取
	} else {
		ch <- urls // 将 urls 放入通道
	}
}

func coordinator(ch chan []string, fetcher Fetcher) {
	n := 1 // n 记录通道中的 urls 数, ch 初始包含一个 url
	fetched := make(map[string]bool)
	for urls := range ch { // 遍历 ch, 唤起 worker goroutine 线程
		for _, u := range urls {
			if fetched[u] == false {
				fetched[u] = true
				n += 1
				// 启动线程来执行 ch, 如果不使用线程启动,
				// 会由于 for 循环的多个 worker 同时等待 ch 引发死锁
				// 尝试
				// worker(u, ch, fetcher)
				go worker(u, ch, fetcher)
			}
		}
		n -= 1
		if n == 0 {
			break
		}
	}
}

func ConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string)
	go func() {
		ch <- []string{url} // 启动一个 channel 线程
	}()
	// coordinator 无法通过线程来启动, 因为不能在线程中启动其它线程
	// 只能在当前进程中直接运行, 尝试:
	// go coordinator(ch, fetcher)
	coordinator(ch, fetcher)
}

func main() {
	fmt.Printf("=== Serial===\n")
	Serial("http://golang.org/", fetcher, make(map[string]bool))

	fmt.Printf("=== ConcurrentMutex ===\n")
	ConcurrentMutex("http://golang.org/", fetcher, makeState())

	fmt.Printf("=== ConcurrentChannel ===\n")
	ConcurrentChannel("http://golang.org/", fetcher)
}

//
// Fetcher
//

type Fetcher interface {
	// Fetch returns a slice of URLs found on the page.
	Fetch(url string) (urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
// By 1uvu
// 保存每个 page 获取到的 content 和 urls
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

// By 1uvu
// 为 fakeFetcher 类型实现 Fetcher 接口
func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("found:   %s\n", url)
		return res.urls, nil
	}
	fmt.Printf("missing: %s\n", url)
	return nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
// By 1uvu
// 初始化一个 fetcher
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
