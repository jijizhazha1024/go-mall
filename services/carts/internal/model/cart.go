package model

import (
	"gorm.io/gorm"
	"time"
)

type Cart struct {
	ID        int32          `gorm:"primaryKey;autoIncrement;comment:主键 自增"`
	CreatedAt time.Time      `gorm:"comment:创建时间"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
	UserID    int32          `gorm:"not null;index;comment:用户ID"`
	ProductID int32          `gorm:"not null;index;comment:商品ID"`
	Quantity  int32          `gorm:"default:0;comment:商品数量"`
	Checked   bool           `gorm:"default:false;comment:商品是否选中"`

	User    User    `gorm:"foreignKey:UserID;references:ID"`
	Product Product `gorm:"foreignKey:ProductID;references:ID"`
}

type User struct {
	ID int32 `gorm:"primaryKey;autoIncrement;comment:用户ID"`
}

type Product struct {
	ID int32 `gorm:"primaryKey;autoIncrement;comment:商品ID"`
}
