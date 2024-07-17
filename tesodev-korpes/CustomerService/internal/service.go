package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"tesodev-korpes/CustomerService/internal/types"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (*types.Customer, error) {
	//var customer *types.Customer
	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	//challenge (everything should be observable somehow in the response or console (print)):
	// 1) do something with using for loop by using customer model and manipulate it (you can add an additional field for it)
	// 2) do something with switch-case
	// 3) do something with goroutines (you should give us an example for both scenarios of not using goroutines and using)
	// 3.1) calculate the elapsed time for both scenarios and show us the gained time
	// 4) add an additional field and use maps
	// 5) add an additional field and use arrays
	// 6) manipulate an existing data to see how pointers and values work
	// Manipulate data using pointer

	return customer, nil
}

// Create method creates a new customer with a custom UUID as the ID
func (s *Service) Create(ctx context.Context, customer *types.Customer) (string, error) {
	// Check if the customer data is valid
	if err := validateCustomer(customer); err != nil {
		fmt.Println("Invalid customer data:", err)
		return "", err
	}
	// Generate a new UUID
	customID := uuid.New().String()
	now := time.Now().Local()
	customer.CreatedAt = now
	// Set the customer's ID to the generated UUID
	customer.Id = customID
	// Insert the customer data into MongoDB

	_, err := s.repo.Create(ctx, customer)
	if err != nil {
		return "", err
	}
	// Return the generated ID if the insertion is successful
	return customID, nil
}

func (s *Service) Update(ctx context.Context, id string, customerUpdateModel types.CustomerUpdateModel) error {
	// Create an update document
	customer, err := s.GetByID(ctx, id)
	now := time.Now().Local()
	if err != nil {
		return err
	}

	customer.FirstName = customerUpdateModel.FirstName
	customer.LastName = customerUpdateModel.LastName
	customer.Phone = customerUpdateModel.Phone
	customer.UpdatedAt = now
	return s.repo.Update(ctx, id, customer)
}
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func containsDigit(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			return true
		}
	}
	return false
}

// startsWithUpperCase checks if a string starts with an uppercase letter.
func startsWithUpperCase(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] >= 'A' && s[0] <= 'Z'
}

// validateCustomer checks if the customer data is valid when it's created.
// 2) do something with switch-case
func validateCustomer(customer *types.Customer) error {
	switch {
	case containsDigit(customer.FirstName):
		return errors.New("First name contains a number")
	case !startsWithUpperCase(customer.FirstName):
		return errors.New("First name does not start with an uppercase letter")
	default:
		fmt.Printf("Customer '%s' is valid.\n", customer.FirstName)
	}
	return nil
}
