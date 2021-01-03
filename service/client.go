/**
 * @Author Hatch
 * @Date 2021/01/01 15:33
**/
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"log"
	"strings"
	"time"
)

var (
	joinRoomText = "%s进入房间，大家欢迎~"
	leaveRoomText = "%s离开房间"
	roomFullText = "%s房间已满员，请选择其它房间"
)

const (

	TextMessage = iota		// 文本消息
	SysJoinRoomMessage		// 加入房间系统消息
	SysLeaveRoomMessage		// 离开房间系统消息
	HandShakeMessage		// 连接握手
	GlobalMessage		  	// 全局系统消息
	RoomFullMessage			// 房间满员消息
)

// 客户端结构
type Client struct {
	*App
	*Room
	// 头像
	Avatar				string
	// 昵称
	NickName			string
	// 客户端UUID
	UUID				uuid.UUID
	// 客户端连接对象
	Conn				*websocket.Conn
	// 消息信道
	Send				chan *Message
	// 最后通信时间
	LastAckTime			time.Time
	// 心跳检测次数
	CheckHeartBeatTimes	int
}

// 消息结构
type Message struct {
	// 发件人
	Sender		uuid.UUID		`json:"sender"`
	// 内容
	Data		string			`json:"data"`
	// 发送日期
	CreatedAt	string			`json:"created_at"`
	// 动作事件
	ActEvent	int				`json:"act_event"`
	// 房间UUID
	From		uuid.UUID		`json:"from"`
	// IP
	RemoteIp	string			`json:"remote_ip"`
}

// 监听"读取"消息
func (c *Client) ListenPushMessage() {
	log.Printf("listen push message for client:%v", c.UUID)
	newLine, space := []byte("\n"), []byte("")
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.App.Disconnector <- c
			}
			break
		}

		data = bytes.TrimSpace(bytes.Replace(data, newLine, space, -1))

		message := &Message{}

		c.LastAckTime = time.Now()
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		if err = decoder.Decode(message); err != nil {
			log.Println("decode data fail:" + err.Error())
		} else {
			c.LastAckTime = time.Now()
			switch message.ActEvent {
			case TextMessage,GlobalMessage:
				message.Sender = c.UUID
				message.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
				message.From = c.Room.UUID
				repData := map[string]interface{}{
					"content": message.Data,
					"nickname": c.NickName,
					"avatar": c.Avatar,
				}

				jsonRepData, err := json.Marshal(repData)
				if err != nil {
					log.Fatalf("encode message data fail:%v", err)
				}
				message.Data = string(jsonRepData)
				if message.ActEvent == TextMessage {
					c.Room.Broadcast <- message
				} else {
					c.App.Broadcast <- message
				}
			case SysJoinRoomMessage:
				c.NickName = message.Data
				c.Avatar = GenAvatar(message.Data, 128, "jpeg")
				c.App.Clients[c.UUID] = c
				if _, ok := c.App.Rooms[c.Room.UUID]; ok {
					repData := map[string]interface{}{
						"content": fmt.Sprintf(joinRoomText, c.NickName),
						"online_users": len(c.Room.Clients),
					}
					joinRoomMessage := &Message{
						Sender: c.UUID,
						ActEvent: SysJoinRoomMessage,
						CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
						From: c.Room.UUID,
					}

					clients := make([]map[string]interface{}, 0)
					for uuid, item := range c.Room.Clients {
						clients = append(clients, map[string]interface{}{
							"nickname": item.NickName,
							"avatar": template.URL(item.Avatar),
							"uuid": uuid,
						})
					}

					repData["clients"] = clients
					jsonData, err := json.Marshal(repData)
					if err != nil {
						log.Fatalf("encode clients data fail:%v", err)
					}

					joinRoomMessage.Data = string(jsonData)
					c.App.Rooms[c.Room.UUID].Broadcast <- joinRoomMessage
				}
			case HandShakeMessage:
				c.NickName = message.Data
				c.Avatar = GenAvatar(message.Data, 128, "jpeg")
				c.App.Clients[c.UUID] = c
			case SysLeaveRoomMessage:
				c.leaveRoom(c.Room.UUID)
			}
		}
	}
}

