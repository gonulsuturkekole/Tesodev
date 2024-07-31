package internal

import (
	"fmt"
	"tesodev-korpes/OrderService/internal/types"
)

func (s *Service) getCustomerByID(customerID string) (*types.OrderResponseModel, error) {
	var res types.OrderResponseModel
	uri := fmt.Sprintf("http://localhost:8001/customer/%s", customerID)
	err := s.client.DoGetRequest(uri, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
