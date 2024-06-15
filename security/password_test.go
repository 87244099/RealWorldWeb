package security

import "testing"

func TestHashPassword(t *testing.T) {
	hashPassword, err := HashPassword("123456") // JDJhJDEwJDhadi90VGpSeFBnRC9nSWUyWkRGanVOSnFabDBhSTRrLy96Q1BDd2dwQkhSQzVUdXFxWkQy
	if err != nil {
		t.Errorf("hash password failed, err:%v", err)
	}

	t.Logf("hash password:%v", hashPassword)

	check := CheckPassword("123456", hashPassword)
	if !check {
		t.Errorf("hash password failed")
		return
	}
	t.Logf("check password:%v", check)
}

func TestCheckPassword(t *testing.T) {

}
