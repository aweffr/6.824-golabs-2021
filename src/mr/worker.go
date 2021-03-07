package mr

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)
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

// sorted by Key, so implement the Sort data interface
type ByKey []KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

var waitSeconds int
var intermediateEncoderMap map[int]*json.Encoder = nil
var intermediateFileNameMap map[int]string = nil

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
	for {
		reply := Request()
		if reply.Done {
			break
		}
		
		switch reply.CurPhase {
		case MapPhase:
			// do map task
			fmt.Println("do map task")
			result := MapTask(reply, mapf)
			if result {
				curStatus:= Finished
				NotifyStatus(curStatus, MapPhase, reply.MapTaskIdx)
			} else {
				curStatus := Failed
				NotifyStatus(curStatus, MapPhase, reply.MapTaskIdx)
			}
		case ReducePhase:
			// do reduce task
			fmt.Println("do reduce task")
			result := ReduceTask(reply, reducef)
			if result {
				curStatus := Finished
				NotifyStatus(curStatus, ReducePhase, reply.ReduceTaskIdx)
			} else {
				curStatus := Failed
				NotifyStatus(curStatus, ReducePhase, reply.ReduceTaskIdx)
			}
		}
		time.Sleep(time.Duration(waitSeconds))
	}
}

func Request() CallReply {
	args := CallArgs{}
	args.CurStatus = Require
	reply := CallReply{}
	call("Coordinator.Response", &args, &reply)
	fmt.Printf("Args: %v, Reply: %v", &args.CurStatus, reply)
	return reply
}

func NotifyStatus(curStatus Status, curPhase Phase, taskIdx int) CallReply {
	args := CallArgs{}
	args.CurStatus = curStatus
	args.CurPhase = curPhase
	args.TaskIdx = taskIdx
	reply := CallReply{}
	call("Coordinator.Response", &args, &reply)
	fmt.Printf("Args: %v, Reply: %v", &args.CurStatus, reply)
	return reply
}

func MapTask(reply CallReply, mapf func(string, string) []KeyValue) bool {
	InitialReduceTask(reply)
	// 读入文件，，使用 mapf 生成中间数据
	fileName := reply.MapFile
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Can't open file %v", fileName)
		return false
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Can't read file [%+v]", fileName)
		return false
	}
	err = file.Close()
	if err != nil {
		log.Fatalf("Close file: [%+v] failed",  fileName)
		return false
	}
	kva := mapf(fileName, string(content))
	
	// 中间数据划分并做持久化处理
	for _, kv := range kva {
		// 划分 kva，并生成 reduce task 的 id
		reduceTaskIdx := ihash(kv.Key) % reply.ReduceNumber
		intermediateEncoder := intermediateEncoderMap[reduceTaskIdx]
		if err := intermediateEncoder.Encode(&kv); err != nil {
			log.Fatalf("Encode kv: [%+v] failed", kv)
			return false
		}
	}
	
	for reduceTaskIdx, fileName := range intermediateFileNameMap {
		err = os.Rename(fileName, fmt.Sprintf("mr-%+v-%+v", reply.MapTaskIdx, reduceTaskIdx))
		if err != nil {
			log.Fatalf("Atomic rename tmp file: [%+v] failed", fileName)
			return false
		}
	}
	intermediateEncoderMap = nil
	intermediateFileNameMap = nil
	return true
}


func InitialReduceTask(reply CallReply) {
	// initialize encoder and filename map for reduce task
	// reduceTaskIdx:json.Encoder
	// reduceTaskIdx:filename
	if intermediateEncoderMap == nil {

		intermediateEncoderMap = make(map[int]*json.Encoder)
		intermediateFileNameMap = make(map[int]string)
		for idx:=0; idx<reply.ReduceNumber; idx++ {
			fileName := fmt.Sprintf("mr-tmp-%+v-%+v", reply.MapTaskIdx, idx)
			file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Can't open the intermediate file: %+v", fileName)
			}
			intermediateFileNameMap[idx] = fileName
			intermediateEncoderMap[idx] = json.NewEncoder(file)
		}
	}
}


func ReduceTask(reply CallReply, reducef func(string, []string) string) bool {
	var intermediateKV []KeyValue
	// summary the map task intermediate output for reduce task form reply
	for idx:=0; idx<reply.MapNumber; idx++ {
		fileName := fmt.Sprintf("mr-%+v-%+v", idx, reply.ReduceTaskIdx)
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("Can't open file %v", fileName)
			return false
		}
		decoder := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := decoder.Decode(&kv); err != nil {
				log.Fatalf("Decode kv: [%+v] failed", kv)
				return false
			}
			intermediateKV = append(intermediateKV, kv)
		}
	}
	// sorted the summary kv by Key
	sort.Sort(ByKey(intermediateKV))

	// store the reduce task output
	oname := fmt.Sprintf("mr-out-%+v", reply.ReduceTaskIdx)
	tname := "mr-tmp-reduce"
	ofile, err := os.Create(tname)
	if err != nil {
		log.Fatalf("Can't create file: [%+v]", tname)
		return false
	}

	//
	// call Reduce on each distinct key in intermediate[],
	// and print the result to mr-out-0.
	//
	i := 0
	for i < len(intermediateKV) {
		j := i + 1
		for j < len(intermediateKV) && intermediateKV[j].Key == intermediateKV[i].Key {
			j++
		}
		var values []string
		for k := i; k < j; k++ {
			values = append(values, intermediateKV[k].Value)
		}
		output := reducef(intermediateKV[i].Key, values)

		// this is the correct format for each line of Reduce output.
		_, err = fmt.Fprintf(ofile, "%v %v\n", intermediateKV[i].Key, output)
		if err != nil {
			log.Fatalf("Write kv: [%+v] failed", intermediateKV[i])
			return false
		}
		i = j
	}
	err = os.Rename(tname, oname)
	if err != nil {
		log.Fatalf("Can't rename file: [%+v] to: [%+v]", tname, oname)
		return false
	}
	for idx:=0; idx<reply.MapNumber; idx++ {
		fileName := fmt.Sprintf("mr-%+v-%+v", idx, reply.ReduceTaskIdx)
		err := os.Remove(fileName)
		if err != nil {
			log.Fatalf("Can't remove file: [%+v]", fileName)
			return false
		}
	}
	err = ofile.Close()
	if err != nil {
		log.Fatalf("Close file: [%+v] failed", ofile.Name())
		return false
	}

	return true
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
