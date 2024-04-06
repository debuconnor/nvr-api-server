package roomapi

import (
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func Run() {
	r := router.New()

	r.POST("/setup", setupHandler)
	r.POST("/add", addHandler)
	r.POST("/cancel", cancelHandler)
	r.POST("/save", saveHandler)
	r.POST("/delete", deleteHandler)

	if err := fasthttp.ListenAndServe(":9090", r.Handler); err != nil {
		log.Println(err)
	}
}

func Test() {

}
