package lock

import (
	m_logger "6.5840/kvsrv1/log"
	"6.5840/kvsrv1/rpc"
	"6.5840/kvtest1"
)

type Lock struct {
	// IKVClerk is a go interface for k/v clerks: the interface hides
	// the specific Clerk type of ck but promises that ck supports
	// Put and Get.  The tester passes the clerk in when calling
	// MakeLock().
	ck kvtest.IKVClerk
	l string
	owner string
	// You may add code here
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// Use l as the key to store the "lock state" (you would have to decide
// precisely what the lock state is).
func MakeLock(ck kvtest.IKVClerk, l string) *Lock {
	lk := &Lock{ck: ck,l: l,owner: kvtest.RandValue(8)}
	// You may add code here
	return lk
}

func (lk *Lock) Acquire() {
	for {
		owner,version,err := lk.ck.Get(lk.l)
		if err == rpc.OK {
			if owner != "" {
				// 有主
				m_logger.Log("The lock is acquired! Try later!")
				if owner == lk.owner {
					break
				} else {
					continue
				}
			} else {
				// 无主
				ok := lk.ck.Put(lk.l,lk.owner,version)
				if ok == rpc.OK {
					m_logger.Log("Lock successfully!")
					break
				}
			}
		} else if err == rpc.ErrNoKey {
			// 无主
			ok := lk.ck.Put(lk.l,lk.owner,version)
			if ok == rpc.OK {
				m_logger.Log("Lock successfully!")
				break
			}
		}
	}
}

func (lk *Lock) Release() {
	owner,version,err := lk.ck.Get(lk.l)
	if err == rpc.ErrNoKey || owner == "" {
		return
	} else if owner == lk.owner {
		ok := lk.ck.Put(lk.l,"",version)
		if ok == rpc.OK {
			m_logger.Log("Lock released!")
		}
	}
}