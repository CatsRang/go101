package example

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the example routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/example", h.Create)
}

func (h *Handler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	resp, err := h.service.Create(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, resp)
}
