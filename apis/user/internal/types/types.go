// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package types

type AddAddressRequest struct {
	RecipientName   string `json:"recipient_name"`
	PhoneNumber     string `json:"phone_number"`
	Province        string `json:"province"`
	City            string `json:"city"`
	DetailedAddress string `json:"detailed_address"`
	IsDefault       bool   `json:"is_default"`
}

type AddAddressResponse struct {
	Data AddressData `json:"data"`
}

type AddressData struct {
	AddressID       int32  `json:"address_id"`
	RecipientName   string `json:"recipient_name"`
	PhoneNumber     string `json:"phone_number"`
	Province        string `json:"province"`
	City            string `json:"city"`
	DetailedAddress string `json:"detailed_address"`
	IsDefault       bool   `json:"is_default"`
	CreatedAt       string `json:"created_at"` // ISO8601 时间转换为时间戳
	UpdatedAt       string `json:"updated_at"` // ISO8601 时间转换为时间戳
}

type AddressListResponse struct {
	Data []AddressData `json:"data"`
}

type AddressResponse struct {
	ID              string `json:"id"`
	RecipientName   string `json:"recipient_name"`
	PhoneNumber     string `json:"phone_number"`
	Province        string `json:"province"`
	City            string `json:"city"`
	DetailedAddress string `json:"detailed_address"`
	IsDefault       string `json:"is_default"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type AllAddressListRequest struct {
}

type DeleteAddressRequest struct {
	AddressID int32 `json:"address_id"`
}

type DeleteAddressResponse struct {
}

type DeleteRequest struct {
}

type DeleteResponse struct {
}

type GetAddressRequest struct {
	AddressID int32 `json:"address_id"`
}

type GetAddressResponse struct {
	Data AddressData `json:"data"`
}

type GetInfoRequest struct {
}

type GetInfoResponse struct {
	UserId    int64  `json:"user_id"`
	LogoutAt  string `json:"logout_at"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Avatar    string `json:"avatar"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
}

type LogoutResponse struct {
	Logout_at int64 `json:"logout_at"`
}

type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateAddressRequest struct {
	RecipientName   string `json:"recipient_name"`
	PhoneNumber     string `json:"phone_number"`
	Province        string `json:"province"`
	City            string `json:"city"`
	DetailedAddress string `json:"detailed_address"`
	IsDefault       bool   `json:"is_default"`
	AddressID       int32  `json:"address_id"`
}

type UpdateAddressResponse struct {
	Data AddressData `json:"data"`
}

type UpdateRequest struct {
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
}

type UpdateResponse struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
}
