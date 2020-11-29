package controllers

import (
	"amazon-scrapper/lib/httperrors"
	"amazon-scrapper/lib/models"
	"amazon-scrapper/lib/web"
	"amazon-scrapper/saver/app/repository"
	"encoding/json"
	"net/http"
)

type SaveController struct {
	productRepo repository.ProductRepositoryInterface
}

func NewSaveController() *SaveController {
	return &SaveController{
		productRepo: repository.NewProductRepository(),
	}
}

type SaveControllerInterface interface {
	SaveProduct(*http.Request) (*web.JSONResponse, *httperrors.HttpError)
}

func (sc *SaveController) SaveProduct(r *http.Request) (*web.JSONResponse, *httperrors.HttpError) {
	requestBodyBytes, err := web.ReadBodyBytes(r)
	if err != nil {
		return nil, httperrors.BadRequestError(err.Error())
	}

	request := models.AmazonProductRequest{}
	err = json.Unmarshal(requestBodyBytes, &request)
	if err != nil {
		return nil, httperrors.BadRequestError(err.Error())
	}

	productJSON, err := json.Marshal(request.Product)
	if err != nil {
		return nil, httperrors.InternalServerError(err.Error())
	}

	product := &repository.Product{
		Link:        request.Link,
		ProductJSON: string(productJSON),
	}
	err = sc.productRepo.CreateProduct(product)
	if err != nil {
		return nil, httperrors.InternalServerError(err.Error())
	}

	return &web.JSONResponse{}, nil
}
