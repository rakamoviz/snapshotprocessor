package notifications

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/streamprocessor"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	streamProcessor streamprocessor.StreamProcessor
}

type getResponse struct {
	Code string `json:"code"`
}

func NewHandler(streamProcessor streamprocessor.StreamProcessor) *handler {
	return &handler{streamProcessor: streamProcessor}
}

func (h *handler) Bind(group *echo.Group) {
	group.GET("", func(c echo.Context) error { return h.get(c) })
}

func (h *handler) get(c echo.Context) error {
	resp := getResponse{Code: "abc"}
	c.JSON(http.StatusOK, &resp)
	return nil
}
