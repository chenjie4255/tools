package web

import "net/http"

type RespWrapper struct {
	w      http.ResponseWriter
	r      *http.Request
	maxAge int
}

func NewRespWrapper(w http.ResponseWriter, r *http.Request) *RespWrapper {
	return &RespWrapper{w, r, 0}
}

func NewCacheRespWrapper(w http.ResponseWriter, r *http.Request, maxAge int) *RespWrapper {
	return &RespWrapper{w, r, maxAge}
}

func (r *RespWrapper) Resp(obj interface{}, err error) {
	if err != nil {
		RespError(r.w, r.r, http.StatusBadRequest, err)
		return
	}

	SetCacheControl(r.w, r.maxAge)
	RespJSON(r.w, r.r, http.StatusOK, obj)
}

func (r *RespWrapper) RespL(obj interface{}, err error) {
	if err != nil {
		Resp40XErrorL(r.w, r.r, err)
		return
	}

	SetCacheControl(r.w, r.maxAge)
	RespJSON(r.w, r.r, http.StatusOK, obj)
}

func (r *RespWrapper) RespError(err error) {
	if err != nil {
		RespError(r.w, r.r, http.StatusBadRequest, err)
		return
	}

	RespStatus(r.w, http.StatusNoContent)
}

func (r *RespWrapper) RespErrorL(err error) {
	if err != nil {
		Resp400ErrorL(r.w, r.r, err)
		return
	}

	RespStatus(r.w, http.StatusNoContent)
}

func (r *RespWrapper) RespCount(count int, err error) {
	if err != nil {
		RespError(r.w, r.r, http.StatusBadRequest, err)
		return
	}

	RespJSON(r.w, r.r, http.StatusOK, JSONMap{
		"count": count,
	})
}
