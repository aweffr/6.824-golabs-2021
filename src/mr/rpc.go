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

type Status int
const (
	Require Status = iota
	Finished
	Failed
)
func (status Status) String() string {
	switch status {
	case Require:
		return "Require Task"
	case Finished:
		return "Task Finished"
	case Failed:
		return "Task Failed"
	default:
		return "Unknown Option"
	}
}

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
	/* call 操作请求类型，包括：
	CurStatus
	1 索要任务 Require
	2 通知任务完成 Finished
	3 notify task failed Failed
	CurPhase
	1 Map
	2 Reduce
	*/
	CurStatus Status
	CurPhase Phase
	TaskIdx int
}

type CallReply struct {
	CurPhase Phase  // 当前处于 Map 阶段还是 Reduce 阶段
	MapFile string  // map input file name
	Done bool  // 是否全部任务已完成
	MapNumber int  // Map 任务总数
	ReduceNumber int  // Reduce 任务总数
	MapTaskIdx int  // 对于 Map 任务来说，需要同时得知 MapTaskIdx 和 ReduceTaskIdx，来对中间文件进行命名
	ReduceTaskIdx int  // 可用来对输出文件命名
}

// Add your RPC definitions here.


// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
