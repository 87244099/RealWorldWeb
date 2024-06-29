package handlers

import (
	"RealWorldWeb/cache"
	"RealWorldWeb/config"
	"RealWorldWeb/logger"
	"RealWorldWeb/middlewares"
	"RealWorldWeb/models"
	"RealWorldWeb/params/request"
	"RealWorldWeb/params/response"
	"RealWorldWeb/security"
	"RealWorldWeb/storage"
	"RealWorldWeb/utils"
	"context"
	redislock "github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// go没有类似godoc的工具

// AddUserHandler r是http服务的实例
func AddUserHandler(r *gin.Engine) {
	usersGroup := r.Group("/api/users")

	usersGroup.POST("", userRegistration)
	usersGroup.POST("/login", userLogin)
	usersGroup.POST("/user-image", uploadUserImage)

	r.GET("/api/profiles/:username", userProfile)

	userGroup := r.Group("/api/user")
	userGroup.Use(middlewares.AuthMiddleware).
		PUT("", editUser)
	//r.Use(middlewares.AuthMiddleware).PUT("/api/user", editUser)
}

func uploadUserImage(ctx *gin.Context) {
	log := logger.New(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		log.WithError(err).Infof("get file failed")
	}
	err = ctx.SaveUploadedFile(file, "./"+file.Filename)
	if err != nil {
		log.WithError(err).Infof("save file failed")
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)

}

func userProfile(ctx *gin.Context) {
	log := logger.New(ctx)
	username := ctx.Param("username")
	log = log.WithField("username", username) //TODO: what is it?
	log.Infof("user profile called, username: %s", username)

	var err error
	var user *models.User
	user, err = cache.GetUserProfile(ctx, username)
	if user == nil {
		user, err = storage.GetUserByUsername(ctx, username)
		if err != nil {
			log.WithError(err).Infof("GetUserByUsername DB failed")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = cache.SetUserProfile(ctx, username, user, 300)
		if err != nil {
			log.WithError(err).Infof("SetUserProfile cache failed")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		log.Infof("GetUserByUsername from cache")

	}

	ctx.JSON(http.StatusOK, response.UserProfileResponse{
		UserProfile: response.UserProfile{
			Username:  user.Username,
			Bio:       user.Bio,
			Image:     user.Image,
			Following: false,
		},
	})
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

	hashPassword, err := security.HashPassword(body.User.Password)
	if err != nil {
		log.WithError(err).Errorln("HashPassword password failed") //output log with error object
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = storage.CreateUser(ctx, &models.User{
		Username: body.User.Username,
		Password: hashPassword,
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

	if !security.CheckPassword(body.User.Password, dbUser.Password) {
		encodedReceived, err := security.HashPassword(body.User.Password)
		if err != nil {
			log.WithError(err).Errorln("HashPassword failed")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.WithError(err).Errorf("password error origin=%v, received=%v, encodedReceived=%v", dbUser.Password, body.User.Password, encodedReceived)
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

	ctx.SetCookie("token", token, 24*3600, "/", "localhost:8000", false, true)

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

func editUser(ctx *gin.Context) {
	log := logger.New(ctx)
	log.Infof("editUser: %v", utils.JsonMarshal(ctx.MustGet("user")))

	editProfileLockKey := cache.UserEditProfileLockKey(security.GetCurrentUsername(ctx))
	lock, err2 := cache.Locker.Obtain(ctx, editProfileLockKey, 30*time.Second, redislock.Options{
		RetryStrategy: redislock.ExponentialBackoff(100*time.Millisecond, 5*time.Second),
		//Metadata:      "",
		//Token:         "",
	})
	if err2 == redislock.ErrNotObtained {
		log.WithError(err2).Infof("user is editing; username=" + security.GetCurrentUsername(ctx))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	} else if err2 != nil {
		log.WithError(err2).Errorln("obtain lock failed")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	defer func(lock *redislock.Lock, ctx context.Context) {
		err := lock.Release(ctx)
		if err != nil {
			log.WithError(err).Errorln("lock release failed")
		}
	}(lock, ctx)

	time.Sleep(5 * time.Second)

	var body = request.EditUserRequest{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		log.WithError(err).Errorln("build json failed")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if body.EditUserBody.Username == "" || body.EditUserBody.Email == "" {
		log.WithError(err).Infof("Username or email is empty; username=%s; email=%s", body.EditUserBody.Username, body.EditUserBody.Email)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if body.EditUserBody.Password != "" {
		var err error
		body.EditUserBody.Password, err = security.HashPassword(body.EditUserBody.Password)
		if err != nil {
			log.WithError(err).Errorln("HashPassword password failed")
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
	username := security.GetCurrentUsername(ctx)

	dbUser := &models.User{
		Username: body.EditUserBody.Username,
		Password: body.EditUserBody.Password,
		Email:    body.EditUserBody.Email,
		Image:    body.EditUserBody.Image,
		Bio:      body.EditUserBody.Bio,
	}
	err = storage.UpdateUserByUsername(ctx, username, dbUser)
	if err != nil {
		log.WithError(err).Errorf("UpdateUserByUsername failed")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = cache.DeleteUserProfile(ctx, username)
	if err != nil {
		log.WithError(err).Infof("DeleteUserProfile failed")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := security.GenerateJWTByHS256(body.EditUserBody.Username, dbUser.Email)
	if err != nil {
		log.WithError(err).Infof("gen jwt failed")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response.UserAuthorizationResponse{
		User: response.UserAuthorizationBody{
			Email:    body.EditUserBody.Email,
			Token:    token,
			Username: dbUser.Username,
			Bio:      dbUser.Bio,
			Image:    dbUser.Image,
		},
	})
}
