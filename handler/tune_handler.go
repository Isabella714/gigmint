package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Isabella714/gigmint/processor"
)

type TuneHandler struct {
	tuneProcessor *processor.TuneProcessor
}

func RegisterTuneHandler(app *fiber.App) {
	_ = &TuneHandler{
		tuneProcessor: processor.NewTuneProcessor(),
	}
}
