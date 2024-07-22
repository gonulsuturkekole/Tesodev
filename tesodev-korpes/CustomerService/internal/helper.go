package internal

import (
	"fmt"
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

/*
func ValidateAge(r *types.CustomerRequestModel) bool{

	age := r.Age
	if age == "" {
	return nil}

}
*/

func Validation(customerRequestModel *types.CustomerRequestModel) {

	var validate *validator.Validate
	errs := validate.Var(customerRequestModel.Age, "gte=18")

	if errs != nil {
		fmt.Println(errs) // output: Key: "" aError:Field validtion for "" failed on the "email" tag
		return
	}
}
