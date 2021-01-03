/**
 * @Author Hatch
 * @Date 2021/01/02 13:49
**/
package service

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var (
	funcMap = template.FuncMap{
		"safeUrl": func(s string) interface{} {
			return template.URL(s)
		},
	}

	homeTpl = template.Must(template.New("home.html").Funcs(funcMap).ParseFiles(
		"resource/template/base.html","resource/template/home.html"))
	roomTpl = template.Must(template.New("room.html").Funcs(funcMap).ParseFiles("resource/template/base.html","resource/template/room.html"))
)

type Room struct {
	// 封面
	Avatar		string
	// 房间名称
	Name		string
	// 房间ID
	UUID		uuid.UUID
	// 内容通道
	Broadcast	chan *Message
	// 房间clients
	Clients		map[uuid.UUID]*Client
}

// 派送消息
func (r *Room) DispatchMessage() {
	for {
		select {
		case message := <-r.Broadcast:
			for _, client := range r.Clients {
				client.Send <- message
			}
		}
	}
}

// 初始化房间
func (a *App) initlizeRoomData() {

	for i := 0; i < a.Config.App.InitRoomCount; i++ {
		id := uuid.NewV4()
		nickname := RandomName(4)
		a.Rooms[id] = &Room{
			Avatar: GenAvatar(nickname, 128, "jpeg"),
			UUID: id,
			Name: nickname,
			Clients: make(map[uuid.UUID]*Client),
			Broadcast: make(chan *Message),
		}

		go a.Rooms[id].DispatchMessage()
	}
}

// 服务连接
func (a *App) ServeConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := a.Upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }

    r.ParseForm()
	inputClientId := r.Form.Get("client_id")
	clientId, err := uuid.FromString(inputClientId)
	if  err != nil {
		clientId = uuid.NewV4()
	}

    client := &Client{
    	UUID: clientId,
    	Conn: conn,
    	Send: make(chan *Message),
    	Room: &Room{},
    	CheckHeartBeatTimes: 0,
    	App: a,
	}

	log.Printf("new client:%v connect to server", client.UUID)
	a.Connector <- client

    inputRoomId := r.Form.Get("room_id")
    if roomId, err := a.checkRoomIsExists(inputRoomId); err == nil {
    	client.joinRoom(roomId)
	}
}

// 检查房间是否存在
func (a *App) checkRoomIsExists(inputRoomId string) (roomId uuid.UUID, err error) {

	if strings.Trim(inputRoomId, " ") != "" {
		roomId, err := uuid.FromString(inputRoomId)
		if err != nil {
			return roomId, errors.New(fmt.Sprintf("The room id %v is invalid", inputRoomId))
		}

		if _, ok := a.Rooms[roomId]; !ok {
			return roomId, errors.New(fmt.Sprintf("The room id %v not found", roomId))
		}

		return roomId, nil
	}

	return roomId, errors.New(fmt.Sprintf("The room id %v is empty", inputRoomId))
}

// 检查是否满员
func (r *Room) checkRoomCapacity(capacity int) bool {
	return len(r.Clients) < capacity
}

func (a *App) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inputRoomId := r.Form.Get("room_id")
	if roomId, err := a.checkRoomIsExists(inputRoomId); err != nil {

		err = homeTpl.Execute(w, map[string]interface{}{
			"host" : "ws://" + r.Host + "/ws",
			"rooms": a.Rooms,
		})
		if err != nil {
			log.Println(err)
		}
	} else {
		err = roomTpl.Execute(w, map[string]interface{}{
			"host" : "ws://" + r.Host + "/ws?room_id=" + roomId.String(),
			"room": a.Rooms[roomId],
		})

		if err != nil {
			log.Fatalf("render template error:%v", err)
		}
	}
}