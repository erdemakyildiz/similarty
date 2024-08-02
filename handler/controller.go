package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"net/http"
	"similarty-engine/model"
	"similarty-engine/service"
	"strconv"
	"strings"
)

type RestController struct {
	simService *service.SimilarityService
}

func NewRestController(simService *service.SimilarityService) RestController {
	return RestController{simService: simService}
}

func (rc *RestController) FilterStrings(ctx fiber.Ctx) error {
	fr := ctx.Query("frequency")
	frequency, err := strconv.ParseFloat(strings.TrimSpace(fr), 64)
	if err != nil {
		//default value
		frequency = 0.8
	}

	request := new(model.FilterLines)
	if err := ctx.Bind().Body(request); err != nil {
		return err
	}

	if len(request.Lines) == 0 {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	result := rc.simService.FilterLines(request.Lines, frequency)

	appsJson, _ := json.Marshal(result)
	ctx.Set("Content-type", "application/json; charset=utf-8")
	return ctx.SendString(string(appsJson))
}
