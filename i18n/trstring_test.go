package i18n

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTranslactObject(t *testing.T) {

	type Follower struct {
		Name TrString
	}
	type PoorCat struct {
		Name  TrString `json:"name"`
		Lover struct {
			Alias TrString
		}
		Followers   []Follower
		Foo         string
		Hehe        int
		xx          TrString
		VIPFollower *Follower
	}

	cat := PoorCat{}
	cat.Name.En = "Chamberlain"
	cat.Name.ZhHans = "张九零"
	cat.Name.ZhHant = "张志中"
	cat.Lover.Alias.En = "Liverpool"
	cat.Lover.Alias.ZhHans = "利物浦"
	cat.Lover.Alias.ZhHant = "马桶"

	follower := Follower{}
	follower.Name.En = "dog"
	follower.Name.ZhHans = "小白"
	follower.Name.ZhHant = "小黑"

	follower2 := Follower{}
	follower2.Name.En = "cat"
	follower2.Name.ZhHans = "小猫"
	follower2.Name.ZhHant = "小小猫"

	cat.Followers = []Follower{follower, follower2}
	cat.xx.En = "Mhxx"
	cat.xx.ZhHant = "怪物猎人"
	cat.xx.ZhHant = "怪物人猎"

	vipFollower := Follower{}
	vipFollower.Name.En = "vipf"
	vipFollower.Name.ZhHans = "重要会员"
	vipFollower.Name.ZhHant = "重要会员？"
	vipFollower.Name.Ja = "Ja"
	vipFollower.Name.Es = "Es"

	cat.VIPFollower = &vipFollower

	err := SetLangForObject(&cat, "zh-hans")
	if err != nil {
		t.Fatalf("failed to set lang for object, %s", err)
	}

	if cat.Name.Tr != "张九零" {
		t.Fatalf("failed to set lang for object, expect:%s, got:%s", cat.Name.ZhHans, cat.Name.Tr)
	}

	if cat.Lover.Alias.Tr != "利物浦" {
		t.Fatalf("failed to set lang for object, expect:%s, got:%s", cat.Lover.Alias.ZhHans, cat.Lover.Alias.Tr)
	}
	if cat.Followers[0].Name.Tr != "小白" {
		t.Fatalf("failed to set lang for object, expect:小白, got:%s", cat.Followers[0].Name.Tr)
	}
	if cat.Followers[1].Name.Tr != "小猫" {
		t.Fatalf("failed to set lang for object, expect:小猫, got:%s", cat.Followers[1].Name.Tr)
	}

	if cat.xx.Tr != "" {
		t.Fatalf("failed to set lang for object, expect empty, got:%s", cat.xx.Tr)
	}
	if cat.VIPFollower.Name.Tr != "重要会员" {
		t.Fatalf("failed to set lang for object, expect 重要会员, got:%s", cat.VIPFollower.Name.Tr)
	}
}

func TestTranslactArray(t *testing.T) {
	type PoolCat struct {
		Name TrString
	}

	cat := PoolCat{}
	cat.Name.En = "sunderlan"
	cat.Name.ZhHans = "黑猫"
	cat.Name.ZhHant = "桑德兰"

	cats := []PoolCat{cat}

	err := SetLangForObject(cats, "en")
	if err != nil {
		t.Fatalf("failed to set lang for object, %s", err)
	}

	if cats[0].Name.Tr != "sunderlan" {
		t.Fatalf("failed to set lang for object, except:sunderlan, got:%s", cat.Name.Tr)
	}

}

func TestLanguageTranslate(t *testing.T) {
	type PoolCat struct {
		Name TrString
	}

	Convey("init transaclate data", t, func() {
		cat := PoolCat{}
		cat.Name.En = "sunderlan"
		cat.Name.ZhHans = "黑猫"
		cat.Name.ZhHant = "桑德兰"
		cat.Name.Ja = "秋田"
		cat.Name.Es = "斗牛"

		Convey("transaclate into spain", func() {
			err := SetLangForObject(&cat, "es")
			So(err, ShouldBeNil)
			So(cat.Name.Tr, ShouldEqual, "斗牛")
		})

		Convey("transaclte into japanese", func() {
			err := SetLangForObject(&cat, "ja")
			So(err, ShouldBeNil)
			So(cat.Name.Tr, ShouldEqual, "秋田")
		})
	})

}
