package main

import (
	initModule "github.com/sheginabo/go-quick-gin/init"
)

func main() {
	initProcess := initModule.NewMainInitProcess("./")
	initProcess.Run()
}
