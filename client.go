package main

import (
	"github.com/valyala/fasthttp/fasthttpproxy"
	"github.com/valyala/fasthttp"
	"crypto/tls"
	"time"
	"net"
)

func createClient(channel chan string) *fasthttp.Client {
	return &fasthttp.Client {
		ReadBufferSize: 4 * 4096,
		WriteBufferSize: 4 * 4096,
		MaxConnsPerHost: 1000000000000,
		MaxConnDuration: time.Second * 2,
		MaxIdleConnDuration: time.Second * 2,
		DisablePathNormalizing: true,
		NoDefaultUserAgentHeader: true,
		DisableHeaderNamesNormalizing: true,
		TLSConfig: &tls.Config {
			InsecureSkipVerify: true, 
			MinVersion: tls.VersionTLS13,
			MaxVersion: tls.VersionTLS13,
		},
		Dial: func(addr string) (net.Conn, error) {
			return fasthttpproxy.FasthttpHTTPDialer(<- channel)(addr)
		},
	}
}

func createRequest(method string, useCompression bool) *fasthttp.Request {
	var request *fasthttp.Request = fasthttp.AcquireRequest()

	request.Header.SetMethod(method)

	request.Header.Set("Connection", "keep-alive")
	request.Header.SetUserAgent("Instagram 361.0.0.46.88 (iPhone14,5; iOS 16_0_2; en_GB; en; scale=3.00; 1170x2532; 521442846) AppleWebKit/420+")
	
	if (useCompression) {
		request.Header.Set("Accept-Encoding", "gzip")
	}
	return request
}