package httpclient

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/chenjie4255/tools/slice"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/chenjie4255/tools/errcode"

	"github.com/chenjie4255/errors"

	"github.com/chenjie4255/tools/log"
)

var logger *log.Logger

func init() {
	logger = log.NewLoggerWithSentry("httpclient")
}

type Client struct {
	Engine Engine
	Logger *log.Logger
}

type Engine interface {
	Do(r *http.Request) (*http.Response, error)
}

//New alloc a new a http clent
func New() *Client {
	tr := &http.Transport{
		// Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 10 * time.Second,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 150,
	}
	engine := &http.Client{Transport: tr, Timeout: 30 * time.Second}
	return &Client{engine, logger}
}

func NewStdClient() *http.Client {
	return newStdClient()
}

func newStdClient() *http.Client {
	tr := &http.Transport{
		// Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 10 * time.Second,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 150,
	}
	return &http.Client{Transport: tr, Timeout: 30 * time.Second}
}

func NewProxyClient() (*http.Client, error) {
	return newDualProxyClient()
}

func newDualProxyClient() (*http.Client, error) {
	dualProxyURL := os.Getenv("DUAL_PROXY_URL")
	if dualProxyURL == "" {
		return nil, errors.New("empty proxy url")
	}

	proxyUrl, err := url.Parse(dualProxyURL)
	if err != nil {
		logger.AddFile().WithFields(log.Fields{
			"proxy_url": dualProxyURL,
			"error":     err,
		}).Error("failed to parse proxy url")
		return nil, err
	}

	proxyTr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 10 * time.Second,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 150,
		Proxy:               http.ProxyURL(proxyUrl),
	}
	return &http.Client{Transport: proxyTr, Timeout: 30 * time.Second}, nil
}

// PostForm post form
func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.Do(req)
}

func NewDualConnClient() *Client {
	dualProxyClient, err := newDualProxyClient()
	if err != nil {
		logger.AddFile().WithField("error", err).Warn("failed to create proxy client, using std http client instead")
		return New()
	}
	stdClient := newStdClient()

	engine := &dualConnEngine{stdClient, dualProxyClient}
	return &Client{engine, logger}
}

// Do send a http request
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	t1 := time.Now()
	resp, err := c.Engine.Do(r)
	useTime := time.Now().Sub(t1)
	if c.Logger != nil {
		if err != nil {
			c.Logger.WithFields(log.Fields{
				"url":            r.URL.String(),
				"error":          err,
				"elapse_time(s)": useTime.Seconds(),
			}).Warn("failed to send http request")
		} else {
			c.Logger.WithFields(log.Fields{
				"url":            r.URL.String(),
				"elapse_time(s)": useTime.Seconds(),
				"resp_code":      resp.StatusCode,
			}).Debug("http request finished")
		}
	}

	return resp, err
}

// DoAndParseJSON 进行一个请求，并且解析其返回值,
// 如果返回code如预期，则结果直接转成JSON从okOutput传出
// 如果返回code非预期，则根据errType的类型，生成相应的ERROR，过程中任何非预期的error,也都从return value中返回
func (c *Client) DoAndParseJSON(r *http.Request, exceptCode int, errFormat, okOutput interface{}) error {
	return c.DoAndParseJSONWithCodes(r, []int{exceptCode}, errFormat, okOutput)
}

func (c *Client) DoAndParseJSONWithCodes(r *http.Request, exceptCodes []int, errFormat, okOutput interface{}) error {
	resp, err := c.Do(r)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return errors.Tag(err, errcode.RemoveServerTimeout)
		}
		return err
	}
	defer resp.Body.Close()

	rd := resp.Body
	// 判断gzip
	if resp.Header.Get("Content-Encoding") == "gzip" {
		rd, err = gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
	}

	match := len(exceptCodes) == 0 || slice.ContainsInt(exceptCodes, resp.StatusCode)

	if !match {
		if errFormat == nil {
			bodyData, _ := ioutil.ReadAll(rd)
			content := fmt.Sprintf("Bad status Code:%d\n%s", resp.StatusCode, string(bodyData))
			return errors.NewWithTag(content, errcode.NotProvisioned)
		}

		ok, err := buildJSONError(rd, errFormat)
		if !ok && c.Logger != nil {
			c.Logger.WithFields(log.Fields{
				"url":   r.URL.String(),
				"error": err,
			}).Error("failed decode respone body")
		}
		return err
	}

	if okOutput != nil {
		return decodeJSON(rd, &okOutput)
	}

	return nil
}

