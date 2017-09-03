package http_bench

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"os"
	"os/signal"
)

//对http请求的封装，能够模拟get请求， post请求， 带上cookie
type benchmarkApp struct {
	param     *Param
	res       *result
	wg        *sync.WaitGroup
	msgChan   chan int
	UrlG      UrlGenerator
	transport *http.Transport
}
type Param struct {
	Concurrent int
	Number     uint64
	FilePath   string
	TargetUrl   string
	Dura       time.Duration
	Wd         bool
}
type result struct {
	total     uint64
	failed    uint64
	bizFailed uint64
}


const TYPE_TOTAL = 0
const TYPE_FAILED = 1
const TYPE_WD_BIZ_FAILED = 2

func (app *benchmarkApp) Start() {
	app.initApp()
	app.Execute()
}

func GetBenchApp() *benchmarkApp {
	app := &benchmarkApp{}
	app.msgChan = make(chan int, 1024)
	app.param = &Param{}
	app.transport = &http.Transport{}
	return app
}

func (app * benchmarkApp) SetConcurrent(c int) {
	app.param.Concurrent = c
}

func (app * benchmarkApp) SetDuration(d int) {
	app.param.Dura = time.Duration(d)
}

func (app * benchmarkApp) SetUrl(url string) {
	app.param.TargetUrl = url
}

func (app * benchmarkApp) SetFilePath(fp string) {
	app.param.FilePath = fp
}

func (app *benchmarkApp) initApp() {
	if app.param.FilePath != "" {
		app.UrlG = GetNewFileUrlGenerator(app.param.FilePath)
	} else if app.param.TargetUrl != "" {
		app.UrlG = GetSpecificUrlGenerator(app.param.TargetUrl)
	} else {
		displayErrors("url not specified\n")
	}

}

//运行一个curl的对象，并将返回信息返回
func (app *benchmarkApp) Execute() {
	go app.receiveResult()
	for i := 0; i < app.param.Concurrent; i++ {
		go app.doOneWorker()
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	select {
	case <- signalChan: {
	}
	case <- time.After(app.param.Dura * time.Second) : {
	}
	}
	app.dispalyResult()
}

func (app *benchmarkApp) dispalyResult() {
	fmt.Printf("qps: %d\n", app.res.total/uint64(app.param.Dura))
	fmt.Printf("total request: %d\n", app.res.total)
	fmt.Printf("failed request: : %d\n", app.res.failed)
	fmt.Printf("total traget urls : %d\n", app.res.total/uint64(app.param.Dura))
}

func (app *benchmarkApp) receiveResult() {
	app.res = &result{}
	for {
		msg := <-app.msgChan
		if msg == TYPE_TOTAL {
			app.res.total++
		} else if msg == TYPE_WD_BIZ_FAILED {
			app.res.bizFailed++
		} else if msg == TYPE_FAILED {
			app.res.failed++
		}
	}
}

func (app *benchmarkApp) doOneWorker() {
	//trans := app.transport
	trans := &http.Transport{}
	for {
		app.doOneRequest(trans)
	}
}

func (app *benchmarkApp) doOneRequest(trans *http.Transport) {
	targetUrl := app.getTargetUrl(app.UrlG)
	client := &http.Client{Transport: trans}
	req, _ := http.NewRequest(http.MethodGet, targetUrl, nil)
	resp, err := client.Do(req)

	if err != nil {
		app.msgChan <- TYPE_FAILED
		return
	}
	if resp.StatusCode != 200 {
		app.msgChan <- TYPE_FAILED
		return
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.msgChan <- TYPE_FAILED
	}
	if app.param.Wd && !strings.HasPrefix(string(r), "{\"errno\":0") {
		app.msgChan <- TYPE_WD_BIZ_FAILED
	}
	app.msgChan <- TYPE_TOTAL
}

func (app *benchmarkApp) getTargetUrl(urlG UrlGenerator) string {
	return urlG.GetUrl()
}
