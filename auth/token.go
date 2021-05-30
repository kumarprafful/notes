package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

var SECRET string = os.Getenv("SECRET")

func GenerateTokenPair(userID uint32) (map[string]string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	at, err := access_token.SignedString([]byte(SECRET))
	if err != nil {
		return nil, err
	}

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rt, err := refresh_token.SignedString([]byte(SECRET))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  at,
		"refresh_token": rt,
	}, nil
}

func GenerateAccessFromRefreshToken(token string) (map[string]string, error) {
	userID, err := ExtractIDFromToken(token)
	if err != nil {
		return nil, err
	}
	token, err = CreateToken(userID)
	if err != nil {
		return nil, err
	}
	tokens := map[string]string{
		"access_token": token,
	}
	return tokens, nil
}

func CreateToken(userID uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET))
}

func TokenValid(ctx *fasthttp.RequestCtx) error {
	tokenString := ExtractToken(ctx)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

func ExtractToken(ctx *fasthttp.RequestCtx) string {
	bearerToken := string(ctx.Request.Header.Peek("Authorization"))
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(ctx *fasthttp.RequestCtx) (uint32, error) {
	tokenString := ExtractToken(ctx)
	id, err := ExtractIDFromToken(tokenString)
	if err != nil {
		return 0, err
	}
	return id, err
}

func ExtractIDFromToken(tokenString string) (uint32, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", "")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}
