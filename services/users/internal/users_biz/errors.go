package users_biz

import "jijizhazha1024/go-mall/services/users/users"

func HandleLoginerror(msg string, code int) (*users.LoginResponse, error) {
	return &users.LoginResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, nil
}
func HandleRegistererror(msg string, code int) (*users.RegisterResponse, error) {
	return &users.RegisterResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, nil
}
func HandleGetUsererror(msg string, code int) (*users.GetUserResponse, error) {
	return &users.GetUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, nil
}
func HandleDeleteUsererror(msg string, code int) (*users.DeleteUserResponse, error) {
	return &users.DeleteUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, nil
}
func HandleUpdateUsererror(msg string, code int) (*users.UpdateUserResponse, error) {
	return &users.UpdateUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, nil
}
func HandleLogoutUsererror(msg string, code int) (*users.LogoutResponse, error) {
	return &users.LogoutResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, nil
}
