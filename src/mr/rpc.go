package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

// 定义 worker 状态
//
type WorkerStatus int
const (
	Idle WorkerStatus = iota
	Finished
	Failed
)

// 定义 task 状态
//
type TaskStatus int
const (
	NotStart TaskStatus = iota
	Doing
	Done
)

// 定义 task 阶段
//
type Phase int
const (
	MapPhase Phase = iota
	ReducePhase
)


type CallArgs struct {
	CurStatus WorkerStatus // 当前 worker 状态
	TaskIdx   int          // 当前 worker task id
}

type CallReply struct {
	CurPhase     Phase  // 当前处于 Map 阶段还是 Reduce 阶段
	MapFile      string // map input file name
	TaskDone     bool   // 是否全部 map 任务已完成
	MapNumber    int    // Map 任务总数
	ReduceNumber int    // Reduce 任务总数
	TaskIdx      int    // 任务 id
}

func (status WorkerStatus) String() string {
	switch status {
	case Idle: // 闲置状态, 代表需要被分配 task
		return "Worker is idle"
	case Finished:
		return "Worker Task Finished"
	case Failed:
		return "Worker Task Failed"
	default:
		return "Unknown Worker Status"
	}
}

func (status TaskStatus) String() string {
	switch status {
	case NotStart:
		return "Task Not Start"
	case Doing:
		return "Task Doing"
	case Done:
		return "Task Done"
	default:
		return "Unknown Task Status"
	}
}

func (phase Phase) String() string {
	switch phase {
	case MapPhase:
		return "Current in Map Phase"
	case ReducePhase:
		return "Current in Reduce Phase"
	default:
		return "Unknown Phase"
	}
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
