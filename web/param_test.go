package web

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseUpdateInfo(t *testing.T) {
	req := httptest.NewRequest("POST", "/123", bytes.NewReader([]byte(`{"age":10086}`)))
	req.Header.Set("Content-Type", "application/json")
	info, err := ParseUpdateInfo(req)
	if err != nil {
		t.Fatalf("failed to parse update info, %s", err)
	}
	val, ok := info["age"]
	if !ok {
		t.Fatal("parse result error, age not found")
	}

	v := reflect.ValueOf(val)
	fmt.Println(v.Kind().String())

	valInt, ok := val.(float64)
	if !ok {
		t.Fatal("age should be a float64 type")
	}

	if valInt != 10086 {
		t.Fatalf("got %f, expect:10086", valInt)
	}

}
