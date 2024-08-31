package clientCon

import (
	"fmt"
	"tesodev-korpes/ConsumerService/internal/types"
	"tesodev-korpes/pkg"
)

// CustomerClient handles requests to the customer service.
type ConsumerClient struct {
	RestClient *pkg.RestClient
}

// NewCustomerClient creates a new instance of CustomerClient.
func NewConsumerClient(restClient *pkg.RestClient) *ConsumerClient {
	return &ConsumerClient{
		RestClient: restClient,
	}
}

// GetCustomerByID fetches a customer by ID using the RestClient.
func (c *ConsumerClient) GetOrderByID(orderID string, token string) (*types.OrderResponseModel, error) {
	var res types.OrderResponseModel
	uri := fmt.Sprintf("http://localhost:1881/order/%s", orderID)
	err := c.RestClient.DoGetRequest(uri, &res, token)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *ConsumerClient) GetCustomerByID(customerID string, token string) (*types.CustomerResponseModel, error) {
	var res types.CustomerResponseModel
	uri := fmt.Sprintf("http://localhost:1907/customer/%s", customerID)
	err := c.RestClient.DoGetRequest(uri, &res, token)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *ConsumerClient) GetToken() (*types.TokenResponseModel, error) {
	var res types.TokenResponseModel
	uri := fmt.Sprintf("http://localhost:1907/login")
	err := c.RestClient.DoGetRequest(uri, &res, "")
	if err != nil {
		return nil, err
	}
	return &res, nil
}
