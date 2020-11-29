package models

// AmazonProduct holds the single product info
type AmazonProduct struct {
	Name         string `json:"name"`
	ImageURL     string `json:"imageURL"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	TotalReviews int    `json:"totalReviews"`
}

type AmazonProductRequest struct {
	Link    string        `json:"link"`
	Product AmazonProduct `json:"product"`
}
