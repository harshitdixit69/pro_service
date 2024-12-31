package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	ID          int32  `json:"id"`
	UserId      int32  `json:"user_id"`
	Upcoming    int32  `json:"upcoming"`
	Completed   int32  `json:"completed"`
	Cancelled   int32  `json:"cancelled"`
	Username    string `json:"username"`
	PhoneNo     string `json:"phone_no"`
	ServiceType string `json:"service_type"`
	Date        string `json:"date"`
	Time        string `json:"time"`
}
type ServiceProvider struct {
	ServiceID         int32   `json:"service_id"`
	ServiceProviderID int32   `json:"service_provider_id"`
	Username          string  `json:"username"`
	Phone             string  `json:"phone_no"`
	Lat               float64 `json:"lat"`
	Lon               float64 `json:"lon"`
	Date              string  `json:"date"`
	Time              string  `json:"time"`
}
type Service struct {
	ID          int32  `json:"id"`
	ServiceType string `json:"service_type"`
}
type DateTime struct {
	ServiceID   int32    `json:"service_id"`
	ServiceType string   `json:"service_type"`
	Date        []string `json:"date"`
	Time        string   `json:"time"`
}
type Availability struct {
	ID                int32 `json:"id"`
	ServiceID         int32 `json:"service_id"`
	ServiceProviderID int32 `json:"service_provider_id"`
}
type TypeAvailability struct {
	ID                int32  `json:"id"`
	ServiceID         int32  `json:"service_id"`
	ServiceProviderID int32  `json:"service_provider_id"`
	Date              string `json:"date"`
	Time              string `json:"time"`
}
type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PhoneNo  string `json:"phone_no"`
	IsUser   int64  `json:"is_user"`
}
type User struct {
	ID                int32     `json:"id"`
	ServiceId         int64     `json:"service_id"`
	Username          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	PrePassword       string    `json:"pre_password"`
	Password          string    `json:"password"`
	EncryptedPasswrod string    `json:"encrypted_passwrod"`
	Address           string    `json:"address"`
	PhoneNo           string    `json:"phone_no"`
	AadhaarNo         string    `json:"aadhaar_no"`
	Status            int32     `json:"status"`
	IsActive          int32     `json:"is_active"`
	IsUser            int64     `json:"is_user"`
	CreatedDate       time.Time `json:"created_date"`
	UpdatedDate       time.Time `json:"updated_date"`
	DeletedDate       time.Time `json:"deleted_date"`
}
type ForgotPass struct {
	Email_Mobile_No string `json:"email_mobile_no"`
	Password        string `json:"password"`
	IsUser          int64  `json:"is_user"`
}
type AccountInfo struct {
	Email_Mobile_No string `json:"email_mobile_no"`
	Code            string `json:"code"`
	IsUser          int64  `json:"is_user"`
}

type Users struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Location string             `json:"location,omitempty" validate:"required"`
	Title    string             `json:"title,omitempty" validate:"required"`
}
type UserResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
