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
	var productIds []int64
	var products []ProductModel
	tx := db.Begin()
	db.Model(&StoreModel{}).Where("store_id = ?", s.StoreId).Limit(limit).Offset(start).Pluck("product_id", &productIds)
	fmt.Println(productIds)
	if productIds == nil {
		fmt.Println("No Products found in the store")
		return nil
	}

	db.Model(&ProductModel{}).Where("id IN ?", productIds).Find(&products)
	//fmt.Println(products)
	tx.Commit()
	return products
}

// todo bulk insert
func (s *StoreModel) AddProducts(db *gorm.DB, products []ProductModel) bool {

	tx := db.Begin()
	for i := 0; i < len(products); i++ {
		res := db.Create(&products[i])
		fmt.Println(products[i].ID, products[i].Name)
		if res.Error != nil {
			tx.Rollback()
			break
		}

		s.ProductId = int64(products[i].ID)
		s.IsAvailable = true
		result := db.Model(&s).Where("store_id = ? and product_id = ?", s.StoreId, s.ProductId).Updates(&s)
		if result.RowsAffected == 0 {
			result = db.Create(&s)
		}

	}
	if tx.Error != nil {
		return false
	}
	err := tx.Commit().Error
	return err == nil
}

func (s *StoreModel) TableName() string {
	return "stores"
}
