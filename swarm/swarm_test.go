package swarm

import (
	"fmt"
	"testing"
	"time"

	"github.com/chenjie4255/tools/dlock"
	"github.com/chenjie4255/tools/testenv"
)

var sum = 0

type testNode struct {
}

func (n *testNode) OnLeaderInterval() error {
	sum++
	return nil
}

func TestSwarm(t *testing.T) {

	if testing.Short() {
		t.Skip("skip integrated test in short mode")
	}
	env := testenv.GetIntegratedTestEnv()
	if env.RedisHost == "" {
		t.Skip("env not configured yet, skip this test")
	}

	dlock := dlock.New(env.RedisHost, env.RedisPassword, 0)
	fw := NewLeaderFramework(1, 2, "gotest", dlock)

	for i := 0; i < 100; i++ {
		node := testNode{}
		fw.AddNode(&node)
	}

	fw.Up()

	time.Sleep(5 * time.Second)

	if sum != 5 {
		t.Fatalf("sum should equal 5, but:%d ", sum)
	}
}

type brokenTestNode struct {
	num int
}

func (n *brokenTestNode) OnLeaderInterval() error {
	if sum < 1 {
		sum++
	} else {
		return fmt.Errorf("[%d]some error", n.num)
	}
	return nil
}

func TestSwarmLeaderBroken(t *testing.T) {

	if testing.Short() {
		t.Skip("skip integrated test in short mode")
	}
	env := testenv.GetIntegratedTestEnv()
	if env.RedisHost == "" {
		t.Skip("env not configured yet, skip this test")
	}

	dlock := dlock.New(env.RedisHost, env.RedisPassword, 0)
	fw := NewLeaderFramework(1, 2, "gotest", dlock)

	for i := 0; i < 100; i++ {
		node := brokenTestNode{i}
		fw.AddNode(&node)
	}

	sum = 0
	fw.Up()

	time.Sleep(5 * time.Second)

	if sum != 1 {
		t.Fatalf("sum should equal 1, but:%d ", sum)
	}
}
