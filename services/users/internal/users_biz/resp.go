package users_biz

import (
	"jijizhazha1024/go-mall/services/users/users"
	"time"
)

func HandleLoginResp(msg string, code int, user_id uint32, user_name string) (*users.LoginResponse, error) {
	return &users.LoginResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,

		UserName: user_name,
	}, nil
}
func HandleRegisterResp(msg string, code int, user_id uint32) (*users.RegisterResponse, error) {
	return &users.RegisterResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,
	}, nil
}
func HandleGetUserResp(msg string, code int, user_id uint32, user_name string, email string, created_at string, updated_at string, logout_at string, avatar_url string) (*users.GetUserResponse, error) {
	return &users.GetUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,
		UserName:   user_name,
		Email:      email,
		CreatedAt:  created_at,
		UpdatedAt:  updated_at,
		LogoutAt:   logout_at,
		AvatarUrl:  avatar_url,
	}, nil
}
func HandleDeleteUserResp(msg string, code int, user_id uint32) (*users.DeleteUserResponse, error) {
	return &users.DeleteUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,
	}, nil
}
func HandleUpdateUserResp(msg string, code int, user_id uint32, email string, user_name string) (*users.UpdateUserResponse, error) {
	return &users.UpdateUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,
		Email:      email,
		UserName:   user_name,
	}, nil
}
func HandleLogoutUserResp(msg string, code int, logout_at time.Time) (*users.LogoutResponse, error) {
	return &users.LogoutResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,

		LogoutTime: logout_at.Unix(),
	}, nil
}
