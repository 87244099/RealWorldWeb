package handlers

import (
	"RealWorldWeb/params/request"
	"RealWorldWeb/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// go没有类似godoc的工具

// AddUserHandler r是http服务的实例
func AddUserHandler(r *gin.Engine) {
	userGroup := r.Group("/api/users")

	userGroup.POST("", userRegistration)
	userGroup.POST("/login", userLogin)
}

func userRegistration(ctx *gin.Context) {
	var body request.UserRegistrationBody
	/*
		其他的创建方式：
			body := new(request.UserRegistrationBody)
			body := &request.UserRegistrationRequest{}
	*/

	err := ctx.ShouldBind(&body) //write data into body
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fmt.Println("userRegistration-body", utils.JsonMarshal(body))

}

func userLogin(ctx *gin.Context) {
	var body = request.UserAuthorizationRequest{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println("userLogin-body", utils.JsonMarshal(body))
}
