package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"tesodev-korpes/CustomerService/authentication"
	"tesodev-korpes/CustomerService/internal/types"
	"time"
	"unicode"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetCustomersWithSecondLetterA(ctx context.Context) ([]types.Customer, error) {
	return s.repo.GetCustomersWithSecondLetterA(ctx)
}

//func (s *Service) GetUser(ctx context.Context, username string) (*types.CustomerResponseModel, error) {
//	user, err := s.repo.GetUser(ctx, username)
//	if err != nil {
//		return nil, err
//	}
//	return user, nil
//}

func (s *Service) GetByID(ctx context.Context, id string) (*types.Customer, error) {
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

	if len(customer.LastName) != 0 {

		start := time.Now()

		done := make(chan struct{})

		go func(name string) {
			defer close(done)
			time.Sleep(1 * time.Nanosecond)
			customer.LastName = strings.ToUpper(name[:1]) + strings.ToLower(name[1:])

		}(customer.LastName)
		<-done
		fmt.Printf("Last Name formatted in : %v nanoseconds by go routine\n", time.Since(start).Nanoseconds())
	}

	if len(customer.LastName) != 0 {

		start := time.Now()
		time.Sleep(1 * time.Nanosecond)

		customer.LastName = func(name string) string {
			return strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
		}(customer.LastName)
		fmt.Printf("Last Name formatted in : %v nanoseconds by function\n", time.Since(start).Nanoseconds())
	}

	if len(customer.Phone) != 0 {
		var formattedPhone string
		for i, char := range customer.Phone {
			if unicode.IsDigit(char) {
				if i != len(customer.Phone)-1 {
					formattedPhone += string(char)
					formattedPhone += "-"
				} else {
					formattedPhone += string(char)
				}

			}
		}
		customer.Phone = formattedPhone
	}

	if customer.AdditionalInfo == nil {
		customer.AdditionalInfo = make(map[string]string)
		switch customer.MembershipType {
		case "standard":
			customer.AdditionalInfo["membership_type"] = "Standard"
			customer.AdditionalInfo["free_shipping"] = "1"
		case "premium":
			customer.AdditionalInfo["membership_type"] = "Premium"
			customer.AdditionalInfo["free_shipping"] = "5"
		case "gold":
			customer.AdditionalInfo["membership_type"] = "Gold"
			customer.AdditionalInfo["free_shipping"] = "100"
			customer.AdditionalInfo["priority_in_customer_line"] = "yes"
			customer.AdditionalInfo["discover"] = "%5"
		default:
			customer.AdditionalInfo["membership_type"] = "None"
			customer.AdditionalInfo["free_shipping"] = "0"
		}
		fmt.Println("Customer ID:", id)
		fmt.Println("Customer Membership Type:", customer.MembershipType)
		fmt.Println("Additional Info:", customer.AdditionalInfo)
	}

	customer.ContactOption = append(customer.ContactOption, "")
	//customer.ContactOption = []string{""}

	return customer, nil
}

// Create method creates a new customer with a custom UUID as the ID
func (s *Service) Create(ctx context.Context, customerRequestModel types.CustomerRequestModel) (string, error) {
	hashedPassword, err := authentication.HashPassword(customerRequestModel.Password)
	if err != nil {
		return "", err
	}
	// Check if the customer data is valid
	//if err := validateCustomer(customerRequestModel); err != nil {
	//fmt.Println("Invalid customer data:", err)
	//return "", err
	// Generate a new UUID
	customID := uuid.New().String()
	now := time.Now().Local()
	customerRequestModel.CreatedAt = now
	// Set the customer's ID to the generated UUID
	//customerRequestModel.Id = customID
	// Insert the customer data into MongoDB

	customer := &types.Customer{
		FirstName: customerRequestModel.FirstName,
		LastName:  customerRequestModel.LastName,
		Age:       customerRequestModel.Age,
		Email:     customerRequestModel.Email,
		CreatedAt: customerRequestModel.CreatedAt,
		Id:        customID,
		Username:  customerRequestModel.Username,
		Password:  hashedPassword, //string(hashPassword)
	}
	//customer.Id = customID

	_, err = s.repo.Create(ctx, customer)
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
	customer.ContactOption = customerUpdateModel.ContactOption
	customer.MembershipType = customerUpdateModel.MembershipType
	customer.UpdatedAt = now
	return s.repo.Update(ctx, id, customer)
}
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetCustomers(ctx context.Context, params types.QueryParams, page, limit int) ([]types.Customer, int64, error) {
	return s.repo.GetCustomersByFilter(ctx, params.FirstName, params.AgeGreaterThan, params.AgeLessThan, page, limit)
}

/*func containsDigit(s string) bool {
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
		return true

	}
	return s[0] >= 'A' && s[0] <= 'Z'
}

// validateCustomer checks if the customer data is valid when it's created.
// 2) do something with switch-case
func validateCustomer(customer types.CustomerRequestModel) error {
	switch {
	case containsDigit(customer.FirstName):
		return errors.New("First name contains a number")
	case !startsWithUpperCase(customer.FirstName):
		return errors.New("First name does not start with an uppercase letter")
	case customer.FirstName == "":
		return errors.New("First name is empty")
	default:
		fmt.Printf("Customer '%s' is valid.\n", customer.FirstName)
	}
	return nil
} */
