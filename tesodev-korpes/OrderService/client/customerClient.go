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
func (c *CustomerClient) GetCustomerByID(customerID string, token string) (*types.CustomerResponse, error) {
	var res types.CustomerResponse
	uri := fmt.Sprintf("http://localhost:8001/customer/%s", customerID)
	err := c.RestClient.DoGetRequest(uri, &res, token)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *CustomerClient) SendFinanceRequest(tokenString string) error {
	// Define the URL to which the POST request will be sent
	postUrl := "http://localhost:8003/finance"

	// Make the POST request using the RestClient, ignoring the response
	err := c.RestClient.DoPostRequest(postUrl, nil, nil, tokenString)
	if err != nil {
		return fmt.Errorf("failed to send finance request: %w", err)
	}

	// No response needed, just return nil on success
	return nil
}
