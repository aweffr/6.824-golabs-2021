# 6.824-golabs-2021
[MIT Course 6.824 (Spring 2021)](https://pdos.csail.mit.edu/6.824/) Lab Source Implementation.

## Source Codes

-   [x] [Lab 1: MapReduce](src/mr/)
-   [ ] [Lab 2: replication for fault-tolerance using Raft](src/raft/)
-   [ ] [Lab 3: fault-tolerant key/value store](src/kvraft/)
-   [ ] [Lab 4: sharded key/value store](src/shardkv/)

## Documents

-   [The Preparation Readings](./docs/readings/)
-   [The Course Notes](./docs/notes/)
-   [The Lab Implementation Records](./docs/labs/)

## Schedule: Spring 2021

<table class="calendar" cellspacing="0" cellpadding="6" width="100%">
 <thead>
  <tr>
   <td width="16%">Monday</td><td width="26%">Tuesday</td>
   <td width="16%">Wednesday</td><td width="26%">Thursday</td>
   <td width="16%">Friday</td>
  </tr>
 </thead>
<tbody><tr> <!-- week of feb 15 -->
  <td id="2021-2-15"><span class="date">feb 15</span></td>
  <td id="2021-2-16" class="lecture"><span class="date">feb 16</span><br>
    <b>LEC 1:</b> <a href="notes/l01.txt">Introduction</a>, <a href="https://youtu.be/WtZ7pcRSkOA">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/mapreduce.pdf">MapReduce (2004)</a></span><br>
    <span class="reading"><b>Assigned:</b> <a href="labs/lab-mr.html">Lab 1: MapReduce</a></span><br>
    <i>First day of classes</i></td>
  <td id="2021-2-17"><span class="date">feb 17</span></td>
  <td id="2021-2-18" class="lecture"><span class="date">feb 18</span><br>
    <b>LEC 2:</b> <a href="notes/l-rpc.txt">RPC and Threads</a>, <a href="notes/crawler.go">crawler.go</a>, <a href="notes/kv.go">kv.go</a>, <a href="notes/condvar.tar.gz">vote examples</a>, <a href="https://youtu.be/oZR76REwSyA">video</a><br>
    <span class="reading"><b>Preparation:</b>  Do <a href="http://tour.golang.org/">Online Go tutorial</a>  (<a href="papers/tour-faq.txt">FAQ</a>) (<a href="questions.html?q=q-gointro&amp;lec=2">Question</a>)</span></td>
  <td id="2021-2-19"><span class="date">feb 19</span></td>
</tr>
<tr> <!-- week of feb 22 -->
  <td id="2021-2-22"><span class="date">feb 22</span></td>
  <td id="2021-2-23" class="lecture"><span class="date">feb 23</span><br>
    <b>LEC 3:</b> <a href="notes/l-gfs.txt">GFS</a>, <a href="https://youtu.be/6ETFk1-53qU">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/gfs.pdf">GFS (2003)</a> (<a href="papers/gfs-faq.txt">FAQ</a>) (<a href="questions.html?q=q-gfs&amp;lec=3">Question</a>)</span><br>
    <span class="reading"><b>Assigned:</b> <a href="labs/lab-raft.html">Lab 2: Raft</a></span></td>
  <td id="2021-2-24"><span class="date">feb 24</span></td>
  <td id="2021-2-25" class="lecture"><span class="date">feb 25</span><br>
    <b>LEC 4:</b> <a href="notes/l-vm-ft.txt">Primary-Backup Replication</a>, <a href="https://youtu.be/gXiDmq1zDq4">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/vm-ft.pdf">Fault-Tolerant Virtual Machines (2010)</a> (<a href="papers/vm-ft-faq.txt">FAQ</a>) (<a href="questions.html?q=q-vm-ft&amp;lec=4">Question</a>)</span></td>
  <td id="2021-2-26" class="due"><span class="date">feb 26</span><br>
    <span class="deadline"><b>DUE:</b> <a href="labs/lab-mr.html">Lab 1</a></span></td>
</tr>
<tr> <!-- week of mar 1 -->
  <td id="2021-3-1"><span class="date">mar 1</span></td>
  <td id="2021-3-2" class="lecture"><span class="date">mar 2</span><br>
    <b>LEC 5:</b> <a href="notes/l-raft.txt">Fault Tolerance: Raft (1)</a>, <a href="https://youtu.be/R2-9bsKmEbo">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/raft-extended.pdf">Raft (extended) (2014), to end of Section 5</a>  (<a href="papers/raft-faq.txt">FAQ</a>) (<a href="questions.html?q=q-raft&amp;lec=5">Question</a>)</span></td>
  <td id="2021-3-3"><span class="date">mar 3</span></td>
  <td id="2021-3-4" class="lecture"><span class="date">mar 4</span><br>
    <b>LEC 6:</b> <a href="notes/mr_solution_lecture.pdf">Q&amp;A Lab 1</a>, <a href="https://youtu.be/QkPiiRQmom8">video</a><br>
    <span class="reading"><b>Preparation:</b>  (<a href="questions.html?q=q-QAlab&amp;lec=6">Question</a>)</span></td>
  <td id="2021-3-5" class="due"><span class="date">mar 5</span><br>
    <span class="deadline"><b>DUE:</b> <a href="labs/lab-raft.html">Lab 2A</a></span></td>
</tr>
<tr> <!-- week of mar 8 -->
  <td id="2021-3-8" class="holiday"><span class="date">mar 8</span><br>
    No Class</td>
  <td id="2021-3-9" class="special" style="border: 2px solid blue;"><span class="date">mar 9</span><br>
    <i>Monday schedule</i></td>
  <td id="2021-3-10"><span class="date">mar 10</span></td>
  <td id="2021-3-11" class="lecture"><span class="date">mar 11</span><br>
    <b>LEC 7:</b> <a href="notes/l-raft2.txt">Fault Tolerance: Raft (2)</a>, <a href="video/7.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/raft-extended.pdf">Raft (extended) (2014), Section 7 to end (but not Section 6)</a>  (<a href="papers/raft2-faq.txt">FAQ</a>) (<a href="questions.html?q=q-raft2&amp;lec=7">Question</a>)</span></td>
  <td id="2021-3-12" class="due"><span class="date">mar 12</span><br>
    <span class="deadline"><b>DUE:</b> <a href="labs/lab-raft.html">Lab 2B</a></span></td>
</tr>
<tr> <!-- week of mar 15 -->
  <td id="2021-3-15"><span class="date">mar 15</span></td>
  <td id="2021-3-16" class="lecture"><span class="date">mar 16</span><br>
    <b>LEC 8:</b> Q&amp;A Lab2 A+B<br>
    <span class="reading"><b>Preparation:</b>  (<a href="questions.html?q=q-QAlab&amp;lec=8">Question</a>)</span><br>
    <span class="reading"><b>Assigned:</b> <a href="project.html">Final Project</a></span></td>
  <td id="2021-3-17"><span class="date">mar 17</span></td>
  <td id="2021-3-18" class="lecture"><span class="date">mar 18</span><br>
    <b>LEC 9:</b> <a href="notes/l-zookeeper.txt">Zookeeper</a>, <a href="video/8.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/zookeeper.pdf">ZooKeeper (2010)</a>  (<a href="papers/zookeeper-faq.txt">FAQ</a>) (<a href="questions.html?q=q-zookeeper&amp;lec=9">Question</a>)</span></td>
  <td id="2021-3-19" class="due"><span class="date">mar 19</span><br>
    <span class="deadline"><b>DUE:</b> <a href="labs/lab-raft.html">Lab 2C</a></span><br>
    <i>ADD DATE</i></td>
</tr>
<tr> <!-- week of mar 22 -->
  <td id="2021-3-22" class="holiday"><span class="date">mar 22</span><br>
    No Class</td>
  <td id="2021-3-23" class="assign"><span class="date">mar 23</span><br>
    <span class="reading"><b>Assigned:</b> <a href="labs/lab-kvraft.html">Lab 3: KV Raft</a></span><br>
    No Class</td>
  <td id="2021-3-24"><span class="date">mar 24</span></td>
  <td id="2021-3-25" class="lecture"><span class="date">mar 25</span><br>
    <b>LEC 10:</b> <a href="notes/gopattern.pdf">Guest lecturer on Go</a> (<a href="http://swtch.com/~rsc/">Russ Cox</a> Google/Go)<br>
    <span class="reading"><b>Preparation:</b>  (<a href="papers/go-faq.txt">FAQ</a>) (<a href="questions.html?q=q-go&amp;lec=10">Question</a>)</span></td>
  <td id="2021-3-26" class="due"><span class="date">mar 26</span><br>
    <span class="deadline"><b>DUE:</b> <a href="labs/lab-raft.html">Lab 2D</a></span></td>
</tr>
<tr> <!-- week of mar 29 -->
  <td id="2021-3-29"><span class="date">mar 29</span></td>
  <td id="2021-3-30" class="lecture"><span class="date">mar 30</span><br>
    <b>LEC 11:</b> <a href="notes/l-cr.txt">Chain Replication</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/cr-osdi04.pdf">CR (2004)</a> (<a href="questions.html?q=q-cr&amp;lec=11">Question</a>)</span></td>
  <td id="2021-3-31"><span class="date">mar 31</span></td>
  <td id="2021-4-1" class="quiz"><span class="date">apr 1</span><br>
    <b>Remote Mid-term Exam</b> <br>
    <b>Materials:</b> Open book, notes, laptop<br>
    <b>Scope:</b> Lectures 1 through 10, Labs 1 and 2<br>
    <a href="https://pdos.csail.mit.edu/6.824/quizzes.html">Old Exams</a></td>
  <td id="2021-4-2" class="due"><span class="date">apr 2</span><br>
    <span class="deadline"><b>DUE:</b> <a href="project.html">Project proposals</a> (if you are doing a project)</span></td>
</tr>
<tr> <!-- week of apr 5 -->
  <td id="2021-4-5"><span class="date">apr 5</span></td>
  <td id="2021-4-6" class="lecture"><span class="date">apr 6</span><br>
    <b>LEC 12:</b> <a href="notes/l-frangipani.txt">Cache Consistency: Frangipani</a>, <a href="video/11.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/thekkath-frangipani.pdf">Frangipani</a>  (<a href="papers/frangipani-faq.txt">FAQ</a>) (<a href="questions.html?q=q-frangipani&amp;lec=12">Question</a>)</span><br>
    <span class="reading"><b>Assigned:</b> <a href="labs/lab-shard.html">Lab 4: Sharded KV</a></span></td>
  <td id="2021-4-7"><span class="date">apr 7</span></td>
  <td id="2021-4-8" class="lecture"><span class="date">apr 8</span><br>
    <b>LEC 13:</b> <a href="notes/l-2pc.txt">Distributed Transactions</a>, <a href="video/12.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="https://ocw.mit.edu/resources/res-6-004-principles-of-computer-system-design-an-introduction-spring-2009/online-textbook/">6.033 Chapter 9</a>, just 9.1.5, 9.1.6, 9.5.2, 9.5.3, 9.6.3 (<a href="papers/chapter9-faq.txt">FAQ</a>) (<a href="questions.html?q=q-chapter9&amp;lec=13">Question</a>)</span></td>
  <td id="2021-4-9" class="due"><span class="date">apr 9</span><br>
    <span class="deadline"><b>DUE:</b> <a href="labs/lab-kvraft.html">Lab 3A</a></span></td>
</tr>
<tr> <!-- week of apr 12 -->
  <td id="2021-4-12"><span class="date">apr 12</span></td>
  <td id="2021-4-13" class="special"><span class="date">apr 13</span><br>
    <i>Hacking day, no lecture</i></td>
  <td id="2021-4-14"><span class="date">apr 14</span></td>
  <td id="2021-4-15" class="lecture"><span class="date">apr 15</span><br>
    <b>LEC 14:</b> <a href="notes/l-spanner.txt">Spanner</a>, <a href="video/13.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/spanner.pdf">Spanner (2012)</a> (<a href="papers/spanner-faq.txt">FAQ</a>) (<a href="questions.html?q=q-spanner&amp;lec=14">Question</a>)</span></td>
  <td id="2021-4-16" class="due"><span class="date">apr 16</span><br>
    <span class="deadline"><b>DUE:</b> <a href="labs/lab-kvraft.html">Lab 3B</a></span></td>
</tr>
<tr> <!-- week of apr 19 -->
  <td id="2021-4-19" class="holiday"><span class="date">apr 19</span><br>
    Patriots day</td>
  <td id="2021-4-20" class="holiday"><span class="date">apr 20</span><br>
    No Class</td>
  <td id="2021-4-21"><span class="date">apr 21</span></td>
  <td id="2021-4-22" class="lecture"><span class="date">apr 22</span><br>
    <b>LEC 15:</b> <a href="notes/l-farm.txt">Optimistic Concurrency Control</a>, <a href="video/14.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/farm-2015.pdf">FaRM (2015)</a>   (<a href="papers/farm-faq.txt">FAQ</a>) (<a href="questions.html?q=q-farm&amp;lec=15">Question</a>)</span></td>
  <td id="2021-4-23" class="due"><span class="date">apr 23</span><br>
    <span class="deadline"><b>DUE:</b> <a href="labs/lab-shard.html">Lab 4A</a></span></td>
</tr>
<tr> <!-- week of apr 26 -->
  <td id="2021-4-26"><span class="date">apr 26</span></td>
  <td id="2021-4-27" class="lecture"><span class="date">apr 27</span><br>
    <b>LEC 16:</b> <a href="notes/l-spark.txt">Big Data: Spark</a>, <a href="video/15.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/zaharia-spark.pdf">Spark (2012)</a> (<a href="papers/spark-faq.txt">FAQ</a>) (<a href="questions.html?q=q-spark&amp;lec=16">Question</a>)</span></td>
  <td id="2021-4-28"><span class="date">apr 28</span></td>
  <td id="2021-4-29" class="lecture"><span class="date">apr 29</span><br>
    <b>LEC 17:</b> <a href="notes/l-memcached.txt">Cache Consistency: Memcached at Facebook</a>, <a href="video/16.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/memcache-fb.pdf">Memcached at Facebook (2013)</a>  (<a href="papers/memcache-faq.txt">FAQ</a>) (<a href="questions.html?q=q-memcached&amp;lec=17">Question</a>)</span><br>
    <i><b class="deadline">DROP DATE</b></i></td>
  <td id="2021-4-30"><span class="date">apr 30</span></td>
</tr>
<tr> <!-- week of may 3 -->
  <td id="2021-5-3"><span class="date">may 3</span></td>
  <td id="2021-5-4" class="lecture"><span class="date">may 4</span><br>
    <b>LEC 18:</b> <a href="notes/l-cops.txt">Causal Consistency, COPS</a>, <a href="video/17.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/cops.pdf">COPS (2011)</a> (<a href="questions.html?q=q-cops&amp;lec=18">Question</a>)</span></td>
  <td id="2021-5-5"><span class="date">may 5</span></td>
  <td id="2021-5-6" class="lecture"><span class="date">may 6</span><br>
    <b>LEC 19:</b> <a href="notes/l-ct.txt">Fork Consistency, Certificate Transparency</a>, <a href="video/18.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="https://www.certificate-transparency.org/what-is-ct">Certificate Transparency</a>, <a href="https://www.certificate-transparency.org/how-ct-works">Also This</a>, <a href="https://research.swtch.com/tlog">And This</a>, but skip the Tiles sections and the appendices. (<a href="papers/ct-faq.txt">FAQ</a>) (<a href="questions.html?q=q-ct&amp;lec=19">Question</a>)</span></td>
  <td id="2021-5-7" class="holiday"><span class="date">may 7</span><br>
    Student holiday</td>
</tr>
<tr> <!-- week of may 10 -->
  <td id="2021-5-10"><span class="date">may 10</span></td>
  <td id="2021-5-11" class="lecture"><span class="date">may 11</span><br>
    <b>LEC 20:</b> Peer-to-peer: <a href="notes/l-bitcoin.txt">Bitcoin</a>, <a href="video/19.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/bitcoin.pdf">Bitcoin (2008)</a>, and <a href="http://www.michaelnielsen.org/ddi/how-the-bitcoin-protocol-actually-works">summary</a> (<a href="papers/bitcoin-faq.txt">FAQ</a>) (<a href="questions.html?q=q-bitcoin&amp;lec=20">Question</a>)</span></td>
  <td id="2021-5-12"><span class="date">may 12</span></td>
  <td id="2021-5-13" class="special"><span class="date">may 13</span><br>
    <i>Hacking day, no lecture</i></td>
  <td id="2021-5-14" class="due"><span class="date">may 14</span><br>
    <span class="deadline"><b>DUE:</b> <a href="labs/lab-shard.html">Lab 4B</a></span><br>
    <span class="deadline"><b>DUE:</b> <a href="project.html">Project reports and code</a></span></td>
</tr>
<tr> <!-- week of may 17 -->
  <td id="2021-5-17"><span class="date">may 17</span></td>
  <td id="2021-5-18" class="lecture"><span class="date">may 18</span><br>
    <b>LEC 21:</b> <a href="notes/l-blockstack.txt">Blockstack</a>, <a href="video/20.html">video</a><br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/blockstack-atc16.pdf">BlockStack (2016)</a> (<a href="papers/blockstack-faq.txt">FAQ</a>) (<a href="questions.html?q=q-blockstack&amp;lec=21">Question</a>)</span></td>
  <td id="2021-5-19"><span class="date">may 19</span></td>
  <td id="2021-5-20" class="lecture"><span class="date">may 20</span><br>
    <b>LEC 22:</b> Project demos<br>
    <span class="reading"><b>Preparation:</b>  Read <a href="papers/katabi-analogicfs.pdf">AnalogicFS experience paper</a> (<a href="papers/analogicfs-faq.txt">FAQ</a>) (<a href="questions.html?q=q-analogic&amp;lec=22">Question</a>)</span><br>
    <i>Last day of classes</i></td>
  <td id="2021-5-21"><span class="date">may 21</span></td>
</tr>
<tr> <!-- week of may 24 -->
  <td id="2021-5-24" class="quiz"><span class="date">may 24</span><br>
    Finals</td>
  <td id="2021-5-25" class="quiz"><span class="date">may 25</span><br>
    Finals</td>
  <td id="2021-5-26" class="quiz"><span class="date">may 26</span><br>
    Finals</td>
  <td id="2021-5-27" class="quiz"><span class="date">may 27</span><br>
    Finals</td>
  <td id="2021-5-28"><span class="date">may 28</span></td>
</tr>
</tbody></table>






