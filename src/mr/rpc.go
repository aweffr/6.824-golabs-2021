package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

//
// example to show how to declare the arguments
// and reply for an RPC.
//

// By 1uvu
// 定义 worker 状态
//
type WorkerStatus int

const (
	Idle WorkerStatus = iota
	Finished
	Failed
)

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

// By 1uvu
// 定义 task 状态
//
type TaskStatus int

const (
	NotStart TaskStatus = iota
	Doing
	Done
)

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

// By 1uvu
// 定义 task 阶段
//
type Phase int

const (
	MapPhase Phase = iota
	ReducePhase
)

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
	TaskIdx      int    // 对于 Map 任务来说，需要同时得知 MapTaskIdx 和 ReduceTaskIdx，来对中间文件进行命名
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
