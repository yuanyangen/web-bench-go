package http_bench

import (
	"fmt"
	"io/ioutil"
	"time"
)

//对http请求的封装，能够模拟get请求， post请求， 带上cookie
type benchmarkApp struct {
	Concurrent *int
	Cookie     *string
	PostField  *string
	Result     *string
}

var BenchmarkApp = &benchmarkApp{}

func (this *benchmarkApp) Bench(url string) {
	defaultRequest.SetUrl(url)
	if "" != *this.PostField {
		defaultRequest.postField = *this.PostField
	}
	if "" != *this.Cookie {
		defaultRequest.cookieJar = *this.Cookie
	}
	fmt.Println("11111111111111")

}

//运行一个curl的对象，并将返回信息返回
func (this *benchmarkApp) Execute() {
	//设置一个channel作为超时使用
	channel := make(chan string, 1)
	httpRequest := defaultRequest.getHttpStream()
	go func(httpRequest string) {
		conn := getConn()
		_, err := conn.Write([]byte(httpRequest))
		if err != nil {
		}
		tmpInfo, _ := ioutil.ReadAll(conn)
		//这里应该将输出的字符串整理成response对象
		defaultResp.processResponse(string(tmpInfo))
		if defaultResp.header == "200" {
			fmt.Println(defaultResp.header)
		}
		channel <- defaultResp.header
	}(httpRequest)
	//如果超时, 返回的响应头中的header为timeout
	select {
	case <-channel:
		{
			return
		}
	case <-time.After(time.Second * defaultRequest.timeout):
		{
			defaultResp.header = "timeout"
		}
	}
	return
}
