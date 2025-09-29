package kvsrv

import (
	// "fmt"
	"log"
	"sync"

	"6.5840/kvsrv1/rpc"
	"6.5840/labrpc"
	tester "6.5840/tester1"

	m_logger "6.5840/kvsrv1/log"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu sync.Mutex
	kv map[string]ValData
}

type ValData struct {
	val   string
	vsion rpc.Tversion
}

func MakeKVServer() *KVServer {
	kv := &KVServer{kv: make(map[string]ValData)}
	// Your code here.
	return kv
}

// Get returns the value and version for args.Key, if args.Key
// exists. Otherwise, Get returns ErrNoKey.
func (kv *KVServer) Get(args *rpc.GetArgs, reply *rpc.GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	val, isExist := kv.kv[args.Key]
	if isExist {
		reply.Value = val.val
		reply.Version = val.vsion
		reply.Err = rpc.OK
	} else {
		reply.Err = rpc.ErrNoKey
	}
	// version_Str := fmt.Sprintln(" version: ",val.vsion)
	// m_logger.Log("func Get called!\n" + "key: " + args.Key + " val: " + val.val + version_Str)
}

// Update the value for a key if args.Version matches the version of
// the key on the server. If versions don't match, return ErrVersion.
// If the key doesn't exist, Put installs the value if the
// args.Version is 0, and returns ErrNoKey otherwise.
func (kv *KVServer) Put(args *rpc.PutArgs, reply *rpc.PutReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	// version := fmt.Sprintln("Version:", args.Version)
	// m_logger.Log("func Put() Called!\n" + "Key: " + args.Key + " Value: " + args.Value + " " + version)
	if kv.kv[args.Key].vsion == args.Version {
		kv.kv[args.Key] = ValData{val: args.Value, vsion: args.Version+1}
		reply.Err = rpc.OK
	} else if kv.kv[args.Key].vsion == 0 && args.Version != 0 {
		reply.Err = rpc.ErrNoKey
	} else {
		reply.Err = rpc.ErrVersion
	}
}

// You can ignore Kill() for this lab
func (kv *KVServer) Kill() {
	m_logger.Log_import("A KVserver 's Kill() method was called!")
}

// You can ignore all arguments; they are \for replicated KVservers
func StartKVServer(ends []*labrpc.ClientEnd, gid tester.Tgid, srv int, persister *tester.Persister) []tester.IService {
	m_logger.Log("KVServer is starting ...")

	kv := MakeKVServer()
	return []tester.IService{kv}
}