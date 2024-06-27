package main

import (
	"RealWorldWeb/cache"
	"RealWorldWeb/server"
)

func main() {
	cache.InitRedis()
	//从这里启动http server
	server.RunHttpServer()
}

/**
https://api.realworld.io/api/users/login
登录之后返回：
{
    "user": {
        "email": "872440899@qq.com",
        "username": "jserWang",
        "bio": null,
        "image": "https://api.realworld.io/images/smiley-cyrus.jpeg",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjozMzYwMX0sImlhdCI6MTcxODM2ODUyNCwiZXhwIjoxNzIzNTUyNTI0fQ.Vux0VKOIattyFs8iF5J0d1oYbf3upVBUIhI0Wa3tDAg"
    }
}

这里的token是一个jwt
可以通过 https://jwt.io/ 进行解密
jwt的内容解密如下：

header:
{
  "alg": "HS256", //加密算法，加密算法是hs256
  "typ": "JWT" //
}
payload:
{
  "user": {
    "id": 33601
  },
  "iat": 1718368754, //sign time
  "exp": 1723552754  //expire time
}
VERIFY SIGNATURE:
	HMACSHA256(
	  base64UrlEncode(header) + "." +
	  base64UrlEncode(payload),
		your-256-bit-secret //保存在后端的密钥
	) secret base64 encoded
*/
