package model

import "time"

type UserSubscribedProduct struct {
	ID          int64     `gorm:"column:id" json:"id"`
	UserID      int64     `gorm:"column:user_id" json:"user_id"`
	ExpiredAt   time.Time `gorm:"column:expired_at" json:"expired_at"`
	IsActive    bool      `gorm:"column:is_active" json:"is_active"`
	ProductCode string    `gorm:"column:product_code" json:"product_code"`
	ProductName string    `gorm:"column:product_name" json:"product_name"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (UserSubscribedProduct) TableName() string {
	return "user_subscribed_products"
}
