package request

// 接口参考：https://www.realworld.how/docs/specs/backend-specs/endpoints#registration

/** 注册
{
  "user":{
    "username": "Jacob",
    "email": "jake@jake.jake",
    "password": "jakejake"
  }
}
*/

type UserRegistrationRequest struct {
	User UserRegistrationBody `json:"user"`
}

type UserRegistrationBody struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

// 登录
/**
{
  "user":{
    "email": "jake@jake.jake",
    "password": "jakejake"
  }
}
*/

type UserAuthorizationRequest struct {
	User UserAuthorizationBody `json:"user"`
}

type UserAuthorizationBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
