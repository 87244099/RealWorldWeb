package response

// 接口参考：https://www.realworld.how/docs/specs/backend-specs/api-response-format/#users-for-authentication

/** 注册
{
  "user": {
    "email": "jake@jake.jake",
    "token": "jwt.token.here",
    "username": "jake",
    "bio": "I work at statefarm",
    "image": null
  }
}
*/

type UserRegistrationResponse struct {
	User UserRegistrationBody `json:"user"`
}

type UserRegistrationBody struct {
	Email    string      `json:"email"`
	Token    string      `json:"token"`
	Username string      `json:"username"`
	Bio      string      `json:"bio"`
	Image    interface{} `json:"image"`
}

// 登录
/** 你没看错，登录和注册都保持了一样的格式返回
{
  "user": {
    "email": "jake@jake.jake",
    "token": "jwt.token.here",
    "username": "jake",
    "bio": "I work at statefarm",
    "image": null
  }
}
*/

type UserAuthorizationResponse struct {
	User UserAuthorizationBody `json:"user"`
}

type UserAuthorizationBody struct {
	Email    string      `json:"email"`
	Token    string      `json:"token"`
	Username string      `json:"username"`
	Bio      string      `json:"bio"`
	Image    interface{} `json:"image"`
}
