package models

import (
	"fmt"
	"github.com/abhikeshri07/go-mux/constants"
	"gorm.io/gorm"
)

type StoreModel struct {
	StoreId     int64        `json:"storeId" gorm:"primaryKey;autoincrement:false"`
	ProductId   int64        `json:"productId" gorm:"primaryKey;autoincrement:false"`
	Product     ProductModel `gorm:"foreignKey:ProductId;references:ID"`
	IsAvailable bool         `json:"isAvailable"`
}

func (s *StoreModel) CheckStoreId(db *gorm.DB) string {
	var store StoreModel
	result := db.Model(&StoreModel{}).Where("store_id = ?", s.StoreId).First(&store)
	if result.RowsAffected == 0 {
		return constants.STORE_NOT_FOUND_ERROR
	}
	return constants.STORE_FOUND_SUCCESS

}
func (s *StoreModel) GetProductsInStore(db *gorm.DB, limit, start int) []ProductModel {
	var productIds []int64
	var products []ProductModel
	tx := db.Begin()
	db.Model(&StoreModel{}).Where("store_id = ?", s.StoreId).Limit(limit).Offset(start).Pluck("product_id", &productIds)
	fmt.Println(productIds)
	if productIds == nil {

		return nil
	}
	db.Model(&ProductModel{}).Where("id IN ?", productIds).Find(&products)

	tx.Commit()
	return products
}

func (s *StoreModel) AddProducts(db *gorm.DB, productIds []int64) string {

	tx := db.Begin()
	for i := 0; i < len(productIds); i++ {

		var product ProductModel
		results := db.Model(&ProductModel{}).Where("id = ?", productIds[i]).First(&product)
		if results.RowsAffected == 0 {
			return constants.STORE_PRODUCT_ENTRY_FOREIGN_KEY_ERROR
		}
		s.ProductId = productIds[i]
		s.IsAvailable = true
		db.Create(&s)
	}
	if tx.Error != nil {
		return constants.DB_TRANSACTION_ERROR
	}
	err := tx.Commit().Error
	if err == nil {
		return constants.STORE_PRODCUT_ENTRY_SUCCESS
	}
	return constants.STORE_PRODCUT_ENTRY_ERROR

}

func (s *StoreModel) TableName() string {
	return "stores"
}
