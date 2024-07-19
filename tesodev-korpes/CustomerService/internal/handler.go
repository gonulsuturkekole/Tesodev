package internal

import (
	"github.com/labstack/echo/v4"
	_ "go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"tesodev-korpes/CustomerService/internal/types"
)

type Handler struct {
	service *Service
}

func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service}

	g := e.Group("/customer")
	g.GET("/:id", handler.GetByID)
	g.POST("/", handler.Create)
	g.PUT("/:id", handler.Update)
	g.PATCH("/:id", handler.PartialUpdate)
	g.DELETE("/:id", handler.Delete)
	e.GET("/customers/", handler.Get)
}

func (h *Handler) GetByID(c echo.Context) error {
	//name := c.QueryParam("name")
	id := c.Param("id")
	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	customerResponse := ToCustomerResponse(customer)
	return c.JSON(http.StatusOK, customerResponse)
}

func (h *Handler) Get(c echo.Context) error {
	name := c.QueryParam("name")
	lastName := c.QueryParam("last_name")
	age := c.QueryParam("age")

	customers, err := h.service.Get(c.Request().Context(), name, lastName, age)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "error fetching customers"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "customer fetch",
		"data":    customers,
	})
}

func (h *Handler) Create(c echo.Context) error {
	var customer *types.Customer

	if err := c.Bind(&customer); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	id, err := h.service.Create(c.Request().Context(), customer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	response := map[string]interface{}{
		"message":   "Successed!",
		"createdId": id,
	}
	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var customer types.CustomerUpdateModel
	if err := c.Bind(&customer); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.service.Update(c.Request().Context(), id, customer); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Customer updated successfully",
	})
}

func (h *Handler) PartialUpdate(c echo.Context) error {
	id := c.Param("id")
	var customer types.CustomerUpdateModel
	if err := c.Bind(&customer); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.service.Update(c.Request().Context(), id, customer); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Customer partially updated successfully",
	})
}
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Customer deleted successfully",
	})
}
