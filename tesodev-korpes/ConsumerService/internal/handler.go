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

	//g := e.Group("/order")
	//g.GET("/:id", handler.GetByID)
	e.POST("/finance", handler.Connect)

}

/*
	func (h *Handler) GetByID(c echo.Context) error {
		id := c.Param("id")
		order, err := h.service.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		orderResponse := ToOrderResponse(order)
		return c.JSON(http.StatusOK, orderResponse)
	}
*/
func (h *Handler) Connect(c echo.Context) error {
	// Kafka'dan gelen mesajları okumaya başla
	err := h.service.ConsumeTopic(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Kafka consumer started successfully",
	})
}

/*func (h *Handler) FinanceCal(c echo.Context) error {
	// Assuming `Consumer` is already initialized and connected
	consumer := h.service.Consumer

	// Define a model that matches the structure of your Kafka messages
	var orderID string

	// Read from the Kafka topic
	consumer.Read(orderID, func(orderID string, err error) {
		if err != nil {
			// Handle the error appropriately
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// Here you can use the orderID to perform further processing
		// For example, you might want to fetch the order from the database using the orderID
		order, err := h.service.GetByID(c.Request().Context(), orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		// Process the order (e.g., calculate finance details)
		// This is just a placeholder, replace with actual finance calculation logic
		financeDetails := CalculateFinanceDetails(order)

		// Return the finance details as a response
		c.JSON(http.StatusOK, financeDetails)
	})

	return nil
}

// Sample function to demonstrate finance calculation logic
func CalculateFinanceDetails(order *Order) map[string]interface{} {
	// Perform some calculation based on the order
	return map[string]interface{}{
		"orderID":        order.Id,
		"calculatedTax":  100.0, // Example value
		"totalWithTax":   order.OrderTotal + 100.0,
		"paymentMethod":  order.PaymentMethod,
		"shipmentStatus": order.ShipmentStatus,
	}
}
*/
