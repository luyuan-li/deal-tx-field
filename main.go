package main

import (
	"github.com/luyuan-li/deal-tx-field/cmd"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 4)
	cmd.Execute()
}
