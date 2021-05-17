package responses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/valyala/fasthttp"
)

func JSON(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
	ctx.SetStatusCode(statusCode)
	err := json.NewEncoder(ctx).Encode(data)
	if err != nil {
		fmt.Fprintf(ctx, "%s", err.Error())
	}
}

func ERROR(ctx *fasthttp.RequestCtx, statusCode int, err error) {
	if err != nil {
		JSON(ctx, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(ctx, http.StatusBadRequest, nil)
}
