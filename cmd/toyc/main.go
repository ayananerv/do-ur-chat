package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	// ⚠️ 记得把这里替换为你 go.mod 中的实际 module 名
	pb "github.com/ayananerv/do-ur-chat/api/gen/message"
)

// 定义命令行参数
var (
	serverAddr = flag.String("addr", "localhost:8080", "服务器地址")
	myUID      = flag.Int64("uid", 1, "当前登录的用户ID")
	targetUID  = flag.Int64("to", 2, "默认发送给谁 (目标用户ID)")
)

func main() {
	flag.Parse() // 解析命令行参数

	// 1. 建立连接
	u := url.URL{Scheme: "ws", Host: *serverAddr, Path: "/ws", RawQuery: fmt.Sprintf("uid=%d", *myUID)}
	log.Printf("正在连接服务器: %s ...", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer c.Close()

	log.Printf("✅ 登录成功! 我是用户 [%d], 默认发给 [%d]", *myUID, *targetUID)
	log.Println("👉 请在控制台输入消息并回车发送 (输入 'exit' 退出):")

	// 2. 启动接收协程 (后台监听)
	done := make(chan struct{})
	go readLoop(c, done)

	// 3. 启动发送循环 (监听键盘输入)
	inputLoop(c)

	// 等待接收协程退出（防止主程序过早结束）
	<-done
}

// 接收消息的逻辑
func readLoop(c *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, messageData, err := c.ReadMessage()
		if err != nil {
			log.Printf("❌ 连接断开: %v", err)
			return
		}

		// 反序列化
		msg := &pb.ChatMessage{}
		if err := proto.Unmarshal(messageData, msg); err != nil {
			log.Printf("数据格式错误: %v", err)
			continue
		}

		// 漂亮的打印格式
		fmt.Printf("\n📩 [收到新消息] From User %d:\n   %s\n👉 请输入: ",
			msg.SenderId,
			string(msg.Content),
		)
	}
}

// 发送消息的逻辑 (读取控制台)
func inputLoop(c *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	// 监听 Ctrl+C 信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		fmt.Print("👉 请输入: ")

		// 使用 select 实现非阻塞监听，防止 Ctrl+C 无法退出
		ch := make(chan string)
		go func() {
			if scanner.Scan() {
				ch <- scanner.Text()
			} else {
				close(ch)
			}
		}()

		select {
		case text, ok := <-ch:
			if !ok {
				return
			} // 读不到输入了
			if strings.TrimSpace(text) == "" {
				continue
			}
			if text == "exit" {
				return
			}

			// 构造 Protobuf 消息
			msg := &pb.ChatMessage{
				MsgId:      fmt.Sprintf("msg_%d", time.Now().UnixNano()), // 模拟唯一ID
				SessionId:  "demo_session",
				SenderId:   *myUID,
				ReceiverId: *targetUID, // 发给命令行指定的那个 ID
				Type:       pb.MessageType_TEXT,
				Content:    []byte(text),
				Timestamp:  time.Now().UnixMilli(),
			}

			// 序列化 & 发送
			data, _ := proto.Marshal(msg)
			err := c.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				log.Printf("发送失败: %v", err)
				return
			}
			// log.Println("✔ 已发送") // 保持界面清爽，发送成功不刷屏

		case <-interrupt:
			log.Println("程序退出")
			return
		}
	}
}
