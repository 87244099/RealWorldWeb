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

/*
*

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
type EditUserRequest struct {
	EditUserBody EditUserBody `json:"user"`
}

type EditUserBody struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Password string `json:"password"`
}
