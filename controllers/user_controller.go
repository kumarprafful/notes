package controllers

import (
	"encoding/json"
	"net/http"
	"notes/auth"
	"notes/formaterror"
	"notes/models"
	"notes/responses"

	"github.com/valyala/fasthttp"
)

func (s *Server) Register(ctx *fasthttp.RequestCtx) {
	user := models.User{}
	err := json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		responses.ERROR(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("register")
	if err != nil {
		responses.ERROR(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(s.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(ctx, http.StatusBadRequest, formattedError)
		return
	}
	tokens, _ := auth.GenerateTokenPair(userCreated.ID)
	responses.JSON(ctx, http.StatusCreated, tokens)
}

func (s *Server) GetUsers(ctx *fasthttp.RequestCtx) {

	user := models.User{}
	users, err := user.FindAllUsers(s.DB)
	if err != nil {
		responses.ERROR(ctx, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(ctx, http.StatusOK, users)
}

func (s *Server) GetUserByID(ctx *fasthttp.RequestCtx) {
	id, _ := ctx.QueryArgs().GetUint("id")
	u := models.User{}
	user, err := u.FindUserByID(s.DB, uint32(id))
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	responses.JSON(ctx, http.StatusOK, user)
}
