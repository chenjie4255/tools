package structfix

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type Names struct {
	Values []string
}

func (n *Names) FixNilArray() {
	if n.Values == nil {
		n.Values = []string{}
	}
}

type EmbedName struct {
	Name         Names
	Meanningless struct {
		Age int
	}
}

type Foo struct {
	Bars []string
	Name Names
}

func (f *Foo) FixNilArray() {
	if f.Bars == nil {
		f.Bars = []string{}
	}
}

func TestStructFix(t *testing.T) {
	Convey("修复普通结构体", t, func() {
		f := Foo{}
		So(f.Bars, ShouldBeNil)
		So(f.Name.Values, ShouldBeNil)
		err := FixNilArray(&f, false)
		So(err, ShouldBeNil)
		So(f.Name.Values, ShouldNotBeNil)
		So(f.Bars, ShouldNotBeNil)
	})

	Convey("修复嵌套结构体", t, func() {
		type Emb struct {
			Foo Foo
		}

		emb := Emb{}
		So(emb.Foo.Bars, ShouldBeNil)
		So(emb.Foo.Name.Values, ShouldBeNil)
		err := FixNilArray(&emb, false)
		So(err, ShouldBeNil)
		So(emb.Foo.Bars, ShouldNotBeNil)
		So(emb.Foo.Name.Values, ShouldNotBeNil)
	})

	Convey("匿名嵌套结构体", t, func() {
		type Emb struct {
			Foo
		}

		emb := Emb{}
		So(emb.Foo.Bars, ShouldBeNil)
		err := FixNilArray(&emb, false)
		So(err, ShouldBeNil)
		So(emb.Foo.Bars, ShouldNotBeNil)
	})

	Convey("数组嵌套结构体", t, func() {
		type Emb struct {
			Foos []Foo
		}

		emb := Emb{}
		emb.Foos = append(emb.Foos, Foo{}, Foo{})
		So(emb.Foos, ShouldHaveLength, 2)
		So(emb.Foos[0].Bars, ShouldBeNil)
		So(emb.Foos[1].Bars, ShouldBeNil)
		err := FixNilArray(&emb, false)
		So(err, ShouldBeNil)
		So(emb.Foos[0].Bars, ShouldNotBeNil)
		So(emb.Foos[1].Bars, ShouldNotBeNil)
	})

	Convey("2层嵌套结构", t, func() {
		type Emb struct {
			Foos []Foo
			EN   EmbedName
		}

		emb := Emb{}

		So(emb.EN.Name.Values, ShouldBeNil)
		err := FixNilArray(&emb, false)
		So(err, ShouldBeNil)
		So(emb.EN.Name.Values, ShouldNotBeNil)
	})
}
