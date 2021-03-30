Course components:

  **lectures**
  **papers**
  two exams
  **labs**
  final project (optional)

Labs:

  Lab 1: MapReduce
  Lab 2: replication for fault-tolerance using Raft
  Lab 3: fault-tolerant key/value store
  Lab 4: sharded key/value store

MAIN TOPICS

This is a course about infrastructure for applications.
  * Storage.
  * Communication.
  * Computation.

The big goal: abstractions that hide the complexity of distribution.

Topic: fault tolerance

  1000s of servers, big network -> always something broken
    We'd like to hide these failures from the application.
  We often want:
    Availability -- app can make progress despite failures
    Recoverability -- app will come back to life when failures are repaired
  Big idea: replicated servers.
    If one server crashes, can proceed using the other(s).
    Very hard to get right
      server may not have crashed, but just unreachable for some
        but still serving requests from clients
    Labs 1, 2 and 3

Topic: consistency
  General-purpose infrastructure needs well-defined behavior.
    E.g. "Get(k) yields the value from the most recent Put(k,v)."
  Achieving good behavior is hard!
    "Replica" servers are hard to keep identical.

Topic: performance
  The goal: scalable throughput
    Nx servers -> Nx total throughput via parallel CPU, disk, net.
  Scaling gets harder as N grows:
    Load im-balance, stragglers, slowest-of-N latency.
    Non-parallelizable code: initialization, interaction.
    Bottlenecks from shared resources, e.g. network.
  Some performance problems aren't easily solved by scaling
    e.g. quick response time for a single user request
    e.g. all users want to update the same data
    often requires better design rather than just more computers
  Lab 4

Topic: Fault-tolerance, consistency, and performance are enemies.
  Strong fault tolerance requires communication
    e.g., send data to backup
  Strong consistency requires communication,
    e.g. Get() must check for a recent Put().
  Many designs provide only weak consistency, to gain speed.
    e.g. Get() does *not* yield the latest Put()!
    Painful for application programmers but may be a good trade-off.
  Many design points are possible in the consistency/performance spectrum!

Topic: implementation
  RPC, threads, concurrency control.
  The labs...



now: MR