package internal

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"tesodev-korpes/CustomerService/internal/types"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
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
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	/*if err := h.validate.Struct(customer); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}*/
	//Validation(&customerRequestModel)
	// Validate customer object
	if err := h.validate.Struct(customerRequestModel); err != nil {
		// Handle validation errors
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)

		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "email":
				errorMessages[fieldError.Field()] = "It is not valid email address"
			case "ageValidation":
				errorMessages[fieldError.Field()] = "Age must be a number greater than or equal to 18"
			default:
				errorMessages[fieldError.Field()] = "Required field"
			}
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation failed",
			"errors":  errorMessages,
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
	firstName := c.QueryParam("first_name")
	ageGreaterThan := c.QueryParam("age_greater_than")
	ageLessThan := c.QueryParam("age_less_than")

	// Call the service method to find customers by first name
	customers, err := h.service.GetCustomers(c.Request().Context(), firstName, ageGreaterThan, ageLessThan)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Error fetching customers"})
	}
	if len(customers) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, map[string]string{"message": "No customers found"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "customer fetch",
		"data":    customers,
	})
}
