package app

import (
	"fmt"

	"github.com/ayananerv/do-ur-chat/api/gen/message" // 确保路径与 go.mod 一致

	"google.golang.org/protobuf/proto"
)

func TestProto() {
	// 创建你在 proto 中定义的消息对象 [cite: 4, 5, 6, 7]
	m := &message.ChatMessage{
		MsgId:      "srv_12345",
		SessionId:  "sess_abc",
		SenderId:   1,
		ReceiverId: 2,
		Type:       message.MessageType_TEXT,
		Content:    []byte("Hello World"),
	}

	// 序列化
	data, _ := proto.Marshal(m)
	fmt.Printf("Encoded size: %d bytes\n", len(data))

	// 反序列化
	newMsg := &message.ChatMessage{}
	proto.Unmarshal(data, newMsg)
	fmt.Printf("Decoded content: %s\n", string(newMsg.Content))
}
