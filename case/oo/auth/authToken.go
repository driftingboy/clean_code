package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"github/litao-2071/clean_code/case/oo/hexutil"
	"strconv"
	"time"
)

/*
AuthToken
- 把 URL、AppID、密码、时间戳拼接为一个字符串
- 对字符串通过加密算法加密生成 token
- 根据时间戳判断 token 是否过期失效
- 验证两个 token 是否匹配。
*/

const TokenLength = 32

type Token [TokenLength]byte

func Hex2Token(tokenHex string) Token {
	return BytesToToken(hexutil.FromHex(tokenHex))
}

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToToken(b []byte) Token {
	var t Token
	t.SetBytes(b)
	return t
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (t *Token) SetBytes(b []byte) {
	if len(b) > len(t) {
		b = b[len(b)-TokenLength:]
	}

	copy(t[TokenLength-len(b):], b)
}

func (t Token) String() string {
	return string(t[:])
}

func (t Token) Hex() string {
	return hex.EncodeToString(t[:])
}

var DefaultExpirationIntervalSec = 1 * 60 * 60

type AuthToken struct {
	token                 Token
	creatTime             time.Time
	expirationIntervalSec int64
}

func NewAuthToken(token Token, createTime time.Time) *AuthToken {
	return &AuthToken{
		token:                 token,
		creatTime:             createTime,
		expirationIntervalSec: int64(DefaultExpirationIntervalSec),
	}
}

func GenerateToken(baseUrl, accId, pwd string, timestamp int64) Token {
	token := baseUrl + accId + pwd + strconv.FormatInt(timestamp, 10)
	return sha256.Sum256([]byte(token))
}

func (a *AuthToken) IsExpired() bool {
	expireTime := a.creatTime.Add(time.Second * time.Duration(a.expirationIntervalSec))
	if time.Now().After(expireTime) {
		return true
	}
	return false
}

func (a *AuthToken) Match(token Token) bool { // TODO token直接使用hex， token和AuthTOken功能重复没有意义，都使用hex传输即可
	if a.token.Hex() != token.Hex() {
		return false
	}
	return true
}
