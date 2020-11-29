package middleware

import (
	"amazon-scrapper/lib/httperrors"
	"amazon-scrapper/lib/web"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

const (
	ApiVersionV1 = 1
)

type ResponseBuilder func(data *web.JSONResponse, responseErr *httperrors.HttpError) *web.JSONResponse

func ServeEndpoint(nextHandler func(request *http.Request) (*web.JSONResponse, *httperrors.HttpError)) httprouter.Handle {
	return serveEndpoint(buildResponseBuilder(ApiVersionV1), nextHandler)
}

func serveEndpoint(responseBuilder ResponseBuilder, nextHandler func(request *http.Request) (*web.JSONResponse, *httperrors.HttpError)) httprouter.Handle {
	return func(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		startTime := time.Now()
		defer func() {
			if recvr := recover(); recvr != nil {
				responseCode := http.StatusInternalServerError
				w.WriteHeader(http.StatusInternalServerError)
				errorMessage := fmt.Sprintf("%v", recvr)
				writeResponse(w, responseBuilder, nil, &httperrors.HttpError{StatusCode: responseCode, Error: errors.New(errorMessage)})
			}
		}()
		setCommonHeaders(w)
		data, responseErr := nextHandler(request)
		responseCode := getResponseCode(responseErr)
		w.WriteHeader(responseCode)
		writeResponse(w, responseBuilder, data, responseErr)
		log.Info("Request processed", map[string]interface{}{
			"status":         responseCode,
			"path":           request.URL.Path,
			"method":         request.Method,
			"request_params": QueryParams(request),
			"duration_ms":    float64(time.Since(startTime).Nanoseconds()) / 1e6,
		})
	}
}

func setCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
}

func writeResponse(w http.ResponseWriter, responseBuilder ResponseBuilder, data *web.JSONResponse, responseErr *httperrors.HttpError) {
	_, err := w.Write(responseBuilder(data, responseErr).ByteArray())
	if err != nil {
		log.Error("Error in writing response", err.Error())
	}
}

func buildResponseBuilder(version int) ResponseBuilder {
	return func(data *web.JSONResponse, responseErr *httperrors.HttpError) *web.JSONResponse {
		if responseErr == nil {
			return &web.JSONResponse{
				"api_version": version,
				"success":     true,
				"data":        data,
			}
		} else {
			return &web.JSONResponse{
				"api_version": version,
				"error": map[string]interface{}{
					"code":    responseErr.StatusCode,
					"message": responseErr.Error.Error(),
				},
				"success": false,
			}
		}
	}
}

func getResponseCode(err *httperrors.HttpError) int {
	if err == nil {
		return http.StatusOK
	}
	return err.StatusCode
}

func QueryParams(r *http.Request) map[string]interface{} {
	var params map[string]interface{}

	for key, val := range r.URL.Query() {
		params[key] = val
	}
	return params
}
