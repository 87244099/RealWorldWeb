package security

import "testing"

func TestGenerateJWT(t *testing.T) {
	token, err := GenerateJWT("jack", "jack@email.com")
	if err != nil {
		t.Error(err)
	}
	t.Logf(token)
}
