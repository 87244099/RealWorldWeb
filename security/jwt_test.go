package security

import (
	"RealWorldWeb/utils"
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

func TestVerifyJwtRS256(t *testing.T) {
	rightJwt := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg1MTQ5NjEsImlhdCI6MTcxODQyODU2MSwidXNlciI6eyJlbWFpbCI6ImphY2tAZW1haWwuY29tIiwidXNlcm5hbWUiOiJqYWNrIn19.j8_FT4pnDYWDVGyk6mYBCvzFiF_nuveljg-bfHirPQnzzG1nXhgow9pXmNECQjth3HbzJg4xbER3YsZuZwQhxjEz_T1td6kcsFjTXk15x53XYmDLqIHVJtEonK4tHn9xQNLNuYMKeaLFg_lEoeWnFWSChEzp54uDkksPFxmVc0NAMzTrsDchPjBbz00AujY3O28cp82Gau_hy_-lHfhNwc1mG8Kfh-aj86AEnv9XsaC-hwwDXBsOvZRQRCBcMIWRTf7UXk9L_5_mbxaz42VES3ymFmiNb9vGT9AsRUUXKrVyu60NEdCe4_mFF162ZOB_hcbCu2qUBhSjnUD4Cme1sQ"

	claim, valid, err := VerifyJwtRS256(rightJwt)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("veriify right jwt: %v, claim: %v\n", valid, utils.JsonMarshal(claim)) //{"exp":1718514961,"iat":1718428561,"user":{"email":"jack@email.com","username":"jack"}}

	//jwtWithWrongBase64(t) //token signature is invalid: crypto/rsa: verification error
	//jwtWithExpire(t) //token has invalid claims: token is expired
}

func TestVerifyJwtHS256(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg1MTQyMTQsImlhdCI6MTcxODQyNzgxNCwidXNlciI6eyJlbWFpbCI6ImphY2tAZW1haWwuY29tIiwidXNlcm5hbWUiOiJqYWNrIn19.JFfSukmM4xH4pKzu1J2nIh57idCWSTogITIfOpHGCzg"
	claim, valid, err := VerifyJwtHS256(token)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("veriify hs256 jwt: %v, claim: %v\n", valid, utils.JsonMarshal(claim)) // {"exp":1718514214,"iat":1718427814,"user":{"email":"jack@email.com","username":"jack"}}

}

func jwtWithExpire(t *testing.T) {
	jwt := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTg1MTQyMTQsImlhdCI6MTcxODQyNzgxNCwidXNlciI6eyJlbWFpbCI6ImphY2tAZW1haWwuY29tIiwidXNlcm5hbWUiOiJqYWNrIn19.PJ-EvYL9PYuYUd24oANypxNfdl0GaYePNqV35qff-yv_6UuN1xDRGstmHSXLTLtIv3W8H4uLMfwb88S38ytYq6LyBmHJrJsReH4v3mGQoLbeZLJ3SmSXq_S-LHwBSCSOvnJ6gC5k5ztqOwN_AM92oK15WPedcH4CQ9hkPrp0JG4Fmf2704z5L7muzMVoR4sruUyrD19PeH3tod2nF9UuxRjET3wIk2OXfUoodmVc6ILa5nQXCRj3-d1R5Km95mOUU4EttqBUfvDU4K__z8ohXRVw4Bk9oXeNoHz2BXsK4kS3ghfg3ONoG2ksTbyT-2QZ2EhLoc7JZTXBQm8qTgB1Eg"

	claim, valid, err := VerifyJwtRS256(jwt)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("veriify wrong jwt: %v, claim: %v\n", valid, utils.JsonMarshal(claim))

}

func jwtWithWrongBase64(t *testing.T) {
	//jwt中间的那一段，是base64编码的结果，你可以解码，篡改，加密生成一个非法jwt，比如从{"exp":1718514961,"iat":1718428561,"user":{"email":"jack@email.com","username":"jack"}}转换成{"exp":1718514961,"iat":1718428561,"user":{"email":"jser@email.com","username":"jser"}}
	wrongJwt := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg1MTQ5NjEsImlhdCI6MTcxODQyODU2MSwidXNlciI6eyJlbWFpbCI6ImpzZXJAZW1haWwuY29tIiwidXNlcm5hbWUiOiJqc2VyIn19.j8_FT4pnDYWDVGyk6mYBCvzFiF_nuveljg-bfHirPQnzzG1nXhgow9pXmNECQjth3HbzJg4xbER3YsZuZwQhxjEz_T1td6kcsFjTXk15x53XYmDLqIHVJtEonK4tHn9xQNLNuYMKeaLFg_lEoeWnFWSChEzp54uDkksPFxmVc0NAMzTrsDchPjBbz00AujY3O28cp82Gau_hy_-lHfhNwc1mG8Kfh-aj86AEnv9XsaC-hwwDXBsOvZRQRCBcMIWRTf7UXk9L_5_mbxaz42VES3ymFmiNb9vGT9AsRUUXKrVyu60NEdCe4_mFF162ZOB_hcbCu2qUBhSjnUD4Cme1sQ"

	claim, valid, err := VerifyJwtRS256(wrongJwt)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("veriify wrong jwt: %v, claim: %v\n", valid, utils.JsonMarshal(claim))
}
