package http_bench

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type benchmarkApp struct {
	param   *Param
	res     *result
	UrlG    UrlGenerator
	workers []*worker
}
type worker struct {
	param    *Param
	res      *result
	UrlG     UrlGenerator
	StopFlag bool //是否停止？
}
type Param struct {
	Concurrent int
	Number     uint64
	FilePath   string
	TargetUrl  string
	dur        time.Duration
	Wd         bool
}
type result struct {
	total     uint64
	failed    uint64
	bizFailed uint64
}

func (app *benchmarkApp) Start() {
	app.initApp()
	app.Execute()
}

func GetBenchApp() *benchmarkApp {
	app := &benchmarkApp{}
	app.param = &Param{}
	app.workers = make([]*worker, 0)
	app.res = &result{}
	return app
}

func (app *benchmarkApp) SetConcurrent(c int) {
	app.param.Concurrent = c
}

func (app *benchmarkApp) SetDuration(d int) {
	app.param.dur = time.Duration(d)
}

func (app *benchmarkApp) SetUrl(url string) {
	app.param.TargetUrl = url
}

func (app *benchmarkApp) SetFilePath(fp string) {
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
	start := time.Now().UnixNano()
	for i := 0; i < app.param.Concurrent; i++ {
		wk := &worker{}
		wk.param = app.param
		wk.res = &result{}
		wk.StopFlag = false
		wk.UrlG = app.UrlG
		app.workers = append(app.workers, wk)
		go wk.startWorker()
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	select {
	case <-signalChan:
	case <-time.After(app.param.dur * time.Second):
	}
	finish := time.Now().UnixNano()
	dur := (finish - start)/ 1000000
	app.param.dur = time.Duration(dur)
	app.stopApp()
}

func (app *benchmarkApp) stopApp() {
	for _, wk := range app.workers {
		wk.StopFlag = true
		app.res.total += wk.res.total
		app.res.failed += wk.res.failed
		app.res.bizFailed += wk.res.bizFailed
	}
	app.dispalyResult()
}

func (app *benchmarkApp) dispalyResult() {
	fmt.Println("wbg test result:")
	fmt.Printf("qps: %.2f request/Second \n", float64(app.res.total) * 1000/float64(app.param.dur))
	fmt.Printf("total request: %d\n", app.res.total)
	fmt.Printf("total target urls : %d\n", app.res.total)
	fmt.Printf("total time :%d ms \n", app.param.dur)
	fmt.Printf("failed request: : %d\n", app.res.failed)
}

func (wk *worker) startWorker() {
	trans := &http.Transport{}
	for !wk.StopFlag {
		wk.doOneRequest(trans)
	}
}

func (wk *worker) doOneRequest(trans *http.Transport) {
	targetUrl := wk.getTargetUrl(wk.UrlG)
	client := &http.Client{Transport: trans}
	req, _ := http.NewRequest(http.MethodGet, targetUrl, nil)
	resp, err := client.Do(req)

	if err != nil {
		wk.res.failed++
		return
	}
	if resp.StatusCode != 200 {
		wk.res.failed++
		return
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wk.res.failed++
	}
	if wk.param.Wd && !strings.HasPrefix(string(r), "{\"errno\":0") {
		wk.res.bizFailed++
	}
	wk.res.total++
}

func (wk *worker) getTargetUrl(urlG UrlGenerator) string {
	return urlG.GetUrl()
}
