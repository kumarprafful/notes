package controllers

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func (s Server) Router() *router.Router {
	r := router.New()
	r.GET("/", func(ctx *fasthttp.RequestCtx) { fmt.Print("as") })
	r.POST("/register/", s.Register)
	r.GET("/users/", s.GetUsers)

	return r
}
