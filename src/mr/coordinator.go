package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)


// mutex 信号量锁住的其实不是资源而是资源的“状态”, 一份资源上一把锁
// 使得每一个获得了这份资源的锁的线程均可以确定地使用和修改资源的状态.
// 也就是说，对上锁了的资源的状态的使用和修改操作是原子的.
// 另外注意: 锁本身无意义,它保证的只是对一系列状态使用和修改操作的原子性,
// map worker 和 reduce worker 并不会一起执行, 因此此处只需要定义一个 mutex 即可.
//
type Coordinator struct {

	// 记录已分配的资源和工作进行情况
	mutex               sync.Mutex
	InputFiles          []string           // map worker 的输入文件
	MapNumber           int                // map task 总数, 一个 task 对应一个 worker
	ReduceNumber        int                // reduce task 总数, 一个 task 对应一个 worker
	MapTaskStatusMap    map[int]TaskStatus // 记录每一个 map worker 的 task 当前执行状态
	ReduceTaskStatusMap map[int]TaskStatus // 记录每一个 reduce worker 的 task 当前执行状态
	CurPhase            Phase              // 当前处于的任务阶段: map/reduce
	TaskDone            bool               // 任务是否已完成
	DoneMapNumber       int                // 已完成的 map task 总数
	DoneReduceNumber    int                // 已完成的 reduce task 总数
}

// 根据 worker 的请求 args 携带的 worker 当前状态, 来做对应的处理
//
func (c *Coordinator) Response(args *CallArgs, reply *CallReply) error {
	switch args.CurStatus {
	case Idle:
		if c.TaskDone {
			// todo fix：任务全部完成也需要通知 worker，而不能直接 return，避免未知错误
			reply.TaskDone = true
			return nil
		}
		switch c.CurPhase {
		case MapPhase:
			return c.HandleMapRequire(reply)
		case ReducePhase:
			return c.HandleReduceRequire(reply)
		}
	case Finished:
		return c.HandleFinished(args)
	case Failed:
		return c.HandleFailed(args)
	default:
		return fmt.Errorf("Invalid worker (%s %d) with %s status", c.CurPhase.String(), args.TaskIdx, args.CurStatus.String())
	}
	return nil
}


// 进行 map worker 的 task 分配工作
//
func (c *Coordinator) HandleMapRequire(reply *CallReply) error {
	c.mutex.Lock()
	defer c.mutex.Unlock() // 此函数退出后解锁信号量
	// task 分配步骤
	// 1 从 MapTaskStatusMap 顺序查找第一个 NotStart 的 task
	// 2 填充 reply
	// 3 更新 MapTaskStatusMap 当前 task 状态为 Doing
	// 4 启动一个定时器 goroutine 来实现简单的容错, 当前 task 的状态如果处于 Doing 超过 10s
	// 	则认为执行它的 map worker crash 掉了, 将其状态修改为 NotStart
	for idx := 0; idx < len(c.MapTaskStatusMap); idx++ {
		if c.MapTaskStatusMap[idx] == NotStart {
			reply.MapFile = c.InputFiles[idx]
			reply.CurPhase = MapPhase
			reply.MapNumber = c.MapNumber
			reply.ReduceNumber = c.ReduceNumber
			reply.TaskIdx = idx
			c.MapTaskStatusMap[idx] = Doing
			// fault tolerance with waiting 10s
			go func(mapTaskIdx int) {
				timer := time.NewTimer(time.Second * 10)
				<-timer.C  // 10s 后继续执行下面的代码
				c.mutex.Lock()
				defer c.mutex.Unlock()
				// timeout and then callback task
				if c.MapTaskStatusMap[mapTaskIdx] == Doing {
					c.MapTaskStatusMap[mapTaskIdx] = NotStart
				}
			}(idx)
			break
		}
	}
	return nil
}


// 与 map task 分配过程类似
func (c *Coordinator) HandleReduceRequire(reply *CallReply) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for idx := 0; idx < len(c.ReduceTaskStatusMap); idx++ {
		if c.ReduceTaskStatusMap[idx] == NotStart {
			reply.CurPhase = ReducePhase
			reply.MapNumber = c.MapNumber
			reply.ReduceNumber = c.ReduceNumber
			reply.TaskIdx = idx
			c.ReduceTaskStatusMap[idx] = Doing
			// fault tolerance with waiting 10s
			go func(reduceTaskIdx int) {
				timer := time.NewTimer(time.Second * 10)
				<-timer.C
				c.mutex.Lock()
				defer c.mutex.Unlock()
				// timeout and then callback task
				if c.ReduceTaskStatusMap[reduceTaskIdx] == Doing {
					c.ReduceTaskStatusMap[reduceTaskIdx] = NotStart
				}
			}(idx)
			break
		}
	}
	return nil
}


// task Finished 时更新 Done*Number, *TaskStatusMap 和 TaskDone
func (c *Coordinator) HandleFinished(args *CallArgs) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	switch c.CurPhase {
	case MapPhase:
		c.MapTaskStatusMap[args.TaskIdx] = Done
		c.DoneMapNumber++
		if c.DoneMapNumber == c.MapNumber {
			c.CurPhase = ReducePhase
		}
		// todo fea：在此处为 worker 实现 Backup 任务，
		// 在 test 时候不允许 worker 额外运行其他任务，
		// 因此 Backup 任务只能在 worker 结束了一段时间后运行，如 10s。
	case ReducePhase:
		c.ReduceTaskStatusMap[args.TaskIdx] = Done
		c.DoneReduceNumber++
		if c.DoneReduceNumber == c.ReduceNumber {
			c.TaskDone = true
		}
	}
	return nil
}


// task Failed 时更新 *TaskStatusMap
func (c *Coordinator) HandleFailed(args *CallArgs) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	switch c.CurPhase {
	case MapPhase:
		c.MapTaskStatusMap[args.TaskIdx] = NotStart
	case ReducePhase:
		c.ReduceTaskStatusMap[args.TaskIdx] = NotStart
	}
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	_ = os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	//fmt.Println(c.TaskDone, c.DoneMapNumber, c.DoneReduceNumber)
	ret := false
	ret = c.TaskDone
	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	c.CurPhase = MapPhase
	c.InputFiles = files
	c.MapNumber = len(files)
	c.ReduceNumber = nReduce
	c.DoneMapNumber = 0
	c.DoneReduceNumber = 0
	c.TaskDone = false
	c.MapTaskStatusMap = make(map[int]TaskStatus, len(files))
	for idx := 0; idx < len(files); idx++ {
		c.MapTaskStatusMap[idx] = NotStart
	}
	c.ReduceTaskStatusMap = make(map[int]TaskStatus, nReduce)
	for idx := 0; idx < nReduce; idx++ {
		c.ReduceTaskStatusMap[idx] = NotStart
	}
	c.server()
	return &c
}
