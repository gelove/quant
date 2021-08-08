package request

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"quant/pkg/utils/errs"
	"quant/pkg/utils/json"
	"quant/pkg/utils/logger"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req"
	"github.com/pkg/errors"
)

type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
)

var methods = map[Method]string{
	GET:    "GET",
	POST:   "POST",
	PUT:    "PUT",
	DELETE: "DELETE",
}

// func printSlice(s []interface{}) {
// 	for _, v := range s {
// 		logger.Infof("printSlice => %#v", v)
// 	}
// }

// New NewRequest
func New(method Method, url string, res interface{}, header req.Header, vs ...interface{}) error {
	// req.Debug = true
	// logger.Info("NewRequest =>", method, url, header)
	// printSlice(vs)
	vs = append(vs, header)
	cli := req.New()
	cli.SetTimeout(60 * time.Second)
	resp, err := cli.Do(methods[method], url, vs...)
	if err != nil {
		return errors.WithStack(err)
	}
	statusCode := resp.Response().StatusCode

	if statusCode == http.StatusUnauthorized {
		return errors.WithStack(errs.Get(http.StatusUnauthorized))
	}
	if statusCode < 200 || statusCode >= 300 {
		return errors.WithStack(errors.New(http.StatusText(statusCode)))
	}
	data, err := resp.ToBytes()
	if err != nil {
		return errors.WithStack(err)
	}
	// logger.Debugf("NewRequest data => %s", string(data))
	if len(data) == 0 {
		return nil
	}
	err = json.Unmarshal(data, res)
	return errors.WithStack(err)
}

// Random 随机数字符串
func Random() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func Encode(values map[string]string) string {
	if values == nil {
		return ""
	}
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	list := make([]string, 0, len(values))
	for _, k := range keys {
		list = append(list, k+"="+values[k])
	}
	return strings.Join(list, "&")
}

// ApiRequest 网络请求
func ApiRequest(uri string, method string, data map[string]interface{}, headers map[string]string, resp interface{}) error {
	logger.Info("ApiRequest =>", uri, method, data, headers)
	var body string
	var header = make(http.Header)
	for k, v := range headers {
		header.Add(k, v)
	}
	val, ok := headers["Content-Type"]
	if !ok || val == "" {
		// 默认为JSON传参
		header.Set("Content-Type", "application/json")
	}
	if val == "" || val == "application/json" {
		body = json.MarshalToString(data)
		return Request(uri, method, body, header, resp)
	}
	//header.Add("Content-Type", "application/x-www-form-urlencoded")
	body = json.MustToString(data)
	return Request(uri, method, body, header, resp)
}

// Request 发起网络请求
func Request(uri string, method string, data string, header http.Header, resp interface{}) error {
	//client := &http.Client{Timeout: time.Second * 10}
	client := &http.Client{}
	request, err := http.NewRequest(method, uri, strings.NewReader(data))
	if err != nil {
		return err
	}

	request.Header = header

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()
	// 30秒超时
	request = request.WithContext(ctx)

	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			logger.Errorf("Request => %+v", errors.WithStack(err))
		}
	}()
	str, err := ioutil.ReadAll(res.Body)
	logger.Debugf("res.Body ==> %s", string(str))
	return json.Unmarshal(str, resp)
}

// ConvertUrlValue 转为url.Values Encode会根据Ascii码排序
func ConvertUrlValue(m req.Param) url.Values {
	vs := make(url.Values)
	for k, v := range m {
		vs.Add(k, fmt.Sprint(v))
	}
	return vs
}
