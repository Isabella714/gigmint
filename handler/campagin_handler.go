package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Isabella714/gigmint/model/dto"
	"github.com/Isabella714/gigmint/processor"
)

type CampaignHandler struct {
	campaignProcessor *processor.CampaignProcessor
}

func RegisterCampaignHandler(app *fiber.App) {
	handler := &CampaignHandler{
		campaignProcessor: processor.NewCampaignProcessor(),
	}

	app.Get("/campaign", handler.PagingCampaign)

}

func (h *CampaignHandler) PagingCampaign(ctx fiber.Ctx) error {
	request := new(dto.PagingCampaignRequest)
	if err := ctx.Bind().Query(request); err != nil {
		return err
	}

	campaigns, err := h.campaignProcessor.PagingCampaign(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(campaigns)
}
