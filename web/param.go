package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

var errMissingCotentTypeHeader = errors.New("missing Content-Type header")
var errInvalidCotentTypeHeader = errors.New("invalid Content-Type header")

// swagger:model PageParam
type PageParam struct {
	Offset int `json:"offset" valid:"range(0|999999999),optional"`
	Limit  int `json:"limit" valid:"range(0|1024),optional"`
}

// DefaultPageParam 获取一个默认的分页参数
func DefaultPageParam() PageParam {
	return PageParam{0, 10}
}

// GetURLVar 获取URL中的参数
func GetURLVar(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// GetURLIntVar 获取一个URL中的整形参数
func GetURLIntVar(r *http.Request, key string) (ret int64, err error) {
	val := GetURLVar(r, key)
	ret, err = strconv.ParseInt(val, 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid param format, want int(%s)", val)
	}
	return
}

// Get2URLIntVar 获取二个URL中的整形参数
func Get2URLIntVar(r *http.Request, key1, key2 string) (ret1 int64, ret2 int64, err error) {
	val1 := GetURLVar(r, key1)
	ret1, err = strconv.ParseInt(val1, 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid %s param format, want int(%s)", key1, val1)
	}

	val2 := GetURLVar(r, key2)
	ret2, err = strconv.ParseInt(val2, 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid %s param format, want int(%s)", key2, val2)
	}
	return
}

// GetURLParam 获取指定的URL参数
func GetURLParam(r *http.Request, key string) string {
	vals := r.URL.Query()
	return vals.Get(key)
}

// GetURLIntParams 获取整形的参数集合
func GetURLIntParams(r *http.Request, key string) ([]int64, error) {
	vals := r.URL.Query()[key]
	ret := []int64{}
	for _, v := range vals {
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, val)
	}

	return ret, nil
}

// GetURLIntParam 获取int类型的指定URL参数
func GetURLIntParam(r *http.Request, key string, acceptEmpty bool) (int64, error) {
	val := GetURLParam(r, key)
	if val == "" && acceptEmpty {
		return 0, nil
	}

	ret, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid param format, want int(%s)", val)
	}
	return ret, err
}

// ParseUpdateInfo 解析更新信息
func ParseUpdateInfo(r *http.Request) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	data, err := GetBodyData(r)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&ret)
	return ret, err
}

// ParsePostJSON 校验一个http post的JSON数据
func ParsePostJSON(r *http.Request, structure interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return errMissingCotentTypeHeader
	} else if !strings.Contains(contentType, "application/json") {
		return errInvalidCotentTypeHeader
	}

	data, err := GetBodyData(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, structure)
}
