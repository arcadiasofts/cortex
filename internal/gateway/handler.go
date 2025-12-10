// internal/gateway/handler.go
package gateway

import (
	"backend/server/pb"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// 웹소켓 업그레이드 미들웨어 (Fiber 필수!)
func (h *Handler) UpgradeMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// 실제 웹소켓 엔드포인트
func (h *Handler) HandleWS(c *websocket.Conn) {
	userID := "user_temp_123" // 테스트용 임시 ID

	// 2. 클라이언트 생성
	client := &Client{
		Hub:    h.hub,
		UserID: userID,
		Conn:   c,
		Send:   make(chan *pb.GatewayPayload, 256),
	}

	// 3. 허브에 등록
	h.hub.Register <- client

	// 4. 펌프 가동
	go client.WritePump()
	client.ReadPump()
}
