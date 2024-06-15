
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
- 
