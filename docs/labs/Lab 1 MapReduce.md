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

【Blank】

提示：

【Blank】