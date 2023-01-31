package models

import (
	"fmt"
	"gorm.io/gorm"
)

type StoreModel struct {
	StoreId     int64        `json:"storeId" gorm:"primaryKey;autoincrement:false"`
	ProductId   int64        `json:"productId" gorm:"primaryKey;autoincrement:false"`
	Product     ProductModel `gorm:"foreignKey:ProductId;references:ID"`
	IsAvailable bool         `json:"isAvailable"`
}

func (s *StoreModel) GetProductsInStore(db *gorm.DB, limit, start int) []ProductModel {
	var productsInStore []StoreModel
	var products []ProductModel
	result := db.Model(&StoreModel{}).Where("store_id = ?", s.StoreId).Limit(limit).Offset(start).Find(&productsInStore)

	if result.Error != nil {
		fmt.Println("Some error occurred")
		return nil
	}
	tx := db.Begin()
	for i := 0; i < len(productsInStore); i++ {
		p := ProductModel{ID: uint(productsInStore[i].ProductId)}
		result := db.First(&p)
		if result.Error != nil {
			tx.Rollback()
			break
		}
		products = append(products, p)

	}
	if tx.Error != nil {
		return nil
	}
	tx.Commit()
	return products
}

// todo AddProducts
func (s *StoreModel) AddProducts() bool {

	return false
}

func (s *StoreModel) TableName() string {
	return "stores"
}
