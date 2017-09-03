package http_bench

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
)

type UrlGenerator interface {
	GetUrl() string
}

type SpecificUrlGenerator struct {
	url string
}
func (u *SpecificUrlGenerator) GetUrl() string{
	return u.url
}

func GetSpecificUrlGenerator(url string) *SpecificUrlGenerator {
	r := &SpecificUrlGenerator{}
	r.url = url
	return r
}

type FileUrlGenerator struct {
	urls []string
	urlSize int
}

func GetNewFileUrlGenerator(filePath string) *FileUrlGenerator {
	fp, err := os.OpenFile(filePath, os.O_RDONLY,0700)
	if err != nil {
		fmt.Println(filePath + " not exist, please check")
		os.Exit(0)
	}
	br := bufio.NewReader(fp)
	urls := make([]string,0)
	for   {
		url, err := br.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("read file error , please check")
			os.Exit(0)
		}
		urls = append(urls, url)
	}
	fg := &FileUrlGenerator{}
	fg.urls = urls
	fg.urlSize = len(urls)
	return fg
}

func (f *FileUrlGenerator) GetUrl() string {
	return f.urls[rand.Int() % f.urlSize]
}
