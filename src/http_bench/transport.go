package http_bench

import (
	"fmt"
	"net"
	"strconv"
)

//对transport的封装，在底层能够获取信息，

//能够支持长连接，
//

//实现curl 功能需要有3个部分：
// 1: Request, 这个是一个结构体，其中包含很多的属性
// 2: connection， 需要根据Request中的host属性及端口号建立连接
// 3: response， 解析返回的响应头，将常用的属性都写入到其中的结构体
//这里含有一个conn，说明对于不同的请求，使用的是短链接
type Ch struct {
	req  request
	resp response
	conn connection
}

type connection struct {
	host string
	port int32
}

var defaultConn = connection{
	"",
	80,
}

//返回的是一个Ch的对象，这个是因为
func Init() Ch {
	return Ch{defaultRequest, defaultResp, defaultConn}
}

func (this *Ch) GetBody() (body string) {
	body = this.resp.body
	return
}
func (this *Ch) GetHeader() string {
	return this.resp.header
}

//根据reqest的中的信息，得到对应的tcp连接
func getConn() (conn *net.TCPConn) {
	host := defaultRequest.host + ":" + strconv.Itoa(defaultRequest.port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		fmt.Println("failed")
	}

	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	return
}
