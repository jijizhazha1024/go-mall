package code

// 用户服务
const (
	UserCreated             = 20001
	UserCreationFailed      = 20002
	UserAlreadyExists       = 20003
	EmailAlreadyExists      = 20004
	LoginSuccess            = 20005
	LoginFailed             = 20006
	InvalidCredentials      = 20007
	LogoutSuccess           = 20008
	LogoutFailed            = 20009
	UserDeleted             = 20010
	UserDeletionFailed      = 20011
	UserUpdated             = 20012
	UserUpdateFailed        = 20013
	UserInfoRetrieved       = 20014
	UserInfoRetrievalFailed = 20015
	UserNotFound            = 20016
	UserHaveDeleted         = 20017
)
const (
	UserCreatedMsg             = "用户创建成功"
	UserCreationFailedMsg      = "用户创建失败"
	UserAlreadyExistsMsg       = "用户已存在"
	EmailAlreadyExistsMsg      = "邮箱已存在"
	LoginSuccessMsg            = "登录成功"
	LoginFailedMsg             = "登录失败"
	InvalidCredentialsMsg      = "无效的凭证"
	LogoutSuccessMsg           = "登出成功"
	LogoutFailedMsg            = "登出失败"
	UserDeletedMsg             = "用户删除成功"
	UserDeletionFailedMsg      = "用户删除失败"
	UserUpdatedMsg             = "用户信息更新成功"
	UserUpdateFailedMsg        = "用户信息更新失败"
	UserInfoRetrievedMsg       = "用户身份信息获取成功"
	UserInfoRetrievalFailedMsg = "用户身份信息获取失败"
	UserNotFoundMsg            = "用户不存在"
	UserHaveDeletedMsg         = "用户已删除"
)
