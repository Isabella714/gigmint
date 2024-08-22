package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/Isabella714/gigmint/component"
	"github.com/Isabella714/gigmint/handler"
)

func main() {
	err := component.Init()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(fiberRecover.New(), cors.New(), logger.New(), requestid.New())

	handler.RegisterTuneHandler(app)
	handler.RegisterCampaignHandler(app)

	//err = scanner.StartJobScheduler()
	//if err != nil {
	//	panic(err)
	//}

	log.Fatal(app.Listen(":3000"))
}
