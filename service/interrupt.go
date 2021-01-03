/**
 * @Author Hatch
 * @Date 2021/01/02 14:10
**/
package service

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// 监听关闭信号
func (a *App) ShutdownTimer() {
	interrupt := make(chan os.Signal)
	// kill -l
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
		case <-interrupt:
			execQuit(a)
	}
}

// 中断处理函数
func execQuit(a *App) {
	log.Println("start interrupt...")
	for _, client := range a.Clients {
		client.Conn.Close()
	}
	log.Println("complete interrupt")
	os.Exit(0)
}
