package structmap

import (
	"testing"
)

func TestConvertJSON2Map(t *testing.T) {

	type CommentTime struct {
		CreatedAt int `json:"created_at"`
		UpdatedAt int `json:"updated_at"`
	}

	type Ts struct {
		Name string `json:"name"`
		Age  int    `json:"-"`
		CommentTime
	}

	testObj := Ts{"Niko", 29, CommentTime{10, 1}}
	result, err := Convert2JSONMap(testObj)
	if err != nil {
		t.Fatalf("error:%s", err.Error())
		return
	}

	t.Logf("result:%v", result)

}
