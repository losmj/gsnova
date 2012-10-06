package proxy

import (
	"bytes"
	"common"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

func dummyReq(method string) *http.Request {
	return &http.Request{Method: method}
}

func indexHandler(req *http.Request) *http.Response {
	res := &http.Response{Status: "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    dummyReq("GET"),
		Header: http.Header{
			"Connection":   {"close"},
			"Content-Type": {"text/html"},
		},
		Close:         true,
		ContentLength: -1}
	hf := common.Home + "/web/html/index.html"
	if content, err := ioutil.ReadFile(hf); nil == err {
		strcontent := string(content)
		strcontent = strings.Replace(strcontent, "${Version}", common.Version, -1)
		strcontent = strings.Replace(strcontent, "${ProxyPort}", common.ProxyPort, -1)
		var buf bytes.Buffer
		buf.WriteString(strcontent)
		res.Body = ioutil.NopCloser(&buf)
	}
	return res
}

func pacHandler(req *http.Request) *http.Response {
	res := &http.Response{Status: "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    dummyReq("GET"),
		Header: http.Header{
			"Connection":          {"close"},
			"Content-Type":        {"application/x-ns-proxy-autoconfig"},
			"Content-Disposition": {"attachment;filename=snova-gfwlist.pac"},
		},
		Close:         true,
		ContentLength: -1}
	hf := common.Home + "/snova-gfwlist.pac"
	if content, err := ioutil.ReadFile(hf); nil == err {
		var buf bytes.Buffer
		buf.Write(content)
		res.Body = ioutil.NopCloser(&buf)
	}
	return res
}

func handleSelfHttpRequest(req *http.Request, conn net.Conn) {
	path := req.URL.Path
	log.Printf("Path is %s\n", path)
	var res *http.Response
	switch path {
	case "/pac/gfwlist":
		res = pacHandler(req)
	case "/":
		res = indexHandler(req)
	}
	if nil != res {
		res.Write(conn)
	}
	conn.Close()
}