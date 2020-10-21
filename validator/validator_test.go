package validator

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEmptyPostJSONValidation(t *testing.T) {
	Convey("Initial a empty post json request...", t, func() {
		jsonStr := []byte("{}")
		r, _ := http.NewRequest("POST", "www.whatever.com/xxx", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")

		Convey("test with all optional struct tag", func() {
			type LoginInfo struct {
				UserName     string `json:"username"`
				UserPassword string `json:"password"`
				UseCookie    bool   `json:"rememberme"`
			}

			var info LoginInfo
			err := ParsePostJSON(r, &info)
			Convey("test should be succeed", func() {
				So(err, ShouldBeNil)
				So(info.UseCookie, ShouldEqual, false)
				So(info.UserName, ShouldEqual, "")
				So(info.UserPassword, ShouldEqual, "")
			})
		})

		Convey("test with a required struct tag", func() {
			type LoginInfo struct {
				UserName     string `json:"username" valid:"required"`
				UserPassword string `json:"password"`
				UseCookie    bool   `json:"rememberme"`
			}

			var info LoginInfo
			err := ParsePostJSON(r, &info)
			Convey("test should be failed", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestPostJSONValidation(t *testing.T) {
	Convey("Initial a post json request...", t, func() {
		jsonStr := []byte(`{"username":"chenjie","password":"youguess","rememberme":true}`)
		r, _ := http.NewRequest("POST", "www.whatever.com/xxx", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")

		Convey("Test correct case", func() {
			type LoginInfo struct {
				UserName     string `json:"username" valid:"required"`
				UserPassword string `json:"password" valid:"required"`
				UseCookie    bool   `json:"rememberme" valid:"required"`
			}

			var info LoginInfo
			err := ParsePostJSON(r, &info)
			Convey("test should be succeed", func() {
				So(err, ShouldBeNil)
				So(info.UseCookie, ShouldEqual, true)
				So(info.UserName, ShouldEqual, "chenjie")
				So(info.UserPassword, ShouldEqual, "youguess")
			})
		})

		Convey("Test uncorrect case", func() {
			type LoginInfo struct {
				Username   string `json:"username" valid:"required"`
				Token      string `json:"token" valid:"required"`
				RememberMe bool   `json:"rememberme" valid:"required"`
			}

			var info LoginInfo
			err := ParsePostJSON(r, &info)
			Convey("Test should be failed", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Test struct with option field", func() {
			type LoginInfo struct {
				Username   string `json:"username" valid:"required"`
				Password   string `json:"password" valid:"required"`
				Token      string `json:"token,omitempty"`
				RememberMe bool   `json:"rememberme" valid:"required"`
			}

			var info LoginInfo
			err := ParsePostJSON(r, &info)
			Convey("Test should be failed", func() {
				So(err, ShouldBeNil)
				So(info.Token, ShouldBeEmpty)
				So(info.Username, ShouldEqual, "chenjie")
				So(info.Password, ShouldEqual, "youguess")
				So(info.RememberMe, ShouldBeTrue)
			})
		})
	})
}

func TestEmbedValid(t *testing.T) {
	Convey("Initial a post json request...", t, func() {
		type Core struct {
			ID   string `json:"id" valid:"numeric,required"`
			Name string `json:"name" valid:"required"`
		}

		type Person struct {
			Core  `json:",squash"`
			Title string `json:"title" valid:"required"`
		}

		correctP := Person{}
		correctP.ID = "123"
		correctP.Name = "123"
		correctP.Title = "321"

		So(ValidateObj(correctP), ShouldBeNil)

		errorP := Person{}
		errorP.Name = "bbb"
		errorP.Title = "ccc"

		So(ValidateObj(errorP), ShouldNotBeNil)
	})
}

func TestPostFormValidation(t *testing.T) {
	Convey("Initial test env...", t, func() {
		type Certification struct {
			Username string `json:"username" valid:"required"`
			Password string `json:"password" valid:"required"`
		}

		var c Certification

		Convey("Test succee case", func() {
			form := url.Values{}
			form.Add("username", "hehe")
			form.Add("password", "haahah")

			r, _ := http.NewRequest("POST", "www.baidu.com/upload", strings.NewReader(form.Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			Convey("validation should be succeed", func() {
				err := ParsePostForm(r, &c)
				So(err, ShouldBeNil)
				So(c.Username, ShouldEqual, "hehe")
				So(c.Password, ShouldEqual, "haahah")
			})
		})

		Convey("Test fail case(no password field)", func() {
			form := url.Values{}
			form.Add("username", "hehe")

			r, _ := http.NewRequest("POST", "www.baidu.com/upload", strings.NewReader(form.Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			Convey("validation should be failed", func() {
				err := ParsePostForm(r, &c)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Test fail case(empty password field)", func() {
			form := url.Values{}
			form.Add("username", "hehe")
			form.Add("password", "")

			r, _ := http.NewRequest("POST", "www.baidu.com/upload", strings.NewReader(form.Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			Convey("validation should be failed", func() {
				err := ParsePostForm(r, &c)
				So(err, ShouldNotBeNil)
			})
		})

	})
}

func TestURLValidation(t *testing.T) {
	Convey("Initial test env...", t, func() {
		u, err := url.Parse("http://bing.com/search?q=dotnet&key=hehe&key=haha&cb=www.baidu.com&timeout=12")
		Convey("initial url should be succeed", func() {
			So(err, ShouldBeNil)
		})

		type Params struct {
			Search   string   `json:"q" valid:"required"`
			Keys     []string `json:"key" valid:"required"`
			Callback string   `json:"cb"`
			Timeout  int      `json:"timeout"`
		}

		var params Params

		Convey("Validation should be succeed", func() {
			err := validateURLParams(u, &params)
			So(err, ShouldBeNil)
			So(params.Search, ShouldEqual, "dotnet")
			So(params.Keys[0], ShouldEqual, "hehe")
			So(params.Keys[1], ShouldEqual, "haha")
			So(params.Callback, ShouldEqual, "www.baidu.com")
			So(params.Timeout, ShouldEqual, 12)
		})

		Convey("When input has no search field", func() {
			u, _ = url.Parse("http://bing.com/search?key=hehe&key=haha&cb=www.google.com")
			err = validateURLParams(u, &params)
			Convey("validation should be failed", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When input has no cb field", func() {
			u, _ = url.Parse("http://bing.com/search?q=dotnet&key=hehe&key=haha")
			err = validateURLParams(u, &params)
			Convey("validation should be succeed", func() {
				So(err, ShouldBeNil)
				So(params.Callback, ShouldBeEmpty)
			})
		})

		Convey("test default value", func() {
			u, _ = url.Parse("http://bing.com/search?q=dotnet&key=hehe&key=haha")
			var params Params
			params.Callback = "123"

			err = validateURLParams(u, &params)
			Convey("validation should be succeed", func() {
				So(err, ShouldBeNil)
				So(params.Callback, ShouldEqual, "123")
			})
		})

		Convey("test int multi param value", func() {
			type IntParams struct {
				Keys []int64 `json:"key"`
			}
			u, _ = url.Parse("http://bing.com/search?q=dotnet&key=1&key=2")
			var params IntParams
			err = validateURLParams(u, &params)
			Convey("validation should be succeed", func() {
				So(err, ShouldBeNil)
				So(params.Keys, ShouldHaveLength, 2)
			})
		})

		Convey("test string multi param value", func() {
			type IntParams struct {
				Status []string `json:"status"`
			}

			var params IntParams

			u, _ = url.Parse("http://bing.com/search?status=published")
			err = validateURLParams(u, &params)
			Convey("validation should be succeed", func() {
				So(err, ShouldBeNil)
				So(params.Status, ShouldHaveLength, 1)
			})
		})

		Convey("test empty int multi param value", func() {
			type IntParams struct {
				Keys []int64 `json:"key"`
			}
			u, _ = url.Parse("http://bing.com/search?q=dotnet")
			var params IntParams
			err = validateURLParams(u, &params)
			Convey("validation should be succeed", func() {
				So(err, ShouldBeNil)
				So(params.Keys, ShouldHaveLength, 0)
			})
		})

		Convey("test one int multi param value", func() {
			type IntParams struct {
				Keys []int64 `json:"key"`
			}
			u, _ = url.Parse("http://bing.com/search?q=dotnet&key=123")
			var params IntParams
			err = validateURLParams(u, &params)
			Convey("validation should be succeed", func() {
				So(err, ShouldBeNil)
				So(params.Keys, ShouldHaveLength, 1)
				So(params.Keys[0], ShouldEqual, 123)
			})
		})
	})
}

func TestPageValidator(t *testing.T) {
	Convey("测试分页类解析", t, func() {
		type Params struct {
			Offset int `json:"offset" valid:"range(0|100),optional"`
			Size   int `json:"size" valid:"range(0|100),optional"`
		}

		Convey("应支持正常用例", func() {
			var params Params
			u, err := url.Parse("http://bing.com/search?offset=50&size=100")
			So(err, ShouldBeNil)
			err = ParseURLParams(u, &params)
			So(err, ShouldBeNil)
			So(params.Offset, ShouldEqual, 50)
			So(params.Size, ShouldEqual, 100)
		})

		Convey("应支持参数为空的URL", func() {
			var params Params
			u, err := url.Parse("http://bing.com/search")
			So(err, ShouldBeNil)
			err = validateURLParams(u, &params)
			So(err, ShouldBeNil)
			Convey("应支持默认值设置", func() {
				params := Params{10, 20}
				err = validateURLParams(u, &params)
				So(err, ShouldBeNil)
				So(params.Offset, ShouldEqual, 10)
				So(params.Size, ShouldEqual, 20)
			})
		})

		Convey("错误的参数应被检测出来", func() {
			type Par struct {
				Offset int `json:"offset" valid:"range(0|9999999)"`
				Size   int `json:"size" valid:"range(0|100)"`
			}

			Convey("错误的范围上限应该被检测出来", func() {
				var params Par
				u, err := url.Parse("http://bing.com/search?offset=-1&size=10")
				So(err, ShouldBeNil)
				err = ParseURLParams(u, &params)
				So(err, ShouldNotBeNil)
			})

			Convey("错误的范围下限应该被检测出来", func() {
				var params Par
				u, err := url.Parse("http://bing.com/search?offset=10&size=1000")
				So(err, ShouldBeNil)
				err = ParseURLParams(u, &params)
				So(err, ShouldNotBeNil)
			})

			Convey("临界值应被允许", func() {
				var params Par
				u, err := url.Parse("http://bing.com/search?offset=0&size=100")
				err = ParseURLParams(u, &params)
				So(err, ShouldBeNil)
			})
		})

		Convey("应支持参数为0的URL", func() {
			var params Params
			u, err := url.Parse("http://bing.com/search?size=0&offset=0")
			err = validateURLParams(u, &params)
			So(err, ShouldBeNil)
		})

		Convey("应支持嵌套的结构体", func() {
			type AdvParam struct {
				Params `json:",squash"`
				Email  string   `json:"email" valid:"email"`
				Types  []string `json:"type" valid:"required"`
			}
			p := AdvParam{}
			u, _ := url.Parse("http://bing.com/search?size=20&offset=10&email=chenjie2@ls.io&type=123")
			err := ParseURLParams(u, &p)
			So(err, ShouldBeNil)
			So(p.Email, ShouldEqual, "chenjie2@ls.io")
			So(p.Offset, ShouldEqual, 10)
			So(p.Size, ShouldEqual, 20)
			So(p.Types, ShouldHaveLength, 1)
			So(p.Types[0], ShouldEqual, "123")

			Convey("嵌套体内错误的范围应该被检测出来", func() {
				u, _ := url.Parse("http://bing.com/search?size=200&offset=10&email=chenjie2@ls.io&type=222")
				err := ParseURLParams(u, &p)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestComplexMap(t *testing.T) {
	Convey("Initial test env...", t, func() {
		type Hero struct {
			Name   string   `json:"name" valid:"required"`
			Titles []string `json:"titles" valid:"required"`
		}

		input := map[string]interface{}{
			"name":   "DK",
			"titles": []string{"Dead Knight", "Human Prince"},
		}

		var h Hero

		Convey("Validation should be succeed", func() {
			err := ParseMap(input, &h)
			So(err, ShouldBeNil)
			So(len(h.Titles), ShouldEqual, 2)
			So(h.Titles[0], ShouldEqual, "Dead Knight")
			So(h.Titles[1], ShouldEqual, "Human Prince")
			So(h.Name, ShouldEqual, "DK")
		})

		Convey("When setting name to a one item array", func() {
			input["name"] = []string{"DK"}

			err := ParseMap(input, &h)
			Convey("validation should be failed", func() {
				So(err, ShouldNotBeNil)
				So(h.Name, ShouldNotEqual, "DK")
			})
		})
	})
}

func TestValidateMap(t *testing.T) {
	Convey("Initial Test Env...", t, func() {
		type Person struct {
			Name string `json:"name" valid:"length(1|2),required"`
			Boy  bool   `json:"boy" valid:"required"`
			Age  int64  `json:"age"`
		}

		input := map[string]interface{}{
			"name": "A3",
			"boy":  true,
			"age":  "123",
		}

		var p Person

		Convey("The validation should be succeed", func() {
			err := ParseMap(input, &p)
			So(err, ShouldBeNil)
		})

		Convey("When setting Name's Length > 2", func() {
			input["name"] = "AXX"
			err := ParseMap(input, &p)

			Convey("The validation should be failed", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When unsetting boy's value", func() {
			delete(input, "boy")
			Convey("The validation should be failed", func() {
				err := ParseMap(input, &p)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When Setting boy's value to valid string", func() {
			input["boy"] = "True"
			err := ParseMap(input, &p)
			Convey("The validation should succeed", func() {
				So(err, ShouldBeNil)
			})

			Convey("The Boy's value should be true", func() {
				So(p.Boy, ShouldBeTrue)
			})
		})

		Convey("When Setting boy's to invalue to string", func() {
			input["boy"] = "hehe"
			err := ParseMap(input, &p)

			Convey("We dont support invalid string to bool", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When Setting boy's value to number", func() {
			input["boy"] = 123123

			err := ParseMap(input, &p)
			Convey("validation should be succeed", func() {
				So(err, ShouldBeNil)
			})

			Convey("should support any number to bool", func() {
				So(p.Boy, ShouldBeTrue)
			})
		})

		Convey("When unsetting a optioanl value", func() {
			delete(input, "age")
			err := ParseMap(input, &p)
			Convey("validation should be succeed", func() {
				So(err, ShouldBeNil)
			})

			Convey("age value should be zero", func() {
				So(p.Age, ShouldBeZeroValue)
			})
			input["age"] = 22
		})

		Convey("When setting age to a valid string value", func() {
			input["age"] = "27"
			err := ParseMap(input, &p)
			Convey("validation should be succeed", func() {
				So(err, ShouldBeNil)
			})

			Convey("age value should be 27", func() {
				So(p.Age, ShouldEqual, 27)
			})
		})

		Convey("When setting age to a invalid string value", func() {
			input["age"] = "27x"
			err := ParseMap(input, &p)
			Convey("validation should be failed", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("age value should be 0", func() {
				So(p.Age, ShouldEqual, 0)
			})
		})
	})
}

func TestEmbedValidateMap(t *testing.T) {
	// 目前mapstruct库不支持嵌套
	Convey("Initial test env...", t, func() {
		type HeroName struct {
			Name string `json:"name"`
		}
		type Hero struct {
			HeroName `json:",squash"`
			Title    string `json:"title"`
		}

		input := map[string]interface{}{
			"name":  "DK",
			"title": "Dead Knight",
		}

		var hero Hero

		config := &mapstructure.DecoderConfig{
			Metadata:         nil,
			Result:           &hero,
			TagName:          "json",
			WeaklyTypedInput: true,
			ZeroFields:       false,
		}

		decoder, _ := mapstructure.NewDecoder(config)
		decoder.Decode(input)
		So(hero.Name, ShouldEqual, "DK")
	})
}

func TestMultipulValue(t *testing.T) {
	Convey("Initial test env...", t, func() {
		type Foo struct {
			Age []int64 `json:"age"`
		}

		input := map[string]interface{}{
			"age": []string{"1", "2"},
		}

		var f Foo

		config := &mapstructure.DecoderConfig{
			Metadata:         nil,
			Result:           &f,
			TagName:          "json",
			WeaklyTypedInput: true,
			ZeroFields:       false,
		}

		decoder, _ := mapstructure.NewDecoder(config)
		decoder.Decode(input)
		So(f.Age, ShouldHaveLength, 2)
	})
}
