package mr

import "fmt"
import "log"
import "net/rpc"
import "hash/fnv"
import "encoding/json"


//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

var waitSeconds int
var intermediateEncoderMap map[int]json.Encoder
var intermediateFileMap map[int]string

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}


//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.
	
	// uncomment to send the Example RPC to the coordinator.
	curOpt := Require
	for {
		reply := RequireTask(curOpt)
		if reply.Done {
			break
		}
		
		switch reply.CurPhase {
		case MapPhase:
			// do map task
			fmt.Println("do map task")
			MapTask(reply, mapf)
		case ReducePhase:
			// do reduce task
			fmt.Println("do reduce task")
		}
		time.Sleep(waitSeconds)
	}
}

func RequireTask(curOpt Opt) CallReply {
	args := CallArgs{}
	args.CurOpt = curOpt
	reply := CallReply{}
	
	call("Coordinator.HandOutTask", &args, &reply)
	fmt.Println("Args: %v, Reply: %v", &args.CurOpt, reply)
	return reply
}

func MapTask(CallReply reply, mapf func(string, string) []KeyValue) bool {
	// 读入文件，，使用 mapf 生成中间数据
	fileName := reply.MapFile
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Can't open file %v", fileName)
	}
	content, err := ioutil.ReadAll(file)
        if err != nil {
                log.Fatalf("Can't read file %v", fileName)
        }
	defer file.Close()
	kva := mapf(filename, string(content))
	
	// 中间数据划分并做持久化处理
	succ := true
	for i, kv := range kva {
		// 划分 kva，并生成 reduce task 的 id
		redecuTaskIdx := ihash(kv.Key) % reply.ReduceNumber
		intermediateEncoder := json.NewEncoder(file)
		intermediateEncoderMap[reduceTaskIdx] = intermediateEncoder
		if err := intermediateEncoder.Encode(&kv); err != nil {
			log.Fatalf("Encode kv: %v failed", kv)
			succ = false
			break
		}
	}

	if succ {
		for reduceTaskIdx, fileName := range intermediateFileMap {
			os.Rename(fileName, fmt.Sprintf("mr-%+v-%+v", reply.MapTaskIdx, reduceTaskidx))
	}
	return succ
}


func ReduceTask(reply CallReply, reducef func(string, []string) string) bool {
	
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
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
