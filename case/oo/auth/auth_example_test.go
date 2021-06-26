package auth_test

import (
	"github/litao-2071/clean_code/case/oo/auth"
	"strconv"
	"testing"
)

type testAuthStorage struct{}

func (as *testAuthStorage) QueryPwdByAccId(id string) string {
	return "123456"
}
func Test_Exmaple(t *testing.T) {
	// client
	// 1 sign to get a token
	baseUrl := "http://lt.fly.com/activity/100"
	accId := "1"
	pwd := "123456"
	timestamp := 1619494512
	token := auth.GenerateToken(baseUrl, accId, pwd, int64(timestamp))
	// 2 take token with request
	fullUrl := "http://lt.fly.com/activity/100?accid=1&token=" + token.Hex() + "&ts=" + strconv.Itoa(1619494512)
	// server
	// auth
	authenticator := auth.NewDefaultAuthenticator(&testAuthStorage{})
	err := authenticator.Auth(fullUrl)
	if err != nil {
		t.Fatal(err)
	}
}
