package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notes/formaterror"
	"notes/models"
	"notes/responses"

	"github.com/valyala/fasthttp"
)

func (s *Server) Register(ctx *fasthttp.RequestCtx) {
	fmt.Println("asdad", s)
	user := models.User{}
	err := json.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		responses.ERROR(ctx, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
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
	responses.JSON(ctx, http.StatusCreated, userCreated)
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
