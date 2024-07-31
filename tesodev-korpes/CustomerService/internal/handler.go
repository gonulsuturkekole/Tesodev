package internal

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	_ "strings"
	"tesodev-korpes/CustomerService/authentication"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/log"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

func NewHandler(e *echo.Echo, service *Service) {

	handler := &Handler{service: service, validate: validator.New()}

	g := e.Group("/customer")
	g.GET("/:id", handler.GetByID, pkg.Authenticate)
	g.POST("/", handler.Create, pkg.Authenticate)
	g.PUT("/:id", handler.Update, pkg.Authenticate)
	g.PATCH("/:id", handler.PartialUpdate, pkg.Authenticate)
	g.DELETE("/:id", handler.Delete, pkg.Authenticate)

	e.POST("/login", handler.Login)
	e.GET("/verify", handler.Verify, pkg.Authenticate)                  //error verdi
	e.GET("/customers", handler.GetCustomersByFilter, pkg.Authenticate) // Get endpoint for filter
}
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
func (h *Handler) Verify(c echo.Context) error {
	userid := c.Get("Id").(string)
	if userid == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	user, err := h.service.GetByID(c.Request().Context(), userid)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Username":   user.Username,
		"First Name": user.FirstName,
		"Last Name":  user.LastName,
		"Email":      user.Email,
		"Age":        user.Age,
	})
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
	////Check if the username exists or not
	//existingUser, err := h.service.GetUserById(c.Request().Context(), customerRequestModel.Username)
	//if existingUser != nil {
	//	return c.JSON(http.StatusBadRequest, map[string]interface{}{
	//		"message": "Username already exists",
	//	})
	//}
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
	//		"message": "Internal server error",
	//	})
	//}
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
