package routes

import (
	"backend/internal/gateway"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	gwHandler *gateway.Handler, // Fx가 주입해줌
) {
	api := app.Group("/api/v1")

	api.Use()

	// 웹소켓 라우트 등록
	app.Use("/gateway", gwHandler.UpgradeMiddleware)
	app.Get("/gateway", websocket.New(
		func(c *websocket.Conn) {
			gwHandler.HandleWS(c)
		}))
}
