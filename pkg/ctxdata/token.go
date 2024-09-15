package ctxdata

import "github.com/golang-jwt/jwt"

const Identify = "Donkor.Easy.Chat"

func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	//过期时间
	claims["exp"] = iat + seconds
	//当前时间
	claims["iat"] = iat
	claims[Identify] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
