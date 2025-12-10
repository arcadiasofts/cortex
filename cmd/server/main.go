package main

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(),
	).Run()
}

func newFiberApp() *fiber.App {
	app := fiber.New(
		fiber.Config{
			AppName: "Cortex Server v01",
		})
	return app
}

func RegisteAndStart(
	lc fx.Lifecycle,
) {

}
