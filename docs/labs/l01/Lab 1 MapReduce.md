可参考 mrsequential.go 和 mrapps/wc.go

任务详情：

实现分布式的 MapReduce，由 coordinator 和 worker 两个程序组成。一个 coordinator  进程，一个或多个 worker 并行的进程。这些进程运行于一个机器来仿真实际的系统。worker 与 coordinator 通过 RPC 来通信。每个 worker 向 coordinator 索取任务，从一个或多个文件读取输入，执行任务并将任务的输出写入一个或多个文件。当一个 worker 超过一定时间（10s）未完成任务，coordinator  应该通知它，并将它的任务交给其它的 worker。

代码实现：

需要完成的代码文件为：mr/coordinator.go, mr/worker.go, 和 mr/rpc.go

运行方式：

```shell
go build -race -buildmode=plugin ../mrapps/wc.go
rm mr-out*
go run -race mrcoordinator.go pg-*.txt
go run -race mrworker.go wc.so
cat mr-out-* | sort | more
```

检查运行正确性：检查是否并行、输出结果是否与串行输出一致、对于 worker 宕机的处理等。

```shell
cd src/main
bash test-mr.sh
```

注意：需要将 coordinator 注册为 [RPC Server](https://golang.org/src/net/rpc/server.go)

规则：

-   map 生成的中间文件需要根据 key 划分为 nReduce 个传入 mr/coordinator.go 作为 MakeCoordinator 的参数，来构建 coordinator。

-   reduce worker 输出命名为：mr-out-X，文件格式为："%v %v"

-   只需要更改 mr/worker.go、mr/coordinator.go、mr/rpc.go 三个文件的代码，其它源代码务必保持不动，只可以临时更改。

-   worker 将输出存放在当前文件夹

-   当 worker 结束任务时，使用 call() 通知 coordinator。

提示：

-   修改 mr/worker.go 的 Worker 通过 RPC 向 coordinator 请求任务；再修改 mr/coordinator.go 中的相应位置来响应请求；之后修改 worker.go 来读取和调用 Map 函数。

-   修改了 mr 中的代码，需要重新 build mrapps/wc.go

-   由于所有 worker 工作在同一个机器，因此不需要 GFS 这类全局的文件系统。

-   中间文件命名为 mr-X-Y，X 为 map 任务序号，Y 为 reduce 任务序号。

-   可以使用 json 来存储中间 KV

    ```golang
    // 写入
    enc := json.NewEncoder(file)
    for _, kv := ... {
    	err := enc.Encode(&kv)
    // 读取
    dec := json.NewDecoder(file)
    for {
    	var kv KeyValue
    	if err := dec.Decode(&kv); err != nil {
    		break
        }
    	kva = append(kva, kv)
      }
    ```

-   可以通过 worker 中的 ihash 函数，来为中间文件挑选 reduce worker，向其提供对应的 key。

-   可以使用 mrsequential.go 中关于读取 Map 输入文件、中间 KV 排序和 Reduce 输出的排序。

-   coordinator 是一个实时的 RPC 服务器，记得给共享数据上锁。

-   使用 race detector 模式来构建和运行

-   为了减少不必要的等待，提高吞吐量，worker 周期性的向 coordinator  请求任务，周期自行设置。

-   coordinator  无法区分 crash worker、因为某些原因停滞不前的 worker、效率过低的 worker，可以让 coordinator 等待一定时间，超时后重新分配任务给 idle worker，此实验中设置超时为 10s。

-   如果想要为 worker 实现 Backup 任务，在 test 时候不允许 worker 额外运行其他任务，因此 Backup 任务只能在 worker 结束了一段时间后运行，如 10s。

-   想要测试 crash 的情况，可以使用 mrapps/crash.go 作为应用插件代替 wc.go，它会随机的终止 Map 和 Reduce 函数。

-   为了确保在崩溃时没有人注意到部分写入的文件，MapReduce论文提到了使用临时文件并在它完全写入后自动重命名它的技巧。可以用 ioutil.TempFile 和 os.Rename 来进行原子重命名。

-   test-mr-many 是 test-mr 的带有超时设置的版本，test-mr 的输出目录为 mr-tmp/。

