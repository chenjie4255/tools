// Package validator 用于检测map[string]interface{}是否符合相应的struct
package validator

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/chenjie4255/errors"
	"github.com/gorilla/schema"
	"github.com/mitchellh/mapstructure"

	"github.com/chenjie4255/tools/web"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(false)
}

var errMissingCotentTypeHeader = errors.New("missing Content-Type header")
var errInvalidCotentTypeHeader = errors.New("invalid Content-Type header")

// ParsePostJSON 校验一个http post的JSON数据
func ParsePostJSON(r *http.Request, structure interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return errMissingCotentTypeHeader
	} else if !strings.Contains(contentType, "application/json") {
		return errInvalidCotentTypeHeader
	}

	data, err := web.GetBodyData(r)
	if err != nil {
		return err
	}

	return ParseJSON(data, structure)
}

// ParseJSON 解析一个JSON
func ParseJSON(data []byte, structure interface{}) error {
	err := json.Unmarshal(data, structure)
	if err != nil {
		return err
	}

	ok, err := govalidator.ValidateStruct(structure)
	if !ok {
		return err
	}
	return nil
}

// ParsePostForm 校验一个http postform
func ParsePostForm(r *http.Request, structure interface{}) error {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("json")

	r.ParseForm()
	err := decoder.Decode(structure, r.PostForm)
	if err != nil {
		return err
	}

	ok, err := govalidator.ValidateStruct(structure)
	if !ok {
		return err
	}

	return nil
}

func strArray2Int(strs []string) ([]int64, bool) {
	ret := []int64{}
	for _, s := range strs {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, false
		}
		ret = append(ret, v)
	}
	return ret, true
}

func validateURLParams(url *url.URL, structure interface{}) error {
	values := url.Query()
	input := make(map[string]interface{})

	for k, v := range values {
		if len(v) == 1 {
			input[k] = v[0]
		} else {
			input[k] = v
		}
	}

	return ParseMap(input, structure)
}

// ParseURLParams 解析并校验URL参数, 如果使用匿名嵌套结构体，需要在匿名的结构体中加上
// json:",squash"，
// 注意数组参数不能初始化为空数组，需要为nil
func ParseURLParams(url *url.URL, structure interface{}) error {
	return validateURLParams(url, structure)
}

// ParseMap 校验一个map是否符合输出,注意,如果使用匿名嵌套结构体，需要在匿名的结构体中加上
// json:",squash"
func ParseMap(m map[string]interface{}, structure interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           structure,
		TagName:          "json",
		WeaklyTypedInput: true,
		ZeroFields:       false,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(m)
	if err != nil {
		return err
	}

	ok, err := govalidator.ValidateStruct(structure)
	if !ok {
		return err
	}

	return nil
}

// ValidateObj 检验一个object
func ValidateObj(obj interface{}) error {
	ok, err := govalidator.ValidateStruct(obj)
	if !ok {
		return err
	}
	return nil
}
