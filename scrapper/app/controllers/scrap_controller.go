package controllers

import (
	"amazon-scrapper/lib/httperrors"
	"amazon-scrapper/lib/models"
	"amazon-scrapper/lib/web"
	"amazon-scrapper/scrapper/app/handlers"
	"encoding/json"
	"net/http"
)

type ScrapController struct {
	scrapHandler handlers.ScrapHandlerInterface
}

func NewScrapController() *ScrapController {
	return &ScrapController{
		scrapHandler: handlers.NewScrapHandler(),
	}
}

type ScrapControllerInterface interface {
	ScrapProduct(*http.Request) (*web.JSONResponse, *httperrors.HttpError)
}

func (sc *ScrapController) ScrapProduct(r *http.Request) (*web.JSONResponse, *httperrors.HttpError) {
	requestBodyBytes, err := web.ReadBodyBytes(r)
	if err != nil {
		return nil, httperrors.BadRequestError(err.Error())
	}

	request := models.AmazonProductRequest{}
	err = json.Unmarshal(requestBodyBytes, &request)
	if err != nil {
		return nil, httperrors.BadRequestError(err.Error())
	}

	if request.Link == "" {
		return nil, httperrors.BadRequestError("link is missing")
	}

	product, err := sc.scrapHandler.ScrapProduct(request.Link)
	if err != nil {
		return nil, httperrors.InternalServerError(err.Error())
	}

	return &web.JSONResponse{
		"link": request.Link,
		"product": map[string]interface{}{
			"name":         product.Name,
			"imageURL":     product.ImageURL,
			"description":  product.Description,
			"price":        product.Price,
			"totalReviews": product.TotalReviews,
		},
	}, nil
}
