package proxy

import (
	"GateWayDemoStudent/demo/proxy/middleware"
	"GateWayDemoStudent/proxy/load_balance"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func NewLoadBalanceReverseProxy(ctx *middleware.SliceRouterContext, lb load_balance.LoadBalance) *httputil.ReverseProxy {
	//请求协调者
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.URL.String())
		if err != nil {
			log.Fatal("get next addr fail")
		}

		target, err := url.Parse(nextAddr)
		if err != nil {
			log.Fatal(err)
		}

		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	//更改内容
	modifyFunc := func(resp *http.Response) error {
		//兼容websocket
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}

		var (
			payload []byte
			readErr error
		)

		//添加zip压缩
		if strings.Contains(resp.Header.Get("Content-Encodeing"), "gzip") {
			gr, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}

			payload, readErr = ioutil.ReadAll(gr)
			resp.Header.Del("Content-Encoding")
		} else {
			payload, readErr = ioutil.ReadAll(resp.Body)
		}

		if readErr != nil {
			return readErr
		}

		//因为预读了数据, 所有内容重新回写
		ctx.Set("status_code", resp.StatusCode)
		ctx.Set("payload", payload)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		resp.ContentLength = int64(len(payload))
		resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))

		return nil
	}

	//错误回调: 关闭real_server 时测试, 错误回调
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		//todo record error log

		fmt.Println(err)
	}

	return &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc,
	}
}
