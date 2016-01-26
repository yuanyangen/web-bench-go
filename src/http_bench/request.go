package http_bench

import (
	"strconv"
	"strings"
	"time"
)

//对http response的封装
//请求类型，描述一个http请求所需要或者含有的所有的资源
type request struct {
	method    string
	host      string
	port      int
	protocal  string
	path      string
	postField string
	cookieJar string
	timeout   time.Duration
}

//得到一个请求的指针，其指向的对象就是一次http请求的所有的信息
var defaultRequest = request{
	"GET",
	"",
	80,
	"HTTP 1.1",
	"",
	"",
	"",
	10,
}

//设置url的对象属性
func (this *request) SetUrl(url string) {
	//将协议头取出来
	t := strings.Split(url, "://")
	//将协议头转换为大写
	this.protocal = strings.ToUpper(t[0])

	//将host及端口取出来，
	tmp := strings.Split(t[1], "/")
	//判断url中是否有写明端口
	if strings.Contains(tmp[0], ":") {
		tmp2 := strings.Split(tmp[0], ":")
		this.host = tmp2[0]
		port, err := strconv.Atoi(tmp2[1])
		if err != nil {
		}
		this.port = port
	} else {
		this.host = tmp[0]
		this.port = 80
	}

	//取出路径
	this.path = strings.TrimPrefix(t[1], tmp[0])
}

//设置允许post方式进行参数的传递
func (this *request) SetPost() {
	this.method = "POST"
}

//根据reqest对象的所有信息，拼装得到http请求头信息
func (this *request) getHttpStream() (httpStream string) {
	header := ""
	splitTag := "\r\n"

	//拼接http的方法， 这个一定要在最开始
	header = header + this.method + " " + this.path + " " + this.protocal + "/1.1" + splitTag

	//拼接host
	header = header + "host: " + this.host + splitTag

	//设置post的头和正文的长度
	if this.method == "POST" {
		this.postField = strings.TrimRight(this.postField, "&")
		header = header + "Content-Length: " + strconv.Itoa(strings.Count(this.postField, "")-1) + splitTag
		header = header + "Content-Type: application/x-www-form-urlencoded" + splitTag
	}

	//拼接cookie头
	if this.cookieJar != "" {
		header += "Cookie: " + this.cookieJar + splitTag
	}

	header += splitTag
	body := ""
	if this.method == "POST" {
		body += this.postField
	}
	httpStream = header + body

	return
}
