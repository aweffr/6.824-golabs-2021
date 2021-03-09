package main

//
// start the coordinator process, which is implemented
// in ../mr/coordinator.go
//
// go run mrcoordinator.go pg*.txt
//
// Please do not change this file.
//

import "6.824-golabs-2021/mr"
import "time"
import "os"
import "fmt"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
		os.Exit(1)
	}

	// By 1uvu
	// 创建 Coordinator，轮询 Done()，true 说明此次运行成功
	// 注意: Coordinator 是一个 RPC Server,一旦创建后再不主动或是触发停止的情况会一直保持执行
	// todo 所以为了 Coordinator 可以进行多次任务, 再一次任务 Done 了后要刷新 Coordinator 的状态信息
	//
	c := mr.MakeCoordinator(os.Args[1:], 10)
	for c.Done() == false {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
