package redis

import (
	"github.com/chenjie4255/tools/testenv"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShareInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrated test in short mode")
	}
	env := testenv.GetIntegratedTestEnv()
	if env.RedisHost == "" {
		t.Skip("env not configured yet, skip this test")
	}

	FlushDB(env.RedisHost, "", 3)
	db := NewDB(env.RedisHost, "", 3)

	shareInfo := NewShareInfo("test_k", db, 100)
	shareInfo1 := NewShareInfo("test_k", db, 100)
	shareInfo2 := NewShareInfo("test_k", db, 100)
	shareInfo3 := NewShareInfo("test_k", db, 100)
	shareInfo4 := NewShareInfo("test_k", db, 100)
	shareInfo5 := NewShareInfo("test_k", db, 100)
	shareInfo6 := NewShareInfo("test_k", db, 100)

	err := shareInfo.Set("hi")
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	val, err := shareInfo1.Get()
	assert.NoError(t, err)
	assert.Equal(t, "hi", string(val))

	val, err = shareInfo2.Get()
	assert.NoError(t, err)
	assert.Equal(t, "hi", string(val))

	val, err = shareInfo3.Get()
	assert.NoError(t, err)
	assert.Equal(t, "hi", string(val))

	val, err = shareInfo4.Get()
	assert.NoError(t, err)
	assert.Equal(t, "hi", string(val))

	val, err = shareInfo5.Get()
	assert.NoError(t, err)
	assert.Equal(t, "hi", string(val))

	val, err = shareInfo6.Get()
	assert.NoError(t, err)
	assert.Equal(t, "hi", string(val))

	err = shareInfo1.Del("hi2")
	assert.Error(t, err)

	err = shareInfo1.Del("hi")
	assert.NoError(t, err)

	val, err = shareInfo1.Get()
	assert.NoError(t, err)
	assert.Equal(t, "", string(val))

	val, err = shareInfo2.Get()
	assert.NoError(t, err)
	assert.Equal(t, "", string(val))

	val, err = shareInfo3.Get()
	assert.NoError(t, err)
	assert.Equal(t, "", string(val))

	val, err = shareInfo4.Get()
	assert.NoError(t, err)
	assert.Equal(t, "", string(val))

	val, err = shareInfo5.Get()
	assert.NoError(t, err)
	assert.Equal(t, "", string(val))

	val, err = shareInfo6.Get()
	assert.NoError(t, err)
	assert.Equal(t, "", string(val))
}
