package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `gorm:"created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updatedAt"`
	DeletedAt time.Time `gorm:"deleted_at" json:"deletedAt"`
	CreateBy  string    `gorm:"create_by" json:"createBy"`
	UpdateBy  string    `gorm:"update_by" json:"updateBy"`
	ID        int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
}
