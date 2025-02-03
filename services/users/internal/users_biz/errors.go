package users_biz

import "jijizhazha1024/go-mall/services/users/users"

func HandleLoginerror(msg string, code int, err error) (*users.LoginResponse, error) {
	return &users.LoginResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, err
}
func HandleRegistererror(msg string, code int, err error) (*users.RegisterResponse, error) {
	return &users.RegisterResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, err
}
func HandleGetUsererror(msg string, code int, err error) (*users.GetUserResponse, error) {
	return &users.GetUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, err
}

func HandleDeleteUsererror(msg string, code int, err error) (*users.DeleteUserResponse, error) {
	return &users.DeleteUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, err
}

func HandleUpdateUsererror(msg string, code int, err error) (*users.UpdateUserResponse, error) {
	return &users.UpdateUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, err
}
func HandleLogoutUsererror(msg string, code int, err error) (*users.LogoutResponse, error) {
	return &users.LogoutResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
	}, err
}
