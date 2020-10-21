package httpclient

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type NoOxygenError struct {
	Count     int    `json:"count"`
	WarnLevel int    `json:"warn_level"`
	Msg       string `json:"msg"`
}

func (e NoOxygenError) Error() string {
	return fmt.Sprintf("warnlevel:%d, no oxygen(%d), %s", e.WarnLevel, e.Count, e.Msg)
}

func TestBuildJSONError(t *testing.T) {

	Convey("should can be handle a struct intput", t, func() {
		reader := bytes.NewBuffer([]byte(`{"count":1, "warn_level":100, "msg":"run"}`))
		_, err := buildJSONError(reader, NoOxygenError{})
		oxyErr, ok := err.(*NoOxygenError)
		So(ok, ShouldBeTrue)
		So(oxyErr.Count, ShouldEqual, 1)
		So(oxyErr.WarnLevel, ShouldEqual, 100)
		So(oxyErr.Msg, ShouldEqual, "run")
	})

	Convey("should can be handle a struct pointer intput", t, func() {
		reader := bytes.NewBuffer([]byte(`{"count":1, "warn_level":100, "msg":"run"}`))
		_, err := buildJSONError(reader, &NoOxygenError{})
		oxyErr, ok := err.(*NoOxygenError)
		So(ok, ShouldBeTrue)
		So(oxyErr.Count, ShouldEqual, 1)
		So(oxyErr.WarnLevel, ShouldEqual, 100)
		So(oxyErr.Msg, ShouldEqual, "run")
	})

}
