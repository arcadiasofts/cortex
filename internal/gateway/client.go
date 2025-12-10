// internal/gateway/client.go
package gateway

import (
	"backend/server/pb"
	"log"

	"github.com/gofiber/contrib/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	Hub    *Hub
	UserID string
	Conn   *websocket.Conn
	Send   chan *pb.GatewayPayload // 보낼 메시지 대기열
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		// 1. 소켓에서 바이트 읽기
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break // 연결 끊김
		}

		// 2. 바이트 -> GatewayPayload 구조체로 변환 (Unmarshal)
		payload := &pb.GatewayPayload{}
		if err := proto.Unmarshal(data, payload); err != nil {
			log.Println("Proto Error:", err)
			continue
		}

		// 3. OpCode에 따른 처리
		c.handleMessage(payload)
	}
}

func (c *Client) handleMessage(payload *pb.GatewayPayload) {
	switch payload.Op {
	case 0:
		c.Hub.Broadcast <- payload

	case 1: // HEARTBEAT
		// 퐁(ACK) 보내기
		c.Send <- &pb.GatewayPayload{Op: 11}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 허브가 채널을 닫음 (강퇴)
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 1. 구조체 -> 바이트 변환 (Marshal)
			data, err := proto.Marshal(message)
			if err != nil {
				continue
			}

			// 2. 바이너리 메시지로 전송 (Text 아님!)
			if err := c.Conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				return
			}
		}
	}
}
