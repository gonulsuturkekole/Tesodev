package internal

import (
	"github.com/labstack/echo/v4"
	_ "go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service}

	g := e.Group("/order")
	g.GET("/:id", handler.GetByID)
	g.POST("/", handler.Create)
	g.PUT("/:id", handler.Update)
	g.PATCH("/:id", handler.PartialUpdate)
	g.DELETE("/:id", handler.Delete)
}

func (h *Handler) GetByID(c echo.Context) error {
	id := c.Param("id")
	order, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	orderResponse := ToOrderResponse(order)
	return c.JSON(http.StatusOK, orderResponse)
}

func (h *Handler) Create(c echo.Context) error {
	var customer interface{}

	if err := c.Bind(&customer); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := h.service.Create(c.Request().Context(), customer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"message":    "Successed!",
		"creadtedId": id.Hex(),
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var update interface{}
	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.service.Update(c.Request().Context(), id, update); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Order updated successfully")
}

func (h *Handler) PartialUpdate(c echo.Context) error {
	id := c.Param("id")
	var update interface{}
	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.service.Update(c.Request().Context(), id, update); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Order partially updated successfully")
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Order deleted successfully")
}
