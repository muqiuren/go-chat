/**
 * @Author Hatch
 * @Date 2021/01/01 20:31
**/
package service

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var router *mux.Router

func (a *App) Router() *mux.Router {
	if router != nil {
		return router
	}
	router = mux.NewRouter()
	// 选择房间
	router.HandleFunc("/", a.ChooseRoom).Name("chooseRoom")
	// 建立websocket连接
	router.HandleFunc("/ws", a.ServeConnection).Name("serveConnection")
	// 静态资源
	router.PathPrefix("/assert/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := a.Config.App.StaticFolder + "/" + r.URL.Path[1:]

		if !FileExist(p) {
			http.NotFound(w, r)
		} else if !IsAllowExtFile(p, a.Config.App.AllowExts) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		} else {
			http.ServeFile(w, r, p)
		}
	})
	log.Println("router initialize complete")

	return router
}
