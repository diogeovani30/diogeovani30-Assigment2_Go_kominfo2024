package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Items struct {
	ID          string `json:"id" gorm:"primary_key"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

func (i *Items) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.NewString()
	return nil
}