// 监听"写入"消息
func (c *Client) ListenPullMessage() {
	log.Printf("listen pull message for client:%v", c.UUID)
	c.LastAckTime = time.Now()
	for {
		select {
		case message := <-c.Send:
			switch message.ActEvent {
			case TextMessage, GlobalMessage, SysJoinRoomMessage, SysLeaveRoomMessage, RoomFullMessage:
				c.writeMessage(message)
			case HandShakeMessage:
				c.writeCustomMessage(&map[string]interface{}{
					"act_event": HandShakeMessage,
					"created_at": message.CreatedAt,
					"client_id": c.UUID,
					"room_id": c.Room.UUID,
				})
			}
		}
	}
}

// 向客户端发送消息通知
func (c *Client) writeMessage(message *Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("encode message error: %v", err)
	}

	if err = c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("write message fail:%v", err)
	}
}

// 向客户端发送自定义格式消息
func (c *Client) writeCustomMessage(message *map[string]interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("errors: %v\n", err)
	}

	if err = c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("write message fail:%v\n", err)
	}
}


// client心跳检测
func (c *Client) HeartBeat() (err error) {
	c.CheckHeartBeatTimes++
	log.Printf("check heartbeat for client:%s nickname: %s times: %v\n", c.UUID, c.NickName, c.CheckHeartBeatTimes)

	// 最后发送/收到消息，小于心跳检测周期，则不检测
	interval := time.Now().Sub(c.LastAckTime)
	if interval < c.App.Config.App.HeartBeatTime * time.Second {
		c.CheckHeartBeatTimes = 0
		return nil
	}

	// 超过最大检测次数
	if c.CheckHeartBeatTimes >= c.App.Config.App.HeartBeatMaxRetry {
		c.App.Disconnector <- c
		return nil
	}

	// 设置写入超时时间
	if err = c.Conn.SetWriteDeadline(time.Now().Add(c.App.Config.App.PingWaitTime * time.Second)); err != nil {
		log.Fatalf("set write deadline fail:%v", err)
	}

	// 发送ping message
	if err = c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		log.Printf("send ping message fail for client:%s", c.UUID)
		c.App.Disconnector <- c
	}

	return err
}

// 离开房间
func (c *Client) leaveRoom(from uuid.UUID) {
	log.Printf("leave room from:%s, client:%s\n", from, c.UUID)

	room := c.App.Rooms[from]
	if room != nil && room.Clients[c.UUID] != nil {
		repData := map[string]interface{}{
			"content": fmt.Sprintf(leaveRoomText, c.NickName),
		}

		leaveRoomMessage := &Message{
			ActEvent: SysLeaveRoomMessage,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			From: from,
		}

		delete(room.Clients, c.UUID)

		clients := make([]map[string]interface{}, 0)
		for uuid, item := range room.Clients {
			clients = append(clients, map[string]interface{}{
				"nickname": item.NickName,
				"avatar": template.URL(item.Avatar),
				"uuid": uuid,
			})
		}

		repData["online_users"] = len(room.Clients)
		repData["clients"] = clients
		jsonData, err := json.Marshal(repData)

		if err != nil {
			log.Fatalf("encode leave room data fail: %v", err)
		}
		leaveRoomMessage.Data = string(jsonData)
		room.Broadcast <- leaveRoomMessage
	}
}

// 进入房间
func (c *Client) joinRoom(from uuid.UUID) {

	room := c.App.Rooms[from]
	c.Room = room
	if room != nil {
		if room.checkRoomCapacity(c.App.Config.App.MaxRoomCapacity) {
			joinRoomMessage := &Message{
				Data: fmt.Sprintf(joinRoomText, c.NickName),
				ActEvent: SysJoinRoomMessage,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
				From: from,
			}
			room.Clients[c.UUID] = c
			room.Broadcast <- joinRoomMessage
		} else {
			roomFullMessage := &Message{
				Data: fmt.Sprintf(roomFullText, room.Name),
				ActEvent: RoomFullMessage,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
				From: from,
			}

			c.Send <- roomFullMessage
		}
	}
}

