package swarm

import (
	"fmt"
	"time"

	"github.com/chenjie4255/tools/dlock"
	"github.com/chenjie4255/tools/gor"
	"github.com/chenjie4255/tools/log"
	"github.com/chenjie4255/tools/worker"
)

var logger *log.Logger

func init() {
	logger = log.NewLoggerWithSentry("swarm_framework")
}

type Node interface {
	OnLeaderInterval() error
}

type LeaderFramework interface {
	AddNode(n Node)

	Up()
}

func JoinSwarm(node Node, swarmName string, dlLock dlock.DLock, leaderWorkTime int) {
	gor.RunWithRecover(func() {
		for {
			key := fmt.Sprintf("swarm_distributed_lock_%s", swarmName)
			lock, err := dlLock.GetLock(key, leaderWorkTime)
			if err == nil {
				if err := node.OnLeaderInterval(); err != nil {
					dlLock.DelLock(key, lock)
				}
			}
			time.Sleep(time.Second)
		}
	})
}

type leaderFrameworkNode struct {
	n                 Node
	leaderLockKey     string
	leaderLockSecret  string
	dLock             dlock.DLock
	leaderOnBoardTime int
}

func (n *leaderFrameworkNode) doLeaderWork() {
	gor.RunWithRecover(func() {
		if err := n.n.OnLeaderInterval(); err != nil {
			// 如果leader worker失败，则退出leader
			logger.AddFile().WithFields(log.Fields{
				"error": err,
			}).Warn("leader failed to finish it's job")
			n.giveupLeader()
		}
	})
}

func (n *leaderFrameworkNode) giveupLeader() {
	secret := n.leaderLockSecret
	n.leaderLockSecret = ""
	n.dLock.DelLock(n.leaderLockKey, secret)
}

func (n *leaderFrameworkNode) onSelectionInterval() {
	if n.leaderLockSecret != "" {
		if err := n.dLock.SetExpiredTime(n.leaderLockKey, n.leaderLockSecret, n.leaderOnBoardTime); err != nil {
			// lose key
			logger.AddFile().WithFields(log.Fields{
				"error": err,
			}).Error("leader has lost it's controlling")
			n.giveupLeader()
		} else {
			n.doLeaderWork()
		}
	} else {
		// try to be leaderele
		secret, err := n.dLock.GetLock(n.leaderLockKey, n.leaderOnBoardTime)
		if err == nil && secret != "" {
			n.leaderLockSecret = secret
			n.doLeaderWork()
			logger.AddFile().Info("new leader was elected")
		}
	}
}

type leaderFramework struct {
	electionInterval  int
	leaderOnBoardTime int
	name              string
	lockKey           string
	dLock             dlock.DLock
	nodes             []*leaderFrameworkNode
}

func NewLeaderFramework(electionInterval, leaderOnBoardTime int, name string, dl dlock.DLock) LeaderFramework {
	if electionInterval >= leaderOnBoardTime {
		panic("electionInterval should be less than leaderOnBoardTime")
	}
	fw := leaderFramework{}
	fw.electionInterval = electionInterval
	fw.leaderOnBoardTime = leaderOnBoardTime
	fw.name = name
	fw.dLock = dl
	fw.lockKey = "swarm_framework_" + name

	return &fw
}

func (f *leaderFramework) AddNode(n Node) {
	node := leaderFrameworkNode{}
	node.n = n
	node.dLock = f.dLock
	node.leaderLockKey = f.lockKey
	node.leaderOnBoardTime = f.leaderOnBoardTime

	f.nodes = append(f.nodes, &node)
}

func (f *leaderFramework) Up() {
	go func() {
		for {
			pool := worker.NewPool(32)
			nodeFn := func(node *leaderFrameworkNode) func() {
				return func() {
					defer pool.JobDone()
					node.onSelectionInterval()
				}
			}
			for i := 0; i < len(f.nodes); i++ {
				pool.AddJob(nodeFn(f.nodes[i]))
			}
			pool.WaitAll()
			pool.Release()
			time.Sleep(time.Duration(f.electionInterval) * time.Second)
		}
	}()
}
