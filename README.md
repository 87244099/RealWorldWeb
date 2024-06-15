
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

### jwt
- docs: https://github.com/golang-jwt/jwt
- install: `go get -u github.com/golang-jwt/jwt/v5`

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
