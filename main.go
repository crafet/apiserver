// main.go
package main

import (
	. "apiserver/server"

	"log"
)

// set log format
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	var server HTTPFramework

	if err := server.Init(); err != nil {
		log.Println("failed to init apiserver")
		return
	}

	log.Println("api server start...")

	server.Run()

}
