package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const UserIdClaim string = "user"

// UserFromToken Validate token, if valid, return user_id, else return error description
func UserFromToken(tokenString string, privateKey string, VerifyExpiresAt bool, expirationDuration time.Duration) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(privateKey), nil
	})

	if err != nil {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// TODO: check
		if claims.VerifyExpiresAt(int64(expirationDuration.Seconds()), VerifyExpiresAt) {
			return claims[UserIdClaim].(string), nil
		}
		return "", errors.New("token expired")
	}
	return "", errors.New("invalid token")
}

// NewToken Issue a new token for user_id
func NewToken(user_id string, privateKey string, expireDuration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		UserIdClaim: user_id,
		//"nbf": expireDuration.Seconds(time.Now().Add(expireDuration).Unix()),
		"exp": time.Now().Add(expireDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}
