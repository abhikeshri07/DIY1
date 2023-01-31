package models

type StoreModel struct {
	StoreId     int64        `json:"storeId" gorm:"primaryKey;autoincrement:false"`
	ProductId   int64        `json:"productId" gorm:"primaryKey;autoincrement:false"`
	Product     ProductModel `gorm:"foreignKey:ProductId;references:ID"`
	IsAvailable bool         `json:"isAvailable"`
}
