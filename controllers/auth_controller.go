package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"notes/auth"
	"notes/formaterror"
	"notes/models"
	"notes/responses"

	"github.com/valyala/fasthttp"
)

// type Block struct {
// 	Try     func()
// 	Catch   func(Exception)
// 	Finally func()
// }

// type Exception interface {
// 	error
// }

// func Throw(e Exception) {
// 	panic(e)
// }

// func (tcf Block) Do() {
// 	if tcf.Finally != nil {

// 		defer tcf.Finally()
// 	}
// 	if tcf.Catch != nil {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				tcf.Catch(r)
// 			}
// 		}()
// 	}
// 	tcf.Try()
// }

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *Login) Validate() error {
	if l.Email == "" {
		return errors.New("email is required")
	}
	if l.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (s *Server) Login(ctx *fasthttp.RequestCtx) {
	userLogin := Login{}
	user := models.User{}
	err := json.Unmarshal(ctx.Request.Body(), &userLogin)
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return

	}
	err = userLogin.Validate()
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return

	}
	userFound, err := user.FindUserByEmail(s.DB, userLogin.Email)
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return

	}
	fmt.Println(userFound.Password, user.Password)
	err = models.VerifyPassword(userFound.Password, userLogin.Password)
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return

	}
	tokens, err := auth.GenerateTokenPair(userFound.ID)
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	responses.JSON(ctx, http.StatusOK, tokens)
}

func (s *Server) GetAccessToken(ctx *fasthttp.RequestCtx) {
	type Token struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenBody := Token{}
	err := json.Unmarshal(ctx.Request.Body(), &tokenBody)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(ctx, http.StatusBadRequest, formattedError)
		return
	}
	access_token, err := auth.GenerateAccessFromRefreshToken(tokenBody.RefreshToken)
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	responses.JSON(ctx, http.StatusOK, access_token)
}
