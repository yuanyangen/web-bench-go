package main

import (
	"flag"
	"fmt"
	"http_bench"
)

//从命令行中获取参数，并将参数写入benchmark的app中
func main() {
	app := http_bench.BenchmarkApp
	app.Concurrent = flag.Int("c", 0, "concurrent")
	app.Cookie = flag.String("cookie", "", "cookie")
	app.PostField = flag.String("post", "", "post info")
	url := flag.String("url", "", "target url")
	flag.Parse()
	if "" == *url {
		fmt.Println("usage:\n")
	}

	app.Bench(*url)
	fmt.Println(*app.Result)

}

//对http请求的封装，能够模拟get请求， post请求， 带上cookie

//对transport的封装，在底层能够获取信息，

//能够支持长连接，
//

//
