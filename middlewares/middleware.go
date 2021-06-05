package middlewares

import (
	"errors"
	"net/http"
	"notes/auth"
	"notes/responses"

	"github.com/valyala/fasthttp"
)

// func SetMiddlewareJSON(next fasthttp.RequestHandler) fasthttp.RequestHandler {
// 	return func(ctx *fasthttp.RequestCtx) {
// 		ctx.Response.Header.Set("Content-Type", "application/json")
// 		next(ctx)
// 	}
// }

func SetMiddlewareAuthentication(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		err := auth.TokenValid(ctx)
		if err != nil {
			responses.ERROR(ctx, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(ctx)
	}
}

var (
	corsAllowCredentials = "true"
	corsAllowHeaders     = "authorization"
	corsAllowOrigin      = "*"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
)

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
		next(ctx)
	}
}
