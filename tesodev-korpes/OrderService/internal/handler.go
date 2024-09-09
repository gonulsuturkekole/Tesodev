package internal

import (
	"github.com/labstack/echo/v4"
	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"tesodev-korpes/OrderService/internal/types"
)

type Handler struct {
	service *Service
}

// Register the endpoint using the repository

// Start the server

func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service}

	g := e.Group("/order")
	g.GET("/:id", handler.GetByID)
	g.POST("/:customer_id", handler.CreateOrder)
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
func (h *Handler) CreateOrder(c echo.Context) error {
	var order types.OrderRequestModel

	customerID := c.Param("customer_id")
	token := c.Request().Header.Get("Authentication")

	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, totalOrders, err := h.service.CreateOrderService(c.Request().Context(), customerID, &order, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	response := map[string]interface{}{
		"message":     "Success!",
		"createdId":   id,
		"totalOrders": totalOrders,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var order types.OrderUpdateModel
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.service.Update(c.Request().Context(), id, order); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order updated successfully",
	})
}

func (h *Handler) PartialUpdate(c echo.Context) error {
	id := c.Param("id")
	var order types.OrderUpdateModel
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.service.Update(c.Request().Context(), id, order); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order partially updated successfully",
	})
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Order deleted successfully",
	})
}
