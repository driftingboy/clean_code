package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
ApiRequest
- 生成新的url，baseurl、token、accid、ts， accid用于获取pwd，ts校验时间窗口
- 从url中解出参数，baseurl、token、accid、ts
*/
const (
	splitFlagParam = "&"
	splitFlagUrl   = "?"
)

var (
	ErrParamsNotEnough = errors.New("params not enough")
)

type ApiRequest struct {
	baseUrl string
	accId   string
	token   Token
	ts      int64
}

func NewApiRequest(baseUrl, accId string, token Token, ts int64) *ApiRequest {
	return &ApiRequest{
		baseUrl: baseUrl,
		accId:   accId,
		token:   token,
		ts:      ts,
	}
}

// TODO unittest 进行优化
// TODO 性能优化
func NewApiRequestWithUrl(fullUrl string) (*ApiRequest, error) {
	// check url formal
	urlEle := strings.Split(fullUrl, splitFlagUrl)
	if len(urlEle) < 2 {
		return nil, fmt.Errorf("invalid url, has no auth info")
	}
	urlParams := strings.Split(urlEle[1], splitFlagParam)
	urlParamsLen := len(urlParams)
	if urlParamsLen < 3 {
		return nil, fmt.Errorf("params not enough, need > %d but %d", 3, urlParamsLen)
	}
	// get baseUrl
	baseUrlEle := make([]string, 3)
	baseUrlEle = append(baseUrlEle, urlEle[0])
	for _, ele := range urlParams[:urlParamsLen-3] {
		baseUrlEle = append(baseUrlEle, ele)
	}
	// get auth params
	accId := strings.Split(urlParams[urlParamsLen-3], "=")[1]
	token := Hex2Token(strings.Split(urlParams[urlParamsLen-2], "=")[1])
	tsStr := strings.Split(urlParams[urlParamsLen-1], "=")[1]
	ts, err := strconv.ParseInt(tsStr, 10, 0)
	if err != nil {
		return nil, err
	}
	return &ApiRequest{
		baseUrl: strings.Join(baseUrlEle, ""),
		accId:   accId,
		token:   token,
		ts:      ts,
	}, nil
}

func (ar ApiRequest) GetBaseUrl() string {
	return ar.baseUrl
}

func (ar ApiRequest) GetAccId() string {
	return ar.accId
}

func (ar ApiRequest) GetToken() Token {
	return ar.token
}

func (ar ApiRequest) GetCreateTimeStamp() int64 {
	return ar.ts
}

func (ar ApiRequest) String() string {
	return ar.baseUrl + ar.accId + ar.token.Hex() + strconv.FormatInt(ar.ts, 10)
}
