package storage

import (
	"RealWorldWeb/models"
	"RealWorldWeb/utils"
	"context"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	user := models.User{
		Username: "xxx123",
		Password: "xxx123",
		Email:    "xxx123@gmail.com",
		Image:    "xxx123",
		Bio:      "xxx123",
	}
	err := CreateUser(ctx, &user)
	if err != nil {
		t.Error(err)
	}
	t.Logf("create user user:%v", utils.JsonMarshal(user))
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	user, err := GetUserByEmail(ctx, "xxx123@gmail.com")
	if err != nil {
		t.Errorf("get user by email err:%v", err)
		return
	}
	t.Logf("get user user:%v", utils.JsonMarshal(user))
}

func TestDeleteUserByEmail(t *testing.T) {
	ctx := context.Background()
	err := DeleteUserByEmail(ctx, "xxx123@gmail.com")
	if err != nil {
		t.Errorf("delete user by email err:%v", err)
		return
	}
	t.Logf("delete user")
}
