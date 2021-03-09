package mr

import (
	"fmt"
	"log"
	"sync"
	"time"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

//
//  for each task, map or reduce, only
//
var mutex sync.Mutex

type Coordinator struct {
	// 记录已分配的资源和工作进行情况
	InputFiles []string
	MapNumber int
	ReduceNumber int
	MapTaskStatusMap map[int]TaskStatus
	ReduceTaskStatusMap map[int]TaskStatus
	CurPhase Phase
	TaskDone bool
	DoneMapNumber int
	DoneReduceNumber int
}

// Your code here -- RPC handlers for the worker to call.
func (c *Coordinator) Response(args *CallArgs, reply *CallReply) error {
	switch args.CurStatus {
	case Idle:
		if c.TaskDone {
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

func (c *Coordinator) HandleMapRequire(reply *CallReply) error {
	// based the worker phase to handout task i.e. input files
	mutex.Lock()
	defer mutex.Unlock()
	for idx:=0; idx<len(c.MapTaskStatusMap); idx++ {
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
				<-timer.C
				mutex.Lock()
				defer mutex.Unlock()
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

func (c *Coordinator) HandleReduceRequire(reply *CallReply) error {
	// based the worker phase to handout task i.e. input files
	mutex.Lock()
	defer mutex.Unlock()
	for idx:=0; idx<len(c.ReduceTaskStatusMap); idx++ {
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
				mutex.Lock()
				defer mutex.Unlock()
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


func (c *Coordinator) HandleFinished(args *CallArgs) error {
	mutex.Lock()
	defer mutex.Unlock()
	switch c.CurPhase {
	case MapPhase:
		c.MapTaskStatusMap[args.TaskIdx] = Done
		c.DoneMapNumber++
		if c.DoneMapNumber == c.MapNumber {
			c.CurPhase = ReducePhase
		}
	case ReducePhase:
		c.ReduceTaskStatusMap[args.TaskIdx] = Done
		c.DoneReduceNumber++
		if c.DoneReduceNumber == c.ReduceNumber {
			c.TaskDone = true
		}
	}
	return nil
}

func (c *Coordinator) HandleFailed(args *CallArgs) error {
	mutex.Lock()
	defer mutex.Unlock()
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
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Println(c.TaskDone, c.DoneMapNumber, c.DoneReduceNumber)
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
	for idx := 0; idx < len(files); idx++ { c.MapTaskStatusMap[idx] = NotStart }
	c.ReduceTaskStatusMap = make(map[int]TaskStatus, nReduce)
	for idx := 0; idx < nReduce; idx++ { c.ReduceTaskStatusMap[idx] = NotStart }
	c.server()
	return &c
}
