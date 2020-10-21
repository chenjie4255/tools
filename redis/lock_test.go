package redis

import (
	"testing"
	"time"

	"github.com/chenjie4255/gouuid"

	"github.com/chenjie4255/tools/testenv"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDLock(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrated test in short mode")
	}
	env := testenv.GetIntegratedTestEnv()
	if env.RedisHost == "" {
		t.Skip("env TEST_REDIS_HOST cannot be found, skip this test")
	}

	Convey("building test env", t, func() {
		FlushDB(env.RedisHost, "", 0)
		lock := NewDB(env.RedisHost, env.RedisPassword, 0)

		uv4, _ := uuid.NewV4()
		randomKey := uv4.NoHyphenString()

		Convey("get key should ok", func() {
			lockID, err := lock.GetLock(randomKey, 2)
			So(err, ShouldBeNil)
			So(lockID, ShouldNotBeBlank)

			Convey("another try get key should fail", func() {
				check, err := lock.GetLock(randomKey, 1)
				So(err, ShouldNotBeNil)
				So(check, ShouldBeEmpty)
			})

			Convey("another key should be ok", func() {
				uv4, _ := uuid.NewV4()
				check, err := lock.GetLock(uv4.NoHyphenString(), 1)
				So(err, ShouldBeNil)
				So(check, ShouldNotBeBlank)
			})

			Convey("delete an unexisted key should be fail", func() {
				err := lock.DelLock("xkey", "aaa")
				So(err, ShouldNotBeNil)
			})

			Convey("key should can be deleted within 2s", func() {
				err := lock.DelLock(randomKey, lockID)
				So(err, ShouldBeNil)
			})

			Convey("reset exipre time as 3s and sleep 2s", func() {
				lock.SetLockTTL(randomKey, lockID, 3)
				time.Sleep(2 * time.Second)

				Convey("key should cannot be deleted", func() {
					err := lock.DelLock(randomKey, lockID)
					So(err, ShouldBeNil)
				})
			})

			Convey("after 3s", func() {
				time.Sleep(3 * time.Second)

				Convey("should cannot be renew expired time", func() {
					err := lock.SetLockTTL(randomKey, lockID, 1)
					So(err, ShouldNotBeNil)
				})

				Convey("key should cannot be delete by old nonce value", func() {
					err := lock.DelLock(randomKey, lockID)
					So(err, ShouldNotBeNil)
				})

				Convey("key should be get again", func() {
					uv4, _ := uuid.NewV4()
					check, err := lock.GetLock(uv4.NoHyphenString(), 1)
					So(err, ShouldBeNil)
					So(check, ShouldNotBeBlank)
				})
			})
		})

	})

}
