本文只从论文中的第三章《实现》选取重要的内容来进行总结。

为了对 MapReduce 有一个大致的认识，这是第二章的节选，可自行阅读理解：

>The **Map** invocations are distributed across multiple machines (worker). The input data are automatically partitioning to *M* **splits**, each split as a machines input. 
>
>**Reduce** invocations are distributed by partitioning the **intermediate key** space into *R* **pieces** using a **partitioning function** (e.g., hash(key) mod R).
>
>......
>
>conceptually the map and reduce functions supplied by the user have associated types:
>
>map(k1,v1) → list(k2,v2)
>
>reduce(k2,list(v2)) → list(v2)

## Implementation

### Execution Overview

![Execution Overview](res/2021-03-01-MapReduce/Execution-Overview.png)

一、input files 首先被一个用户程序分割为 M 个划分，通常每个划分16 MB 到64 MB，然后在集群中的一组机器上启动此程序的副本。

二、有一个机器上的程序副本是特殊的，称为 Master，其它机器称为 Worker，完成由 Master 所分配的工作。这里有 M 个 map 任务和 R 个 reduce 任务待分配。一个 Worker 对应一个 map 任务+一个 input file 划分，或是一个 reduce 任务+一个 R output file 划分。

三、Worker 读取其对应的 input data split 并从中解析出 key/value pairs，传入 map 函数（用户定义），map 函数生成 intermediate key/value pairs 存入内存缓存。

四、buffer pairs 会定期写入本地磁盘，并通过 partitioning 函数划分成 R 个分区。之后这些 buffer pairs 的磁盘地址会传递回 Master 并由其分配给 reduce workers。

五、reduce worker 收到 Master 传来的 buffer pairs 磁盘地址，使用远程过程调用从 map worker 的本地磁盘读取 buffer pairs，然后对缓存中的 kv pairs 根据 key 来进行排序，来将相同 key 的 kv pairs 组织在一起。当磁盘中的 buffer pairs 太大时会进行外排序。

六、reduce worker 遍历已排序的中间数据，并对遇到的每个唯一的 intermediate key 进行遍历，将 key 和 intermediate value set 传递给 reduce 函数。函数的输出追加到这个 reduce 分区的最终输出文件中。

七、所以 map 和 reduce 任务结束时，Master 唤醒用户程序，MapReduce 调用返回用户代码。

### Master Data Structure

每个任务具有三种状态：闲置 idle、进行中 in-progress、完成 completed，非闲置的任务还保存着 worker 的标识符 id

master 是一个通道，在 map worker 和 reduce worker 之间传递 intermediate file 区域的地址。因此，master 存储了每一个已完成 map 任务所生成的 R 个中间文件区域的地址（M 个 map 任务汇总的中间文件存到磁盘缓冲区，再划分为 R 个区域传递给 reduce worker）。

>   注意：在 6.824 Lab1 的实现中简化了上面关于中间文件地址的获取过程—— reduce worker 可直接根据 map number 和 reduce task id 来获取中间输出的文件地址。

### Fault Tolerance

**Worker Faliure**

Master 周期性的探测每个 worker，一定时间收不到回应则将此 worker 判定为 failed。完成 map 任务的 worker 会进入初始的 idle 状态，之后可以调度新的 worker，同理，当进行任务的 worker failed 时也会进入 idle，等待被重新调度。

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

而当操作是非确定的时，依然可以提供较弱但是仍旧合理的语义。下面是原文的解释，感兴趣可以读一读：

>   In the presence of non-deterministic operators, the output of a particular reduce task R1 is equivalent to the output for R1 produced by a sequential execution of the non-deterministic program. However, the output for a different reduce task R2 may correspond to the output for R2 produced by a different sequential execution of the non-deterministic program.
>
>   Consider map task M and reduce tasks R1 and R2. Let e(Ri) be the execution of Ri that committed (there is exactly one such execution). The weaker semantics arise because e(R1) may have read the output produced by one execution of M and e(R2) may have read the output produced by a different execution of M.

### Locality



### Task Granularity



### Backup Tasks



