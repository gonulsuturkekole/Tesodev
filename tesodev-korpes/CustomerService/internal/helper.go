package internal

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"tesodev-korpes/CustomerService/internal/types"
)

// here is an example what this helper method that does data casting from db model to response model
// the return statement that I commented out repreents an introduction that how you can implement it
// you can delete after you'd completed the helper method, its a placeholder put here just to prevent getting errors at
// the beginning
func ToCustomerResponse(customer *types.Customer) *types.CustomerResponseModel {
	return &types.CustomerResponseModel{
		FirstName:      customer.FirstName,
		LastName:       customer.LastName,
		Email:          customer.Email,
		Phone:          customer.Phone,
		Address:        customer.Address,
		AdditionalInfo: customer.AdditionalInfo,
		ContactOption:  customer.ContactOption,
	}

}

func ValidateAge(r *types.CustomerRequestModel) error {
	age := r.Age
	if age == 0 {
		return errors.New("age is required")
	}

	if age < 18 {
		return errors.New("age must be 18 or older")
	}

	return nil
}

func ValidateCustomer(customer *types.CustomerRequestModel, validate *validator.Validate) error {
	validationErrors := make(map[string]string)

	if err := ValidateAge(customer); err != nil {
		validationErrors["Age"] = err.Error()
	}

	if err := validate.Struct(customer); err != nil {
		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range fieldErrors {
				switch fieldError.Tag() {
				case "email":
					validationErrors[fieldError.Field()] = "It is not a valid email address"
				case "age":
					validationErrors[fieldError.Field()] = "Age must be a number greater than or equal to 18"
				case "required":
					validationErrors[fieldError.Field()] = "This field is required"
				default:
					validationErrors[fieldError.Field()] = "Validation failed"
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return &ValidationError{Errors: validationErrors}
	}

	return nil
}

// Custom validation error structure
type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "Validation failed"
}
