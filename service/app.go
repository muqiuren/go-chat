/**
 * @Author Hatch
 * @Date 2021/01/01 11:09
**/
package service

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

type App struct {
	// 应用配置
	Config			*Config
	// 连接器
	Connector		chan *Client
	// 断开连接器
	Disconnector	chan *Client
	// 客户端集合
	Clients			map[uuid.UUID]*Client
	// 房间集合
	Rooms			map[uuid.UUID]*Room
	// 消息广播
	Broadcast		chan *Message
	// 连接升级器
	Upgrader		*websocket.Upgrader
}

// 应用实例
func NewApp() *App {
	return &App {
		Config:       InitConfig(),
		Connector:    make(chan *Client),
		Disconnector: make(chan *Client),
		Clients:      make(map[uuid.UUID]*Client),
		Rooms:		  make(map[uuid.UUID]*Room),
		Broadcast:    make(chan *Message),
		Upgrader:		&websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
		},
	}
}

// 运行消息监听
func (a *App) Run() {
	go a.ShutdownTimer()
	go a.HeartBeat()
	a.initlizeRoomData()

	for {
		select {
		case client := <-a.Connector:	// 连接客户端
			log.Printf("send handshake message for client:%v", client.UUID)
			go client.ListenPushMessage()
			go client.ListenPullMessage()
			// 发送handshake
			message := &Message{
				ActEvent: HandShakeMessage,
				RemoteIp: client.Conn.RemoteAddr().String(),
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
				From: client.Room.UUID,
			}
			client.Send <- message
		case client := <-a.Disconnector: // 客户端断连
			log.Printf("disconnect client:%v", client.UUID)
			client.leaveRoom(client.Room.UUID)
			err := client.Conn.Close()
			if err != nil {
				log.Printf("disconnect client fail:%v", err)
			}
			delete(a.Clients, client.UUID)
		case message := <-a.Broadcast:	// 全局消息广播
			if message.ActEvent == GlobalMessage {
				for _, client := range a.Clients {
					client.Send <- message
				}
			}
		}
	}
}

// 心跳检测
func (a *App) HeartBeat() {

	ticker := time.NewTicker(a.Config.App.HeartBeatTime * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, client := range a.Clients {
				client.HeartBeat()
			}
		}
	}
}
