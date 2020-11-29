package handlers

import (
	"amazon-scrapper/lib/models"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ScrapHandler struct {
}

func NewScrapHandler() *ScrapHandler {
	return &ScrapHandler{}
}

type ScrapHandlerInterface interface {
	ScrapProduct(url string) (*models.AmazonProduct, error)
}

func (sh *ScrapHandler) ScrapProduct(url string) (*models.AmazonProduct, error) {
	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
	)

	product := &models.AmazonProduct{}

	c.OnHTML("#titleSection", func(element *colly.HTMLElement) {
		product.Name = strings.Trim(element.ChildText("#productTitle"), " ")
	})

	c.OnHTML("#landingImage", func(element *colly.HTMLElement) {
		if imgURL := element.Attr("data-old-hires"); imgURL != "" {
			product.ImageURL = imgURL
			return
		}

		data := element.Attr("data-a-dynamic-image")
		var imgData map[string]interface{}
		json.Unmarshal([]byte(data), &imgData)
		for key := range imgData {
			product.ImageURL = key
			return
		}
	})

	c.OnHTML("#feature-bullets", func(element *colly.HTMLElement) {
		product.Description = strings.ReplaceAll(element.ChildText("#feature-bullets .a-unordered-list > li .a-list-item"), "\n", "")
	})

	c.OnHTML(".priceBlockBuyingPriceString", func(element *colly.HTMLElement) {
		product.Price = element.Text
	})

	if product.Price == "" {
		c.OnHTML("#edition_0_price", func(element *colly.HTMLElement) {
			str := strings.ReplaceAll(element.Text, "\n", "")
			strArr := strings.Split(str, " ")
			product.Price = strArr[len(strArr)-2]
		})
	}

	c.OnHTML("#acrCustomerReviewText", func(element *colly.HTMLElement) {
		str := regexp.MustCompile(`\d[\d,]*`).FindString(element.Text)
		totalReviews, _ := strconv.Atoi(strings.ReplaceAll(str, ",", ""))
		product.TotalReviews = totalReviews
	})

	c.Visit(url)
	c.Wait()
	return product, nil
}
