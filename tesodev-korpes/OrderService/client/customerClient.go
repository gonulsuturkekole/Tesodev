package client

import (
	"fmt"
	"tesodev-korpes/OrderService/internal/types"
	"tesodev-korpes/pkg"
)

// CustomerClient handles requests to the customer service.
type CustomerClient struct {
	RestClient *pkg.RestClient
}

// NewCustomerClient creates a new instance of CustomerClient.
func NewCustomerClient(restClient *pkg.RestClient) *CustomerClient {
	return &CustomerClient{
		RestClient: restClient,
	}
}

// GetCustomerByID fetches a customer by ID using the RestClient.
func (c *CustomerClient) GetCustomerByID(customerID string, token string) (*types.OrderResponseModel, error) {
	var res types.OrderResponseModel
	uri := fmt.Sprintf("http://localhost:8001/customer/%s", customerID)
	err := c.RestClient.DoGetRequest(uri, &res, token)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
