package web

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
)

type contextKey int

const (
	keyBodyData contextKey = iota
)

// GetBodyData 获取http request中的body参数（注：由于http request中提供的body是reader，因此多次读取时需要一处缓存）
func GetBodyData(r *http.Request) ([]byte, error) {
	ctxBodyData := r.Context().Value(keyBodyData)
	if ctxBodyData != nil {
		bodyData, ok := ctxBodyData.([]byte)
		if !ok {
			panic("the type of keyBodyData should be []byte")
		}
		return bodyData, nil
	}

	return readBodyData(r)
}

func readBodyData(r *http.Request) ([]byte, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	return data, nil
}

const (
	mediaTypeJSON = "application/json"
)

// BodyReaderMiddleware 读取body中内容，写入context
func BodyReaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		mt, _, err := mime.ParseMediaType(ct)
		if mt != mediaTypeJSON && !strings.Contains(ct, mediaTypeJSON) {
			next.ServeHTTP(w, r)
			return
		}
		data, err := readBodyData(r)
		if err == nil {
			ctx := context.WithValue(r.Context(), keyBodyData, data)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AllBodyReaderMiddleware读取body中内容，写入context
func AllBodyReaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct == "" {
			next.ServeHTTP(w, r)
			return
		}

		data, err := readBodyData(r)
		if err == nil {
			ctx := context.WithValue(r.Context(), keyBodyData, data)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}
