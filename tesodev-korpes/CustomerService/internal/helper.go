/*
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
*/
package internal

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"strings"
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

func ValidateEmail(r *types.CustomerRequestModel) error {

	email := r.Email

	if email == "" {
		return errors.New("Email is required")
	}

	if !strings.Contains(email, "@") {
		return errors.New("Email must contain @")
	}
	return nil

}

func ValidateAge(r *types.CustomerRequestModel) error {
	age := r.Age
	if age == 0 {
		return errors.New("Age is required")
	}

	if age < 18 {
		return errors.New("Age must be 18 or older")
	}

	return nil
}
func ValidateFirstLetterUpperCase(customer *types.CustomerRequestModel) error {
	errors := make(map[string]string)
	if customer.FirstName != "" {
		if !isFirstLetterUpperCase(customer.FirstName) {
			errors["FirstName"] = "First name must start with an uppercase letter"
		}
		if containsDigit(customer.FirstName) {
			errors["FirstName"] = "First name contains a number"
		}
	}
	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}
	return nil
}

func containsDigit(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			return true
		}
	}
	return false
}

// Helper function to check if the first letter is uppercase
func isFirstLetterUpperCase(s string) bool {
	if len(s) > 0 {
		return strings.ToUpper(s[:1]) == s[:1]
	}
	return false
}
func ValidateCustomer(customer *types.CustomerRequestModel, validate *validator.Validate) error {
	validationErrors := make(map[string]string)

	if err := ValidateAge(customer); err != nil {
		validationErrors["Age"] = err.Error()
	}

	if err := ValidateEmail(customer); err != nil {
		validationErrors["Email"] = err.Error()
	}

	if err := ValidateFirstLetterUpperCase(customer); err != nil {
		// Use the errors from ValidateFirstLetterUpperCase directly
		if valErr, ok := err.(*ValidationError); ok {
			for field, msg := range valErr.Errors {
				validationErrors[field] = msg
			}
		}
	}

	if err := validate.Struct(customer); err != nil {
		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range fieldErrors {
				switch fieldError.Tag() {
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
