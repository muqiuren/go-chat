/**
 * @Author Hatch
 * @Date 2021/01/01 11:09
**/
package main

import (
	"flag"
	"go-chat/service"
	"net/http"

	"log"
)

func main() {
	Run()
}

func Run() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("catch panic: %s\n", err)
		}
	}()
	app := service.NewApp()
	go app.Run()
	// 解析命令行参数
	addr := flag.String("addr", app.Config.App.Port, "http server host")
	flag.Parse()

	log.Printf("connect to server: http://localhost:%s", *addr)
	log.Fatal(http.ListenAndServe(":" + *addr, app.Router()))
}