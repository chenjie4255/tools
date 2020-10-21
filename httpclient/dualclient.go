package httpclient

import (
	"net/http"
	"time"

	"github.com/chenjie4255/tools/errcode"
	"github.com/chenjie4255/tools/gor"
	"github.com/chenjie4255/tools/log"
	"github.com/chenjie4255/errors"
)

// DualConnEngine 双通道HTTP CLIENT, 一个请求将同时往2个通道发送，任何一个连接成功，则忽略另一个
// 如果2者都失败，则返回值以默认的为准
type dualConnEngine struct {
	defClient   *http.Client
	proxyClient *http.Client
}

type requestResp struct {
	resp *http.Response
	err  error
}

// Do send a http request
func (c *dualConnEngine) Do(r *http.Request) (*http.Response, error) {
	t1 := time.Now()

	defChan := make(chan *requestResp, 1)
	proxyChan := make(chan *requestResp, 1)

	gor.RunWithRecover(func() {
		resp, err := c.defClient.Do(r)
		useTime := time.Now().Sub(t1)
		if err != nil {
			logger.WithFields(log.Fields{
				"url":            r.URL.String(),
				"error":          err,
				"elapse_time(s)": useTime.Seconds(),
			}).Debug("default client failed to send http request")
		} else {
			logger.WithFields(log.Fields{
				"url":            r.URL.String(),
				"elapse_time(s)": useTime.Seconds(),
				"resp_code":      resp.StatusCode,
			}).Debug("default client http request finished")
		}

		defChan <- &requestResp{resp, err}
	})

	gor.RunWithRecover(func() {
		if c.proxyClient == nil {
			proxyChan <- nil
			return
		}
		resp, err := c.proxyClient.Do(r)
		useTime := time.Now().Sub(t1)
		if err != nil {
			logger.WithFields(log.Fields{
				"url":            r.URL.String(),
				"error":          err,
				"elapse_time(s)": useTime.Seconds(),
			}).Debug("proxy client failed to send http request")
		} else {
			logger.WithFields(log.Fields{
				"url":            r.URL.String(),
				"elapse_time(s)": useTime.Seconds(),
				"resp_code":      resp.StatusCode,
			}).Debug("proxy client http request finished")
		}

		proxyChan <- &requestResp{resp, err}
	})

	defResp := &requestResp{}

	retChan := make(chan *requestResp, 1)

	gor.RunWithRecover(func() {
		got := false
		// 2个请求一定要读完，并且废弃的那个请求，需要关闭Body
		for i := 0; i < 2; i++ {
			select {
			case defResp = <-defChan:
				if defResp.err == nil && defResp.resp != nil {
					if got {
						defResp.resp.Body.Close()
						continue
					}

					logger.WithFields(log.Fields{
						"url": r.URL.String(),
					}).Debug("request using default request's responese")
					got = true
					retChan <- defResp
				}
			case proxyResp := <-proxyChan:
				if proxyResp == nil {
					continue
				}

				if proxyResp.err == nil && proxyResp.resp != nil {
					if got {
						proxyResp.resp.Body.Close()
						continue
					}

					logger.WithFields(log.Fields{
						"url": r.URL.String(),
					}).Debug("request using proxy request's responese")
					got = true
					retChan <- proxyResp
				}
			case <-time.After(60 * time.Second):
				// 主要防止在以上2个出现超时情况下，这里永远在select block
				err := errors.NewWithTag("request timeout", errcode.RemoveServerTimeout)
				logger.AddFile().WithFields(log.Fields{
					"url":   r.URL.String(),
					"error": err,
				}).Error("dual conn request timeout")
				got = true
				retChan <- &requestResp{nil, err}
			}
		}
	})

	ret := <-retChan
	return ret.resp, ret.err // 函数内保证ret不为nil
}
