package handlers

import (
	"RealWorldWeb/config"
	"RealWorldWeb/logger"
	"RealWorldWeb/models"
	"RealWorldWeb/params/request"
	"RealWorldWeb/params/response"
	"RealWorldWeb/security"
	"RealWorldWeb/storage"
	"RealWorldWeb/utils"
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
	log := logger.New(ctx) // create log by ctx

	var body request.UserRegistrationBody
	/*
		其他的创建方式：
			body := new(request.UserRegistrationBody)
			body := &request.UserRegistrationRequest{}
	*/

	err := ctx.ShouldBind(&body) //write data into body
	if err != nil {
		log.WithError(err).Errorln("build json failed") //output log with error object
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.WithField("user", utils.JsonMarshal(body)).Infof("user registration called") //output log with custom info

	err = storage.CreateUser(ctx, &models.User{
		Username: body.User.Username,
		Password: body.User.Password,
		Email:    body.User.Email,
		Image:    config.GetDefaultPortrait(),
		Bio:      "",
	})
	if err != nil {
		log.WithError(err).Errorf("create user failed")
		//ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{
		//	"code": 0,
		//	"msg": ""
		//})
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := security.GenerateJWTByHS256(body.User.Username, body.User.Email)
	if err != nil {
		log.WithError(err).Errorln("generate jwt failed")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//todo: insert new user into db
	//todo: get new user from db

	ctx.JSON(http.StatusOK, response.UserRegistrationResponse{
		User: response.UserRegistrationBody{
			Email:    body.User.Email,
			Token:    token,
			Username: body.User.Username,
			Bio:      "",
			Image:    nil,
		},
	})
	return
}

func userLogin(ctx *gin.Context) {
	log := logger.New(ctx)
	var body = request.UserAuthorizationRequest{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		log.WithError(err).Errorln("build json failed")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	dbUser, err := storage.GetUserByEmail(ctx, body.User.Email)
	if err != nil {
		log.WithError(err).Errorf("get user failed")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if dbUser.Password != body.User.Password {
		log.WithError(err).Errorf("password error origin=%v, received=%v", dbUser.Password, body.User.Password)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	username := dbUser.Username
	token, err := security.GenerateJWTByHS256(username, body.User.Email)
	if err != nil {
		log.WithError(err).Errorln("generate jwt failed")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resUser := response.UserAuthorizationResponse{
		User: response.UserAuthorizationBody{
			Email:    body.User.Email,
			Token:    token,
			Username: username,
			Bio:      dbUser.Bio,
			Image:    dbUser.Image,
		},
	}
	log.WithField("resUser", utils.JsonMarshal(resUser))
	ctx.JSON(http.StatusOK, &resUser)

	log.WithField("user", utils.JsonMarshal(body)).Infof("user login called")
}
