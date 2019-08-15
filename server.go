package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

//客户端管理
type ClientManager struct {
	//客户端 map 储存并管理所有的长连接client，在线的为true，不在的为false
	clients map[int]*Client
	//新创建的长连接client
	register chan *Client
	//新注销的长连接client
	unregister chan int
	// 已经注册人数
	index	int
}

//客户端 Client
type Client struct {
	//用户id
	id string
	//连接的socket
	socket *websocket.Conn
}

// 对战信息
type Result struct{
	// 落子坐标
	xy			string	`json:"xy,omitempty"`
	// 发送消息
	message		string	`json:"message,omitempty"`
	// 允许下棋
	bout		bool	`json:"bout,omitempty"`
	// color
	color		string	`json:"color,omitempty"`
}

//创建客户端管理者
var manager = ClientManager{
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[int]*Client),
	index:		0,
}

func (manager *ClientManager) start() {
	for {
		select {
		//如果有新的连接接入,就通过channel把连接传递给conn
		case conn := <-manager.register:
			manager.index++
			// 判断有无对手
			if (manager.index%2 == 0){
				// 或许有对手了，确认对手存在
				if(manager.clients[manager.index-1] != nil){
					// 把自己载入队列先
					manager.clients[manager.index] = conn
					// 让先进来的人先下
					jsonMessage, _ := json.Marshal(&Result{bout: true, message:"系统：游戏开始，请您先落子", color:"black"})
					manager.clients[manager.index-1].send(jsonMessage)
					// 给自己的提示
					jsonMessage, _ := json.Marshal(&Result{bout: false, message:"系统：游戏开始，请等待对手落子！", color:"white"})
					conn.send(jsonMessage)
				} else {
					// 对手不存在，把自己加入奇数位置
					manager.index--
					manager.clients[manager.index] = conn
					jsonMessage, _ := json.Marshal(&Result{message:"系统：等待玩家匹配！"})
					conn.send(jsonMessage)
				}
			} else {
				// 自己本身就是奇数位置，没有对手，加入到队列并且返回信息
				jsonMessage, _ := json.Marshal(&Result{message:"系统：等待玩家匹配！"})
				conn.send(jsonMessage)
				manager.clients[manager.index] = conn
			}

		// 该玩家退出房间
		case conn := <-manager.unregister:
			//判断连接的状态，如果是true,就关闭send，删除连接client的值
			if(manager.clients[conn] != nil) {
				// 先找到对手
				if (conn % 2 == 0) {
					opp := manager.clients[conn - 1]
				}else{
					opp := manager.clients[conn + 1]
				}
				// 删除自己
				delete(manager.clients, manager.clients[conn])
				if (opp != nil) {
					jsonMessage, _ := json.Marshal(&Result{message:"系统：你的对手已离开！"})
					conn.send(jsonMessage)
				}
			}
		// 有信息进来
		}
	}
}

func (c *Client) send(message []byte){
	if (message != nil) {
		c.socket.WriteMessage(websocket.TextMessage, message)
	} else {
		c.socket.WriteMessage(websocket.CloseMessage, []byte{})
	}
}



func main() {
	fmt.Println("Starting application...")
	//开一个goroutine执行开始程序
	go manager.start()
	//注册默认路由为 /ws ，并使用wsHandler这个方法
	http.HandleFunc("/ws", wsHandler)
	//监听本地的8011端口
	http.ListenAndServe(":8011", nil)
}

func wsHandler(res http.ResponseWriter, req *http.Request) {
	//将http协议升级成websocket协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	//每一次连接都会新开一个client，client.id通过uuid生成保证每次都是不同的
	client := &Client{id: uuid.Must(uuid.NewV4()).String(), socket: conn, send: make(chan []byte)}
	//注册一个新的链接
	manager.register <- client

	// //启动协程收web端传过来的消息
	// go client.read()
	// //启动协程把消息返回给web端
	// go client.write()
}
