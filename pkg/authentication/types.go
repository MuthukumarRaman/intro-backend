package authentication

import "time"

// "github.com/sirupsen/logrus"
// "introme-api/pkg/shared/utils"

// var loggerName string = "authentication"
// var log *logrus.Entry = utils.GetLogEntry(loggerName)

// LoginRequest
type LoginRequest struct {
	Id       string `json:"id" validate:"required"`
	Password string `json:"pwd" validate:"required"`
}

// LoginResponse - for Login Response
// type LoginResponse struct {
// 	Name        string      `json:"name"`
// 	UserRole    interface{}      `json:"role"`
// 	UserOrg     interface{} `json:"org" bson:"org"`
// 	Token       string      `json:"token"`
// }
type LoginResponse struct {
	Name        string      `json:"name"`
	UserRole    string      `json:"role"`
	UserProfile interface{} `json:"profile" bson:"profile"`
	Token       string      `json:"token"`
	Status      int         `json:"status" bson:"status"`
	Points      float64     `json:"points" bson:"points"`
}

// ResetPasswordRequestDto - Dto for reset password Request
type ResetPasswordRequest struct {
	Id     string `json:"id,omitempty"`
	OldPwd string `json:"old_pwd,omitempty" bson:"old_pwd,omitempty"`
	NewPwd string `json:"new_pwd" bson:"new_pwd" validate:"required"`
}

// OTPGenerateResponse - To return AuthKey
type OTP struct {
	AuthKey string `json:"auth_key"`
	Otp     int    `json:"otp"`
}

type OTPResponse struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

type Organization struct {
	Id        string      `json:"_id" bson:"_id"`
	Name      string      `json:"name" bson:"name"`
	Type      string      `json:"type" bson:"type"`
	SubDomain string      `json:"sub_domain" bson:"sub_domain"`
	Style     interface{} `json:"style" bson:"style"`
	Logo      string      `json:"logo" bson:"logo"`
	Group     string      `json:"group" bson:"group"`
	AppName   string      `json:"app_name" bson:"app_name"`
	LocOption bool        `json:"loc" bson:"loc"`
}

type UserRegister struct {
	Id           string    `json:"_id" bson:"_id"`
	Role         string    `json:"role" bson:"role"`
	EmailId      string    `json:"email_id" bson:"email_id"`
	FirstName    string    `json:"first_name" bson:"first_name"`
	LastName     string    `json:"last_name" bson:"last_name"`
	MobileNumber string    `json:"mobile_number" bson:"mobile_number"`
	Status       string    `json:"status" bson:"status"`
	CreatedOn    time.Time `json:"created_on" bson:"created_on"`
	Pwd          string    `json:"pwd" bson:"pwd"`
}

type SSOUser struct {
	Id            string    `json:"_id" bson:"_id"`
	EmailId       string    `json:"email_id" bson:"email_id"`
	FirstName     string    `json:"first_name" bson:"first_name"`
	LastName      string    `json:"last_name" bson:"last_name"`
	ProfileImage  string    `json:"profile_image" bson:"profile_image"`
	EmailVerified bool      `json:"email_verified" bson:"email_verified"`
	ProvideBy     string    `json:"provide_by" bson:"provide_by"`
	MobileNumber  string    `json:"mobile_number" bson:"mobile_number"`
	Status        string    `json:"status" bson:"status"`
	CreatedOn     time.Time `json:"created_on" bson:"created_on"`
	Role          string    `json:"role" bson:"role"`
}

type Wallet struct {
	ID                string    `json:"_id" bson:"_id"`
	User_ID           string    `json:"user_id" bson:"user_id"`
	Available_Credits float64   `json:"available_credits" bson:"available_credits"`
	Status            string    `json:"status" bson:"status"`
	CreatedON         time.Time `json:"created_on" bson:"created_on"`
	UpdatedON         time.Time `json:"updated_on" bson:"updated_on"`
}
type Transaction struct {
	ID                string    `json:"_id" bson:"_id"`
	WalletId          string    `json:"wallet_id" bson:"wallet_id"`
	Type              string    `json:"type" bson:"type"`
	OpenAmount        float64   `json:"open_amount" bson:"open_amount"`
	TransactionAmount float64   `json:"transaction_amount" bson:"transaction_amount"`
	AvailableAmount   float64   `json:"available_amount" bson:"available_amount"`
	TransactionTime   time.Time `json:"transaction_datetime" bson:"transaction_datetime"`
}
