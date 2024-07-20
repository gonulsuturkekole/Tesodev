package types

import "time"

type CustomerRequestModel struct {
	Id        string `bson:"_id" json:"id"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email"`
	Age       string `bson:"age" json:"age"`
}

type CustomerResponseModel struct {
	FirstName      string            `bson:"first_name" json:"first_name"`
	LastName       string            `bson:"last_name" json:"last_name"`
	Email          string            `bson:"email" json:"email"`
	Phone          string            `bson:"phone" json:"phone"`
	Address        string            `bson:"address" json:"address"`
	AdditionalInfo map[string]string `bson:"additional_info" json:"additional_info"`
	ContactOption  []string          `bson:"contact_option" json:"contact_option"`
	Age            string            `bson:"age" json:"age"`
}

type CustomerUpdateModel struct {
	FirstName      string            `bson:"first_name" json:"first_name"`
	LastName       string            `bson:"last_name" json:"last_name"`
	Age            string            `bson:"age" json:"age"`
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
