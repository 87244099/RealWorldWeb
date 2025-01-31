package storage

import (
	"RealWorldWeb/models"
	"context"
	"github.com/gin-gonic/gin"
)

func CreateUser(ctx context.Context, user *models.User) error {
	_, err := db.ExecContext(ctx, "insert user(username, password, email, image, bio) values(?,?,?,?,?)", user.Username, user.Password, user.Email, user.Image, user.Bio)

	return err
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	//一般推荐穷举字段，避免表结构变更导致报错
	err := db.GetContext(ctx, &user, "select username, password, email, image, bio from user where email=?", email)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	//一般推荐穷举字段，避免表结构变更导致报错
	err := db.GetContext(ctx, &user, "select username, password, email, image, bio from user where username=?", username)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func DeleteUserByEmail(ctx context.Context, email string) error {
	_, err := db.ExecContext(ctx, "delete from user where email=?", email)
	return err
}

func UpdateUserByUsername(ctx *gin.Context, username string, user *models.User) error {
	if user.Password != "" {
		_, err := db.ExecContext(ctx, "update user set username=?, password=?, email=?, image=?, bio=? where username=?", username, user.Password, user.Email, user.Image, user.Bio, user.Username)
		if err != nil {
			return err
		}

	}
	_, err := db.ExecContext(ctx, "update user set username=?, email=?, image=?, bio=? where username=?", username, user.Email, user.Image, user.Bio, user.Username)
	return err
}
