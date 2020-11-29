package repository

import (
	"amazon-scrapper/saver/config/db"
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID          string    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Link        string    `json:"link" gorm:"not null;column:link;default:null"`
	ProductJSON string    `json:"product" gorm:"not null;column:product;default:null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRepositoryInterface interface {
	CreateProduct(*Product) error
}

type ProductRepository struct {
	Db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{Db: db.Get()}
}

func (productRepo *ProductRepository) CreateProduct(product *Product) error {
	err := productRepo.Db.Create(&product).Error
	if err != nil {
		return err
	}

	return nil
}
