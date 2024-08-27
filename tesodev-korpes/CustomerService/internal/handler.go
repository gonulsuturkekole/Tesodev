package internal

import (
	_ "fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	_ "strconv"
	"tesodev-korpes/CustomerService/authentication"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/pkg/log"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

// NewHandler initializes the routes and sets up the handlers.
// @title Customer Service API
// @version 2.0
// @description API documentation for Customer Service.
// @BasePath /api/v1
func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service, validate: validator.New()}

	g := e.Group("/customer")
	g.GET("/:id", handler.GetByID)
	g.POST("/", handler.Create)
	g.PUT("/:id", handler.Update)
	g.PATCH("/:id", handler.PartialUpdate)
	g.DELETE("/:id", handler.Delete)

	e.POST("/login", handler.Login)
	e.GET("/verify", handler.Verify)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

}

// Login handles user authentication and returns a JWT token.
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param user body types.Customer true "User credentials"
// @Success 200 {object} types.Customer
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (h *Handler) Login(c echo.Context) error {
	var user types.Customer
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	result, err := h.service.GetByID(c.Request().Context(), user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if result == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID"})
	}
	if !authentication.CheckPasswordHash(user.Password, result.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid password"})
	}
	result.Token = authentication.JwtGenerator(result.Id, result.FirstName, result.LastName, "secret")

	resp := c.JSON(http.StatusOK, result)
	log.Info("Status Ok")
	return resp
}

// Verify validates the JWT token and checks if the user exists.
// @Summary Verify JWT token
// @Description Verify JWT token and check user existence
// @Tags authentication
// @Produce  json
// @Param Authentication header string true "JWT token"
// @Success 200 {string} string "Token verified and user exists"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /verify [get]
func (h *Handler) Verify(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authentication")
	token, err := jwt.ParseWithClaims(authHeader, &authentication.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return authentication.SecretKey, nil
	})
	if err != nil || !token.Valid {
		c.Logger().Error("Token verification failed: ", err)
		return echo.ErrUnauthorized
	}
	claims := token.Claims.(*authentication.Claims)
	userID := claims.ID

	exists, err := h.service.ExistsbyID(c.Request().Context(), userID)
	if err != nil {
		c.Logger().Error("Error checking user existence: ", err)
		return echo.ErrInternalServerError
	}

	if !exists {
		c.Logger().Error("User does not exist")
		return echo.ErrUnauthorized
	}

	return c.String(http.StatusOK, "Token verified and user exists")
}

// GetByID retrieves a customer by their ID.
// @Summary Get customer by ID
// @Description Get customer details by ID
// @Tags customers
// @Produce  json
// @Param id path string true "Customer ID"
// @Success 200 {object} types.CustomerResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	id := c.Param("id")

	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	customerResponse := ToCustomerResponse(customer)

	return c.JSON(http.StatusOK, customerResponse)
}

// Create adds a new customer to the database.
// @Summary Create a new customer
// @Description Create a new customer with the provided details
// @Tags customers
// @Accept  json
// @Produce  json
// @Param customer body types.CustomerRequestModel true "Customer data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customer [post]
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
	log.Info("customer created")

	response := map[string]interface{}{
		"message":   "Succeeded!",
		"createdId": id,
	}
	return c.JSON(http.StatusCreated, response)
}

// Update modifies an existing customer's details.
// @Summary Update customer details
// @Description Update customer details with the given data
// @Tags customers
// @Accept  json
// @Produce  json
// @Param id path string true "Customer ID"
// @Param customer body types.CustomerUpdateModel true "Customer data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer/{id} [put]
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

// PartialUpdate modifies specific fields of an existing customer.
// @Summary Partially update customer details
// @Description Partially update customer details with the given data
// @Tags customers
// @Accept  json
// @Produce  json
// @Param id path string true "Customer ID"
// @Param customer body types.CustomerUpdateModel true "Customer data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer/{id} [patch]
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

// Delete removes a customer from the database.
// @Summary Delete customer
// @Description Delete a customer by their ID
// @Tags customers
// @Produce  json
// @Param id path string true "Customer ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /customer/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Customer deleted successfully",
	})
}

// GetCustomersByFilter retrieves customers based on filtering criteria.
// @Summary Get customers by filter
// @Description Get customers based on various filtering options
// @Tags customers
// @Produce  json
// @Param first
