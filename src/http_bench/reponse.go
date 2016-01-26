package http_bench

import "strings"

//对http request的封装
type response struct {
	//status int32
	header string
	body   string
	/* setCookie        string*/
	//contentType      string
	//date             string
	//transferEncoding string
	//connection       string
	//server           string
	/*cacheControl     string*/
}

var defaultResp = response{
	"",
	"",
}

//解析http的响应，这里只是将http的头和body分开
func (this *response) processResponse(res string) {
	tmpResponse := strings.Split(res, "\r\n\r\n")
	//解析http response头
	this.header = tmpResponse[0]

	/*
		//如果返回的transfer-encoding是chunked，对响应的body进行解析
		if strings.Contains(this.header, "Transfer-Encoding: chunked") {
			//这里临时只取一个
			this.parserChunkedBody(tmpResponse[1])
		} else {
			this.body = tmpResponse[1]
		}
	*/
}

//这里的临时解决办法是只考虑只有一个分片的时候， 如果拥有多个分片， 目前的做法会有问题
func (this *response) parserChunkedBody(body string) {
	tmpBody := strings.Split(body, "\r\n")
	this.body = tmpBody[1]
}
