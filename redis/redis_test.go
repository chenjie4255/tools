package redis

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/chenjie4255/tools/errcode"

	"github.com/chenjie4255/tools/testenv"
	"github.com/chenjie4255/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCacheOp(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrated test in short mode")
	}
	env := testenv.GetIntegratedTestEnv()
	if env.RedisHost == "" {
		t.Skip("env not configured yet, skip this test")
	}

	type Item struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	Convey("create cache", t, func() {
		FlushDB(env.RedisHost, "", 3)
		cache := NewDB(env.RedisHost, "", 3)

		Convey("test decr", func() {
			So(cache.Set("kk1", 100, 1000), ShouldBeNil)

			ret, err := cache.CmpGTDecr("kk1", 98)
			So(err, ShouldBeNil)
			So(ret, ShouldEqual, 99)

			ret, err = cache.CmpGTDecr("kk1", 98)
			So(err, ShouldBeNil)
			So(ret, ShouldEqual, 98)

			for i := 0; i < 100; i++ {
				_, err := cache.CmpGTDecr("k1", 98)
				So(err, ShouldNotBeNil)
			}
		})

		Convey("test decrby", func() {
			So(cache.Set("kk1", 100, 1000), ShouldBeNil)

			ret, err := cache.DecrByUint64("kk1", 1)
			So(err, ShouldBeNil)
			So(ret, ShouldEqual, 99)

			ret, err = cache.DecrByUint64("kk1", 1)
			So(err, ShouldBeNil)
			So(ret, ShouldEqual, 98)
		})

		Convey("test getset op", func() {
			err := cache.SetNotExists("k1", 1, 10)
			So(err, ShouldBeNil)

			exists, err := cache.Exists("k1")
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)

			exists, err = cache.Exists("k2")
			So(err, ShouldBeNil)
			So(exists, ShouldBeFalse)

			err = cache.SetNotExists("k1", 2, 10)
			So(err, ShouldNotBeNil)
			So(errors.FindTag(err, errcode.ResExisted), ShouldBeTrue)

			check := 0
			err = cache.Get("k1", &check)
			So(err, ShouldBeNil)
			So(check, ShouldEqual, 1)
		})

		Convey("add cache item", func() {
			item := Item{"item", 1.99}
			err := cache.Set("item_01", item, 10)
			So(err, ShouldBeNil)

			Convey("delete item by key", func() {
				err := cache.Del("item_01")
				So(err, ShouldBeNil)

				check := Item{}
				err = cache.Get("item_01", &check)
				So(err, ShouldNotBeNil)
				So(errors.FindTag(err, errcode.ResNotFound), ShouldBeTrue)
			})

			Convey("delete item by keys", func() {
				err := cache.DelKeys([]string{"item_01"})
				So(err, ShouldBeNil)

				check := Item{}
				err = cache.Get("item_01", &check)
				So(err, ShouldNotBeNil)
				So(errors.FindTag(err, errcode.ResNotFound), ShouldBeTrue)
			})

			Convey("added item should can be retrieved", func() {
				check := Item{}
				err := cache.Get("item_01", &check)
				So(err, ShouldBeNil)
				So(check, ShouldResemble, item)
			})

			Convey("support delete key with value", func() {
				err := cache.DelKeyForValue("item_01", item)
				So(err, ShouldBeNil)

				check := Item{}
				err = cache.Get("item_01", &check)
				So(err, ShouldNotBeNil)
				So(errors.FindTag(err, errcode.ResNotFound), ShouldBeTrue)
			})

			Convey("delete key with value must have same value", func() {
				item.Name = "noname"
				err := cache.DelKeyForValue("item_01", item)
				So(err, ShouldNotBeNil)
			})

			Convey("get an not existed item should be failed", func() {
				check := Item{}
				err := cache.Get("whatever", &check)
				So(err, ShouldNotBeNil)
				So(errors.FindTag(err, errcode.ResNotFound), ShouldBeTrue)
			})

			Convey("item should can be removed by keys*", func() {
				err := cache.DelByKeys("ite*")
				So(err, ShouldBeNil)

				check := Item{}
				err = cache.Get("item_01", &check)
				So(err, ShouldNotBeNil)
				So(errors.FindTag(err, errcode.ResNotFound), ShouldBeTrue)
			})

			Convey("iter should can be removed by scan keys*", func() {
				err := cache.DelKeysByScan("ite*")
				So(err, ShouldBeNil)

				check := Item{}
				err = cache.Get("item_01", &check)
				So(err, ShouldNotBeNil)
				So(errors.FindTag(err, errcode.ResNotFound), ShouldBeTrue)
			})

			Convey("item should can be removed by it's key", func() {
				err := cache.Del("item_01")
				So(err, ShouldBeNil)

				check := Item{}
				err = cache.Get("item_01", &check)
				So(err, ShouldNotBeNil)
				So(errors.FindTag(err, errcode.ResNotFound), ShouldBeTrue)
			})
		})

		Convey("int operator ", func() {
			err := cache.Set("123", 123, 0)
			So(err, ShouldBeNil)

			check := 0
			err = cache.Get("123", &check)
			So(err, ShouldBeNil)
			So(check, ShouldEqual, 123)
		})

		Convey("test cache get", func() {
			getCount := 0
			getFn := func() (interface{}, error) {
				if getCount == 0 {
					getCount++
					return 100, nil
				}
				return 0, errors.New("whatever")
			}

			check := 0

			err := cache.CacheGet("a", &check, getFn, 100)
			So(err, ShouldBeNil)
			So(check, ShouldEqual, 100)

			time.Sleep(1 * time.Second)

			check2 := 0
			err = cache.CacheGet("a", &check2, getFn, 100)
			So(err, ShouldBeNil)
			So(check2, ShouldEqual, 100)
		})

		Convey("test set not exists", func() {
			err := cache.SetNotExists("k1", nil, 1)
			So(err, ShouldBeNil)

			err = cache.SetNotExists("k1", nil, 100)
			So(err, ShouldNotBeNil)

			time.Sleep(2 * time.Second)

			err = cache.SetNotExists("k1", nil, 100)
			So(err, ShouldBeNil)
		})

		Convey("test sortset op", func() {
			err := cache.AddSortSetStr("ks", "1", 100)
			So(err, ShouldBeNil)

			count, err := cache.GetSortSetCount("ks", 0, 100)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 1)

			err = cache.AddSortSetStr("ks", "2", 200)
			So(err, ShouldBeNil)

			count, err = cache.GetSortSetCount("ks", 0, 200)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 2)

			result, err := cache.GetSortSetRangeStr("ks", 0, 200)
			So(err, ShouldBeNil)
			So(result, ShouldHaveLength, 2)
			So(result[0], ShouldEqual, "1")
			So(result[1], ShouldEqual, "2")

			result, err = cache.GetSortSetRangeStr("ks", 110, 200)
			So(err, ShouldBeNil)
			So(result, ShouldHaveLength, 1)
			So(result[0], ShouldEqual, "2")

			count, err = cache.RemoveSortSet("ks", 0, 100)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 1)

			count, err = cache.GetSortSetCount("ks", 0, 200)
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 1)
		})

		Convey("test list op", func() {
			err := cache.PushStringList("kl", "v1", 100)
			So(err, ShouldBeNil)

			checks, err := cache.GetStringList("kl")
			So(err, ShouldBeNil)
			So(checks, ShouldHaveLength, 1)

			err = cache.PushStringList("kl", "v2", 100)
			So(err, ShouldBeNil)

			checks, err = cache.GetStringList("kl")
			So(err, ShouldBeNil)
			So(checks, ShouldHaveLength, 2)

			checks, err = cache.GetStringList("kl2")
			So(err, ShouldBeNil)
			So(checks, ShouldHaveLength, 0)
		})

		Convey("test ttl", func() {
			err := cache.Set("k1", 200, 100)
			So(err, ShouldBeNil)

			ttl, err := cache.TTL("k1")
			So(err, ShouldBeNil)
			So(ttl, ShouldBeGreaterThan, 98)

			time.Sleep(2 * time.Second)
			ttl, err = cache.TTL("k1")
			So(err, ShouldBeNil)
			So(ttl, ShouldBeLessThan, 100)

			ttl, err = cache.TTL("kxxx")
			So(err, ShouldNotBeNil)
			So(errors.FindTag(err, errcode.ResNotFound), ShouldBeTrue)

			err = cache.Set("k2", 200, 0)
			So(err, ShouldBeNil)

			ttl, err = cache.TTL("k2")
			So(err, ShouldBeNil)
			So(ttl, ShouldEqual, -1)
		})

		Convey("test replace value", func() {
			err := cache.Set("k1", 100, 0)
			So(err, ShouldBeNil)

			err = cache.ReplaceValue("k1", 100, 101)
			So(err, ShouldBeNil)

			err = cache.ReplaceValue("k1", 101, 102)
			So(err, ShouldBeNil)

			err = cache.ReplaceValue("k1", 103, 104)
			So(err, ShouldNotBeNil)
		})

		Convey("test cmp and set", func() {
			err := cache.Set("k1", 100, 0)
			So(err, ShouldBeNil)

			err = cache.CmpAndSet("k1", 100, "k2", 101)
			So(err, ShouldBeNil)

			var check int
			err = cache.Get("k2", &check)
			So(err, ShouldBeNil)
			So(check, ShouldEqual, 101)

			err = cache.CmpAndSet("k2", 102, "k3", 103)
			So(err, ShouldNotBeNil)
		})

		Convey("set uint64 op", func() {
			val, err := cache.IncrByUint64("k64", 100)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100)

			val, err = cache.IncrByUint64("k64", 101)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 201)

			val, err = cache.IncrToUint64("k64", 200)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)

			val, err = cache.IncrToUint64("k64", 202)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 202)
		})

		Convey("GET TIME", func() {
			val, err := cache.Time()
			So(err, ShouldBeNil)
			So(val, ShouldBeGreaterThan, 0)
		})
	})
}

