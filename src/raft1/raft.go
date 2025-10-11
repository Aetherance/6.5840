package raft

// 文件 raftapi/raft.go 定义了 Raft 必须暴露给服务器（或测试程序）的接口，
// 但有关每个函数的更多细节请参阅下面的注释。
//
// Make() 创建一个新的 Raft 对等节点，该节点实现了 Raft 接口。

import (
	//	"bytes"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	//	"6.5840/labgob"
	"6.5840/labrpc"
	"6.5840/raftapi"
	"6.5840/tester1"
)

// 实现单个 Raft 对等节点的 Go 对象。
type Raft struct {
	mu        sync.Mutex          // 保护对此对等节点状态共享访问的锁
	peers     []*labrpc.ClientEnd // 所有对等节点的 RPC 端点
	persister *tester.Persister   // 用于保存此对等节点持久化状态的对象
	me        int                 // 此对等节点在 peers[] 中的索引
	dead      int32               // 由 Kill() 设置
	isLeader  bool

	// All
	currentTerm int
	voteFor 	int
	logs 		[]LogEntry

	// All changable
	commitIndex int
	lastApplied int

	// Leader
	nextIndex  []int
	matchIndex []int

	// 你的数据在这里 (3A, 3B, 3C)。
	// 查看论文中的图 2 了解 Raft 服务器必须维护的状态描述。
}

type LogEntry struct {
	Command interface{}
	Term int
}

// 返回 currentTerm 以及此服务器是否认为自己是领导者。
func (rf *Raft) GetState() (int, bool) {
	var term int
	var isleader bool
	term = rf.currentTerm
	isleader = rf.isLeader
	return term, isleader
}

// 将 Raft 的持久状态保存到稳定存储中，
// 以便在崩溃和重启后可以恢复。
// 有关应持久化哪些内容的描述，请参阅论文中的图 2。
// 在实现快照之前，你应该将 nil 作为第二个参数传递给 persister.Save()。
// 实现快照后，传递当前快照（如果还没有快照则为 nil）。
func (rf *Raft) persist() {
	// 你的代码在这里 (3C)。
	// 示例：
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// raftstate := w.Bytes()
	// rf.persister.Save(raftstate, nil)
}

// 恢复之前持久化的状态。
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // 没有任何状态的引导？
		return
	}
	// 你的代码在这里 (3C)。
	// 示例：
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

// Raft 持久化日志中有多少字节？
func (rf *Raft) PersistBytes() int {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.persister.RaftStateSize()
}

// 服务层表示已创建包含直到并包括 index 的所有信息的快照。
// 这意味着服务层不再需要直到（并包括）该索引的日志。
// Raft 现在应尽可能修剪其日志。
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// 你的代码在这里 (3D)。
}

// 示例 RequestVote RPC 参数结构。
// 字段名必须以大写字母开头！
type RequestVoteArgs struct {
	// 你的数据在这里 (3A, 3B)。
}

// 示例 RequestVote RPC 回复结构。
// 字段名必须以大写字母开头！
type RequestVoteReply struct {
	// 你的数据在这里 (3A)。
}

// 示例 RequestVote RPC 处理程序。
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// 你的代码在这里 (3A, 3B)。
}

// 发送 RequestVote RPC 到服务器的示例代码。
// server 是目标服务器在 rf.peers[] 中的索引。
// 期望 RPC 参数在 args 中。
// 使用 RPC 回复填充 *reply，因此调用方应传递 &reply。
// 传递给 Call() 的 args 和 reply 类型必须与处理函数中声明的参数类型相同（包括它们是否是指针）。
//
// labrpc 包模拟了一个有损网络，其中服务器可能无法访问，请求和回复可能丢失。
// Call() 发送请求并等待回复。如果在超时间隔内收到回复，Call() 返回 true；否则返回 false。
// 因此 Call() 可能一段时间不会返回。返回 false 可能是由于服务器死亡、无法访问的活跃服务器、丢失的请求或丢失的回复。
//
// 除非服务器端的处理函数不返回，否则 Call() 保证会返回（可能延迟后）。因此不需要在 Call() 周围实现自己的超时。
//
// 有关更多详细信息，请查看 ../labrpc/labrpc.go 中的注释。
//
// 如果遇到 RPC 工作问题，请检查是否已将通过 RPC 传递的结构中的所有字段名大写，
// 并且调用方使用 & 传递回复结构的地址，而不是结构本身。
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

// 使用 Raft 的服务（例如 k/v 服务器）希望就下一个要追加到 Raft 日志中的命令达成一致。
// 如果此服务器不是领导者，返回 false。否则启动协议并立即返回。
// 不能保证此命令将永远提交到 Raft 日志，因为领导者可能失败或输掉选举。
// 即使 Raft 实例已被杀死，此函数也应优雅返回。
//
// 第一个返回值是命令如果被提交将出现在的索引。
// 第二个返回值是当前任期。
// 第三个返回值是此服务器是否认为自己是领导者。
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// 你的代码在这里 (3B)。

	return index, term, isLeader
}

// 测试程序不会在每次测试后停止由 Raft 创建的 goroutine，
// 但会调用 Kill() 方法。你的代码可以使用 killed() 来检查是否已调用 Kill()。
// 使用 atomic 可以避免需要锁。
//
// 问题是长时间运行的 goroutine 会使用内存并可能消耗 CPU 时间，
// 可能导致后续测试失败并生成令人困惑的调试输出。
// 任何具有长时间运行循环的 goroutine 应调用 killed() 检查是否应停止。
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// 如有需要，你的代码在这里。
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

func (rf *Raft) ticker() {
	for rf.killed() == false {

		// 你的代码在这里 (3A)
		// 检查是否应开始领导者选举。

		// 暂停随机时间，介于 50 到 350 毫秒之间。
		ms := 50 + (rand.Int63() % 300)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

// 服务或测试程序希望创建一个 Raft 服务器。
// 所有 Raft 服务器（包括此服务器）的端口都在 peers[] 中。
// 此服务器的端口是 peers[me]。所有服务器的 peers[] 数组顺序相同。
// persister 是此服务器保存其持久状态的地方，并且最初保存最近保存的状态（如果有）。
// applyCh 是一个通道，测试程序或服务期望 Raft 通过它发送 ApplyMsg 消息。
// Make() 必须快速返回，因此它应为任何长时间运行的工作启动 goroutine。
func Make(peers []*labrpc.ClientEnd, me int,
	persister *tester.Persister, applyCh chan raftapi.ApplyMsg) raftapi.Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// 你的初始化代码在这里 (3A, 3B, 3C)。

	// 从崩溃前持久化的状态初始化
	rf.readPersist(persister.ReadRaftState())

	// 启动 ticker goroutine 以开始选举
	go rf.ticker()

	return rf
}