
## 日志
- logrus: 最常用的日志库
- zap: 性能比较高

### logrus 使用指南
- docs: https://github.com/sirupsen/logrus
- install:
  - 如果文档里面没有，参考install go get
  - `go get github.com/sirupsen/logrus`
- usages
  - create instance by gin context: `log := logger.New(ctx)`
  - output with error: `log.WithError(err).Errorln("build json failed")`
  - output with custom info: `log.WithField("user", utils.JsonMarshal(body)).Infof("user login called")`
  - config custom log formater: `logrus.SetFormatter(&logrus.JSONFormatter{})`

## token生成

### 对称加密

#### jwt
- docs: https://github.com/golang-jwt/jwt
- install: `go get -u github.com/golang-jwt/jwt/v5`

### 非堆成加密

#### 生成密钥
- 电脑看是否有`openssl`
- 生成私钥:`openssl genrsa -out private.pem 2048`，当前路径下会出现一个`private.pem`的文件
- 由私钥生成公钥：`openssl rsa -in private.pem --outform PEM -pubout -out public.pem`，目录下会出现`public.pem`
#### 验证
- usage
```go
func VerifyJwt(token string) (bool, error) {
	var claim jwt.MapClaims
	claims, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {// 比如过期，内容被串改等都有提示
		return false, err
	}
	if claims.Valid {//最终看这个值处理
		return true, nil
	}
	return false, nil
}
```
- error type: 
  - 内容不对： `token signature is invalid: crypto/rsa: verification error`
  - 过期了：`token has invalid claims: token is expired`
- QA:
  - Q: 为什么jwt中间那一段用base64？
#### 生成token



## 配置文件
- background: 
  - 一些配置信息可能存放在go文件，配置文件或环境变量
  - 如果要统一这些方式，需要第三方库解决

### viper
- docs: https://github.com/spf13/viper
- usages
  - init
  
    ```go
     viper.SetConfigName("config") // name of config file (without extension)
  	  viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
  	  viper.AddConfigPath("..")     // optionally look for config in the working directory
  	  err := viper.ReadInConfig()   // Find and read the config file
  	  if err != nil {               // Handle errors reading the config file
  	  	panic(fmt.Errorf("fatal error config file: %w", err))
  	  }
    ```

  - read config: `viper.Get("secret")`
  - read struct: `viper.Unmarshal(&_config)`


