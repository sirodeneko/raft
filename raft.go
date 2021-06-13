package raft

import (
	"sync"
	"time"
)

type ApplyMsg struct {
	CommandValid bool // true为log，false为snapshot

	// 向application层提交日志
	Command      interface{}
	CommandIndex int
	CommandTerm  int

	// 向application层安装快照
	Snapshot          []byte
	LastIncludedIndex int
	LastIncludedTerm  int
}

// 日志项
type LogEntry struct {
	Command interface{}
	Term    int
}

// 当前角色
const ROLE_LEADER = "Leader"
const ROLE_FOLLOWER = "Follower"
const ROLE_CANDIDATES = "Candidates"

type Raft struct {
	mu        sync.Mutex
	peers     []*ClientEnd
	persister *Persister
	me        int
	dead      int32

	// 所有服务器，持久化状态（lab-2A不要求持久化）
	currentTerm       int        // 见过的最大任期
	votedFor          int        // 记录在currentTerm任期投票给谁了
	log               []LogEntry // 操作日志
	lastIncludedIndex int        // snapshot最后1个logEntry的index，没有snapshot则为0
	lastIncludedTerm  int        // snapthost最后1个logEntry的term，没有snaphost则无意义

	// 所有服务器，易失状态
	commitIndex int // 已知的最大已提交索引
	lastApplied int // 当前应用到状态机的索引

	// 仅Leader，易失状态（成为leader时重置）
	nextIndex  []int //	每个follower的log同步起点索引（初始为leader log的最后一项）
	matchIndex []int // 每个follower的log同步进度（初始为0），和nextIndex强关联

	// 所有服务器，选举相关状态
	role              string    // 身份
	leaderId          int       // leader的id
	lastActiveTime    time.Time // 上次活跃时间（刷新时机：收到leader心跳、给其他candidates投票、请求其他节点投票）
	lastBroadcastTime time.Time // 作为leader，上次的广播时间

	applyCh chan ApplyMsg // 应用层的提交队列
}
