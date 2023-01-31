package services

import "gorm.io/gorm"

type Products struct {
	conn *gorm.DB
}

type IProducts interface {
	GetProduct(id uint) int
}

func NewProduct(conn *gorm.DB) *Products {
	return &Products{conn}
}

func (p *Products) GetProduct(id uint) int {
	return 0
}
