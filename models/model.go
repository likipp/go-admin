package models

import "time"

type BaseModel struct {
	CreatedAt *time.Time `gorm:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"updated_at" json:"updatedAt"`
	CreateBy  string     `gorm:"create_by"  json:"createBy"`
	UpdateBy  string     `gorm:"update_by"  json:"updateBy"`
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
}
