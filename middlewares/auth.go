package middlewares

import (
	"RealWorldWeb/logger"
	"RealWorldWeb/security"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(ctx *gin.Context) {
	log := logger.New(ctx)
	token := ctx.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ") //这里后面不能少掉空格
	claims, ok, err := security.VerifyJwtHS256(token)
	if err != nil || !ok {
		log.WithError(err).Infof("verify jwt failed")
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx.Set("user", claims) //缓存获取到的用户信息
	ctx.Next()
}
