package users_biz

import "jijizhazha1024/go-mall/services/users/users"

func HandleLoginResp(msg string, code int, user_id uint32, token string, user_name string) (*users.LoginResponse, error) {
	return &users.LoginResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,
		Token:      token,
		UserName:   user_name,
	}, nil
}
func HandleRegisterResp(msg string, code int, user_id uint32, token string) (*users.RegisterResponse, error) {
	return &users.RegisterResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,
		Token:      token,
	}, nil
}
func HandleGetUserResp(msg string, code int, user_id uint32, user_name string) (*users.GetUserResponse, error) {
	return &users.GetUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,
		UserName:   user_name,
	}, nil
}
func HandleDeleteUserResp(msg string, code int, user_id uint32) (*users.DeleteUserResponse, error) {
	return &users.DeleteUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,
	}, nil
}
func HandleUpdateUserResp(msg string, code int, user_id uint32, token string) (*users.UpdateUserResponse, error) {
	return &users.UpdateUserResponse{
		StatusCode: uint32(code),
		StatusMsg:  msg,
		UserId:     user_id,
		Token:      token,
	}, nil
}