// DoParseJSONData 进行一个请求，并且解析其返回值,
// 如果返回code如预期，则结果直接转成JSON从okOutput传出
// 如果返回code非预期，则根据errType的类型，生成相应的ERROR，过程中任何非预期的error,也都从return value中返回
func (c *Client) DoParseJSONData(r *http.Request, exceptCode int, errFormat, okOutput interface{}) ([]byte, error) {
	return c.DoParseJSONDataWithCodes(r, []int{exceptCode}, errFormat, okOutput)
}

func (c *Client) DoParseJSONDataWithCodes(r *http.Request, exceptCodes []int, errFormat, okOutput interface{}) ([]byte, error) {
	resp, err := c.Do(r)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, errors.Tag(err, errcode.RemoveServerTimeout)
		}
		return nil, err
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 判断gzip
	if resp.Header.Get("Content-Encoding") == "gzip" {
		rd, err := gzip.NewReader(bytes.NewReader(bodyData))
		if err != nil {
			return nil, err
		}

		bodyData, err = ioutil.ReadAll(rd)
		if err != nil {
			return nil, err
		}
	}

	match := len(exceptCodes) == 0 || slice.ContainsInt(exceptCodes, resp.StatusCode)

	if !match {
		if errFormat == nil {
			content := fmt.Sprintf("Bad status Code:%d\n%s", resp.StatusCode, string(bodyData))
			return bodyData, errors.NewWithTag(content, errcode.NotProvisioned)
		}

		ok, err := buildJSONError(bytes.NewReader(bodyData), errFormat)
		if !ok && c.Logger != nil {
			c.Logger.WithFields(log.Fields{
				"url":   r.URL.String(),
				"error": err,
			}).Error("failed decode respone body")
		}
		return bodyData, err
	}

	if okOutput != nil {
		return bodyData, decodeJSON(bytes.NewBuffer(bodyData), &okOutput)
	}

	return bodyData, nil
}

func buildJSONError(data io.Reader, errFormat interface{}) (bool, error) {
	t := reflect.TypeOf(errFormat)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	newVal := reflect.New(t)
	if err := decodeJSON(data, newVal.Interface()); err != nil {
		return false, err
	}

	formatError, ok := newVal.Interface().(error)
	if !ok {
		return false, errors.New("reflect logic fatal")
	}

	return true, formatError
}

type unexpectError struct {
	code int
	body string
}

func (e *unexpectError) toError() error {
	return errors.NewWithTag(fmt.Sprintf("unexpect resp code %d: %s", e.code, e.body), errcode.UnexpectRemoteResponse)
}

// ParseExpectRespJSON 解析一个预期中的JSON，如果code不一致，则直接输出body
func ParseExpectRespJSON(resp *http.Response, exceptCode int, okOutput interface{}) error {
	return ParseExpectsRespJSON(resp, []int{exceptCode}, okOutput)
}

// ParseExpectsRespJSON 解析一个预期中的JSON，如果code不一致，则直接输出body
func ParseExpectsRespJSON(resp *http.Response, exceptCodes []int, okOutput interface{}) error {
	defer resp.Body.Close()

	match := len(exceptCodes) == 0 || slice.ContainsInt(exceptCodes, resp.StatusCode)

	if !match {
		bodyData, _ := ioutil.ReadAll(resp.Body)
		unexpectError := unexpectError{resp.StatusCode, string(bodyData)}
		return unexpectError.toError()
	}

	return decodeJSON(resp.Body, &okOutput)
}

func ParseRespJSON(resp *http.Response, exceptCode int, errFormat, okOutput interface{}) error {
	defer resp.Body.Close()
	if resp.StatusCode != exceptCode {
		if errFormat == nil {
			bodyData, _ := ioutil.ReadAll(resp.Body)
			return errors.NewWithTag(string(bodyData), errcode.NotProvisioned)
		}

		ok, err := buildJSONError(resp.Body, errFormat)
		if !ok {
			logger.WithFields(log.Fields{
				"error": err,
			}).Error("failed decode respone body")
		}
		return err
	}

	return decodeJSON(resp.Body, &okOutput)
}

// GetAndParseJSON 发起GET请求，并解析JSON
func (c *Client) GetAndParseJSON(url string, exceptCode int, errFormat, okOutput interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return c.DoAndParseJSON(req, exceptCode, errFormat, okOutput)
}

func NewPostRequest(url string, payload interface{}) (req *http.Request, err error) {
	if payload != nil {
		jsonData, jsonErr := json.Marshal(payload)
		if err != nil {
			return nil, jsonErr
		}
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(http.MethodPost, url, nil)
	}

	return
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func decodeJSON(reader io.Reader, output interface{}) error {
	bodyData, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bodyData, output); err != nil {
		str := fmt.Sprintf("error:%s\ndata:\n%s", err.Error(), string(bodyData))
		return errors.NewWithTag(str, errcode.DecodeJSONError)
	}

	return nil

}
