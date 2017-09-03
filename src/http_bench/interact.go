package http_bench

import (
	"fmt"
	"os"
)

func displayErrors(msg string) {
    fmt.Println(msg)
	os.Exit(1)
}
