package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/chenjie4255/errors"

	"github.com/chenjie4255/tools/errcode"
	"github.com/chenjie4255/tools/i18n"
	"github.com/chenjie4255/tools/log"
	"github.com/chenjie4255/tools/structfix"
)

// RespOKJSON 返回 StatusOK JSON
func RespOKJSON(w http.ResponseWriter, r *http.Request, o interface{}) {
	RespJSON(w, r, http.StatusOK, o)
}

// RespOKJSONL 返回 StatusOK JSON
func RespOKJSONL(w http.ResponseWriter, r *http.Request, o interface{}) {
	RespJSONL(w, r, http.StatusOK, o)
}

// RespJSON 返回JSON
func RespJSON(w http.ResponseWriter, r *http.Request, status int, o interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if r.Context().Value("trLang") != nil {
		lang := r.Context().Value("trLang").(string)
		if lang != "" {
			i18n.SetLangForObject(o, lang)
		}
	}

	structfix.FixNilArray(o, true)

	if err := encoder.Encode(o); err != nil {
		log.Default().Info(err) // 这里的错误一般是broken pipe。。。
	}
}

// RespJSON 返回JSON
func RespJSONL(w http.ResponseWriter, r *http.Request, status int, o interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if r.Context().Value("trLang") != nil {
		lang := r.Context().Value("trLang").(string)
		if lang != "" {
			i18n.SetLangForObject(o, lang)
		}
	}

	structfix.FixNilArray(o, true)

	if err := encoder.Encode(o); err != nil {
		log.Default().Info(err) // 这里的错误一般是broken pipe。。。
	}

	objectJSON, _ := json.Marshal(o)

	requestBody, _ := GetBodyData(r)
	log.Default().WithFields(log.Fields{
		"request_uri":    r.RequestURI,
		"request_body":   string(requestBody),
		"request_header": fmt.Sprintf("%v", r.Header),
		"status":         status,
		"resp":           string(objectJSON),
	}).Info("RespJSONL")
}

// RespParamError 参数错误MSG
func RespParamError(w http.ResponseWriter, r *http.Request, msg string) {
	RespJSON(w, r, http.StatusBadRequest, newRespErr(errcode.ParamError, msg))
}

func getErrorCode(err error) int {
	tag := errors.GetTag(err)
	if tag == 0 {
		return errcode.Undefined
	}
	return tag
}

// RespError 返回一个错误
func RespError(w http.ResponseWriter, r *http.Request, code int, err error) {
	RespJSON(w, r, code, newErrorJSON(err))
}

// Resp400Error 返回一个400错误
func Resp400Error(w http.ResponseWriter, r *http.Request, err error) {
	RespJSON(w, r, http.StatusBadRequest, newErrorJSON(err))
}

func Resp400ErrorL(w http.ResponseWriter, r *http.Request, err error) {
	RespJSON(w, r, http.StatusBadRequest, newErrorJSON(err))

	requestBody, _ := GetBodyData(r)
	log.Default().WithFields(log.Fields{
		"request_uri":    r.RequestURI,
		"request_body":   string(requestBody),
		"request_header": fmt.Sprintf("%v", r.Header),
		"error":          err.Error(),
	}).Warn("web resp 400 error")
}

// Resp40XErrorL if err.code == errcode.Unauthorized, will return http 401 status code instead 400
func Resp40XErrorL(w http.ResponseWriter, r *http.Request, err error) {
	tag := errors.GetTag(err)
	httpCode := http.StatusBadRequest
	if tag == errcode.Unauthorized {
		httpCode = http.StatusUnauthorized
	}
	RespJSON(w, r, httpCode, newErrorJSON(err))

	requestBody, _ := GetBodyData(r)
	log.Default().WithFields(log.Fields{
		"request_uri":    r.RequestURI,
		"request_body":   string(requestBody),
		"request_header": fmt.Sprintf("%v", r.Header),
		"error":          err.Error(),
	}).Warn("web resp 40x error")
}

func RespXML(w http.ResponseWriter, r *http.Request, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/xml; charset=UTF-8")
	w.WriteHeader(status)
	encoder := xml.NewEncoder(w)

	if err := encoder.Encode(body); err != nil {
		log.Default().Info(err) // 这里的错误一般是broken pipe。。。
	}
}

type partialError struct {
	JSONMap
	Succeed interface{} `json:"succeed"`
}

// Resp500Error 返回一个400错误
func Resp500Error(w http.ResponseWriter, r *http.Request, err error) {
	RespJSON(w, r, http.StatusInternalServerError, newErrorJSON(err))
}

type JSONMap map[string]interface{}

func newErrorJSON(err error) JSONMap {
	if err == nil {
		return newRespErr(errcode.Undefined, "no err")
	}
	return newRespErr(getErrorCode(err), err.Error())
}

// ErrorResp 服务端标准的错误输出
type ErrorResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e ErrorResp) Error() string {
	return fmt.Sprintf("code: %d, %s", e.Code, e.Msg)
}

func newRespErr(code int, msg string) JSONMap {
	return JSONMap{
		"code": code,
		"msg":  msg}
}

// SetCacheControl 设置请求的缓存
func SetCacheControl(w http.ResponseWriter, maxAge int) {
	if maxAge > 0 {
		val := fmt.Sprintf("public, max-age=%d, must-revalidate", maxAge)
		w.Header().Set("Cache-Control", val)
	}
}

func RespStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func RespNoContent(w http.ResponseWriter) {
	RespStatus(w, http.StatusNoContent)
}

func DefaultOptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
