package routes

import (
	"amazon-scrapper/lib/httperrors"
	"amazon-scrapper/lib/utils/middleware"
	"amazon-scrapper/lib/web"
	"amazon-scrapper/saver/app/controllers"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Init(router *httprouter.Router) {
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, "{ \"message\":\"Hello world!. I am Scrap saver.\",\"success\":true,\"api_version\": 1 }")
	})

	router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(404)
		fmt.Fprint(rw, "{ \"message\":\"Not Found.\",\"success\":true,\"api_version\": 1 }")
	})

	saveController := controllers.NewSaveController()
	router.POST("/api/save", ServeEndpoint(saveController.SaveProduct))
}

func ServeEndpoint(endpointHandler func(request *http.Request) (*web.JSONResponse, *httperrors.HttpError)) httprouter.Handle {
	return middleware.ServeEndpoint(endpointHandler)
}