func TestSliceOp(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrated test in short mode")
	}
	env := testenv.GetIntegratedTestEnv()
	if env.RedisHost == "" {
		t.Skip("env not configured yet, skip this test")
	}

	type Item struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	Convey("create cache", t, func() {
		FlushDB(env.RedisHost, "", 3)
		db := NewDB(env.RedisHost, "", 3)

		item1 := Item{Name: "1", Price: 1.2}
		item2 := Item{Name: "2", Price: 1.3}

		items := []interface{}{item1, item2}

		itemData, _ := MarshalJSONSlice(items)

		err := db.SAdd("skey", itemData...)
		So(err, ShouldBeNil)

		count, err := db.SCard("skey")
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 2)

		checks, err := db.SRandMember("skey", 2)
		So(err, ShouldBeNil)
		So(checks, ShouldHaveLength, 2)

		itemRet := make([]Item, len(checks))
		itemRetIfs := []interface{}{}
		for i := range itemRet {
			itemRetIfs = append(itemRetIfs, &itemRet[i])
		}

		err = UnMarshalJSONSlice(checks, itemRetIfs)
		So(err, ShouldBeNil)

		fmt.Println(itemRet)
	})
}

func TestPubSub(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integrated test in short mode")
	}
	env := testenv.GetIntegratedTestEnv()
	if env.RedisHost == "" {
		t.Skip("env not configured yet, skip this test")
	}

	type Item struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	FlushDB(env.RedisHost, "", 3)
	db := NewDB(env.RedisHost, "", 3)

	sendData := "abc"

	count := 0

	wg := sync.WaitGroup{}
	wg.Add(1)

	done := make(chan bool)

	go db.Subscribe([]string{"s1", "s2"}, done, func(name string, data []byte) {
		t.Logf("%s, %s", name, string(data))
		if name != "s1" {
			t.Errorf("invalid name, %s", name)
		}

		if string(data) != sendData {
			t.Errorf("invalid data, %s", string(data))
		}
		count++
		wg.Done()
	})

	val := <-done
	if !val {
		t.Errorf("should done")
	}

	cc, err := db.Publish("s1", []byte(sendData))
	wg.Wait()

	if cc == 0 {
		t.Errorf("cc %d", cc)
	}

	if err != nil {
		t.Errorf("pub err:%s", err)
	}

	if count != 1 {
		t.Errorf("count should equal 1, %d", count)
	}
}
