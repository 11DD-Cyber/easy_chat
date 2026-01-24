package ctxdata

import (
	"github.com/golang-jwt/jwt"
)

const Identify = "gougouxuegao"

func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[Identify] = uid
	// 使用 HS256 算法创建 Token
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	// 用密钥签名，生成最终 Token 字符串
	return token.SignedString([]byte(secretKey))
}
