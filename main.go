package main

import (
	"log"
	"notes/controllers"
	"notes/middlewares"

	"github.com/valyala/fasthttp"
)

func main() {
	s := controllers.Server{}
	s.DBConnect()
	r := s.Router()

	log.Fatal(fasthttp.ListenAndServe(":3000", middlewares.CORS(r.Handler)))
}
