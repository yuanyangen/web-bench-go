package http_bench

import "testing"

func Test_DoOneWorker(t *testing.T) {
	BenchmarkApp.initParam()
	BenchmarkApp.doOneRequest()
}
func Test_Start(t *testing.T) {
	BenchmarkApp.Start()
}
