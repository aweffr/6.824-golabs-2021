## Abstract

MapReduce is a programming model and an associated implementation for processing and generating large data sets. Users specify a map function that processes a key/value pair to generate a set of intermediate key/value pairs, and a reduce function that merges all intermediate values associated with the same intermediate key. 

## Introduction

The run-time system takes care of the details of partitioning the input data, scheduling the program’s execution across a set of machines, handling machine failures, and managing the required inter-machine communication. 

process large amounts of raw data, such as crawled documents, web request logs, etc., to compute various kinds of derived data, such as inverted indices, various representations of the graph structure of web documents, summaries of the number of pages crawled per host, the set of most frequent queries in a
given day, etc. 

However, the input data is usually large and the computations have to be distributed across hundreds or thousands of machines in order to finish in a reasonable amount of time. 

The issues of how to parallelize the computation, distribute the data, and handle failures conspire to obscure the original simple computation with large amounts of complex code to deal with these issues.

express the simple computations we were trying to perform but hides the messy details of parallelization, fault-tolerance, data distribution and load balancing in a library. Our abstraction is inspired by the map and reduce primitives present in Lisp
and many other functional languages.



We realized that most of our computations involved applying a map operation to each logical “record” in our input in order to compute a set of intermediate key/value pairs, and then applying a reduce operation to all the values that shared the same key, in order to combine the derived data appropriately.



Our use of a functional model with user-specified map and reduce operations allows us to parallelize large computations easily and to use re-execution as the primary mechanism for fault tolerance.

## Programming Model

The computation takes a set of input key/value pairs, and produces a set of output key/value pairs. The user of the MapReduce library expresses the computation as two functions: Map and Reduce.
Map, written by the user, takes an input pair and produces a set of intermediate key/value pairs. The MapReduce library groups together all intermediate values associated with the same intermediate key I and passes them to the Reduce function.
The Reduce function, also written by the user, accepts an intermediate key I and a set of values for that key. It merges together these values to form a possibly smaller set of values. Typically just zero or one output value is produced per Reduce invocation. The intermediate values are supplied to the user’s reduce function via an iterator. This allows us to handle lists of values that are too large to fit in memory.

#### Example

```c++
#include "mapreduce/mapreduce.h"
// User’s map function
class WordCounter : public Mapper {
	public:
		virtual void Map(const MapInput& input) {
			const string& text = input.value();
			const int n = text.size();
			for (int i = 0; i < n; ) {
				// Skip past leading whitespace
				while ((i < n) && isspace(text[i]))
					i++;
				// Find word end
				int start = i;
				while ((i < n) && !isspace(text[i]))
					i++;
				if (start < i)
					Emit(text.substr(start,i-start),"1");
			}
		}
};
REGISTER_MAPPER(WordCounter);
// User’s reduce function
class Adder : public Reducer {
	virtual void Reduce(ReduceInput* input) {
		// Iterate over all entries with the
		// same key and add the values
		int64 value = 0;
		while (!input->done()) {
			value += StringToInt(input->value());
			input->NextValue();
		}
		// Emit sum for input->key()
		Emit(IntToString(value));
	}
};
REGISTER_REDUCER(Adder);
int main(int argc, char** argv) {
	ParseCommandLineFlags(argc, argv);
	MapReduceSpecification spec;
	// Store list of input files into "spec"
	for (int i = 1; i < argc; i++) {
		MapReduceInput* input = spec.add_input();
		input->set_format("text");
		input->set_filepattern(argv[i]);
		input->set_mapper_class("WordCounter");
	}
	// Specify the output files:
	//
	/gfs/test/freq-00000-of-00100
	//
	/gfs/test/freq-00001-of-00100
	//
	...
	MapReduceOutput* out = spec.output();
	out->set_filebase("/gfs/test/freq");
	out->set_num_tasks(100);
	out->set_format("text");
	out->set_reducer_class("Adder");
	// Optional: do partial sums within map
	// tasks to save network bandwidth
	out->set_combiner_class("Adder");
	// Tuning parameters: use at most 2000
	// machines and 100 MB of memory per task
	spec.set_machines(2000);
	spec.set_map_megabytes(100);
	spec.set_reduce_megabytes(100);
	// Now run it
	MapReduceResult result;
	if (!MapReduce(spec, &result)) 				abort();
	// Done: ’result’ structure contains info
	// about counters, time taken, number of
	// machines used, etc.
	return 0;
}
```

#### Types

conceptually the map and reduce functions supplied by the user have associated types:

map(k1,v1) → list(k2,v2)

reduce(k2,list(v2)) → list(v2)

I.e., the input keys and values are drawn from a different domain than the output keys and values. Furthermore, the intermediate keys and values are from the same domain as the output keys and values.

Our C++ implementation passes strings to and from the user-defined functions and leaves it to the user code to convert between strings and appropriate types.

#### More Examples

Distributed Grep

Count of URL Access Frequency

Reverse Web-Link Graph

Term-Vector per Host

Inverted Index

Distributed Sort

## Implementation

![Execution Overview](res/MapReduce/Execution-Overview.png)

The **Map** invocations are distributed across multiple machines (worker). The input data are automatically partitioning to *M* **splits**, each split as a machines input. 

**Reduce** invocations are distributed by partitioning the **intermediate key** space into *R* **pieces** using a **partitioning function** (e.g., hash(key) mod R).

#### Execution Overview

