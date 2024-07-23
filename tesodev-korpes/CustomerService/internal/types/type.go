package types

import "time"

// CustomerRequestModel represents a model for customer requests.
type CustomerRequestModel struct {
	FirstName string    `bson:"first_name" json:"first_name" validate:"required"`
	LastName  string    `bson:"last_name" json:"last_name" validate:"required"`
	Email     string    `bson:"email" json:"email" validate:"required,email"`
	Age       int       `bson:"age" json:"age"  validate:"required"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type CustomerResponseModel struct {
	FirstName      string            `bson:"first_name" json:"first_name"`
	LastName       string            `bson:"last_name" json:"last_name"`
	Email          string            `bson:"email" json:"email"`
	Phone          string            `bson:"phone" json:"phone"`
	Address        string            `bson:"address" json:"address"`
	AdditionalInfo map[string]string `bson:"additional_info" json:"additional_info"`
	ContactOption  []string          `bson:"contact_option" json:"contact_option"`
	Age            int               `bson:"age" json:"age"`
}

type CustomerUpdateModel struct {
	FirstName      string            `bson:"first_name" json:"first_name"`
	LastName       string            `bson:"last_name" json:"last_name"`
	Age            int               `bson:"age" json:"age"`
	Phone          string            `bson:"phone" json:"phone"`
	Address        string            `bson:"address" json:"address"`
	City           string            `bson:"city" json:"city"`
	State          string            `bson:"state" json:"state"`
	ZipCode        string            `bson:"zip_code" json:"zip_code"`
	AdditionalInfo map[string]string `bson:"additional_info" json:"additional_info"`
	MembershipType string            `bson:"membership_type" json:"membership_type"`
	ContactOption  []string          `bson:"contact_option" json:"contact_option"`
	CreatedAt      time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time         `bson:"updated_at" json:"updated_at"`
}

/*type  QueryParams interface {

	firstName string
	age_greater_than int
	age_less_than    int

}
*/
