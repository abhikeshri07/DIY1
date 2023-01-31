package models

import (
	"gorm.io/gorm"
)

type ProductModel struct {
	ID    uint    `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (*ProductModel) TableName() string {
	return "products"
}

func (p *ProductModel) GetProduct(db *gorm.DB) *gorm.DB {
	result := db.First(&p)

	return result
}

func (p *ProductModel) GetProducts(db *gorm.DB, limit, start int) ([]ProductModel, *gorm.DB) {
	var products []ProductModel
	result := db.Model(ProductModel{}).Offset(start).Limit(limit).Find(&products)

	return products, result
}

func (p *ProductModel) CreateProduct(db *gorm.DB) *gorm.DB {
	result := db.Create(&p)

	return result
}

func (p *ProductModel) UpdateProduct(db *gorm.DB, newProduct *ProductModel) *gorm.DB {
	result := db.Model(&p).Updates(newProduct)

	return result
}

func (p *ProductModel) DeleteProduct(db *gorm.DB) *gorm.DB {
	result := db.Delete(&p)

	return result
}
