package security

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

const salt = "addSalt"

func HashPassword(password string) (string, error) {
	password += salt
	//deal with password
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10) //加盐的强度
	if err != nil {
		return "", err
	}

	//TODO 为什么存储的时候最好密码base64一次？
	return base64.StdEncoding.EncodeToString([]byte(bcryptPassword)), nil
}

func CheckPassword(plain, hash string) bool {
	plain += salt
	bcryptedPassword, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(bcryptedPassword, []byte(plain))
	return err == nil
}
