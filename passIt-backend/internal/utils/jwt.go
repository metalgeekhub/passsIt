package utils

import (
	"crypto/rand"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	key            []byte
	t              *jwt.Token
	s              string
	expirationTime int
)

func init() {
	key = make([]byte, 32) // HS256 expects a []byte key, 256 bits
	_, err := rand.Read(key)
	if err != nil {
		panic("failed to generate JWT key: " + err.Error())
	}
	expStr := os.Getenv("JWT_EXPIRATION_TIME")
	if expStr == "" {
		expirationTime = 60 // default to 60 minutes if not set
	} else {
		if expInt, err := strconv.Atoi(expStr); err == nil {
			expirationTime = expInt
		} else {
			expirationTime = 60 // fallback to default if conversion fails
		}
	}
}

func GenerateJWTToken(userID uuid.UUID) (string, error) {
	var err error
	claims := jwt.MapClaims{}
	claims["uer_id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(expirationTime)).Unix() // Token expires in 1 Hour

	t = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err = t.SignedString(key)
	if err != nil {
		return "", err
	}
	return s, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, jwt.ErrTokenMalformed
		}

		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
