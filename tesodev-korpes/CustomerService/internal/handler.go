package internal

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"tesodev-korpes/CustomerService/authentication"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/pkg/log"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

// @title Customer Service API
// @version 1.0
// @description API documentation for Customer Service.
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8001
// @BasePath /api/v1
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func NewHandler(e *echo.Echo, service *Service) {
	handler := &Handler{service: service, validate: validator.New()}

	g := e.Group("/customer")
	g.GET("/:id", handler.GetByID)
	g.POST("/", handler.Create)
	g.PUT("/:id", handler.Update)
	g.PATCH("/:id", handler.PartialUpdate)
	g.DELETE("/:id", handler.Delete)

	e.GET("/customers", handler.GetCustomersByFilter)
	e.POST("/login", handler.Login)
	e.GET("/verify", handler.Verify)
}

// Login handles user authentication and returns a JWT token.
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param user body types.Customer true "User credentials"
// @Success 200 {object} types.Customer
// @Failure 400 {object} string "Invalid input"
// @Failure 401 {object} string "Invalid credentials"
// @Failure 500 {object} string "Internal server error"
// @Security BearerAuth
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
// @Param Authorization header string true "JWT token"
// @Success 200 {string} string "Token verified and user exists"
// @Failure 401 {object} map[string]string "Invalid or expired token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
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
// @Tags customer
// @Produce  json
// @Param id path string true "Customer ID"
// @Param Authentication header string true "JWT token"
// @Success 200 {object} types.CustomerResponseModel
// @Failure 400 {object} string  "Invalid customer ID"
// @Failure 404 {object}  string "Customer not found"
// @Failure 500 {object} string "Internal server error"
// @Router /customer/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	id := c.Param("id")

	customer, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if customer == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Customer not found"})
	}

	customerResponse := ToCustomerResponse(customer)

	return c.JSON(http.StatusOK, customerResponse)
}

// Create adds a new customer to the database.
// @Summary Create a new customer
// @Description Create a new customer with the provided details
// @Tags customer
// @Accept  json
// @Produce  json
// @Param customer body types.CustomerRequestModel true "Customer data"
// @Success 201 {object} map[string]interface{} "Customer created"
// @Failure 400 {object} string "Invalid customer data"
// @Failure 500 {object} string "Internal server error"
// @Router /customer/ [post]
func (h *Handler) Create(c echo.Context) error {
	var customerRequestModel types.CustomerRequestModel

	if err := c.Bind(&customerRequestModel); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := ValidateCustomer(&customerRequestModel, h.validate); err != nil {
		if valErr, ok := err.(*ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
				"errors":  valErr.Errors,
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	id, err := h.service.Create(c.Request().Context(), customerRequestModel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	log.Info("Customer created")

	response := map[string]interface{}{
		"message":   "Succeeded!",
		"createdId": id,
	}
	return c.JSON(http.StatusCreated, response)
}

// Update modifies an existing customer's details.
// @Summary Update customer details
// @Description Update customer details with the given data
// @Tags customer
// @Accept  json
// @Produce  json
// @Param id path string true "Customer ID"
// @Param customer body types.CustomerUpdateModel true "Customer data"
// @Param Authentication header string true "JWT token"
// @Success 200 {object} map[string]string "Customer updated successfully"
// @Failure 400 {object} string "Invalid input"
// @Failure 404 {object} string "Customer not found"
// @Failure 500 {object} string "Internal server error"
// @Router /customer/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")

	var customer types.CustomerUpdateModel
	if err := c.Bind(&customer); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if err := h.service.Update(c.Request().Context(), id, customer); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Customer updated successfully",
	})
}

// PartialUpdate modifies specific fields of an existing customer.
// @Summary Partially update customer details
// @Description Partially update customer details with the given data
// @Tags customer
// @Accept  json
// @Produce  json
// @Param id path string true "Customer ID"
// @Param customer body types.CustomerUpdateModel true "Customer data"
// @Param Authentication header string true "JWT token"
// @Success 200 {object} map[string]string "Customer partially updated successfully"
// @Failure 400 {object} string "Invalid input"
// @Failure 404 {object} string "Customer not found"
// @Failure 500 {object} string "Internal server error"
// @Router /customer/{id} [patch]
func (h *Handler) PartialUpdate(c echo.Context) error {
	id := c.Param("id")

	var customer types.CustomerUpdateModel
	if err := c.Bind(&customer); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if err := h.service.Update(c.Request().Context(), id, customer); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Customer partially updated successfully",
	})
}

// Delete removes a customer from the database.
// @Summary Delete customer
// @Description Delete a customer by their ID
// @Tags customer
// @Produce  json
// @Param id path string true "Customer ID"
// @Param Authentication header string true "JWT token"
// @Success 200 {object} map[string]string "Customer deleted successfully"
// @Failure 404 {object} string "Customer not found"
// @Failure 500 {object} string "Internal server error"
// @Router /customer/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Customer deleted successfully",
	})
}

// GetCustomersByFilter retrieves customers based on query parameters such as first name and age range.
// @Summary Get customers by filter
// @Description Retrieve a list of customers based on optional filters like first name, age greater than, and age less than. Pagination is supported.
// @Tags customer
// @Accept  json
// @Produce  json
// @Param first_name query string false "Filter by first name"
// @Param agt query int false "Filter by age greater than"
// @Param alt query int false "Filter by age less than"
// @Param page query int true "Page number for pagination"
// @Param limit query int true "Number of items per page"
// @Success 200 {object} map[string]interface{} "Customer data retrieved successfully"
// @Failure 400 {object} string "Invalid page or limit parameter"
// @Failure 404 {object} string "No customers found"
// @Failure 500 {object} string "Error fetching customers"
// @Router /customers [get]
func (h *Handler) GetCustomersByFilter(c echo.Context) error {
	params := types.QueryParams{
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
