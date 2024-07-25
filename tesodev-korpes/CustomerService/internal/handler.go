package internal

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/pkg/log"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

type QueryParams struct {
	FirstName      string `json:"first_name"`
	AgeGreaterThan string `json:"agt"`
	AgeLessThan    string `json:"alt"`
}

func NewHandler(e *echo.Echo, service *Service) {

	handler := &Handler{service: service, validate: validator.New()}

	//handler.validate.RegisterValidation("ageValidation", ageValidation)
	//handler.validate.RegisterValidation("email", validateEmail)
	g := e.Group("/customer")
	g.GET("/:id", handler.GetByID)
	g.POST("/", handler.Create)
	g.PUT("/:id", handler.Update)
	g.PATCH("/:id", handler.PartialUpdate)
	g.DELETE("/:id", handler.Delete)
	e.GET("/customers", handler.GetCustomersByFilter) // Get endpoint for filter
}

func (h *Handler) GetByID(c echo.Context) error {
	id := c.Param("id")
	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	customerResponse := ToCustomerResponse(customer)
	return c.JSON(http.StatusOK, customerResponse)
}

func (h *Handler) Create(c echo.Context) error {
	var customerRequestModel types.CustomerRequestModel

	if err := c.Bind(&customerRequestModel); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := ValidateCustomer(&customerRequestModel, h.validate); err != nil {
		if valErr, ok := err.(*ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
				"errors":  valErr.Errors,
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}
	id, err := h.service.Create(c.Request().Context(), customerRequestModel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"message":   "Succeeded!",
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

func (h *Handler) GetCustomersByFilter(c echo.Context) error {
	params := QueryParams{
		FirstName:      c.QueryParam("first_name"),
		AgeGreaterThan: c.QueryParam("agt"),
		AgeLessThan:    c.QueryParam("alt"),
	}

	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid page parameter"})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid limit parameter"})
	}

	// Call the service method to find customers by first name
	customers, totalCount, err := h.service.GetCustomers(c.Request().Context(), params, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Error fetching customers"})
	}
	fmt.Printf("Total Customers: %d\n", totalCount)
	fmt.Printf("Customers: %v\n", customers)
	if len(customers) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, map[string]string{"message": "No customers found"})
	}

	return echo.NewHTTPError(http.StatusOK, map[string]interface{}{
		"message":     "customer fetch",
		"data":        customers,
		"total_count": totalCount,
	})
}
