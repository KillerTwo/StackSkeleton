package Core

import (
	"GinSkeleton/App/Global/MyErrors"
	"GinSkeleton/App/Global/Variable"
	"GinSkeleton/App/Utils/Config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Client struct {
	Hub                *Hub            // 负责处理客户端注册、注销、在线管理
	Conn               *websocket.Conn // 一个ws连接
	Send               chan []byte     // 一个ws连接存储自己的消息管道
	PingPeriod         time.Duration
	PongWait           time.Duration
	WriteWait          time.Duration
	HeartbeatFailTimes int
}

// 处理握手+协议升级
func (c *Client) OnOpen(context *gin.Context) (*Client, bool) {
	// 1.升级连接,从http--->websocket

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  Config.CreateYamlFactory().GetInt("Websocket.WriteReadBufferSize"),
		WriteBufferSize: Config.CreateYamlFactory().GetInt("Websocket.WriteReadBufferSize"),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 2.将http协议升级到websocket协议.初始化一个有效的websocket长连接客户端
	if ws_conn, err := upgrader.Upgrade(context.Writer, context.Request, nil); err != nil {
		log.Panic(MyErrors.Errors_Websocket_OnOpen_Fail, err.Error())
		return nil, false
	} else {
		if ws_hub, ok := Variable.Websocket_Hub.(*Hub); ok {
			c.Hub = ws_hub
		}

		c.Conn = ws_conn
		c.Send = make(chan []byte, Config.CreateYamlFactory().GetInt("Websocket.WriteReadBufferSize"))
		c.PingPeriod = Config.CreateYamlFactory().GetDuration("Websocket.PingPeriod")
		c.PongWait = Config.CreateYamlFactory().GetDuration("Websocket.PingPeriod") * 10 / 9
		c.WriteWait = Config.CreateYamlFactory().GetDuration("Websocket.WriteWait")
		c.Hub.Register <- c
		ws_conn.WriteMessage(websocket.TextMessage, []byte(Variable.Websocket_Handshake_Success))
		return c, true
	}

}

// 主要功能主要是实时接收消息
func (c *Client) ReadPump(callback_on_message func(messageType int, p []byte), callback_on_error func(err error), callback_on_close func()) {

	// 回调 onclose 事件
	defer func() {
		callback_on_close()
	}()

	// OnMessage事件
	c.Conn.SetReadDeadline(time.Now().Add(c.PongWait))                                   // 设置最大读取时间
	c.Conn.SetReadLimit(Config.CreateYamlFactory().GetInt64("Websocket.MaxMessageSize")) // 设置最大读取长度
	for {
		messageType, byte_message, err := c.Conn.ReadMessage()
		if err == nil {
			callback_on_message(messageType, byte_message)
		} else {
			callback_on_error(err)
			//c.HeartbeatFailTimes++
			break
			// 关闭客户端以心跳检测为准，并不是发生一次错误就立刻关闭
			if c.HeartbeatFailTimes > Config.CreateYamlFactory().GetInt("Websocket.HeartbeatFailMaxTimes") {
				break
			}
		}
	}
}

// 按照websocket标准协议实现隐式心跳,Server端向Client远端发送ping格式数据包
func (c *Client) Heartbeat() {

	//2.浏览器收到服务器的ping格式消息，会自动原路返回
	c.Conn.SetPongHandler(func(pong string) error {
		c.Conn.SetReadDeadline(time.Now().Add(c.PongWait))
		fmt.Println("浏览器自动响应服务器发出去的ping格式消息：", pong)
		return nil
	})

	//  1. 设置一个时钟，周期性的向client远端发送心跳数据包
	ticker := time.NewTicker(c.PingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(c.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("Server->Ping->Client")); err != nil {
				// 这里可以计算累积出现错误的次数，超过某个值，就关闭连接
				c.HeartbeatFailTimes++
				if c.HeartbeatFailTimes > Config.CreateYamlFactory().GetInt(Variable.Websocket_Server_Ping_Msg) {
					return
				}
			} else {
				if c.HeartbeatFailTimes > 0 {
					c.HeartbeatFailTimes++
				}
			}
		}
	}
}
