package auth

import (
	"time"

	"github.com/pkg/errors"
)

type ApiAuthenticator interface {
	Auth(fullUrl string)
	AuthApi(api ApiRequest)
}

type DefaultAuthenticator struct {
	storage Storage
}

func NewDefaultAuthenticator(storage Storage) *DefaultAuthenticator {
	return &DefaultAuthenticator{
		storage: storage,
	}
}

func (d *DefaultAuthenticator) Auth(fullUrl string) error {
	api, err := NewApiRequestWithUrl(fullUrl)
	if err != nil {
		return errors.Wrap(err, "Auth NewApiRequestWithUrl")
	}
	return d.AuthApi(api)
}

func (d *DefaultAuthenticator) AuthApi(api *ApiRequest) error {
	// 1 参数准备
	token := api.GetToken()
	baseUrl := api.GetBaseUrl()
	createTime := api.GetCreateTimeStamp()
	accId := api.GetAccId()
	pwd := d.storage.QueryPwdByAccId(accId)

	// 2 校验
	clientAuthToken := NewAuthToken(token, time.Unix(createTime, 0))
	if clientAuthToken.IsExpired() {
		return errors.New("token expired!")
	}
	serverAuthToken := GenerateToken(baseUrl, accId, pwd, clientAuthToken.creatTime.Unix())
	if !clientAuthToken.Match(serverAuthToken) {
		return errors.New("token not match!")
	}
	return nil
}

// - 把 URL、AppID、密码、时间戳拼接为一个字符串
// - 对字符串通过加密算法加密生成 token
// - 根据时间戳判断 token 是否过期失效
// - 验证两个 token 是否匹配。

// func authMap(urlParams []string) (authMap map[string]string, err error) {
// 	urlParamsLen := len(urlParams)
// 	if urlParamsLen < 3 {
// 		return nil, fmt.Errorf("params not enough, need > %d but %d", 3, urlParamsLen)
// 	}
// 	authMap = map[string]string{
// 		"accId": "",
// 		"token": "0x",
// 		"ts":    "0",
// 	}
// 	for _, ele := range urlParams {
// 		kv := strings.Split(ele, "=")
// 		if _, ok := authMap[kv[0]]; ok {
// 			authMap[kv[0]] = kv[1]
// 		}
// 	}
// 	return
// }
