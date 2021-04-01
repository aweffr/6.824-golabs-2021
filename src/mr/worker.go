package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
	"time"
)

// todo 1 worker 实现 Backup: mr-wc-all-initial
// todo 2 use sync.Cond to replace time.Sleep, and use sync channel to replace the mutex, passed the test-mr-many.sh
// todo 3 improve the code based the official solution

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}
// sorted by Key, 实现 sort.Sort() 的 data 接口
type ByKey []KeyValue
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// 根据 ihash(key) % NReduce 的值将生成的每一个中间 kv 分配给对应的 reduce
// (以存入 intermediate file named [mr-mapTaskIdx-reduceTaskIdx] 的形式)
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// 保存每个中间文件的 json Encoder
var intermediateEncoderMap map[int]*json.Encoder = nil

// 保存每个中间文件的 file name
var intermediateFileNameMap map[int]string = nil

// 根据 coordinator 的响应 reply 携带的 worker 当前状态及其他信息, 来做对应的处理
//
func Worker(mapf func(string, string) []KeyValue, reducef func(string, []string) string) {
	for {
		reply := Request()
		if reply.TaskDone {
			break
		}
		result := true // 当前 task 是否正常完成
		switch reply.CurPhase {
		case MapPhase:
			result = MapTask(reply, mapf)
		case ReducePhase:
			result = ReduceTask(reply, reducef)
		}
		var curStatus WorkerStatus
		if result {
			curStatus = Finished
		} else {
			curStatus = Failed
		}
		NotifyStatus(curStatus, reply.TaskIdx)
		time.Sleep(time.Second)
	}
}

// 请求 task
func Request() CallReply {
	args := CallArgs{}
	args.CurStatus = Idle
	reply := CallReply{}
	call("Coordinator.Response", &args, &reply)
	return reply
}

// 通知状态
func NotifyStatus(curStatus WorkerStatus, idx int) {
	args := CallArgs{}
	args.CurStatus = curStatus
	args.TaskIdx = idx
	reply := CallReply{}
	call("Coordinator.Response", &args, &reply)
}

func MapTask(reply CallReply, mapf func(string, string) []KeyValue) bool {
	InitialReduceTask(reply)
	fileName := reply.MapFile
	file, err := os.Open(fileName)
	if err != nil {
		os.Exit(1)
		return false
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", fileName)
		return false
	}
	defer file.Close()

	kva := mapf(fileName, string(content))
	success := true
	for i := 0; i < len(kva); i++ {
		reduceTaskIdx := ihash(kva[i].Key) % reply.ReduceNumber
		intermediateEncoder := intermediateEncoderMap[reduceTaskIdx]
		if err := intermediateEncoder.Encode(kva[i]); err != nil {
			success = false
			intermediateEncoderMap = nil
			intermediateFileNameMap = nil
			log.Fatalf("encode kv:%v failed", kva[i])
			return false
		}
	}
	if success {
		for reduceIdx, fileName := range intermediateFileNameMap {
			_ = os.Rename(fileName, fmt.Sprintf("mr-%+v-%+v", reply.TaskIdx, reduceIdx))
		}
		intermediateEncoderMap = nil
		intermediateFileNameMap = nil
	}

	return true
}

func InitialReduceTask(reply CallReply) {
	if intermediateEncoderMap == nil {
		intermediateEncoderMap = make(map[int]*json.Encoder)
		intermediateFileNameMap = make(map[int]string)
		for j := 0; j < reply.ReduceNumber; j++ {
			fileName := fmt.Sprintf("mr-tmp-%+v-%+v", reply.TaskIdx, j)
			file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("cannot open intermediateEncoderFile:%v mapIdx:%+v reduceIdx:%+v", fileName, reply.TaskIdx, j)
			}
			intermediateEncoderMap[j] = json.NewEncoder(file)
			intermediateFileNameMap[j] = fileName
		}
	}
}

func ReduceTask(reply CallReply, reducef func(string, []string) string) bool {
	var intermediate []KeyValue
	for i := 0; i < reply.MapNumber; i++ {
		fileName := fmt.Sprintf("mr-%+v-%+v", i, reply.TaskIdx)
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("cannot open intermediateDecoderFile:%v mapIdx:%+v reduceIdx:%+v", fileName, i, reply.TaskIdx)
		}
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}
	}
	sort.Sort(ByKey(intermediate))

	oname := fmt.Sprintf("mr-out-%+v", reply.TaskIdx)
	tname := fmt.Sprintf("mr-tmp-reduce")
	ofile, err := os.OpenFile(tname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("cannot open temp fileName:%v", tname)
		return false
	}
	defer ofile.Close()

	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		_, _ = fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}
	_ = os.Rename(tname, oname)
	for i := 0; i < reply.MapNumber; i++ {
		fileName := fmt.Sprintf("mr-+%v-%+v", i, reply.TaskIdx)
		err := os.Remove(fileName)
		if err != nil {
			log.Fatalf("cannot remove tmp file:%v", fileName)
		}
	}
	return true
}

func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
