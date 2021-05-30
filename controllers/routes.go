package controllers

import (
	"notes/middlewares"

	"github.com/fasthttp/router"
)

func (s Server) Router() *router.Router {
	r := router.New()
	r.POST("/login/", s.Login)
	r.POST("/access_token/", s.GetAccessToken)
	r.POST("/register/", s.Register)

	r.GET("/users/", s.GetUsers)
	r.GET("/user/", s.GetUserByID)
	r.POST("/note/", middlewares.SetMiddlewareAuthentication(s.CreateANote))
	r.GET("/note/", middlewares.SetMiddlewareAuthentication(s.FetchANote))
	r.POST("/note/content/", middlewares.SetMiddlewareAuthentication(s.CreateContent))
	r.GET("/note/content/list/", middlewares.SetMiddlewareAuthentication(s.FetchContentOfNote))

	return r
}
