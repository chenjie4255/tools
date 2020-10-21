package signature

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func CalcSignature(key, method, uri, contentMD5, contentType, date string) string {
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, uri, contentMD5, contentType, date)
	// fmt.Println(stringToSign)
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(stringToSign))
	originSig := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(originSig)
}

func SigRequest(key string, r *http.Request, payload []byte) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	contentMD5 := ""

	if len(payload) != 0 {
		hash := md5.New()
		hash.Write(payload)
		contentMD5 = base64.StdEncoding.EncodeToString(hash.Sum(nil))
	}

	ct := r.Header.Get("Content-Type")
	sig := CalcSignature(key, r.Method, r.Host+r.RequestURI, contentMD5, ct, ts)
	if contentMD5 != "" {
		r.Header.Set("Content-MD5", contentMD5)
	}

	r.Header.Set("ML-Authorization", sig)
	r.Header.Set("ML-Timestamp", ts)
}
