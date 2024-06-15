package security

import (
	"testing"
)

func TestGenerateJWTByHS256(t *testing.T) {
	token, err := GenerateJWTByHS256("jack", "jack@email.com")
	if err != nil {
		t.Error(err)
	}
	t.Logf(token)
}
func TestGenerateJWTByRS256(t *testing.T) {
	token, err := GenerateJWTByRS256("jack", "jack@email.com")
	if err != nil {
		t.Error(err)
	}
	t.Logf(token)
}
