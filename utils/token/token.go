package token

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"vision/config"
	"vision/types"

	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

func Generate(userID *types.ID) (string, error) {
	jwtConfig := config.GetConfig().JWT

	//key, err := jwt.ParseRSAPrivateKeyFromPEM(jwtConfig.GetSecret())
	//if err != nil {
	//	return "", err
	//}

	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"exp":        time.Now().Add(time.Hour * 24 * time.Duration(jwtConfig.Expiry)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.Secret))
}

func Extract(ctx *gin.Context) string {
	token := ctx.Query("token")
	if token != "" {
		return token
	}
	bearerToken := ctx.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func getToken(ctx *gin.Context) (*jwt.Token, error) {
	jwtConfig := config.GetConfig().JWT

	tokenString := Extract(ctx)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.Secret), nil
	})
}

func ExtractUserID(ctx *gin.Context) (userID *types.ID, err error) {
	token, err := getToken(ctx)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		strUserID, ok := claims["user_id"].(string)
		if !ok {
			err = errors.New("error while typecast userID")
			return
		}
		userID, err = types.ParseID(strUserID)
	}
	return
}

func Validate(ctx *gin.Context) error {
	_, err := getToken(ctx)
	return err
}