一、input files 首先被一个用户程序分割为 M 个划分，通常每个划分16 MB 到64 MB，然后在集群中的一组机器上启动此程序的副本。

二、有一个机器上的程序副本是特殊的，称为 Master，其它机器称为 Worker，完成由 Master 所分配的工作。这里有 M 个 map 任务和 R 个 reduce 任务待分配。一个 Worker 对应一个 map 任务+一个 input file 划分，或是一个 reduce 任务+一个 R output file 划分。

三、Worker 读取其对应的 input data split 并从中解析出 key/value pairs，传入 map 函数（用户定义），map 函数生成 intermediate key/value pairs 存入内存缓存。

四、buffer pairs 会定期写入本地磁盘，并通过 partitioning 函数划分成 R 个分区。之后这些 buffer pairs 的磁盘地址会传递回 Master 并由其分配给 reduce workers。

五、reduce worker 收到 Master 传来的 buffer pairs 磁盘地址，使用远程过程调用从 map worker 的本地磁盘读取 buffer pairs，然后对缓存中的 kv pairs 根据 key 来进行排序，来将相同 key 的 kv pairs 组织在一起。当磁盘中的 buffer pairs 太大时会进行外排序。

六、reduce worker 遍历已排序的中间数据，并对遇到的每个唯一的 intermediate key 进行遍历，将 key 和 intermediate value set 传递给 reduce 函数。函数的输出追加到这个 reduce 分区的最终输出文件中。

七、所以 map 和 reduce 任务结束时，Master 唤醒用户程序，MapReduce 调用返回用户代码。

#### Master Data Structure

每个任务具有三种状态：闲置 idle、进行中 in-progress、完成 completed，且保存着 worker （non-idle）的标识符

master 是一个通道，在 map worker 和 reduce worker 之间传递 intermediate file 区域的地址。因此，master 存储了每一个已完成 map 任务所生成的 R 个中间文件区域的地址（M 个 map 任务汇总的中间文件存到磁盘缓冲区，再划分为 R 个区域传递给 reduce worker）。

#### Fault Tolerance

**Worker Faliure**

Master 周期性的 ping 每个 worker，一定时间收不到回应则将此 worker 判定为 failed。完成 map 任务的 worker 会进入初始的 idle 状态，之后可以调度新的 worker，同理，当进行任务的 worker failed 时也会进入 idle，等待被重新调度。

已完成的 map 任务不再可运行，因为它的输出文件已存入磁盘，被占用。

已完成的 reduce 任务不需再运行，因为它的输出文件已存入磁盘，被占用。

当 map worker B 代替失败的 A 执行 map 任务时，所有 reduce worker 会接收到 re-execution 通知，任何尚未从 worker A 读取数据的 reduce 任务都将从 worker B 读取数据。

MapReduce 使用简单的 re-execution 来处理大范围的 worker failure。

**Master Failure**

根据 Master 数据结构来建立检查点，待到恢复时或是切换 Master 时，将已保存的检查点重新加载。

**Semantics in the Presence of Failures**

依赖于 map 和 reduce 操作的原子性提交，只要用户指定的 map 和 reduce 函数确定，那么无论是分布式还是普通顺序执行的结果是一致的。

注：这里的语义指的是数据处理的逻辑，即输入与输出的对应关系。

每个正在进行的任务都将其输出写入私有临时文件。reduce 任务生成一个这样的文件，map 任务生成R个这样的文件(对应于每个 reduce 任务一个)。map 任务结束时会向 Master 发送消息，携带着 R 个输出文件的名称。如果 master 接收到一个已经完成的map任务的完成消息，它将忽略该消息。否则，会将 R 个文件名存入 Master 的数据结构中。

当 reduce 任务完成时，reduce worker 以原子方式将其临时输出文件重命名为最终输出文件。如果在多台机器上执行相同的reduce任务，那么将对同一个最终输出文件执行多个重命名调用。

依赖于底层文件管理系统提供的原子重命名操作，确保了最终的文件系统状态只包含由 reduce worker 执行一次所生产的数据。

大量的 map 和 reduce 操作是确定的，且执行语义与普通的顺序执行一致。这使得程序员很容易的去研究所编写程序的行为，只需考虑确定情况下的正常语义即可。

而当操作是非确定的时，依然可以提供较弱但是仍旧合理的语义。In the presence of non-deterministic operators, the output of a particular reduce task R1 is equivalent to the output for R1 produced by a sequential execution of the non-deterministic program. However, the output for a different reduce task R2 may correspond to the output for R2 produced by a different sequential execution of the non-deterministic program.

Consider map task M and reduce tasks R1 and R2. Let e(Ri) be the execution of Ri that committed (there is exactly one such execution). The weaker semantics arise because e(R1) may have read the output produced by one execution of M and e(R2) may have read the output produced by a different execution of M.

注：这里的非确定性，应该指的是 map 的输出与 reduce 的输入可能不是一一对应的。

。。。

由于单词频率倾向于遵循 Zipf 分布，每个 map 任务将产生数百或数千个 <the, 1> 形式的记录。所有这些计数将通过网络发送到一个 Reduce 任务，然后由 Reduce 函数相加产生一个数字。因此我们允许用户指定一个可选的 Combiner 函数，该函数在通过网络发送数据之前对数据进行部分合并。

reduce 函数和 combiner 函数之间的唯一区别是 MapReduce 库如何处理函数的输出。reduce 函数的输出被写入最终的输出文件。combiner 函数的输出被写入一个中间文件，该文件将被发送到 reduce 任务。

