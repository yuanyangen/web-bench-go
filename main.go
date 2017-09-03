package main

import (
	"fmt"
	"http_bench"
	"os"
	"runtime"
	"runtime/pprof"
	"flag"
)

//从命令行中获取参数，并将参数写入benchmark的app中
func main() {
	//app.Concurrent = flag.Int("c", 0, "concurrent")
	//app.Number = flag.Uint64("n", 0, "number")
	//flag.Parse()
	//	startCPUProfile()
	//	startMemProfile()
	//	startBlockProfile()
	app := http_bench.GetBenchApp()
	c := flag.Int("c", 1, "concurrent")
	u := flag.String("u", "", "url")
	f := flag.String("f", "", "filePath")
	d := flag.Int("d", 10, "duration")
	flag.Parse()
	app.SetConcurrent(*c)
	app.SetFilePath(*f)
	app.SetUrl(*u)
	app.SetDuration(*d)
	app.Start()
	//	stopCPUProfile()
	//	stockMemProfile()
	//	stopBlockProfile()
}
func startCPUProfile() {
	f, err := os.OpenFile("/tmp/go_pprof_cpu.out", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create cpu profile output file: %s", err)
		return
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Fprintf(os.Stderr, "Can not start cpu profile: %s", err)
		f.Close()
		return
	}
}

func stopCPUProfile() {
	pprof.StopCPUProfile()
}

func startMemProfile() {
	//每多少内存进行一次检查内存的记录， 默认是512k
	runtime.MemProfileRate = 1024 * 64
}
func stockMemProfile() {
	f, err := os.OpenFile("/tmp/go_pprof_mem.out", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create mem profile output file: %s", err)
		return
	}
	if err = pprof.WriteHeapProfile(f); err != nil {
		fmt.Fprintf(os.Stderr, "Can not write %s: ", err)
	}
	f.Close()
}

func startBlockProfile() {
	//如果我们不显式的使用runtime.SetBlockProfileRate函数设置取样间隔，那么取样间隔就为1
	runtime.SetBlockProfileRate(1)
}
func stopBlockProfile() {
	f, err := os.OpenFile("/tmp/go_pprof_block.out", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create mem profile output file: %s", err)
		return
	}
	if err = pprof.Lookup("block").WriteTo(f, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Can not write %s:", err)
	}
	f.Close()
}
